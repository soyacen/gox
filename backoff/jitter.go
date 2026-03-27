package backoff

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/soyacen/gox/randx"
)

// JitterUp returns a backoff function that adds random jitter to the interval.
// This adds or subtracts time from the interval within a given jitter fraction,
// helping to prevent thundering herd problems by distributing retry attempts.
//
// For example, for 10s interval and jitter 0.1, it will return a time within [9s, 11s].
//
// Parameters:
//   - backoff: The base backoff function to add jitter to
//   - jitter: The jitter fraction (e.g., 0.1 means ±10% variation)
//
// Returns:
//   - Func: A backoff function with randomized jitter applied
//
// Example:
//
//	// Add 10% jitter to exponential backoff
//	backoff := JitterUp(Exponential(time.Second), 0.1)
func JitterUp(backoff Func, jitter float64) Func {
	r, err := randx.NewChaCha8()
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, attempt uint) time.Duration {
		interval := backoff(ctx, attempt)
		return jitterUp(r, interval, jitter)
	}
}

// jitterUp applies jitter to a single interval value.
// It randomly adjusts the interval by the specified jitter fraction.
//
// Parameters:
//   - r: Random number generator
//   - interval: The base interval duration
//   - jitter: The jitter fraction (e.g., 0.1 means ±10%)
//
// Returns:
//   - time.Duration: The adjusted interval with jitter applied
func jitterUp(r *rand.Rand, interval time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (r.Float64()*2 - 1)
	return time.Duration(float64(interval) * (1 + multiplier))
}

// JitterUpFactory returns a factory that creates JitterUp backoff functions.
// This wraps another factory to automatically add jitter to its backoff functions.
//
// Parameters:
//   - factory: The base backoff factory to wrap
//   - jitter: The jitter fraction to apply (e.g., 0.1 means ±10% variation)
//
// Returns:
//   - Factory: A factory function that creates backoff functions with jitter
//
// Example:
//
//	// Create exponential backoff with 10% jitter
//	factory := JitterUpFactory(ExponentialFactory(), 0.1)
//	backoff := factory(time.Second)
func JitterUpFactory(factory Factory, jitter float64) Factory {
	return func(delta time.Duration) Func {
		return JitterUp(factory(delta), jitter)
	}
}
