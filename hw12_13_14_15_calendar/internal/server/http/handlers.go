package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	id, err := parseParamUuid(r, "id")
	if err != nil {
		logger.Error(fmt.Sprintf("error to get id: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		logger.Error(fmt.Sprintf("error get event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetEventsForPeriod(w http.ResponseWriter, r *http.Request) {
	var filter struct {
		DateFrom time.Time `json:"dateFrom"`
		DateTo   time.Time `json:"dateTo"`
	}

	err := parseBody(r, &filter)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get filter: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.app.GetEventsForPeriod(r.Context(), filter.DateFrom, filter.DateTo)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		logger.Error(fmt.Sprintf("error get events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ListEvents(w http.ResponseWriter, r *http.Request) {
	limit, err := parseParamUint64(r, "limit")
	if err != nil {
		logger.Error(fmt.Sprintf("error to get param: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	low, err := parseParamUint64(r, "low")
	if err != nil {
		logger.Error(fmt.Sprintf("error to get param: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := s.app.ListEvents(r.Context(), limit, low)
	if err != nil {
		logger.Error(fmt.Sprintf("error to get events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		logger.Error(fmt.Sprintf("error get events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
