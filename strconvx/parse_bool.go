package strconvx

import (
	"strconv"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

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

func ParseWrapperBool(s string) (*wrapperspb.BoolValue, error) {
	v, err := ParseBool(s)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Bool(v), nil
}
