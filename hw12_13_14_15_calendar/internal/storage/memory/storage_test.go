package memorystorage

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testEventsCount   = 3
	testEventDuration = time.Minute
	testEventGap      = 10 * time.Second
)

type MemoryStorageTestSuite struct {
	suite.Suite
	events  []models.Event
	storage *Storage
	start   time.Time
}

func (s *MemoryStorageTestSuite) SetupTest() {
	s.start = time.Now()
	s.events = make([]models.Event, testEventsCount)
	t := s.start
	for i := 0; i < testEventsCount; i++ {
		s.events[i] = makeTestEvent(fmt.Sprintf("Event_%d", i), t)
		t = t.Add(testEventDuration).Add(testEventGap)
	}
	s.storage = &Storage{events: makeEventMap(s.events)}
}

func (s *MemoryStorageTestSuite) TestListEvents() {
	ctx := context.Background()
	s.Run("simple case", func() {
		events, err := s.storage.ListEvents(ctx, 2, 0)
		s.NoError(err)
		s.Equal(s.events[:2], events)

		events, err = s.storage.ListEvents(ctx, 2, 1)
		s.NoError(err)
		s.Equal(s.events[1:3], events)

		events, err = s.storage.ListEvents(ctx, 2, 2)
		s.NoError(err)
		s.Equal(s.events[2:3], events)

		events, err = s.storage.ListEvents(ctx, 2, 3)
		s.NoError(err)
		s.Empty(events)
	})
	s.Run("get all", func() {
		events, err := s.storage.ListEvents(ctx, math.MaxUint64, 0)
		s.NoError(err)
		s.Equal(s.events, events)
	})
	s.Run("empty storage", func() {
		storage := Storage{}
		events, err := storage.ListEvents(ctx, 100, 0)
		s.NoError(err)
		s.Empty(events)
	})
}

func (s *MemoryStorageTestSuite) TestGetEventsForPeriod() {
	ctx := context.Background()
	s.Run("simple case", func() {
		s.Run("all events", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start, s.start.Add(150*time.Second))
			s.NoError(err)
			s.Equal(s.events, events)
		})
		s.Run("1st and 2nd events", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start, s.start.Add(120*time.Second))
			s.NoError(err)
			s.Equal(s.events[:2], events)
		})
		s.Run("2st and 3rd events", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start.Add(100*time.Second), s.start.Add(150*time.Second))
			s.NoError(err)
			s.Equal(s.events[1:], events)
		})
		s.Run("2st event", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start.Add(100*time.Second), s.start.Add(70*time.Second))
			s.NoError(err)
			s.Equal(s.events[1:1], events)
		})
	})
	s.Run("no events", func() {
		s.Run("before first", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start.Add(-time.Second), s.start)
			s.NoError(err)
			s.Empty(events)
		})
		s.Run("after last", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start.Add(211*time.Second), s.start.Add(time.Hour))
			s.NoError(err)
			if len(events) > 0 {
				println(len(events))
			}
			s.Empty(events)
		})
		s.Run("in a gap", func() {
			events, err := s.storage.GetEventsForPeriod(ctx, s.start.Add(65*time.Second), s.start.Add(70*time.Second))
			s.NoError(err)
			s.Empty(events)
		})
	})
}

func (s *MemoryStorageTestSuite) TestGetEvent() {
	ctx := context.Background()
	s.Run("simple case", func() {
		for i := 0; i < testEventsCount; i++ {
			eventID := s.events[i].ID
			event, err := s.storage.GetEvent(ctx, eventID)
			s.NoError(err)
			s.Equal(&s.events[i], event)
		}
	})
	s.Run("non-existent event", func() {
		event, err := s.storage.GetEvent(ctx, uuid.New())
		s.ErrorIs(err, models.ErrEventNotFound)
		s.Nil(event)
	})
	s.Run("err event id", func() {
		event, err := s.storage.GetEvent(ctx, uuid.Nil)
		s.ErrorIs(err, models.ErrEventID)
		s.Nil(event)
	})
}

func (s *MemoryStorageTestSuite) TestInsertEvent() {
	ctx := context.Background()
	s.Run("simple case", func() {
		testEvent := makeTestEvent("test_event", time.Now())
		err := s.storage.AddEvent(ctx, &testEvent)
		s.NoError(err)

		event, err := s.storage.GetEvent(ctx, testEvent.ID)
		s.NoError(err)
		s.Equal(&testEvent, event)
	})
	s.Run("already existent event", func() {
		for i := 0; i < testEventsCount; i++ {
			eventID := s.events[i].ID
			changedEvent := makeTestEvent("changed_event", time.Now())
			changedEvent.ID = eventID
			err := s.storage.AddEvent(ctx, &changedEvent)
			s.ErrorIs(err, models.ErrEventAlreadyExists)
		}
	})
}

func (s *MemoryStorageTestSuite) TestUpdateEvent() {
	ctx := context.Background()
	s.Run("simple case", func() {
		eventID := s.events[0].ID
		changedEvent := makeTestEvent("changed_event", time.Now())
		changedEvent.ID = eventID
		err := s.storage.UpdateEvent(ctx, &changedEvent)
		s.NoError(err)

		event, err := s.storage.GetEvent(ctx, eventID)
		s.NoError(err)
		s.Equal(&changedEvent, event)
	})
	s.Run("non-existent event", func() {
		testEvent := makeTestEvent("test_event", time.Now())
		err := s.storage.UpdateEvent(ctx, &testEvent)
		s.ErrorIs(err, models.ErrEventNotFound)
	})
}

func (s *MemoryStorageTestSuite) TestDeleteEvent() {
	ctx := context.Background()
	s.Run("simple case", func() {
		for i := 0; i < testEventsCount; i++ {
			eventID := s.events[i].ID
			err := s.storage.DeleteEvent(ctx, eventID)
			s.NoError(err)

			event, err := s.storage.GetEvent(ctx, eventID)
			s.ErrorIs(err, models.ErrEventNotFound)
			s.Nil(event)
		}
	})
	s.Run("non-existent event", func() {
		err := s.storage.DeleteEvent(ctx, uuid.New())
		s.ErrorIs(err, models.ErrEventNotFound)
	})
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(MemoryStorageTestSuite))
}

func TestStorage_Concurrency(t *testing.T) {
	ctx := context.Background()
	s := &Storage{}
	s.InitStorage(&config.Storage{})
	count := 100

	var wg sync.WaitGroup
	wg.Add(2 * count)
	for i := 0; i < count; i++ {
		eventID := uuid.New()

		// NOTE: insert
		go func(i int) {
			defer wg.Done()
			event := makeTestEvent(fmt.Sprintf("Event_%d", i), time.Now().Add(time.Duration(i)*time.Second))
			event.ID = eventID
			require.NoError(t, s.AddEvent(ctx, &event))
		}(i)

		// NOTE: get
		go func() {
			defer wg.Done()
			require.Eventually(t, func() bool {
				event, err := s.GetEvent(ctx, eventID)
				return err == nil && event.ID == eventID
			}, time.Minute, time.Millisecond)
		}()
	}
}

func makeTestEvent(title string, dateStart time.Time) models.Event {
	return models.Event{
		ID:               uuid.New(),
		Title:            title,
		DateStart:        dateStart,
		DateEnd:          dateStart.Add(testEventDuration),
		Description:      "some description",
		UserID:           uuid.New(),
		DateNotification: dateStart.Add(-time.Duration(rand.Intn(10)) * time.Second),
	}
}

func makeEventMap(events []models.Event) map[uuid.UUID]models.Event {
	m := make(map[uuid.UUID]models.Event, len(events))
	for _, event := range events {
		m[event.ID] = event
	}
	return m
}
