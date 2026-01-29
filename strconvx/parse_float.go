package strconvx

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ParseFloat[Float constraints.Float](s string, bitSize int) (Float, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	return Float(f), err
}

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

func ParseWrapperFloat32(s string) (*wrapperspb.FloatValue, error) {
	v, err := ParseFloat[float32](s, 32)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Float(v), nil
}

func ParseWrapperFloat64(s string) (*wrapperspb.DoubleValue, error) {
	v, err := ParseFloat[float64](s, 64)
	if err != nil {
		return nil, err
	}
	return wrapperspb.Double(v), nil
}
