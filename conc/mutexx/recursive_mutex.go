package mutexx

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// RecursiveMutex wraps a sync.Mutex to provide reentrant locking.
// It tracks the goroutine ID of the current owner and the recursion count.
// The same goroutine can acquire the lock multiple times without deadlocking.
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // current goroutine id that holds the lock
	recursion int32 // recursion count for the owner goroutine
}

// Lock acquires the lock. If the current goroutine already holds the lock,
// it increments the recursion count instead of blocking.
func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

// Unlock releases the lock. If the current goroutine has acquired the lock
// multiple times, it decrements the recursion count instead of fully unlocking.
// Panics if the calling goroutine does not hold the lock.
func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("syncx: wrong the owner(%d): %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

// TokenRecursiveMutex is a token-based recursive mutex.
// It allows reentrant locking by matching a user-provided token rather than
// tracking goroutine IDs. This is useful when the same logical owner may
// execute across multiple goroutines.
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

// Lock acquires the lock using the provided token.
// If the token matches the current holder, the recursion count is incremented.
// Otherwise, it blocks until the lock is available and then claims it with
// the new token.
//
// Parameters:
//   - token: the token identifying the lock owner
func (m *TokenRecursiveMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// Unlock releases the lock using the provided token.
// The token must match the current holder. Panics if it does not.
//
// Parameters:
//   - token: the token that currently holds the lock
func (m *TokenRecursiveMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("syncx: wrong the owner(%d): %d!", m.token, token))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.token, -1)
	m.Mutex.Unlock()
}
