package reflectx

import (
	"reflect"
)

// New creates a new zero value of type T, initializing maps, slices, channels,
// and pointers recursively.
//
// Returns:
//   - T: A new zero value of type T with initialized complex types.
func New[T any]() T {
	var t T
	value := reflectNew(reflect.TypeOf(t)).Interface()
	return value.(T)
}

func reflectNew(typ reflect.Type) reflect.Value {
	if typ == nil {
		return reflect.Value{}
	}
	switch typ.Kind() {
	case reflect.Invalid:
		return reflect.Value{}
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Array,
		reflect.Interface,
		reflect.String,
		reflect.Struct,
		reflect.UnsafePointer:
		return reflect.Zero(typ)
	case reflect.Chan:
		return reflect.MakeChan(typ, 0)
	case reflect.Func:
		return reflect.MakeFunc(typ, func([]reflect.Value) []reflect.Value {
			results := make([]reflect.Value, typ.NumOut())
			for i := 0; i < typ.NumOut(); i++ {
				results[i] = reflect.Zero(typ.Out(i))
			}
			return results
		})
	case reflect.Map:
		return reflect.MakeMap(typ)
	case reflect.Slice:
		return reflect.MakeSlice(typ, 0, 0)
	case reflect.Ptr:
		value := reflect.New(typ.Elem())
		value.Elem().Set(reflectNew(typ.Elem()))
		return value
	default:
		return reflect.Value{}
	}
}
