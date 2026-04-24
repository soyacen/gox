package lazyload

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/singleflight"
)

// ErrNilFunction is returned when the New function is nil.
var ErrNilFunction = errors.New("lazyload: New function is nil")

// entry wraps the actual object and its cleanup function.
type entry[Obj any] struct {
	Obj       Obj                                      // The actual stored object
	CloseFunc func(ctx context.Context, obj Obj) error // Object cleanup function
}

// Group implements a concurrency-safe lazy-loading cache ensuring values are created only once per key.
//
// It supports custom creation and cleanup functions.
type Group[Obj any] struct {
	m        sync.Map                                 // Stores key-value pairs where value is of type entry
	g        singleflight.Group                       // Concurrency control to prevent duplicate creation
	New      func(key string) (Obj, error)            // Function to create new values
	Finalize func(ctx context.Context, obj Obj) error // Function to cleanup resources
}

// Option allows customizing LoadOrNew behavior.
type Option[Obj any] func(*options[Obj])

// options stores LoadOrNew configuration options.
type options[Obj any] struct {
	factory   func(key string) (Obj, error)            // Custom creation function
	finalizer func(ctx context.Context, obj Obj) error // Custom cleanup function
}

// WithFactory specifies a custom object creation function.
func WithFactory[Obj any](factory func(key string) (Obj, error)) Option[Obj] {
	return func(opts *options[Obj]) {
		opts.factory = factory
	}
}

// WithFinalizer specifies a custom object cleanup function.
func WithFinalizer[Obj any](finalizer func(ctx context.Context, obj Obj) error) Option[Obj] {
	return func(opts *options[Obj]) {
		opts.finalizer = finalizer
	}
}

// Load gets the value for the given key, creating it using the default New function if it doesn't exist.
//
// Returns:
//   - Obj: the value
//   - error: any error encountered
//   - bool: true if the value already existed, false if it was newly created
func (g *Group[Obj]) Load(key string) (Obj, error, bool) {
	return g.LoadOrNew(key)
}

// LoadOrNew loads an existing value or creates a new one using the specified options.
//
// Parameters:
//   - key: the key to look up
//   - opts: optional configuration options
//
// Returns:
//   - Obj: the value
//   - error: any error encountered
//   - bool: true if the value already existed, false if it was newly created
func (g *Group[Obj]) LoadOrNew(key string, opts ...Option[Obj]) (Obj, error, bool) {
	// Fast path: if value already exists, return it directly
	if value, ok := g.m.Load(key); ok {
		entry := value.(*entry[Obj])
		return entry.Obj, nil, true
	}

	// Apply options
	opt := &options[Obj]{}
	for _, o := range opts {
		o(opt)
	}

	// Determine creation function
	factory := opt.factory
	if factory == nil {
		factory = g.New
	}
	if factory == nil {
		var obj Obj
		return obj, ErrNilFunction, false
	}

	// Use singleflight to prevent duplicate creation
	value, err, _ := g.g.Do(key, func() (any, error) {
		// Double-check to prevent creation by other goroutines during waiting
		if value, ok := g.m.Load(key); ok {
			return value, nil
		}
		value, err := factory(key)
		if err != nil {
			return nil, err
		}

		// Store value with appropriate cleanup function
		finalizer := g.Finalize
		if opt.finalizer != nil {
			finalizer = opt.finalizer
		}
		g.m.Store(key, &entry[Obj]{
			Obj:       value,
			CloseFunc: finalizer,
		})
		return value, nil
	})
	if err != nil {
		var obj Obj
		return obj, err, false
	}

	obj := value.(Obj)
	return obj, nil, false
}

// Close cleans up all cached items and releases resources.
//
// Returns:
//   - error: a combination of all cleanup errors
func (g *Group[Obj]) Close(ctx context.Context) error {
	var errs []error
	// Iterate and cleanup all cached items
	g.m.Range(func(key, value any) bool {
		g.m.Delete(key.(string))
		entry := value.(*entry[Obj])
		if entry.CloseFunc != nil {
			errs = append(errs, entry.CloseFunc(ctx, entry.Obj))
		}
		return true
	})
	return errors.Join(errs...)
}
