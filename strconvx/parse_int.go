package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ParseInt[Signed constraints.Signed](s string, base int, bitSize int) (Signed, error) {
	i, err := strconv.ParseInt(s, base, bitSize)
	return Signed(i), err
}

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

func ParseWrapperInt32(s string) (*wrapperspb.Int32Value, error) {
	i, err := ParseInt[int32](s, 10, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Int32(i), nil
}

func ParseWrapperInt64(s string) (*wrapperspb.Int64Value, error) {
	i, err := ParseInt[int64](s, 10, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Int64(i), nil
}
