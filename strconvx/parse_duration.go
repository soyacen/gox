package strconvx

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func ParseWrapperDuration(s string) (*durationpb.Duration, error) {
	v, err := ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return durationpb.New(v), nil
}
