package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger    Logger
	storage   Storage
	validator *validator.Validate
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
		logger:    logger,
		storage:   storage,
		validator: validator.New(),
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	err := a.validator.Struct(event)
	if err != nil {
		return err
	}
	return a.storage.Create(event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	err := a.validator.Struct(event)
	if err != nil {
		return err
	}
	return a.storage.Update(event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.Delete(id)
}

func (a *App) ListDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListDay(date)
}

func (a *App) ListWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListWeek(date)
}

func (a *App) ListMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListMonth(date)
}
