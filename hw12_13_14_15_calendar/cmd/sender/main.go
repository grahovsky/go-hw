package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/rmq"
)

func main() {
	config.InitSenderSettings()
	logger.SetLogLevel(config.SenderSettings.Log.Level)

	queue, err := rmq.NewQueue(&config.SenderSettings.Rmq)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create RMQ queue: %v", err))
		os.Exit(1)
	}
	defer func() {
		if err := queue.Close(); err != nil {
			logger.Error(fmt.Sprintf("failed to close RMQ queue: %v", err))
		}
	}()

	// change to rmq notify
	cfgN := config.SenderSettings.Rmq
	cfgN.Queue = cfgN.SendTo
	notifier, err := rmq.NewNotifier(&cfgN)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create RMQ notifier: %v", err))
		os.Exit(1) //nolint:gocritic
	}
	defer func() {
		if err := notifier.Close(); err != nil {
			logger.Error(fmt.Sprintf("failed to close RMQ notifier: %v", err))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	ch, err := queue.ConsumeChannel(ctx, "calendar_sender")
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create consume channel: %v", err))
		os.Exit(1)
	}

	logger.Info("sender is running...")
	for msg := range ch {
		var notification models.Notification
		if err := json.Unmarshal(msg, &notification); err != nil {
			logger.Error(fmt.Sprintf("failed to unmarshal notification: %v", err))
			continue
		}
		if err := notifier.Notify(ctx, &notification); err != nil {
			logger.Error(fmt.Sprintf("failed to send notificaion %v: %v", notification.ID, err))
		}
	}

	<-ctx.Done()
	logger.Info("sender is stopping...")
}
