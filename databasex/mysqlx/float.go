package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
)

// Float represents a MySQL FLOAT value, a single precision floating point
// number.
type Float struct {
	// Float32 is the floating point value.
	Float32 float32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Float or exceeds the
//     single precision range.
func (dst *Float) Scan(src any) error {
	v, isNull, err := coerceFloat64(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Float{}
		return nil
	}

	if !math.IsNaN(v) && !math.IsInf(v, 0) && (v > math.MaxFloat32 || v < -math.MaxFloat32) {
		return fmt.Errorf("%w: %v does not fit in %s", ErrOutOfRange, v, TypeFloat)
	}

	*dst = Float{Float32: float32(v), Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the float64 representation of the value, or nil when the
//     value is NULL.
//   - error: always nil.
func (src Float) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return float64(src.Float32), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Float) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Float32)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number or null.
func (dst *Float) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Float{}
		return nil
	}

	var v float32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*dst = Float{Float32: v, Valid: true}
	return nil
}

// Double represents a MySQL DOUBLE value, a double precision floating point
// number. DOUBLE, REAL and DOUBLE PRECISION are synonyms.
type Double struct {
	// Float64 is the floating point value.
	Float64 float64
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Double.
func (dst *Double) Scan(src any) error {
	v, isNull, err := coerceFloat64(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Double{}
		return nil
	}

	*dst = Double{Float64: v, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the float64 value, or nil when the value is NULL.
//   - error: always nil.
func (src Double) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Float64, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Double) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Float64)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number or null.
func (dst *Double) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Double{}
		return nil
	}

	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*dst = Double{Float64: v, Valid: true}
	return nil
}
