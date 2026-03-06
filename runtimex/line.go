package runtimex

import (
	"fmt"
	"runtime"
)

// Line returns the current file and line number in the format "file:line".
// It uses runtime.Caller(1) to get the caller's position.
func Line() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "???:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
