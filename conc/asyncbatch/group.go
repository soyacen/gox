// Package asyncbatch provides asynchronous batch processing functionality
// that can process objects either by batch size or time interval.
package asyncbatch

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// Package-level error variables
var (
	// ErrTaskInvalid is returned when the processing task function is nil
	ErrTaskInvalid = errors.New("asyncbatch: task is nil")
	// ErrClosed is returned when trying to operate on a closed Group
	ErrClosed = errors.New("asyncbatch: group is closed")
)

// options holds configuration options for the Group
type options struct {
	// Size is the threshold batch size that triggers processing
	Size int
	// Interval is the time interval that triggers processing
	Interval time.Duration
	// Recover is the panic recovery handler function
	Recover func(p any, stack []byte)
}

// Option is a function that configures options
type Option func(*options)

// Size returns an Option that sets the batch size threshold
func Size(size int) Option {
	return func(o *options) {
		o.Size = size
	}
}

// Interval returns an Option that sets the time interval for processing
func Interval(interval time.Duration) Option {
	return func(o *options) {
		o.Interval = interval
	}
}

// Recover returns an Option that sets the panic recovery handler
func Recover(f func(p any, stack []byte)) Option {
	return func(o *options) {
		o.Recover = f
	}
}

// Apply applies the given options to the options struct
func (o *options) Apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Correct validates and corrects the options with default values if needed
func (o *options) Correct() *options {
	// Set default size if not specified or invalid
	if o.Size <= 0 {
		o.Size = 64
	}
	// Set default interval if not specified or invalid
	if o.Interval <= 0 {
		o.Interval = 128 * time.Millisecond
	}
	// Set default recovery handler if not specified
	if o.Recover == nil {
		// Default error handler that prints panic info and stack trace
		o.Recover = func(p any, stack []byte) {
			fmt.Printf("asyncbatch: panic trigger, %v, stack: %s", p, stack)
		}
	}
	return o
}

// Group is a generic struct for asynchronously batch processing objects of type Obj
type Group[Obj any] struct {
	// options contains the configuration options for this Group
	options *options
	// mu mutex protects shared resources like buf
	mu sync.Mutex
	// buf stores objects waiting to be processed
	buf []Obj
	// submitCh is a signal channel for notifying new batches are ready
	submitCh chan struct{}
	// closed atomic boolean flag indicating whether the Group is closed
	closed atomic.Bool
	// closedCh is a notification channel for signaling the loop to exit
	closedCh chan struct{}
	// wg wait group for waiting the loop goroutine to finish
	wg sync.WaitGroup
	// task is the actual function that processes batches of objects
	task func(objs []Obj)
}

// New creates a new Group instance
// task: the actual function that processes batches of objects
// opts: optional configuration functions
// Returns the created Group instance and possible error
func New[Obj any](task func(objs []Obj), opts ...Option) (*Group[Obj], error) {
	// Validate task function
	if task == nil {
		return nil, ErrTaskInvalid
	}

	// Apply and correct options
	opt := new(options).Apply(opts...).Correct()

	// Create Group instance
	g := &Group[Obj]{
		mu:       sync.Mutex{},
		buf:      make([]Obj, 0, opt.Size),
		submitCh: make(chan struct{}, 1),
		closed:   atomic.Bool{},
		closedCh: make(chan struct{}),
		wg:       sync.WaitGroup{},
		options:  opt,
		// Wrap task function with panic recovery mechanism
		task: func(objs []Obj) {
			defer func() {
				// Catch panics in task execution and call recovery handler
				if p := recover(); p != nil {
					opt.Recover(p, debug.Stack())
				}
			}()
			task(objs)
		},
	}

	// Start the processing loop goroutine
	g.wg.Add(1)
	go g.loop()
	return g, nil
}

// Submit submits an object to the Group
// obj: the object to submit
// Returns nil on success, ErrClosed if the Group is closed
func (g *Group[Obj]) Submit(obj Obj) error {
	// Check if Group is closed (fast path)
	if g.closed.Load() {
		return ErrClosed
	}

	// Lock to protect shared resources
	g.mu.Lock()
	// Double-check if Group is closed
	if g.closed.Load() {
		g.mu.Unlock()
		return ErrClosed
	}

	// Add object to buffer
	g.buf = append(g.buf, obj)

	// Send submit signal if buffer reaches threshold
	if len(g.buf) >= g.options.Size {
		g.mu.Unlock()
		// Non-blocking send signal
		select {
		case g.submitCh <- struct{}{}:
		default:
		}
	} else {
		g.mu.Unlock()
	}
	return nil
}

// Close closes the Group, processes remaining objects and waits for completion
// Returns nil on successful close, ErrClosed if already closed
func (g *Group[Obj]) Close() error {
	// Check if Group is closed (fast path)
	if g.closed.Load() {
		return ErrClosed
	}

	// Lock to protect shared resources
	g.mu.Lock()
	// Double-check if Group is closed
	if g.closed.Load() {
		g.mu.Unlock()
		return ErrClosed
	}

	// Mark Group as closed
	g.closed.Store(true)
	g.mu.Unlock()

	// Notify loop to exit
	close(g.closedCh)
	// Wait for loop to finish
	g.wg.Wait()
	return nil
}

// loop is the main processing loop running in a separate goroutine
// It listens for submit signals, timer signals, and close signals
func (g *Group[Obj]) loop() {
	// Mark wait group as done when function exits
	defer g.wg.Done()

	// Create ticker
	ticker := time.NewTicker(g.options.Interval)
	defer ticker.Stop()

	// Infinite loop listening for various signals
	for {
		select {
		// Process on submit signal
		case <-g.submitCh:
			g.onSubmit()
		// Process on timer tick
		case <-ticker.C:
			g.onTick()
		// Process remaining objects and exit on close signal
		case <-g.closedCh:
			g.onClose()
			return
		}
	}
}

// onSubmit handles submit signals and processes batches when buffer reaches size threshold
func (g *Group[Obj]) onSubmit() {
	// Lock to protect shared resources
	g.mu.Lock()
	// Return if buffer hasn't reached size threshold
	if len(g.buf) < g.options.Size {
		g.mu.Unlock()
		return
	}

	// Extract a complete batch
	batch := g.buf[0:g.options.Size]
	// Remove processed objects
	buf := make([]Obj, len(g.buf)-g.options.Size)
	copy(buf, g.buf[g.options.Size:])
	g.buf = buf
	g.mu.Unlock()

	// Execute task
	g.task(batch)
}

// onTick handles timer signals and processes batches at time intervals
func (g *Group[Obj]) onTick() {
	// Lock to protect shared resources
	g.mu.Lock()
	// Return if buffer is empty
	if len(g.buf) <= 0 {
		g.mu.Unlock()
		return
	}

	var batch []Obj
	// If buffer size exceeds threshold, extract a complete batch
	if len(g.buf) >= g.options.Size {
		batch = g.buf[0:g.options.Size]
		buf := make([]Obj, len(g.buf)-g.options.Size)
		copy(buf, g.buf[g.options.Size:])
		g.buf = buf
	} else {
		// Otherwise process all remaining objects
		batch = g.buf
		g.buf = make([]Obj, 0, g.options.Size)
	}

	g.mu.Unlock()
	// Execute task
	g.task(batch)
}

// onClose handles close signals and processes all remaining objects
func (g *Group[Obj]) onClose() {
	// Lock to protect shared resources
	g.mu.Lock()
	// Chunk remaining objects for processing
	batches := chunk(g.buf, g.options.Size)
	g.buf = nil
	g.mu.Unlock()

	// Process each batch
	for _, batch := range batches {
		g.task(batch)
	}
}

// chunk splits a slice into chunks of specified size
// s: the slice to split
// size: the size of each chunk
// Returns array of slice chunks
func chunk[S ~[]E, E any](s S, size int) []S {
	l := len(s)
	// Pre-allocate capacity
	ss2 := make([]S, 0, (l+size)/size)

	// Split slice by size
	for i := 0; i < l; i += size {
		if i+size < l {
			ss2 = append(ss2, s[i:i+size])
		} else {
			ss2 = append(ss2, s[i:l])
		}
	}
	return ss2
}
