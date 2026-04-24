package reflectx

import (
	"reflect"
)

// Indirect returns the value after dereferencing as many times
// as necessary to reach the base type (or nil).
//
// Parameters:
//   - a: The value to dereference.
//
// Returns:
//   - any: The base value after dereferencing, or nil if the input was nil.
func Indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	return IndirectValue(reflect.ValueOf(a)).Interface()
}

// IndirectValue returns the reflect.Value after dereferencing as many times
// as necessary to reach the base type (or nil).
//
// Parameters:
//   - v: The reflect.Value to dereference.
//
// Returns:
//   - reflect.Value: The base reflect.Value after dereferencing.
func IndirectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

// IndirectType returns the reflect.Type after dereferencing as many times
// as necessary to reach the base type.
//
// Parameters:
//   - t: The reflect.Type to dereference.
//
// Returns:
//   - reflect.Type: The base reflect.Type after dereferencing.
func IndirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// IndirectOrImplements returns the reflect.Value after dereferencing until it
// implements one of the given interface types or reaches the base type.
//
// Parameters:
//   - v: The reflect.Value to dereference.
//   - ifaceTypes: A variadic list of interface types to check against.
//
// Returns:
//   - reflect.Value: The dereferenced reflect.Value that implements an interface or the base value.
func IndirectOrImplements(v reflect.Value, ifaceTypes ...reflect.Type) reflect.Value {
	for {
		for _, ifaceType := range ifaceTypes {
			if v.Type().Implements(ifaceType) {
				return v
			}
		}
		if v.Kind() == reflect.Pointer && !v.IsNil() {
			v = v.Elem()
			continue
		}
		return v
	}
}
