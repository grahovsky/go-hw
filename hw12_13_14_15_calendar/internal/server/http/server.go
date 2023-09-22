package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

type Application interface {
	AddEvent(context.Context, *models.Event) error
	GetEvent(context.Context, uuid.UUID) (*models.Event, error)
	GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]models.Event, error)
	ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type Server struct {
	app Application
	srv *http.Server
}

func NewServer(app Application, addr string) *Server {
	serv := &Server{app: app}

	router := mux.NewRouter()
	router.HandleFunc("/hello", serv.SayHello)
	router.HandleFunc("/AddEvent", serv.AddEvent).Methods("POST")
	router.HandleFunc("/GetEvent", serv.GetEvent).Methods("GET")
	router.HandleFunc("/GetEventsForPeriod", serv.GetEventsForPeriod).Methods("POST")
	router.HandleFunc("/ListEvents", serv.ListEvents).Methods("GET")
	router.HandleFunc("/UpdateEvent", serv.UpdateEvent).Methods("POST")
	router.HandleFunc("/DeleteEvent", serv.DeleteEvent).Methods("DELETE")

	serv.srv = &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     loggingMiddleware(router),
	}

	logger.Info(fmt.Sprintf("create server: %v", addr))

	return serv
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
	logger.Info("HTTP server stopping..")
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
