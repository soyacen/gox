package backoff

import (
	"context"
	"time"
)

type key struct{}

// Inject injects a backoff function into the given context and returns a new context.
// This allows carrying backoff configuration through the context chain.
//
// Parameters:
//   - ctx: The parent context to inject the backoff into
//   - backoff: The backoff function to store in the context
//
// Returns:
//   - context.Context: A new context containing the backoff function
//
// Example:
//
//	ctx := Inject(context.Background(), Exponential(time.Second))
func Inject(ctx context.Context, backoff Func) context.Context {
	return context.WithValue(ctx, key{}, backoff)
}

// Context returns a backoff function that retrieves and invokes a backoff from the context.
// If a backoff function exists in the context, it calls it and returns the duration; otherwise returns 0.
// This is primarily used for dynamically configuring retry delays through context values.
//
// Returns:
//   - Func: A backoff function that extracts and executes the backoff from context
//
// Example:
//
//	// Use backoff configured in context, or no backoff if not present
//	backoffFn := Context()
//	duration := backoffFn(ctx, attempt)
func Context() Func {
	return func(ctx context.Context, attempt uint) time.Duration {
		backoff, ok := ctx.Value(key{}).(Func)
		if ok {
			return backoff(ctx, attempt)
		}
		return 0
	}
}
