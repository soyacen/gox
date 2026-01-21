// Package gopgpool 实现了基于go-playground/pool库的Gofer接口
package gopgpool

import (
	"context"
	"sync"

	"github.com/soyacen/gox/conc/gofer"
	"gopkg.in/go-playground/pool.v3"
)

// 确保Gofer实现了gofer.Gofer接口
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer 是基于go-playground/pool的异步任务执行器实现
type Gofer struct {
	// Pool 底层的go-playground工作池
	Pool pool.Pool
	// m 用于存储工作单元，确保所有任务完成
	m sync.Map
}

// Go 提交一个任务到go-playground工作池中执行
// f: 要执行的任务函数
// 返回nil，因为任务提交总是成功
func (g *Gofer) Go(f func()) error {
	// 将任务加入队列
	unit := g.Pool.Queue(func(unit pool.WorkUnit) (interface{}, error) {
		// 执行任务
		f()
		// 任务完成后从map中删除
		g.m.Delete(unit)
		return nil, nil
	})
	// 将工作单元存储到map中
	g.m.Store(unit, struct{}{})
	return nil
}

// Close 关闭go-playground工作池并等待所有任务完成
// ctx: 上下文参数（当前实现未使用）
// 返回nil，因为pool.Close不会返回错误
func (g *Gofer) Close(ctx context.Context) error {
	// 关闭工作池
	g.Pool.Close()
	// 等待所有任务完成
	g.m.Range(func(key, value any) bool {
		unit, _ := key.(pool.WorkUnit)
		unit.Wait()
		return true
	})
	return nil
}
