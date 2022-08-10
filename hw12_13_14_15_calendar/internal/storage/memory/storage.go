package memorystorage

import (
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/dates"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	events []storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make([]storage.Event, 0),
	}
}

func (s *Storage) Create(event storage.Event) error {
	if event.ID == "" {
		return storage.ErrMissingId
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.events {
		if s.events[i].ID == event.ID {
			return storage.ErrDuplicateId
		}
	}
	s.events = append(s.events, event)
	return nil
}

func (s *Storage) Update(event storage.Event) error {
	for i := range s.events {
		if s.events[i].ID == event.ID {
			s.mu.Lock()
			s.events[i] = event
			s.mu.Unlock()
			return nil
		}
	}
	return storage.ErrNotFoundId
}

func (s *Storage) Delete(id string) error {
	for i := range s.events {
		if s.events[i].ID == id {
			s.mu.Lock()
			s.events = append(s.events[:i], s.events[i+1:]...)
			s.mu.Unlock()
			return nil
		}
	}
	return storage.ErrNotFoundId
}

func (s *Storage) ListDay(date time.Time) ([]storage.Event, error) {
	year, month, day := date.Date()
	r := make([]storage.Event, 0)
	for i := range s.events {
		y, m, d := s.events[i].StartTime.Date()
		if y == year && m == month && d == day {
			r = append(r, s.events[i])
		}
	}
	return r, nil
}

func (s *Storage) ListWeek(date time.Time) ([]storage.Event, error) {
	start, end := dates.WeekRange(date)
	return s.findInRange(start, end)
}

func (s *Storage) ListMonth(date time.Time) ([]storage.Event, error) {
	start, end := dates.MonthRange(date)
	return s.findInRange(start, end)
}

func (s *Storage) findInRange(start time.Time, end time.Time) ([]storage.Event, error) {
	r := make([]storage.Event, 0)
	for i := range s.events {
		d := s.events[i].StartTime
		if (d == start || d.After(start)) && (d == end || d.Before(end)) {
			r = append(r, s.events[i])
		}
	}
	return r, nil
}
