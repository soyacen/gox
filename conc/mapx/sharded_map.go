package mapx

var _ MapInterface = (*ShardedMap)(nil)

// ShardedMap is an implementation of MapInterface that distributes keys across
// multiple segments to reduce lock contention.
type ShardedMap struct {
	Segments []MapInterface
	Hash     func(key any) int
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
func (m *ShardedMap) Load(key any) (value any, loaded bool) {
	return m.bucket(key).Load(key)
}

// Store sets the value for a key.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (m *ShardedMap) Store(key, value any) {
	m.bucket(key).Store(key, value)
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
func (m *ShardedMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	return m.bucket(key).LoadOrStore(key, value)
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
func (m *ShardedMap) LoadAndDelete(key any) (value any, loaded bool) {
	return m.bucket(key).LoadAndDelete(key)
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (m *ShardedMap) Delete(key any) {
	m.bucket(key).Delete(key)
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
func (m *ShardedMap) Swap(key, value any) (previous any, loaded bool) {
	return m.bucket(key).Swap(key, value)
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
func (m *ShardedMap) CompareAndSwap(key, old, new any) (swapped bool) {
	return m.bucket(key).CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the value for a key if it matches the old value.
//
// Parameters:
//   - key: the key to compare and delete
//   - old: the expected current value
//
// Returns:
//   - deleted: true if the deletion was performed, false otherwise
func (m *ShardedMap) CompareAndDelete(key, old any) (deleted bool) {
	return m.bucket(key).CompareAndDelete(key, old)
}

// Range calls a function for each key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (m *ShardedMap) Range(f func(key any, value any) (shouldContinue bool)) {
	for _, segment := range m.Segments {
		segment.Range(f)
	}
}

// bucket returns the segment that should contain the given key.
//
// Parameters:
//   - k: the key to hash
//
// Returns:
//   - MapInterface: the segment for the key
func (m *ShardedMap) bucket(k any) MapInterface {
	return m.Segments[m.Hash(k)%len(m.Segments)]
}
