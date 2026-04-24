package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// FormatBool converts a boolean value to its string representation.
//
// Parameters:
//   - b: the boolean value to format.
//
// Returns:
//   - string: the string "true" or "false".
func FormatBool[Bool ~bool](b Bool) string {
	return strconv.FormatBool(bool(b))
}

// FormatUint converts an unsigned integer to a string in the specified base.
//
// Parameters:
//   - i: the unsigned integer value to format.
//   - base: the base (2–36) for the string representation.
//
// Returns:
//   - string: the formatted string representation.
func FormatUint[Unsigned constraints.Unsigned](i Unsigned, base int) string {
	return strconv.FormatUint(uint64(i), base)
}

// FormatInt converts a signed integer to a string in the specified base.
//
// Parameters:
//   - i: the signed integer value to format.
//   - base: the base (2–36) for the string representation.
//
// Returns:
//   - string: the formatted string representation.
func FormatInt[Signed constraints.Signed](i Signed, base int) string {
	return strconv.FormatInt(int64(i), base)
}

// FormatFloat converts a floating-point number to a string.
//
// Parameters:
//   - f: the floating-point value to format.
//   - fmt: the format ('b', 'e', 'E', 'f', 'g', 'G', 'x', 'X').
//   - prec: the precision.
//   - bitSize: the bit size of the original type (32 or 64).
//
// Returns:
//   - string: the formatted string representation.
func FormatFloat[Float constraints.Float](f Float, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(float64(f), fmt, prec, bitSize)
}

// FormatBoolSlice converts a slice of boolean values to a slice of strings.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of boolean values to format.
//
// Returns:
//   - []string: the formatted string slice.
func FormatBoolSlice[Bool ~bool](s []Bool) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, b := range s {
		r = append(r, FormatBool(b))
	}
	return r
}

// FormatUintSlice converts a slice of unsigned integers to a slice of strings.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of unsigned integer values to format.
//   - base: the base (2–36) for the string representation.
//
// Returns:
//   - []string: the formatted string slice.
func FormatUintSlice[Unsigned constraints.Unsigned](s []Unsigned, base int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, i := range s {
		r = append(r, FormatUint(i, base))
	}
	return r
}

// FormatIntSlice converts a slice of signed integers to a slice of strings.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of signed integer values to format.
//   - base: the base (2–36) for the string representation.
//
// Returns:
//   - []string: the formatted string slice.
func FormatIntSlice[Signed constraints.Signed](s []Signed, base int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, i := range s {
		r = append(r, FormatInt(i, base))
	}
	return r
}

// FormatFloatSlice converts a slice of floating-point numbers to a slice of strings.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of floating-point values to format.
//   - fmt: the format ('b', 'e', 'E', 'f', 'g', 'G', 'x', 'X').
//   - prec: the precision.
//   - bitSize: the bit size of the original type (32 or 64).
//
// Returns:
//   - []string: the formatted string slice.
func FormatFloatSlice[Float constraints.Float](s []Float, fmt byte, prec, bitSize int) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, f := range s {
		r = append(r, FormatFloat(float64(f), fmt, prec, bitSize))
	}
	return r
}
