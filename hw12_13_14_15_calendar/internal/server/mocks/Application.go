// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// AddEvent provides a mock function with given fields: _a0, _a1
func (_m *Application) AddEvent(_a0 context.Context, _a1 *models.Event) (uuid.UUID, error) {
	ret := _m.Called(_a0, _a1)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event) (uuid.UUID, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event) uuid.UUID); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Event) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEvent provides a mock function with given fields: ctx, id
func (_m *Application) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEvent provides a mock function with given fields: _a0, _a1
func (_m *Application) GetEvent(_a0 context.Context, _a1 uuid.UUID) (*models.Event, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*models.Event, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *models.Event); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEventsForPeriod provides a mock function with given fields: ctx, from, to
func (_m *Application) GetEventsForPeriod(ctx context.Context, from time.Time, to time.Time) ([]models.Event, error) {
	ret := _m.Called(ctx, from, to)

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) ([]models.Event, error)); ok {
		return rf(ctx, from, to)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) []models.Event); ok {
		r0 = rf(ctx, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListEvents provides a mock function with given fields: ctx, limit, low
func (_m *Application) ListEvents(ctx context.Context, limit uint64, low uint64) ([]models.Event, error) {
	ret := _m.Called(ctx, limit, low)

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) ([]models.Event, error)); ok {
		return rf(ctx, limit, low)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) []models.Event); ok {
		r0 = rf(ctx, limit, low)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, limit, low)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEvent provides a mock function with given fields: ctx, event
func (_m *Application) UpdateEvent(ctx context.Context, event *models.Event) error {
	ret := _m.Called(ctx, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewApplication creates a new instance of Application. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplication(t interface {
	mock.TestingT
	Cleanup(func())
}) *Application {
	mock := &Application{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
