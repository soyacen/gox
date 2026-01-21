// Package ants 实现了基于ants库的Gofer接口
package ants

import (
	"context"
	"runtime"
	"time"

	"github.com/soyacen/gox/conc/gofer"
	ants "github.com/panjf2000/ants/v2"
)

// 确保Gofer实现了gofer.Gofer接口
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer 是基于ants池的异步任务执行器实现
type Gofer struct {
	// Pool 底层的ants工作池
	Pool *ants.Pool
}

// Go 提交一个任务到ants池中执行
// f: 要执行的任务函数
// 返回ants.Pool.Submit的错误结果
func (g *Gofer) Go(f func()) error {
	return g.Pool.Submit(f)
}

// Close 释放ants池资源
// ctx: 上下文参数，用于控制关闭超时
// 返回错误信息，如果关闭过程中出现错误则返回具体错误
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
