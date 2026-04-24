package stringx

// ReplaceFunc replaces each rune in str with the result of f.
//
// Parameters:
//   - str: the source string.
//   - f: the replacement function.
//
// Returns:
//   - string: the replaced string.
func ReplaceFunc(str string, f func(rune) string) string {
	var b Builder
	for _, r := range str {
		_, _ = b.WriteString(f(r))
	}
	return b.String()
}
