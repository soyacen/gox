package mutexx

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRecursiveMutex(t *testing.T) {
	var mu RecursiveMutex

	// Test single lock/unlock
	mu.Lock()
	mu.Unlock()

	// Test reentrant lock (recursion)
	mu.Lock()
	mu.Lock() // Should not deadlock
	mu.Unlock()
	mu.Unlock()

	// Test panic when unlocking without holding
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic when unlocking unowned mutex")
			}
		}()
		mu.Unlock() // This should panic
	}()
}

func TestTokenRecursiveMutex(t *testing.T) {
	const token1 = 1234
	const token2 = 5678

	var mu TokenRecursiveMutex

	// Test single lock/unlock
	mu.Lock(token1)
	mu.Unlock(token1)

	// Test reentrant lock
	mu.Lock(token1)
	mu.Lock(token1)
	mu.Unlock(token1)
	mu.Unlock(token1)

	// Test panic when unlocking with wrong token
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic when unlocking with wrong token")
			}
		}()
		mu.Lock(token1)
		mu.Unlock(token2) // Should panic
	}()
}

func TestChanMutex(t *testing.T) {
	mu := NewChanMutex()

	// Test basic Lock/Unlock
	mu.Lock()
	if !mu.IsLocked() {
		t.Fatal("mutex should be locked")
	}
	mu.Unlock()
	if mu.IsLocked() {
		t.Fatal("mutex should be unlocked")
	}

	// Test TryLock
	if !mu.TryLock() {
		t.Fatal("TryLock should succeed when unlocked")
	}
	if mu.TryLock() {
		t.Fatal("TryLock should fail when already locked")
	}
	mu.Unlock()

	// Test LockContext with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if mu.LockContext(ctx) {
		t.Fatal("LockContext should fail with cancelled context")
	}

	// Test LockContext that succeeds
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if !mu.LockContext(ctx) {
		t.Fatal("LockContext should succeed")
	}
	mu.Unlock()

	// Test panic when unlocking unlocked ChanMutex
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic when unlocking unlocked ChanMutex")
			}
		}()
		mu.Unlock() // Should panic
	}()
}

func TestSpinMutex(t *testing.T) {
	var mu SpinMutex
	var counter int
	var wg sync.WaitGroup

	const goroutines = 10
	const iterations = 1000

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if counter != goroutines*iterations {
		t.Fatalf("counter mismatch: got %d, want %d", counter, goroutines*iterations)
	}
}

func TestGroupMutex(t *testing.T) {
	var gm GroupMutex

	key1 := "key1"
	key2 := "key2"

	var counter1, counter2 int
	var wg sync.WaitGroup

	const goroutines = 10
	const iterations = 1000

	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				gm.Lock(key1)
				counter1++
				gm.Unlock(key1)
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				gm.Lock(key2)
				counter2++
				gm.Unlock(key2)
			}
		}()
	}

	wg.Wait()
	if counter1 != goroutines*iterations {
		t.Fatalf("counter1 mismatch: got %d, want %d", counter1, goroutines*iterations)
	}
	if counter2 != goroutines*iterations {
		t.Fatalf("counter2 mismatch: got %d, want %d", counter2, goroutines*iterations)
	}
}
