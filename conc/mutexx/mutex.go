package mutexx

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift      = iota
	starvationThresholdNs = 1e6
)

// TryLock attempts to acquire the lock without blocking.
// It returns true if the lock was successfully acquired, false otherwise.
// This is useful for non-blocking synchronization scenarios.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to lock
//
// Returns:
//   - bool: True if lock acquired, false if already locked
func TryLock(m *sync.Mutex) bool {
	// 如果能成功抢到锁
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m)), 0, mutexLocked) {
		return true
	}
	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回false
	oldState := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	if oldState&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}
	// 尝试在竞争的状态下请求锁
	newState := oldState | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(m)), oldState, newState)
}

// WaiterCount returns the number of goroutines waiting to acquire the lock.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to inspect
//
// Returns:
//   - int: Number of waiting goroutines
func WaiterCount(m *sync.Mutex) int {
	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m)))
	v = v >> mutexWaiterShift //得到等待者的数值
	return int(v)
}

// HolderCount returns the number of lock holders (0 or 1).
// A mutex can only be held by at most one goroutine at a time.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to inspect
//
// Returns:
//   - int: 1 if locked, 0 if unlocked
func HolderCount(m *sync.Mutex) int {
	v := WaiterCount(m)
	return v & mutexLocked
}

// IsLocked checks whether the mutex is currently locked.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to check
//
// Returns:
//   - bool: True if the mutex is locked, false otherwise
func IsLocked(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&spinMutexLocked == spinMutexLocked
}

// IsWoken checks whether a waiting goroutine has been woken up.
// This indicates that the mutex has signaled a waiting goroutine.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to check
//
// Returns:
//   - bool: True if a waiter has been woken, false otherwise
func IsWoken(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&mutexWoken == mutexWoken
}

// IsStarving checks whether the mutex is in starvation mode.
// Starvation mode occurs when waiters have been waiting for a long time,
// causing the mutex to prioritize fairness over throughput.
//
// Parameters:
//   - m: Pointer to the sync.Mutex to check
//
// Returns:
//   - bool: True if in starvation mode, false otherwise
func IsStarving(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&mutexStarving == mutexStarving
}
