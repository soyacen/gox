package slogx

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiHandler(t *testing.T) {
	// Create two buffers to capture output
	var buf1, buf2 bytes.Buffer

	// Create two handlers
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{AddSource: false})
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{AddSource: false})

	// Create MultiHandler
	multiHandler := NewMultiHandler(handler1, handler2)

	// Create logger
	logger := slog.New(multiHandler)

	// Log message
	logger.Info("test message", "key", "value")

	// Verify both buffers have output
	assert.Contains(t, buf1.String(), "test message")
	assert.Contains(t, buf1.String(), "key")
	assert.Contains(t, buf1.String(), "value")

	assert.Contains(t, buf2.String(), "test message")
	assert.Contains(t, buf2.String(), "key")
	assert.Contains(t, buf2.String(), "value")
}

func TestMultiHandlerWithNil(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{AddSource: false})

	// Test with nil handlers included
	multiHandler := NewMultiHandler(nil, handler, nil)
	assert.Len(t, multiHandler.handlers, 1)

	logger := slog.New(multiHandler)
	logger.Info("test with nil")
	assert.Contains(t, buf.String(), "test with nil")
}

func TestMultiHandlerEnabled(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	// handler1 only logs Info level and above
	handler1 := slog.NewTextHandler(&buf1, &slog.HandlerOptions{Level: slog.LevelInfo})
	// handler2 only logs Warn level and above
	handler2 := slog.NewTextHandler(&buf2, &slog.HandlerOptions{Level: slog.LevelWarn})

	multiHandler := NewMultiHandler(handler1, handler2)

	// Debug level should not be enabled for either
	assert.False(t, multiHandler.Enabled(context.Background(), slog.LevelDebug))
	// Info level should be enabled (because handler1 is enabled)
	assert.True(t, multiHandler.Enabled(context.Background(), slog.LevelInfo))
	// Warn level should be enabled
	assert.True(t, multiHandler.Enabled(context.Background(), slog.LevelWarn))
}
