package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ParseFloat converts a string to a floating-point value.
//
// Parameters:
//   - s: the string to parse.
//   - bitSize: the bit size of the target type (32 or 64).
//
// Returns:
//   - Float: the parsed floating-point value.
//   - error: an error if parsing fails.
func ParseFloat[Float constraints.Float](s string, bitSize int) (Float, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	return Float(f), err
}

// ParseFloatSlice converts a slice of strings to a slice of floating-point values.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to parse.
//   - bitSize: the bit size of the target type (32 or 64).
//
// Returns:
//   - []Float: the parsed floating-point values.
//   - error: an error if any element fails to parse.
func ParseFloatSlice[Float constraints.Float](s []string, bitSize int) ([]Float, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]Float, 0, len(s))
	for _, str := range s {
		f, err := ParseFloat[Float](str, bitSize)
		if err != nil {
			return nil, err
		}
		r = append(r, f)
	}
	return r, nil
}

// ParseWrapperFloat32 converts a string to a protobuf FloatValue wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.FloatValue: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperFloat32(s string) (*wrapperspb.FloatValue, error) {
	v, err := ParseFloat[float32](s, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Float(v), nil
}

// ParseWrapperFloat64 converts a string to a protobuf DoubleValue wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.DoubleValue: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperFloat64(s string) (*wrapperspb.DoubleValue, error) {
	v, err := ParseFloat[float64](s, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Double(v), nil
}
