// Package grpool implements the Gofer interface using the grpool library.
package grpool

import (
	"context"

	"github.com/ivpusic/grpool"
	"github.com/soyacen/gox/conc/gofer"
)

// Ensure Gofer implements the gofer.Gofer interface.
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer is an asynchronous task executor implementation based on grpool.
type Gofer struct {
	// Pool is the underlying grpool worker pool.
	Pool *grpool.Pool
}

// Go submits a task to the grpool job queue.
//
// Parameters:
//   - f: the task function to execute.
//
// Returns:
//   - error: always nil as sending to the channel does not return an error.
func (g *Gofer) Go(f func()) error {
	g.Pool.JobQueue <- f
	return nil
}

// Close releases the grpool resources.
//
// Parameters:
//   - ctx: the context used to control shutdown timeout (not used in the current implementation).
//
// Returns:
//   - error: always nil as grpool.Release does not return an error.
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.Release()
	return nil
}
