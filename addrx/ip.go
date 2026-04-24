package addrx

import (
	"errors"
	"fmt"
	"github.com/soyacen/gox/errorx"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Addrs returns a list of unicast interface addresses for all network interfaces.
//
// Returns:
//   - []net.Addr: List of network interface addresses
//   - error: Error if failed to get addresses
func Addrs() ([]net.Addr, error) {
	// Get all network interfaces
	ifaces, err := net.Interfaces()
	// Return error if no network interfaces found
	if err != nil {
		return nil, err
	}
	var res []net.Addr
	var errs []error
	// Iterate over all network interfaces
	for _, iface := range ifaces {
		// Get all addresses for the network interface
		addrs, err := iface.Addrs()
		if err != nil {
			// Record error and continue if failed to get addresses
			errs = append(errs, err)
			continue
		}
		// Add all addresses to result
		res = append(res, addrs...)
	}
	// Return error if no addresses found
	if len(res) <= 0 {
		return nil, errors.Join(errs...)
	}
	// Return all network interface addresses
	return res, nil
}

// IPs returns a list of IP addresses for all network interfaces.
//
// Returns:
//   - []net.IP: List of IP addresses
//   - error: Error if failed to get IPs
func IPs() ([]net.IP, error) {
	// Get all addresses
	addrs, err := Addrs()
	// Return error if no addresses found
	if len(addrs) == 0 {
		return nil, err
	}

	var res []net.IP
	errs := errorx.UnwrapMultiErr(err)

	// Iterate over all addresses
	for _, addr := range addrs {
		// Parse the address
		ip, _, err := SplitHostPort(addr)
		if err != nil {
			// Record error and continue if parsing fails
			errs = append(errs, err)
			continue
		}
		// Add IP to result
		res = append(res, ip)
	}
	if len(res) <= 0 {
		return nil, errors.Join(errs...)
	}
	return res, nil
}

// GlobalUnicastIPs returns a list of global unicast IP addresses for all network interfaces.
//
// Returns:
//   - []net.IP: List of global unicast IP addresses
//   - error: Error if no global unicast IP found
func GlobalUnicastIPs() ([]net.IP, error) {
	ips, err := IPs()
	if len(ips) == 0 {
		return nil, err
	}
	errs := errorx.UnwrapMultiErr(err)
	var res []net.IP
	for _, ip := range ips {
		if IsGlobalUnicastIP(ip) {
			res = append(res, ip)
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	errs = append(errs, errors.New("not found global unicast IP"))
	return nil, errors.Join(errs...)
}

// GlobalUnicastAddr extracts the global unicast IP and port from a network address.
// If the address is not global unicast, it attempts to find one from all interfaces.
//
// Parameters:
//   - address: The network address to extract from
//
// Returns:
//   - net.IP: The global unicast IP
//   - int: The port number
//   - error: Error if extraction fails
func GlobalUnicastAddr(address net.Addr) (net.IP, int, error) {
	ip, port, err := SplitHostPort(address)
	if err != nil {
		return nil, 0, err
	}
	if IsGlobalUnicastIP(ip) {
		return ip, port, nil
	}
	if !ip.IsUnspecified() {
		return nil, 0, errors.New("failed to get global unicast ip")
	}
	ips, err := GlobalUnicastIPs()
	if err != nil {
		return nil, 0, err
	}
	return ips[0], port, err
}

// SplitHostPort splits a network address into IP and port components.
// It handles various address types including IPAddr, IPNet, TCPAddr, and UDPAddr.
//
// Parameters:
//   - addr: Network address to split
//
// Returns:
//   - net.IP: The IP address component
//   - int: The port number
//   - error: Error if failed to parse the address
func SplitHostPort(addr net.Addr) (net.IP, int, error) {
	switch v := addr.(type) {
	case *net.IPAddr:
		return v.IP, 0, nil
	case *net.IPNet:
		return v.IP, 0, nil
	case *net.TCPAddr:
		return v.IP, v.Port, nil
	case *net.UDPAddr:
		return v.IP, v.Port, nil
	default:
		host, port, err := net.SplitHostPort(addr.String())
		if err != nil {
			return net.IP{}, 0, err
		}
		portNum, err := strconv.Atoi(port)
		return net.ParseIP(host), portNum, err
	}
}

// IsGlobalUnicastIP checks whether the IP is a global unicast IP address.
// A global unicast IP is routable on the internet and not a private, loopback, or multicast address.
//
// Parameters:
//   - ip: The IP address to check
//
// Returns:
//   - bool: True if the IP is a global unicast address, false otherwise
func IsGlobalUnicastIP(ip net.IP) bool {
	if ip.IsUnspecified() {
		// 这个方法用于检查给定的 IP 地址是否是未指定的地址。
		// 未指定的地址: 在 IP 地址中，未指定的地址表示没有特定的网络接口或地址。具体来说：
		// IPv4: "0.0.0.0"。
		// IPv6: "::"。
		return false
	}
	if ip.IsLoopback() {
		// 这个方法用于检查给定的 IP 地址是否是回环地址。
		// 回环地址（Loopback Address）是网络中的一种特殊IP地址，主要用于测试和本地通信。
		// IPv4: "127.0.0.0/8"，最常用的回环地址是 127.0.0.1。
		// IPv6: "::1"。
		return false
	}
	if ip.IsMulticast() {
		// 这个方法用于检查给定的 IP 地址是多播地址。
		// 多播地址用于向一组主机发送数据包，而不是单个主机。
		// IPv4: 224.0.0.0 到 239.255.255.255。
		// IPv6: ff00::/8
		return false
	}
	if ip.IsLinkLocalMulticast() {
		// 这个方法用于检查给定的 IP 地址是链路本地多播地址。
		// 链路本地多播地址用于在同一网络段内的设备之间进行通信。
		// IPv4: 224.0.0.0 到 224.0.0.255。
		// IPv6: ff02::/16
		return false
	}
	if ip.IsInterfaceLocalMulticast() {
		// 这个方法用于检查给定的 IP 地址是接口本地多播地址。
		// 接口本地多播地址用于在同一网络接口内的设备之间进行通信。
		// IPv4: 224.0.0.0 到 224.0.0.255（具体到 224.0.0.252 和 224.0.0.253）。
		// IPv6: ff01::/16
		return false
	}
	if ip.IsLinkLocalUnicast() {
		// 这个方法用于检查给定的 IP 地址是链路本地单播地址。
		// 链路本地单播地址用于在同一网络段内的设备之间进行单点通信。
		// IPv4: 169.254.0.0/16。
		// IPv6: fe80::/10
		return false
	}
	// 其他情况，认为是全局单播地址
	// 全局单播地址用于在互联网上进行通信，而不是在本地网络或特定的网络段内。
	return ip.IsGlobalUnicast()
}

// GlobalUnicastIPString returns a global unicast IP address as a string.
//
// Returns:
//   - string: The global unicast IP address
//   - error: Error if no global unicast IP is found
func GlobalUnicastIPString() (string, error) {
	ips, err := GlobalUnicastIPs()
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

// InterfaceIPs returns public IP addresses for the specified network interface name.
//
// Parameters:
//   - name: The name of the network interface
//
// Returns:
//   - []net.IP: List of IP addresses for the interface
//   - error: Error if the interface is not found
func InterfaceIPs(name string) ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		if iface.Name != name {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := SplitHostPort(addr)
			if err != nil {
				continue
			}
			ips = append(ips, ip)
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("not found the ip of interface %s", name)
	}
	return ips, nil
}

// InterfaceIPv4 returns public IPv4 addresses for the specified network interface name.
//
// Parameters:
//   - name: The name of the network interface
//
// Returns:
//   - []net.IP: List of IPv4 addresses for the interface
//   - error: Error if the interface is not found or has no IPv4 addresses
func InterfaceIPv4(name string) ([]net.IP, error) {
	ips, err := InterfaceIPs(name)
	if err != nil {
		return nil, err
	}
	var r []net.IP
	for _, ip := range ips {
		ip = ip.To4()
		if len(ip) == 0 {
			continue
		}
		r = append(r, ip)
	}
	return r, nil
}

// IsLocalIPAddr checks whether the given IP address string is a local (private) address.
//
// Parameters:
//   - ip: The IP address string to check
//
// Returns:
//   - bool: True if the IP is a local address, false otherwise
func IsLocalIPAddr(ip string) bool {
	return IsLocalIP(net.ParseIP(ip))
}

// IsLocalIP checks whether the given IP address is a local (private) address.
// It checks for loopback, private, and link-local addresses.
//
// Parameters:
//   - ip: The IP address to check
//
// Returns:
//   - bool: True if the IP is a local address, false otherwise
func IsLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	if ip.IsPrivate() {
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return (ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// ClientIP returns the client IP address from the HTTP request.
// It checks X-Forwarded-For and X-Real-IP headers first, then falls back to RemoteAddr.
// This supports reverse proxies such as nginx or haproxy.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The client IP address
func ClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientPublicIP returns the client's public IP address from the HTTP request.
// It checks X-Forwarded-For and X-Real-IP headers, filtering out local addresses.
// This supports reverse proxies such as nginx or haproxy.
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip = strings.TrimSpace(ip); ip != "" && !IsLocalIPAddr(ip) {
			return ip
		}
	}

	if ip = strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" && !IsLocalIPAddr(ip) {
		return ip
	}

	if ip = RemoteIP(r); !IsLocalIPAddr(ip) {
		return ip
	}

	return ""
}

// RemoteIP returns the IP address from the request's RemoteAddr.
// It performs a quick parse of the remote address.
//
// Parameters:
//   - r: The HTTP request
//
// Returns:
//   - string: The remote IP address
func RemoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// IPString2Long converts an IPv4 string to a numeric value.
//
// Parameters:
//   - ip: The IPv4 address string
//
// Returns:
//   - uint: The numeric representation of the IP
//   - error: Error if the IP format is invalid
func IPString2Long(ip string) (uint, error) {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IPString converts a numeric value to an IPv4 string.
//
// Parameters:
//   - i: The numeric value to convert
//
// Returns:
//   - string: The IPv4 address string
//   - error: Error if the value exceeds the IPv4 range
func Long2IPString(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}

// IP2Long converts a net.IP to a numeric value.
//
// Parameters:
//   - ip: The IP address to convert
//
// Returns:
//   - uint: The numeric representation of the IP
//   - error: Error if the IP format is invalid
func IP2Long(ip net.IP) (uint, error) {
	b := ip.To4()
	if b == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(b[3]) | uint(b[2])<<8 | uint(b[1])<<16 | uint(b[0])<<24, nil
}

// Long2IP converts a numeric value to a net.IP.
//
// Parameters:
//   - i: The numeric value to convert
//
// Returns:
//   - net.IP: The IP address
//   - error: Error if the value exceeds the IPv4 range
func Long2IP(i uint) (net.IP, error) {
	if i > math.MaxUint32 {
		return nil, errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip, nil
}
