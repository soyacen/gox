package backoff

import (
	"context"
	"time"
)

// Func returns the backoff duration between call retries.
// The context.Context can be used to extract context values.
//
// Parameters:
//   - ctx: Context that may contain backoff configuration
//   - attempt: The retry attempt number (0-indexed)
//
// Returns:
//   - time.Duration: The duration to wait before the next retry
type Func func(ctx context.Context, attempt uint) time.Duration

// Factory returns a Func.
// It creates a backoff function based on a delta parameter.
// See: ConstantFactory, ExponentialFactory, Exponential2Factory, FibonacciFactory, LinearFactory, JitterUpFactory.
//
// Parameters:
//   - delta: Base duration or scaling factor for the backoff strategy
//
// Returns:
//   - Func: A backoff function configured with the given delta
type Factory func(delta time.Duration) Func
