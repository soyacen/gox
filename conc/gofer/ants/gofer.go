// Package ants implements the Gofer interface using the ants library.
package ants

import (
	"context"
	"runtime"
	"time"

	ants "github.com/panjf2000/ants/v2"
	"github.com/soyacen/gox/conc/gofer"
)

// Ensure Gofer implements the gofer.Gofer interface.
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer is an asynchronous task executor implementation based on ants pool.
type Gofer struct {
	// Pool is the underlying ants worker pool.
	Pool *ants.Pool
}

// Go submits a task to the ants pool for execution.
//
// Parameters:
//   - f: the task function to execute.
//
// Returns:
//   - error: the error returned by ants.Pool.Submit.
func (g *Gofer) Go(f func()) error {
	return g.Pool.Submit(f)
}

// Close releases the ants pool resources.
//
// Parameters:
//   - ctx: the context used to control shutdown timeout.
//
// Returns:
//   - error: an error if the shutdown process encounters any issues.
func (g *Gofer) Close(ctx context.Context) error {
	// 获取上下文的截止时间
	deadline, ok := ctx.Deadline()
	if !ok {
		// 如果没有设置截止时间，直接释放资源
		g.Pool.Release()
		// 等待所有任务执行完成
		for g.Pool.Running()+g.Pool.Waiting() > 0 {
			runtime.Gosched()
		}
		return nil
	}

	// 如果设置了截止时间，使用带超时的释放方法
	return g.Pool.ReleaseTimeout(time.Until(deadline))
}
