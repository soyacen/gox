package mapx

// MapInterface is a concurrent-safe map interface matching the API of sync.Map.
// It is copied from sync/map_reference_test.go.
//
// Methods:
//   - Load: returns the value stored for a key, or nil if no value is present
//   - Store: sets the value for a key
//   - LoadOrStore: returns the existing value for the key if present, otherwise stores and returns the given value
//   - LoadAndDelete: returns the value stored for a key and deletes it if present
//   - Delete: deletes the value for a key
//   - Swap: swaps the value for a key and returns the previous value
//   - CompareAndSwap: swaps the old and new values for a key if the value matches
//   - CompareAndDelete: deletes the value for a key if it matches the old value
//   - Range: calls a function for each key-value pair in the map
type MapInterface interface {
	Load(key any) (value any, loaded bool)
	Store(key, value any)
	LoadOrStore(key, value any) (actual any, loaded bool)
	LoadAndDelete(key any) (value any, loaded bool)
	Delete(key any)
	Swap(key, value any) (previous any, loaded bool)
	CompareAndSwap(key, old, new any) (swapped bool)
	CompareAndDelete(key, old any) (deleted bool)
	Range(func(key, value any) (shouldContinue bool))
}
