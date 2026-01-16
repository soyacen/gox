package slogx_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/soyacen/gox/slogx"
)

func TestAppendContext(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("key1", "value1")
	attr2 := slog.Int("key2", 123)

	ctx = slogx.AppendContext(ctx, attr1)
	ctx = slogx.AppendContext(ctx, attr2)

	attrs, ok := slogx.FromContext(ctx)
	if !ok {
		t.Fatal("Failed to retrieve attributes from context")
	}

	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}

	if attrs[0].Key != "key1" || attrs[0].Value.String() != "value1" {
		t.Errorf("Expected key1=value1, got %s=%v", attrs[0].Key, attrs[0].Value)
	}

	if attrs[1].Key != "key2" || attrs[1].Value.Int64() != 123 {
		t.Errorf("Expected key2=123, got %s=%v", attrs[1].Key, attrs[1].Value)
	}
}

func TestNewContext(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("key1", "value1")
	attr2 := slog.Int("key2", 123)

	ctx = slogx.NewContext(ctx, attr1, attr2)

	attrs, ok := slogx.FromContext(ctx)
	if !ok {
		t.Fatal("Failed to retrieve attributes from context")
	}

	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}

	if attrs[0].Key != "key1" || attrs[0].Value.String() != "value1" {
		t.Errorf("Expected key1=value1, got %s=%v", attrs[0].Key, attrs[0].Value)
	}

	if attrs[1].Key != "key2" || attrs[1].Value.Int64() != 123 {
		t.Errorf("Expected key2=123, got %s=%v", attrs[1].Key, attrs[1].Value)
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("key1", "value1")
	attr2 := slog.Int("key2", 123)

	ctx = slogx.NewContext(ctx, attr1, attr2)

	attrs, ok := slogx.FromContext(ctx)
	if !ok {
		t.Fatal("Failed to retrieve attributes from context")
	}

	if len(attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(attrs))
	}

	if attrs[0].Key != "key1" || attrs[0].Value.String() != "value1" {
		t.Errorf("Expected key1=value1, got %s=%v", attrs[0].Key, attrs[0].Value)
	}

	if attrs[1].Key != "key2" || attrs[1].Value.Int64() != 123 {
		t.Errorf("Expected key2=123, got %s=%v", attrs[1].Key, attrs[1].Value)
	}

	// Test with empty context
	ctx = context.Background()
	attrs, ok = slogx.FromContext(ctx)
	if ok {
		t.Fatal("Expected no attributes in empty context")
	}
}

func TestContextAttrHandler_Handle(t *testing.T) {
	h := slogx.WithContext(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	ctx := context.Background()
	attr1 := slog.String("key1", "value1")
	ctx = slogx.AppendContext(ctx, attr1)

	record := slog.Record{
		Time:    time.Now(),
		Message: "this is message",
		Level:   slog.LevelDebug,
	}

	err := h.Handle(ctx, record)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSlog(t *testing.T) {
	h := slogx.WithContext(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	logger := slog.New(h)
	ctx := context.Background()
	attr1 := slog.String("key1", "value1")
	ctx = slogx.AppendContext(ctx, attr1)
	logger.InfoContext(ctx, "this is message")
}
