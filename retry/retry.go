package retry

import (
	"context"
	"time"

	"github.com/soyacen/gox/backoff"
)

// BackoffStrategy defines the interface for configuring backoff behavior.
// It provides methods to set backoff functions and chain to retry conditions.
type BackoffStrategy interface {
	// Backoff sets the backoff function for retry intervals.
	//
	// Parameters:
	//   - backoffFunc: Function that calculates backoff duration based on attempt number
	//
	// Returns:
	//   - RetryOnStrategy: Strategy configured to set retry conditions
	Backoff(backoffFunc backoff.Func) RetryOnStrategy
}

// RetryOnStrategy defines the interface for configuring retry conditions.
// It provides methods to set error-based retry logic and execute commands.
type RetryOnStrategy interface {
	// RetryOn sets the condition function to determine whether to retry on error.
	//
	// Parameters:
	//   - retryOnFunc: Function that returns true if the error should trigger a retry
	//
	// Returns:
	//   - Executor: Strategy configured to execute commands with retry logic
	RetryOn(retryOnFunc func(err error) bool) Executor
}

// Executor defines the interface for executing commands with retry logic.
// It executes a command and handles retries based on configured strategies.
type Executor interface {
	// Exec executes a command with retry logic, accepting context and current attempt number.
	// The command will be retried based on the configured backoff and retry conditions.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control
	//   - cmd: Function to execute, receiving context and attempt number
	//
	// Returns:
	//   - error: nil if successful, or the last error after all retries exhausted
	Exec(ctx context.Context, cmd func(ctx context.Context, attempt uint) error) error
}

// defaultStrategy implements BackoffStrategy, RetryOnStrategy, and Executor interfaces.
// It holds configuration for maximum attempts, backoff function, and retry condition.
type defaultStrategy struct {
	maxAttempts uint
	backoffFunc backoff.Func
	retryOnFunc func(err error) bool
}

// Backoff sets the backoff function for this strategy.
//
// Parameters:
//   - backoffFunc: Function that calculates backoff duration
//
// Returns:
//   - RetryOnStrategy: Self for method chaining
func (r *defaultStrategy) Backoff(backoffFunc backoff.Func) RetryOnStrategy {
	r.backoffFunc = backoffFunc
	return r
}

// RetryOn sets the condition function to determine whether to retry on error.
//
// Parameters:
//   - retryOnFunc: Function that returns true if retry should occur
//
// Returns:
//   - Executor: Self for method chaining
func (r *defaultStrategy) RetryOn(retryOnFunc func(err error) bool) Executor {
	r.retryOnFunc = retryOnFunc
	return r
}

// Exec executes a command with retry logic.
// It attempts to execute the command and retries based on configured backoff and retry conditions.
// The loop continues until success, max attempts reached, context cancelled, or retry condition fails.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - cmd: Function to execute, receives context and current attempt number
//
// Returns:
//   - error: nil if successful, context error if cancelled, or command error after all retries
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
//   - maxAttempts: The maximum number of retry attempts (must be > 0)
//
// Returns:
//   - BackoffStrategy: A strategy configured with the given max attempts
func MaxAttempts(maxAttempts uint) BackoffStrategy {
	return &defaultStrategy{
		maxAttempts: maxAttempts,
	}
}
