package memorystorage

import (
	"context"
	"sync"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/google/uuid"
)

type (
	Events map[uuid.UUID]storage.Event

	Storage struct {
		mu     sync.RWMutex //nolint:unused
		events Events
	}
)

func (s *Storage) InitStorage() {
	s.events = make(Events)
}

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.contains(event.ID) {
		return storage.ErrEventAlreadyExists
	}

	s.events[event.ID] = *event
	return nil
}

func (store *Storage) DeleteEvent(id uuid.UUID) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if id == uuid.Nil {
		logger.Error(storage.ErrEventID.Error())
		return storage.ErrEventID
	}

	if _, ok := store.events[id]; ok {
		delete(store.events, id)
	}

	logger.Info("Delete event with ID: " + id.String())

	return nil
}

func (s *Storage) GetEventsById(id uuid.UUID) storage.Event {
	if id == uuid.Nil {
		return storage.Event{}
	}

	return s.events[id]
}

func (s *Storage) contains(id uuid.UUID) bool {
	_, ok := s.events[id]
	return ok
}
