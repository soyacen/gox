package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ParseUint[Unsigned constraints.Unsigned](s string, base int, bitSize int) (Unsigned, error) {
	i, err := strconv.ParseUint(s, base, bitSize)
	return Unsigned(i), err
}

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

func ParseWrapperUint32(s string) (*wrapperspb.UInt32Value, error) {
	v, err := ParseUint[uint32](s, 10, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.UInt32(v), nil
}

func ParseWrapperUint64(s string) (*wrapperspb.UInt64Value, error) {
	v, err := ParseUint[uint64](s, 10, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.UInt64(uint64(v)), nil
}
