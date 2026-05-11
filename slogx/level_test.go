// file: level_test.go
package slogx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/soyacen/gox/slogx"
)

func TestLevelVarHandler_Enabled(t *testing.T) {
	// 创建一个 LevelVar 并设置初始级别为 INFO
	levelVar := &slog.LevelVar{}
	levelVar.Set(slog.LevelInfo)

	// 创建一个基础的文本处理器
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)

	// 创建 levelVarHandler
	lvh := slogx.WithLevel(handler, levelVar)

	tests := []struct {
		name     string
		level    slog.Level
		expected bool
	}{
		{
			name:     "Debug level disabled",
			level:    slog.LevelDebug,
			expected: false,
		},
		{
			name:     "Info level enabled",
			level:    slog.LevelInfo,
			expected: true,
		},
		{
			name:     "Warn level enabled",
			level:    slog.LevelWarn,
			expected: true,
		},
		{
			name:     "Error level enabled",
			level:    slog.LevelError,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lvh.Enabled(context.Background(), tt.level)
			if result != tt.expected {
				t.Errorf("Enabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWithLevelHandler(t *testing.T) {
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	levelVar := &slog.LevelVar{}

	levelHandler := slogx.WithLevel(handler, levelVar)

	// 检查返回的是否是正确的类型
	if levelHandler == nil {
		t.Error("WithLevelHandler() returned nil")
	}

	// 检查是否实现了所需的接口
	if _, ok := levelHandler.(slog.Handler); !ok {
		t.Error("Returned handler does not implement slog.Handler")
	}

	if _, ok := levelHandler.(http.Handler); !ok {
		t.Error("Returned handler does not implement http.Handler")
	}
}

func TestLevelVarHandler_ServeHTTP_Success(t *testing.T) {
	// 创建 LevelVar 和 handler
	levelVar := &slog.LevelVar{}
	levelVar.Set(slog.LevelInfo) // 初始设置为 INFO

	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	lvh := slogx.WithLevel(handler, levelVar)

	// handler 通过 req.PathValue("level") 读取新级别,
	// 因此必须经过带路径参数的 ServeMux 才能被正确填充。
	mux := http.NewServeMux()
	mux.Handle("POST /level/{level}", lvh)

	// 创建 HTTP 请求(路径中携带级别)
	req := httptest.NewRequest(http.MethodPost, "/level/DEBUG", nil)

	// 创建响应记录器
	rr := httptest.NewRecorder()

	// 经由 mux 调度到 handler
	mux.ServeHTTP(rr, req)

	// 检查响应
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// 检查级别是否已更新
	if levelVar.Level() != slog.LevelDebug {
		t.Errorf("level was not updated: got %v want %v", levelVar.Level(), slog.LevelDebug)
	}
}

func TestLevelVarHandler_ServeHTTP_WrongMethod(t *testing.T) {
	levelVar := &slog.LevelVar{}
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	lvh := slogx.WithLevel(handler, levelVar)

	// handler GET 用于回读当前级别, POST 用于更新,
	// 其它方法返回 405; 这里用 DELETE 验证 405 分支。
	req := httptest.NewRequest(http.MethodDelete, "/level", nil)
	rr := httptest.NewRecorder()

	lvh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	expected := "Method Not Allowed"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLevelVarHandler_ServeHTTP_InvalidJSON(t *testing.T) {
	levelVar := &slog.LevelVar{}
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	lvh := slogx.WithLevel(handler, levelVar)

	// 发送无效的 JSON
	req := httptest.NewRequest(http.MethodPost, "/level", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	lvh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Bad Request"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLevelVarHandler_ServeHTTP_InvalidLevel(t *testing.T) {
	levelVar := &slog.LevelVar{}
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	lvh := slogx.WithLevel(handler, levelVar)

	// 发送无效的日志级别
	body := map[string]string{"level": "INVALID_LEVEL"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/level", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	lvh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Bad Request"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLevelVarHandler_ServeHTTP_MissingLevelField(t *testing.T) {
	levelVar := &slog.LevelVar{}
	handler := slog.NewTextHandler(&bytes.Buffer{}, nil)
	lvh := slogx.WithLevel(handler, levelVar)

	// 发送缺少 level 字段的 JSON
	body := map[string]string{"other_field": "DEBUG"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/level", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	lvh.ServeHTTP(rr, req)

	// 这里应该返回 Bad Request，因为 UnmarshalText 会失败
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
