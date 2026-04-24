package slogx

import (
	"context"
	"log/slog"
)

// attrKey is a private type used as a key for storing attributes in the context.
// Using a private struct type ensures no key collisions in the context.
type attrKey struct{}

// AppendContext appends attributes to the context.
// It retrieves existing attributes from the context, creates a new slice with
// the existing attributes and the new ones, and returns a new context with
// the combined attributes.
//
// Parameters:
//   - ctx: the context to append attributes to
//   - attrs: the attributes to append
//
// Returns:
//   - context.Context: a new context with the combined attributes
func AppendContext(ctx context.Context, attrs ...slog.Attr) context.Context {
	src, _ := FromContext(ctx)
	dst := make([]slog.Attr, len(src), len(src)+len(attrs))
	copy(dst, src)
	dst = append(dst, attrs...)
	return context.WithValue(ctx, attrKey{}, dst)
}

// NewContext creates a new context with attributes.
// It stores the provided attributes in the context using attrKey as the key.
//
// Parameters:
//   - ctx: the base context
//   - attrs: the attributes to store
//
// Returns:
//   - context.Context: a new context with the stored attributes
func NewContext(ctx context.Context, attrs ...slog.Attr) context.Context {
	return context.WithValue(ctx, attrKey{}, attrs)
}

// FromContext returns the attributes from the context.
// It retrieves attributes stored with attrKey and returns them along with
// a boolean indicating whether attributes were found.
//
// Parameters:
//   - ctx: the context to retrieve attributes from
//
// Returns:
//   - []slog.Attr: the retrieved attributes
//   - bool: true if attributes were found
func FromContext(ctx context.Context) ([]slog.Attr, bool) {
	attr, ok := ctx.Value(attrKey{}).([]slog.Attr)
	return attr, ok
}

// WithContext wraps a slog.Handler with context support.
// It returns a new handler that can extract and log attributes from the context.
//
// Parameters:
//   - handler: the base slog.Handler to wrap
//
// Returns:
//   - slog.Handler: the wrapped handler with context support
func WithContext(handler slog.Handler) slog.Handler {
	return &contextHandler{Handler: handler}
}

// contextHandler is a private type that implements slog.Handler and adds context support.
// It embeds the original handler to delegate all other methods.
type contextHandler struct {
	slog.Handler
}

// Handle processes a log record and adds context attributes if present.
// It extracts attributes from the context using FromContext and adds them to
// the log record before delegating to the embedded handler.
//
// Parameters:
//   - ctx: the log context
//   - record: the log record to process
//
// Returns:
//   - error: any error encountered during processing
func (h *contextHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs, ok := FromContext(ctx)
	if ok {
		record.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, record)
}
