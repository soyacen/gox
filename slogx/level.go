// Package slogx provides extensions for the slog logging library.
package slogx

import (
	"context"
	"log/slog"
	"net/http"
)

// DebugLevel returns a LevelVar set to DEBUG level.
//
// Returns:
//   - *slog.LevelVar: a LevelVar set to DEBUG level
func DebugLevel() *slog.LevelVar {
	return NewLevel(slog.LevelDebug)
}

// InfoLevel returns a LevelVar set to INFO level.
//
// Returns:
//   - *slog.LevelVar: a LevelVar set to INFO level
func InfoLevel() *slog.LevelVar {
	return NewLevel(slog.LevelInfo)
}

// WarnLevel returns a LevelVar set to WARN level.
//
// Returns:
//   - *slog.LevelVar: a LevelVar set to WARN level
func WarnLevel() *slog.LevelVar {
	return NewLevel(slog.LevelWarn)
}

// ErrorLevel returns a LevelVar set to ERROR level.
//
// Returns:
//   - *slog.LevelVar: a LevelVar set to ERROR level
func ErrorLevel() *slog.LevelVar {
	return NewLevel(slog.LevelError)
}

// NewLevel creates a LevelVar instance with the specified initial level.
//
// Parameters:
//   - l: the initial slog.Level
//
// Returns:
//   - *slog.LevelVar: the created LevelVar
func NewLevel(l slog.Level) *slog.LevelVar {
	level := &slog.LevelVar{}
	level.Set(l)
	return level
}

// LevelHandler combines slog.Handler and http.Handler interfaces.
// It allows the same handler to process logs and dynamically adjust
// the log level via an HTTP endpoint.
type LevelHandler interface {
	slog.Handler
	http.Handler
}

// WithLevel creates a LevelHandler instance.
//
// Parameters:
//   - handler: the base slog.Handler
//   - levelVar: the dynamically adjustable log level variable
//
// Returns:
//   - LevelHandler: the created level handler
func WithLevel(handler slog.Handler, levelVar *slog.LevelVar) LevelHandler {
	return &levelVarHandler{
		Handler:  handler,
		levelVar: levelVar,
	}
}

// levelVarHandler is the concrete implementation of LevelHandler.
// It wraps a slog.Handler and controls the log level via slog.LevelVar.
type levelVarHandler struct {
	slog.Handler
	levelVar *slog.LevelVar
}

// Enabled checks whether a log record at the given level should be logged.
//
// Parameters:
//   - ctx: the log context
//   - level: the log level to check
//
// Returns:
//   - bool: true if the level is enabled
func (h *levelVarHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.levelVar.Level()
}

// ServeHTTP implements http.Handler to allow dynamic log level adjustment via HTTP.
// GET returns the current log level.
// POST updates the log level from the "level" path value.
//
// Parameters:
//   - resp: the HTTP response writer
//   - req: the HTTP request
func (h *levelVarHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		text, err := h.levelVar.MarshalText()
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			_, _ = resp.Write([]byte("Internal Server Error"))
			return
		}

		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write(text)
	case http.MethodPost:
		level := req.PathValue("level")

		if err := h.levelVar.UnmarshalText([]byte(level)); err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			_, _ = resp.Write([]byte("Bad Request"))
			return
		}

		resp.WriteHeader(http.StatusOK)
		_, _ = resp.Write([]byte("OK"))
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = resp.Write([]byte("Method Not Allowed"))
	}
}
