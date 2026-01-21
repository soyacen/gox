// Package waiter provides utilities to convert blocking Wait operations into channel-based notifications.
// This allows for easier integration with select statements and other channel-based concurrency patterns.
package waiter

import (
	"context"
)

// WaitNotify converts a Wait() operation into a channel notification.
// It takes any object that implements the Wait() method (like sync.WaitGroup) and returns a channel
// that will be closed when the Wait() method completes.
// This is useful for integrating blocking Wait operations with select statements.
// waiter: any object that implements Wait() method
// Returns a channel that will be closed when waiter.Wait() completes
func WaitNotify(waiter interface{ Wait() }) <-chan struct{} {
	// Create an unbuffered channel to signal completion
	c := make(chan struct{})

	// Start a goroutine to perform the wait operation
	go func() {
		// Ensure the channel is closed when this goroutine exits
		defer close(c)
		// Block until the waiter completes
		waiter.Wait()
	}()

	return c
}

// WaitNotifyE converts a Wait() error operation into a channel notification.
// It takes any object that implements the Wait() error method (like errgroup.Group) and returns a channel
// that will receive any error when the Wait() method completes.
// If the Wait() method returns nil, the channel will be closed without sending any value.
// waiter: any object that implements Wait() error method
// Returns a buffered channel that will receive an error if waiter.Wait() returns one
func WaitNotifyE(waiter interface{ Wait() error }) <-chan error {
	// Create a buffered channel to send potential errors
	c := make(chan error, 1)

	// Start a goroutine to perform the wait operation
	go func() {
		// Ensure the channel is closed when this goroutine exits
		defer close(c)
		// Execute the wait operation and check for errors
		if err := waiter.Wait(); err != nil {
			// Send the error to the channel
			c <- err
		}
	}()

	return c
}

// WaitContentNotifyE converts a Wait(context.Context) error operation into a channel notification.
// It takes any object that implements the Wait(context.Context) error method and returns a channel
// that will receive any error when the Wait() method completes.
// If the Wait() method returns nil, the channel will be closed without sending any value.
// ctx: the context to pass to the waiter's Wait method
// waiter: any object that implements Wait(context.Context) error method
// Returns a buffered channel that will receive an error if waiter.Wait(ctx) returns one
func WaitContentNotifyE(ctx context.Context, waiter interface{ Wait(context.Context) error }) <-chan error {
	// Create a buffered channel to send potential errors
	c := make(chan error, 1)

	// Start a goroutine to perform the wait operation
	go func() {
		// Ensure the channel is closed when this goroutine exits
		defer close(c)
		// Execute the wait operation with context and check for errors
		if err := waiter.Wait(ctx); err != nil {
			// Send the error to the channel
			c <- err
		}
	}()

	return c
}
