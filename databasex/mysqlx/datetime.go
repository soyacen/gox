package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// dateTimeLayouts are the text layouts accepted for DATETIME and TIMESTAMP
// values, in the order they are tried.
var dateTimeLayouts = []string{
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	time.RFC3339Nano,
}

// scanDateTime scans a driver value into a time.Time using dateTimeLayouts.
//
// Parameters:
//   - src: a value produced by the database driver.
//   - dst: the destination time.
//
// Returns:
//   - bool: true when a non-NULL value was scanned.
//   - error: non-nil when src cannot be scanned.
func scanDateTime(src any, dst *time.Time) (bool, error) {
	v, isNull, err := coerceTime(src, dateTimeLayouts...)
	if err != nil {
		return false, err
	}

	if isNull {
		*dst = time.Time{}
		return false, nil
	}

	*dst = v
	return true, nil
}

// unmarshalDateTimeJSON parses a JSON string using dateTimeLayouts.
//
// Parameters:
//   - data: the raw JSON document.
//   - dst: the destination time.
//
// Returns:
//   - bool: true when a non-null value was parsed.
//   - error: non-nil when data cannot be parsed.
func unmarshalDateTimeJSON(data []byte, dst *time.Time) (bool, error) {
	if isJSONNull(data) {
		*dst = time.Time{}
		return false, nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return false, err
	}

	t, isNull, err := parseTimeString(s, dateTimeLayouts...)
	if err != nil {
		return false, err
	}
	if isNull {
		return false, nil
	}

	*dst = t
	return true, nil
}

// DateTime represents a MySQL DATETIME value. DATETIME carries no time zone
// information; text values are parsed as UTC.
type DateTime struct {
	// Time is the date and time value.
	Time time.Time
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface. The all-zero MySQL value
// "0000-00-00 00:00:00" is scanned as NULL.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a DateTime.
func (dst *DateTime) Scan(src any) error {
	valid, err := scanDateTime(src, &dst.Time)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the time.Time value, or nil when the value is NULL.
//   - error: always nil.
func (src DateTime) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Time, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface using RFC 3339
// with nanosecond precision.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src DateTime) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Time.Format(time.RFC3339Nano))
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. Both
// "2006-01-02 15:04:05" and RFC 3339 strings are accepted.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a recognized time string or null.
func (dst *DateTime) UnmarshalJSON(data []byte) error {
	valid, err := unmarshalDateTimeJSON(data, &dst.Time)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Timestamp represents a MySQL TIMESTAMP value. TIMESTAMP carries no time zone
// information in the protocol; text values are parsed as UTC.
type Timestamp struct {
	// Time is the date and time value.
	Time time.Time
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface. The all-zero MySQL value
// "0000-00-00 00:00:00" is scanned as NULL.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Timestamp.
func (dst *Timestamp) Scan(src any) error {
	valid, err := scanDateTime(src, &dst.Time)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the time.Time value, or nil when the value is NULL.
//   - error: always nil.
func (src Timestamp) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Time, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface using RFC 3339
// with nanosecond precision.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Timestamp) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Time.Format(time.RFC3339Nano))
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. Both
// "2006-01-02 15:04:05" and RFC 3339 strings are accepted.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a recognized time string or null.
func (dst *Timestamp) UnmarshalJSON(data []byte) error {
	valid, err := unmarshalDateTimeJSON(data, &dst.Time)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}
