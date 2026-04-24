package timex

import (
	"time"
)

func month(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// ThisMonth returns the first day of the month for the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the first day of the month at midnight.
func ThisMonth(ts ...time.Time) time.Time {
	if len(ts) <= 0 {
		return month(time.Now())
	}
	return month(ts[0])
}

// LastMonth returns the first day of the previous month for the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the first day of the previous month at midnight.
func LastMonth(ts ...time.Time) time.Time {
	return ThisMonth(ts...).AddDate(0, -1, 0)
}

// NextMonth returns the first day of the next month for the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the first day of the next month at midnight.
func NextMonth(ts ...time.Time) time.Time {
	return ThisMonth(ts...).AddDate(0, 1, 0)
}

// MonthOfLastDay returns the last day of the month for the given time.
// If no time is provided, it uses the current time.
//
// Parameters:
//   - ts: optional time value (defaults to time.Now()).
//
// Returns:
//   - time.Time: the last day of the month at midnight.
func MonthOfLastDay(ts ...time.Time) time.Time {
	return NextMonth(ts...).AddDate(0, 0, -1)
}
