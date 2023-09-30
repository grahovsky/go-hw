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
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/notify"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/schedule"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

func main() {
	config.InitSchedulerSettings()
	logger.SetLogLevel(config.SchedulerSettings.Log.Level)

	st, err := storage.New(config.SchedulerSettings.Storage)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := st.Close(); err != nil {
			logger.Error(fmt.Sprintf("faield to close storage: %v", err))
		}
	}()

	notifier := notify.NewLogNotifier()

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

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			scheduler.LoadSchedule(ctx)
		}
	}
	<-ctx.Done()
	scheduler.Wait()

	logger.Info("scheduler is stopping...")
}
