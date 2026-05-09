// Package slogx provides extensions to the standard log/slog package.
package slogx

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

// MultiHandler is a slog.Handler that dispatches log records to multiple handlers simultaneously.
//
// It implements slog.Handler and forwards all calls to each underlying handler.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a new MultiHandler that dispatches to the given handlers.
//
// Nil handlers are filtered out.
//
// Parameters:
//   - handlers: The handlers to dispatch log records to
//
// Returns:
//   - *MultiHandler: A new MultiHandler instance
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	// Filter out nil handlers
	validHandlers := make([]slog.Handler, 0, len(handlers))
	for _, h := range handlers {
		if h != nil {
			validHandlers = append(validHandlers, h)
		}
	}
	return &MultiHandler{handlers: validHandlers}
}

// Enabled reports whether any handler is enabled for the given level.
// It returns true if at least one handler is enabled.
//
// Parameters:
//   - ctx: Context for the log record
//   - level: The log level to check
//
// Returns:
//   - bool: true if any handler is enabled for the level
func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle dispatches the log record to all handlers concurrently.
// It waits for all handlers to complete and joins any errors returned.
//
// Parameters:
//   - ctx: Context for the log record
//   - record: The log record to handle
//
// Returns:
//   - error: Joined errors from handlers, if any
func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var wg sync.WaitGroup
	errs := make([]error, len(m.handlers))

	for i, h := range m.handlers {
		wg.Add(1)
		go func(idx int, handler slog.Handler) {
			defer wg.Done()
			// Copy record to avoid concurrent issues
			recCopy := record
			errs[idx] = handler.Handle(ctx, recCopy)
		}(i, h)
	}

	wg.Wait()

	// Join all errors
	var resultErr error
	for _, err := range errs {
		if err != nil {
			resultErr = errors.Join(resultErr, err)
		}
	}
	return resultErr
}

// WithAttrs returns a new MultiHandler with the given attributes added to all handlers.
//
// Parameters:
//   - attrs: The attributes to add
//
// Returns:
//   - slog.Handler: A new MultiHandler with attributes added
func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(newHandlers...)
}

// WithGroup returns a new MultiHandler with the given group added to all handlers.
//
// Parameters:
//   - name: The group name to add
//
// Returns:
//   - slog.Handler: A new MultiHandler with group added
func (m *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(newHandlers...)
}
