package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage Storage
}

/* change to package use
type Logger interface {
	Error()
	Warn()
	Info()
	Debug()
}
*/

type Storage interface { // TODO
	InitStorage()
	AddEvent(ctx context.Context, event *storage.Event) error
	// UpdateEvent(ctx context.Context, event *storage.Event) error
	// DeleteEvent(context.Context, string) error
	// GetEventByDay(context.Context, int64, time.Time) ([]storage.Event, error)
	// GetEventByWeek(context.Context, int64, time.Time) ([]storage.Event, error)
	// GetEventByMonth(context.Context, int64, time.Time)

	GetEvent(id uuid.UUID) (*storage.Event, error)
}

func New(storage Storage) *App {
	logger.Info("app start")
	return &App{storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	// TODO
	return a.storage.AddEvent(ctx, event)
}

// TODO
