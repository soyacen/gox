// Package gofer provides a simple asynchronous task executor interface.
package gofer

import "context"

// Gofer defines the interface for an asynchronous task executor.
type Gofer interface {
	// Go starts an asynchronous task.
	//
	// Parameters:
	//   - f: the task function to execute.
	//
	// Returns:
	//   - error: an error if the task fails to start.
	Go(f func()) error

	// Close shuts down the executor and waits for all tasks to complete.
	//
	// Parameters:
	//   - ctx: the context used to control shutdown timeout.
	//
	// Returns:
	//   - error: an error if the shutdown process encounters any issues.
	Close(ctx context.Context) error
}