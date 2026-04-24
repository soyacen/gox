package stringx

import (
	"strings"
	"sync"
)

const (
	// DefaultBuilderSize is the default size of the Builder in bytes.
	DefaultBuilderSize = 1024

	// MaxBuilderSize is the maximum size of the Builder in bytes.
	MaxBuilderSize = 4096
)

// defaultBuilderPool is a pool of strings.Builder.
var defaultBuilderPool = sync.Pool{
	New: func() any {
		return &strings.Builder{}
	},
}

// GetBuilder retrieves a strings.Builder from the pool.
//
// Returns:
//   - *strings.Builder: a Builder from the pool.
func GetBuilder() *strings.Builder {
	return defaultBuilderPool.Get().(*strings.Builder)
}

// PutBuilder returns a strings.Builder to the pool.
//
// Parameters:
//   - buf: the Builder to return.
func PutBuilder(buf *strings.Builder) {
	if buf.Cap() > MaxBuilderSize {
		return
	}
	buf.Reset()
	defaultBuilderPool.Put(buf)
}
