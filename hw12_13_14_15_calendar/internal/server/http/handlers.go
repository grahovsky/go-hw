package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
)

func (s *HTTPServer) SayHello(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Hello World!\n")); err != nil {
		logger.Error(fmt.Sprintf("failed to write response: %v", err))
	}
	w.WriteHeader(http.StatusOK)
}

func (s *HTTPServer) AddEvent(w http.ResponseWriter, r *http.Request) {
	ev := models.Event{}
	err := parseBody(r, &ev)
	if checkError(w, err) {
		return
	}
	id, err := s.app.AddEvent(r.Context(), &ev)
	if checkError(w, err) {
		return
	}
	// status need set before write json data!!!
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if checkError(w, err) {
		return
	}
}

func (s *HTTPServer) GetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseParamUUID(r, "id")
	if checkError(w, err) {
		return
	}
	event, err := s.app.GetEvent(r.Context(), id)
	if checkError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(event)
	if checkError(w, err) {
		return
	}
}

func (s *HTTPServer) GetEventsForPeriod(w http.ResponseWriter, r *http.Request) {
	var filter struct {
		DateFrom time.Time `json:"dateFrom"`
		DateTo   time.Time `json:"dateTo"`
	}
	err := parseBody(r, &filter)
	if checkError(w, err) {
		return
	}
	events, err := s.app.GetEventsForPeriod(r.Context(), filter.DateFrom, filter.DateTo)
	if checkError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(events)
	if checkError(w, err) {
		return
	}
}

func (s *HTTPServer) ListEvents(w http.ResponseWriter, r *http.Request) {
	limit, err := parseParamUint64(r, "limit")
	if checkError(w, err) {
		return
	}
	low, err := parseParamUint64(r, "low")
	if checkError(w, err) {
		return
	}
	events, err := s.app.ListEvents(r.Context(), limit, low)
	if checkError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(events)
	if checkError(w, err) {
		return
	}
}

func (s *HTTPServer) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ev := models.Event{}
	err := parseBody(r, &ev)
	if checkError(w, err) {
		return
	}
	err = s.app.UpdateEvent(r.Context(), &ev)
	if checkError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *HTTPServer) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := parseParamUUID(r, "id")
	if checkError(w, err) {
		return
	}
	err = s.app.DeleteEvent(r.Context(), id)
	if checkError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
}
