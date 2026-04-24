// Package brave provides utilities for handling panics in Go programs,
// allowing recovery from panics and converting them to errors or custom handling.
package brave

import (
	"fmt"
	"runtime/debug"
)

// Do executes the given function f and captures any panics that may occur during execution.
// If a panic occurs, it optionally calls the first recovery function from rs to handle the panic.
// If no recovery function is provided, it defaults to logging the panic information.
//
// Parameters:
//   - f: the function to execute
//   - rs: optional recovery functions to handle panics
func Do(f func(), rs ...func(p any, stack []byte)) {
	defer func() {
		if p := recover(); p != nil {
			var r func(p any, stack []byte)
			// Use the first provided recovery function if available
			if len(rs) > 0 {
				r = rs[0]
			}
			// Default recovery function if none provided
			if r == nil {
				r = func(p any, stack []byte) {
					fmt.Printf("brave: panic triggered: %v\nstack: %s\n", p, stack)
				}
			}
			// Execute the recovery function
			r(p, debug.Stack())
		}
	}()
	f()
}

// DoE is functionally equivalent to Do but returns an error.
// It executes the given function f and captures any panics, converting them to errors.
//
// Parameters:
//   - f: the function to execute that may return an error
//   - rs: optional recovery functions to handle panics and return errors
//
// Returns:
//   - error: the error from f or a panic-converted error
func DoE(f func() error, rs ...func(p any, stack []byte) error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			var r func(p any, stack []byte) error
			// Use the first provided recovery function if available
			if len(rs) > 0 {
				r = rs[0]
			}
			// Default recovery function if none provided
			if r == nil {
				r = func(p any, stack []byte) error {
					return fmt.Errorf("brave: panic triggered: %v, stack: %s", p, stack)
				}
			}
			// Set the error from the recovery function
			err = r(p, debug.Stack())
		}
	}()
	return f()
}

// DoRE is functionally equivalent to DoE but returns both a result and an error.
// It executes the given function f and captures any panics, converting them to errors.
//
// Type Parameters:
//   - R: the type of the result
//
// Parameters:
//   - f: the function to execute that returns a result and an error
//   - rs: optional recovery functions to handle panics and return errors
//
// Returns:
//   - R: the result from f
//   - error: the error from f or a panic-converted error
func DoRE[R any](f func() (R, error), rs ...func(p any, stack []byte) error) (_ R, err error) {
	defer func() {
		if p := recover(); p != nil {
			var r func(p any, stack []byte) error
			// Use the first provided recovery function if available
			if len(rs) > 0 {
				r = rs[0]
			}
			// Default recovery function if none provided
			if r == nil {
				r = func(p any, stack []byte) error {
					return fmt.Errorf("brave: panic triggered: %v, stack: %s", p, stack)
				}
			}
			// Set the error from the recovery function
			err = r(p, debug.Stack())
		}
	}()
	return f()
}

// Go asynchronously executes the Do function in a new goroutine.
//
// Parameters:
//   - f: the function to execute
//   - rs: optional recovery functions to handle panics
func Go(f func(), rs ...func(p any, stack []byte)) {
	go Do(f, rs...)
}

// GoE asynchronously executes the DoE function in a new goroutine and returns an error channel.
//
// Parameters:
//   - f: the function to execute that may return an error
//   - rs: optional recovery functions to handle panics and return errors
//
// Returns:
//   - <-chan error: a channel that will receive the error if any
func GoE(f func() error, rs ...func(p any, stack []byte) error) <-chan error {
	// Create error channel with buffer of 1
	errC := make(chan error, 1)
	go func() {
		// Close channel when goroutine exits
		defer close(errC)
		// Execute function with panic handling
		err := DoE(f, rs...)
		// Send error to channel if one occurred
		if err != nil {
			errC <- err
		}
	}()
	return errC
}

// GoRE asynchronously executes the DoRE function in a new goroutine and returns result and error channels.
//
// Type Parameters:
//   - R: the type of the result
//
// Parameters:
//   - f: the function to execute that returns a result and an error
//   - rs: optional recovery functions to handle panics and return errors
//
// Returns:
//   - <-chan R: a channel that will receive the result if successful
//   - <-chan error: a channel that will receive the error if any
func GoRE[R any](f func() (R, error), rs ...func(p any, stack []byte) error) (<-chan R, <-chan error) {
	// Create result and error channels with buffer of 1
	retC := make(chan R, 1)
	errC := make(chan error, 1)
	go func() {
		// Close channels when goroutine exits
		defer close(errC)
		defer close(retC)
		// Execute function with panic handling
		ret, err := DoRE[R](f, rs...)
		// Send error to channel if one occurred
		if err != nil {
			errC <- err
			return
		}
		// Send result to channel if no error
		retC <- ret
	}()
	return retC, errC
}
