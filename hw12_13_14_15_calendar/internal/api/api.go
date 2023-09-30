package api

import (
	"context"
	"fmt"

	pb "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/api/apppb"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/calendar"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/mapper"
)

func NewAPI(app *calendar.App) pb.AppServer {
	return &api{app: app}
}

type api struct {
	pb.UnimplementedAppServer
	app *calendar.App
}

func (a *api) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	cmd, err := mapper.AddEventCommand(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	eventID, err := a.app.AddEvent(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	return &pb.AddEventResponse{EventId: eventID.String()}, nil
}

func (a *api) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	event, err := mapper.Event(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if err = a.app.UpdateEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("update event: %w", err)
	}
	return &pb.UpdateEventResponse{}, nil
}

func (a *api) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	eventID, err := mapper.EventID(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}
	if err = a.app.DeleteEvent(ctx, eventID); err != nil {
		return nil, fmt.Errorf("delete event: %w", err)
	}
	return &pb.DeleteEventResponse{}, nil
}

func (a *api) GetEventsOfDay(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := a.app.GetEventsOfDay(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}

func (a *api) GetEventsOfWeek(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := a.app.GetEventsOfWeek(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}

func (a *api) GetEventsOfMonth(ctx context.Context, req *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	events, err := a.app.GetEventsOfMonth(ctx, mapper.BeginOfDay(req))
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return mapper.GetEventsResponse(events), nil
}
