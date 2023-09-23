package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

//go:generate mockery --name Application

/* other way - mockgen
go:generate mockgen -destination=mocks/Application.go
-package=mocks github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server Application
*/

type Application interface {
	AddEvent(ctx context.Context, event *models.Event) (uuid.UUID, error)
	GetEvent(ctx context.Context, event uuid.UUID) (*models.Event, error)
	GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]models.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}
