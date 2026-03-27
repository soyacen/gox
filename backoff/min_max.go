package backoff

import (
	"context"
	"time"
)

// Min returns a new backoff function that ensures the backoff duration never falls below the specified minimum value.
// When the underlying backoff strategy returns an interval less than minDuration, minDuration is used as the actual wait time.
// This is useful for ensuring sufficient minimum intervals between retries to avoid excessively frequent retry attempts.
//
// Parameters:
//   - backoff: The base backoff strategy function
//   - minDuration: The minimum backoff duration (lower bound)
//
// Returns:
//   - Func: A wrapped backoff function that guarantees return value >= minDuration
//
// Example:
//
//	// Ensure backoff time is at least 100 milliseconds
//	minBackoff := Min(Exponential(), 100*time.Millisecond)
func Min(backoff Func, minDuration time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return max(backoff(ctx, attempt), minDuration)
	}
}

// Max returns a new backoff function that ensures the backoff duration never exceeds the specified maximum value.
// When the underlying backoff strategy returns an interval greater than maxDuration, maxDuration is used as the actual wait time.
// This is useful for setting an upper bound on retry wait times to prevent excessively long delays.
//
// Parameters:
//   - backoff: The base backoff strategy function
//   - maxDuration: The maximum backoff duration (upper bound)
//
// Returns:
//   - Func: A wrapped backoff function that guarantees return value <= maxDuration
//
// Example:
//
//	// Ensure backoff time does not exceed 30 seconds
//	maxBackoff := Max(Exponential(), 30*time.Second)
func Max(backoff Func, maxDuration time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return min(backoff(ctx, attempt), maxDuration)
	}
}
