package timex

import "time"

// UnixZero represents the Unix epoch (1970-01-01 00:00:00 UTC).
var UnixZero = time.Unix(0, 0)

// IsUnixZero reports whether t is exactly the Unix epoch.
//
// Parameters:
//   - t: the time to check.
//
// Returns:
//   - bool: true if t equals UnixZero.
func IsUnixZero(t time.Time) bool {
	return t.Equal(UnixZero)
}

// IsZero reports whether t is either the zero time or the Unix epoch.
//
// Parameters:
//   - t: the time to check.
//
// Returns:
//   - bool: true if t is zero or UnixZero.
func IsZero(t time.Time) bool {
	return t.IsZero() || IsUnixZero(t)
}

// IsNotZero reports whether t is neither the zero time nor the Unix epoch.
//
// Parameters:
//   - t: the time to check.
//
// Returns:
//   - bool: true if t is not zero and not UnixZero.
func IsNotZero(t time.Time) bool {
	return !IsZero(t)
}
