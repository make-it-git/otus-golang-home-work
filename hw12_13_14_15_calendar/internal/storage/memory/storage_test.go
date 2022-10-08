package memorystorage

import (
	"context"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logic"
	"testing"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func ymd(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t
}

func event(id string, title string) storage.Event {
	return storage.Event{
		ID:               id,
		Title:            title,
		StartTime:        time.Time{},
		Duration:         0,
		Description:      nil,
		OwnerID:          0,
		NotificationTime: nil,
	}
}

func dateEvent(id string, date string) storage.Event {
	e := event(id, "event for "+date)
	e.StartTime = ymd(date)
	return e
}

func TestStorage(t *testing.T) {
	t.Run("Should store event", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), event("1", "test"))
		require.NoError(t, err)
	})
	t.Run("Should fail store event missing id", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), event("", ""))
		require.ErrorIs(t, err, logic.ErrMissingID)
	})
	t.Run("Should fail store event duplicate id", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), event("1", "test 1"))
		require.NoError(t, err)
		err = s.Create(context.Background(), event("1", "test 2"))
		require.ErrorIs(t, err, logic.ErrDuplicateID)
	})
	t.Run("Should update event", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), event("1", "test 1"))
		require.NoError(t, err)
		require.Equal(t, "test 1", s.events[0].Title)
		err = s.Update(context.Background(), event("1", "test 2"))
		require.NoError(t, err)
		require.Equal(t, "test 2", s.events[0].Title)
	})
	t.Run("Should delete event", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), event("1", "test 1"))
		require.NoError(t, err)
		require.Equal(t, 1, len(s.events))
		err = s.Delete(context.Background(), "1")
		require.NoError(t, err)
		require.Equal(t, 0, len(s.events))
	})
	t.Run("Should fail delete event", func(t *testing.T) {
		s := New()
		err := s.Delete(context.Background(), "1")
		require.ErrorIs(t, err, logic.ErrNotFoundID)
	})
	t.Run("Should fail update event", func(t *testing.T) {
		s := New()
		err := s.Update(context.Background(), event("1", "test 1"))
		require.ErrorIs(t, err, logic.ErrNotFoundID)
	})
	t.Run("Should list day", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), dateEvent("1", "2010-12-01"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("2", "2010-12-02"))
		require.NoError(t, err)
		events, err := s.ListDay(context.Background(), ymd("2010-12-02"))
		require.NoError(t, err)
		require.Equal(t, 1, len(events))
		require.Equal(t, "event for 2010-12-02", events[0].Title)
	})
	t.Run("Should list week", func(t *testing.T) {
		// One week: monday=2010-08-02, sunday=2010-08-08
		// Next week: monday=2010-08-09
		s := New()
		err := s.Create(context.Background(), dateEvent("1", "2010-08-02"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("2", "2010-08-04"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("3", "2010-08-08"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("4", "2010-08-09"))
		require.NoError(t, err)
		events, err := s.ListWeek(context.Background(), ymd("2010-08-05"))
		require.NoError(t, err)
		require.Equal(t, 3, len(events))
		require.Equal(t, "event for 2010-08-02", events[0].Title)
		require.Equal(t, "event for 2010-08-04", events[1].Title)
		require.Equal(t, "event for 2010-08-08", events[2].Title)
	})
	t.Run("Should list month", func(t *testing.T) {
		s := New()
		err := s.Create(context.Background(), dateEvent("1", "2010-08-01"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("2", "2010-08-15"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("3", "2010-08-31"))
		require.NoError(t, err)
		err = s.Create(context.Background(), dateEvent("4", "2010-09-01"))
		require.NoError(t, err)
		events, err := s.ListMonth(context.Background(), ymd("2010-08-29"))
		require.NoError(t, err)
		require.Equal(t, 3, len(events))
		require.Equal(t, "event for 2010-08-01", events[0].Title)
		require.Equal(t, "event for 2010-08-15", events[1].Title)
		require.Equal(t, "event for 2010-08-31", events[2].Title)
	})
}
