package storage

import (
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
