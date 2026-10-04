package mysqlx

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Bit represents a MySQL BIT value. MySQL limits BIT to 64 bits, so the value
// is stored as a uint64.
//
// Because a []byte may be either the raw big-endian bit string produced by a
// MySQL driver or a textual decimal representation, the two forms are
// distinguished by their Go type: []byte is always interpreted as the raw
// big-endian bit string, while string is always interpreted as a decimal
// literal.
type Bit struct {
	// Uint64 is the bit value.
	Uint64 uint64
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Bit or holds more than
//     64 bits.
func (dst *Bit) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = Bit{}
		return nil
	case []byte:
		if len(v) > 8 {
			return fmt.Errorf("%w: %d bytes does not fit in %s", ErrOutOfRange, len(v), TypeBit)
		}
		var buf [8]byte
		copy(buf[len(buf)-len(v):], v)
		*dst = Bit{Uint64: binary.BigEndian.Uint64(buf[:]), Valid: true}
		return nil
	case string:
		return dst.scanString(v)
	}

	valid, err := scanUnsigned[uint64](src, &dst.Uint64)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// scanString parses a textual decimal bit value.
//
// Parameters:
//   - s: the decimal literal.
//
// Returns:
//   - error: non-nil when s is not a valid unsigned decimal literal.
func (dst *Bit) scanString(s string) error {
	u, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: cannot parse %q as %s", ErrCannotScan, s, TypeBit)
	}
	*dst = Bit{Uint64: u, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface. The value is
// encoded as a big-endian byte slice with leading zero bytes removed.
//
// Returns:
//   - driver.Value: the big-endian bit string, or nil when the value is NULL.
//   - error: always nil.
func (src Bit) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return uint64Bytes(src.Uint64), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: always nil.
func (src Bit) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return strconv.AppendUint(nil, src.Uint64, 10), nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. Both JSON
// numbers and JSON strings containing a decimal literal are accepted.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not an unsigned decimal number or null.
func (dst *Bit) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Bit{}
		return nil
	}

	raw := strings.TrimSpace(string(data))
	if len(raw) > 1 && raw[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		return dst.scanString(s)
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*dst = Bit{Uint64: v, Valid: true}
	return nil
}

// uint64Bytes returns the big-endian representation of v with leading zero
// bytes removed, keeping at least one byte.
//
// Parameters:
//   - v: the value to encode.
//
// Returns:
//   - []byte: the big-endian bit string.
func uint64Bytes(v uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)

	start := 0
	for start < len(buf)-1 && buf[start] == 0 {
		start++
	}
	return cloneBytes(buf[start:])
}
