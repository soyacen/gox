package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
)

// scanBlob scans a driver value into a byte slice field. The stored bytes are
// always a copy, so they remain valid after the driver reuses its row buffer.
//
// Parameters:
//   - src: a value produced by the database driver.
//   - dst: the destination byte slice.
//
// Returns:
//   - bool: true when a non-NULL value was scanned.
//   - error: non-nil when src cannot be scanned.
func scanBlob(src any, dst *[]byte) (bool, error) {
	v, isNull, err := coerceBytes(src)
	if err != nil {
		return false, err
	}

	if isNull {
		*dst = nil
		return false, nil
	}

	*dst = v
	return true, nil
}

// unmarshalBlob parses a base64 encoded JSON string into a byte slice field.
//
// Parameters:
//   - data: the raw JSON document.
//   - dst: the destination byte slice.
//
// Returns:
//   - bool: true when a non-null JSON string was parsed.
//   - error: non-nil when data is not a base64 encoded JSON string.
func unmarshalBlob(data []byte, dst *[]byte) (bool, error) {
	if isJSONNull(data) {
		*dst = nil
		return false, nil
	}

	var b []byte
	if err := json.Unmarshal(data, &b); err != nil {
		return false, err
	}

	*dst = b
	return true, nil
}

// Blob represents a MySQL binary value. BINARY, VARBINARY, TINYBLOB,
// MEDIUMBLOB and LONGBLOB share the exact same Go representation and are
// exposed as aliases of Blob.
type Blob struct {
	// Bytes is the binary value.
	Bytes []byte
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Blob.
func (dst *Blob) Scan(src any) error {
	valid, err := scanBlob(src, &dst.Bytes)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the binary value, or nil when the value is NULL.
//   - error: always nil.
func (src Blob) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Bytes, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The bytes are
// encoded as a base64 JSON string.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Blob) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Bytes)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The value
// is expected to be a base64 JSON string.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a base64 JSON string or null.
func (dst *Blob) UnmarshalJSON(data []byte) error {
	valid, err := unmarshalBlob(data, &dst.Bytes)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Binary is a MySQL BINARY value. It is an alias for Blob.
type Binary = Blob

// VarBinary is a MySQL VARBINARY value. It is an alias for Blob.
type VarBinary = Blob

// TinyBlob is a MySQL TINYBLOB value. It is an alias for Blob.
type TinyBlob = Blob

// MediumBlob is a MySQL MEDIUMBLOB value. It is an alias for Blob.
type MediumBlob = Blob

// LongBlob is a MySQL LONGBLOB value. It is an alias for Blob.
type LongBlob = Blob
