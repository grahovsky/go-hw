package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/schedule"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

func main() {
	config.InitSchedulerSettings()
	logger.SetLogLevel(config.SchedulerSettings.Log.Level)

	st, err := storage.New(&config.SchedulerSettings.Storage)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to creqte storage: %v", err))
		os.Exit(1)
	}
	defer func() {
		if err := st.Close(); err != nil {
			logger.Error(fmt.Sprintf("faield to close storage: %v", err))
		}
	}()

	notifier, err := rmq.NewNotifier(&config.SchedulerSettings.Rmq)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to creqte RMQ notifier: %v", err))
		os.Exit(1) //nolint:gocritic
	}
	defer func() {
		if err := notifier.Close(); err != nil {
			logger.Error(fmt.Sprintf("faield to close RMQ notifier: %v", err))
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	logger.Info("scheduler is running...")
	scheduler := schedule.NewScheduler(st, notifier)
	scheduler.LoadSchedule(ctx)

	ticker := time.NewTicker(config.SchedulerSettings.Schedule)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				scheduler.LoadSchedule(ctx)
			}
		}
	}()
	<-ctx.Done()

	logger.Info("scheduler is stopping...")
	scheduler.Wait()
}
