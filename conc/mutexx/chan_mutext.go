package mutexx

import (
	"context"
	"errors"
)

// ChanMutex is a mutex implementation based on a buffered channel.
// It provides Lock, Unlock, TryLock, LockContext, and IsLocked methods.
type ChanMutex struct {
	state chan struct{}
}

// NewChanMutex creates a new ChanMutex.
// The mutex is initially unlocked.
//
// Returns:
//   - *ChanMutex: a new ChanMutex instance
func NewChanMutex() *ChanMutex {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &ChanMutex{
		state: ch,
	}
}

// Lock acquires the lock. It blocks until the lock is available.
func (m *ChanMutex) Lock() {
	<-m.state
}

// Unlock releases the lock.
// Panics if the mutex is not locked when Unlock is called.
func (m *ChanMutex) Unlock() {
	select {
	case m.state <- struct{}{}:
	default:
		panic(errors.New("mutexx: unlock of unlocked mutex"))
	}

}

// TryLock attempts to acquire the lock without blocking.
//
// Returns:
//   - bool: true if the lock was acquired, false otherwise
func (m *ChanMutex) TryLock() bool {
	select {
	case <-m.state:
		return true
	default:
		return false
	}
}

// LockContext attempts to acquire the lock with a context.
// It returns true if the lock was acquired, or false if the context is done.
//
// Parameters:
//   - ctx: the context for cancellation
//
// Returns:
//   - bool: true if the lock was acquired, false if the context was cancelled
func (m *ChanMutex) LockContext(ctx context.Context) bool {
	select {
	case <-m.state:
		return true
	case <-ctx.Done():
		return false
	}
}

// IsLocked returns true if the lock is currently held.
//
// Returns:
//   - bool: true if locked, false otherwise
func (m *ChanMutex) IsLocked() bool {
	return len(m.state) == 0
}
