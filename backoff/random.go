package backoff

import (
	"context"
	"github.com/soyacen/gox/randx"
	"time"
)

// Random returns a backoff function that waits for a random time in the range [alpha, beta) between calls.
// This generates uniformly distributed random backoff intervals within the specified bounds.
// Useful for scenarios where unpredictable retry timing is desired to avoid synchronization.
//
// Parameters:
//   - alpha: The minimum backoff duration (inclusive)
//   - beta: The maximum backoff duration (exclusive)
//
// Returns:
//   - Func: A backoff function that returns random durations in [alpha, beta)
//
// Example:
//
//	// Random backoff between 1 and 5 seconds
//	backoff := Random(time.Second, 5*time.Second)
func Random(alpha, beta time.Duration) Func {
	r, err := randx.NewChaCha8()
	if err != nil {
		panic(err)
	}
	return func(ctx context.Context, attempt uint) time.Duration {
		return alpha + time.Duration(r.Int64N(beta.Nanoseconds()-alpha.Nanoseconds()))
	}
}

// RandomFactory returns a factory that creates Random backoff functions.
// The factory ignores the delta parameter and always creates Random backoff with fixed alpha and beta.
//
// Parameters:
//   - alpha: The minimum backoff duration (inclusive)
//   - beta: The maximum backoff duration (exclusive)
//
// Returns:
//   - Factory: A factory function that creates Random backoff functions
//
// Example:
//
//	factory := RandomFactory(time.Second, 5*time.Second)
//	backoff := factory(0) // Creates random backoff between 1s and 5s
func RandomFactory(alpha, beta time.Duration) Factory {
	return func(delta time.Duration) Func {
		return func(ctx context.Context, attempt uint) time.Duration {
			return Random(alpha, beta)(ctx, attempt)
		}
	}
}
