package stringx

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

// GoCamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the camel-cased string.
func GoCamelCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
		case c == '.':
			b = append(b, '_')
		case c == '_' && (i == 0 || s[i-1] == '.'):
			b = append(b, 'X')
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}

// GoSanitized converts a string to a valid Go identifier.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: a valid Go identifier.
func GoSanitized(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}

// JSONCamelCase converts a snake_case identifier to a camelCase identifier,
// according to the protobuf JSON specification.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the camel-cased string.
func JSONCamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}

// JSONSnakeCase converts a camelCase identifier to a snake_case identifier,
// according to the protobuf JSON specification.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the snake-cased string.
func JSONSnakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isASCIIUpper(c) {
			b = append(b, '_')
			c += 'a' - 'A'
		}
		b = append(b, c)
	}
	return string(b)
}

// MapEntryName derives the name of the map entry message given the field name.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:254-276,6057
//
// Parameters:
//   - s: the field name.
//
// Returns:
//   - string: the map entry message name.
func MapEntryName(s string) string {
	var b []byte
	upperNext := true
	for _, c := range s {
		switch {
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(c))
		}
	}
	b = append(b, "Entry"...)
	return string(b)
}

// EnumValueName derives the camel-cased enum value name.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:297-313
//
// Parameters:
//   - s: the enum value.
//
// Returns:
//   - string: the camel-cased enum value name.
func EnumValueName(s string) string {
	var b []byte
	upperNext := true
	for _, c := range s {
		switch {
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(unicode.ToLower(c)))
			upperNext = false
		}
	}
	return string(b)
}

// TrimEnumPrefix trims the enum name prefix from an enum value name,
// where the prefix is all lowercase without underscores.
// See protoc v3.8.0: src/google/protobuf/descriptor.cc:330-375
//
// Parameters:
//   - s: the enum value name.
//   - prefix: the prefix to trim.
//
// Returns:
//   - string: the trimmed enum value name.
func TrimEnumPrefix(s, prefix string) string {
	s0 := s
	for len(s) > 0 && len(prefix) > 0 {
		if s[0] == '_' {
			s = s[1:]
			continue
		}
		if unicode.ToLower(rune(s[0])) != rune(prefix[0]) {
			return s0
		}
		s, prefix = s[1:], prefix[1:]
	}
	if len(prefix) > 0 {
		return s0
	}
	s = strings.TrimLeft(s, "_")
	if len(s) == 0 {
		return s0
	}
	return s
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// KebabCase converts an arbitrary string to a kebab-case string that complies with domain name specifications.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the kebab-cased string.
func KebabCase(s string) string {
	return formatCase(s, '-')
}

// DotCase converts an arbitrary string to a dot-case string (e.g., service.name).
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the dot-cased string.
func DotCase(s string) string {
	return formatCase(s, '.')
}

func formatCase(s string, sep byte) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))

	var lastWasSep bool
	var prevRune rune

	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			if unicode.IsLower(prevRune) || unicode.IsDigit(prevRune) {
				if b.Len() > 0 && !lastWasSep {
					b.WriteByte(sep)
					lastWasSep = true
				}
			} else if unicode.IsUpper(prevRune) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				if b.Len() > 0 && !lastWasSep {
					b.WriteByte(sep)
					lastWasSep = true
				}
			}
		}

		if unicode.IsLower(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastWasSep = false
		} else if unicode.IsUpper(r) {
			b.WriteRune(unicode.ToLower(r))
			lastWasSep = false
		} else {
			if b.Len() > 0 && !lastWasSep {
				b.WriteByte(sep)
				lastWasSep = true
			}
		}

		if b.Len() >= 63 {
			break
		}
		prevRune = r
	}

	res := b.String()
	res = strings.Trim(res, string(sep))

	if len(res) > 63 {
		res = res[:63]
		res = strings.TrimRight(res, string(sep))
	}

	return res
}
