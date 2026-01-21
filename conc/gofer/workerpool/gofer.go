// Package workerpool 实现了基于workerpool库的Gofer接口
package workerpool

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/soyacen/gox/conc/gofer"
)

// 确保Gofer实现了gofer.Gofer接口
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer 是基于workerpool的异步任务执行器实现
type Gofer struct {
	// Pool 底层的worker池
	Pool *workerpool.WorkerPool
}

// Go 提交一个任务到worker池中执行
// f: 要执行的任务函数
// 返回nil，因为workerpool.Submit不会返回错误
func (g *Gofer) Go(f func()) error {
	g.Pool.Submit(f)
	return nil
}

// Close 关闭worker池并等待所有任务执行完成
// ctx: 上下文参数（当前实现未使用）
// 返回nil，因为workerpool.StopWait不会返回错误
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.StopWait()
	return nil
}
