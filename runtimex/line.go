package runtimex

import (
	"fmt"
	"runtime"
)

// Line returns the current file and line number in the format "file:line".
// It uses runtime.Caller(1) to get the caller's position if no argument is provided.
// If skip is provided, it uses runtime.Caller(skip[0]).
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
