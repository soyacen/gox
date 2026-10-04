package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
)

// Enum represents a MySQL ENUM value.
//
// The Go type has no access to the member list declared on the column, so
// scanned and encoded values are not validated against it. MySQL itself rejects
// values outside the declared member list.
type Enum struct {
	// String is the enum member value.
	String string
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an Enum.
func (dst *Enum) Scan(src any) error {
	valid, err := scanText(src, &dst.String)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the enum member, or nil when the value is NULL.
//   - error: always nil.
func (src Enum) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.String, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Enum) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.String)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON string or null.
func (dst *Enum) UnmarshalJSON(data []byte) error {
	valid, err := unmarshalText(data, &dst.String)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}
