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

type Storage interface { // TODO
	InitStorage(ctx context.Context)
	AddEvent(ctx context.Context, event *storage.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]storage.Event, error)
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
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
