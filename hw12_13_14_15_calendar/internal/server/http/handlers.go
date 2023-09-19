package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	err := json.NewDecoder(r.Body).Decode(&ev)
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
	var idEv struct {
		id uuid.UUID `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&idEv)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := s.app.GetEvent(r.Context(), idEv.id)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
