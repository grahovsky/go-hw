package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	DateStart        time.Time `db:"date_start" json:"dateStart"`
	DateEnd          time.Time `db:"date_end" json:"dateEnd"`
	UserID           uuid.UUID `db:"user_id" json:"userId"`
	Description      string    `db:"description" json:"description"`
	DateNotification time.Time `db:"date_notification" json:"dateNotification"`
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
