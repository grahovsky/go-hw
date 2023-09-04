package memorystorage

import (
	"context"
	"sort"
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

func (s *Storage) InitStorage(_ context.Context) {
	s.events = make(Events)
}

func (s *Storage) AddEvent(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.contains(event.ID) {
		return storage.ErrEventAlreadyExists
	}

	s.events[event.ID] = *event
	return nil
}

func (s *Storage) GetEvent(_ context.Context, id uuid.UUID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id == uuid.Nil || !s.contains(id) {
		return &storage.Event{}, storage.ErrEventID
	}

	event, _ := s.events[id]

	return &event, nil
}

func (s *Storage) ListEvents(_ context.Context, limit, low uint64) ([]storage.Event, error) {
	events := make([]storage.Event, 0, len(s.events))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})

	high := min(low+limit, uint64(len(events)))
	return events[low:high], nil
}

func (s *Storage) UpdateEvent(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.contains(event.ID) {
		return storage.ErrEventNotFound
	}
	s.events[event.ID] = *event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == uuid.Nil || !s.contains(id) {
		return storage.ErrEventID
	}

	delete(s.events, id)
	logger.Info("Delete event with ID: " + id.String())

	return nil
}

func (s *Storage) contains(id uuid.UUID) bool {
	_, ok := s.events[id]
	return ok
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
