package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

type (
	Events map[uuid.UUID]models.Event

	Storage struct {
		mu     sync.RWMutex
		events Events
	}
)

func (s *Storage) InitStorage(_ config.Storage) error {
	s.events = make(Events)
	return nil
}

func (s *Storage) Close() error {
	for u := range s.events {
		delete(s.events, u)
	}
	return nil
}

func (s *Storage) AddEvent(_ context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.contains(event.ID) {
		return models.ErrEventAlreadyExists
	}

	s.events[event.ID] = *event
	return nil
}

func (s *Storage) GetEvent(_ context.Context, id uuid.UUID) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id == uuid.Nil {
		return nil, models.ErrEventID
	}

	if !s.contains(id) {
		return nil, models.ErrEventNotFound
	}

	event := s.events[id]

	return &event, nil
}

func (s *Storage) GetEventsForPeriod(_ context.Context, from, to time.Time) ([]models.Event, error) {
	events := make([]models.Event, 0, len(s.events))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.InPeriod(from, to) {
			events = append(events, event)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})

	return events, nil
}

func (s *Storage) ListEvents(_ context.Context, limit, low uint64) ([]models.Event, error) {
	events := make([]models.Event, 0, len(s.events))

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

func (s *Storage) UpdateEvent(_ context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.contains(event.ID) {
		return models.ErrEventNotFound
	}
	s.events[event.ID] = *event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == uuid.Nil {
		return models.ErrEventID
	}

	if !s.contains(id) {
		return models.ErrEventNotFound
	}

	delete(s.events, id)
	logger.Debug("Delete event with ID: " + id.String())

	return nil
}

func (s *Storage) DeleteEventsBefore(_ context.Context, before time.Time) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	deleted := int64(0)
	for _, event := range s.events {
		if event.DateEnd.Before(before) {
			deleted++
		}
	}
	return deleted, nil
}

func (s *Storage) GetEventsToNotify(_ context.Context, from, to time.Time) ([]models.Event, error) {
	events := make([]models.Event, 0, len(s.events))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.events {
		if event.IsToNotify(from, to) {
			events = append(events, event)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].DateStart.Before(events[j].DateStart)
	})

	return events, nil
}
