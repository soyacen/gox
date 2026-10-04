package mysqlx

import (
	"bytes"
	"errors"
	"math"
	"net/netip"
	"testing"
)

func TestIPScan(t *testing.T) {
	packedV4 := []byte{192, 168, 1, 1}
	v6 := MustParseIP("2001:db8:85a3::1")
	packedV6 := v6.Addr.As16()

	testCases := []struct {
		name      string
		src       any
		want      string
		wantValid bool
		wantErr   error
	}{
		{name: "ipv4 text", src: "192.168.1.1", want: "192.168.1.1", wantValid: true},
		{name: "ipv4 text bytes", src: []byte("192.168.1.1"), want: "192.168.1.1", wantValid: true},
		{name: "ipv4 packed", src: packedV4, want: "192.168.1.1", wantValid: true},
		{name: "ipv4 array", src: [4]byte{192, 168, 1, 1}, want: "192.168.1.1", wantValid: true},
		{name: "ipv4 integer", src: int64(3232235777), want: "192.168.1.1", wantValid: true},
		{name: "ipv4 unsigned", src: uint64(3232235777), want: "192.168.1.1", wantValid: true},
		{name: "ipv4 zero integer", src: int64(0), want: "0.0.0.0", wantValid: true},
		{name: "ipv6 text", src: "2001:db8::1", want: "2001:db8::1", wantValid: true},
		{name: "ipv6 packed", src: packedV6[:], want: "2001:db8:85a3::1", wantValid: true},
		{name: "ipv6 array", src: packedV6, want: "2001:db8:85a3::1", wantValid: true},
		// The text form of this address is exactly 16 characters long, so a
		// printable 16 byte source must be decoded as text, not as packed bytes.
		{name: "sixteen char text", src: []byte("2001:db8:85a3::1"), want: "2001:db8:85a3::1", wantValid: true},
		{name: "zoned", src: "fe80::1%eth0", want: "fe80::1%eth0", wantValid: true},
		{name: "mapped", src: "::ffff:192.168.1.1", want: "::ffff:192.168.1.1", wantValid: true},
		{name: "addr", src: netip.MustParseAddr("10.0.0.1"), want: "10.0.0.1", wantValid: true},
		{name: "nil", src: nil},
		{name: "zero addr", src: netip.Addr{}},
		{name: "invalid text", src: "not-an-ip", wantErr: ErrInvalidIP},
		{name: "empty text", src: "", wantErr: ErrInvalidIP},
		{name: "prefix text", src: "192.168.1.0/24", wantErr: ErrInvalidIP},
		{name: "negative integer", src: int64(-1), wantErr: ErrOutOfRange},
		{name: "integer overflow", src: uint64(math.MaxUint32) + 1, wantErr: ErrOutOfRange},
		{name: "unsupported", src: 1.5, wantErr: ErrCannotScan},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v IP
			err := v.Scan(tc.src)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Scan(%#v) error = %v, want %v", tc.src, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan(%#v) unexpected error: %v", tc.src, err)
			}
			if v.Valid != tc.wantValid {
				t.Fatalf("Scan(%#v).Valid = %v, want %v", tc.src, v.Valid, tc.wantValid)
			}
			if got := v.String(); got != tc.want {
				t.Fatalf("Scan(%#v).String() = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

func TestIPValueAndHelpers(t *testing.T) {
	v4 := MustParseIP("192.168.1.1")

	value, err := v4.Value()
	if err != nil || value != "192.168.1.1" {
		t.Fatalf("Value() = %#v, %v", value, err)
	}

	if got := v4.Bin(); !bytes.Equal(got, []byte{192, 168, 1, 1}) {
		t.Fatalf("Bin() = %v, want [192 168 1 1]", got)
	}

	integer, ok := v4.Uint32()
	if !ok || integer != 3232235777 {
		t.Fatalf("Uint32() = %d, %v, want 3232235777, true", integer, ok)
	}

	v6 := MustParseIP("2001:db8::1")
	if got := v6.Bin(); len(got) != ipBinarySize16 {
		t.Fatalf("Bin() for IPv6 length = %d, want %d", len(got), ipBinarySize16)
	}
	if _, ok := v6.Uint32(); ok {
		t.Fatal("Uint32() succeeded for an IPv6 address")
	}

	mapped := MustParseIP("::ffff:192.168.1.1")
	if got := mapped.Bin(); len(got) != ipBinarySize16 {
		t.Fatalf("Bin() for a mapped address length = %d, want %d", len(got), ipBinarySize16)
	}
	if _, ok := mapped.Uint32(); ok {
		t.Fatal("Uint32() succeeded for an IPv4-mapped address")
	}

	var null IP
	if value, err := null.Value(); err != nil || value != nil {
		t.Fatalf("Value() on NULL = %#v, %v", value, err)
	}
	if got := null.Bin(); got != nil {
		t.Fatalf("Bin() on NULL = %v, want nil", got)
	}
	if integer, ok := null.Uint32(); ok || integer != 0 {
		t.Fatalf("Uint32() on NULL = %d, %v, want 0, false", integer, ok)
	}
	if got := null.String(); got != "" {
		t.Fatalf("String() on NULL = %q, want an empty string", got)
	}
}

func TestIPJSON(t *testing.T) {
	v := MustParseIP("2001:db8::1")

	jsonBytes, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}
	if want := `"2001:db8::1"`; string(jsonBytes) != want {
		t.Fatalf("MarshalJSON() = %s, want %s", jsonBytes, want)
	}

	var parsed IP
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || !parsed.Valid || parsed.String() != v.String() {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"not-an-ip"`)); !errors.Is(err, ErrInvalidIP) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrInvalidIP", err)
	}

	if err := parsed.UnmarshalJSON([]byte("1")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("null")); err != nil || parsed.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", parsed, err)
	}

	nullJSON, err := (IP{}).MarshalJSON()
	if err != nil || string(nullJSON) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", nullJSON, err)
	}
}

func TestIPPrefixScan(t *testing.T) {
	testCases := []struct {
		name      string
		src       any
		want      string
		wantValid bool
		wantErr   error
	}{
		{name: "text", src: "192.168.1.0/24", want: "192.168.1.0/24", wantValid: true},
		{name: "text bytes", src: []byte("192.168.1.0/24"), want: "192.168.1.0/24", wantValid: true},
		{name: "host bits kept", src: "192.168.1.5/24", want: "192.168.1.5/24", wantValid: true},
		{name: "ipv6", src: "2001:db8::/32", want: "2001:db8::/32", wantValid: true},
		// The shortest prefix text is 4 characters long, so a 4 byte source must
		// be decoded as text rather than as a packed IPv4 address.
		{name: "four char text", src: []byte("::/0"), want: "::/0", wantValid: true},
		{name: "prefix", src: netip.MustParsePrefix("10.0.0.0/8"), want: "10.0.0.0/8", wantValid: true},
		{name: "nil", src: nil},
		{name: "zero prefix", src: netip.Prefix{}},
		{name: "missing bits", src: "192.168.1.0", wantErr: ErrInvalidIP},
		{name: "bad bits", src: "192.168.1.0/33", wantErr: ErrInvalidIP},
		{name: "empty", src: "", wantErr: ErrInvalidIP},
		{name: "unsupported", src: int64(1), wantErr: ErrCannotScan},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v IPPrefix
			err := v.Scan(tc.src)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Scan(%#v) error = %v, want %v", tc.src, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan(%#v) unexpected error: %v", tc.src, err)
			}
			if v.Valid != tc.wantValid {
				t.Fatalf("Scan(%#v).Valid = %v, want %v", tc.src, v.Valid, tc.wantValid)
			}
			if got := v.String(); got != tc.want {
				t.Fatalf("Scan(%#v).String() = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

func TestIPPrefixValueAndJSON(t *testing.T) {
	v := MustParseIPPrefix("192.168.1.5/24")

	if got := v.Prefix.Masked().String(); got != "192.168.1.0/24" {
		t.Fatalf("Masked() = %q, want 192.168.1.0/24", got)
	}

	if got := v.Bin(); !bytes.Equal(got, []byte{192, 168, 1, 5}) {
		t.Fatalf("Bin() = %v, want [192 168 1 5]", got)
	}

	value, err := v.Value()
	if err != nil || value != "192.168.1.5/24" {
		t.Fatalf("Value() = %#v, %v", value, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `"192.168.1.5/24"` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed IPPrefix
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.String() != v.String() {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"192.168.1.0"`)); !errors.Is(err, ErrInvalidIP) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrInvalidIP", err)
	}

	if err := parsed.UnmarshalJSON([]byte("1")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("null")); err != nil || parsed.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", parsed, err)
	}

	var null IPPrefix
	if value, err := null.Value(); err != nil || value != nil {
		t.Fatalf("Value() on NULL = %#v, %v", value, err)
	}
	if got := null.Bin(); got != nil {
		t.Fatalf("Bin() on NULL = %v, want nil", got)
	}
	if got := null.String(); got != "" {
		t.Fatalf("String() on NULL = %q, want an empty string", got)
	}
	if jsonBytes, err := null.MarshalJSON(); err != nil || string(jsonBytes) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", jsonBytes, err)
	}
}

func TestIPPortScan(t *testing.T) {
	testCases := []struct {
		name      string
		src       any
		want      string
		wantValid bool
		wantErr   error
	}{
		{name: "text", src: "192.168.1.1:3306", want: "192.168.1.1:3306", wantValid: true},
		{name: "text bytes", src: []byte("192.168.1.1:3306"), want: "192.168.1.1:3306", wantValid: true},
		{name: "ipv6", src: "[2001:db8::1]:3306", want: "[2001:db8::1]:3306", wantValid: true},
		{name: "addrPort", src: netip.MustParseAddrPort("10.0.0.1:80"), want: "10.0.0.1:80", wantValid: true},
		{name: "nil", src: nil},
		{name: "zero addrPort", src: netip.AddrPort{}},
		{name: "missing port", src: "192.168.1.1", wantErr: ErrInvalidIP},
		{name: "unbracketed ipv6", src: "2001:db8::1:3306", wantErr: ErrInvalidIP},
		{name: "bad port", src: "192.168.1.1:70000", wantErr: ErrInvalidIP},
		{name: "empty", src: "", wantErr: ErrInvalidIP},
		{name: "unsupported", src: true, wantErr: ErrCannotScan},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v IPPort
			err := v.Scan(tc.src)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Scan(%#v) error = %v, want %v", tc.src, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan(%#v) unexpected error: %v", tc.src, err)
			}
			if v.Valid != tc.wantValid {
				t.Fatalf("Scan(%#v).Valid = %v, want %v", tc.src, v.Valid, tc.wantValid)
			}
			if got := v.String(); got != tc.want {
				t.Fatalf("Scan(%#v).String() = %q, want %q", tc.src, got, tc.want)
			}
		})
	}
}

func TestIPPortValueAndJSON(t *testing.T) {
	v := MustParseIPPort("192.168.1.1:3306")

	if got := v.Bin(); !bytes.Equal(got, []byte{192, 168, 1, 1}) {
		t.Fatalf("Bin() = %v, want [192 168 1 1]", got)
	}

	value, err := v.Value()
	if err != nil || value != "192.168.1.1:3306" {
		t.Fatalf("Value() = %#v, %v", value, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `"192.168.1.1:3306"` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed IPPort
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.String() != v.String() {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"192.168.1.1"`)); !errors.Is(err, ErrInvalidIP) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrInvalidIP", err)
	}

	if err := parsed.UnmarshalJSON([]byte("1")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("null")); err != nil || parsed.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", parsed, err)
	}

	var null IPPort
	if value, err := null.Value(); err != nil || value != nil {
		t.Fatalf("Value() on NULL = %#v, %v", value, err)
	}
	if got := null.Bin(); got != nil {
		t.Fatalf("Bin() on NULL = %v, want nil", got)
	}
	if got := null.String(); got != "" {
		t.Fatalf("String() on NULL = %q, want an empty string", got)
	}
	if jsonBytes, err := null.MarshalJSON(); err != nil || string(jsonBytes) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", jsonBytes, err)
	}
}

func TestIPMustParsePanics(t *testing.T) {
	testCases := []struct {
		name string
		fn   func()
	}{
		{name: "IP", fn: func() { MustParseIP("not-an-ip") }},
		{name: "IPPrefix", fn: func() { MustParseIPPrefix("not-a-prefix") }},
		{name: "IPPort", fn: func() { MustParseIPPort("not-an-address-port") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("MustParse%s did not panic", tc.name)
				}
			}()
			tc.fn()
		})
	}
}
