// Package gopgpool implements the Gofer interface using the go-playground/pool library.
package gopgpool

import (
	"context"
	"sync"

	"github.com/soyacen/gox/conc/gofer"
	"gopkg.in/go-playground/pool.v3"
)

// Ensure Gofer implements the gofer.Gofer interface.
var _ gofer.Gofer = (*Gofer)(nil)

// Gofer is an asynchronous task executor implementation based on go-playground/pool.
type Gofer struct {
	// Pool is the underlying go-playground worker pool.
	Pool pool.Pool
	// m is used to store work units to ensure all tasks are completed.
	m sync.Map
}

// Go submits a task to the go-playground pool for execution.
//
// Parameters:
//   - f: the task function to execute.
//
// Returns:
//   - error: always nil as task submission always succeeds.
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

// Close shuts down the go-playground pool and waits for all tasks to complete.
//
// Parameters:
//   - ctx: the context used to control shutdown timeout (not used in the current implementation).
//
// Returns:
//   - error: always nil as pool.Close does not return an error.
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
