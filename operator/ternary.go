package operator

// Ternary is akin to a ternary operator in Go.
// Based on the boolean condition, it returns either exprIfTrue or exprIfFalse.
// It utilizes generics, making it applicable for expressions of any type.
//
// Parameters:
//   - condition: The boolean condition to evaluate
//   - exprIfTrue: The value to return if condition is true
//   - exprIfFalse: The value to return if condition is false
//
// Returns:
//   - T: Either exprIfTrue or exprIfFalse based on the condition
func Ternary[T any](condition bool, exprIfTrue, exprIfFalse T) T {
	if condition {
		return exprIfTrue
	}
	return exprIfFalse
}

// TernaryFunc is akin to a ternary operator in Go.
// It conditionally executes either exprIfTrue or exprIfFalse based on the value of condition,
// and returns the result of the executed function.
// It leverages Go's generics feature, making it applicable for any type.
//
// Parameters:
//   - condition: The boolean condition to evaluate
//   - exprIfTrue: The function to execute if condition is true
//   - exprIfFalse: The function to execute if condition is false
//
// Returns:
//   - T: The result of the executed function
func TernaryFunc[T any](condition bool, exprIfTrue, exprIfFalse func() T) T {
	if condition {
		return exprIfTrue()
	}
	return exprIfFalse()
}
