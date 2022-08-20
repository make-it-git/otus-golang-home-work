package app

import (
	"context"
	v "github.com/go-playground/validator/v10"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logic"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger    Logger
	storage   Storage
	validator *v.Validate
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
		validator: v.New(),
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	err := a.validator.Struct(event)
	if err != nil {
		return err
	}
	if event.Duration < 0 {
		return logic.ErrEndTimeBeforeStartTime
	}
	return a.storage.Create(event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	err := a.validator.Struct(event)
	if err != nil {
		return err
	}
	if event.Duration < 0 {
		return logic.ErrEndTimeBeforeStartTime
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
