package backoff

import (
	"context"
	"time"
)

// Fibonacci returns a backoff function that waits for "delta * fibonacci(attempt)" time between calls.
// This implements Fibonacci backoff where wait times grow according to the Fibonacci sequence.
// Provides moderate growth rate between linear and exponential backoff.
//
// Parameters:
//   - delta: The base duration multiplier
//
// Returns:
//   - Func: A backoff function implementing Fibonacci sequence growth
//
// Example:
//
//	// Fibonacci backoff: 1s, 1s, 2s, 3s, 5s, 8s...
//	backoff := Fibonacci(time.Second)
func Fibonacci(delta time.Duration) Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		return fibonacci(delta, attempt)
	}
}

// fibonacci calculates the Fibonacci sequence value at the given attempt.
// Sequence: 0, 1, 1, 2, 3, 5, 8, 13...
//
// Parameters:
//   - delta: Base duration multiplier
//   - attempt: The retry attempt number (index in Fibonacci sequence)
//
// Returns:
//   - time.Duration: delta multiplied by the attempt-th Fibonacci number
func fibonacci(delta time.Duration, attempt uint) time.Duration {
	var (
		pre int64
		cur int64
		i   uint
	)
	for pre, cur, i = 0, 1, 0; i < attempt; i++ {
		pre, cur = cur, pre+cur
	}
	return delta * time.Duration(pre)
}

// FibonacciFactory returns a factory that creates Fibonacci backoff functions.
// The factory takes a delta parameter and returns a Fibonacci backoff function.
//
// Returns:
//   - Factory: A factory function that creates Fibonacci backoff functions
//
// Example:
//
//	factory := FibonacciFactory()
//	backoff := factory(500 * time.Millisecond) // Creates Fibonacci backoff with 500ms base
func FibonacciFactory() Factory {
	return func(delta time.Duration) Func {
		return Fibonacci(delta)
	}
}
