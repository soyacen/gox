package stringx

import (
	"strings"
)

// Indices returns all starting indices of substr in str.
//
// Parameters:
//   - str: the source string.
//   - substr: the substring to find.
//
// Returns:
//   - []int: all starting indices of substr.
func Indices(str, substr string) []int {
	s := str
	var indices []int
	subStrLen := len(substr)
	var index int
	for {
		i := strings.Index(s, substr)
		if i < 0 {
			break
		}
		indices = append(indices, i+index)
		s = s[i+subStrLen:]
		index = i + index + subStrLen
	}
	return indices
}
