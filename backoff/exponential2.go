package backoff

import (
	"context"
	"math"
	"time"
)

// Exponential2 returns a backoff function that waits for "delta * 2^attempt" time between calls.
// This implements binary exponential backoff, doubling the wait time with each retry attempt.
// Commonly used in network protocols and distributed systems for congestion control.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing binary exponential growth
//
// Example:
//
//	// Binary exponential backoff: 1s, 2s, 4s, 8s, 16s...
//	backoff := Exponential2(time.Second)
func Exponential2(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return exponential2(delta, attempt)
	}
}

// exponential2 calculates "delta * 2^attempt" duration.
// This is the core implementation of binary exponential backoff.
// Returns 0 when attempt is 0.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: The calculated backoff duration
func exponential2(delta time.Duration, attempt uint) time.Duration {
	if attempt == 0 {
		return 0
	}
	return delta * time.Duration(math.Exp2(float64(attempt)))
}

// Exponential2Factory returns a factory that creates Exponential2 backoff functions.
// The factory takes a delta parameter and returns an Exponential2 backoff function.
//
// Returns:
//   - Factory: A factory function that creates Exponential2 backoff functions
//
// Example:
//
//	factory := Exponential2Factory()
//	backoff := factory(1 * time.Second) // Creates binary exponential backoff with 1s base
func Exponential2Factory() Factory {
	return func(delta time.Duration) Func {
		return Exponential2(delta)
	}
}
