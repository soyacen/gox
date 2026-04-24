package slicex

import (
	"golang.org/x/exp/slices"
	"sync"
)

// RWMutexSlice is a thread-safe slice protected by a read-write mutex.
type RWMutexSlice[S ~[]E, E any] struct {
	mu  sync.RWMutex
	raw S
}

// MakeSlice creates a new RWMutexSlice with the given length and capacity.
//
// Parameters:
//   - length: the initial length.
//   - capacity: the initial capacity.
//
// Returns:
//   - *RWMutexSlice[S, E]: the newly created thread-safe slice.
func MakeSlice[S ~[]E, E any](length, capacity int) *RWMutexSlice[S, E] {
	return &RWMutexSlice[S, E]{raw: make(S, length, capacity)}
}

// WrapSlice wraps an existing slice in a RWMutexSlice.
//
// Parameters:
//   - raw: the slice to wrap.
//
// Returns:
//   - *RWMutexSlice[S, E]: the wrapped thread-safe slice.
func WrapSlice[S ~[]E, E any](raw []E) *RWMutexSlice[S, E] {
	return &RWMutexSlice[S, E]{raw: raw}
}

// Range iterates over each element and calls f. If f returns false, iteration stops.
//
// Parameters:
//   - f: a function receiving the index and element; returning false stops iteration.
func (s *RWMutexSlice[S, E]) Range(f func(index int, elem E) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for index, elem := range s.raw {
		if !f(index, elem) {
			break
		}
	}
}

// Append adds one or more elements to the end of the slice.
//
// Parameters:
//   - elems: the elements to append.
//
// Returns:
//   - *RWMutexSlice[S, E]: the receiver for chaining.
func (s *RWMutexSlice[S, E]) Append(elems ...E) *RWMutexSlice[S, E] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.raw = append(s.raw, elems...)
	return s
}

// Prepend inserts one or more elements at the beginning of the slice.
//
// Parameters:
//   - elems: the elements to prepend.
//
// Returns:
//   - *RWMutexSlice[S, E]: the receiver for chaining.
func (s *RWMutexSlice[S, E]) Prepend(elems ...E) *RWMutexSlice[S, E] {
	s.mu.Lock()
	defer s.mu.Unlock()
	slice := make([]E, len(elems)+len(s.raw))
	copy(slice, elems)
	copy(slice[len(elems):], s.raw)
	s.raw = slice
	return s
}

// Slice returns a new RWMutexSlice containing elements from low to high.
// If max is provided, the capacity of the returned slice is set to max.
//
// Parameters:
//   - low: the start index (inclusive).
//   - high: the end index (exclusive).
//   - max: optional capacity limit.
//
// Returns:
//   - *RWMutexSlice[S, E]: the sub-slice.
func (s *RWMutexSlice[S, E]) Slice(low int, high int, max ...int) *RWMutexSlice[S, E] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if low < 0 || high > len(s.raw) || low > high {
		panic("slicex: slice bounds out of range")
	}
	if len(max) == 0 {
		return WrapSlice[S, E](slices.Clone(s.raw[low:high]))
	}
	if len(max) == 1 {
		if max[0] < high || max[0] > cap(s.raw) {
			panic("slicex: slice capacity out of range")
		}
		return WrapSlice[S, E](slices.Clone(s.raw[low:high:max[0]]))
	}
	panic("slicex: invalid argument")
}

// Index returns the element at the specified index.
//
// Parameters:
//   - x: the index to access.
//
// Returns:
//   - E: the element at the index.
func (s *RWMutexSlice[S, E]) Index(x int) E {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if x < 0 || x >= len(s.raw) {
		panic("slicex: index out of range")
	}
	return s.raw[x]
}

// Len returns the length of the slice.
//
// Returns:
//   - int: the number of elements.
func (s *RWMutexSlice[S, E]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.raw)
}

// Cap returns the capacity of the slice.
//
// Returns:
//   - int: the capacity.
func (s *RWMutexSlice[S, E]) Cap() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cap(s.raw)
}

// Unwrap returns a copy of the underlying slice.
//
// Returns:
//   - []E: a cloned copy of all elements.
func (s *RWMutexSlice[S, E]) Unwrap() []E {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return slices.Clone(s.raw)
}

// ClearAndUnwrap clears the slice and returns a copy of the previous contents.
//
// Returns:
//   - []E: a cloned copy of the elements before clearing.
func (s *RWMutexSlice[S, E]) ClearAndUnwrap() []E {
	s.mu.Lock()
	defer s.mu.Unlock()
	raw := slices.Clone(s.raw)
	s.raw = s.raw[:]
	return raw
}
