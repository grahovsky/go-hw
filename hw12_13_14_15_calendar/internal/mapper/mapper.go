package mapper

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api/eventservice"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AddEventCommand(req *eventservice.AddEventRequest) (*models.Event, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid UserID: %w", err)
	}

	cmd := models.Event{
		Title:            req.Title,
		DateStart:        req.DateStart.AsTime(),
		DateEnd:          req.DateEnd.AsTime(),
		Description:      req.Description,
		UserID:           userID,
		DateNotification: req.DateNotification.AsTime(),
	}

	return &cmd, nil
}

func Event(req *eventservice.UpdateEventRequest) (*models.Event, error) {
	if req.Event == nil {
		return nil, errors.New("event field is empty")
	}

	eventID, err := uuid.Parse(req.Event.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid EventID: %w", err)
	}

	userID, err := uuid.Parse(req.Event.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid UserID: %w", err)
	}

	return &models.Event{
		ID:               eventID,
		Title:            req.Event.Title,
		DateStart:        req.Event.DateStart.AsTime(),
		DateEnd:          req.Event.DateEnd.AsTime(),
		Description:      req.Event.Description,
		UserID:           userID,
		DateNotification: req.Event.DateNotification.AsTime(),
	}, nil
}

func EventID(req *eventservice.DeleteEventRequest) (uuid.UUID, error) {
	eventID, err := uuid.Parse(req.EventId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid eventID: %w", err)
	}
	return eventID, nil
}

func BeginOfDay(req *eventservice.GetEventsRequest) time.Time {
	return req.Since.AsTime().Truncate(24 * time.Hour)
}

func GetEventsResponse(events []models.Event) *eventservice.GetEventsResponse {
	mapped := make([]*eventservice.Event, len(events))
	for i, event := range events {
		mapped[i] = &eventservice.Event{
			Id:               event.ID.String(),
			DateStart:        timestamppb.New(event.DateStart),
			DateEnd:          timestamppb.New(event.DateEnd),
			Description:      event.Description,
			UserId:           event.UserID.String(),
			DateNotification: timestamppb.New(event.DateNotification),
		}
	}
	return &eventservice.GetEventsResponse{Events: mapped}
}
