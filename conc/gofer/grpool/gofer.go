// Package grpool 实现了基于grpool库的Gofer接口
package grpool

import (
	"context"

	"github.com/soyacen/gox/conc/gofer"
	"github.com/ivpusic/grpool"
)

// 确保Gofer实现了gofer.Gofer接口
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer 是基于grpool的异步任务执行器实现
type Gofer struct {
	// Pool 底层的grpool工作池
	Pool *grpool.Pool
}

// Go 提交一个任务到grpool的工作队列中
// f: 要执行的任务函数
// 返回nil，因为向通道发送任务不会返回错误
func (g *Gofer) Go(f func()) error {
	g.Pool.JobQueue <- f
	return nil
}

// Close 释放grpool资源
// ctx: 上下文参数（当前实现未使用）
// 返回nil，因为grpool.Release不会返回错误
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.Release()
	return nil
}
