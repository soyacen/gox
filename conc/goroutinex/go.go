package goroutinex

import (
	"github.com/petermattis/goid"
)

// GoID is used to retrieve the ID of the current Goroutine.
func GoID() int64 {
	return goid.Get()
}
