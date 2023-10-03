// go:build integration

package integration_test

import (
	"context"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api/eventservice"
	_ "github.com/jackc/pgx/stdlib" // justifying
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type calendarTestSuite struct {
	suite.Suite
	calendarConn   *grpc.ClientConn
	calendarClient eventservice.CalendarClient
	ctx            context.Context
	db             *sqlx.DB
}

func TestMain(m *testing.M) {
	// config.InitCalendarSettings()
	os.Exit(m.Run())
}

func (s *calendarTestSuite) SetupSuite() {
	// config.InitCalendarSettings()

	calendarConn, err := grpc.Dial(net.JoinHostPort("localhost", "8889"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.calendarConn = calendarConn
	s.calendarClient = eventservice.NewCalendarClient(s.calendarConn)

	// s.ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	s.ctx = context.Background()

	dsn := "postgres://postgres:postgres@localhost:5432/calendar"
	// dsn := sqlstorage.GetDsn(&config.CalendarSettings.Storage)
	db, err := sqlx.Connect("pgx", dsn)
	s.Require().NoError(err)
	s.db = db
}

func (s *calendarTestSuite) SetupTest() {
	seed := time.Now().UnixNano()
	rand.NewSource(seed)
	s.T().Log("seed:", seed)
}

func (s *calendarTestSuite) TearDownSuite() {
	s.NoError(s.calendarConn.Close())
	s.NoError(s.db.Close())
}

func (s *calendarTestSuite) TearDownTest() {
	_, err := s.db.ExecContext(s.ctx, `TRUNCATE TABLE events`)
	s.NoError(err)
	s.T().Logf("%s - done", s.T().Name())
}

func (s *calendarTestSuite) TestAddGetEvent() {
	s.Run("invalid request", func() {
		_, err := s.calendarClient.GetEvent(s.ctx, &eventservice.GetEventRequest{EventId: "not exist"})
		s.ErrorContains(err, "invalid eventID")
	})
	s.Run("simple case", func() {
		res, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
			Title:     "some testing event",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Second)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)

		eventID, err := uuid.Parse(res.EventId)
		s.NoError(err)

		created, err := s.calendarClient.GetEvent(s.ctx, &eventservice.GetEventRequest{EventId: eventID.String()})
		s.NoError(err)
		s.NotNil(created.Event)
		s.Equal(res.EventId, created.Event.Id)
		s.Equal("some testing event", created.Event.Title)
	})
	s.Run("create unique events id", func() {
		res1, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
			Title:     "some testing event",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Second)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)
		res2, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
			Title:     "some testing event 2",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Second)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)
		s.NotEqual(res1.EventId, res2.EventId)
	})
	s.Run("empty event", func() {
		_, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{Title: "some testing event"})
		s.ErrorContains(err, "invalid request")
	})
}

func (s *calendarTestSuite) TestUpdateEvent() {
	s.Run("invalid request", func() {
		_, err := s.calendarClient.UpdateEvent(s.ctx, &eventservice.UpdateEventRequest{
			Event: &eventservice.Event{Title: "some event "},
		})
		s.ErrorContains(err, "invalid UUID length")
	})
	s.Run("error update case", func() {
		_, err := s.calendarClient.UpdateEvent(s.ctx, &eventservice.UpdateEventRequest{
			Event: &eventservice.Event{
				Id:               "some id",
				Title:            "some updated event",
				DateStart:        timestamppb.New(time.Now()),
				DateEnd:          timestamppb.New(time.Now().Add(time.Second)),
				UserId:           uuid.New().String(),
				DateNotification: timestamppb.New(time.Now().Add(-10 * time.Second)),
			},
		})
		s.ErrorContains(err, "invalid UUID")
	})
	s.Run("simple case", func() {
		eventID := s.addOneEvent()
		_, err := s.calendarClient.UpdateEvent(s.ctx, &eventservice.UpdateEventRequest{
			Event: &eventservice.Event{
				Id:               eventID,
				Title:            "some updated event",
				DateStart:        timestamppb.New(time.Now()),
				DateEnd:          timestamppb.New(time.Now().Add(time.Second)),
				UserId:           uuid.New().String(),
				DateNotification: timestamppb.New(time.Now().Add(-10 * time.Second)),
			},
		})
		s.NoError(err)

		updated, err := s.calendarClient.GetEvent(s.ctx, &eventservice.GetEventRequest{EventId: eventID})
		s.NoError(err)
		s.NotNil(updated.Event)
		s.Equal("some updated event", updated.Event.Title)
	})
	s.Run("update not exists", func() {
		eventId := uuid.New().String()
		_, err := s.calendarClient.UpdateEvent(s.ctx, &eventservice.UpdateEventRequest{
			Event: &eventservice.Event{
				Id:               eventId,
				Title:            "some updated event",
				DateStart:        timestamppb.New(time.Now()),
				DateEnd:          timestamppb.New(time.Now().Add(time.Second)),
				UserId:           uuid.New().String(),
				DateNotification: timestamppb.New(time.Now().Add(-10 * time.Second)),
			},
		})
		s.ErrorContains(err, "event not found")
	})
}

func (s *calendarTestSuite) addOneEvent() string {
	event, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
		DateStart:        timestamppb.New(time.Now()),
		DateEnd:          timestamppb.New(time.Now().Add(time.Second)),
		UserId:           uuid.New().String(),
		DateNotification: timestamppb.New(time.Now().Add(-10 * time.Second)),
	})
	s.NoError(err)
	return event.EventId
}

func (s *calendarTestSuite) addFewEvents(since time.Time, d time.Duration, count int) []string {
	eventsId := make([]string, 0)

	t := since
	for i := 0; i < count; i++ {
		event, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
			DateStart:        timestamppb.New(time.Now()),
			DateEnd:          timestamppb.New(time.Now().Add(time.Second)),
			UserId:           uuid.New().String(),
			DateNotification: timestamppb.New(time.Now().Add(-10 * time.Second)),
		})
		s.NoError(err)
		t = t.Add(d)
		eventsId = append(eventsId, event.EventId)
	}

	return eventsId
}

func TestSystemStatusSuite(t *testing.T) {
	suite.Run(t, new(calendarTestSuite))
}
