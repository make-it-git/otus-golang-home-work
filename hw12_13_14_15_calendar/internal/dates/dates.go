package dates

import (
	"time"
)

func makeStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
}

func makeEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func DayRange(date time.Time) (time.Time, time.Time) {
	return makeStart(date), makeEnd(date)
}

func MonthRange(date time.Time) (time.Time, time.Time) {
	start := date.AddDate(0, 0, -date.Day()+1)
	end := date.AddDate(0, 1, -date.Day())
	return makeStart(start), makeEnd(end)
}

func WeekRange(date time.Time) (time.Time, time.Time) {
	w := int(date.Weekday())
	start := date.AddDate(0, 0, -w+1)
	end := date.AddDate(0, 0, 7-w)
	return makeStart(start), makeEnd(end)
}
