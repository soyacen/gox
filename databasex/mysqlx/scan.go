package mysqlx

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

// signed is the constraint for the Go signed integer types used by the MySQL
// integer types.
type signed interface {
	~int8 | ~int16 | ~int32 | ~int64
}

// unsigned is the constraint for the Go unsigned integer types used by the
// MySQL unsigned integer types.
type unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

// isJSONNull reports whether data is the JSON null literal, ignoring
// surrounding whitespace.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - bool: true when data represents JSON null.
func isJSONNull(data []byte) bool {
	return bytes.Equal(bytes.TrimSpace(data), []byte("null"))
}

// cloneBytes returns a copy of b so that values retained from a database driver
// remain valid after the underlying row buffer is reused.
//
// Parameters:
//   - b: the byte slice to copy.
//
// Returns:
//   - []byte: a copy of b, or nil when b is nil.
func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

// cannotScan builds the error returned when a driver value cannot be scanned
// into a destination.
//
// Parameters:
//   - src: the value produced by the driver.
//   - dst: a typed nil pointer describing the destination type.
//
// Returns:
//   - error: an error wrapping ErrCannotScan.
func cannotScan(src any, dst any) error {
	return fmt.Errorf("%w: %T into %T", ErrCannotScan, src, dst)
}

// coerceBytes converts a driver value into a byte slice. The returned bytes are
// always a copy owned by the caller.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - []byte: the decoded bytes.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src cannot be converted.
func coerceBytes(src any) ([]byte, bool, error) {
	switch v := src.(type) {
	case nil:
		return nil, true, nil
	case []byte:
		return cloneBytes(v), false, nil
	case string:
		return []byte(v), false, nil
	default:
		return nil, false, cannotScan(src, (*[]byte)(nil))
	}
}

// coerceString converts a driver value into a string.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - string: the decoded string.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src cannot be converted.
func coerceString(src any) (string, bool, error) {
	switch v := src.(type) {
	case nil:
		return "", true, nil
	case string:
		return v, false, nil
	case []byte:
		return string(v), false, nil
	default:
		return "", false, cannotScan(src, (*string)(nil))
	}
}

// coerceInt64 converts a driver value into an int64.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - int64: the decoded integer.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src cannot be converted or does not fit in an int64.
func coerceInt64(src any) (int64, bool, error) {
	switch v := src.(type) {
	case nil:
		return 0, true, nil
	case int:
		return int64(v), false, nil
	case int8:
		return int64(v), false, nil
	case int16:
		return int64(v), false, nil
	case int32:
		return int64(v), false, nil
	case int64:
		return v, false, nil
	case uint:
		return uintToInt64(uint64(v))
	case uint8:
		return int64(v), false, nil
	case uint16:
		return int64(v), false, nil
	case uint32:
		return int64(v), false, nil
	case uint64:
		return uintToInt64(v)
	case bool:
		if v {
			return 1, false, nil
		}
		return 0, false, nil
	case string:
		return parseInt64(v, "")
	case []byte:
		return parseInt64(string(v), "[]byte")
	default:
		return 0, false, cannotScan(src, (*int64)(nil))
	}
}

// uintToInt64 converts a uint64 that must fit in an int64.
//
// Parameters:
//   - v: the value to convert.
//
// Returns:
//   - int64: the converted value.
//   - bool: always false.
//   - error: non-nil when v exceeds math.MaxInt64.
func uintToInt64(v uint64) (int64, bool, error) {
	if v > math.MaxInt64 {
		return 0, false, fmt.Errorf("%w: %d does not fit in int64", ErrOutOfRange, v)
	}
	return int64(v), false, nil
}

// parseInt64 parses a decimal integer literal.
//
// Parameters:
//   - s: the literal to parse.
//   - kind: a label used in error messages to describe the source.
//
// Returns:
//   - int64: the parsed value.
//   - bool: always false.
//   - error: non-nil when s is not a valid int64 literal.
func parseInt64(s string, kind string) (int64, bool, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if kind == "" {
			kind = "string"
		}
		return 0, false, fmt.Errorf("%w: cannot parse %s %q as int64", ErrCannotScan, kind, s)
	}
	return v, false, nil
}

// coerceUint64 converts a driver value into a uint64.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - uint64: the decoded integer.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src is negative or cannot be converted.
func coerceUint64(src any) (uint64, bool, error) {
	switch v := src.(type) {
	case nil:
		return 0, true, nil
	case uint:
		return uint64(v), false, nil
	case uint8:
		return uint64(v), false, nil
	case uint16:
		return uint64(v), false, nil
	case uint32:
		return uint64(v), false, nil
	case uint64:
		return v, false, nil
	case int:
		return intToUint64(int64(v))
	case int8:
		return intToUint64(int64(v))
	case int16:
		return intToUint64(int64(v))
	case int32:
		return intToUint64(int64(v))
	case int64:
		return intToUint64(v)
	case bool:
		if v {
			return 1, false, nil
		}
		return 0, false, nil
	case string:
		return parseUint64(v, "")
	case []byte:
		return parseUint64(string(v), "[]byte")
	default:
		return 0, false, cannotScan(src, (*uint64)(nil))
	}
}

// intToUint64 converts an int64 that must be non-negative.
//
// Parameters:
//   - v: the value to convert.
//
// Returns:
//   - uint64: the converted value.
//   - bool: always false.
//   - error: non-nil when v is negative.
func intToUint64(v int64) (uint64, bool, error) {
	if v < 0 {
		return 0, false, fmt.Errorf("%w: %d does not fit in uint64", ErrOutOfRange, v)
	}
	return uint64(v), false, nil
}

// parseUint64 parses a decimal unsigned integer literal.
//
// Parameters:
//   - s: the literal to parse.
//   - kind: a label used in error messages to describe the source.
//
// Returns:
//   - uint64: the parsed value.
//   - bool: always false.
//   - error: non-nil when s is not a valid uint64 literal.
func parseUint64(s string, kind string) (uint64, bool, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		if kind == "" {
			kind = "string"
		}
		return 0, false, fmt.Errorf("%w: cannot parse %s %q as uint64", ErrCannotScan, kind, s)
	}
	return v, false, nil
}

// coerceFloat64 converts a driver value into a float64.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - float64: the decoded number.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src cannot be converted.
func coerceFloat64(src any) (float64, bool, error) {
	switch v := src.(type) {
	case nil:
		return 0, true, nil
	case float64:
		return v, false, nil
	case float32:
		return float64(v), false, nil
	case int:
		return float64(v), false, nil
	case int8:
		return float64(v), false, nil
	case int16:
		return float64(v), false, nil
	case int32:
		return float64(v), false, nil
	case int64:
		return float64(v), false, nil
	case uint:
		return float64(v), false, nil
	case uint8:
		return float64(v), false, nil
	case uint16:
		return float64(v), false, nil
	case uint32:
		return float64(v), false, nil
	case uint64:
		return float64(v), false, nil
	case string:
		return parseFloat64(v, "")
	case []byte:
		return parseFloat64(string(v), "[]byte")
	default:
		return 0, false, cannotScan(src, (*float64)(nil))
	}
}

// parseFloat64 parses a floating point literal.
//
// Parameters:
//   - s: the literal to parse.
//   - kind: a label used in error messages to describe the source.
//
// Returns:
//   - float64: the parsed value.
//   - bool: always false.
//   - error: non-nil when s is not a valid float64 literal.
func parseFloat64(s string, kind string) (float64, bool, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if kind == "" {
			kind = "string"
		}
		return 0, false, fmt.Errorf("%w: cannot parse %s %q as float64", ErrCannotScan, kind, s)
	}
	return v, false, nil
}

// coerceBool converts a driver value into a bool. Any non-zero number is
// treated as true.
//
// Parameters:
//   - src: the value produced by the driver.
//
// Returns:
//   - bool: the decoded boolean.
//   - bool: true when src is SQL NULL.
//   - error: non-nil when src cannot be converted.
func coerceBool(src any) (bool, bool, error) {
	switch v := src.(type) {
	case nil:
		return false, true, nil
	case bool:
		return v, false, nil
	case int:
		return v != 0, false, nil
	case int8:
		return v != 0, false, nil
	case int16:
		return v != 0, false, nil
	case int32:
		return v != 0, false, nil
	case int64:
		return v != 0, false, nil
	case uint:
		return v != 0, false, nil
	case uint8:
		return v != 0, false, nil
	case uint16:
		return v != 0, false, nil
	case uint32:
		return v != 0, false, nil
	case uint64:
		return v != 0, false, nil
	case float32:
		return v != 0, false, nil
	case float64:
		return v != 0, false, nil
	case string:
		return parseBool(v, "")
	case []byte:
		return parseBool(string(v), "[]byte")
	default:
		return false, false, cannotScan(src, (*bool)(nil))
	}
}

// parseBool parses a boolean literal using the strconv.ParseBool rules.
//
// Parameters:
//   - s: the literal to parse.
//   - kind: a label used in error messages to describe the source.
//
// Returns:
//   - bool: the parsed value.
//   - bool: always false.
//   - error: non-nil when s is not a valid boolean literal.
func parseBool(s string, kind string) (bool, bool, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		if kind == "" {
			kind = "string"
		}
		return false, false, fmt.Errorf("%w: cannot parse %s %q as bool", ErrCannotScan, kind, s)
	}
	return v, false, nil
}

// coerceTime converts a driver value into a time.Time by trying each layout in
// order. The all-zero MySQL date literals are reported as NULL.
//
// Parameters:
//   - src: the value produced by the driver.
//   - layouts: the accepted time layouts, tried in order.
//
// Returns:
//   - time.Time: the decoded time.
//   - bool: true when src is SQL NULL or an all-zero date.
//   - error: non-nil when src cannot be converted.
func coerceTime(src any, layouts ...string) (time.Time, bool, error) {
	switch v := src.(type) {
	case nil:
		return time.Time{}, true, nil
	case time.Time:
		return v, false, nil
	case string:
		return parseTimeString(v, layouts...)
	case []byte:
		return parseTimeString(string(v), layouts...)
	default:
		return time.Time{}, false, cannotScan(src, (*time.Time)(nil))
	}
}

// parseTimeString parses a time literal using the first matching layout.
//
// Parameters:
//   - s: the literal to parse.
//   - layouts: the accepted time layouts, tried in order.
//
// Returns:
//   - time.Time: the parsed time.
//   - bool: true when s is an all-zero date.
//   - error: non-nil when no layout matches.
func parseTimeString(s string, layouts ...string) (time.Time, bool, error) {
	if isZeroDate(s) {
		return time.Time{}, true, nil
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			return t, false, nil
		}
	}

	return time.Time{}, false, fmt.Errorf("%w: cannot parse %q as time", ErrCannotScan, s)
}

// isZeroDate reports whether s is an all-zero MySQL date literal, which cannot
// be represented by time.Time.
//
// Parameters:
//   - s: the literal to inspect.
//
// Returns:
//   - bool: true when s starts with the "0000-00-00" date.
func isZeroDate(s string) bool {
	return len(s) >= len(zeroDatePrefix) && s[:len(zeroDatePrefix)] == zeroDatePrefix
}

// zeroDatePrefix is the MySQL all-zero date prefix.
const zeroDatePrefix = "0000-00-00"

// signedRange returns the inclusive value range of the signed integer type T.
//
// Type Parameters:
//   - T: a signed integer type.
//
// Returns:
//   - int64: the minimum value.
//   - int64: the maximum value.
func signedRange[T signed]() (int64, int64) {
	var zero T
	switch any(zero).(type) {
	case int8:
		return math.MinInt8, math.MaxInt8
	case int16:
		return math.MinInt16, math.MaxInt16
	case int32:
		return math.MinInt32, math.MaxInt32
	default:
		return math.MinInt64, math.MaxInt64
	}
}

// unsignedRange returns the inclusive value range of the unsigned integer type T.
//
// Type Parameters:
//   - T: an unsigned integer type.
//
// Returns:
//   - uint64: the minimum value (always 0).
//   - uint64: the maximum value.
func unsignedRange[T unsigned]() (uint64, uint64) {
	var zero T
	switch any(zero).(type) {
	case uint8:
		return 0, math.MaxUint8
	case uint16:
		return 0, math.MaxUint16
	case uint32:
		return 0, math.MaxUint32
	default:
		return 0, math.MaxUint64
	}
}

// scanSigned scans a driver value into a signed integer destination, enforcing
// the destination range.
//
// Type Parameters:
//   - T: the destination signed integer type.
//
// Parameters:
//   - src: the value produced by the driver.
//   - dst: the destination value.
//
// Returns:
//   - bool: true when a non-NULL value was scanned.
//   - error: non-nil when src cannot be converted or does not fit in T.
func scanSigned[T signed](src any, dst *T) (bool, error) {
	v, isNull, err := coerceInt64(src)
	if err != nil {
		return false, err
	}

	var zero T
	if isNull {
		*dst = zero
		return false, nil
	}

	lo, hi := signedRange[T]()
	if v < lo || v > hi {
		return false, fmt.Errorf("%w: %d does not fit in %T", ErrOutOfRange, v, zero)
	}

	*dst = T(v)
	return true, nil
}

// scanUnsigned scans a driver value into an unsigned integer destination,
// enforcing the destination range.
//
// Type Parameters:
//   - T: the destination unsigned integer type.
//
// Parameters:
//   - src: the value produced by the driver.
//   - dst: the destination value.
//
// Returns:
//   - bool: true when a non-NULL value was scanned.
//   - error: non-nil when src cannot be converted or does not fit in T.
func scanUnsigned[T unsigned](src any, dst *T) (bool, error) {
	v, isNull, err := coerceUint64(src)
	if err != nil {
		return false, err
	}

	var zero T
	if isNull {
		*dst = zero
		return false, nil
	}

	_, hi := unsignedRange[T]()
	if v > hi {
		return false, fmt.Errorf("%w: %d does not fit in %T", ErrOutOfRange, v, zero)
	}

	*dst = T(v)
	return true, nil
}

// signedValue encodes a signed integer as a driver value.
//
// Type Parameters:
//   - T: the source signed integer type.
//
// Parameters:
//   - v: the value to encode.
//   - valid: whether the value is not NULL.
//
// Returns:
//   - driver.Value: the encoded value, or nil for NULL.
//   - error: always nil.
func signedValue[T signed](v T, valid bool) (driver.Value, error) {
	if !valid {
		return nil, nil
	}
	return int64(v), nil
}

// unsignedValue encodes an unsigned integer as a driver value.
//
// Type Parameters:
//   - T: the source unsigned integer type.
//
// Parameters:
//   - v: the value to encode.
//   - valid: whether the value is not NULL.
//
// Returns:
//   - driver.Value: the encoded value, or nil for NULL.
//   - error: always nil.
func unsignedValue[T unsigned](v T, valid bool) (driver.Value, error) {
	if !valid {
		return nil, nil
	}
	return uint64(v), nil
}

// marshalSignedJSON marshals a signed integer as a JSON number.
//
// Type Parameters:
//   - T: the source signed integer type.
//
// Parameters:
//   - v: the value to marshal.
//   - valid: whether the value is not NULL.
//
// Returns:
//   - []byte: the JSON encoding.
//   - error: always nil.
func marshalSignedJSON[T signed](v T, valid bool) ([]byte, error) {
	if !valid {
		return []byte("null"), nil
	}
	return strconv.AppendInt(nil, int64(v), 10), nil
}

// marshalUnsignedJSON marshals an unsigned integer as a JSON number.
//
// Type Parameters:
//   - T: the source unsigned integer type.
//
// Parameters:
//   - v: the value to marshal.
//   - valid: whether the value is not NULL.
//
// Returns:
//   - []byte: the JSON encoding.
//   - error: always nil.
func marshalUnsignedJSON[T unsigned](v T, valid bool) ([]byte, error) {
	if !valid {
		return []byte("null"), nil
	}
	return strconv.AppendUint(nil, uint64(v), 10), nil
}

// unmarshalSignedJSON parses a JSON number into a signed integer.
//
// Type Parameters:
//   - T: the destination signed integer type.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - T: the decoded value.
//   - bool: true when data is a non-null JSON number.
//   - error: non-nil when data is not a number or does not fit in T.
func unmarshalSignedJSON[T signed](data []byte) (T, bool, error) {
	var zero T
	if isJSONNull(data) {
		return zero, false, nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return zero, false, err
	}

	lo, hi := signedRange[T]()
	if v < lo || v > hi {
		return zero, false, fmt.Errorf("%w: %d does not fit in %T", ErrOutOfRange, v, zero)
	}

	return T(v), true, nil
}

// unmarshalUnsignedJSON parses a JSON number into an unsigned integer.
//
// Type Parameters:
//   - T: the destination unsigned integer type.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - T: the decoded value.
//   - bool: true when data is a non-null JSON number.
//   - error: non-nil when data is not a number or does not fit in T.
func unmarshalUnsignedJSON[T unsigned](data []byte) (T, bool, error) {
	var zero T
	if isJSONNull(data) {
		return zero, false, nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return zero, false, err
	}

	_, hi := unsignedRange[T]()
	if v > hi {
		return zero, false, fmt.Errorf("%w: %d does not fit in %T", ErrOutOfRange, v, zero)
	}

	return T(v), true, nil
}
