// Package poolx provides a generic object pool implementation that extends the standard library's sync.Pool.
// It allows for efficient reuse of objects with custom creation and reset functions.
package poolx

import (
	"errors"
	"sync"
)

// Package-level error variables for pool creation validation
var (
	// ErrNilNewFunction is returned when a nil new function is passed to NewPool
	ErrNilNewFunction = errors.New("poolx: new function is nil")
	// ErrNilResetFunction is returned when a nil reset function is passed to NewPool
	ErrNilResetFunction = errors.New("poolx: reset function is nil")
)

// Pool is a generic wrapper around sync.Pool that provides type safety and custom reset functionality.
// It allows for efficient reuse of objects of type T, reducing garbage collection pressure.
// T: the type of objects to be pooled
type Pool[T any] struct {
	// p is the underlying sync.Pool that manages the object storage
	p sync.Pool
	// reset is a function that cleans up an object before it's returned to the pool
	reset func(o T)
}

// NewPool creates a new Pool instance with the specified new and reset functions.
// The new function is used to create new instances of T when the pool is empty.
// The reset function is called on objects before they are put back into the pool.
// newFunc: function that creates new instances of T
// resetFunc: function that resets/cleans up instances of T before pooling
// Returns a new Pool instance and an error if either function is nil
func NewPool[T any](newFunc func() T, resetFunc func(o T)) (*Pool[T], error) {
	// Validate that the new function is not nil
	if newFunc == nil {
		return nil, ErrNilNewFunction
	}

	// Validate that the reset function is not nil
	if resetFunc == nil {
		return nil, ErrNilResetFunction
	}

	// Create and return the new Pool
	return &Pool[T]{
		p: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
		reset: resetFunc,
	}, nil
}

// Put returns an object to the pool after calling the reset function on it.
// This method should be called when an object is no longer needed so it can be reused.
// o: the object to return to the pool
func (p *Pool[T]) Put(o T) {
	// Call the reset function to clean up the object before pooling
	p.reset(o)
	// Put the object into the underlying sync.Pool
	p.p.Put(o)
}

// Get retrieves an object from the pool.
// If the pool is empty, the new function provided to NewPool will be called to create a new instance.
// Returns an object of type T from the pool or a newly created instance
func (p *Pool[T]) Get() T {
	// Get an object from the underlying sync.Pool and cast it to type T
	return p.p.Get().(T)
}
