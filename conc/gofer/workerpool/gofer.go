// Package workerpool implements the Gofer interface using the workerpool library.
package workerpool

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/soyacen/gox/conc/gofer"
)

// Ensure Gofer implements the gofer.Gofer interface.
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer is an asynchronous task executor implementation based on workerpool.
type Gofer struct {
	// Pool is the underlying worker pool.
	Pool *workerpool.WorkerPool
}

// Go submits a task to the worker pool for execution.
//
// Parameters:
//   - f: the task function to execute.
//
// Returns:
//   - error: always nil as workerpool.Submit does not return an error.
func (g *Gofer) Go(f func()) error {
	g.Pool.Submit(f)
	return nil
}

// Close shuts down the worker pool and waits for all tasks to complete.
//
// Parameters:
//   - ctx: the context used to control shutdown timeout (not used in the current implementation).
//
// Returns:
//   - error: always nil as workerpool.StopWait does not return an error.
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.StopWait()
	return nil
}
