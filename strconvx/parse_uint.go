package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ParseUint converts a string to an unsigned integer.
//
// Parameters:
//   - s: the string to parse.
//   - base: the base (0, 2–36) for the string representation.
//   - bitSize: the bit size of the target type (0, 8, 16, 32, 64).
//
// Returns:
//   - Unsigned: the parsed unsigned integer value.
//   - error: an error if parsing fails.
func ParseUint[Unsigned constraints.Unsigned](s string, base int, bitSize int) (Unsigned, error) {
	i, err := strconv.ParseUint(s, base, bitSize)
	return Unsigned(i), err
}

// ParseUintSlice converts a slice of strings to a slice of unsigned integers.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to parse.
//   - base: the base (0, 2–36) for the string representation.
//   - bitSize: the bit size of the target type (0, 8, 16, 32, 64).
//
// Returns:
//   - []Unsigned: the parsed unsigned integer values.
//   - error: an error if any element fails to parse.
func ParseUintSlice[Unsigned constraints.Unsigned](s []string, base int, bitSize int) ([]Unsigned, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Unsigned, 0, len(s))
	for _, str := range s {
		i, err := ParseUint[Unsigned](str, base, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, i)
	}
	return r, nil
}

// ParseWrapperUint32 converts a string to a protobuf UInt32Value wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.UInt32Value: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperUint32(s string) (*wrapperspb.UInt32Value, error) {
	v, err := ParseUint[uint32](s, 10, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.UInt32(v), nil
}

// ParseWrapperUint64 converts a string to a protobuf UInt64Value wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.UInt64Value: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperUint64(s string) (*wrapperspb.UInt64Value, error) {
	v, err := ParseUint[uint64](s, 10, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.UInt64(uint64(v)), nil
}
