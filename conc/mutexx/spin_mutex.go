package mutexx

import (
	"runtime"
	"sync/atomic"
)

const (
	spinMutexLocked = 1 << iota // mutex is locked
	spinMaxBackoff  = 16
)

// SpinMutex is a spin-lock based mutex that busy-waits until the lock becomes
// available. It uses exponential backoff to reduce contention under high load.
// Suitable for very short critical sections where blocking would be more expensive
// than spinning.
type SpinMutex struct {
	state int32
}

// Lock acquires the lock. It spins using runtime.Gosched() with exponential
// backoff until the lock is available.
func (m *SpinMutex) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapInt32(&m.state, 0, spinMutexLocked) {
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < spinMaxBackoff {
			backoff <<= 1
		}
	}
}

// Unlock releases the lock.
func (m *SpinMutex) Unlock() {
	atomic.StoreInt32(&m.state, 0)
}
