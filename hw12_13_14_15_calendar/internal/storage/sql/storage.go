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
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/models"
	_ "github.com/jackc/pgx/stdlib" // justifying
	"github.com/jmoiron/sqlx"
)

type Storage struct { // TODO
	db *sqlx.DB
}

func (s *Storage) InitStorage(settings *config.Storage) error {
	return s.Connect(settings)
}

func (s *Storage) Connect(settings *config.Storage) error {
	dsn := GetDsn(settings)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create DB connection: %v", err))
		os.Exit(1)
	}
	s.db = db
	return nil
}

func GetDsn(settings *config.Storage) string {
	dbURL := &url.URL{
		Scheme:   "postgres",
		Host:     settings.DB.Host,
		User:     url.UserPassword(settings.DB.User, settings.DB.Password),
		Path:     settings.DB.Name,
		RawQuery: "sslmode=disable",
	}
	logger.Debug(fmt.Sprintf("database: %v", dbURL.String()))
	return dbURL.String()
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(ctx context.Context, event *models.Event) error {
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

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	getEventQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events 
	WHERE id = $1
	`
	var event models.Event
	if err := s.db.GetContext(ctx, &event, getEventQuery, id); err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return &event, nil
}

func (s *Storage) GetEventsForPeriod(ctx context.Context, from, to time.Time) ([]models.Event, error) {
	getEventsForPeriodQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events 
	WHERE date_start < $2 AND date_end >= $1
	ORDER BY date_start
	`
	events := make([]models.Event, 0)
	if err := s.db.SelectContext(ctx, &events, getEventsForPeriodQuery, from, to); err != nil {
		return nil, fmt.Errorf("get events for period: %w", err)
	}
	return events, nil
}

func (s *Storage) ListEvents(ctx context.Context, limit, low uint64) ([]models.Event, error) {
	listEventsQuery := `
	SELECT id, title, date_start, date_end, user_id, coalesce(description, '') as description, date_notification 
	FROM events
	ORDER BY date_start
	LIMIT $1 OFFSET $2
	`
	events := make([]models.Event, 0)
	if err := s.db.SelectContext(ctx, &events, listEventsQuery, limit, low); err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) error {
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
	r, err := s.db.NamedExecContext(ctx, updateEventQuery, event)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	updated, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if updated == 0 {
		return fmt.Errorf("event not found")
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

func (s *Storage) DeleteEventsBefore(ctx context.Context, before time.Time) (int64, error) {
	deletePastEventsQuery := `DELETE FROM events WHERE date_end < $1 IS TRUE RETURNING id`
	r, err := s.db.ExecContext(ctx, deletePastEventsQuery, before)
	if err != nil {
		return 0, fmt.Errorf("delete old event: %w", err)
	}
	return r.RowsAffected()
}

func (s *Storage) GetEventsToNotify(ctx context.Context, from, to time.Time) ([]models.Event, error) {
	getEventToNotifyQuery := `
	SELECT id, title, date_start, date_end, coalesce(description, '') as description, user_id, date_notification 
	FROM events 
	WHERE date_notification >= $1 AND date_notification < $2
	ORDER BY date_start
	`
	events := make([]models.Event, 0)
	if err := s.db.SelectContext(ctx, &events, getEventToNotifyQuery, from, to); err != nil {
		return nil, fmt.Errorf("get events to notify: %w", err)
	}
	return events, nil
}
