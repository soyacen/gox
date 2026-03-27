package backoff

import (
	"context"
	"math"
	"time"
)

// Exponential returns a backoff function that waits for "delta * e^attempt" time between calls.
// This implements exponential backoff using Euler's number (e ≈ 2.718), providing rapid growth in wait times.
// Suitable for scenarios where aggressive backoff is needed to reduce system load.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing exponential growth
//
// Example:
//
//	// Exponential backoff starting at 1 second: 1s, 2.7s, 7.4s, 20.1s...
//	backoff := Exponential(time.Second)
func Exponential(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return exponential(delta, attempt)
	}
}

// exponential calculates "delta * e^attempt" duration.
// This is the core implementation of exponential backoff using natural exponentiation.
// Returns 0 when attempt is 0.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: The calculated backoff duration
func exponential(delta time.Duration, attempt uint) time.Duration {
	if attempt == 0 {
		return 0
	}
	return delta * time.Duration(math.Exp(float64(attempt)))
}

// ExponentialFactory returns a factory that creates Exponential backoff functions.
// The factory takes a delta parameter and returns an Exponential backoff function.
//
// Returns:
//   - Factory: A factory function that creates Exponential backoff functions
//
// Example:
//
//	factory := ExponentialFactory()
//	backoff := factory(500 * time.Millisecond) // Creates exponential backoff with 500ms base
func ExponentialFactory() Factory {
	return func(delta time.Duration) Func {
		return Exponential(delta)
	}
}
