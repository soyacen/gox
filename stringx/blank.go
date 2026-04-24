package stringx

import "strings"

// IsBlank checks if a string is empty ("") or whitespace only.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is blank.
func IsBlank[S ~string](s S) bool {
	return len(strings.TrimSpace(string(s))) == 0
}

// IsNotBlank checks if a string is not empty ("") and not whitespace only.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is not blank.
func IsNotBlank[S ~string](s S) bool {
	return !IsBlank(s)
}

// IsAllBlank checks if all of the CharSequences are empty ("") or whitespace only.
//
// Parameters:
//   - ss: the strings to check.
//
// Returns:
//   - bool: true if all strings are blank.
func IsAllBlank[S ~string](ss ...S) bool {
	for _, s := range ss {
		if IsNotBlank(s) {
			return false
		}
	}
	return true
}

// IsAnyBlank checks if any of the string are empty ("") or whitespace only.
//
// Parameters:
//   - ss: the strings to check.
//
// Returns:
//   - bool: true if any string is blank.
func IsAnyBlank[S ~string](ss ...S) bool {
	for _, s := range ss {
		if IsBlank(s) {
			return true
		}
	}
	return false
}
