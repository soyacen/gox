// Package slogx_test 提供了对 slogx 扩展功能的单元测试
package slogx_test

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/soyacen/gox/slogx"
)

// TestWithCallerSkipHandle 测试 WithCallerSkipHandle 函数是否正确包装 Handler
func TestWithCallerSkipHandle(t *testing.T) {
	// 创建一个内存缓冲区作为日志输出目标
	var buf bytes.Buffer

	// 创建基础的 slog.Handler
	baseHandler := slog.NewTextHandler(&buf, nil)

	// 使用 WithCallerSkipHandle 包装基础 Handler
	wrappedHandler := slogx.WithCallerSkip(baseHandler, 2)

	// 创建 logger 并记录一条日志
	logger := slog.New(wrappedHandler)
	logger.Info("test message")

	// 验证日志是否被正确记录
	output := buf.String()
	if !strings.Contains(output, "level=INFO") || !strings.Contains(output, "msg=\"test message\"") {
		t.Errorf("Expected log output to contain level and message, got: %s", output)
	}
}

// TestCallerSkipHandler_Handle 测试 callerSkipHandler 的 Handle 方法是否正确处理日志记录
func TestCallerSkipHandler_Handle(t *testing.T) {
	// 创建一个自定义的测试 Handler 来验证 PC 值
	testHandler := &testHandlerWithPC{}

	// 创建 callerSkipHandler 实例
	callerSkipHandler := slogx.WithCallerSkip(testHandler, 1)

	// 创建日志记录
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	// 处理日志记录
	err := callerSkipHandler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	// 验证 PC 是否被设置
	if testHandler.lastRecord == nil {
		t.Fatal("Expected record to be handled")
	}

	if testHandler.lastRecord.PC == 0 {
		t.Error("Expected PC to be set")
	}

	// 验证消息内容是否正确传递
	if testHandler.lastRecord.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", testHandler.lastRecord.Message)
	}
}

// testHandlerWithPC 是一个用于测试的 Handler 实现，用来捕获处理的记录
type testHandlerWithPC struct {
	lastRecord *slog.Record
}

// Enabled 实现 slog.Handler 接口
func (h *testHandlerWithPC) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// Handle 实现 slog.Handler 接口，并保存最后处理的记录
func (h *testHandlerWithPC) Handle(_ context.Context, record slog.Record) error {
	h.lastRecord = &record
	return nil
}

// WithAttrs 实现 slog.Handler 接口
func (h *testHandlerWithPC) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup 实现 slog.Handler 接口
func (h *testHandlerWithPC) WithGroup(_ string) slog.Handler {
	return h
}

// TestCallerSkipValue 测试不同的 skip 值是否影响调用栈跟踪
func TestCallerSkipValue(t *testing.T) {
	tests := []struct {
		name string
		skip int
	}{
		{"skip 0", 0},
		{"skip 1", 1},
		{"skip 2", 2},
		{"skip 5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler := &testHandlerWithPC{}
			callerSkipHandler := slogx.WithCallerSkip(testHandler, tt.skip)

			record := slog.NewRecord(time.Now(), slog.LevelInfo, "test", 0)
			err := callerSkipHandler.Handle(context.Background(), record)
			if err != nil {
				t.Fatalf("Handle returned error: %v", err)
			}

			// 验证记录是否被处理
			if testHandler.lastRecord == nil {
				t.Fatal("Expected record to be handled")
			}
		})
	}
}

// 示例测试展示如何使用 WithCallerSkipHandle
func ExampleWithCallerSkip() {
	// 创建一个带有调用者跳过的 Handler
	handler := slogx.WithCallerSkip(slog.NewTextHandler(os.Stdout, nil), 2)

	// 创建 logger
	logger := slog.New(handler)

	// 记录日志
	logger.Info("example message")
}
