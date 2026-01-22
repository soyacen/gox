package addrx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicIP(t *testing.T) {
	ips, err := GlobalUnicastIPs()
	assert.Nil(t, err)
	t.Log(ips)
	//assert.Equal(t, "172.16.40.45", ips[0])
}

func TestInterfaceIP(t *testing.T) {
	ip, err := InterfaceIPs("en7")
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Log(ip)
}
