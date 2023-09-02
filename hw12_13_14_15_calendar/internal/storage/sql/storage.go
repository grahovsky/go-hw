package sqlstorage

import (
	"context"

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

func (s *Storage) Create() {
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	return nil
}

func (s *Storage) GetSortedEventsById(id string) storage.Event {
	return storage.Event{}
}
