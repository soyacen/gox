package oncex

import (
	"sync"
	"testing"
	"time"
)

func TestGroup_Do(t *testing.T) {
	t.Run("single execution", func(t *testing.T) {
		group := &Group{}
		executionCount := 0

		group.Do("key1", func() {
			executionCount++
		})

		if executionCount != 1 {
			t.Errorf("expected function to be executed once, got %d", executionCount)
		}
	})

	t.Run("multiple calls with same key", func(t *testing.T) {
		group := &Group{}
		executionCount := 0

		// Call Do multiple times with the same key
		group.Do("key1", func() {
			executionCount++
		})
		group.Do("key1", func() {
			executionCount++
		})
		group.Do("key1", func() {
			executionCount++
		})

		// Function should only be executed once
		if executionCount != 1 {
			t.Errorf("expected function to be executed once, got %d", executionCount)
		}
	})

	t.Run("different keys", func(t *testing.T) {
		group := &Group{}
		executionCount := 0
		mutex := sync.Mutex{}

		// Call Do with different keys
		group.Do("key1", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})
		group.Do("key2", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})
		group.Do("key3", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		// Each key should execute the function once
		if executionCount != 3 {
			t.Errorf("expected function to be executed 3 times, got %d", executionCount)
		}
	})

	t.Run("concurrent calls with same key", func(t *testing.T) {
		group := &Group{}
		executionCount := 0
		mutex := sync.Mutex{}

		// Number of concurrent goroutines
		concurrentCalls := 10
		done := make(chan bool, concurrentCalls)

		// Launch multiple goroutines calling Do with the same key
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				group.Do("concurrent_key", func() {
					mutex.Lock()
					executionCount++
					mutex.Unlock()
				})
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < concurrentCalls; i++ {
			<-done
		}

		// Function should only be executed once despite concurrent calls
		if executionCount != 1 {
			t.Errorf("expected function to be executed once, got %d", executionCount)
		}
	})

	t.Run("concurrent calls with different keys", func(t *testing.T) {
		group := &Group{}
		executionCount := 0
		mutex := sync.Mutex{}

		// Number of concurrent goroutines
		concurrentCalls := 10
		done := make(chan bool, concurrentCalls)

		// Launch multiple goroutines calling Do with different keys
		for i := 0; i < concurrentCalls; i++ {
			key := i // Capture loop variable
			go func() {
				group.Do(key, func() {
					mutex.Lock()
					executionCount++
					mutex.Unlock()
				})
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < concurrentCalls; i++ {
			<-done
		}

		// Each key should execute the function once
		if executionCount != concurrentCalls {
			t.Errorf("expected function to be executed %d times, got %d", concurrentCalls, executionCount)
		}
	})

	t.Run("nil function", func(t *testing.T) {
		group := &Group{}

		// Should not panic when function is nil
		group.Do("key1", nil)

		// Test passes if no panic occurs
	})

	t.Run("different key types", func(t *testing.T) {
		group := &Group{}
		executionCount := 0
		mutex := sync.Mutex{}

		// Test with different key types
		group.Do("string_key", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		group.Do(42, func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		group.Do(struct{ A, B int }{1, 2}, func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		// Each different key type should execute the function once
		if executionCount != 3 {
			t.Errorf("expected function to be executed 3 times, got %d", executionCount)
		}
	})

	t.Run("execution order", func(t *testing.T) {
		group := &Group{}
		var executionOrder []int
		mutex := sync.Mutex{}

		// Start multiple goroutines to test execution order
		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			group.Do("order_key", func() {
				mutex.Lock()
				executionOrder = append(executionOrder, 1)
				mutex.Unlock()
			})
		}()

		go func() {
			defer wg.Done()
			group.Do("order_key", func() {
				mutex.Lock()
				executionOrder = append(executionOrder, 2)
				mutex.Unlock()
			})
		}()

		go func() {
			defer wg.Done()
			group.Do("order_key", func() {
				mutex.Lock()
				executionOrder = append(executionOrder, 3)
				mutex.Unlock()
			})
		}()

		// Wait for all goroutines to complete
		wg.Wait()

		// Only one execution should happen, and it should be the first one
		if len(executionOrder) != 1 {
			t.Errorf("expected only one execution, got %d", len(executionOrder))
		}

		// Note: In practice, which goroutine executes first is not guaranteed,
		// but sync.Once ensures only one will execute
	})

	t.Run("reuse of sync.Once instances", func(t *testing.T) {
		group := &Group{}
		executionCount := 0
		mutex := sync.Mutex{}

		// Execute with a key
		group.Do("reuse_key", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		// Execute with the same key again (should not execute)
		group.Do("reuse_key", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		// Execute with a different key (should execute)
		group.Do("reuse_key2", func() {
			mutex.Lock()
			executionCount++
			mutex.Unlock()
		})

		// First call for each unique key should execute
		if executionCount != 2 {
			t.Errorf("expected function to be executed 2 times, got %d", executionCount)
		}
	})
}

func TestGroup_ConcurrentPerformance(t *testing.T) {
	group := &Group{}
	executionCount := 0
	mutex := sync.Mutex{}

	const goroutines = 100
	const executionsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	start := time.Now()

	// Launch many goroutines making repeated calls
	for i := 0; i < goroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < executionsPerGoroutine; j++ {
				key := goroutineID % 10 // Create some key contention
				group.Do(key, func() {
					mutex.Lock()
					executionCount++
					mutex.Unlock()
				})
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(start)

	// We should have 10 unique keys, each executed once
	if executionCount != 10 {
		t.Errorf("expected function to be executed 10 times, got %d", executionCount)
	}

	// Log performance metrics
	t.Logf("Executed %d operations in %v with %d goroutines",
		goroutines*executionsPerGoroutine, duration, goroutines)
}
