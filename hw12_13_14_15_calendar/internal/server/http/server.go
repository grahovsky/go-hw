package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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
	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
	}

	srv.Handler = initRoutes()

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

func initRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(loggingMiddleware())

	router.Get("/hello", SayHello())

	return router
}

func SayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "hello world")
	}
}
