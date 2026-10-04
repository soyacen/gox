package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON represents a MySQL JSON value. The raw JSON document is preserved
// byte-for-byte.
//
// The SQL NULL value and the JSON null literal are indistinguishable in the
// JSON encoding: both are marshaled as null, and unmarshaling null yields the
// SQL NULL value.
type JSON struct {
	// Bytes is the raw JSON document.
	Bytes []byte
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface. The scanned document is
// validated as JSON.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a JSON value or is not a
//     valid JSON document.
func (dst *JSON) Scan(src any) error {
	v, isNull, err := coerceBytes(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = JSON{}
		return nil
	}

	if !json.Valid(v) {
		return fmt.Errorf("%w: invalid JSON document", ErrInvalidJSON)
	}

	*dst = JSON{Bytes: v, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the raw JSON document, or nil when the value is NULL.
//   - error: non-nil when the stored bytes are not a valid JSON document.
func (src JSON) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	if !json.Valid(src.Bytes) {
		return nil, fmt.Errorf("%w: invalid JSON document", ErrInvalidJSON)
	}
	return src.Bytes, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The raw JSON
// document is emitted unchanged.
//
// Returns:
//   - []byte: the raw JSON document, or null when the value is NULL.
//   - error: non-nil when the stored bytes are not a valid JSON document.
func (src JSON) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	if !json.Valid(src.Bytes) {
		return nil, fmt.Errorf("%w: invalid JSON document", ErrInvalidJSON)
	}
	return src.Bytes, nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The JSON
// null literal is decoded as the SQL NULL value.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a valid JSON document.
func (dst *JSON) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = JSON{}
		return nil
	}

	if !json.Valid(data) {
		return fmt.Errorf("%w: invalid JSON document", ErrInvalidJSON)
	}

	*dst = JSON{Bytes: cloneBytes(data), Valid: true}
	return nil
}
