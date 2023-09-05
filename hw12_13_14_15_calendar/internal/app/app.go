package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
}

type Storage interface {
	InitStorage(ctx context.Context)
	AddEvent(ctx context.Context, event *storage.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error)
	GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]storage.Event, error)
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	Close(ctx context.Context) error
}

func New(storage Storage) *App {
	logger.Info("app start")
	return &App{storage: storage}
}

func (a *App) AddEvent(ctx context.Context, event *storage.Event) error {
	// TODO
	return a.storage.AddEvent(ctx, event)
}

// TODO
