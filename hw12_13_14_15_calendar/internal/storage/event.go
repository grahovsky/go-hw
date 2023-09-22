package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	DateStart        time.Time `db:"date_start" json:"date_start"`
	DateEnd          time.Time `db:"date_end" json:"date_end"`
	UserID           uuid.UUID `db:"user_id" json:"user_id"`
	Description      string    `db:"description" json:"description"`
	DateNotification time.Time `db:"date_notification" json:"date_notification"`
}

var (
	ErrEventID            = errors.New("Event ID is not correct")
	ErrEventNotFound      = errors.New("Event not found")
	ErrEventAlreadyExists = errors.New("Event already exists")
)

func (e *Event) String() string {
	return fmt.Sprintf("%v %v %v", e.ID, e.Title, e.DateStart)
}

func (e *Event) InPeriod(from time.Time, to time.Time) bool {
	return e.DateStart.Before(to) && !e.DateEnd.Before(from)
}
