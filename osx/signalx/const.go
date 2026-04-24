package signalx

import (
	"os"
	"syscall"
)

// ShutdownSignals returns a slice of signals that are typically used to request process shutdown.
// The returned signals include SIGTERM, SIGINT, SIGQUIT, and SIGKILL.
//
// Returns:
//   - []os.Signal: A slice of shutdown signals.
func ShutdownSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
	}
}
