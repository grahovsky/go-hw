package notify

import (
	"context"
	"fmt"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

type Notifier interface {
	Notify(ctx context.Context, notification *models.Notification) error
}

type logNotifier struct{}

func NewLogNotifier() Notifier {
	return &logNotifier{}
}

func (n *logNotifier) Notify(_ context.Context, notification *models.Notification) error {
	logger.Info(fmt.Sprintf("[UPCOMING EVENT NOTIFICATION]: %v", *notification))
	return nil
}
