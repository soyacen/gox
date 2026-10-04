package mysqlx

import (
	"database/sql/driver"
	"fmt"
)

// mediumIntMin is the smallest value of a signed 24-bit MEDIUMINT.
const mediumIntMin = -8388608

// mediumIntMax is the largest value of a signed 24-bit MEDIUMINT.
const mediumIntMax = 8388607

// mediumIntUnsignedMax is the largest value of an unsigned 24-bit MEDIUMINT.
const mediumIntUnsignedMax = 16777215

// TinyInt represents a MySQL TINYINT value.
type TinyInt struct {
	// Int8 is the integer value.
	Int8 int8
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a TinyInt.
func (dst *TinyInt) Scan(src any) error {
	valid, err := scanSigned[int8](src, &dst.Int8)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src TinyInt) Value() (driver.Value, error) {
	return signedValue(src.Int8, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src TinyInt) MarshalJSON() ([]byte, error) {
	return marshalSignedJSON(src.Int8, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is out of range, or is not
//     null.
func (dst *TinyInt) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int8](data)
	if err != nil {
		return err
	}
	*dst = TinyInt{Int8: v, Valid: valid}
	return nil
}

// TinyIntUnsigned represents a MySQL TINYINT UNSIGNED value.
type TinyIntUnsigned struct {
	// Uint8 is the integer value.
	Uint8 uint8
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a TinyIntUnsigned or is
//     negative.
func (dst *TinyIntUnsigned) Scan(src any) error {
	valid, err := scanUnsigned[uint8](src, &dst.Uint8)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src TinyIntUnsigned) Value() (driver.Value, error) {
	return unsignedValue(src.Uint8, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src TinyIntUnsigned) MarshalJSON() ([]byte, error) {
	return marshalUnsignedJSON(src.Uint8, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is negative, is out of
//     range, or is not null.
func (dst *TinyIntUnsigned) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalUnsignedJSON[uint8](data)
	if err != nil {
		return err
	}
	*dst = TinyIntUnsigned{Uint8: v, Valid: valid}
	return nil
}

// SmallInt represents a MySQL SMALLINT value.
type SmallInt struct {
	// Int16 is the integer value.
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
//   - error: non-nil when src cannot be scanned into a SmallInt.
func (dst *SmallInt) Scan(src any) error {
	valid, err := scanSigned[int16](src, &dst.Int16)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src SmallInt) Value() (driver.Value, error) {
	return signedValue(src.Int16, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src SmallInt) MarshalJSON() ([]byte, error) {
	return marshalSignedJSON(src.Int16, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is out of range, or is not
//     null.
func (dst *SmallInt) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int16](data)
	if err != nil {
		return err
	}
	*dst = SmallInt{Int16: v, Valid: valid}
	return nil
}

// SmallIntUnsigned represents a MySQL SMALLINT UNSIGNED value.
type SmallIntUnsigned struct {
	// Uint16 is the integer value.
	Uint16 uint16
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a SmallIntUnsigned or is
//     negative.
func (dst *SmallIntUnsigned) Scan(src any) error {
	valid, err := scanUnsigned[uint16](src, &dst.Uint16)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src SmallIntUnsigned) Value() (driver.Value, error) {
	return unsignedValue(src.Uint16, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src SmallIntUnsigned) MarshalJSON() ([]byte, error) {
	return marshalUnsignedJSON(src.Uint16, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is negative, is out of
//     range, or is not null.
func (dst *SmallIntUnsigned) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalUnsignedJSON[uint16](data)
	if err != nil {
		return err
	}
	*dst = SmallIntUnsigned{Uint16: v, Valid: valid}
	return nil
}

// MediumInt represents a MySQL MEDIUMINT value, a signed 24-bit integer.
type MediumInt struct {
	// Int32 is the integer value.
	Int32 int32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a MediumInt or is outside
//     the signed 24-bit range.
func (dst *MediumInt) Scan(src any) error {
	valid, err := scanSigned[int32](src, &dst.Int32)
	if err != nil {
		return err
	}

	if valid && (dst.Int32 < mediumIntMin || dst.Int32 > mediumIntMax) {
		v := dst.Int32
		*dst = MediumInt{}
		return fmt.Errorf("%w: %d does not fit in %s", ErrOutOfRange, v, TypeMediumInt)
	}

	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: non-nil when the value is outside the signed 24-bit range.
func (src MediumInt) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	if src.Int32 < mediumIntMin || src.Int32 > mediumIntMax {
		return nil, fmt.Errorf("%w: %d does not fit in %s", ErrOutOfRange, src.Int32, TypeMediumInt)
	}
	return int64(src.Int32), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value is outside the signed 24-bit range.
func (src MediumInt) MarshalJSON() ([]byte, error) {
	if src.Valid && (src.Int32 < mediumIntMin || src.Int32 > mediumIntMax) {
		return nil, fmt.Errorf("%w: %d does not fit in %s", ErrOutOfRange, src.Int32, TypeMediumInt)
	}
	return marshalSignedJSON(src.Int32, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is outside the signed
//     24-bit range, or is not null.
func (dst *MediumInt) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int32](data)
	if err != nil {
		return err
	}
	if valid && (v < mediumIntMin || v > mediumIntMax) {
		return fmt.Errorf("%w: %d does not fit in %s", ErrOutOfRange, v, TypeMediumInt)
	}
	*dst = MediumInt{Int32: v, Valid: valid}
	return nil
}

// MediumIntUnsigned represents a MySQL MEDIUMINT UNSIGNED value, an unsigned
// 24-bit integer.
type MediumIntUnsigned struct {
	// Uint32 is the integer value.
	Uint32 uint32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a MediumIntUnsigned or is
//     outside the unsigned 24-bit range.
func (dst *MediumIntUnsigned) Scan(src any) error {
	valid, err := scanUnsigned[uint32](src, &dst.Uint32)
	if err != nil {
		return err
	}

	if valid && dst.Uint32 > mediumIntUnsignedMax {
		v := dst.Uint32
		*dst = MediumIntUnsigned{}
		return fmt.Errorf("%w: %d does not fit in %s unsigned", ErrOutOfRange, v, TypeMediumInt)
	}

	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: non-nil when the value is outside the unsigned 24-bit range.
func (src MediumIntUnsigned) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	if src.Uint32 > mediumIntUnsignedMax {
		return nil, fmt.Errorf("%w: %d does not fit in %s unsigned", ErrOutOfRange, src.Uint32, TypeMediumInt)
	}
	return uint64(src.Uint32), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value is outside the unsigned 24-bit range.
func (src MediumIntUnsigned) MarshalJSON() ([]byte, error) {
	if src.Valid && src.Uint32 > mediumIntUnsignedMax {
		return nil, fmt.Errorf("%w: %d does not fit in %s unsigned", ErrOutOfRange, src.Uint32, TypeMediumInt)
	}
	return marshalUnsignedJSON(src.Uint32, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is outside the unsigned
//     24-bit range, or is not null.
func (dst *MediumIntUnsigned) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalUnsignedJSON[uint32](data)
	if err != nil {
		return err
	}
	if valid && v > mediumIntUnsignedMax {
		return fmt.Errorf("%w: %d does not fit in %s unsigned", ErrOutOfRange, v, TypeMediumInt)
	}
	*dst = MediumIntUnsigned{Uint32: v, Valid: valid}
	return nil
}

// Int represents a MySQL INT value. INTEGER is a synonym for INT.
type Int struct {
	// Int32 is the integer value.
	Int32 int32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an Int.
func (dst *Int) Scan(src any) error {
	valid, err := scanSigned[int32](src, &dst.Int32)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src Int) Value() (driver.Value, error) {
	return signedValue(src.Int32, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Int) MarshalJSON() ([]byte, error) {
	return marshalSignedJSON(src.Int32, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is out of range, or is not
//     null.
func (dst *Int) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int32](data)
	if err != nil {
		return err
	}
	*dst = Int{Int32: v, Valid: valid}
	return nil
}

// IntUnsigned represents a MySQL INT UNSIGNED value.
type IntUnsigned struct {
	// Uint32 is the integer value.
	Uint32 uint32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into an IntUnsigned or is
//     negative.
func (dst *IntUnsigned) Scan(src any) error {
	valid, err := scanUnsigned[uint32](src, &dst.Uint32)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src IntUnsigned) Value() (driver.Value, error) {
	return unsignedValue(src.Uint32, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src IntUnsigned) MarshalJSON() ([]byte, error) {
	return marshalUnsignedJSON(src.Uint32, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is negative, is out of
//     range, or is not null.
func (dst *IntUnsigned) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalUnsignedJSON[uint32](data)
	if err != nil {
		return err
	}
	*dst = IntUnsigned{Uint32: v, Valid: valid}
	return nil
}

// BigInt represents a MySQL BIGINT value.
type BigInt struct {
	// Int64 is the integer value.
	Int64 int64
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a BigInt.
func (dst *BigInt) Scan(src any) error {
	valid, err := scanSigned[int64](src, &dst.Int64)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src BigInt) Value() (driver.Value, error) {
	return signedValue(src.Int64, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src BigInt) MarshalJSON() ([]byte, error) {
	return marshalSignedJSON(src.Int64, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number or is not null.
func (dst *BigInt) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalSignedJSON[int64](data)
	if err != nil {
		return err
	}
	*dst = BigInt{Int64: v, Valid: valid}
	return nil
}

// BigIntUnsigned represents a MySQL BIGINT UNSIGNED value.
type BigIntUnsigned struct {
	// Uint64 is the integer value.
	Uint64 uint64
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a BigIntUnsigned or is
//     negative.
func (dst *BigIntUnsigned) Scan(src any) error {
	valid, err := scanUnsigned[uint64](src, &dst.Uint64)
	if err != nil {
		return err
	}
	dst.Valid = valid
	return nil
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the integer value, or nil when the value is NULL.
//   - error: always nil.
func (src BigIntUnsigned) Value() (driver.Value, error) {
	return unsignedValue(src.Uint64, src.Valid)
}

// MarshalJSON implements the encoding/json.Marshaler interface.
//
// Returns:
//   - []byte: the JSON number, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src BigIntUnsigned) MarshalJSON() ([]byte, error) {
	return marshalUnsignedJSON(src.Uint64, src.Valid)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON number, is negative, or is not
//     null.
func (dst *BigIntUnsigned) UnmarshalJSON(data []byte) error {
	v, valid, err := unmarshalUnsignedJSON[uint64](data)
	if err != nil {
		return err
	}
	*dst = BigIntUnsigned{Uint64: v, Valid: valid}
	return nil
}
