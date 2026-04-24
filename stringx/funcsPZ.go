package stringx

import (
	"fmt"
	"html"

	//"log"
	"math"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Pad pads string s on both sides with c until it has length of n.
//
// Parameters:
//   - s: the string to pad.
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - string: the padded string.
func Pad(s, c string, n int) string {
	L := len(s)
	if L >= n {
		return s
	}
	n -= L

	left := strings.Repeat(c, int(math.Ceil(float64(n)/2)))
	right := strings.Repeat(c, int(math.Floor(float64(n)/2)))
	return left + s + right
}

// PadF is the filter form of Pad.
//
// Parameters:
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - func(string) string: a function that pads the string.
func PadF(c string, n int) func(string) string {
	return func(s string) string {
		return Pad(s, c, n)
	}
}

// PadLeft pads s on left side with c until it has length of n.
//
// Parameters:
//   - s: the string to pad.
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - string: the left-padded string.
func PadLeft(s, c string, n int) string {
	L := len(s)
	if L > n {
		return s
	}
	return strings.Repeat(c, (n-L)) + s
}

// PadLeftF is the filter form of PadLeft.
//
// Parameters:
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - func(string) string: a function that left-pads the string.
func PadLeftF(c string, n int) func(string) string {
	return func(s string) string {
		return PadLeft(s, c, n)
	}
}

// PadRight pads s on right side with c until it has length of n.
//
// Parameters:
//   - s: the string to pad.
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - string: the right-padded string.
func PadRight(s, c string, n int) string {
	L := len(s)
	if L > n {
		return s
	}
	return s + strings.Repeat(c, n-L)
}

// PadRightF is the filter form of PadRight.
//
// Parameters:
//   - c: the padding character.
//   - n: the desired length.
//
// Returns:
//   - func(string) string: a function that right-pads the string.
func PadRightF(c string, n int) func(string) string {
	return func(s string) string {
		return PadRight(s, c, n)
	}
}

// Pipe pipes s through one or more string filters.
//
// Parameters:
//   - s: the source string.
//   - funcs: the filter functions.
//
// Returns:
//   - string: the piped string.
func Pipe(s string, funcs ...func(string) string) string {
	for _, fn := range funcs {
		s = fn(s)
	}
	return s
}

// QuoteItems quotes all items in array, mostly for debugging.
//
// Parameters:
//   - arr: the array to quote.
//
// Returns:
//   - []string: the quoted items.
func QuoteItems(arr []string) []string {
	return Map(arr, func(s string) string {
		return strconv.Quote(s)
	})
}

// ReplaceF is the filter form of strings.Replace.
//
// Parameters:
//   - old: the substring to replace.
//   - new: the replacement substring.
//   - n: the number of replacements.
//
// Returns:
//   - func(string) string: a function that replaces substrings.
func ReplaceF(old, new string, n int) func(string) string {
	return func(s string) string {
		return strings.Replace(s, old, new, n)
	}
}

// ReplacePattern replaces string with regexp string.
// ReplacePattern returns a copy of src, replacing matches of the Regexp with the replacement string repl. Inside repl, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.
//
// Parameters:
//   - s: the source string.
//   - pattern: the regexp pattern.
//   - repl: the replacement string.
//
// Returns:
//   - string: the replaced string.
func ReplacePattern(s, pattern, repl string) string {
	r := regexp.MustCompile(pattern)
	return r.ReplaceAllString(s, repl)
}

// ReplacePatternF is the filter form of ReplacePattern.
//
// Parameters:
//   - pattern: the regexp pattern.
//   - repl: the replacement string.
//
// Returns:
//   - func(string) string: a function that replaces patterns.
func ReplacePatternF(pattern, repl string) func(string) string {
	return func(s string) string {
		return ReplacePattern(s, pattern, repl)
	}
}

// Reverse reverses a string.
//
// Parameters:
//   - s: the string to reverse.
//
// Returns:
//   - string: the reversed string.
func Reverse(s string) string {
	cs := make([]rune, utf8.RuneCountInString(s))
	i := len(cs)
	for _, c := range s {
		i--
		cs[i] = c
	}
	return string(cs)
}

// Right returns the right substring of length n.
//
// Parameters:
//   - s: the source string.
//   - n: the length.
//
// Returns:
//   - string: the right substring.
func Right(s string, n int) string {
	if n < 0 {
		return Left(s, -n)
	}
	return Substr(s, len(s)-n, n)
}

// RightF is the filter form of Right.
//
// Parameters:
//   - n: the length.
//
// Returns:
//   - func(string) string: a function that returns the right substring.
func RightF(n int) func(string) string {
	return func(s string) string {
		return Right(s, n)
	}
}

// RightOf returns the substring to the right of prefix.
//
// Parameters:
//   - s: the source string.
//   - prefix: the prefix.
//
// Returns:
//   - string: the substring to the right of prefix.
func RightOf(s string, prefix string) string {
	return Between(s, prefix, "")
}

// SetTemplateDelimiters sets the delimiters for Template function. Defaults to "{{" and "}}".
//
// Parameters:
//   - opening: the opening delimiter.
//   - closing: the closing delimiter.
func SetTemplateDelimiters(opening, closing string) {
	templateOpen = opening
	templateClose = closing
}

// Slice slices a string. If end is negative then it is from the end of the string.
//
// Parameters:
//   - s: the source string.
//   - start: the start index.
//   - end: the end index.
//
// Returns:
//   - string: the sliced string.
func Slice(s string, start, end int) string {
	if end > -1 {
		return s[start:end]
	}
	L := len(s)
	if L+end > 0 {
		return s[start : L-end]
	}
	return s[start:]
}

// SliceF is the filter form of Slice.
//
// Parameters:
//   - start: the start index.
//   - end: the end index.
//
// Returns:
//   - func(string) string: a function that slices the string.
func SliceF(start, end int) func(string) string {
	return func(s string) string {
		return Slice(s, start, end)
	}
}

// SliceContains determines whether val is an element in slice.
//
// Parameters:
//   - slice: the slice to search.
//   - val: the value to find.
//
// Returns:
//   - bool: true if val is in the slice.
func SliceContains(slice []string, val string) bool {
	if slice == nil {
		return false
	}

	for _, it := range slice {
		if it == val {
			return true
		}
	}
	return false
}

// SliceIndexOf gets the index of val in slice. Returns -1 if not found.
//
// Parameters:
//   - slice: the slice to search.
//   - val: the value to find.
//
// Returns:
//   - int: the index of val, or -1 if not found.
func SliceIndexOf(slice []string, val string) int {
	if slice == nil {
		return -1
	}

	for i, it := range slice {
		if it == val {
			return i
		}
	}
	return -1
}

// Slugify converts s into a dasherized string suitable for URL segment.
//
// Parameters:
//   - s: the string to slugify.
//
// Returns:
//   - string: the slugified string.
func Slugify(s string) string {
	sl := slugifyRe.ReplaceAllString(s, "")
	sl = strings.ToLower(sl)
	sl = Dasherize(sl)
	return sl
}

// StripPunctuation strips punctuation from string.
//
// Parameters:
//   - s: the string to strip.
//
// Returns:
//   - string: the string without punctuation.
func StripPunctuation(s string) string {
	s = stripPuncRe.ReplaceAllString(s, "")
	s = nWhitespaceRe.ReplaceAllString(s, " ")
	return s
}

// StripTags strips all of the html tags or tags specified by the parameters.
//
// Parameters:
//   - s: the source string.
//   - tags: the tags to strip (optional).
//
// Returns:
//   - string: the string with tags stripped.
func StripTags(s string, tags ...string) string {
	if len(tags) == 0 {
		tags = append(tags, "")
	}
	for _, tag := range tags {
		stripTagsRe := regexp.MustCompile(`(?i)<\/?` + tag + `[^()<>]*>`)
		s = stripTagsRe.ReplaceAllString(s, "")
	}
	return s
}

// Substr returns a substring of s starting at index of length n.
//
// Parameters:
//   - s: the source string.
//   - index: the start index.
//   - n: the length.
//
// Returns:
//   - string: the substring.
func Substr(s string, index int, n int) string {
	L := len(s)
	if index < 0 || index >= L || s == "" {
		return ""
	}
	end := index + n
	if end >= L {
		end = L
	}
	if end <= index {
		return ""
	}
	return s[index:end]
}

// SubstrF is the filter form of Substr.
//
// Parameters:
//   - index: the start index.
//   - n: the length.
//
// Returns:
//   - func(string) string: a function that returns the substring.
func SubstrF(index, n int) func(string) string {
	return func(s string) string {
		return Substr(s, index, n)
	}
}

// Template is a string template which replaces template placeholders delimited by "{{" and "}}" with values from map. The global delimiters may be set with SetTemplateDelimiters.
//
// Parameters:
//   - s: the template string.
//   - values: the map of values.
//
// Returns:
//   - string: the templated string.
func Template(s string, values map[string]interface{}) string {
	return TemplateWithDelimiters(s, values, templateOpen, templateClose)
}

// TemplateDelimiters returns the opening and closing delimiters for Template.
//
// Returns:
//   - opening: the opening delimiter.
//   - closing: the closing delimiter.
func TemplateDelimiters() (opening string, closing string) {
	return templateOpen, templateClose
}

// TemplateWithDelimiters is a string template with user-definable opening and closing delimiters.
//
// Parameters:
//   - s: the template string.
//   - values: the map of values.
//   - opening: the opening delimiter.
//   - closing: the closing delimiter.
//
// Returns:
//   - string: the templated string.
func TemplateWithDelimiters(s string, values map[string]interface{}, opening, closing string) string {
	escapeDelimiter := func(delim string) string {
		result := templateRe.ReplaceAllString(delim, "\\$1")
		return templateRe2.ReplaceAllString(result, "\\$")
	}

	openingDelim := escapeDelimiter(opening)
	closingDelim := escapeDelimiter(closing)
	r := regexp.MustCompile(openingDelim + `(.+?)` + closingDelim)
	matches := r.FindAllStringSubmatch(s, -1)
	for _, submatches := range matches {
		match := submatches[0]
		key := submatches[1]
		if values[key] != nil {
			v := fmt.Sprintf("%v", values[key])
			s = strings.Replace(s, match, v, -1)
		}
	}

	return s
}

// ToArgv converts string s into an argv for exec.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - []string: the argv.
func ToArgv(s string) []string {
	const (
		InArg = iota
		InArgQuote
		OutOfArg
	)
	currentState := OutOfArg
	currentQuoteChar := "\x00"
	currentArg := ""
	argv := []string{}

	isQuote := func(c string) bool {
		return c == `"` || c == `'`
	}

	isEscape := func(c string) bool {
		return c == `\`
	}

	isWhitespace := func(c string) bool {
		return c == " " || c == "\t"
	}

	L := len(s)
	for i := 0; i < L; i++ {
		c := s[i : i+1]

		if isQuote(c) {
			switch currentState {
			case OutOfArg:
				currentArg = ""
				fallthrough
			case InArg:
				currentState = InArgQuote
				currentQuoteChar = c

			case InArgQuote:
				if c == currentQuoteChar {
					currentState = InArg
				} else {
					currentArg += c
				}
			}
		} else if isWhitespace(c) {
			switch currentState {
			case InArg:
				argv = append(argv, currentArg)
				currentState = OutOfArg
			case InArgQuote:
				currentArg += c
			case OutOfArg:
				// nothing
			}
		} else if isEscape(c) {
			switch currentState {
			case OutOfArg:
				currentArg = ""
				currentState = InArg
				fallthrough
			case InArg:
				fallthrough
			case InArgQuote:
				if i == L-1 {
					if runtime.GOOS == "windows" {
						currentArg += c
					} else {
						panic("Escape character at end string")
					}
				} else {
					if runtime.GOOS == "windows" {
						peek := s[i+1 : i+2]
						if peek != `"` {
							currentArg += c
						}
					} else {
						i++
						c = s[i : i+1]
						currentArg += c
					}
				}
			}
		} else {
			switch currentState {
			case InArg, InArgQuote:
				currentArg += c

			case OutOfArg:
				currentArg = ""
				currentArg += c
				currentState = InArg
			}
		}
	}

	if currentState == InArg {
		argv = append(argv, currentArg)
	} else if currentState == InArgQuote {
		panic("Starting quote has no ending quote.")
	}

	return argv
}

// ToBool fuzzily converts truthy values.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - bool: true if the string is truthy.
func ToBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "yes" || s == "on" || s == "1"
}

// ToBoolOr parses s as a bool or returns defaultValue.
//
// Parameters:
//   - s: the string to parse.
//   - defaultValue: the default value.
//
// Returns:
//   - bool: the parsed bool or defaultValue.
func ToBoolOr(s string, defaultValue bool) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return b
}

// ToIntOr parses s as an int or returns defaultValue.
//
// Parameters:
//   - s: the string to parse.
//   - defaultValue: the default value.
//
// Returns:
//   - int: the parsed int or defaultValue.
func ToIntOr(s string, defaultValue int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return n
}

// ToFloat32Or parses s as a float32 or returns defaultValue on error.
//
// Parameters:
//   - s: the string to parse.
//   - defaultValue: the default value.
//
// Returns:
//   - float32: the parsed float32 or defaultValue.
func ToFloat32Or(s string, defaultValue float32) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return defaultValue
	}
	return float32(f)
}

// ToFloat64Or parses s as a float64 or returns defaultValue.
//
// Parameters:
//   - s: the string to parse.
//   - defaultValue: the default value.
//
// Returns:
//   - float64: the parsed float64 or defaultValue.
func ToFloat64Or(s string, defaultValue float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

// ToFloatOr parses s as a float64 or returns defaultValue.
//
// Parameters:
//   - s: the string to parse.
//   - defaultValue: the default value.
//
// Returns:
//   - float64: the parsed float64 or defaultValue.
var ToFloatOr = ToFloat64Or

// TODO This is not working yet. Go's regexp package does not have some
// of the niceities in JavaScript
//
// Truncate truncates the string, accounting for word placement and chars count
// adding a morestr (defaults to ellipsis)
// func Truncate(s, morestr string, n int) string {
// 	L := len(s)
// 	if L <= n {
// 		return s
// 	}
//
// 	if morestr == "" {
// 		morestr = "..."
// 	}
//
// 	tmpl := func(c string) string {
// 		if strings.ToUpper(c) != strings.ToLower(c) {
// 			return "A"
// 		}
// 		return " "
// 	}
// 	template := s[0 : n+1]
// 	var truncateRe = regexp.MustCompile(`.(?=\W*\w*$)`)
// 	truncateRe.ReplaceAllStringFunc(template, tmpl) // 'Hello, world' -> 'HellAA AAAAA'
// 	var wwRe = regexp.MustCompile(`\w\w`)
// 	var whitespaceRe2 = regexp.MustCompile(`\s*\S+$`)
// 	if wwRe.MatchString(template[len(template)-2:]) {
// 		template = whitespaceRe2.ReplaceAllString(template, "")
// 	} else {
// 		template = strings.TrimRight(template, " \t\n")
// 	}
//
// 	if len(template+morestr) > L {
// 		return s
// 	}
// 	return s[0:len(template)] + morestr
// }
//
//     truncate: function(length, pruneStr) { //from underscore.string, author: github.com/rwz
//       var str = this.s;
//
//       length = ~~length;
//       pruneStr = pruneStr || '...';
//
//       if (str.length <= length) return new this.constructor(str);
//
//       var tmpl = function(c){ return c.toUpperCase() !== c.toLowerCase() ? 'A' : ' '; },
//         template = str.slice(0, length+1).replace(/.(?=\W*\w*$)/g, tmpl); // 'Hello, world' -> 'HellAA AAAAA'
//
//       if (template.slice(template.length-2).match(/\w\w/))
//         template = template.replace(/\s*\S+$/, '');
//       else
//         template = new S(template.slice(0, template.length-1)).trimRight().s;
//
//       return (template+pruneStr).length > str.length ? new S(str) : new S(str.slice(0, template.length)+pruneStr);
//     },

// Underscore returns converted camel cased string into a string delimited by underscores.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - string: the underscored string.
func Underscore(s string) string {
	if s == "" {
		return ""
	}
	u := strings.TrimSpace(s)

	u = underscoreRe.ReplaceAllString(u, "${1}_$2")
	u = dashSpaceRe.ReplaceAllString(u, "_")
	u = strings.ToLower(u)
	if IsUpper(s[0:1]) {
		return "_" + u
	}
	return u
}

// UnescapeHTML is an alias for html.UnescapeString.
//
// Parameters:
//   - s: the string to unescape.
//
// Returns:
//   - string: the unescaped string.
func UnescapeHTML(s string) string {
	if Verbose {
		fmt.Println("Use html.UnescapeString instead of UnescapeHTML")
	}
	return html.UnescapeString(s)
}

// WrapHTML wraps s within HTML tag having attributes attrs. Note,
// WrapHTML does not escape s value.
//
// Parameters:
//   - s: the string to wrap.
//   - tag: the HTML tag.
//   - attrs: the HTML attributes.
//
// Returns:
//   - string: the wrapped HTML string.
func WrapHTML(s string, tag string, attrs map[string]string) string {
	escapeHTMLAttributeQuotes := func(v string) string {
		v = strings.Replace(v, "<", "&lt;", -1)
		v = strings.Replace(v, "&", "&amp;", -1)
		v = strings.Replace(v, "\"", "&quot;", -1)
		return v
	}
	if tag == "" {
		tag = "div"
	}
	el := "<" + tag
	for name, val := range attrs {
		el += " " + name + "=\"" + escapeHTMLAttributeQuotes(val) + "\""
	}
	el += ">" + s + "</" + tag + ">"
	return el
}

// WrapHTMLF is the filter form of WrapHTML.
//
// Parameters:
//   - tag: the HTML tag.
//   - attrs: the HTML attributes.
//
// Returns:
//   - func(string) string: a function that wraps the string in HTML.
func WrapHTMLF(tag string, attrs map[string]string) func(string) string {
	return func(s string) string {
		return WrapHTML(s, tag, attrs)
	}
}
