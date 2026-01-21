// Package gofer 提供了一个简单的异步任务执行器接口
package gofer

import "context"

// Gofer 定义了异步任务执行器的接口
type Gofer interface {
	// Go 启动一个异步任务
	// f: 要执行的任务函数
	// 返回错误信息，如果启动失败则返回具体错误
	Go(f func()) error
	
	// Close 关闭执行器，等待所有任务完成
	// ctx: 上下文，用于控制关闭超时
	// 返回错误信息，如果关闭过程中出现错误则返回具体错误
	Close(ctx context.Context) error
}