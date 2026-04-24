package mapx

import (
	"golang.org/x/exp/maps"
	"sync"
)

var _ MapInterface = (*RWMutexMap)(nil)

// RWMutexMap is an implementation of MapInterface using a sync.RWMutex.
type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[any]any
}

// Load returns the value stored for a key, or nil if no value is present.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to look up
//
// Returns:
//   - value: the value stored for the key, or nil if not found
//   - loaded: true if the key was found, false otherwise
func (m *RWMutexMap) Load(key any) (value any, loaded bool) {
	m.mu.RLock()
	value, loaded = m.dirty[key]
	m.mu.RUnlock()
	return
}

// Store sets the value for a key.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (m *RWMutexMap) Store(key, value any) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[any]any)
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to look up or store
//   - value: the value to store if the key is not present
//
// Returns:
//   - actual: the value stored for the key
//   - loaded: true if the key was found, false otherwise
func (m *RWMutexMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[any]any)
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

// LoadAndDelete returns the value stored for a key and deletes it if present.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to load and delete
//
// Returns:
//   - value: the value stored for the key, or nil if not found
//   - loaded: true if the key was found, false otherwise
func (m *RWMutexMap) LoadAndDelete(key any) (value any, loaded bool) {
	m.mu.Lock()
	value, loaded = m.dirty[key]
	if !loaded {
		m.mu.Unlock()
		return nil, false
	}
	delete(m.dirty, key)
	m.mu.Unlock()
	return value, loaded
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (m *RWMutexMap) Delete(key any) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to swap
//   - value: the new value to store
//
// Returns:
//   - previous: the previous value stored for the key, or nil if not found
//   - loaded: true if the key was found, false otherwise
func (m *RWMutexMap) Swap(key, value any) (previous any, loaded bool) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[any]any)
	}

	previous, loaded = m.dirty[key]
	m.dirty[key] = value
	m.mu.Unlock()
	return
}

// CompareAndSwap swaps the old and new values for a key if the value matches.
//
// Parameters:
//   - key: the key to compare and swap
//   - old: the expected current value
//   - new: the new value to store
//
// Returns:
//   - swapped: true if the swap was performed, false otherwise
func (m *RWMutexMap) CompareAndSwap(key, old, new any) (swapped bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dirty == nil {
		return false
	}

	value, loaded := m.dirty[key]
	if loaded && value == old {
		m.dirty[key] = new
		return true
	}
	return false
}

// CompareAndDelete deletes the value for a key if it matches the old value.
//
// Parameters:
//   - key: the key to compare and delete
//   - old: the expected current value
//
// Returns:
//   - deleted: true if the deletion was performed, false otherwise
func (m *RWMutexMap) CompareAndDelete(key, old any) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.dirty == nil {
		return false
	}

	value, loaded := m.dirty[key]
	if loaded && value == old {
		delete(m.dirty, key)
		return true
	}
	return false
}

// Range calls a function for each key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (m *RWMutexMap) Range(f func(key, value any) (shouldContinue bool)) {
	m.mu.RLock()
	keys := maps.Keys(m.dirty)
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}
