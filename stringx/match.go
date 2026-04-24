package stringx

import "regexp"

// MatchAll returns all matches of pattern in str.
//
// Parameters:
//   - str: the source string.
//   - pattern: the regexp pattern.
//
// Returns:
//   - []string: all matches.
func MatchAll(str string, pattern string) []string {
	return regexp.MustCompile(pattern).FindAllString(str, -1)
}

// Match returns true if pattern matches the string.
//
// Parameters:
//   - s: the source string.
//   - pattern: the regexp pattern.
//
// Returns:
//   - bool: true if pattern matches.
func Match(s, pattern string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(s)
}
