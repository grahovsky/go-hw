package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type (
	id = string

	Events map[id]storage.Event
	Dates  map[time.Time]map[id]struct{}

	Storage struct {
		mu     sync.RWMutex //nolint:unused
		events Events
	}
)

func New() *Storage {
	return &Storage{
		events: make(Events),
	}
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return nil
}

func (store *Storage) DeleteEvent(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if id == "" {
		logger.Error(storage.ErrEventID.Error())
		return storage.ErrEventID
	}

	if _, ok := store.events[id]; ok {
		delete(store.events, id)
	}

	logger.Info("Delete event with ID: " + string(id))

	return nil
}

func (s *Storage) GetSortedEventsById(id string) storage.Event {
	if id == "" {
		return storage.Event{}
	}

	return s.events[id]
}
