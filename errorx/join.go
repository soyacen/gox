package errorx

import "errors"

// Join combines multiple errors into a single error.
// It unwraps the first error if it is a multi-error and appends the remaining errors.
//
// Parameters:
//   - err: The primary error to join
//   - errs: Additional errors to join
//
// Returns:
//   - error: A single error containing all provided errors
func Join(err error, errs ...error) error {
	return errors.Join(append(UnwrapMultiErr(err), errs...)...)
}

// UnwrapMultiErr unwraps a multi-error into its constituent errors.
// It supports standard errors.Join, hashicorp/go-multierror, and uber-go/multierr.
// If the error is not a multi-error, it returns a slice containing the error itself.
//
// Parameters:
//   - err: The error to unwrap
//
// Returns:
//   - []error: The constituent errors, or nil if err is nil
func UnwrapMultiErr(err error) []error {
	if err == nil {
		return nil
	}
	switch multiErr := err.(type) {
	case interface{ Unwrap() []error }:
		// standard join error
		return multiErr.Unwrap()
	case interface{ WrappedErrors() []error }:
		// https://github.com/hashicorp/go-multierror/blob/main/multierror.go
		return multiErr.WrappedErrors()
	case interface{ Errors() []error }:
		// https://github.com/uber-go/multierr/blob/master/error.go
		return multiErr.Errors()
	default:
		return []error{err}
	}
}
