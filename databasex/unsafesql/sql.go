// Package unsafesql provides a set of utilities for building SQL query statements.
// It helps developers safely construct dynamic SQL queries and avoid errors from manual string concatenation.
package unsafesql

import (
	"errors"
	"fmt"
	"strings"
)

// Combiner is an interface that defines how to combine multiple SQL conditions together.
// For example, using AND or OR operators to join conditions.
type Combiner interface {
	// Combine joins multiple condition strings into one according to specific logic.
	//
	// Parameters:
	//   - conditions: A slice of condition strings to combine.
	//
	// Returns:
	//   - string: The combined condition string.
	Combine(conditions []string) string
}

// combinerFunc is a function type that implements the Combiner interface.
// It allows regular functions to be used as Combiner interface implementations.
type combinerFunc func(conditions []string) string

// Combine implements the Combiner interface for combinerFunc.
//
// Parameters:
//   - conditions: A slice of condition strings to combine.
//
// Returns:
//   - string: The combined condition string.
func (f combinerFunc) Combine(conditions []string) string {
	return f(conditions)
}

// And returns a Combiner that uses the AND operator to combine conditions.
// Returns an empty string when the conditions array is empty.
//
// Returns:
//   - Combiner: A combiner that joins conditions with AND.
func And() Combiner {
	return combinerFunc(func(conditions []string) string {
		if len(conditions) <= 0 {
			return ""
		}
		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND "))
	})
}

// Or returns a Combiner that uses the OR operator to combine conditions.
// Returns an empty string when the conditions array is empty.
//
// Returns:
//   - Combiner: A combiner that joins conditions with OR.
func Or() Combiner {
	return combinerFunc(func(conditions []string) string {
		if len(conditions) <= 0 {
			return ""
		}
		return fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
	})
}

// Where constructs a WHERE clause, combining conditions with the specified combiner.
// Returns an empty string if the conditions array is empty.
//
// Parameters:
//   - conditions: A slice of condition strings.
//   - op: A Combiner used to join the conditions.
//
// Returns:
//   - string: The constructed WHERE clause, or an empty string if no conditions.
//
// Example:
//
//	Where([]string{"age > 18", "status = 'active'"}, And())
//	// Returns: "WHERE age > 18 AND status = 'active'"
func Where(conditions []string, op Combiner) string {
	if len(conditions) <= 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", op.Combine(conditions))
}

// MustWhere constructs a WHERE clause, but panics if the conditions array is empty.
// Used in scenarios where query conditions are mandatory.
//
// Parameters:
//   - conditions: A slice of condition strings.
//   - op: A Combiner used to join the conditions.
//
// Returns:
//   - string: The constructed WHERE clause.
//
// Panics:
//   - If conditions is empty.
func MustWhere(conditions []string, op Combiner) string {
	if len(conditions) <= 0 {
		panic(errors.New("sqls: conditions is empty"))
	}
	return Where(conditions, op)
}

// Having constructs a HAVING clause for filtering conditions in aggregate queries.
// Returns an empty string if the conditions array is empty.
//
// Parameters:
//   - conditions: A slice of condition strings.
//   - op: A Combiner used to join the conditions.
//
// Returns:
//   - string: The constructed HAVING clause, or an empty string if no conditions.
func Having(conditions []string, op Combiner) string {
	if len(conditions) <= 0 {
		return ""
	}
	return fmt.Sprintf("HAVING %s", op.Combine(conditions))
}

// MustHaving constructs a HAVING clause, but panics if the conditions array is empty.
//
// Parameters:
//   - conditions: A slice of condition strings.
//   - op: A Combiner used to join the conditions.
//
// Returns:
//   - string: The constructed HAVING clause.
//
// Panics:
//   - If conditions is empty.
func MustHaving(conditions []string, op Combiner) string {
	if len(conditions) <= 0 {
		panic(errors.New("sqls: conditions is empty"))
	}
	return Having(conditions, op)
}

// Eq adds an equals (=) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	Eq(conditions, "name", "'john'")
//	// Appends "(name = 'john')" to the conditions array.
func Eq(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s = %s)", field, value))
}

// MustEq adds an equals (=) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustEq(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Eq(conditions, field, value)
}

// Ne adds a not equals (<>) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
func Ne(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s <> %s)", field, value))
}

// MustNe adds a not equals (<>) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustNe(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Ne(conditions, field, value)
}

// Gt adds a greater than (>) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
func Gt(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s > %s)", field, value))
}

// MustGt adds a greater than (>) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustGt(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Gt(conditions, field, value)
}

// Lt adds a less than (<) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
func Lt(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s < %s)", field, value))
}

// MustLt adds a less than (<) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustLt(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Lt(conditions, field, value)
}

// Ge adds a greater than or equal (>=) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
func Ge(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s >= %s)", field, value))
}

// MustGe adds a greater than or equal (>=) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustGe(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Ge(conditions, field, value)
}

// Le adds a less than or equal (<=) condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
func Le(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s <= %s)", field, value))
}

// MustLe adds a less than or equal (<=) condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The value to compare.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustLe(conditions []string, field, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Le(conditions, field, value)
}

// Between adds a BETWEEN condition to the conditions array.
// Used to check if a field value is within a specified range.
// Ignores the condition if field name or min/max values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - minValue: The minimum value.
//   - maxValue: The maximum value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	Between(conditions, "age", "18", "65")
//	// Appends "(age BETWEEN 18 AND 65)"
func Between(conditions []string, field, minValue, maxValue string) []string {
	if len(field) <= 0 || len(minValue) <= 0 || len(maxValue) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s BETWEEN %s AND %s)", field, minValue, maxValue))
}

// MustBetween adds a BETWEEN condition, but panics if field name or min/max values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - minValue: The minimum value.
//   - maxValue: The maximum value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field, minValue, or maxValue is empty.
func MustBetween(conditions []string, field, minValue, maxValue string) []string {
	if len(field) <= 0 || len(minValue) <= 0 || len(maxValue) <= 0 {
		panic(errors.New("sqls: field or min or max value is empty"))
	}
	return Between(conditions, field, minValue, maxValue)
}

// NotBetween adds a NOT BETWEEN condition to the conditions array.
// Ignores the condition if field name or min/max values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - minValue: The minimum value.
//   - maxValue: The maximum value.
//
// Returns:
//   - []string: The updated conditions slice.
func NotBetween(conditions []string, field, minValue, maxValue string) []string {
	if len(field) <= 0 || len(minValue) <= 0 || len(maxValue) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s NOT BETWEEN %s AND %s)", field, minValue, maxValue))
}

// MustNotBetween adds a NOT BETWEEN condition, but panics if field name or min/max values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - minValue: The minimum value.
//   - maxValue: The maximum value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field, minValue, or maxValue is empty.
func MustNotBetween(conditions []string, field, minValue, maxValue string) []string {
	if len(field) <= 0 || len(minValue) <= 0 || len(maxValue) <= 0 {
		panic(errors.New("sqls: field or min or max value is empty"))
	}
	return NotBetween(conditions, field, minValue, maxValue)
}

// Like adds a LIKE condition to the conditions array.
// Used for pattern matching queries.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The pattern value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	Like(conditions, "name", "%john%")
//	// Appends "(name LIKE '%john%')"
func Like(conditions []string, field string, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s LIKE %s)", field, value))
}

// MustLike adds a LIKE condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The pattern value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustLike(conditions []string, field string, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return Like(conditions, field, value)
}

// NotLike adds a NOT LIKE condition to the conditions array.
// Ignores the condition if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The pattern value.
//
// Returns:
//   - []string: The updated conditions slice.
func NotLike(conditions []string, field string, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s NOT LIKE %s)", field, value))
}

// MustNotLike adds a NOT LIKE condition, but panics if field name or value is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - value: The pattern value.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or value is empty.
func MustNotLike(conditions []string, field string, value string) []string {
	if len(field) <= 0 || len(value) <= 0 {
		panic(errors.New("sqls: field or value is empty"))
	}
	return NotLike(conditions, field, value)
}

// In adds an IN condition to the conditions array.
// Used to check if a field value exists in a specified list of values.
// Ignores the condition if field name or values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - values: A slice of values.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	In(conditions, "id", []string{"1", "2", "3"})
//	// Appends "(id IN (1, 2, 3))"
func In(conditions []string, field string, values []string) []string {
	if len(field) <= 0 || len(values) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s IN (%s))", field, strings.Join(values, ", ")))
}

// MustIn adds an IN condition, but panics if field name or values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - values: A slice of values.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or values is empty.
func MustIn(conditions []string, field string, values []string) []string {
	if len(field) <= 0 || len(values) <= 0 {
		panic(errors.New("sqls: field or values is empty"))
	}
	return In(conditions, field, values)
}

// NotIn adds a NOT IN condition to the conditions array.
// Ignores the condition if field name or values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - values: A slice of values.
//
// Returns:
//   - []string: The updated conditions slice.
func NotIn(conditions []string, field string, values []string) []string {
	if len(field) <= 0 || len(values) <= 0 {
		return conditions
	}
	return append(conditions, fmt.Sprintf("(%s NOT IN (%s))", field, strings.Join(values, ", ")))
}

// MustNotIn adds a NOT IN condition, but panics if field name or values are empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - values: A slice of values.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field or values is empty.
func MustNotIn(conditions []string, field string, values []string) []string {
	if len(field) <= 0 || len(values) <= 0 {
		panic(errors.New("sqls: field or values is empty"))
	}
	return NotIn(conditions, field, values)
}

// IsNull adds an IS NULL condition to the conditions array.
// Used to check if a field value is NULL.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	IsNull(conditions, "email")
//	// Appends "(email IS NULL)"
func IsNull(conditions []string, field string) []string {
	return append(conditions, fmt.Sprintf("(%s IS NULL)", field))
}

// MustIsNull adds an IS NULL condition, but panics if field name is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field is empty.
func MustIsNull(conditions []string, field string) []string {
	if len(field) <= 0 {
		panic(errors.New("sqls: field is empty"))
	}
	return IsNull(conditions, field)
}

// IsNotNull adds an IS NOT NULL condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//
// Returns:
//   - []string: The updated conditions slice.
func IsNotNull(conditions []string, field string) []string {
	return append(conditions, fmt.Sprintf("(%s IS NOT NULL)", field))
}

// MustIsNotNull adds an IS NOT NULL condition, but panics if field name is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If field is empty.
func MustIsNotNull(conditions []string, field string) []string {
	if len(field) <= 0 {
		panic(errors.New("sqls: field is empty"))
	}
	return IsNotNull(conditions, field)
}

// Exists adds an EXISTS subquery condition to the conditions array.
// Used to check if a subquery returns results.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Example:
//
//	Exists(conditions, "SELECT 1 FROM orders WHERE user_id = users.id")
//	// Appends "(EXISTS (SELECT 1 FROM orders WHERE user_id = users.id))"
func Exists(conditions []string, subquery string) []string {
	return append(conditions, fmt.Sprintf("(EXISTS (%s))", subquery))
}

// MustExists adds an EXISTS subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustExists(conditions []string, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return Exists(conditions, subquery)
}

// NotExists adds a NOT EXISTS subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func NotExists(conditions []string, subquery string) []string {
	return append(conditions, fmt.Sprintf("(NOT EXISTS (%s))", subquery))
}

// MustNotExists adds a NOT EXISTS subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustNotExists(conditions []string, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return NotExists(conditions, subquery)
}

// EqAll adds an = ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func EqAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s = ALL (%s))", field, subquery))
}

// MustEqAll adds an = ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustEqAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return EqAll(conditions, field, subquery)
}

// EqAny adds an = ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func EqAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s = ANY (%s))", field, subquery))
}

// MustEqAny adds an = ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustEqAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return EqAny(conditions, field, subquery)
}

// NeAll adds a <> ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func NeAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s <> ALL (%s))", field, subquery))
}

// MustNeAll adds a <> ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustNeAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return NeAll(conditions, field, subquery)
}

// NeAny adds a <> ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func NeAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s <> ANY (%s))", field, subquery))
}

// MustNeAny adds a <> ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustNeAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return NeAny(conditions, field, subquery)
}

// GtAll adds a > ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func GtAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s > ALL (%s))", field, subquery))
}

// MustGtAll adds a > ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustGtAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return GtAll(conditions, field, subquery)
}

// GtAny adds a > ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func GtAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s > ANY (%s))", field, subquery))
}

// MustGtAny adds a > ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustGtAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return GtAny(conditions, field, subquery)
}

// LtAll adds a < ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func LtAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s < ALL (%s))", field, subquery))
}

// MustLtAll adds a < ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustLtAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return LtAll(conditions, field, subquery)
}

// LtAny adds a < ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func LtAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s < ANY (%s))", field, subquery))
}

// MustLtAny adds a < ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustLtAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return LtAny(conditions, field, subquery)
}

// GeAll adds a >= ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func GeAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s >= ALL (%s))", field, subquery))
}

// MustGeAll adds a >= ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustGeAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return GeAll(conditions, field, subquery)
}

// GeAny adds a >= ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func GeAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s >= ANY (%s))", field, subquery))
}

// MustGeAny adds a >= ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustGeAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return GeAny(conditions, field, subquery)
}

// LeAll adds a <= ALL subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func LeAll(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s <= ALL (%s))", field, subquery))
}

// MustLeAll adds a <= ALL subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustLeAll(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return LeAll(conditions, field, subquery)
}

// LeAny adds a <= ANY subquery condition to the conditions array.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
func LeAny(conditions []string, field, subquery string) []string {
	return append(conditions, fmt.Sprintf("(%s <= ANY (%s))", field, subquery))
}

// MustLeAny adds a <= ANY subquery condition, but panics if subquery is empty.
//
// Parameters:
//   - conditions: The current slice of conditions.
//   - field: The field name.
//   - subquery: The subquery string.
//
// Returns:
//   - []string: The updated conditions slice.
//
// Panics:
//   - If subquery is empty.
func MustLeAny(conditions []string, field, subquery string) []string {
	if len(subquery) <= 0 {
		panic(errors.New("sqls: subquery is empty"))
	}
	return LeAny(conditions, field, subquery)
}

// Select constructs a SELECT clause.
// Defaults to SELECT * if no fields are specified.
//
// Parameters:
//   - fields: The field names to select.
//
// Returns:
//   - string: The constructed SELECT clause.
//
// Example:
//
//	Select("id", "name", "email")
//	// Returns: "SELECT id, name, email"
func Select(fields ...string) string {
	if len(fields) <= 0 {
		return "SELECT *"
	}
	return fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
}

// From constructs a FROM clause, specifying the query table name.
// Panics if the table name is empty.
//
// Parameters:
//   - table: The table name.
//
// Returns:
//   - string: The constructed FROM clause.
//
// Example:
//
//	From("users")
//	// Returns: "FROM users"
func From(table string) string {
	if len(table) == 0 {
		panic(errors.New("sqls: table is empty"))
	}
	return fmt.Sprintf("FROM %s", table)
}

// LeftJoin constructs a LEFT JOIN clause.
// Requires specifying the join table name and join conditions.
//
// Parameters:
//   - table: The table name to join.
//   - conditions: The join conditions.
//
// Returns:
//   - string: The constructed LEFT JOIN clause.
//
// Example:
//
//	LeftJoin("orders", "users.id = orders.user_id")
//	// Returns: "LEFT JOIN orders ON users.id = orders.user_id"
func LeftJoin(table string, conditions ...string) string {
	if len(table) <= 0 || len(conditions) <= 0 {
		panic(errors.New("sqls: conditions is empty"))
	}
	return fmt.Sprintf("LEFT JOIN %s ON %s", table, strings.Join(conditions, " AND "))
}

// RightJoin constructs a RIGHT JOIN clause.
// Requires specifying the join table name and join conditions.
//
// Parameters:
//   - table: The table name to join.
//   - conditions: The join conditions.
//
// Returns:
//   - string: The constructed RIGHT JOIN clause.
func RightJoin(table string, conditions ...string) string {
	if len(table) <= 0 || len(conditions) <= 0 {
		panic(errors.New("sqls: conditions is empty"))
	}
	return fmt.Sprintf("RIGHT JOIN %s ON %s", table, strings.Join(conditions, " AND "))
}

// GroupBy constructs a GROUP BY clause for grouping queries.
// Returns an empty string if the fields array is empty.
//
// Parameters:
//   - fields: A slice of field names to group by.
//
// Returns:
//   - string: The constructed GROUP BY clause, or an empty string if no fields.
//
// Example:
//
//	GroupBy([]string{"department", "status"})
//	// Returns: "GROUP BY department, status"
func GroupBy(fields []string) string {
	if len(fields) <= 0 {
		return ""
	}
	return fmt.Sprintf("GROUP BY %s", strings.Join(fields, ", "))
}

// MustGroupBy constructs a GROUP BY clause, but panics if the fields array is empty.
//
// Parameters:
//   - fields: A slice of field names to group by.
//
// Returns:
//   - string: The constructed GROUP BY clause.
//
// Panics:
//   - If fields is empty.
func MustGroupBy(fields []string) string {
	if len(fields) <= 0 {
		panic(errors.New("sqls: fields is empty"))
	}
	return GroupBy(fields)
}

// OrderBy constructs an ORDER BY clause for sorting.
// Returns an empty string if the fields array is empty.
//
// Parameters:
//   - fields: A slice of field names with optional sort direction.
//
// Returns:
//   - string: The constructed ORDER BY clause, or an empty string if no fields.
//
// Example:
//
//	OrderBy([]string{"created_at DESC", "name ASC"})
//	// Returns: "ORDER BY created_at DESC, name ASC"
func OrderBy(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(fields, ", "))
}

// MustOrderBy constructs an ORDER BY clause, but panics if the fields array is empty.
//
// Parameters:
//   - fields: A slice of field names with optional sort direction.
//
// Returns:
//   - string: The constructed ORDER BY clause.
//
// Panics:
//   - If fields is empty.
func MustOrderBy(fields []string) string {
	if len(fields) <= 0 {
		panic(errors.New("sqls: fields is empty"))
	}
	return OrderBy(fields)
}

// Limit constructs a LIMIT clause to limit the number of returned records.
// Returns an empty string if the value is less than 0.
//
// Parameters:
//   - n: The maximum number of records to return.
//
// Returns:
//   - string: The constructed LIMIT clause, or an empty string if n < 0.
//
// Example:
//
//	Limit(10)
//	// Returns: "LIMIT 10"
func Limit(n int) string {
	if n < 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", n)
}

// MustLimit constructs a LIMIT clause, but panics if the value is less than 0.
//
// Parameters:
//   - n: The maximum number of records to return.
//
// Returns:
//   - string: The constructed LIMIT clause.
//
// Panics:
//   - If n is less than 0.
func MustLimit(n int) string {
	if n < 0 {
		panic(errors.New("sqls: n is less than 0"))
	}
	return Limit(n)
}

// Offset constructs an OFFSET clause to specify the query offset.
// Returns an empty string if the value is less than 0.
//
// Parameters:
//   - n: The offset value.
//
// Returns:
//   - string: The constructed OFFSET clause, or an empty string if n < 0.
//
// Example:
//
//	Offset(20)
//	// Returns: "OFFSET 20"
func Offset(n int) string {
	if n < 0 {
		return ""
	}
	return fmt.Sprintf("OFFSET %d", n)
}

// MustOffset constructs an OFFSET clause, but panics if the value is less than 0.
//
// Parameters:
//   - n: The offset value.
//
// Returns:
//   - string: The constructed OFFSET clause.
//
// Panics:
//   - If n is less than 0.
func MustOffset(n int) string {
	if n < 0 {
		panic(errors.New("sqls: n is less than 0"))
	}
	return Offset(n)
}
