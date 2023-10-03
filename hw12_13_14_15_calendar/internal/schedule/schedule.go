package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/notify"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	storage  storage.Storage
	notifier notify.Notifier
	wg       sync.WaitGroup
}

func NewScheduler(storage storage.Storage, notifier notify.Notifier) *Scheduler {
	return &Scheduler{
		storage:  storage,
		notifier: notifier,
	}
}

func (s *Scheduler) LoadSchedule(ctx context.Context) {
	now := time.Now()
	deleted, err := s.storage.DeleteEventsBefore(ctx, now.AddDate(-1, 0, 0))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete past events: %v", err))
		return
	}
	if deleted > 0 {
		logger.Info(fmt.Sprintf("deleted %d past events", deleted))
	}

	events, err := s.storage.GetEventsToNotify(ctx, now, now.Add(config.SchedulerSettings.Schedule))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get events: %v", err))
		return
	}
	s.scheduleNotification(ctx, events)
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func (s *Scheduler) scheduleNotification(ctx context.Context, events []models.Event) {
	s.wg.Add(len(events))
	for _, event := range events {
		logger.Debug(fmt.Sprintf("event %v will to notify", event))
		go func(event models.Event) {
			defer s.wg.Done()
			select {
			case <-ctx.Done():
				logger.Warn(fmt.Sprintf("event %s was scheduled, but no notified", event.ID))
				return
			case <-time.After(time.Until(event.DateNotification)):
				s.notifyAboutEvent(ctx, &event)
				return
			}
		}(event)
	}
}

func (s *Scheduler) notifyAboutEvent(ctx context.Context, event *models.Event) {
	notification := models.Notification{
		ID:        uuid.New(),
		EventID:   event.ID,
		Title:     event.Title,
		EventDate: event.DateStart,
		UserID:    event.UserID,
	}
	if err := s.notifier.Notify(ctx, &notification); err != nil {
		logger.Error(fmt.Sprintf("failed to notify about event %s", event.ID))
	}
}
