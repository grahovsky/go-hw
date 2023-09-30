package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	InitStorage() error
	AddEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]models.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	Close() error
}

func New(sType string) (Storage, error) {
	var st Storage
	if sType == "sql" {
		st = &sqlstorage.Storage{}
	} else {
		st = &memorystorage.Storage{}
	}
	err := st.InitStorage()

	return st, err
}
