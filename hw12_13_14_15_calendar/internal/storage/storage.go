package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	memorystorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	InitStorage(settings *config.Storage) error
	AddEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error)
	GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]models.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	Close() error
	DeleteEventsBefore(ctx context.Context, before time.Time) (int64, error)
	GetEventsToNotify(ctx context.Context, from, to time.Time) ([]models.Event, error)
}

func New(settings *config.Storage) (Storage, error) {
	var st Storage

	if settings.Type == "sql" {
		st = &sqlstorage.Storage{}
	} else {
		st = &memorystorage.Storage{}
	}
	err := st.InitStorage(settings)

	return st, err
}
