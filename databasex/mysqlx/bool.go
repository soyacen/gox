package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
)

// Bool represents a MySQL BOOL value. BOOL is stored as TINYINT(1), so any
// non-zero numeric value is scanned as true.
type Bool struct {
	// Bool is the boolean value.
	Bool bool
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Bool.
func (dst *Bool) Scan(src any) error {
	v, isNull, err := coerceBool(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Bool{}
		return nil
	}

	*dst = Bool{Bool: v, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the boolean value, or nil when the value is NULL.
//   - error: always nil.
func (src Bool) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Bool, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON boolean, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Bool) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Bool)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON boolean or null.
func (dst *Bool) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Bool{}
		return nil
	}

	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*dst = Bool{Bool: v, Valid: true}
	return nil
}
