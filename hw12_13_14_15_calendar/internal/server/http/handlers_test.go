package internalhttp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/stretchr/testify/require"
)

func TestAppMockGetEvent(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.New()
	existingEvent := models.Event{
		ID:          eventID,
		Title:       "some title",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(1 * time.Hour),
		UserID:      uuid.UUID{},
		Description: "some desc for test",
	}

	tests := []struct {
		title      string
		event      *models.Event
		err        error
		GetAppMock func(t *testing.T) *mocks.Application
	}{
		{
			title: "positive case",
			event: &existingEvent,
			GetAppMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("GetEvent", ctx, eventID).Return(&existingEvent, nil)
				return appMock
			},
		},
		{
			title: "event not found",
			err:   models.ErrEventNotFound,
			GetAppMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("GetEvent", ctx, eventID).Return(&models.Event{}, models.ErrEventNotFound)
				return appMock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			appMock := tt.GetAppMock(t)
			server := NewServer(appMock, fmt.Sprintf("%v:%v", config.Settings.Server.Host, config.Settings.Server.HTTPPort))

			res, err := server.app.GetEvent(ctx, eventID)
			if tt.err != nil {
				require.Equal(t, tt.err, err)
			} else {
				require.Equal(t, tt.event, res)
				require.Equal(t, nil, err)
			}

			appMock.AssertExpectations(t)
		})
	}
}
