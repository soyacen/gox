// Package tunny implements the Gofer interface using the tunny library.
package tunny

import (
	"context"

	"github.com/Jeffail/tunny"
	"github.com/soyacen/gox/conc/gofer"
)

// Ensure Gofer implements the gofer.Gofer interface.
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer is an asynchronous task executor implementation based on tunny pool.
type Gofer struct {
	// Pool is the underlying tunny worker pool.
	Pool *tunny.Pool
}

// Go submits a task to the tunny pool for execution.
//
// Parameters:
//   - f: the task function to execute.
//
// Returns:
//   - error: always nil as tunny.Pool.Process does not return an error.
func (g *Gofer) Go(f func()) error {
	g.Pool.Process(f)
	return nil
}

// Close shuts down the tunny pool.
//
// Parameters:
//   - ctx: the context used to control shutdown timeout (not used in the current implementation).
//
// Returns:
//   - error: always nil as tunny.Pool.Close does not return an error.
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.Close()
	return nil
}
