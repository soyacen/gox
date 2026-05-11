// file: logx/slogx/logf_test.go
package slogx_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/soyacen/gox/slogx"
)

func TestFormatLogger(t *testing.T) {
	// 创建一个缓冲区来捕获日志输出
	var buf bytes.Buffer

	// 创建一个新的 slog.Logger 实例
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	// 设置全局格式化日志记录器
	slogx.SetFormatLogger(logger)
	// SetFormatLogger 会用全局 levelVar(默认 INFO) 包装 handler,即使
	// 底层 handler 配置为 LevelDebug, Debug 记录仍会被外层包装过滤掉。
	// 测试 Debugf 前先把全局级别下调到 Debug, 测试结束后恢复 INFO。
	slogx.SetFormatLoggerLevel(slog.LevelDebug)
	t.Cleanup(func() { slogx.SetFormatLoggerLevel(slog.LevelInfo) })

	// 测试上下文
	ctx := context.Background()

	// 测试 Debugf
	buf.Reset()
	slogx.Debugf(ctx, "debug message: %s", "test")
	if !strings.Contains(buf.String(), "debug message: test") {
		t.Errorf("Debugf did not produce expected output. Got: %s", buf.String())
	}

	// 测试 Infof
	buf.Reset()
	slogx.Infof(ctx, "info message: %d", 42)
	if !strings.Contains(buf.String(), "info message: 42") {
		t.Errorf("Infof did not produce expected output. Got: %s", buf.String())
	}

	// 测试 Warnf
	buf.Reset()
	slogx.Warnf(ctx, "warn message: %v", true)
	if !strings.Contains(buf.String(), "warn message: true") {
		t.Errorf("Warnf did not produce expected output. Got: %s", buf.String())
	}

	// 测试 ErrorContext
	buf.Reset()
	slogx.Errorf(ctx, "error message: %f", 3.14)
	if !strings.Contains(buf.String(), "error message: 3.14") {
		t.Errorf("ErrorContext did not produce expected output. Got: %s", buf.String())
	}

	// 测试 Logf
	buf.Reset()
	slogx.Logf(ctx, slog.LevelInfo, "log message: %s", "hello")
	if !strings.Contains(buf.String(), "log message: hello") {
		t.Errorf("Logf did not produce expected output. Got: %s", buf.String())
	}
}

func TestSetFormatLogger(t *testing.T) {
	// 测试设置 logger 是否正常工作
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	logger := slog.New(handler)

	slogx.SetFormatLogger(logger)

	// 验证 logger 被正确设置
	ctx := context.Background()
	slogx.Infof(ctx, "test message")
	if !strings.Contains(buf.String(), "test message") {
		t.Error("SetFormatLogger did not set the logger correctly")
	}
}
