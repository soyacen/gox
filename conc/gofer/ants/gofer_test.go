package ants

import (
	"context"
	"testing"
	"time"

	"github.com/panjf2000/ants/v2"
)

func TestGofer_Go(t *testing.T) {
	// 创建ants池
	pool, err := NewDefaultAntsPool()
	if err != nil {
		t.Fatalf("failed to create ants pool: %v", err)
	}
	defer pool.Release()

	gofer := &Gofer{Pool: pool}

	// 测试提交任务
	taskExecuted := make(chan bool, 1)
	err = gofer.Go(func() {
		taskExecuted <- true
	})
	if err != nil {
		t.Errorf("Go() returned unexpected error: %v", err)
	}

	// 验证任务被执行
	select {
	case <-taskExecuted:
		// 任务已执行
	case <-time.After(1 * time.Second):
		t.Error("task was not executed within timeout")
	}
}

func TestGofer_Close(t *testing.T) {
	// 创建ants池
	pool, err := NewDefaultAntsPool()
	if err != nil {
		t.Fatalf("failed to create ants pool: %v", err)
	}

	gofer := &Gofer{Pool: pool}

	// 提交一个长时间运行的任务
	taskStarted := make(chan bool, 1)
	taskFinished := make(chan bool, 1)
	gofer.Go(func() {
		taskStarted <- true
		time.Sleep(100 * time.Millisecond)
		taskFinished <- true
	})

	// 等待任务开始执行
	<-taskStarted

	// 关闭gofer并验证所有任务完成
	ctx := context.Background()
	err = gofer.Close(ctx)
	if err != nil {
		t.Errorf("Close() returned unexpected error: %v", err)
	}

	// 验证任务已完成
	select {
	case <-taskFinished:
		// 任务已完成
	default:
		t.Error("not all tasks were completed after Close()")
	}
}

// NewDefaultAntsPool 创建默认配置的ants池用于测试
func NewDefaultAntsPool() (*ants.Pool, error) {
	return ants.NewPool(10)
}
