package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID               string        `validate:"required,uuid4"`
	Title            string        `validate:"required"`
	StartTime        time.Time     `validate:"required"`
	Duration         time.Duration `validate:"required"`
	Description      *string       `validate:"omitempty,max=256"`
	OwnerID          int32         `validate:"required,min=1"`
	NotificationTime *time.Time
}

var ErrMissingID = errors.New("id not provided")

var ErrDuplicateID = errors.New("id duplicate")

var ErrNotFoundID = errors.New("id not found")
