package storage

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	ID               string    `db:"id"`
	Title            string    `db:"title"`
	DateStart        time.Time `db:"date_start"`
	DateEnd          time.Time `db:"date_end"`
	UserID           int64     `db:"user_id"`
	Description      string    `db:"description"`
	DateNotification time.Time `db:"date_notification"`
}

var (
	ErrEventID  = errors.New("Event ID is not correct")
	ErrDateBusy = errors.New("Event Date is bussy")
)

func (e Event) String() string {
	return fmt.Sprintf("%v %v %v", e.ID, e.Title, e.DateStart)
}
