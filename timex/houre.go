package timex

import "time"

// Hour returns a new time.Time with the minute, second, and nanosecond
// set to zero, preserving only the hour portion of t.
//
// Parameters:
//   - t: the input time.
//
// Returns:
//   - time.Time: the hour-truncated representation of t.
func Hour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

// ThisHour returns the hour-truncated representation of the current hour.
// If a time is provided, it uses that time; otherwise, it uses time.Now().
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the current hour with minutes and below truncated.
func ThisHour(ts ...time.Time) time.Time {
	if len(ts) <= 0 {
		return Hour(time.Now())
	}
	return Hour(ts[0])
}

// LastHour returns the hour-truncated representation of the previous hour.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the previous hour with minutes and below truncated.
func LastHour(ts ...time.Time) time.Time {
	return ThisHour(ts...).Add(-time.Hour)
}

// NextHour returns the hour-truncated representation of the next hour.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the next hour with minutes and below truncated.
func NextHour(ts ...time.Time) time.Time {
	return ThisHour(ts...).Add(time.Hour)
}
