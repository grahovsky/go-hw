package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Application interface {
	CreateEvent(context.Context, *storage.Event) error
	// UpdateEvent(context.Context, app.Event) error
	// DeleteEvent(context.Context, string) error
	// GetEventByDay(context.Context, int64, time.Time) ([]app.Event, error)
	// GetEventByWeek(context.Context, int64, time.Time) ([]app.Event, error)
	// GetEventByMonth(context.Context, int64, time.Time) ([]app.Event, error)
}

type Server struct {
	app *Application
	srv *http.Server
}

func NewServer(app *Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", SayHello())

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     loggingMiddleware(mux),
	}

	return &Server{
		app: app,
		srv: srv,
	}
}

func (s *Server) Start(_ context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func SayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!\n")); err != nil {
			logger.Error(fmt.Sprintf("failed to write response: %v", err))
		}
	}
}
