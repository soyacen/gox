// Package slogx 提供了对 slog 日志库的扩展功能
package slogx

import (
	"context" // 用于上下文处理
	// 用于JSON数据解析
	"log/slog" // Go标准库的日志包
	"net/http" // 用于HTTP服务处理
)

// DebugLevel 返回一个设置为 DEBUG 级别的 LevelVar
func DebugLevel() *slog.LevelVar {
	return NewLevel(slog.LevelDebug)
}

// InfoLevel 返回一个设置为 INFO 级别的 LevelVar
func InfoLevel() *slog.LevelVar {
	return NewLevel(slog.LevelInfo)
}

// WarnLevel 返回一个设置为 WARN 级别的 LevelVar
func WarnLevel() *slog.LevelVar {
	return NewLevel(slog.LevelWarn)
}

// ErrorLevel 返回一个设置为 ERROR 级别的 LevelVar
func ErrorLevel() *slog.LevelVar {
	return NewLevel(slog.LevelError)
}

// NewLevel 创建一个 LevelVar 实例，初始级别为指定的 slog.Level
func NewLevel(l slog.Level) *slog.LevelVar {
	level := &slog.LevelVar{}
	level.Set(l)
	return level
}

// LevelHandler 接口组合了 slog.Handler 和 http.Handler
// 允许同一个处理器既处理日志，又可以通过HTTP接口动态调整日志级别
type LevelHandler interface {
	slog.Handler // slog日志处理器接口
	http.Handler // HTTP请求处理器接口
}

// WithLevel 创建一个 LevelHandler 实例
// 参数 handler: 基础的 slog.Handler
// 参数 levelVar: 可动态调整的日志级别变量
// 返回值: 实现了 LevelHandler 接口的处理器
func WithLevel(handler slog.Handler, levelVar *slog.LevelVar) LevelHandler {
	return &levelVarHandler{
		Handler:  handler,  // 嵌入基础日志处理器
		levelVar: levelVar, // 关联动态级别变量
	}
}

// levelVarHandler 是 LevelHandler 接口的具体实现
// 它包装了一个 slog.Handler 并通过 slog.LevelVar 控制日志级别
type levelVarHandler struct {
	slog.Handler                // 嵌入的基础日志处理器
	levelVar     *slog.LevelVar // 动态日志级别变量
}

// Enabled 检查指定级别的日志是否应该被记录
// 参数 ctx: 日志上下文
// 参数 level: 要检查的日志级别
// 返回值: 如果日志级别大于等于当前设置的级别则返回true，否则返回false
func (h *levelVarHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// 只有当日志级别大于等于当前设置的级别时才启用日志记录
	return level >= h.levelVar.Level()
}

// ServeHTTP 实现 http.Handler 接口，允许通过HTTP请求动态修改日志级别
// 参数 resp: HTTP响应写入器
// 参数 req: HTTP请求对象
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
