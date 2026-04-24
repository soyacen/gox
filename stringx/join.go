package stringx

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// JoinInt joins a slice of signed integers into a string.
//
// Parameters:
//   - es: the slice of integers.
//   - sep: the separator.
//
// Returns:
//   - string: the joined string.
func JoinInt[E constraints.Signed](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatInt(int64(e), 10))
	}
	return strings.Join(s, sep)
}

// JoinUint joins a slice of unsigned integers into a string.
//
// Parameters:
//   - es: the slice of unsigned integers.
//   - sep: the separator.
//
// Returns:
//   - string: the joined string.
func JoinUint[E constraints.Unsigned](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatUint(uint64(e), 10))
	}
	return strings.Join(s, sep)
}

// JoinFloat joins a slice of floats into a string.
//
// Parameters:
//   - es: the slice of floats.
//   - sep: the separator.
//
// Returns:
//   - string: the joined string.
func JoinFloat[E constraints.Float](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatFloat(float64(e), 'f', -1, 64))
	}
	return strings.Join(s, sep)
}

// JoinBool joins a slice of booleans into a string.
//
// Parameters:
//   - es: the slice of booleans.
//   - sep: the separator.
//
// Returns:
//   - string: the joined string.
func JoinBool[E ~bool](es []E, sep string) string {
	s := make([]string, 0, len(es))
	for _, e := range es {
		s = append(s, strconv.FormatBool(bool(e)))
	}
	return strings.Join(s, sep)
}
