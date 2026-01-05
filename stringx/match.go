package stringx

import "regexp"

func MatchAll(str string, pattern string) []string {
	return regexp.MustCompile(pattern).FindAllString(str, -1)
}

// Match returns true if patterns matches the string
func Match(s, pattern string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(s)
}
