package reflectx

import (
	"reflect"
)

// BaseKind returns the base reflect.Kind of a value after dereferencing pointers.
//
// Parameters:
//   - a: The value to determine the kind of.
//
// Returns:
//   - reflect.Kind: The base kind of the value.
func BaseKind(a any) reflect.Kind {
	a = Indirect(a)
	return reflect.TypeOf(a).Kind()
}
