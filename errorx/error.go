package errorx

// Must returns the value if the error is nil, otherwise it panics with the error message.
// This is useful for operations that should never fail in normal circumstances.
//
// Parameters:
//   - v: The value to return if no error
//   - err: The error to check
//
// Returns:
//   - T: The value if no error occurred
func Must[T any](v T, err error) T {
	if err != nil {
		panic("must: " + err.Error())
	}
	return v
}

// Ignore returns the value and ignores the error.
// Use this when you intentionally want to discard an error.
//
// Parameters:
//   - v: The value to return
//   - _: The error to ignore
//
// Returns:
//   - T: The value
func Ignore[T any](v T, _ error) T {
	return v
}

// Concern returns the error and discards the value.
// This is useful when you only care about the error.
//
// Parameters:
//   - _: The value to discard
//   - err: The error to return
//
// Returns:
//   - error: The error
func Concern[T any](_ T, err error) error {
	return err
}

// Quiet silently consumes the error without doing anything.
// Use this when you intentionally want to suppress an error.
//
// Parameters:
//   - _: The error to suppress
func Quiet(_ error) {}

// Stringfy converts an error to its string representation.
// Returns an empty string if the error is nil.
//
// Parameters:
//   - err: The error to convert
//
// Returns:
//   - string: The error message or empty string
func Stringfy(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
