package mapx

// AppendValue appends values from multiple maps into a map of slices.
//
// Parameters:
//   - r: The result map to append to
//   - ms: The maps to append values from
//
// Returns:
//   - R: The result map with appended values
func AppendValue[M ~map[K]V, R ~map[K]S, S ~[]V, K comparable, V any](r R, ms ...M) R {
	if r == nil {
		r = make(R)
	}
	for _, elem := range ms {
		for k, v := range elem {
			r[k] = append(r[k], v)
		}
	}
	return r
}

// IsEmpty checks if a map is empty.
//
// Parameters:
//   - m: The map to check
//
// Returns:
//   - bool: True if the map is empty, false otherwise
func IsEmpty[M ~map[K]V, K comparable, V any](m M) bool {
	return len(m) <= 0
}

// IsNotEmpty checks if a map is not empty.
//
// Parameters:
//   - m: The map to check
//
// Returns:
//   - bool: True if the map is not empty, false otherwise
func IsNotEmpty[M ~map[K]V, K comparable, V any](m M) bool {
	return len(m) > 0
}

// Entries returns a slice of all key-value entries in the map.
//
// Parameters:
//   - m: The map to extract entries from
//
// Returns:
//   - []Entry[K, V]: A slice of all key-value pairs
func Entries[M ~map[K]V, K comparable, V any](m M) []Entry[K, V] {
	r := make([]Entry[K, V], 0, len(m))
	for k, v := range m {
		r = append(r, Entry[K, V]{Key: k, Value: v})
	}
	return r
}

// KeySet returns a set of all keys in the map.
//
// Parameters:
//   - m: The map to extract keys from
//
// Returns:
//   - R: A set containing all keys from the map
func KeySet[M ~map[K]V, R map[K]struct{}, K comparable, V any](m M) R {
	r := make(R, len(m))
	for k := range m {
		r[k] = struct{}{}
	}
	return r
}

// FromRanger creates a map from an object that implements the Range method.
//
// Parameters:
//   - ranger: An object with a Range method
//
// Returns:
//   - M: A map containing all key-value pairs from the ranger
func FromRanger[M ~map[K]V, K comparable, V any](ranger interface {
	Range(func(key, value any) (shouldContinue bool))
}) M {
	m := make(M)
	ranger.Range(func(k any, v any) bool {
		m[k.(K)] = v.(V)
		return true
	})
	return m
}
