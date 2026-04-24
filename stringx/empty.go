package stringx

// IsEmpty returns true if the string is solely composed of whitespace.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is empty.
func IsEmpty[S ~string](s S) bool {
	if s == "" {
		return true
	}
	return whitespaceRe.MatchString(string(s))
}

// IsNotEmpty checks if a string is not empty ("").
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is not empty.
func IsNotEmpty[S ~string](s S) bool {
	return !IsEmpty(s)
}

// IsAllEmpty checks if all of the strings are empty ("").
//
// Parameters:
//   - ss: the strings to check.
//
// Returns:
//   - bool: true if all strings are empty.
func IsAllEmpty[S ~string](ss ...S) bool {
	for _, s := range ss {
		if IsNotEmpty(s) {
			return false
		}
	}
	return true
}

// IsAnyEmpty checks if any of the strings are empty ("").
//
// Parameters:
//   - ss: the strings to check.
//
// Returns:
//   - bool: true if any string is empty.
func IsAnyEmpty[S ~string](ss ...S) bool {
	for _, s := range ss {
		if IsEmpty(s) {
			return true
		}
	}
	return false
}
