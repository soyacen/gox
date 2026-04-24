package errorx

import "reflect"

// ErrorType is the reflect.Type of the built-in error interface.
//
// It represents the interface type for all error values in Go.
var ErrorType = reflect.TypeOf((*error)(nil)).Elem()
