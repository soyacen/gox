package timex

import "time"

// Year returns the first day of the year for the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: January 1st of the year at midnight.
func Year(ts ...time.Time) time.Time {
	if len(ts) <= 0 {
		return year(time.Now())
	}
	if ts[0].IsZero() {
		return year(time.Now())
	}
	return year(ts[0])
}

func year(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}
