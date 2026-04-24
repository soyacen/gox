package ringx

// defaultRingSize defines the default ring size if not specified.
const defaultRingSize = 16

// RingGrowingOptions sets parameters for RingGrowing and TypedRingGrowing.
type RingGrowingOptions struct {
	// InitialSize is the number of pre-allocated elements in the
	// initial underlying storage buffer.
	InitialSize int
}

// TypedRingGrowing is a growing ring buffer.
// The zero value has an initial size of 0 and is ready to use.
// Not thread-safe.
type TypedRingGrowing[T any] struct {
	data     []T
	n        int // Size of Data
	beg      int // First available element
	readable int // Number of data items available
}

// NewTypedRingGrowing constructs a new TypedRingGrowing instance with provided options.
//
// Parameters:
//   - opts: The options for the ring buffer
//
// Returns:
//   - *TypedRingGrowing[T]: A new TypedRingGrowing instance
func NewTypedRingGrowing[T any](opts RingGrowingOptions) *TypedRingGrowing[T] {
	return &TypedRingGrowing[T]{
		data: make([]T, opts.InitialSize),
		n:    opts.InitialSize,
	}
}

// ReadOne reads (consumes) the first item from the buffer if it is available.
//
// Returns:
//   - T: The read item (zero value if not available)
//   - bool: True if an item was available, false otherwise
func (r *TypedRingGrowing[T]) ReadOne() (data T, ok bool) {
	if r.readable == 0 {
		return
	}
	r.readable--
	element := r.data[r.beg]
	var zero T
	r.data[r.beg] = zero // Remove reference to the object to help GC
	if r.beg == r.n-1 {
		// Was the last element
		r.beg = 0
	} else {
		r.beg++
	}
	return element, true
}

// WriteOne adds an item to the end of the buffer, growing it if it is full.
//
// Parameters:
//   - data: The item to add
func (r *TypedRingGrowing[T]) WriteOne(data T) {
	if r.readable == r.n {
		// Time to grow
		newN := r.n * 2
		if newN == 0 {
			newN = defaultRingSize
		}
		newData := make([]T, newN)
		to := r.beg + r.readable
		if to <= r.n {
			copy(newData, r.data[r.beg:to])
		} else {
			copied := copy(newData, r.data[r.beg:])
			copy(newData[copied:], r.data[:(to%r.n)])
		}
		r.beg = 0
		r.data = newData
		r.n = newN
	}
	r.data[(r.readable+r.beg)%r.n] = data
	r.readable++
}

// Len returns the number of items in the buffer.
//
// Returns:
//   - int: The number of items
func (r *TypedRingGrowing[T]) Len() int {
	return r.readable
}

// Cap returns the capacity of the buffer.
//
// Returns:
//   - int: The capacity
func (r *TypedRingGrowing[T]) Cap() int {
	return r.n
}

// RingOptions sets parameters for Ring.
type RingOptions struct {
	// InitialSize is the number of pre-allocated elements in the
	// initial underlying storage buffer.
	InitialSize int
	// NormalSize is the number of elements to allocate for new storage
	// buffers once the Ring is consumed and
	// can shrink again.
	NormalSize int
}

// Ring is a dynamically-sized ring buffer which can grow and shrink as needed.
// The zero value has an initial size and normal size of 0 and is ready to use.
// Not thread-safe.
type Ring[T any] struct {
	growing    TypedRingGrowing[T]
	normalSize int // Limits the size of the buffer that is kept for reuse. Read-only.
}

// NewRing constructs a new Ring instance with provided options.
//
// Parameters:
//   - opts: The options for the ring buffer
//
// Returns:
//   - *Ring[T]: A new Ring instance
func NewRing[T any](opts RingOptions) *Ring[T] {
	return &Ring[T]{
		growing:    *NewTypedRingGrowing[T](RingGrowingOptions{InitialSize: opts.InitialSize}),
		normalSize: opts.NormalSize,
	}
}

// ReadOne reads (consumes) the first item from the buffer if it is available.
// When the buffer has been totally consumed and has grown beyond its normal size, it shrinks.
//
// Returns:
//   - T: The read item (zero value if not available)
//   - bool: True if an item was available, false otherwise
func (r *Ring[T]) ReadOne() (data T, ok bool) {
	element, ok := r.growing.ReadOne()

	if r.growing.readable == 0 && r.growing.n > r.normalSize {
		// The buffer is empty. Reallocate a new buffer so the old one can be
		// garbage collected.
		r.growing.data = make([]T, r.normalSize)
		r.growing.n = r.normalSize
		r.growing.beg = 0
	}

	return element, ok
}

// WriteOne adds an item to the end of the buffer, growing it if it is full.
//
// Parameters:
//   - data: The item to add
func (r *Ring[T]) WriteOne(data T) {
	r.growing.WriteOne(data)
}

// Len returns the number of items in the buffer.
//
// Returns:
//   - int: The number of items
func (r *Ring[T]) Len() int {
	return r.growing.Len()
}

// Cap returns the capacity of the buffer.
//
// Returns:
//   - int: The capacity
func (r *Ring[T]) Cap() int {
	return r.growing.Cap()
}
