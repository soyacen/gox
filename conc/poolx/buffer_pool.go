package poolx

import (
	"bytes"
	"errors"
	"math/bits"
)

// ErrInvalidSize is returned when the maximum size is less than the minimum size in NewBucketBufferPool
var ErrInvalidSize = errors.New("poolx: maxSize can't be less than minSize")

// BufferPool is a type alias for Pool[*bytes.Buffer], representing a pool of bytes.Buffer objects
type BufferPool = Pool[*bytes.Buffer]

// NewBufferPool creates a new BufferPool instance with the specified buffer size.
// The pool creates new bytes.Buffer instances with the given size and resets them before returning to the pool.
// size: the initial capacity of each bytes.Buffer in the pool
// Returns a new BufferPool instance and an error if creation fails
func NewBufferPool(size int) (*BufferPool, error) {
	return NewPool(
		func() *bytes.Buffer {
			return bytes.NewBuffer(make([]byte, 0, size))
		},
		func(b *bytes.Buffer) {
			b.Reset() // Reset clears the buffer content but retains the underlying byte slice
		},
	)
}

// BucketBufferPool is the main structure for bucket-based memory pool.
// It manages multiple buffer pools of different sizes to efficiently handle buffers of various capacities.
type BucketBufferPool struct {
	minSize int           // Minimum buffer size supported by this pool
	maxSize int           // Maximum buffer size supported by this pool
	buckets []*BufferPool // Array of buffer pools for different sizes
}

// NewBucketBufferPool creates a new bucket-based memory pool with the specified size range.
// It creates a series of buffer pools with sizes ranging from minSize to maxSize,
// where each pool's size is double the previous one.
// minSize: the minimum buffer size supported by this pool
// maxSize: the maximum buffer size supported by this pool
// Returns a new BucketBufferPool instance and an error if maxSize is less than minSize
func NewBucketBufferPool(minSize, maxSize int) (*BucketBufferPool, error) {
	// Validate that maxSize is not less than minSize
	if maxSize < minSize {
		return nil, ErrInvalidSize
	}

	// Multiplier factor for each bucket size (doubles each time)
	multiplier := 2
	var buckets []*BufferPool
	curSize := minSize

	// Create a series of buffer pools from minSize to maxSize,
	// where each pool's size is multiplier times the previous one
	for curSize < maxSize {
		p, err := NewBufferPool(curSize)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, p)
		curSize *= multiplier
	}

	// Create and add the final pool for the maximum size
	p, err := NewBufferPool(curSize)
	if err != nil {
		return nil, err
	}
	buckets = append(buckets, p)

	// Return the newly created BucketBufferPool
	return &BucketBufferPool{
		minSize: minSize,
		maxSize: maxSize,
		buckets: buckets,
	}, nil
}

// Get retrieves a suitable bytes.Buffer based on the requested size.
// If the requested size exceeds maxSize, a new buffer is directly created without pooling.
// size: the requested buffer size
// Returns a bytes.Buffer with at least the requested capacity
func (p *BucketBufferPool) Get(size int) *bytes.Buffer {
	sp := p.findPool(size) // Find the appropriate pool
	if sp == nil {
		// If no suitable pool is found (requested size exceeds maxSize),
		// directly create a new buffer
		return bytes.NewBuffer(make([]byte, size))
	}
	return sp.Get() // Get buffer from the found pool
}

// Put returns a bytes.Buffer to the appropriate pool based on its capacity.
// Buffers with capacity exceeding maxSize are not returned to any pool.
// b: the bytes.Buffer to return to the pool
func (p *BucketBufferPool) Put(b *bytes.Buffer) {
	sp := p.findPool(b.Cap()) // Find the pool based on buffer capacity
	if sp == nil {
		return // Do not return to pool if capacity exceeds maxSize
	}
	sp.Put(b) // Return buffer to the pool
}

// findPool finds the most suitable buffer pool based on the given size.
// It uses bit operations to efficiently determine which bucket should be used.
// size: the size for which to find an appropriate pool
// Returns the appropriate buffer pool or nil if size exceeds maxSize
func (p *BucketBufferPool) findPool(size int) *BufferPool {
	if size > p.maxSize {
		return nil // Return nil if requested size exceeds maximum
	}

	// Calculate quotient and remainder of size divided by minSize
	div, rem := bits.Div64(0, uint64(size), uint64(p.minSize))

	// Calculate the number of bits in the binary representation of quotient
	// to determine the pool index
	idx := bits.Len64(div)

	// Adjust index if division is exact and quotient is a power of 2
	if rem == 0 && div != 0 && (div&(div-1)) == 0 {
		idx = idx - 1
	}
	return p.buckets[idx]
}
