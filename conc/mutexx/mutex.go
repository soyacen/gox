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

// TryLock 尝试获取锁
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

// WaiterCount 锁的等待者数量
func WaiterCount(m *sync.Mutex) int {
	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m)))
	v = v >> mutexWaiterShift //得到等待者的数值
	return int(v)
}

// HolderCount 锁的持有者数量 0或者1
func HolderCount(m *sync.Mutex) int {
	v := WaiterCount(m)
	return v & mutexLocked
}

// IsLocked 锁是否被持有
func IsLocked(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&spinMutexLocked == spinMutexLocked
}

// IsWoken 是否有等待者被唤醒
func IsWoken(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func IsStarving(m *sync.Mutex) bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(m)))
	return state&mutexStarving == mutexStarving
}
