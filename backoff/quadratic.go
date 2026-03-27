package backoff

import (
	"context"
	"math"
	"time"
)

// Quadratic returns a backoff function that waits for "delta * incremental(attempt)" time between calls.
// This implements quadratic backoff using arithmetic progression: 1, 2, 4, 7, 11, 16...
// The growth rate is between linear and Fibonacci, providing moderate backoff behavior.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing quadratic growth
//
// Example:
//
//	// Quadratic backoff: 1s, 2s, 4s, 7s, 11s...
//	backoff := Quadratic(time.Second)
func Quadratic(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return incremental(delta, attempt)
	}
}

// incremental calculates quadratic backoff using arithmetic progression.
// Sequence: 1, 2, 4, 7, 11, 16... (triangular numbers)
// Formula: sum(1..attempt) = attempt*(attempt+1)/2
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: delta multiplied by the triangular number of attempt
func incremental(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(attempt*(attempt+1)/2)
}

// SubQuadratic returns a backoff function that implements sub-quadratic (1.5 power) growth.
// This waits for "delta * attempt^1.5" time between calls.
// Provides growth rate between linear and quadratic, offering balanced retry delays.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing sub-quadratic growth
//
// Example:
//
//	// Sub-quadratic backoff: 1s, 2.8s, 5.2s, 8s, 11.2s...
//	backoff := SubQuadratic(time.Second)
func SubQuadratic(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return subQuadratic(delta, attempt)
	}
}

// subQuadratic calculates "attempt^1.5" duration.
// Implements power-law backoff with exponent 1.5.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: delta multiplied by attempt raised to power 1.5
func subQuadratic(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(math.Pow(float64(attempt), 1.5))
}
