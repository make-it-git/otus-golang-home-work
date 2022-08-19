package sqlstorage

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/dates"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	config config.ConnectionConf
	pool   *pgxpool.Pool
}

func New(c config.ConnectionConf) *Storage {
	return &Storage{
		config: c,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	pool, err := pgxpool.Connect(
		ctx,
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s pool_max_conns=10",
			s.config.User, s.config.Password, s.config.Host, s.config.Port, s.config.Database,
		),
	)
	if err != nil {
		return err
	}
	s.pool = pool
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.pool.Close()
	return nil
}

func (s *Storage) Create(event storage.Event) error {
	_, err := s.pool.Exec(
		context.Background(),
		`INSERT INTO events
    	(id, title, description, start_time, end_time, owner_id, notification_time) VALUES($1, $2, $3, $4, $5, $6, $7)`,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.StartTime.Add(event.Duration),
		event.OwnerID,
		event.NotificationTime,
	)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			return storage.ErrDuplicateID
		}
	}
	return err
}

func (s *Storage) Update(event storage.Event) error {
	_, err := s.pool.Exec(
		context.Background(),
		`UPDATE events
				SET title = $1,
				description = $2,
				start_time = $3,
				end_time = $4,
				owner_id = $5,
				notification_time = $6
				WHERE id = $7`,
		event.Title,
		event.Description,
		event.StartTime,
		event.StartTime.Add(event.Duration),
		event.OwnerID,
		event.NotificationTime,
		event.ID,
	)
	return err
}

func (s *Storage) Delete(id string) error {
	_, err := s.pool.Exec(context.Background(), "DELETE FROM events WHERE id = $1", id)
	return err
}

func (s *Storage) ListDay(date time.Time) ([]storage.Event, error) {
	start, end := dates.DayRange(date)
	return s.findInRange(start, end)
}

func (s *Storage) ListWeek(date time.Time) ([]storage.Event, error) {
	start, end := dates.WeekRange(date)
	return s.findInRange(start, end)
}

func (s *Storage) ListMonth(date time.Time) ([]storage.Event, error) {
	start, end := dates.MonthRange(date)
	return s.findInRange(start, end)
}

func (s *Storage) findInRange(start time.Time, end time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	rows, err := s.pool.Query(
		context.Background(),
		`SELECT id, title, description, start_time, end_time, owner_id, notification_time
			FROM events WHERE start_time >= $1 and start_time <= $2`,
		start,
		end,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ev storage.Event
		var endTime time.Time
		err = rows.Scan(&ev.ID, &ev.Title, &ev.Description, &ev.StartTime, &endTime, &ev.OwnerID, &ev.NotificationTime)
		if err != nil {
			return nil, err
		}
		ev.Duration = endTime.Sub(ev.StartTime)
		events = append(events, ev)
	}
	return events, nil
}
