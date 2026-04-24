package mapx

// GenericMap is a generic wrapper around MapInterface that provides type-safe
// access to the underlying map.
type GenericMap[K comparable, V any] struct {
	MapInterface MapInterface
}

// Load returns the value stored for a key, or the zero value if no value is present.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to look up
//
// Returns:
//   - value: the value stored for the key, or the zero value if not found
//   - loaded: true if the key was found, false otherwise
func (m *GenericMap[K, V]) Load(key K) (V, bool) {
	var v V
	load, ok := m.MapInterface.Load(key)
	if !ok {
		return v, false
	}
	v = load.(V)
	return v, true
}

// Store sets the value for a key.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (m *GenericMap[K, V]) Store(key K, value V) {
	m.MapInterface.Swap(key, value)
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
func (m *GenericMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	var v V
	actual, loaded := m.MapInterface.LoadOrStore(key, value)
	if !loaded {
		return v, false
	}
	v = actual.(V)
	return v, true
}

// LoadAndDelete returns the value stored for a key and deletes it if present.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to load and delete
//
// Returns:
//   - value: the value stored for the key, or the zero value if not found
//   - loaded: true if the key was found, false otherwise
func (m *GenericMap[K, V]) LoadAndDelete(key K) (V, bool) {
	var v V
	value, loaded := m.MapInterface.LoadAndDelete(key)
	if !loaded {
		return v, false
	}
	v = value.(V)
	return v, true
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (m *GenericMap[K, V]) Delete(key K) {
	m.MapInterface.Delete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was found.
//
// Parameters:
//   - key: the key to swap
//   - value: the new value to store
//
// Returns:
//   - previous: the previous value stored for the key, or the zero value if not found
//   - loaded: true if the key was found, false otherwise
func (m *GenericMap[K, V]) Swap(key K, value V) (V, bool) {
	var v V
	previous, loaded := m.MapInterface.Swap(key, value)
	if !loaded {
		return v, false
	}
	v = previous.(V)
	return v, true
}

// CompareAndSwap swaps the old and new values for a key if the value matches.
//
// Parameters:
//   - key: the key to compare and swap
//   - old: the expected current value
//   - new: the new value to store
//
// Returns:
//   - bool: true if the swap was performed, false otherwise
func (m *GenericMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.MapInterface.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the value for a key if it matches the old value.
//
// Parameters:
//   - key: the key to compare and delete
//   - old: the expected current value
//
// Returns:
//   - bool: true if the deletion was performed, false otherwise
func (m *GenericMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.MapInterface.CompareAndDelete(key, old)
}

// Range calls a function for each key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (m *GenericMap[K, V]) Range(f func(key K, value V) bool) {
	m.MapInterface.Range(func(key, value any) bool {
		k := key.(K)
		v := value.(V)
		return f(k, v)
	})
}
