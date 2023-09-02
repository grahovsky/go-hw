package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	var used_storage app.Storage

	if config.Settings.Storage.Type == "sql" {
		used_storage = &sqlstorage.Storage{}
	} else {
		used_storage = &memorystorage.Storage{}
	}
	used_storage.Create()

	calendar := app.New(used_storage)

	server := internalhttp.NewServer(calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	event := storage.Event{
		ID: "111",
	}
	used_storage.AddEvent(ctx, event)
	calendar.CreateEvent(ctx, "222", "some title")
	fmt.Println(used_storage.GetSortedEventsById("111"))
	fmt.Println(used_storage.GetSortedEventsById("222"))

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
