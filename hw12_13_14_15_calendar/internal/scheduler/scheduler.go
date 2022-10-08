package scheduler

import (
	"context"
	"encoding/json"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	logger      logger.ILogger
	storage     Storage
	rabbit      *rabbit.Rabbit
	ticker      *time.Ticker
	done        chan bool
	cleanupDays uint
}

type Storage interface {
	FindDue(ctx context.Context, date time.Time) ([]storage.Event, error)
	CleanupTill(ctx context.Context, date time.Time) error
	Notified(ctx context.Context, id string, date time.Time) error
}

func New(logger logger.ILogger, storage Storage, rabbit *rabbit.Rabbit, conf config.RabbitmqTimerConf, cleanup config.CleanupConf) *Scheduler {
	ticker := time.NewTicker(time.Duration(conf.Wait) * time.Second)
	done := make(chan bool)
	return &Scheduler{
		logger:      logger,
		storage:     storage,
		rabbit:      rabbit,
		ticker:      ticker,
		done:        done,
		cleanupDays: cleanup.Days,
	}
}

func (s *Scheduler) Run() {
	defer s.ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-s.ticker.C:
			err := s.send()
			if err != nil {
				s.logger.Error(err)
				return
			}
			err = s.cleanup()
			if err != nil {
				s.logger.Error(err)
				return
			}
		}
	}
}

func (s *Scheduler) Stop() {
	s.logger.Info("Stopping scheduler")
	s.done <- true
}

func (s *Scheduler) send() error {
	ctx := context.Background()

	events, err := s.storage.FindDue(ctx, time.Now())
	if err != nil {
		return err
	}

	for _, ev := range events {
		b, err := json.Marshal(ev)
		if err != nil {
			return err
		}
		err = s.rabbit.Publish(b)
		if err != nil {
			return err
		}
		err = s.storage.Notified(ctx, ev.ID, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) cleanup() error {
	before := time.Duration(s.cleanupDays) * time.Hour * 24
	t := time.Now().Add(-before)
	return s.storage.CleanupTill(context.Background(), t)
}
