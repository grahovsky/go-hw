package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
)

func (s *Server) SayHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World!\n")); err != nil {
		logger.Error(fmt.Sprintf("failed to write response: %v", err))
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) AddEvent(w http.ResponseWriter, r *http.Request) {
	ev := storage.Event{}
	err := parseBody(r, &ev)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.AddEvent(r.Context(), &ev)
	if err != nil {
		logger.Error(fmt.Sprintf("error to create event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get id: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		logger.Error(fmt.Sprintf("error get event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
