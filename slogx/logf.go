// Package slogx provides extensions for the slog logging library,
// primarily adding formatted log recording capabilities.
package slogx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

// formatLogger is a global atomic pointer to a slog.Logger instance.
// It uses atomic.Pointer for safe concurrent access.
var formatLogger atomic.Pointer[slog.Logger]

// formatLoggerLevel is a global log level variable.
var formatLoggerLevel *slog.LevelVar

// formatLoggerLevelHttpHandler is an HTTP handler for dynamically adjusting the log level.
var formatLoggerLevelHttpHandler LevelHandler

// init initializes the package defaults.
func init() {
	formatLoggerLevel = InfoLevel()
	SetFormatLogger(slog.Default())
}

// SetFormatLogger sets the global formatted logger.
// It wraps the given logger with caller skip and context handling.
//
// Parameters:
//   - l: the slog.Logger instance to set
func SetFormatLogger(l *slog.Logger) {
	formatLoggerLevelHttpHandler = WithLevel(WithCallerSkip(WithContext(l.Handler()), 5), formatLoggerLevel)
	formatLogger.Store(slog.New(formatLoggerLevelHttpHandler))
}

// SetFormatLoggerLevel sets the log level for the formatted logger.
//
// Parameters:
//   - l: the slog.Level to set
func SetFormatLoggerLevel(l slog.Level) {
	formatLoggerLevel.Set(l)
}

// FormatLoggerLevelHandler returns the HTTP handler for dynamically adjusting the log level.
//
// Returns:
//   - http.Handler: the HTTP handler for log level adjustment
func FormatLoggerLevelHandler() http.Handler {
	return formatLoggerLevelHttpHandler
}

// Debugf logs a formatted DEBUG level message.
//
// Parameters:
//   - ctx: the context
//   - format: the format string
//   - a: the format arguments
func Debugf(ctx context.Context, format string, a ...any) {
	formatLogger.Load().DebugContext(ctx, fmt.Sprintf(format, a...))
}

// Infof logs a formatted INFO level message.
//
// Parameters:
//   - ctx: the context
//   - format: the format string
//   - a: the format arguments
func Infof(ctx context.Context, format string, a ...any) {
	formatLogger.Load().InfoContext(ctx, fmt.Sprintf(format, a...))
}

// Warnf logs a formatted WARN level message.
//
// Parameters:
//   - ctx: the context
//   - format: the format string
//   - a: the format arguments
func Warnf(ctx context.Context, format string, a ...any) {
	formatLogger.Load().WarnContext(ctx, fmt.Sprintf(format, a...))
}

// Errorf logs a formatted ERROR level message.
//
// Parameters:
//   - ctx: the context
//   - format: the format string
//   - a: the format arguments
func Errorf(ctx context.Context, format string, a ...any) {
	formatLogger.Load().ErrorContext(ctx, fmt.Sprintf(format, a...))
}

// Logf logs a formatted message at the specified level.
//
// Parameters:
//   - ctx: the context
//   - level: the log level
//   - format: the format string
//   - a: the format arguments
func Logf(ctx context.Context, level slog.Level, format string, a ...any) {
	formatLogger.Load().Log(ctx, level, fmt.Sprintf(format, a...))
}
