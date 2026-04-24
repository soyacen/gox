package operator

// SwitchCases selects a value from values based on a matching key from keys.
// It searches for the key in the keys slice and returns the corresponding value.
//
// Parameters:
//   - key: The key to search for
//   - keys: The slice of keys
//   - values: The slice of values corresponding to keys
//
// Returns:
//   - V: The value corresponding to the matching key (zero value if not found)
//   - bool: True if the key was found, false otherwise
func SwitchCases[K comparable, V any](key K, keys []K, values []V) (V, bool) {
	var v V
	if len(keys) != len(values) {
		return v, false
	}
	for i := range keys {
		if keys[i] == key {
			return values[i], true
		}
	}
	return v, false
}
