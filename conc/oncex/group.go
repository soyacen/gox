// Package oncex provides a mechanism to ensure functions are executed only once per key.
// It extends the functionality of sync.Once by allowing keyed execution.
package oncex

import (
	"sync"
)

// Group ensures that for each key, the associated function is executed only once.
// It uses sync.Map to store sync.Once instances keyed by any comparable type.
// It also uses a sync.Pool to reuse sync.Once instances for better performance.
type Group struct {
	// m is a concurrent map that stores sync.Once instances keyed by the provided keys
	m sync.Map
	// p is a pool of sync.Once instances to reduce memory allocations
	p sync.Pool
}

// Do executes the function f only once for the given key.
// If multiple goroutines call Do with the same key, only the first call will execute f,
// and all other calls will wait for the first call to complete.
// key: the unique identifier for the function execution (can be any comparable type)
// f: the function to execute only once
func (o *Group) Do(key any, f func()) {
	// If f is nil, there's nothing to execute, so return early
	if f == nil {
		return
	}

	// Try to get a sync.Once instance from the pool
	v := o.p.Get()

	// If no instance was available from the pool, create a new one
	if v == nil {
		v = &sync.Once{}
	}

	// Attempt to load or store the sync.Once instance for the given key
	// If loaded is true, it means another goroutine has already stored an instance for this key
	actual, loaded := o.m.LoadOrStore(key, v)

	// If another goroutine already stored an instance for this key,
	// put the retrieved sync.Once instance back into the pool for reuse
	if loaded {
		o.p.Put(v)
	}

	// Cast the stored value to *sync.Once
	once := actual.(*sync.Once)

	// Execute the function through sync.Once.Do, which ensures it runs only once
	// Even if multiple goroutines reach this point, sync.Once guarantees f is executed only once
	once.Do(f)
}
