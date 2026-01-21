package mapx

import (
	"runtime"
	"sync"
	"time"
)

type ExpiredMap struct {
	items map[any]expiredMapItem
	mu    sync.RWMutex

	expireAfter func(key any) time.Duration
	onEvicted   func(any, any)
	janitor     *janitor
}

func (c *ExpiredMap) Load(key any) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expired() {
		return nil, false
	}
	return item.Object, true
}

func (c *ExpiredMap) Store(key, value any) {
	expiration := c.expiration(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = expiredMapItem{
		Object:     value,
		Expiration: expiration,
	}
}

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

func (c *ExpiredMap) Delete(key any) {
	_, _ = c.LoadAndDelete(key)
}

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

func (c *ExpiredMap) expiration(key any) int64 {
	var expiration int64
	expireAfter := c.expireAfter(key)
	if expireAfter > 0 {
		expiration = time.Now().Add(expireAfter).UnixNano()
	}
	return expiration
}

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

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

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

type expiredMapItem struct {
	Object     any
	Expiration int64
}

// Expired Returns true if the item has expired.
func (item expiredMapItem) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type ExpiredMapOption func(*ExpiredMap)

func ExpireAfter(f func(key any) time.Duration) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.expireAfter = f
	}
}

func CleanupInterval(interval time.Duration) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.janitor.Interval = interval
	}
}

func OnEvicted(f func(any, any)) ExpiredMapOption {
	return func(expiredMap *ExpiredMap) {
		expiredMap.onEvicted = f
	}
}

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
