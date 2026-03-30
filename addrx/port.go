package addrx

import (
	"net"
)

// PickFreePort automatically selects a free port and returns it.
// It binds to localhost (127.0.0.1) and lets the OS assign an available port.
//
// Returns:
//   - int: An available port number
//   - error: Error if failed to allocate a port
func PickFreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	_, port, err := SplitHostPort(l.Addr())
	return port, err
}
