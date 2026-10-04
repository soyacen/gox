package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// dateLayout is the MySQL DATE text layout.
const dateLayout = "2006-01-02"

// Date represents a MySQL DATE value.
type Date struct {
	// Time is the date at midnight UTC when scanned from text, or the original
	// time when the driver returned a time.Time.
	Time time.Time
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface. The all-zero MySQL date
// "0000-00-00" is scanned as NULL.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Date.
func (dst *Date) Scan(src any) error {
	v, isNull, err := coerceTime(src, dateLayout)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Date{}
		return nil
	}

	*dst = Date{Time: v, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the time.Time value, or nil when the value is NULL.
//   - error: always nil.
func (src Date) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Time, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface using the
// "2006-01-02" layout.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Date) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Time.Format(dateLayout))
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The
// accepted layout is "2006-01-02".
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a date string or null.
func (dst *Date) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Date{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.ParseInLocation(dateLayout, s, time.UTC)
	if err != nil {
		return err
	}

	*dst = Date{Time: t, Valid: true}
	return nil
}
