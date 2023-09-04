package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID               uuid.UUID `db:"id"`
	Title            string    `db:"title"`
	DateStart        time.Time `db:"date_start"`
	DateEnd          time.Time `db:"date_end"`
	UserID           uuid.UUID `db:"user_id"`
	Description      string    `db:"description"`
	DateNotification time.Time `db:"date_notification"`
}

var (
	ErrEventID            = errors.New("Event ID is not correct")
	ErrDateBusy           = errors.New("Event Date is bussy")
	ErrEventAlreadyExists = errors.New("Event already exists")
)

func (e Event) String() string {
	return fmt.Sprintf("%v %v %v", e.ID, e.Title, e.DateStart)
}
