package sqlstorage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	db string
}

func New() *Storage {
	return &Storage{
		db: "some",
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) InitStorage(ctx context.Context) {
}

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) error {
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	return &storage.Event{}, nil
}

func (s *Storage) GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) ListEvents(ctx context.Context, limit, low uint64) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return nil
}
