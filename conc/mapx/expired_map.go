package mapx

import (
	"runtime"
	"sync"
	"time"
)

// ExpiredMap is a concurrent-safe map that supports key expiration.
// Each key can have an individual expiration time.
type ExpiredMap struct {
	items map[any]expiredMapItem
	mu    sync.RWMutex

	expireAfter func(key any) time.Duration
	onEvicted   func(any, any)
	janitor     *janitor
}

// Load returns the value stored for a key, or nil if no value is present
// or the value has expired.
// The loaded result reports whether the key was found and not expired.
//
// Parameters:
//   - key: the key to look up
//
// Returns:
//   - value: the value stored for the key, or nil if not found or expired
//   - loaded: true if the key was found and not expired, false otherwise
func (c *ExpiredMap) Load(key any) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expired() {
		return nil, false
	}
	return item.Object, true
}

// Store sets the value for a key with expiration.
//
// Parameters:
//   - key: the key to store
//   - value: the value to store
func (c *ExpiredMap) Store(key, value any) {
	expiration := c.expiration(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = expiredMapItem{
		Object:     value,
		Expiration: expiration,
	}
}

// LoadOrStore returns the existing value for the key if present and not expired.
// Otherwise, it stores and returns the given value.
// The loaded result reports whether the key was found and not expired.
//
// Parameters:
//   - key: the key to look up or store
//   - value: the value to store if the key is not present or expired
//
// Returns:
//   - value: the value stored for the key
//   - loaded: true if the key was found and not expired, false otherwise
func (c *ExpiredMap) LoadOrStore(key, value any) (any, bool) {
	expiration := c.expiration(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if !found || item.Expired() {
		c.items[key] = expiredMapItem{
			Object:     value,
			Expiration: expiration,
		}
		return value, false
	}
	return item.Object, true
}

// LoadAndDelete returns the value stored for a key and deletes it if present
// and not expired.
// The loaded result reports whether the key was found and not expired.
//
// Parameters:
//   - key: the key to load and delete
//
// Returns:
//   - value: the value stored for the key, or nil if not found or expired
//   - loaded: true if the key was found and not expired, false otherwise
func (c *ExpiredMap) LoadAndDelete(key any) (any, bool) {
	c.mu.Lock()
	item, found := c.items[key]
	if !found || item.Expired() {
		c.mu.Unlock()
		return nil, false
	}
	delete(c.items, key)
	c.mu.Unlock()
	if c.onEvicted == nil {
		c.onEvicted(key, item.Object)
	}
	return item.Object, true
}

// Delete deletes the value for a key.
//
// Parameters:
//   - key: the key to delete
func (c *ExpiredMap) Delete(key any) {
	_, _ = c.LoadAndDelete(key)
}

// Swap swaps the value for a key and returns the previous value if present
// and not expired.
// The loaded result reports whether the key was found and not expired.
//
// Parameters:
//   - key: the key to swap
//   - value: the new value to store
//
// Returns:
//   - value: the previous value stored for the key, or nil if not found or expired
//   - loaded: true if the key was found and not expired, false otherwise
func (c *ExpiredMap) Swap(key, value any) (any, bool) {
	expiration := c.expiration(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if !found || item.Expired() {
		return nil, false
	}
	c.items[key] = expiredMapItem{
		Object:     value,
		Expiration: expiration,
	}
	return item.Object, true
}

// CompareAndSwap swaps the old and new values for a key if the value matches
// and the key has not expired.
//
// Parameters:
//   - key: the key to compare and swap
//   - oldValue: the expected current value
//   - newValue: the new value to store
//
// Returns:
//   - bool: true if the swap was performed, false otherwise
func (c *ExpiredMap) CompareAndSwap(key, oldValue, newValue any) bool {
	expiration := c.expiration(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	previous, found := c.items[key]
	if !found || previous.Expired() || previous.Object != oldValue {
		return false
	}
	c.items[key] = expiredMapItem{
		Object:     newValue,
		Expiration: expiration,
	}
	return true
}

// CompareAndDelete deletes the value for a key if it matches the old value
// and the key has not expired.
//
// Parameters:
//   - key: the key to compare and delete
//   - oldValue: the expected current value
//
// Returns:
//   - bool: true if the deletion was performed, false otherwise
func (c *ExpiredMap) CompareAndDelete(key, oldValue any) bool {
	c.mu.Lock()
	item, found := c.items[key]
	if !found || item.Expired() || item.Object != oldValue {
		c.mu.Unlock()
		return false
	}
	delete(c.items, key)
	c.mu.Unlock()
	if c.onEvicted == nil {
		c.onEvicted(key, item.Object)
	}
	return true
}

// Range calls a function for each non-expired key-value pair in the map.
// If the function returns false, Range stops the iteration.
//
// Parameters:
//   - f: the function to call for each key-value pair
func (c *ExpiredMap) Range(f func(key any, value any) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for key, value := range c.items {
		if value.Expired() {
			continue
		}
		if !f(key, value) {
			break
		}
	}
}

// expiration calculates the expiration timestamp for a key.
//
// Parameters:
//   - key: the key to calculate expiration for
//
// Returns:
//   - int64: the expiration timestamp in Unix nanoseconds, or 0 if no expiration
func (c *ExpiredMap) expiration(key any) int64 {
	var expiration int64
	expireAfter := c.expireAfter(key)
	if expireAfter > 0 {
		expiration = time.Now().Add(expireAfter).UnixNano()
	}
	return expiration
}

// deleteExpired removes all expired items from the map.
func (c *ExpiredMap) deleteExpired() {
	type keyAndValue struct {
		key   any
		value any
	}
	var evictedItems []keyAndValue
	c.mu.Lock()
	for key, value := range c.items {
		if !value.Expired() {
			continue
		}
		delete(c.items, key)
		if c.onEvicted != nil {
			evictedItems = append(evictedItems, keyAndValue{key: key, value: value.Object})
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

// janitor is a background goroutine that periodically cleans up expired items.
type janitor struct {
	Interval time.Duration
	stop     chan bool
}

// Run starts the janitor's cleanup loop.
//
// Parameters:
//   - c: the ExpiredMap to clean up
func (j *janitor) Run(c *ExpiredMap) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

// expiredMapItem is an item stored in the ExpiredMap with an expiration time.
type expiredMapItem struct {
	Object     any
	Expiration int64
}

// Expired returns true if the item has expired.
//
// Returns:
//   - bool: true if the item has expired, false otherwise
func (item expiredMapItem) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// ExpiredMapOption is a function that configures an ExpiredMap.
type ExpiredMapOption func(*ExpiredMap)

// ExpireAfter sets the function that determines the expiration duration for a key.
//
// Parameters:
//   - f: the function that returns the expiration duration for a key
//
// Returns:
//   - ExpiredMapOption: an option that sets the expiration function
func ExpireAfter(f func(key any) time.Duration) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.expireAfter = f
	}
}

// CleanupInterval sets the interval at which the janitor cleans up expired items.
//
// Parameters:
//   - interval: the cleanup interval
//
// Returns:
//   - ExpiredMapOption: an option that sets the cleanup interval
func CleanupInterval(interval time.Duration) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.janitor.Interval = interval
	}
}

// OnEvicted sets the callback function that is called when an item is evicted.
//
// Parameters:
//   - f: the callback function called with the evicted key and value
//
// Returns:
//   - ExpiredMapOption: an option that sets the eviction callback
func OnEvicted(f func(any, any)) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.onEvicted = f
	}
}

// NewExpiredMap creates a new ExpiredMap with the given options.
//
// Parameters:
//   - options: optional configuration options
//
// Returns:
//   - *ExpiredMap: a new expired map
func NewExpiredMap(options ...ExpiredMapOption) *ExpiredMap {
	c := &ExpiredMap{
		items:       make(map[any]expiredMapItem),
		mu:          sync.RWMutex{},
		expireAfter: func(key any) time.Duration { return 0 },
		onEvicted:   func(key any, val any) {},
		janitor: &janitor{
			stop: make(chan bool),
		},
	}
	for _, option := range options {
		option(c)
	}
	if c.janitor.Interval > 0 {
		go c.janitor.Run(c)
		runtime.SetFinalizer(c, func(c *ExpiredMap) {
			c.janitor.stop <- true
		})
	}
	return c
}
