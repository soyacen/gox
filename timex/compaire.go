package timex

import "time"

// Max returns the maximum time.Time value from a set of times.
//
// Parameters:
//   - x: the first time value.
//   - y: additional time values to compare.
//
// Returns:
//   - time.Time: the latest time among all provided values.
func Max(x time.Time, y ...time.Time) time.Time {
	r := x
	for _, t := range y {
		if t.After(r) {
			r = t
		}
	}
	return r
}

// Min returns the smallest time.Time value from a set of times.
//
// Parameters:
//   - x: the first time value.
//   - y: additional time values to compare.
//
// Returns:
//   - time.Time: the earliest time among all provided values.
func Min(x time.Time, y ...time.Time) time.Time {
	r := x
	for _, t := range y {
		if t.Before(r) {
			r = t
		}
	}
	return r
}
