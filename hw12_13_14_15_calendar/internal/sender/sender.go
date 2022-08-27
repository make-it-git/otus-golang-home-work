package sender

import (
	"encoding/json"
	"fmt"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Sender struct {
	logger logger.ILogger
	rabbit *rabbit.Rabbit
	done   chan struct{}
}

func New(logger logger.ILogger, rabbit *rabbit.Rabbit) *Sender {
	done := make(chan struct{})
	return &Sender{
		logger: logger,
		rabbit: rabbit,
		done:   done,
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
			}
		}
	}()

	return nil
}

func (s *Sender) Stop() {
	s.logger.Info("Stopping sender")
	close(s.done)
}
