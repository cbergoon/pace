package main

import "time"

const DISPLAY_DATE_FORMAT string = "2006-01-02 15:04:05"

func beginningOfDay() time.Time {
	t := time.Now()
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func beginningOfWeek() time.Time {
	t := time.Now()
	weekday := time.Duration(t.Weekday())
	year, month, day := t.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return currentZeroDay.Add(-1 * (weekday) * 24 * time.Hour)
}

func beginningOfMonth() time.Time {
	t := time.Now()
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func maximumTime(times ...time.Time) time.Time {
	if len(times) <= 0 {
		return time.Time{}
	}
	maxTime := times[0]
	for _, t := range times {
		if t.After(maxTime) {
			maxTime = t
		}
	}
	return maxTime
}
