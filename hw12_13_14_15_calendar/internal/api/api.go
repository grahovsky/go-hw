package api

import (
	"context"
	"fmt"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api/eventservice"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/mapper"
)

func NewAPI(calendar *calendar.App) eventservice.CalendarServer {
	return &api{calendar: calendar}
}

type api struct {
	eventservice.UnimplementedCalendarServer
	calendar *calendar.App
}

func (a *api) AddEvent(ctx context.Context, req *eventservice.AddEventRequest) (*eventservice.AddEventResponse, error) {
	cmd, err := mapper.AddEventCommand(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	eventID, err := a.calendar.AddEvent(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return &eventservice.AddEventResponse{EventId: eventID.String()}, nil
}

func (a *api) UpdateEvent(ctx context.Context,
	req *eventservice.UpdateEventRequest,
) (*eventservice.UpdateEventResponse, error) {
	event, err := mapper.Event(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if err = a.calendar.UpdateEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("update event: %w", err)
	}
	return &eventservice.UpdateEventResponse{}, nil
}

func (a *api) DeleteEvent(ctx context.Context,
	req *eventservice.DeleteEventRequest,
) (*eventservice.DeleteEventResponse, error) {
	eventID, err := mapper.EventID(req.EventId)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if err = a.calendar.DeleteEvent(ctx, eventID); err != nil {
		return nil, fmt.Errorf("delete event: %w", err)
	}
	return &eventservice.DeleteEventResponse{}, nil
}

func (a *api) GetEventsOfDay(ctx context.Context,
	req *eventservice.GetEventsRequest,
) (*eventservice.GetEventsResponse, error) {
	events, err := a.calendar.GetEventsOfDay(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}

func (a *api) GetEventsOfWeek(ctx context.Context,
	req *eventservice.GetEventsRequest,
) (*eventservice.GetEventsResponse, error) {
	events, err := a.calendar.GetEventsOfWeek(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}

func (a *api) GetEventsOfMonth(ctx context.Context,
	req *eventservice.GetEventsRequest,
) (*eventservice.GetEventsResponse, error) {
	events, err := a.calendar.GetEventsOfMonth(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}

func (a *api) GetEvent(ctx context.Context,
	req *eventservice.GetEventRequest,
) (*eventservice.GetEventResponse, error) {
	eventID, err := mapper.EventID(req.EventId)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	event, err := a.calendar.GetEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return mapper.GetEventResponse(event), nil
}
