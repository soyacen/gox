package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
)

// scanText scans a driver value into a string field.
//
// Parameters:
//   - src: a value produced by the database driver.
//   - dst: the destination string.
//
// Returns:
//   - bool: true when a non-NULL value was scanned.
//   - error: non-nil when src cannot be scanned.
func scanText(src any, dst *string) (bool, error) {
	v, isNull, err := coerceString(src)
	if err != nil {
		return false, err
	}

	if isNull {
		*dst = ""
		return false, nil
	}

	*dst = v
	return true, nil
}

// unmarshalText parses a JSON string into a string field.
//
// Parameters:
//   - data: the raw JSON document.
//   - dst: the destination string.
//
// Returns:
//   - bool: true when a non-null JSON string was parsed.
//   - error: non-nil when data is not a JSON string.
func unmarshalText(data []byte, dst *string) (bool, error) {
	if isJSONNull(data) {
		*dst = ""
		return false, nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return false, err
	}

	*dst = s
	return true, nil
}

// Text represents a MySQL text value. CHAR, VARCHAR, TINYTEXT, MEDIUMTEXT and
// LONGTEXT share the exact same Go representation and are exposed as aliases of
// Text.
type Text struct {
	// String is the text value.
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
//   - error: non-nil when src cannot be scanned into a Text.
func (dst *Text) Scan(src any) error {
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
//   - driver.Value: the text value, or nil when the value is NULL.
//   - error: always nil.
func (src Text) Value() (driver.Value, error) {
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
func (src Text) MarshalJSON() ([]byte, error) {
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
func (dst *Text) UnmarshalJSON(data []byte) error {
	valid, err := unmarshalText(data, &dst.String)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Char is a MySQL CHAR value. It is an alias for Text.
type Char = Text

// VarChar is a MySQL VARCHAR value. It is an alias for Text.
type VarChar = Text

// TinyText is a MySQL TINYTEXT value. It is an alias for Text.
type TinyText = Text

// MediumText is a MySQL MEDIUMTEXT value. It is an alias for Text.
type MediumText = Text

// LongText is a MySQL LONGTEXT value. It is an alias for Text.
type LongText = Text
