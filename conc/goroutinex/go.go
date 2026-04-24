package goroutinex

import (
	"github.com/petermattis/goid"
)

// GoID returns the ID of the current goroutine.
//
// Returns:
//   - int64: The unique identifier of the calling goroutine.
func GoID() int64 {
	return goid.Get()
}
