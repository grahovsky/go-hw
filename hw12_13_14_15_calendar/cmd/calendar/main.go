package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	if config.Settings.PrintVersion {
		PrintVersion()
		return
	}

	logger.SetLogLevel(config.Settings.Log.Level)
	logger.Debug(config.Settings.DebugMessage)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var usedStorage app.Storage
	if config.Settings.Storage.Type == "sql" {
		usedStorage = &sqlstorage.Storage{}
	} else {
		usedStorage = &memorystorage.Storage{}
	}
	usedStorage.InitStorage()
	defer usedStorage.Close()

	// var calendar internalhttp.Application
	calendar := app.New(usedStorage)
	server := internalhttp.NewServer(calendar, "localhost:8080")

	uid := uuid.New()

	newEvent := storage.Event{
		ID:          uid,
		Title:       "first event",
		DateStart:   time.Date(2023, 9, 4, 22, 0, 0, 0, time.Local),
		DateEnd:     time.Date(2023, 9, 4, 23, 40, 0, 0, time.Local),
		UserID:      uuid.New(),
		Description: "some description",
	}
	usedStorage.AddEvent(ctx, &newEvent)
	calendar.AddEvent(ctx,
		&storage.Event{
			ID:        uuid.New(),
			Title:     "some title",
			DateStart: time.Now().Add(4 * time.Hour),
			DateEnd:   time.Now().Add(5 * time.Hour),
		},
	)

	if events, err := usedStorage.ListEvents(ctx, 10, 0); err == nil {
		for _, event := range events {
			logger.Info(event.String())
		}
	}

	if events, err := usedStorage.GetEventsForPeriod(ctx,
		time.Date(2023, 9, 4, 23, 0, 0, 0, time.Local),
		time.Date(2023, 9, 4, 23, 59, 0, 0, time.Local),
	); err == nil {
		for _, event := range events {
			logger.Info(event.String())
		}
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
