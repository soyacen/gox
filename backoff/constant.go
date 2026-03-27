package backoff

import (
	"context"
	"time"
)

// Constant returns a backoff function that waits for a fixed period of time between calls.
// This is useful when you want consistent retry intervals regardless of the attempt number.
// Returns 0 when attempt is 0.
//
// Parameters:
//   - delta: The fixed duration to wait between each retry attempt
//
// Returns:
//   - Func: A backoff function that always returns the same duration (or 0 for attempt 0)
//
// Example:
//
//	// Wait 5 seconds between each retry
//	backoff := Constant(5 * time.Second)
func Constant(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		if attempt == 0 {
			return 0
		}
		return delta
	}
}

// Zero returns a backoff function that waits for zero time between calls.
// This effectively means retries happen immediately without any delay.
// Useful for scenarios where immediate retry is preferred over waiting.
//
// Returns:
//   - Func: A backoff function that returns 0 duration
//
// Example:
//
//	// Retry immediately without waiting
//	backoff := Zero()
func Zero() Func {
	return Constant(0)
}

// ConstantFactory returns a factory that creates Constant backoff functions.
// The factory takes a delta parameter and returns a Constant backoff function.
//
// Returns:
//   - Factory: A factory function that creates Constant backoff functions
//
// Example:
//
//	factory := ConstantFactory()
//	backoff := factory(5 * time.Second) // Creates a 5-second constant backoff
func ConstantFactory() Factory {
	return func(delta time.Duration) Func {
		return Constant(delta)
	}
}
