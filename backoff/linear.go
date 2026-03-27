package backoff

import (
	"context"
	"math"
	"time"
)

// Linear returns a backoff function that waits for "delta * attempt" time between calls.
// This implements linear backoff where wait times grow proportionally with the attempt number.
// Provides steady, predictable growth in retry delays.
//
// Parameters:
//   - delta: The base duration multiplier (time per attempt)
//
// Returns:
//   - Func: A backoff function implementing linear growth
//
// Example:
//
//	// Linear backoff: 0s, 1s, 2s, 3s, 4s...
//	backoff := Linear(time.Second)
func Linear(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linear(delta, attempt)
	}
}

// linear calculates "delta * attempt" duration.
// Simple linear growth where each attempt adds one more delta unit.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: delta multiplied by attempt number
func linear(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(attempt)
}

// Linear15 returns a backoff function that implements 1.5x linear backoff.
// This waits for "delta * attempt * 1.5" time between calls.
// Provides slightly faster growth than standard linear backoff.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing 1.5x linear growth
//
// Example:
//
//	// 1.5x linear backoff: 0s, 1.5s, 3s, 4.5s, 6s...
//	backoff := Linear15(time.Second)
func Linear15(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linear15(delta, attempt)
	}
}

// linear15 calculates "delta * attempt * 1.5" duration.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: delta multiplied by attempt and 1.5
func linear15(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(float64(attempt)*1.5)
}

// LinearSqrt returns a backoff function that implements linear-square root backoff.
// This waits for "delta * attempt * sqrt(attempt)" time between calls.
// Provides growth rate between linear and quadratic, offering balanced retry behavior.
// This is considered one of the most optimal backoff strategies for many scenarios.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing linear-square root growth
//
// Example:
//
//	// Linear-sqrt backoff: 0s, 1s, 2.8s, 5.2s, 7.9s...
//	backoff := LinearSqrt(time.Second)
func LinearSqrt(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linearSqrt(delta, attempt)
	}
}

// linearSqrt calculates "delta * attempt * sqrt(attempt)" duration.
// Combines linear and square root growth for moderate backoff behavior.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: Calculated backoff using linear-sqrt formula
func linearSqrt(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(float64(attempt)*math.Sqrt(float64(attempt)))
}

// LinearPlus returns a backoff function that implements linear plus half-linear backoff.
// This waits for "delta * (attempt + attempt/2)" time between calls.
// Equivalent to "delta * attempt * 1.5", similar to Linear15 but using integer arithmetic.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing linear-plus growth
//
// Example:
//
//	// Linear-plus backoff: 0s, 1.5s, 3s, 4.5s, 6s...
//	backoff := LinearPlus(time.Second)
func LinearPlus(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return linearPlus(delta, attempt)
	}
}

// linearPlus calculates "delta * (attempt + attempt/2)" duration.
// Uses integer arithmetic to achieve 1.5x linear growth.
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number
//
// Returns:
//   - time.Duration: Calculated backoff using linear-plus formula
func linearPlus(delta time.Duration, attempt uint) time.Duration {
	return delta * time.Duration(int64(attempt)+int64(attempt)/2)
}

// LinearFactory returns a factory that creates Linear backoff functions.
// The factory takes a delta parameter and returns a Linear backoff function.
//
// Returns:
//   - Factory: A factory function that creates Linear backoff functions
//
// Example:
//
//	factory := LinearFactory()
//	backoff := factory(500 * time.Millisecond) // Creates linear backoff with 500ms increment
func LinearFactory() Factory {
	return func(delta time.Duration) Func {
		return Linear(delta)
	}
}
