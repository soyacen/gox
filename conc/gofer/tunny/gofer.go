// Package tunny 实现了基于tunny库的Gofer接口
package tunny

import (
	"context"

	"github.com/Jeffail/tunny"
	"github.com/soyacen/gox/conc/gofer"
)

// 确保Gofer实现了gofer.Gofer接口
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer 是基于tunny池的异步任务执行器实现
type Gofer struct {
	// Pool 底层的tunny工作池
	Pool *tunny.Pool
}

// Go 提交一个任务到tunny池中执行
// f: 要执行的任务函数
// 返回nil，因为tunny.Pool.Process不会返回错误
func (g *Gofer) Go(f func()) error {
	g.Pool.Process(f)
	return nil
}

// Close 关闭tunny池
// ctx: 上下文参数（当前实现未使用）
// 返回nil，因为tunny.Pool.Close不会返回错误
func (g *Gofer) Close(ctx context.Context) error {
	g.Pool.Close()
	return nil
}
