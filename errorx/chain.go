package errorx

// Break is a higher-order function that processes errors and determines whether to continue executing another function.
// If the pre parameter is not nil, it returns a zero value and the pre error, thus interrupting subsequent execution;
// If pre is nil, it calls and returns the result of f().
//
// Parameters:
//   - pre: The previous error to check
//
// Returns:
//   - func(func() (T, error)) (T, error): A function that takes a function and returns its result or an error
func Break[T any](pre error) func(f func() (T, error)) (T, error) {
	return func(f func() (T, error)) (T, error) {
		if pre != nil {
			var v T       // Declare a zero value of type T
			return v, pre // If pre is not nil, return the zero value of T and the error
		}
		return f() // If pre is nil, execute the wrapped function and return its result
	}
}

// Continue allows checking the pre error before executing f(), and if f() returns an error, it merges both errors and returns them.
// This is useful for continuing execution even when errors occur, collecting all errors along the way.
//
// Parameters:
//   - pre: The previous error to check
//
// Returns:
//   - func(func() (T, error)) (T, error): A function that executes f() and combines any errors with pre
func Continue[T any](pre error) func(f func() (T, error)) (T, error) {
	return func(f func() (T, error)) (T, error) {
		v, err := f()
		if err != nil {
			var v T                  // Declare a zero value of type T
			return v, Join(pre, err) // If there's an error from f, combine the previous error pre with the new error using Join, and return the zero value of T

		}
		return v, pre // If f executes without error, return the value and the original error pre	}
	}
}
