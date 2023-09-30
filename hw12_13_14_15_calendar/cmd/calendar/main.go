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
	internalgrpc "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	config.InitCalendarSettings()
	if config.CalendarSettings.PrintVersion {
		PrintVersion()
		return
	}

	logger.SetLogLevel(config.CalendarSettings.Log.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var usedStorage storage.Storage
	if config.CalendarSettings.Storage.Type == "sql" {
		usedStorage = &sqlstorage.Storage{}
	} else {
		usedStorage = &memorystorage.Storage{}
	}
	usedStorage.InitStorage()
	defer usedStorage.Close()

	// var calendar internalhttp.Application
	calendar := app.New(usedStorage)
	httpSrv := internalhttp.NewServer(calendar,
		fmt.Sprintf("%v:%v", config.CalendarSettings.Server.Host, config.CalendarSettings.Server.HTTPPort))
	grpcSrv := internalgrpc.NewServer(calendar,
		fmt.Sprintf("%v:%v", config.CalendarSettings.Server.Host, config.CalendarSettings.Server.GRPCPort))

	go func() {
		if err := httpSrv.Start(ctx); err != nil {
			logger.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
		logger.Info("http is running...")
	}()

	go func() {
		if err := grpcSrv.Start(); err != nil {
			logger.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
		logger.Info("grpc is running...")
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	grpcSrv.Stop()
	if err := httpSrv.Stop(ctx); err != nil {
		logger.Error("failed to stop http server: " + err.Error())
	}
}
