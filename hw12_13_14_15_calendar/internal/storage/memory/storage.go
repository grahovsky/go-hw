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
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.contains(event.ID) {
		logger.Error(storage.ErrEventAlreadyExists.Error())
		return storage.ErrEventAlreadyExists
	}

	s.events[event.ID] = *event
	return nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == uuid.Nil || !s.contains(id) {
		logger.Error(storage.ErrEventID.Error())
		return storage.ErrEventID
	}

	delete(s.events, id)
	logger.Info("Delete event with ID: " + id.String())

	return nil
}

func (s *Storage) GetEvent(id uuid.UUID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id == uuid.Nil || !s.contains(id) {
		logger.Error(storage.ErrEventID.Error())
		return &storage.Event{}, storage.ErrEventID
	}

	event, _ := s.events[id]

	return &event, nil
}

func (s *Storage) contains(id uuid.UUID) bool {
	_, ok := s.events[id]
	return ok
}
