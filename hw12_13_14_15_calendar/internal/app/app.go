package app

import (
	"context"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(data interface{})
	Error(data interface{})
	Debug(data interface{})
	Warn(data interface{})
}

type Storage interface {
	Create(event storage.Event) error
	Update(event storage.Event) error
	Delete(id string) error
	ListDay(date time.Time) ([]storage.Event, error)
	ListWeek(date time.Time) ([]storage.Event, error)
	ListMonth(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
