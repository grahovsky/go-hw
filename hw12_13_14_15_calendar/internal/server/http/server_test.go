package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/server/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServerHandlers(t *testing.T) {
	ctx := context.Background()
	eventID := uuid.MustParse("19317604-a433-4761-953a-c0ca625bfa26")
	testingEvent := models.Event{
		ID:          eventID,
		Title:       "test",
		Description: "some desc",
		UserID:      uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		DateStart:   time.Date(2023, 9, 23, 14, 0, 0, 0, time.UTC),
		DateEnd:     time.Date(2023, 9, 23, 15, 0, 0, 0, time.UTC),
	}

	cases := []struct {
		name      string
		url       string
		event     models.Event
		code      int
		exValue   interface{}
		respValue interface{}
		respError string
		mockError error
		appMock   func(t *testing.T) *mocks.Application
	}{
		{
			name:      "AddEvent. Success",
			url:       "/AddEvent",
			event:     testingEvent,
			code:      http.StatusCreated,
			exValue:   &eventID,
			respValue: &uuid.UUID{},
			appMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("AddEvent", mock.Anything, &testingEvent).Return(eventID, nil)
				return appMock
			},
		},
		{
			name:      "AddEvent. Err event ID error",
			url:       "/AddEvent",
			event:     testingEvent,
			code:      http.StatusInternalServerError,
			mockError: models.ErrEventID,
			exValue:   &eventID,
			respValue: &uuid.UUID{},
			appMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("AddEvent", mock.Anything, &testingEvent).Return(eventID, models.ErrEventID)
				return appMock
			},
		},
		{
			name:  "UpdateEvent. Success",
			url:   "/UpdateEvent",
			event: testingEvent,
			code:  http.StatusOK,
			appMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("UpdateEvent", mock.Anything, &testingEvent).Return(nil)
				return appMock
			},
		},
		{
			name:  "UpdateEvent. Success",
			url:   "/UpdateEvent",
			event: testingEvent,
			code:  http.StatusInternalServerError,
			appMock: func(t *testing.T) *mocks.Application {
				t.Helper()

				appMock := mocks.NewApplication(t)
				appMock.On("UpdateEvent", mock.Anything, &testingEvent).Return(models.ErrEventID)
				return appMock
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jsonData, err := json.Marshal(tc.event)
			require.NoError(t, err)
			reqBody := bytes.NewBuffer(jsonData)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.url, reqBody)
			require.NoError(t, err)

			appMock := tc.appMock(t)
			serv := NewServer(appMock, fmt.Sprintf("%v:%v", "0.0.0.0", "8081"))

			rr := httptest.NewRecorder()
			serv.srv.Handler.ServeHTTP(rr, req)
			require.Equal(t, tc.code, rr.Code)

			if tc.respError == "" && tc.mockError == nil && tc.exValue != nil && tc.respValue != nil {
				err = json.Unmarshal(rr.Body.Bytes(), tc.respValue)
				require.NoError(t, err)

				require.Equal(t, tc.exValue, tc.respValue)
			}

			appMock.AssertExpectations(t)
		})
	}
}
