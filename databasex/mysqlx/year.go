package mysqlx

import (
	"database/sql/driver"
	"fmt"
)

// yearMin is the smallest non-zero MySQL YEAR value.
const yearMin = 1901

// yearMax is the largest MySQL YEAR value.
const yearMax = 2155

// Year represents a MySQL YEAR value. The value 0 represents the MySQL zero
// year, and valid non-zero years are between 1901 and 2155.
type Year struct {
	// Int16 is the year value.
	Int16 int16
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Year or is outside the
//     MySQL YEAR range.
func (dst *Year) Scan(src any) error {
	v, isNull, err := coerceInt64(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Year{}
		return nil
	}

	if v != 0 && (v < yearMin || v > yearMax) {
		return fmt.Errorf("%w: %d is not a valid %s", ErrInvalidYear, v, TypeYear)
	}

	*dst = Year{Int16: int16(v), Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the year as an int64, or nil when the value is NULL.
//   - error: non-nil when the value is outside the MySQL YEAR range.
func (src Year) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	if !validYear(src.Int16) {
		return nil, fmt.Errorf("%w: %d is not a valid %s", ErrInvalidYear, src.Int16, TypeYear)
	}
	return int64(src.Int16), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value is outside the MySQL YEAR range.
func (src Year) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	if !validYear(src.Int16) {
		return nil, fmt.Errorf("%w: %d is not a valid %s", ErrInvalidYear, src.Int16, TypeYear)
	}
	return marshalSignedJSON(src.Int16, true)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is outside the MySQL YEAR
//     range, or is not null.
func (dst *Year) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int16](data)
	if err != nil {
		return err
	}
	if valid && !validYear(v) {
		return fmt.Errorf("%w: %d is not a valid %s", ErrInvalidYear, v, TypeYear)
	}
	*dst = Year{Int16: v, Valid: valid}
	return nil
}

// validYear reports whether v is 0 or within the supported MySQL YEAR range.
//
// Parameters:
//   - v: the year value.
//
// Returns:
//   - bool: true when v is a valid YEAR value.
func validYear(v int16) bool {
	return v == 0 || (v >= yearMin && v <= yearMax)
}
