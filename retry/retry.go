package retry

import (
	"context"
	"time"

	"github.com/soyacen/gox/backoff"
)

// Strategy defines the interface for configuring retry behavior with backoff strategies.
// It provides methods to set backoff functions, retry conditions, and execute commands with retries.
type Strategy interface {

	// Backoff sets the backoff function for retry intervals.
	Backoff(backoffFunc backoff.Func) Strategy

	// RetryOn sets the condition function to determine whether to retry on error.
	RetryOn(retryOnFunc func(err error) bool) Strategy

	// Exec executes a command with retry logic, accepting context and current attempt number.
	Exec(ctx context.Context, cmd func(ctx context.Context, attempt uint) error) error
}

// defaultStrategy implements the Strategy interface with configurable max attempts,
// backoff function, and retry condition function.
type defaultStrategy struct {
	maxAttempts uint
	backoffFunc backoff.Func
	retryOnFunc func(err error) bool
}

func (r *defaultStrategy) Backoff(backoffFunc backoff.Func) Strategy {
	r.backoffFunc = backoffFunc
	return r
}

func (r *defaultStrategy) RetryOn(retryOnFunc func(err error) bool) Strategy {
	r.retryOnFunc = retryOnFunc
	return r
}

func (r *defaultStrategy) Exec(ctx context.Context, cmd func(ctx context.Context, attempt uint) error) error {
	var attempt uint
	for attempt < r.maxAttempts {
		// execute cmd
		err := cmd(ctx, attempt)
		if err == nil {
			// return if err is nil.
			return nil
		}
		if r.retryOnFunc != nil && !r.retryOnFunc(err) {
			// return if retryOnFunc is not nil and retryOnFunc returns false
			return err
		}
		// increase the number of attempts
		attempt++
		select {
		case <-ctx.Done(): // return if context is done, return context error
			return ctx.Err()
		case <-time.After(r.backoffFunc(ctx, attempt)): // sleep and wait retry
			continue
		}
	}
	// perform the execution
	return cmd(ctx, attempt)
}

// MaxAttempts creates a new defaultStrategy instance with the specified maximum retry attempts.
//
// Parameters:
//   - maxAttempts: The maximum number of retry attempts
//
// Returns:
//   - Strategy: A strategy configured with the given max attempts
func MaxAttempts(maxAttempts uint) Strategy {
	return &defaultStrategy{
		maxAttempts: maxAttempts,
		backoffFunc: backoff.Zero(),
		retryOnFunc: func(err error) bool {
			return true
		},
	}
}

// Call executes a command with context, max attempts, and backoff function.
// Deprecated: Do not use. Use MaxAttempts instead.
//
// Parameters:
//   - ctx: Context for cancellation
//   - maxAttempts: Maximum number of retry attempts
//   - backoffFunc: Backoff function for retry intervals
//   - method: Function to execute
//
// Returns:
//   - error: Error from execution or nil if successful
func Call(ctx context.Context, maxAttempts uint, backoffFunc backoff.Func, method func(attemptTime int) error) error {
	return MaxAttempts(maxAttempts).Backoff(backoffFunc).Exec(ctx, func(ctx context.Context, attempt uint) error {
		return method(int(attempt))
	})
}
