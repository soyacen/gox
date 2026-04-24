package mapx

import (
	"sync"
	"sync/atomic"
)

var _ MapInterface = (*DeepCopyMap)(nil)

// DeepCopyMap is an implementation of MapInterface using a sync.Mutex and
// atomic.Value. It makes deep copies of the map on every write to avoid
// acquiring the Mutex in Load.
type DeepCopyMap struct {
	mu    sync.Mutex
	clean atomic.Value
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
func (m *DeepCopyMap) Load(key any) (value any, loaded bool) {
	clean, _ := m.clean.Load().(map[any]any)
	value, loaded = clean[key]
	return value, loaded
}

// Store sets the value for a key.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (m *DeepCopyMap) Store(key, value any) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
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
func (m *DeepCopyMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	clean, _ := m.clean.Load().(map[any]any)
	actual, loaded = clean[key]
	if loaded {
		return actual, loaded
	}

	m.mu.Lock()
	// Reload clean in case it changed while we were waiting on MapInterface.mu.
	clean, _ = m.clean.Load().(map[any]any)
	actual, loaded = clean[key]
	if !loaded {
		dirty := m.dirty()
		dirty[key] = value
		actual = value
		m.clean.Store(dirty)
	}
	m.mu.Unlock()
	return actual, loaded
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
func (m *DeepCopyMap) Swap(key, value any) (previous any, loaded bool) {
	m.mu.Lock()
	dirty := m.dirty()
	previous, loaded = dirty[key]
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
	return
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
func (m *DeepCopyMap) LoadAndDelete(key any) (value any, loaded bool) {
	m.mu.Lock()
	dirty := m.dirty()
	value, loaded = dirty[key]
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
	return
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (m *DeepCopyMap) Delete(key any) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
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
func (m *DeepCopyMap) CompareAndSwap(key, old, new any) (swapped bool) {
	clean, _ := m.clean.Load().(map[any]any)
	if previous, ok := clean[key]; !ok || previous != old {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	dirty := m.dirty()
	value, loaded := dirty[key]
	if loaded && value == old {
		dirty[key] = new
		m.clean.Store(dirty)
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
func (m *DeepCopyMap) CompareAndDelete(key, old any) (deleted bool) {
	clean, _ := m.clean.Load().(map[any]any)
	if previous, ok := clean[key]; !ok || previous != old {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	dirty := m.dirty()
	value, loaded := dirty[key]
	if loaded && value == old {
		delete(dirty, key)
		m.clean.Store(dirty)
		return true
	}
	return false
}

// Range calls a function for each key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (m *DeepCopyMap) Range(f func(key, value any) (shouldContinue bool)) {
	clean, _ := m.clean.Load().(map[any]any)
	for k, v := range clean {
		if !f(k, v) {
			break
		}
	}
}

// dirty creates a deep copy of the current clean map.
//
// Returns:
//   - map[any]any: a deep copy of the current map
func (m *DeepCopyMap) dirty() map[any]any {
	clean, _ := m.clean.Load().(map[any]any)
	dirty := make(map[any]any, len(clean)+1)
	for k, v := range clean {
		dirty[k] = v
	}
	return dirty
}
