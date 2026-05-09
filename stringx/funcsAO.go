package stringx

import (
	"fmt"
	"html"

	//"log"
	"regexp"
	"strings"
)

// Verbose enables console output for functions that have counterparts in Go's standard packages.
var (
	Verbose       = false
	templateOpen  = "{{"
	templateClose = "}}"
)

var (
	beginEndSpacesRe  = regexp.MustCompile("^\\s+|\\s+$")
	camelizeRe        = regexp.MustCompile(`(\-|_|\s)+(.)?`)
	camelizeRe2       = regexp.MustCompile(`(\-|_|\s)+`)
	capitalsRe        = regexp.MustCompile("([A-Z])")
	dashSpaceRe       = regexp.MustCompile(`[-\s]+`)
	dashesRe          = regexp.MustCompile("-+")
	isAlphaNumericRe  = regexp.MustCompile(`[^0-9a-z\xC0-\xFF]`)
	isAlphaRe         = regexp.MustCompile(`[^a-z\xC0-\xFF]`)
	nWhitespaceRe     = regexp.MustCompile(`\s+`)
	notDigitsRe       = regexp.MustCompile(`[^0-9]`)
	slugifyRe         = regexp.MustCompile(`[^\w\s\-]`)
	spaceUnderscoreRe = regexp.MustCompile("[_\\s]+")
	spacesRe          = regexp.MustCompile(`[\s ]+`)
	stripPuncRe       = regexp.MustCompile(`[^\w\s]|_`)
	templateRe        = regexp.MustCompile(`([\-\[\]()*\s])`)
	templateRe2       = regexp.MustCompile(`\$`)
	underscoreRe      = regexp.MustCompile(`([a-z\d])([A-Z]+)`)
	whitespaceRe      = regexp.MustCompile(`^[\s ]*$`)
)

// Between extracts a string between left and right strings.
//
// Parameters:
//   - s: the source string.
//   - left: the left delimiter.
//   - right: the right delimiter.
//
// Returns:
//   - string: the substring between left and right.
func Between(s, left, right string) string {
	l := len(left)
	startPos := strings.Index(s, left)
	if startPos < 0 {
		return ""
	}
	endPos := IndexOf(s, right, startPos+l)
	if endPos < 0 {
		return ""
	} else if right == "" {
		return s[endPos:]
	} else {
		return s[startPos+l : endPos]
	}
}

// BetweenF is the filter form for Between.
//
// Parameters:
//   - left: the left delimiter.
//   - right: the right delimiter.
//
// Returns:
//   - func(string) string: a function that extracts the substring.
func BetweenF(left, right string) func(string) string {
	return func(s string) string {
		return Between(s, left, right)
	}
}

// Camelize returns a new string which removes any underscores or dashes and converts a string into camel casing.
//
// Parameters:
//   - s: the string to camelize.
//
// Returns:
//   - string: the camelized string.
func Camelize(s string) string {
	return camelizeRe.ReplaceAllStringFunc(s, func(val string) string {
		val = strings.ToUpper(val)
		val = camelizeRe2.ReplaceAllString(val, "")
		return val
	})
}

// Capitalize uppercases the first char of s and lowercases the rest.
//
// Parameters:
//   - s: the string to capitalize.
//
// Returns:
//   - string: the capitalized string.
func Capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

// CharAt returns a string from the character at the specified position.
//
// Parameters:
//   - s: the source string.
//   - index: the position.
//
// Returns:
//   - string: the character at the specified position.
func CharAt(s string, index int) string {
	l := len(s)
	shortcut := index < 0 || index > l-1 || l == 0
	if shortcut {
		return ""
	}
	return s[index : index+1]
}

// CharAtF is the filter form of CharAt.
//
// Parameters:
//   - index: the position.
//
// Returns:
//   - func(string) string: a function that returns the character at the specified position.
func CharAtF(index int) func(string) string {
	return func(s string) string {
		return CharAt(s, index)
	}
}

// ChompLeft removes prefix at the start of a string.
//
// Parameters:
//   - s: the source string.
//   - prefix: the prefix to remove.
//
// Returns:
//   - string: the string with the prefix removed.
func ChompLeft(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// ChompLeftF is the filter form of ChompLeft.
//
// Parameters:
//   - prefix: the prefix to remove.
//
// Returns:
//   - func(string) string: a function that removes the prefix.
func ChompLeftF(prefix string) func(string) string {
	return func(s string) string {
		return ChompLeft(s, prefix)
	}
}

// ChompRight removes suffix from end of s.
//
// Parameters:
//   - s: the source string.
//   - suffix: the suffix to remove.
//
// Returns:
//   - string: the string with the suffix removed.
func ChompRight(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// ChompRightF is the filter form of ChompRight.
//
// Parameters:
//   - suffix: the suffix to remove.
//
// Returns:
//   - func(string) string: a function that removes the suffix.
func ChompRightF(suffix string) func(string) string {
	return func(s string) string {
		return ChompRight(s, suffix)
	}
}

// Classify returns a camelized string with the first letter upper cased.
//
// Parameters:
//   - s: the string to classify.
//
// Returns:
//   - string: the classified string.
func Classify(s string) string {
	return Camelize("-" + s)
}

// ClassifyF is the filter form of Classify.
//
// Parameters:
//   - s: the string to classify.
//
// Returns:
//   - func(string) string: a function that classifies the string.
func ClassifyF(s string) func(string) string {
	return func(s string) string {
		return Classify(s)
	}
}

// Clean compresses all adjacent whitespace to a single space and trims s.
//
// Parameters:
//   - s: the string to clean.
//
// Returns:
//   - string: the cleaned string.
func Clean(s string) string {
	s = spacesRe.ReplaceAllString(s, " ")
	s = beginEndSpacesRe.ReplaceAllString(s, "")
	return s
}

// Dasherize converts a camel cased string into a string delimited by dashes.
//
// Parameters:
//   - s: the string to dasherize.
//
// Returns:
//   - string: the dasherized string.
func Dasherize(s string) string {
	s = strings.TrimSpace(s)
	s = spaceUnderscoreRe.ReplaceAllString(s, "-")
	s = capitalsRe.ReplaceAllString(s, "-$1")
	s = dashesRe.ReplaceAllString(s, "-")
	s = strings.ToLower(s)
	return s
}

// EscapeHTML is an alias for html.EscapeString.
//
// Parameters:
//   - s: the string to escape.
//
// Returns:
//   - string: the escaped string.
func EscapeHTML(s string) string {
	if Verbose {
		fmt.Println("Use html.EscapeString instead of EscapeHTML")
	}
	return html.EscapeString(s)
}

// DecodeHTMLEntities decodes HTML entities into their proper string representation.
// DecodeHTMLEntities is an alias for html.UnescapeString.
//
// Parameters:
//   - s: the string to decode.
//
// Returns:
//   - string: the decoded string.
func DecodeHTMLEntities(s string) string {
	if Verbose {
		fmt.Println("Use html.UnescapeString instead of DecodeHTMLEntities")
	}
	return html.UnescapeString(s)
}

// EnsurePrefix ensures s starts with prefix.
//
// Parameters:
//   - s: the source string.
//   - prefix: the prefix to ensure.
//
// Returns:
//   - string: the string with the prefix.
func EnsurePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// EnsurePrefixF is the filter form of EnsurePrefix.
//
// Parameters:
//   - prefix: the prefix to ensure.
//
// Returns:
//   - func(string) string: a function that ensures the prefix.
func EnsurePrefixF(prefix string) func(string) string {
	return func(s string) string {
		return EnsurePrefix(s, prefix)
	}
}

// EnsureSuffix ensures s ends with suffix.
//
// Parameters:
//   - s: the source string.
//   - suffix: the suffix to ensure.
//
// Returns:
//   - string: the string with the suffix.
func EnsureSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// EnsureSuffixF is the filter form of EnsureSuffix.
//
// Parameters:
//   - suffix: the suffix to ensure.
//
// Returns:
//   - func(string) string: a function that ensures the suffix.
func EnsureSuffixF(suffix string) func(string) string {
	return func(s string) string {
		return EnsureSuffix(s, suffix)
	}
}

// Humanize transforms s into a human friendly form.
//
// Parameters:
//   - s: the string to humanize.
//
// Returns:
//   - string: the humanized string.
func Humanize(s string) string {
	if s == "" {
		return s
	}
	s = Underscore(s)
	humanizeRe := regexp.MustCompile(`_id$`)
	s = humanizeRe.ReplaceAllString(s, "")
	s = strings.Replace(s, "_", " ", -1)
	s = strings.TrimSpace(s)
	s = Capitalize(s)
	return s
}

// Iif is short for immediate if. If condition is true return truthy else falsey.
//
// Parameters:
//   - condition: the condition.
//   - truthy: the value to return if condition is true.
//   - falsey: the value to return if condition is false.
//
// Returns:
//   - string: truthy or falsey.
func Iif(condition bool, truthy string, falsey string) string {
	if condition {
		return truthy
	}
	return falsey
}

// IndexOf finds the index of needle in s starting from start.
//
// Parameters:
//   - s: the source string.
//   - needle: the substring to find.
//   - start: the starting index.
//
// Returns:
//   - int: the index of needle, or -1 if not found.
func IndexOf(s string, needle string, start int) int {
	l := len(s)
	if needle == "" {
		if start < 0 {
			return 0
		} else if start < l {
			return start
		} else {
			return l
		}
	}
	if start < 0 || start > l-1 {
		return -1
	}
	pos := strings.Index(s[start:], needle)
	if pos == -1 {
		return -1
	}
	return start + pos
}

// IsAlpha returns true if a string contains only letters from ASCII (a-z,A-Z). Other letters from other languages are not supported.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is alphabetic.
func IsAlpha(s string) bool {
	return !isAlphaRe.MatchString(strings.ToLower(s))
}

// IsAlphaNumeric returns true if a string contains letters and digits.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is alphanumeric.
func IsAlphaNumeric(s string) bool {
	return !isAlphaNumericRe.MatchString(strings.ToLower(s))
}

// IsLower returns true if s is comprised of all lower case characters.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is all lowercase.
func IsLower(s string) bool {
	return IsAlpha(s) && s == strings.ToLower(s)
}

// IsNumeric returns true if a string contains only digits from 0-9. Other digits not in Latin (such as Arabic) are not currently supported.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is numeric.
func IsNumeric(s string) bool {
	return !notDigitsRe.MatchString(s)
}

// IsUpper returns true if s contains all upper case characters.
//
// Parameters:
//   - s: the string to check.
//
// Returns:
//   - bool: true if the string is all uppercase.
func IsUpper(s string) bool {
	return IsAlpha(s) && s == strings.ToUpper(s)
}

// Left returns the left substring of length n.
//
// Parameters:
//   - s: the source string.
//   - n: the length.
//
// Returns:
//   - string: the left substring.
func Left(s string, n int) string {
	if n < 0 {
		return Right(s, -n)
	}
	return Substr(s, 0, n)
}

// LeftF is the filter form of Left.
//
// Parameters:
//   - n: the length.
//
// Returns:
//   - func(string) string: a function that returns the left substring.
func LeftF(n int) func(string) string {
	return func(s string) string {
		return Left(s, n)
	}
}

// LeftOf returns the substring left of needle.
//
// Parameters:
//   - s: the source string.
//   - needle: the delimiter.
//
// Returns:
//   - string: the substring left of needle.
func LeftOf(s string, needle string) string {
	return Between(s, "", needle)
}

// Letters returns an array of runes as strings so it can be indexed into.
//
// Parameters:
//   - s: the source string.
//
// Returns:
//   - []string: an array of runes as strings.
func Letters(s string) []string {
	result := []string{}
	for _, r := range s {
		result = append(result, string(r))
	}
	return result
}

// Lines converts windows newlines to unix newlines then converts to an array of lines.
//
// Parameters:
//   - s: the source string.
//
// Returns:
//   - []string: an array of lines.
func Lines(s string) []string {
	s = strings.Replace(s, "\r\n", "\n", -1)
	return strings.Split(s, "\n")
}

// Map maps an array's items through an iterator.
//
// Parameters:
//   - arr: the array to map.
//   - iterator: the mapping function.
//
// Returns:
//   - []string: the mapped array.
func Map(arr []string, iterator func(string) string) []string {
	r := []string{}
	for _, item := range arr {
		r = append(r, iterator(item))
	}
	return r
}
