package app

import (
	"context"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(storage Storage) *App {
	logger.Info("app start")
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
