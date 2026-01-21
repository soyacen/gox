// Package atomicx provides enhanced atomic operations that extend the standard library's atomic package
package atomicx

import "sync/atomic"

// SubUint32 atomically subtracts delta from *addr and returns the new value
// addr: pointer to the uint32 value to modify
// delta: the value to subtract (can be negative to add)
// Returns the new value after subtraction
func SubUint32(addr *uint32, delta int32) (new uint32) {
	return atomic.AddUint32(addr, ^uint32(delta-1))
}

// SubUint64 atomically subtracts delta from *addr and returns the new value
// addr: pointer to the uint64 value to modify
// delta: the value to subtract (can be negative to add)
// Returns the new value after subtraction
func SubUint64(addr *uint64, delta int64) (new uint64) {
	return atomic.AddUint64(addr, ^uint64(delta-1))
}

// DecrUint32 atomically decrements *addr by 1 and returns the new value
// addr: pointer to the uint32 value to decrement
// Returns the new value after decrementing
func DecrUint32(addr *uint32) (new uint32) {
	return SubUint32(addr, 1)
}

// DecrUint64 atomically decrements *addr by 1 and returns the new value
// addr: pointer to the uint64 value to decrement
// Returns the new value after decrementing
func DecrUint64(addr *uint64) (new uint64) {
	return SubUint64(addr, 1)
}

// IncrUint32 atomically increments *addr by 1 and returns the new value
// addr: pointer to the uint32 value to increment
// Returns the new value after incrementing
func IncrUint32(addr *uint32) (new uint32) {
	return atomic.AddUint32(addr, 1)
}

// IncrUint64 atomically increments *addr by 1 and returns the new value
// addr: pointer to the uint64 value to increment
// Returns the new value after incrementing
func IncrUint64(addr *uint64) (new uint64) {
	return atomic.AddUint64(addr, 1)
}
