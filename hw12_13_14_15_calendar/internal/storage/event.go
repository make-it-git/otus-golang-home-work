package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID               string
	Title            string
	StartTime        time.Time
	Duration         time.Duration
	Description      *string
	OwnerID          int32
	NotificationTime *time.Time
}

var ErrMissingID = errors.New("id not provided")

var ErrDuplicateID = errors.New("id duplicate")

var ErrNotFoundID = errors.New("id not found")
