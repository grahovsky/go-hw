package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

func main() {
	config.InitCalendarSettings()
	if config.CalendarSettings.PrintVersion {
		PrintVersion()
		return
	}

	logger.SetLogLevel(config.CalendarSettings.Log.Level)

	st, err := storage.New(config.CalendarSettings.Storage)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := st.Close(); err != nil {
			logger.Error(fmt.Sprintf("faield to close storage: %v", err))
		}
	}()

	calendar := calendar.New(st)
	httpSrv := internalhttp.NewServer(calendar)
	grpcSrv := internalgrpc.NewServer(calendar)

	logger.Info("calendar is running...")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	go func() {
		logger.Info("http is running...")
		if err := httpSrv.Start(ctx); err != nil {
			logger.Error(fmt.Sprintf("failed to start http server: %v", err))
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		logger.Info("grpc is running...")
		if err := grpcSrv.Start(); err != nil {
			logger.Error(fmt.Sprintf("failed to start grpc server: %v", err))
			cancel()
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	logger.Info("calendar is stopping...")
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := httpSrv.Stop(ctx); err != nil {
		logger.Error(fmt.Sprintf("failed to stop http server: %v", err))
	}
	grpcSrv.Stop()
}
