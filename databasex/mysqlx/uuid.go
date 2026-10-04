package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"uuid"
)

// uuidBinarySize is the number of bytes in the packed BINARY(16) representation
// of a UUID.
const uuidBinarySize = 16

// UUID represents a MySQL UUID value.
//
// MySQL has no dedicated UUID column type, so identifiers are stored either in a
// CHAR(36) column using the canonical hyphenated form or in a BINARY(16) column
// using the packed form produced by the UUID_TO_BIN function. Scan accepts both
// representations, and Value writes the canonical form.
//
// The identifier is backed by the standard library [uuid.UUID] type introduced
// in Go 1.27.
type UUID struct {
	// UUID is the decoded identifier.
	UUID uuid.UUID
	// Valid reports whether the value is not NULL.
	Valid bool
}

// NewUUID returns a new version 4 UUID.
//
// Returns:
//   - UUID: a valid random UUID. It is equivalent to NewUUIDV4.
func NewUUID() UUID {
	return NewUUIDV4()
}

// NewUUIDV4 returns a new version 4 UUID.
//
// Version 4 UUIDs contain 122 bits of random data, which makes them suitable for
// identifiers that must not leak creation time.
//
// Returns:
//   - UUID: a valid random UUID.
func NewUUIDV4() UUID {
	return UUID{UUID: uuid.NewV4(), Valid: true}
}

// NewUUIDV7 returns a new version 7 UUID.
//
// Version 7 UUIDs embed a Unix millisecond timestamp in their most significant
// bits, so they sort in creation order. They are a good fit for MySQL primary
// keys because sequential inserts reduce B-tree page splits.
//
// Returns:
//   - UUID: a valid time ordered UUID.
func NewUUIDV7() UUID {
	return UUID{UUID: uuid.NewV7(), Valid: true}
}

// ParseUUID parses the string representation of a UUID.
//
// The accepted forms are the same as [uuid.Parse]: the canonical hyphenated
// form, the bare 32 character hexadecimal form, the brace wrapped form and the
// "urn:uuid:" prefixed form. Hexadecimal digits may use any case.
//
// Parameters:
//   - s: the string representation of a UUID.
//
// Returns:
//   - UUID: the parsed UUID with Valid set to true.
//   - error: non-nil when s is not a valid UUID, wrapping [ErrInvalidUUID].
func ParseUUID(s string) (UUID, error) {
	u, err := parseUUID(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{UUID: u, Valid: true}, nil
}

// MustParseUUID is like ParseUUID but panics when s is not a valid UUID.
//
// Parameters:
//   - s: the string representation of a UUID.
//
// Returns:
//   - UUID: the parsed UUID with Valid set to true.
func MustParseUUID(s string) UUID {
	u, err := ParseUUID(s)
	if err != nil {
		panic(err)
	}
	return u
}

// parseUUID parses s and wraps the standard library parse error.
//
// Parameters:
//   - s: the string representation of a UUID.
//
// Returns:
//   - uuid.UUID: the parsed identifier.
//   - error: non-nil when s is not a valid UUID, wrapping [ErrInvalidUUID].
func parseUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: %q", ErrInvalidUUID, s)
	}
	return u, nil
}

// Scan implements the database/sql.Scanner interface.
//
// A string or []byte input is parsed as text unless it is exactly 16 bytes long,
// in which case it is treated as the packed BINARY(16) representation.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a UUID.
func (dst *UUID) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = UUID{}
		return nil
	case uuid.UUID:
		dst.UUID, dst.Valid = v, true
		return nil
	case [16]byte:
		dst.UUID, dst.Valid = uuid.UUID(v), true
		return nil
	case string:
		return dst.scanUUIDText(v)
	case []byte:
		if len(v) == uuidBinarySize {
			var raw [16]byte
			copy(raw[:], v)
			dst.UUID, dst.Valid = uuid.UUID(raw), true
			return nil
		}
		return dst.scanUUIDText(string(v))
	default:
		return cannotScan(src, dst)
	}
}

// scanUUIDText parses the textual representation s into dst.
//
// Parameters:
//   - s: the textual representation of a UUID.
//
// Returns:
//   - error: non-nil when s is not a valid UUID, wrapping [ErrInvalidUUID].
func (dst *UUID) scanUUIDText(s string) error {
	u, err := parseUUID(s)
	if err != nil {
		return err
	}
	dst.UUID, dst.Valid = u, true
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// The canonical hyphenated form is written, which is the representation expected
// by CHAR(36) columns and by MySQL functions such as UUID_TO_BIN and
// BIN_TO_UUID.
//
// Returns:
//   - driver.Value: the canonical string, or nil when the value is NULL.
//   - error: always nil.
func (src UUID) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.UUID.String(), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the canonical string as a JSON string, or null when the value is
//     NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src UUID) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.UUID.String())
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON string accepted by ParseUUID or
//     null.
func (dst *UUID) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = UUID{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return dst.scanUUIDText(s)
}

// String returns the canonical lowercase hyphenated representation.
//
// Returns:
//   - string: the canonical representation, or an empty string when the value is
//     NULL. A valid Nil UUID yields
//     "00000000-0000-0000-0000-000000000000", which is not the same as NULL.
func (src UUID) String() string {
	if !src.Valid {
		return ""
	}
	return src.UUID.String()
}

// Bytes returns the packed representation of the UUID.
//
// Returns:
//   - [16]byte: the big-endian bytes of the identifier, or the zero array when
//     the value is NULL.
func (src UUID) Bytes() [16]byte {
	return [16]byte(src.UUID)
}

// Bin returns the packed representation of the UUID as a slice.
//
// The result can be compared against BINARY(16) columns or passed to MySQL
// functions such as BIN_TO_UUID.
//
// Returns:
//   - []byte: a newly allocated 16 byte slice, or nil when the value is NULL.
func (src UUID) Bin() []byte {
	if !src.Valid {
		return nil
	}
	b := make([]byte, uuidBinarySize)
	copy(b, src.UUID[:])
	return b
}
