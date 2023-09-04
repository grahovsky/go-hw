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
	AddEvent(context.Context, *storage.Event) error
	// GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error)
	// GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error)
	// ListEvents(ctx context.Context, limit, low uint64) ([]storage.Event, error)
	// UpdateEvent(ctx context.Context, event *storage.Event) error
	// DeleteEvent(ctx context.Context, id uuid.UUID) error
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
