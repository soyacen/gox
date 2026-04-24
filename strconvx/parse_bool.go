package strconvx

import (
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ParseBool converts a string to a boolean value.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - bool: the parsed boolean value.
//   - error: an error if parsing fails.
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ParseBoolSlice converts a slice of strings to a slice of boolean values.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to parse.
//
// Returns:
//   - []bool: the parsed boolean values.
//   - error: an error if any element fails to parse.
func ParseBoolSlice(s []string) ([]bool, error) {
	if s == nil {
		return nil, nil
	}
	r := make([]bool, 0, len(s))
	for _, str := range s {
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		r = append(r, b)
	}
	return r, nil
}

// ParseWrapperBool converts a string to a protobuf BoolValue wrapper.
//
// Parameters:
//   - s: the string to parse.
//
// Returns:
//   - *wrapperspb.BoolValue: the protobuf wrapper containing the parsed value.
//   - error: an error if parsing fails.
func ParseWrapperBool(s string) (*wrapperspb.BoolValue, error) {
	v, err := ParseBool(s)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Bool(v), nil
}
