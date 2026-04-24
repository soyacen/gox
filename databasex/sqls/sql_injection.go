package sqls

import (
	"regexp"
	"strings"

	"github.com/soyacen/gox/stringx"
)

// sqlSyntaxPattern is a regex that detects SQL syntax by matching two keywords in sequence.
// Reference: mybatis-plus-core SQL injection detection.
var sqlSyntaxPattern = regexp.MustCompile("(?i)(insert|delete|update|select|create|drop|truncate|grant|alter|deny|revoke|call|execute|exec|declare|show|rename|set).+(into|from|set|where|table|database|view|index|on|cursor|procedure|trigger|for|password|union|and|or)")

// sqlCommentPattern is a regex that detects SQL truncation using quotes, semicolons, or comments.
// Reference: mybatis-plus-core SQL injection detection.
var sqlCommentPattern = regexp.MustCompile("(?i)'.*(or|union|--|#|/*|;)")

// CheckSqlInjection checks whether the given value contains SQL injection patterns.
// When funcAllowed is false, any presence of '(' will be flagged as suspicious
// to prevent injection like `order by id,if(1=2,1,(sleep(100)))`.
//
// Parameters:
//   - value: the string to inspect.
//   - funcAllowed: if true, allows SQL function calls (parentheses); if false, rejects them.
//
// Returns:
//   - bool: true if the value is suspicious (potentially SQL injection), false otherwise.
func CheckSqlInjection(value string, funcAllowed bool) bool {
	if stringx.IsBlank(value) {
		return false
	}
	if funcAllowed {
		return sqlCommentPattern.MatchString(value) || sqlSyntaxPattern.MatchString(value)
	}
	// Reject any function calls (parentheses) when funcAllowed is false,
	// otherwise injection like `order by id,if(1=2,1,(sleep(100)));` cannot be detected.
	return strings.Contains(value, "(") || sqlCommentPattern.MatchString(value) || sqlSyntaxPattern.MatchString(value)
}
