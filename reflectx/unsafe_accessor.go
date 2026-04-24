package reflectx

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Accessor provides unsafe field access to a struct using field names or tag values.
// It uses unsafe.Pointer to access fields directly, including unexported fields.
type Accessor struct {
	fields  map[string]reflect.StructField
	address unsafe.Pointer
}

// Get retrieves the value of a field by name as any.
//
// Parameters:
//   - field: The name of the field to retrieve.
//
// Returns:
//   - any: The value of the field.
//   - error: An error if the field is not found.
func (accessor *Accessor) Get(field string) (any, error) {
	structField, ok := accessor.fields[field]
	if !ok {
		return nil, fmt.Errorf("reflectx: failed to get %s field", field)
	}
	return accessor.get(structField)
}

// Set sets the value of a field by name.
//
// Parameters:
//   - field: The name of the field to set.
//   - val: The new value to assign.
//
// Returns:
//   - error: An error if the field is not found.
func (accessor *Accessor) Set(field string, val any) error {
	structField, ok := accessor.fields[field]
	if !ok {
		return fmt.Errorf("reflectx: failed to get %s field", field)
	}
	accessor.set(structField, val)
	return nil
}

func (accessor *Accessor) fieldAddress(structField reflect.StructField) unsafe.Pointer {
	// field address = objAddress + fieldOffset
	return unsafe.Pointer(uintptr(accessor.address) + structField.Offset)
}

func (accessor *Accessor) get(structField reflect.StructField) (any, error) {
	return reflect.NewAt(structField.Type, accessor.fieldAddress(structField)).Elem().Interface(), nil
}

func (accessor *Accessor) set(structField reflect.StructField, val any) {
	reflect.NewAt(structField.Type, accessor.fieldAddress(structField)).Elem().Set(reflect.ValueOf(val))
}

func checkValue(objValue reflect.Value) (reflect.Value, error) {
	if !objValue.IsValid() {
		return reflect.Value{}, fmt.Errorf("reflectx: invalid reflect.Value")
	}
	if objValue.Kind() != reflect.Pointer {
		return reflect.Value{}, fmt.Errorf("reflectx: %T is not a pointer", objValue.Interface())
	}
	objStructValue := IndirectValue(objValue)

	// Ensure the value points to a struct or struct field.
	if objStructValue.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("reflectx: %T does not point to a struct", objValue.Interface())
	}
	return objStructValue, nil
}

func extractFields(objStructValue reflect.Value) map[string]reflect.StructField {
	objStructType := objStructValue.Type()
	numField := objStructType.NumField()
	fields := make(map[string]reflect.StructField, numField)
	for i := 0; i < numField; i++ {
		structField := objStructType.Field(i)
		fields[structField.Name] = structField
	}
	return fields
}

func extractFieldsByTag(objStructValue reflect.Value, key string) (map[string]reflect.StructField, error) {
	objStructType := objStructValue.Type()
	numField := objStructType.NumField()
	fields := make(map[string]reflect.StructField, numField)
	for i := 0; i < numField; i++ {
		structField := objStructType.Field(i)
		value, ok := structField.Tag.Lookup(key)
		if !ok {
			continue
		}
		if _, ok := fields[value]; ok {
			return nil, fmt.Errorf("reflectx: %s tag value %s is duplicated", key, value)
		}
		fields[value] = structField
	}
	return fields, nil
}

// FieldAccessorOf creates an Accessor for a struct using field names as keys.
//
// Parameters:
//   - objValue: A pointer reflect.Value to the struct.
//
// Returns:
//   - *Accessor: An Accessor that can read and write fields by name.
//   - error: An error if objValue is invalid, not a pointer, or not a struct.
func FieldAccessorOf(objValue reflect.Value) (*Accessor, error) {
	objStructValue, err := checkValue(objValue)
	if err != nil {
		return nil, err
	}
	fields := extractFields(objStructValue)
	address := objStructValue.Addr().UnsafePointer()
	return &Accessor{fields: fields, address: address}, nil
}

// TagAccessorOf creates an Accessor for a struct using tag values as keys.
//
// Parameters:
//   - objValue: A pointer reflect.Value to the struct.
//   - key: The tag key to use for field lookup.
//
// Returns:
//   - *Accessor: An Accessor that can read and write fields by tag value.
//   - error: An error if objValue is invalid, not a pointer, not a struct, or if tag values are duplicated.
func TagAccessorOf(objValue reflect.Value, key string) (*Accessor, error) {
	objStructValue, err := checkValue(objValue)
	if err != nil {
		return nil, err
	}

	fields, err := extractFieldsByTag(objStructValue, key)
	if err != nil {
		return nil, err
	}
	address := objStructValue.Addr().UnsafePointer()
	return &Accessor{fields: fields, address: address}, nil
}

// GetByAccessor retrieves the typed value of a field from an Accessor.
//
// Parameters:
//   - accessor: The Accessor to read from.
//   - field: The name of the field to retrieve.
//
// Returns:
//   - T: The typed value of the field.
//   - error: An error if the field is not found.
func GetByAccessor[T any](accessor *Accessor, field string) (T, error) {
	var res T
	structField, ok := accessor.fields[field]
	if !ok {
		return res, fmt.Errorf("reflectx: failed to get %s field", field)
	}
	res = *(*T)(accessor.fieldAddress(structField))
	return res, nil
}

// SetByAccessor sets the typed value of a field on an Accessor.
//
// Parameters:
//   - accessor: The Accessor to write to.
//   - field: The name of the field to set.
//   - val: The new typed value to assign.
//
// Returns:
//   - error: An error if the field is not found.
func SetByAccessor[T any](accessor *Accessor, field string, val T) error {
	structField, ok := accessor.fields[field]
	if !ok {
		return fmt.Errorf("reflectx: failed to get %s field", field)
	}
	*(*T)(accessor.fieldAddress(structField)) = val
	return nil
}
