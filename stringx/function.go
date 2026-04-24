package stringx

import "strings"

// Max returns the greater of a and b.
//
// Parameters:
//   - a: the first string.
//   - b: the second string.
//
// Returns:
//   - S: the greater string.
func Max[S ~string](a, b S) S {
	if strings.Compare(string(a), string(b)) >= 0 {
		return a
	}
	return b
}

// Min returns the lesser of a and b.
//
// Parameters:
//   - a: the first string.
//   - b: the second string.
//
// Returns:
//   - S: the lesser string.
func Min[S ~string](a, b S) S {
	if strings.Compare(string(a), string(b)) <= 0 {
		return a
	}
	return b
}
