// go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"net"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api/eventservice"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/rmq"
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
	rmqCfg         *config.RMQ
}

func TestMain(m *testing.M) {
	// config.InitCalendarSettings()
	// not working relative default config path
	// to do - load test settings from env
	os.Exit(m.Run())
}

func (s *calendarTestSuite) SetupSuite() {
	calendarConn, err := grpc.Dial(net.JoinHostPort("calendar", "8082"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.calendarConn = calendarConn
	s.calendarClient = eventservice.NewCalendarClient(s.calendarConn)

	// s.ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	s.ctx = context.Background()

	dsn := "postgres://postgres:postgres@postgres:5432/calendar"
	// dsn := sqlstorage.GetDsn(&config.CalendarSettings.Storage)
	db, err := sqlx.Connect("pgx", dsn)
	s.Require().NoError(err)
	s.db = db

	s.rmqCfg = &config.RMQ{
		Host:     "rabbit",
		Port:     "5672",
		User:     "guest",
		Password: "guest",
		Queue:    "notifications_out",
	}
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
	s.Run("standard case", func() {
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
	s.Run("standard case", func() {
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
		eventID := uuid.New().String()
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
		s.ErrorContains(err, "event not found")
	})
}

func (s *calendarTestSuite) TestDeleteEvent() {
	eventID := s.addOneEvent()
	s.Run("standard case", func() {
		_, err := s.calendarClient.DeleteEvent(s.ctx, &eventservice.DeleteEventRequest{EventId: eventID})
		s.NoError(err)

		_, err = s.calendarClient.GetEvent(s.ctx, &eventservice.GetEventRequest{EventId: eventID})
		s.ErrorContains(err, "event not found")
	})
	s.Run("non existent event", func() {
		_, err := s.calendarClient.DeleteEvent(s.ctx, &eventservice.DeleteEventRequest{EventId: uuid.New().String()})
		s.ErrorContains(err, "event not found")
	})
}

func (s *calendarTestSuite) TestSendNotifications() {
	queue, err := rmq.NewQueue(s.rmqCfg)
	s.NoError(err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := queue.ConsumeChannel(ctx, "calendar_sender_test")
	s.NoError(err)

	count := 5
	since := time.Now().Add(5 * time.Second)
	events := s.addFewEvents(since, 500*time.Millisecond, count)
	notifications := make([]models.Notification, 0, count)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer cancel()
		defer wg.Done()
		s.Eventually(func() bool {
			return len(notifications) == count
		}, time.Minute, time.Millisecond)
	}()
	go func() {
		defer wg.Done()
		for msg := range ch {
			var notification models.Notification
			s.NoError(json.Unmarshal(msg, &notification))

			if contains(events, notification.EventID.String()) {
				notifications = append(notifications, notification)
			}
		}
	}()
	<-ctx.Done()
	wg.Wait()

	s.Equal(len(events), len(notifications))
	sort.Strings(events)
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].EventID.String() < notifications[j].EventID.String()
	})
	for i := 0; i < len(events); i++ {
		s.Equal(events[i], notifications[i].EventID.String())
	}
}

func (s *calendarTestSuite) TestGetEventsOfDay() {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	s.addFewEvents(yesterday, 4*time.Hour, 18)

	s.Run("in period day", func() {
		for _, date := range []time.Time{yesterday, today, tomorrow} {
			res, err := s.calendarClient.GetEventsOfDay(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("no in period day", func() {
		for _, date := range []time.Time{yesterday.AddDate(0, 0, -5), tomorrow.AddDate(1, 0, 5)} {
			res, err := s.calendarClient.GetEventsOfDay(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
	})
}

func (s *calendarTestSuite) TestGetEventsOfWeek() {
	s.addFewEvents(time.Now(), 24*time.Hour, 10)

	s.Run("in period week", func() {
		for _, date := range []time.Time{time.Now(), time.Now().AddDate(0, 0, 7)} {
			res, err := s.calendarClient.GetEventsOfWeek(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("no in period week", func() {
		for _, date := range []time.Time{time.Now().AddDate(0, 0, -8), time.Now().AddDate(0, 0, 15)} {
			res, err := s.calendarClient.GetEventsOfWeek(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
	})
}

func (s *calendarTestSuite) TestGetEventsOfMonth() {
	s.addFewEvents(time.Now(), 24*time.Hour, 35)

	s.Run("in period month", func() {
		for _, date := range []time.Time{time.Now(), time.Now().AddDate(0, 1, 0)} {
			res, err := s.calendarClient.GetEventsOfMonth(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("no in period month", func() {
		for _, date := range []time.Time{time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 2, 0)} {
			res, err := s.calendarClient.GetEventsOfMonth(s.ctx, &eventservice.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
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
	eventsID := make([]string, 0)

	t := since
	for i := 0; i < count; i++ {
		event, err := s.calendarClient.AddEvent(s.ctx, &eventservice.AddEventRequest{
			DateStart:        timestamppb.New(t),
			DateEnd:          timestamppb.New(t.Add(d)),
			UserId:           uuid.New().String(),
			DateNotification: timestamppb.New(t),
		})
		s.NoError(err)
		t = t.Add(d)
		eventsID = append(eventsID, event.EventId)
	}

	return eventsID
}

func contains[T comparable](collection []T, value T) bool {
	for _, v := range collection {
		if value == v {
			return true
		}
	}
	return false
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(calendarTestSuite))
}
