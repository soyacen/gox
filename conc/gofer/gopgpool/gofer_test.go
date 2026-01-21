package gopgpool

import (
	"context"
	"testing"
	"time"

	"gopkg.in/go-playground/pool.v3"
)

func TestGofer_Go(t *testing.T) {
	// 创建go-playground池
	pool := newDefaultGoPgPool()
	defer pool.Close()

	gofer := &Gofer{Pool: pool}

	// 测试提交任务
	taskExecuted := make(chan bool, 1)
	err := gofer.Go(func() {
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
	// 创建go-playground池
	pool := newDefaultGoPgPool()
	gofer := &Gofer{Pool: pool}

	// 提交一个任务
	taskStarted := make(chan bool, 1)
	taskFinished := make(chan bool, 1)
	gofer.Go(func() {
		taskStarted <- true
		time.Sleep(100 * time.Millisecond)
		taskFinished <- true
	})

	// 等待任务开始执行
	<-taskStarted

	// 关闭gofer
	ctx := context.Background()
	err := gofer.Close(ctx)
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

// newDefaultGoPgPool 创建默认配置的go-playground池用于测试
func newDefaultGoPgPool() pool.Pool {
	return pool.NewLimited(10)
}
