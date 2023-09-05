package sqlstorage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct { // TODO
	db *sqlx.DB
}

func (s *Storage) InitStorage(ctx context.Context) {
	s.Connect(ctx)
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := getDsn()
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create DB connection: %v", err))
		os.Exit(1)
	}
	s.db = db
	return nil
}

func getDsn() string {
	dbURL := &url.URL{
		Scheme:   "postgres",
		Host:     config.Settings.Db.Host,
		User:     url.UserPassword(config.Settings.Db.User, config.Settings.Db.Password),
		Path:     config.Settings.Db.Name,
		RawQuery: "sslmode=disable",
	}
	logger.Debug(fmt.Sprintf("database: %v", dbURL.String()))
	return dbURL.String()
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) error {
	insertEventQuery := `
	INSERT INTO events 
	VALUES (:id, :title, :date_start, :date_end, :user_id, :description, :date_notification)
	`
	if _, err := s.db.NamedExecContext(ctx, insertEventQuery, event); err != nil {
		logger.Debug(fmt.Sprintf("insert event: %v", err))
		return fmt.Errorf("insert event: %w", err)
	}
	logger.Debug(fmt.Sprintf("insert event: %v", event))
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	getEventQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events 
	WHERE id = $1
	`
	var event storage.Event
	if err := s.db.GetContext(ctx, &event, getEventQuery, id); err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return &event, nil
}

func (s *Storage) GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error) {
	getEventsForPeriodQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events 
	WHERE date_start < $2 AND date_end >= $1
	ORDER BY date_start
	`
	events := make([]storage.Event, 0)
	if err := s.db.SelectContext(ctx, &events, getEventsForPeriodQuery, from, to); err != nil {
		return nil, fmt.Errorf("get events for period: %w", err)
	}
	return events, nil
}

func (s *Storage) ListEvents(ctx context.Context, limit, low uint64) ([]storage.Event, error) {
	listEventsQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events
	ORDER BY date_start
	LIMIT $1 OFFSET $2
	`
	events := make([]storage.Event, 0)
	if err := s.db.SelectContext(ctx, &events, listEventsQuery, limit, low); err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	updateEventQuery := `
	UPDATE events 
	SET 
		title=:title, 
		date_start=:date_start, 
		date_end=:date_end, 
		user_id=:user_id, 
		description=nullif(:description,''),
		date_notification=:date_notification
	WHERE id = :id
	`
	if _, err := s.db.NamedExecContext(ctx, updateEventQuery, event); err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	deleteEventQuery := `DELETE FROM events WHERE id = $1`
	if _, err := s.db.ExecContext(ctx, deleteEventQuery, id); err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}
