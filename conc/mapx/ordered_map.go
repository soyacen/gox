package mapx

import (
	"container/list"
	"sync"
)

var _ MapInterface = (*OrderedMap)(nil)

type entry[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap is an implementation of MapInterface that maintains the insertion
// order of keys. It is safe for concurrent use.
type OrderedMap struct {
	mu    sync.RWMutex
	items map[any]*list.Element
	list  *list.List
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
func (m *OrderedMap) Load(key any) (value any, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	elem, loaded := m.items[key]
	if !loaded {
		return nil, loaded
	}
	return elem.Value.(*entry[any, any]).Value, loaded
}

// Store sets the value for a key.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (m *OrderedMap) Store(key, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.items[key]
	if !ok {
		entry := &entry[any, any]{Key: key, Value: value}
		elem = m.list.PushBack(entry)
		m.items[key] = elem
		return
	}
	elem.Value.(*entry[any, any]).Value = value
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
func (m *OrderedMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		entry := &entry[any, any]{Key: key, Value: value}
		elem = m.list.PushBack(entry)
		m.items[key] = elem
		return value, loaded
	}
	return elem.Value.(*entry[any, any]).Value, loaded
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
func (m *OrderedMap) LoadAndDelete(key any) (value any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return nil, loaded
	}
	delete(m.items, key)
	m.list.Remove(elem)
	return elem.Value.(*entry[any, any]).Value, loaded
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (m *OrderedMap) Delete(key any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, ok := m.items[key]
	if !ok {
		return
	}
	delete(m.items, key)
	m.list.Remove(elem)
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
func (m *OrderedMap) Swap(key, value any) (previous any, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return previous, loaded
	}
	entry := elem.Value.(*entry[any, any])
	previous = entry.Value
	entry.Value = value
	return previous, loaded
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
func (m *OrderedMap) CompareAndSwap(key, old, new any) (swapped bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return false
	}
	entry := elem.Value.(*entry[any, any])
	if entry.Value != old {
		return false
	}
	entry.Value = new
	return true
}

// CompareAndDelete deletes the value for a key if it matches the old value.
//
// Parameters:
//   - key: the key to compare and delete
//   - old: the expected current value
//
// Returns:
//   - deleted: true if the deletion was performed, false otherwise
func (m *OrderedMap) CompareAndDelete(key, old any) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	elem, loaded := m.items[key]
	if !loaded {
		return false
	}
	entry := elem.Value.(*entry[any, any])
	if entry.Value != old {
		return false
	}
	delete(m.items, key)
	m.list.Remove(elem)
	return true
}

// Range calls a function for each key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (m *OrderedMap) Range(f func(key any, value any) (shouldContinue bool)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for elem := m.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*entry[any, any])
		if !f(entry.Key, entry.Value) {
			break
		}
	}
}

// Front returns the first element in the ordered map.
//
// Returns:
//   - *list.Element: the first element in the map
func (m *OrderedMap) Front() *list.Element {
	return m.list.Front()
}

// Back returns the last element in the ordered map.
//
// Returns:
//   - *list.Element: the last element in the map
func (m *OrderedMap) Back() *list.Element {
	return m.list.Back()
}

// NewOrderedMap creates a new OrderedMap.
//
// Returns:
//   - *OrderedMap: a new ordered map
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items: make(map[any]*list.Element),
		list:  list.New(),
	}
}
