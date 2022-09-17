package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Sender struct {
	logger  logger.ILogger
	rabbit  *rabbit.Rabbit
	storage Storage
	done    chan struct{}
}

type Storage interface {
	NotificationHandled(ctx context.Context, id string, date time.Time) error
}

func New(logger logger.ILogger, storage Storage, rabbit *rabbit.Rabbit) *Sender {
	done := make(chan struct{})
	return &Sender{
		logger:  logger,
		storage: storage,
		rabbit:  rabbit,
		done:    done,
	}
}

func (s *Sender) Run() error {
	ch, err := s.rabbit.Consume()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-s.done:
				return
			case event, ok := <-ch:
				if !ok {
					return
				}
				ev := new(storage.Event)
				err := json.Unmarshal(event, ev)
				if err != nil {
					s.logger.Error(err)
					continue
				}
				d := ""
				if ev.Description != nil {
					d = *ev.Description
				}
				s.logger.Info(
					fmt.Sprintf(
						"Received notification for event %s: %s",
						ev.Title,
						d,
					),
				)
				err = s.storage.NotificationHandled(context.Background(), ev.ID, time.Now())
				if err != nil {
					s.logger.Error(err)
				}
			}
		}
	}()

	return nil
}

func (s *Sender) Stop() {
	s.logger.Info("Stopping sender")
	close(s.done)
}
