package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server"
)

type HTTPServer struct {
	app server.Application
	srv *http.Server
}

func NewServer(app server.Application) *HTTPServer {
	serv := &HTTPServer{app: app}

	router := mux.NewRouter()
	router.HandleFunc("/hello", serv.SayHello)
	router.HandleFunc("/AddEvent", serv.AddEvent).Methods("POST")
	router.HandleFunc("/GetEvent", serv.GetEvent).Methods("GET")
	router.HandleFunc("/GetEventsForPeriod", serv.GetEventsForPeriod).Methods("POST")
	router.HandleFunc("/GetEventsOfDay", serv.GetEventsOfDay).Methods("GET")
	router.HandleFunc("/GetEventsOfWeek", serv.GetEventsOfWeek).Methods("GET")
	router.HandleFunc("/GetEventsOfMonth", serv.GetEventsOfMonth).Methods("GET")
	router.HandleFunc("/ListEvents", serv.ListEvents).Methods("GET")
	router.HandleFunc("/UpdateEvent", serv.UpdateEvent).Methods("POST")
	router.HandleFunc("/DeleteEvent", serv.DeleteEvent).Methods("DELETE")

	addr := net.JoinHostPort(config.CalendarSettings.Server.Host, config.CalendarSettings.Server.HTTPPort)
	serv.srv = &http.Server{
		Addr:        addr,
		ReadTimeout: 5 * time.Second,
		Handler:     loggingMiddleware(router),
	}

	return serv
}

func (s *HTTPServer) Start(_ context.Context) error {
	logger.Info(fmt.Sprintf("http server is starting on %s", s.srv.Addr))
	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	logger.Info("HTTP server stopping...")
	err := s.srv.Shutdown(ctx)
	return err
}
