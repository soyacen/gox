package contextx

import (
	"context"
	"fmt"
)

// Error represents an error that occurred in a context.
// It contains both the context error and the underlying cause error.
type Error struct {
	// CtxErr is the error from the context (e.g., context cancellation or timeout).
	CtxErr error
	// CauseErr is the underlying cause of the context error.
	CauseErr error
}

// Error implements the error interface.
// It formats the error string to show both the context error and its cause.
//
// Returns:
//   - string: The formatted error message
func (c Error) Error() string {
	return fmt.Sprintf("%v, because %v", c.CtxErr, c.CauseErr)
}

// WrapContextError extracts and wraps errors from the given context.
// It returns a custom Error type containing both the context error and cause error.
// If no errors are present in the context, it returns nil.
//
// Parameters:
//   - ctx: The context to extract errors from
//
// Returns:
//   - error: A wrapped Error if errors exist, nil otherwise
func WrapContextError(ctx context.Context) error {
	ctxErr := ctx.Err()
	causeErr := context.Cause(ctx)
	if ctxErr == nil && causeErr == nil {
		return nil
	}
	return Error{CtxErr: ctxErr, CauseErr: causeErr}
}
