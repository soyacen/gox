package addrx

import (
	"net"
)

// PickFreePort automatically chose a free port and return it
func PickFreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	_, port, err := SplitHostPort(l.Addr())
	return port, err
}
