package rmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

type Notifier struct {
	queue *Queue
}

func (n *Notifier) Notify(ctx context.Context, notification *models.Notification) error {
	msg, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal to json: %w", err)
	}
	logger.Debug(fmt.Sprintf("notify push to %s %s", n.queue.queue.Name, n.queue.contentType))
	return n.queue.Push(ctx, msg, "application/json")
}

func (n *Notifier) Close() error {
	return n.queue.Close()
}

func NewNotifier(rmqCf *config.RMQ) (*Notifier, error) {
	queue, err := NewQueue(rmqCf)
	if err != nil {
		return nil, fmt.Errorf("create queue: %w", err)
	}
	return &Notifier{queue: queue}, nil
}
