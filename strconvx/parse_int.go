package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ParseInt converts a string to a signed integer.
//
// Parameters:
//   - s: the string to parse.
//   - base: the base (0, 2–36) for the string representation.
//   - bitSize: the bit size of the target type (0, 8, 16, 32, 64).
//
// Returns:
//   - Signed: the parsed signed integer value.
//   - error: an error if parsing fails.
func ParseInt[Signed constraints.Signed](s string, base int, bitSize int) (Signed, error) {
	i, err := strconv.ParseInt(s, base, bitSize)
	return Signed(i), err
}

// ParseIntSlice converts a slice of strings to a slice of signed integers.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to parse.
//   - base: the base (0, 2–36) for the string representation.
//   - bitSize: the bit size of the target type (0, 8, 16, 32, 64).
//
// Returns:
//   - []Signed: the parsed signed integer values.
//   - error: an error if any element fails to parse.
func ParseIntSlice[Signed constraints.Signed](s []string, base int, bitSize int) ([]Signed, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Signed, 0, len(s))
	for _, str := range s {
		i, err := ParseInt[Signed](str, base, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, i)
	}
	return r, nil
}

// ParseWrapperInt32 converts a string to a protobuf Int32Value wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.Int32Value: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperInt32(s string) (*wrapperspb.Int32Value, error) {
	i, err := ParseInt[int32](s, 10, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Int32(i), nil
}

// ParseWrapperInt64 converts a string to a protobuf Int64Value wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.Int64Value: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperInt64(s string) (*wrapperspb.Int64Value, error) {
	i, err := ParseInt[int64](s, 10, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Int64(i), nil
}
