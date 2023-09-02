package app

import (
	"context"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	storage Storage
}

/* change to package use
type Logger interface {
	Info()
	Debug()
	Error()
	Fatal()
}
*/

type Storage interface { // TODO
	Create()
	AddEvent(ctx context.Context, event storage.Event) error
	GetSortedEventsById(id string) storage.Event
}

func New(storage Storage) *App {
	logger.Info("app start")
	return &App{storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return a.storage.AddEvent(ctx, storage.Event{ID: id, Title: title})
}

// TODO
