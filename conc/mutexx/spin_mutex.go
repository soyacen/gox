package mutexx

import (
	"runtime"
	"sync/atomic"
)

const (
	spinMutexLocked = 1 << iota // mutex is locked
	spinMaxBackoff  = 16
)

type SpinMutex struct {
	state int32
}

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

func (m *SpinMutex) Unlock() {
	atomic.StoreInt32(&m.state, 0)
}
