package addrx

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicIP(t *testing.T) {
	ips, err := GlobalUnicastIPs()
	assert.Nil(t, err)
	t.Log(ips)
}

// pickInterfaceWithAddrs returns the name of the first network interface that
// has at least one assigned address. Returns an empty string if none is
// available, so callers can skip on environments without usable interfaces.
func pickInterfaceWithAddrs(t *testing.T) string {
	t.Helper()
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("net.Interfaces failed: %v", err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}
		return iface.Name
	}
	return ""
}

func TestInterfaceIP(t *testing.T) {
	name := pickInterfaceWithAddrs(t)
	if name == "" {
		t.Skip("no network interface with addresses available")
	}
	ips, err := InterfaceIPs(name)
	assert.NoError(t, err)
	assert.NotEmpty(t, ips)
	t.Logf("interface %q ips=%v", name, ips)
}

func TestInterfaceIP_NotFound(t *testing.T) {
	_, err := InterfaceIPs("definitely-not-an-interface-xyz")
	assert.Error(t, err)
}
