// Package slogx provides extensions for the slog logging library.
package slogx

import (
	"context"
	"log/slog"
	"runtime"
)

// WithCallerSkip wraps the given handler with a callerSkipHandler
// that skips the specified number of stack frames when determining
// the caller location. This allows adjusting the reported source
// location in log output.
//
// Parameters:
//   - handler: the base slog.Handler to wrap
//   - skip: the number of stack frames to skip
//
// Returns:
//   - slog.Handler: the wrapped handler
func WithCallerSkip(handler slog.Handler, skip int) slog.Handler {
	return &callerSkipHandler{Handler: handler, skip: skip}
}

// callerSkipHandler is a decorator for slog.Handler that modifies
// the program counter (PC) used to determine the caller location
// by skipping the specified number of stack frames.
type callerSkipHandler struct {
	slog.Handler
	skip int
}

// Handle processes a log record by setting the record's PC (program counter)
// field based on the current call stack, skipping the configured number of
// frames, then delegating to the wrapped handler.
//
// Parameters:
//   - ctx: the log context
//   - record: the log record to process
//
// Returns:
//   - error: any error encountered during processing
func (h *callerSkipHandler) Handle(ctx context.Context, record slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(h.skip, pcs[:])
	record.PC = pcs[0]
	return h.Handler.Handle(ctx, record)
}
