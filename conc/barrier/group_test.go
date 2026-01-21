package barrier

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// TestBasicBarrier 测试基本屏障功能
func TestBasicBarrier(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	var wg sync.WaitGroup
	wg.Add(parties)

	// 启动多个goroutine测试屏障
	for i := 0; i < parties; i++ {
		go func(id int) {
			defer wg.Done()
			err := barrier.Wait(context.Background())
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}

// TestBarrierWithAction 测试带动作的屏障
func TestBarrierWithAction(t *testing.T) {
	const parties = 3
	actionExecuted := false
	actionMutex := sync.Mutex{}

	barrier := NewGroup(parties, func() {
		actionMutex.Lock()
		actionExecuted = true
		actionMutex.Unlock()
	})

	var wg sync.WaitGroup
	wg.Add(parties)

	for i := 0; i < parties; i++ {
		go func(id int) {
			defer wg.Done()
			err := barrier.Wait(context.Background())
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	actionMutex.Lock()
	if !actionExecuted {
		t.Error("barrier action was not executed")
	}
	actionMutex.Unlock()
}

// TestContextCancellation 测试上下文取消
func TestContextCancellation(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	// 先让一个goroutine到达屏障
	subCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := barrier.Wait(subCtx)
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got: %v", err)
		}
	}()

	// 等待一点时间确保第一个goroutine已经等待
	time.Sleep(10 * time.Millisecond)

	// 创建一个会立即取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	err := barrier.Wait(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}

	// 验证屏障是否损坏
	if !barrier.IsBroken() {
		t.Error("barrier should be broken after context cancellation")
	}
}

// TestReset 测试重置功能
func TestReset(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	subCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 让一个goroutine等待
	go func() {
		err := barrier.Wait(subCtx)
		if err == nil {
			t.Error("expected error after reset, got nil")
		}
	}()

	// 等待一点时间确保goroutine已经开始等待
	time.Sleep(10 * time.Millisecond)

	// 重置屏障
	barrier.Reset()

	// 验证屏障已重置
	if barrier.IsBroken() {
		t.Error("barrier should not be broken after reset")
	}

	// 验证计数是否正确重置
	if barrier.GetNumberWaiting() != 0 {
		t.Errorf("expected 0 waiting, got %d", barrier.GetNumberWaiting())
	}
}

// TestConcurrentResets 测试并发重置
func TestConcurrentResets(t *testing.T) {
	const parties = 5
	barrier := NewGroup(parties, nil)

	var wg sync.WaitGroup
	wg.Add(parties)

	// 多个goroutine同时调用Reset
	for i := 0; i < parties; i++ {
		go func() {
			defer wg.Done()
			barrier.Reset()
		}()
	}

	wg.Wait()

	// 验证屏障处于正常状态
	if barrier.IsBroken() {
		t.Error("barrier should not be broken after concurrent resets")
	}
}

// TestPanicInBarrierAction 测试屏障动作中的panic
func TestPanicInBarrierAction(t *testing.T) {
	const parties = 2
	barrier := NewGroup(parties, func() {
		panic("intentional panic in barrier action")
	})

	var wg sync.WaitGroup
	results := make([]error, parties)
	wg.Add(parties)

	for i := 0; i < parties; i++ {
		go func(id int) {
			defer wg.Done()
			results[id] = barrier.Wait(context.Background())
		}(i)
	}

	wg.Wait()

	// 检查所有结果都是BrokenBarrierError
	for i, err := range results {
		if _, ok := err.(*BrokenBarrierError); !ok {
			t.Errorf("goroutine %d: expected BrokenBarrierError, got: %v", i, err)
		}
	}
}

// TestMultipleGenerations 测试多轮屏障使用
func TestMultipleGenerations(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	for round := 0; round < 3; round++ {
		var wg sync.WaitGroup
		wg.Add(parties)

		for i := 0; i < parties; i++ {
			go func(id int, round int) {
				defer wg.Done()
				err := barrier.Wait(context.Background())
				if err != nil {
					t.Errorf("round %d, goroutine %d: unexpected error: %v", round, id, err)
				}
			}(i, round)
		}

		wg.Wait()

		// 检查每轮结束后屏障状态
		if barrier.IsBroken() {
			t.Errorf("round %d: barrier should not be broken", round)
		}
	}
}

// TestGetParties 测试获取参与者数量
func TestGetParties(t *testing.T) {
	const parties = 5
	barrier := NewGroup(parties, nil)

	if barrier.GetParties() != parties {
		t.Errorf("expected %d parties, got %d", parties, barrier.GetParties())
	}
}

// TestGetNumberWaiting 测试获取等待数量
func TestGetNumberWaiting(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	// 初始状态应该没有等待者
	if barrier.GetNumberWaiting() != 0 {
		t.Errorf("expected 0 waiting, got %d", barrier.GetNumberWaiting())
	}

	go func() {
		barrier.Wait(context.Background())
	}()

	// 等待一点时间确保goroutine已经开始等待
	time.Sleep(10 * time.Millisecond)

	if barrier.GetNumberWaiting() != 1 {
		t.Errorf("expected 1 waiting, got %d", barrier.GetNumberWaiting())
	}
}

// TestInvalidParties 测试无效的参与者数量
func TestInvalidParties(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for parties <= 0")
		}
	}()

	NewGroup(0, nil)
}

// TestWaitOnBrokenBarrier 测试在已损坏屏障上等待
func TestWaitOnBrokenBarrier(t *testing.T) {
	barrier := NewGroup(2, func() {
		panic("intentional panic")
	})

	// 第一次使用使屏障损坏
	go barrier.Wait(context.Background())
	time.Sleep(10 * time.Millisecond)
	barrier.Wait(context.Background())

	// 再次等待应该立即返回错误
	err := barrier.Wait(context.Background())
	if _, ok := err.(*BrokenBarrierError); !ok {
		t.Errorf("expected BrokenBarrierError on broken barrier, got: %v", err)
	}
}

// TestConcurrentWaits 测试高并发等待
func TestConcurrentWaits(t *testing.T) {
	const parties = 100
	barrier := NewGroup(parties, nil)

	var wg sync.WaitGroup
	errors := make([]error, parties)
	wg.Add(parties)

	start := time.Now()

	for i := 0; i < parties; i++ {
		go func(id int) {
			defer wg.Done()
			errors[id] = barrier.Wait(context.Background())
		}(i)
	}

	wg.Wait()

	duration := time.Since(start)
	t.Logf("100 goroutines passed barrier in %v", duration)

	// 检查是否有错误
	for i, err := range errors {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %v", i, err)
		}
	}
}

// BenchmarkBarrier 基准测试
func BenchmarkBarrier(b *testing.B) {
	const parties = 10
	barrier := NewGroup(parties, nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			barrier.Wait(context.Background())
		}
	})
}

// TestTimeoutContext 测试超时上下文
func TestTimeoutContext(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	// 创建一个很短超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err := barrier.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got: %v", err)
	}

	// 验证屏障是否损坏
	if !barrier.IsBroken() {
		t.Error("barrier should be broken after timeout")
	}
}

// TestMixedSuccessAndFailure 测试成功和失败混合的情况
func TestMixedSuccessAndFailure(t *testing.T) {
	const parties = 3
	barrier := NewGroup(parties, nil)

	// 成功一轮
	var wg sync.WaitGroup
	wg.Add(parties)
	for i := 0; i < parties; i++ {
		go func() {
			defer wg.Done()
			err := barrier.Wait(context.Background())
			if err != nil {
				t.Errorf("unexpected error in successful round: %v", err)
			}
		}()
	}
	wg.Wait()

	// 损坏屏障
	barrier.Reset()

	// 在损坏的屏障上等待应该失败
	err := barrier.Wait(context.Background())
	var brokenErr *BrokenBarrierError
	if !errors.As(err, &brokenErr) {
		t.Errorf("expected BrokenBarrierError after reset, got: %v", err)
	}
}
