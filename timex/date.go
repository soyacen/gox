package timex

import "time"

// Date returns a new time.Time with the hour, minute, second, and nanosecond
// set to zero, preserving only the date portion of t.
//
// Parameters:
//   - t: the input time.
//
// Returns:
//   - time.Time: the date-only representation of t.
func Date(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Today returns the date-only representation of the current day.
// If a time is provided, it uses that time; otherwise, it uses time.Now().
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: today's date with time truncated to midnight.
func Today(ts ...time.Time) time.Time {
	if len(ts) <= 0 {
		return Date(time.Now())
	}
	if ts[0].IsZero() {
		return Date(time.Now())
	}
	return Date(ts[0])
}

// Tomorrow returns the date-only representation of the day after the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: tomorrow's date with time truncated to midnight.
func Tomorrow(ts ...time.Time) time.Time {
	return Today(ts...).AddDate(0, 0, 1)
}

// Yesterday returns the date-only representation of the day before the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: yesterday's date with time truncated to midnight.
func Yesterday(ts ...time.Time) time.Time {
	return Today(ts...).AddDate(0, 0, -1)
}
