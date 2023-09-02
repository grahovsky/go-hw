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
	UserID           string    `db:"user_id"`
	Description      string    `db:"description"`
	DateNotification time.Time `db:"date_notification"`
}

var ErrEventID = errors.New("Event ID is not correct")

func (e Event) String() string {
	return fmt.Sprintf("%v %v", e.ID, e.Title)
}
