package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage storage.Storage
}

func New(storage storage.Storage) *App {
	logger.Info("app start")
	return &App{storage: storage}
}

func (a *App) AddEvent(ctx context.Context, event *models.Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	return a.storage.AddEvent(ctx, event)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) GetEventsForPeriod(ctx context.Context, dateFrom, dateTo time.Time) ([]models.Event, error) {
	return a.storage.GetEventsForPeriod(ctx, dateFrom, dateTo)
}

func (a *App) ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error) {
	return a.storage.ListEvents(ctx, limit, low)
}

func (a *App) UpdateEvent(ctx context.Context, event *models.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, id)
}
