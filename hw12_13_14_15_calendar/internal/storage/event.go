package storage

import (
	"errors"
	"time"
)

type Storage interface {
	Create(event Event) error
	Update(event Event) error
	Delete(id string) error
	ListDay(date time.Time) ([]Event, error)
	ListWeek(date time.Time) ([]Event, error)
	ListMonth(date time.Time) ([]Event, error)
}

type Event struct {
	ID               string
	Title            string
	StartTime        time.Time
	Duration         time.Duration
	Description      *string
	OwnerId          int
	NotificationTime *time.Time
}

var ErrMissingId = errors.New("id not provided")
var ErrDuplicateId = errors.New("id duplicate")
var ErrNotFoundId = errors.New("id not found")
