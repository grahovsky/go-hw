package api

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	pb "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/api/apppb"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufferSize = 1024 * 1024

type apiTestSuite struct {
	suite.Suite
	storage  storage.Storage
	srv      *grpc.Server
	listener *bufconn.Listener

	conn   *grpc.ClientConn
	client pb.AppClient
}

func (s *apiTestSuite) SetupSuite() {
	s.storage = &memorystorage.Storage{}
	s.storage.InitStorage()
	s.srv = grpc.NewServer()
	s.listener = bufconn.Listen(bufferSize)

	pb.RegisterAppServer(s.srv, NewAPI(app.New(s.storage)))
	go func() {
		s.NoError(s.srv.Serve(s.listener))
	}()
}

func (s *apiTestSuite) SetupTest() {
	conn, err := grpc.DialContext(
		context.Background(), "api_test",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoError(err)
	s.conn = conn
	s.client = pb.NewAppClient(conn)
}

func (s *apiTestSuite) TearDownTest() {
	s.NoError(s.conn.Close())
}

func (s *apiTestSuite) TearDownSuite() {
	s.srv.Stop()
}

func (s *apiTestSuite) TestAddEvent() {
	ctx := context.Background()
	s.Run("simple case add event", func() {
		res, err := s.client.AddEvent(ctx, &pb.AddEventRequest{
			Title:     "some testing event",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Hour)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)

		eventID, err := uuid.Parse(res.EventId)
		s.NoError(err)

		event, err := s.storage.GetEvent(ctx, eventID)
		s.NoError(err)
		s.Equal(eventID, event.ID)
		s.Equal("some testing event", event.Title)
	})
	s.Run("add unique events id", func() {
		res1, err := s.client.AddEvent(ctx, &pb.AddEventRequest{
			Title:     "some testing event",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Hour)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)
		res2, err := s.client.AddEvent(ctx, &pb.AddEventRequest{
			Title:     "event_2",
			DateStart: timestamppb.New(time.Now()),
			DateEnd:   timestamppb.New(time.Now().Add(time.Hour)),
			UserId:    uuid.New().String(),
		})
		s.NoError(err)
		s.NotEqual(res1.EventId, res2.EventId)
	})
	s.Run("invalid request", func() {
		_, err := s.client.AddEvent(ctx, &pb.AddEventRequest{Title: "some testing event"})
		s.ErrorContains(err, "invalid request")
	})
}

func (s *apiTestSuite) TestDeleteEvent() {
	ctx := context.Background()
	eventID := s.addTestEvent()
	s.Run("simple case", func() {
		_, err := s.client.DeleteEvent(ctx, &pb.DeleteEventRequest{EventId: eventID.String()})
		s.NoError(err)

		_, err = s.storage.GetEvent(ctx, eventID)
		s.ErrorContains(err, "Event not found")
	})
	s.Run("non existent event", func() {
		_, err := s.client.DeleteEvent(ctx, &pb.DeleteEventRequest{EventId: uuid.New().String()})
		s.ErrorContains(err, "Event not found")
	})
}

func (s *apiTestSuite) TestUpdateEvent() {
	ctx := context.Background()
	eventID := s.addTestEvent()
	s.Run("simple case", func() {
		_, err := s.client.UpdateEvent(ctx, &pb.UpdateEventRequest{
			Event: &pb.Event{
				Id:               eventID.String(),
				Title:            "some changing event",
				DateStart:        timestamppb.New(time.Now()),
				DateEnd:          timestamppb.New(time.Now().Add(time.Hour)),
				UserId:           uuid.New().String(),
				DateNotification: timestamppb.New(time.Now().Add(-time.Hour)),
			},
		})
		s.NoError(err)

		event, err := s.storage.GetEvent(ctx, eventID)
		s.NoError(err)
		s.Equal("some changing event", event.Title)
	})
	s.Run("non existent event", func() {
		_, err := s.client.UpdateEvent(ctx, &pb.UpdateEventRequest{
			Event: &pb.Event{
				Id:               uuid.New().String(),
				Title:            "some changing event",
				DateStart:        timestamppb.New(time.Now()),
				DateEnd:          timestamppb.New(time.Now().Add(time.Hour)),
				UserId:           uuid.New().String(),
				DateNotification: timestamppb.New(time.Now().Add(-time.Hour)),
			},
		})
		s.ErrorContains(err, "Event not found")
	})
	s.Run("invalid request", func() {
		_, err := s.client.UpdateEvent(ctx, &pb.UpdateEventRequest{
			Event: &pb.Event{Title: "some testing event"},
		})
		s.ErrorContains(err, "invalid request")
	})
}

func (s *apiTestSuite) TestGetEventsOfDay() {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	s.fillEventStorage(yesterday, 4*time.Hour, 18)

	ctx := context.Background()
	s.Run("in period day", func() {
		for _, date := range []time.Time{yesterday, today, tomorrow} {
			res, err := s.client.GetEventsOfDay(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("not in period day", func() {
		for _, date := range []time.Time{yesterday.AddDate(0, 0, -1), tomorrow.AddDate(1, 0, 5)} {
			res, err := s.client.GetEventsOfDay(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
	})
}

func (s *apiTestSuite) TestGetEventsOfWeek() {
	s.fillEventStorage(time.Now(), 24*time.Hour, 10)

	ctx := context.Background()
	s.Run("in period week", func() {
		for _, date := range []time.Time{time.Now(), time.Now().AddDate(0, 0, 7)} {
			res, err := s.client.GetEventsOfWeek(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("not in period week", func() {
		for _, date := range []time.Time{time.Now().AddDate(0, 0, -8), time.Now().AddDate(0, 0, 15)} {
			res, err := s.client.GetEventsOfWeek(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
	})
}

func (s *apiTestSuite) GetEventsOfMonth() {
	s.fillEventStorage(time.Now(), 24*time.Hour, 35)

	ctx := context.Background()
	s.Run("in period month", func() {
		for _, date := range []time.Time{time.Now(), time.Now().AddDate(0, 1, 0)} {
			res, err := s.client.GetEventsOfMonth(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.NotEmpty(res.Events)
		}
	})
	s.Run("not in period month", func() {
		for _, date := range []time.Time{time.Now().AddDate(0, -1, 0), time.Now().AddDate(0, 2, 0)} {
			res, err := s.client.GetEventsOfMonth(ctx, &pb.GetEventsRequest{Since: timestamppb.New(date)})
			s.NoError(err)
			s.Empty(res.Events)
		}
	})
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(apiTestSuite))
}

func (s *apiTestSuite) addTestEvent() uuid.UUID {
	eventID := uuid.New()
	event := models.Event{
		ID:               eventID,
		DateStart:        time.Now(),
		DateEnd:          time.Now().Add(time.Hour),
		UserID:           uuid.New(),
		DateNotification: time.Now().Add(-time.Hour),
	}
	s.NoError(s.storage.AddEvent(context.Background(), &event))
	return eventID
}

func (s *apiTestSuite) fillEventStorage(since time.Time, d time.Duration, count int) {
	t := since
	for i := 0; i < count; i++ {
		event := models.Event{
			ID:               uuid.New(),
			DateStart:        t,
			DateEnd:          t.Add(d),
			UserID:           uuid.New(),
			DateNotification: t,
		}
		s.NoError(s.storage.AddEvent(context.Background(), &event))
		t = t.Add(d)
	}
}
