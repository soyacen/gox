package runtimex

import (
	"fmt"
	"runtime"
)

// Line returns the current file and line number in the format "file:line".
// It uses runtime.Caller to get the caller's position.
//
// Parameters:
//   - skip: Optional number of stack frames to skip (default is 1)
//
// Returns:
//   - string: The file and line number in "file:line" format
func Line(skip ...int) string {
	s := 1
	if len(skip) > 0 {
		s = skip[0]
	}
	_, file, line, ok := runtime.Caller(s)
	if !ok {
		return "???:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
