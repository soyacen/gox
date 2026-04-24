package mapx

// Entry represents a key-value pair in a map.
//
// It is used to hold a single entry extracted from a map,
// preserving both the key and its associated value.
//
// Type Parameters:
//   - K: The type of the key, which must be comparable.
//   - V: The type of the value.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
