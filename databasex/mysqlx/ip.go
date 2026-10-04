package mysqlx

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/netip"
)

// ipBinarySize4 and ipBinarySize16 are the lengths of the packed representation
// produced by INET6_ATON for IPv4 and IPv6 addresses.
const (
	ipBinarySize4  = 4
	ipBinarySize16 = 16
)

// parseIPLiteral parses s as an IPv4 or IPv6 address.
//
// Parameters:
//   - s: the textual representation of an address.
//
// Returns:
//   - netip.Addr: the parsed address.
//   - error: non-nil when s is not a valid address, wrapping [ErrInvalidIP].
func parseIPLiteral(s string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("%w: %q", ErrInvalidIP, s)
	}
	return addr, nil
}

// parseIPPrefixLiteral parses s as a CIDR network.
//
// Parameters:
//   - s: the textual representation of a prefix.
//
// Returns:
//   - netip.Prefix: the parsed prefix.
//   - error: non-nil when s is not a valid prefix, wrapping [ErrInvalidIP].
func parseIPPrefixLiteral(s string) (netip.Prefix, error) {
	prefix, err := netip.ParsePrefix(s)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("%w: %q", ErrInvalidIP, s)
	}
	return prefix, nil
}

// parseIPPortLiteral parses s as an address with a port.
//
// Parameters:
//   - s: the textual representation of an address and port.
//
// Returns:
//   - netip.AddrPort: the parsed address and port.
//   - error: non-nil when s is not a valid address and port, wrapping
//     [ErrInvalidIP].
func parseIPPortLiteral(s string) (netip.AddrPort, error) {
	addrPort, err := netip.ParseAddrPort(s)
	if err != nil {
		return netip.AddrPort{}, fmt.Errorf("%w: %q", ErrInvalidIP, s)
	}
	return addrPort, nil
}

// isPrintableASCII reports whether every byte of b is printable ASCII.
//
// Parameters:
//   - b: the bytes to inspect.
//
// Returns:
//   - bool: true when no byte falls outside the range 0x20 to 0x7e.
func isPrintableASCII(b []byte) bool {
	for _, c := range b {
		if c < 0x20 || c > 0x7e {
			return false
		}
	}
	return true
}

// ipFromBytes decodes the packed representation produced by INET6_ATON.
//
// A 16 byte value that is also printable text is decoded as text, because the
// text form of an IPv6 address can be exactly 16 characters long.
//
// Parameters:
//   - b: the packed bytes.
//
// Returns:
//   - netip.Addr: the decoded address.
//   - bool: false when b is neither a 4 byte nor a 16 byte address.
func ipFromBytes(b []byte) (netip.Addr, bool) {
	switch len(b) {
	case ipBinarySize4:
		var raw [4]byte
		copy(raw[:], b)
		return netip.AddrFrom4(raw), true
	case ipBinarySize16:
		if isPrintableASCII(b) {
			if addr, err := netip.ParseAddr(string(b)); err == nil {
				return addr, true
			}
		}

		var raw [16]byte
		copy(raw[:], b)
		return netip.AddrFrom16(raw), true
	default:
		return netip.Addr{}, false
	}
}

// addrBin returns the packed representation of addr as produced by INET6_ATON.
//
// Parameters:
//   - addr: a valid address.
//
// Returns:
//   - []byte: 4 bytes for an IPv4 address, 16 bytes otherwise.
func addrBin(addr netip.Addr) []byte {
	if addr.Is4() {
		b := addr.As4()
		return b[:]
	}

	b := addr.As16()
	return b[:]
}

// ipv4Bytes converts an INET_ATON integer into its big-endian byte form.
//
// Parameters:
//   - v: the unsigned 32 bit integer form of an IPv4 address.
//
// Returns:
//   - [4]byte: the packed address.
func ipv4Bytes(v uint32) [4]byte {
	return [4]byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

// IP represents a MySQL IP address value.
//
// MySQL has no native IP type. Addresses are stored as text in a VARCHAR(45)
// column using the form returned by INET6_NTOA, as packed bytes in a
// VARBINARY(16) column using INET6_ATON, or as an unsigned 32 bit integer in an
// INT UNSIGNED column using INET_ATON. Scan accepts all three representations
// and Value writes the text form.
//
// The address is backed by the standard library [netip.Addr] type.
type IP struct {
	// Addr is the decoded address.
	Addr netip.Addr
	// Valid reports whether the value is not NULL.
	Valid bool
}

// ParseIP parses the string representation of an IPv4 or IPv6 address.
//
// Parameters:
//   - s: the textual representation of an address.
//
// Returns:
//   - IP: the parsed address with Valid set to true.
//   - error: non-nil when s is not a valid address, wrapping [ErrInvalidIP].
func ParseIP(s string) (IP, error) {
	addr, err := parseIPLiteral(s)
	if err != nil {
		return IP{}, err
	}
	return IP{Addr: addr, Valid: true}, nil
}

// MustParseIP is like ParseIP but panics when s is not a valid address.
//
// Parameters:
//   - s: the textual representation of an address.
//
// Returns:
//   - IP: the parsed address with Valid set to true.
func MustParseIP(s string) IP {
	addr, err := ParseIP(s)
	if err != nil {
		panic(err)
	}
	return addr
}

// Scan implements the database/sql.Scanner interface.
//
// A string is parsed as text. A []byte is decoded as the packed INET6_ATON
// representation when it is exactly 4 or 16 bytes long, and parsed as text
// otherwise. An int64 or uint64 is decoded as the INET_ATON value of an IPv4
// address.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an IP, wrapping
//     [ErrInvalidIP], [ErrOutOfRange] or [ErrCannotScan].
func (dst *IP) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = IP{}
		return nil
	case netip.Addr:
		if !v.IsValid() {
			*dst = IP{}
			return nil
		}
		dst.Addr, dst.Valid = v, true
		return nil
	case [4]byte:
		dst.Addr, dst.Valid = netip.AddrFrom4(v), true
		return nil
	case [16]byte:
		dst.Addr, dst.Valid = netip.AddrFrom16(v), true
		return nil
	case string:
		return dst.scanIPText(v)
	case []byte:
		if addr, ok := ipFromBytes(v); ok {
			dst.Addr, dst.Valid = addr, true
			return nil
		}
		return dst.scanIPText(string(v))
	case int64:
		if v < 0 || v > math.MaxUint32 {
			return fmt.Errorf("%w: %d does not fit into an IPv4 address", ErrOutOfRange, v)
		}
		dst.Addr, dst.Valid = netip.AddrFrom4(ipv4Bytes(uint32(v))), true
		return nil
	case uint64:
		if v > math.MaxUint32 {
			return fmt.Errorf("%w: %d does not fit into an IPv4 address", ErrOutOfRange, v)
		}
		dst.Addr, dst.Valid = netip.AddrFrom4(ipv4Bytes(uint32(v))), true
		return nil
	default:
		return cannotScan(src, dst)
	}
}

// scanIPText parses the textual representation s into dst.
//
// Parameters:
//   - s: the textual representation of an address.
//
// Returns:
//   - error: non-nil when s is not a valid address, wrapping [ErrInvalidIP].
func (dst *IP) scanIPText(s string) error {
	addr, err := parseIPLiteral(s)
	if err != nil {
		return err
	}
	dst.Addr, dst.Valid = addr, true
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// The text form written by INET6_NTOA is returned, which is what VARCHAR(45)
// columns and functions such as INET6_ATON expect.
//
// Returns:
//   - driver.Value: the textual address, or nil when the value is NULL.
//   - error: always nil.
func (src IP) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Addr.String(), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the textual address as a JSON string, or null when the value is
//     NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src IP) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Addr.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON string accepted by ParseIP or null.
func (dst *IP) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = IP{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return dst.scanIPText(s)
}

// String returns the textual representation of the address.
//
// Returns:
//   - string: the text form written by INET6_NTOA, or an empty string when the
//     value is NULL.
func (src IP) String() string {
	if !src.Valid {
		return ""
	}
	return src.Addr.String()
}

// Bin returns the packed representation of the address.
//
// The result is the value accepted by a BINARY(16) column and the value produced
// by INET6_ATON.
//
// Returns:
//   - []byte: 4 bytes for an IPv4 address, 16 bytes otherwise, or nil when the
//     value is NULL.
func (src IP) Bin() []byte {
	if !src.Valid {
		return nil
	}
	return addrBin(src.Addr)
}

// Uint32 returns the INET_ATON integer form of an IPv4 address.
//
// The result is the value stored in an INT UNSIGNED column and the value
// produced by INET_ATON.
//
// Returns:
//   - uint32: the big-endian integer form of the address.
//   - bool: false when the value is NULL or the address is not IPv4.
func (src IP) Uint32() (uint32, bool) {
	if !src.Valid || !src.Addr.Is4() {
		return 0, false
	}

	b := src.Addr.As4()
	return binary.BigEndian.Uint32(b[:]), true
}

// IPPrefix represents a MySQL CIDR network value.
//
// MySQL has no native network type, so prefixes are stored as text in a
// VARCHAR(49) column, for example "192.168.1.0/24" or "2001:db8::/32". Scan
// treats every string and []byte source as text, because the shortest prefix
// form ("::/0") is only 4 characters long and would otherwise be mistaken for a
// packed address.
//
// The host bits of the textual form are preserved, matching [netip.ParsePrefix].
// Call [netip.Prefix.Masked] to obtain the network address.
//
// The prefix is backed by the standard library [netip.Prefix] type.
type IPPrefix struct {
	// Prefix is the decoded network.
	Prefix netip.Prefix
	// Valid reports whether the value is not NULL.
	Valid bool
}

// ParseIPPrefix parses the string representation of a CIDR network.
//
// Parameters:
//   - s: the textual representation of a prefix.
//
// Returns:
//   - IPPrefix: the parsed prefix with Valid set to true.
//   - error: non-nil when s is not a valid prefix, wrapping [ErrInvalidIP].
func ParseIPPrefix(s string) (IPPrefix, error) {
	prefix, err := parseIPPrefixLiteral(s)
	if err != nil {
		return IPPrefix{}, err
	}
	return IPPrefix{Prefix: prefix, Valid: true}, nil
}

// MustParseIPPrefix is like ParseIPPrefix but panics when s is invalid.
//
// Parameters:
//   - s: the textual representation of a prefix.
//
// Returns:
//   - IPPrefix: the parsed prefix with Valid set to true.
func MustParseIPPrefix(s string) IPPrefix {
	prefix, err := ParseIPPrefix(s)
	if err != nil {
		panic(err)
	}
	return prefix
}

// Scan implements the database/sql.Scanner interface.
//
// Every string and []byte source is parsed as a CIDR literal.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an IPPrefix, wrapping
//     [ErrInvalidIP] or [ErrCannotScan].
func (dst *IPPrefix) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = IPPrefix{}
		return nil
	case netip.Prefix:
		if !v.IsValid() {
			*dst = IPPrefix{}
			return nil
		}
		dst.Prefix, dst.Valid = v, true
		return nil
	case string:
		return dst.scanIPPrefixText(v)
	case []byte:
		return dst.scanIPPrefixText(string(v))
	default:
		return cannotScan(src, dst)
	}
}

// scanIPPrefixText parses the textual representation s into dst.
//
// Parameters:
//   - s: the textual representation of a prefix.
//
// Returns:
//   - error: non-nil when s is not a valid prefix, wrapping [ErrInvalidIP].
func (dst *IPPrefix) scanIPPrefixText(s string) error {
	prefix, err := parseIPPrefixLiteral(s)
	if err != nil {
		return err
	}
	dst.Prefix, dst.Valid = prefix, true
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the CIDR text, or nil when the value is NULL.
//   - error: always nil.
func (src IPPrefix) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Prefix.String(), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the CIDR text as a JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src IPPrefix) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Prefix.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON string accepted by ParseIPPrefix
//     or null.
func (dst *IPPrefix) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = IPPrefix{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return dst.scanIPPrefixText(s)
}

// String returns the CIDR text of the prefix.
//
// Returns:
//   - string: the text form, or an empty string when the value is NULL.
func (src IPPrefix) String() string {
	if !src.Valid {
		return ""
	}
	return src.Prefix.String()
}

// Bin returns the packed representation of the address part of the prefix.
//
// MySQL has no standard binary encoding for a CIDR network, so the prefix length
// is not included in the result. Store it in a separate column when the network
// boundary must survive a round trip.
//
// Returns:
//   - []byte: 4 bytes when the prefix addresses IPv4, 16 bytes otherwise, or nil
//     when the value is NULL.
func (src IPPrefix) Bin() []byte {
	if !src.Valid {
		return nil
	}
	return addrBin(src.Prefix.Addr())
}

// IPPort represents a MySQL address and port value.
//
// MySQL has no native type for a socket address, so the value is stored as text
// in a VARCHAR(53) column, for example "192.168.1.1:3306" or "[2001:db8::1]:3306".
// Scan treats every string and []byte source as text.
//
// The address and port are backed by the standard library [netip.AddrPort] type.
type IPPort struct {
	// AddrPort is the decoded address and port.
	AddrPort netip.AddrPort
	// Valid reports whether the value is not NULL.
	Valid bool
}

// ParseIPPort parses the string representation of an address and port.
//
// Parameters:
//   - s: the textual representation of an address and port.
//
// Returns:
//   - IPPort: the parsed value with Valid set to true.
//   - error: non-nil when s is not a valid address and port, wrapping
//     [ErrInvalidIP].
func ParseIPPort(s string) (IPPort, error) {
	addrPort, err := parseIPPortLiteral(s)
	if err != nil {
		return IPPort{}, err
	}
	return IPPort{AddrPort: addrPort, Valid: true}, nil
}

// MustParseIPPort is like ParseIPPort but panics when s is invalid.
//
// Parameters:
//   - s: the textual representation of an address and port.
//
// Returns:
//   - IPPort: the parsed value with Valid set to true.
func MustParseIPPort(s string) IPPort {
	addrPort, err := ParseIPPort(s)
	if err != nil {
		panic(err)
	}
	return addrPort
}

// Scan implements the database/sql.Scanner interface.
//
// Every string and []byte source is parsed as an address and port literal.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an IPPort, wrapping
//     [ErrInvalidIP] or [ErrCannotScan].
func (dst *IPPort) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = IPPort{}
		return nil
	case netip.AddrPort:
		if !v.IsValid() {
			*dst = IPPort{}
			return nil
		}
		dst.AddrPort, dst.Valid = v, true
		return nil
	case string:
		return dst.scanIPPortText(v)
	case []byte:
		return dst.scanIPPortText(string(v))
	default:
		return cannotScan(src, dst)
	}
}

// scanIPPortText parses the textual representation s into dst.
//
// Parameters:
//   - s: the textual representation of an address and port.
//
// Returns:
//   - error: non-nil when s is not a valid address and port, wrapping
//     [ErrInvalidIP].
func (dst *IPPort) scanIPPortText(s string) error {
	addrPort, err := parseIPPortLiteral(s)
	if err != nil {
		return err
	}
	dst.AddrPort, dst.Valid = addrPort, true
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the textual address and port, or nil when the value is NULL.
//   - error: always nil.
func (src IPPort) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.AddrPort.String(), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the textual address and port as a JSON string, or null when the
//     value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src IPPort) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.AddrPort.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON string accepted by ParseIPPort or
//     null.
func (dst *IPPort) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = IPPort{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return dst.scanIPPortText(s)
}

// String returns the textual representation of the address and port.
//
// Returns:
//   - string: the text form, or an empty string when the value is NULL.
func (src IPPort) String() string {
	if !src.Valid {
		return ""
	}
	return src.AddrPort.String()
}

// Bin returns the packed representation of the address part of the value.
//
// The port is not included in the result because MySQL has no standard binary
// encoding for an address and port pair.
//
// Returns:
//   - []byte: 4 bytes for an IPv4 address, 16 bytes otherwise, or nil when the
//     value is NULL.
func (src IPPort) Bin() []byte {
	if !src.Valid {
		return nil
	}
	return addrBin(src.AddrPort.Addr())
}
