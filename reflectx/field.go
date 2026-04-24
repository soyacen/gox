package reflectx

import (
	"fmt"
	"reflect"
	"unsafe"
)

// FindFieldByTag searches through the fields of a given struct value for a field
// that has a specific tag and whose tag value satisfies a given condition.
// It returns the reflect.Value of the matched field along with a boolean indicating
// whether a match was found.
//
// Parameters:
//   - objValue: The reflect.Value representing the struct to search.
//   - tagKey: The key of the tag to look for.
//   - match: A function that takes the tag value as a string and returns true if it satisfies the condition.
//
// Returns:
//   - reflect.Value: The reflect.Value of the matched field.
//   - bool: A boolean indicating if a match was found.
func FindFieldByTag(objValue reflect.Value, tagKey string, match func(tagVal string) bool) (reflect.Value, bool) {
	// Indirect the value to get the underlying value.
	structValue := IndirectValue(objValue)

	// Check if the value is a struct.
	if structValue.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	// Directly access the type once instead of on each iteration.
	structType := structValue.Type()

	// Iterate over all fields in the given struct.
	for i := 0; i < structType.NumField(); i++ {
		// Get the current field and its corresponding struct field type.
		field := structValue.Field(i)
		structField := structType.Field(i)

		// Check if the field has the specified tag with the given value.
		if tagVal, ok := structField.Tag.Lookup(tagKey); ok && match(tagVal) {
			// Return the field value if the tag matches.
			return field, true
		}
	}

	// Return zero Value if no matching field is found.
	return reflect.Value{}, false
}

// GetField retrieves the value of a named field from a struct.
//
// Parameters:
//   - objValue: The reflect.Value representing the struct.
//   - field: The name of the field to retrieve.
//
// Returns:
//   - any: The value of the field.
//   - error: An error if objValue is not a struct or the field is not found.
func GetField(objValue reflect.Value, field string) (any, error) {
	structValue := IndirectValue(objValue)
	for structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("reflectx: %T is not struct", structValue.Interface())
	}
	fieldVal := structValue.FieldByName(field)
	if fieldVal.IsZero() {
		return nil, fmt.Errorf("reflectx: field %s not found", field)
	}
	return fieldVal.Interface(), nil
}

// GetUnexportedField retrieves the value of an unexported field from a struct using reflection,
// supporting nested field access.
//
// Parameters:
//   - objValue: The reflect.Value of the object containing the target field.
//   - fields: A sequence of field names to access, supporting multi-level nested access.
//
// Returns:
//   - reflect.Value: The reflect.Value of the target field. Returns an invalid reflect.Value if the field is not found.
func GetUnexportedField(objValue reflect.Value, fields ...string) reflect.Value {
	// 遍历字段路径，逐级深入获取目标字段
	for _, field := range fields {
		v := IndirectValue(objValue)
		field := v.FieldByName(field)
		if !field.IsValid() {
			return reflect.Value{}
		}
		objValue = field
	}
	// 通过不安全指针创建可访问的字段值副本
	return reflect.NewAt(objValue.Type(), unsafe.Pointer(objValue.UnsafeAddr())).Elem()
}

// SetField sets the value of a named field on a struct.
//
// Parameters:
//   - objValue: A pointer reflect.Value to the struct.
//   - field: The name of the field to set.
//   - newValue: The new value to assign to the field.
//
// Returns:
//   - error: An error if objValue is not a pointer, not a struct, the field is not found, or the field cannot be set.
func SetField(objValue reflect.Value, field string, newValue any) error {
	if objValue.Kind() != reflect.Pointer {
		return fmt.Errorf("reflectx: %T is not pointer", objValue.Interface())
	}
	structValue := IndirectValue(objValue)
	for structValue.Kind() != reflect.Struct {
		return fmt.Errorf("reflectx: %T is not struct", structValue.Interface())
	}
	fieldVal := structValue.FieldByName(field)
	if fieldVal.IsZero() {
		return fmt.Errorf("reflectx: field %s not found", field)
	}
	if !fieldVal.CanSet() {
		return fmt.Errorf("reflectx: cannot set field %s", field)
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil
}

// RangeFields iterates over all fields of a struct and returns a map of field names to values.
// Exported fields return their actual values; unexported fields return zero values.
//
// Parameters:
//   - objValue: The reflect.Value representing the struct.
//
// Returns:
//   - map[string]any: A map where keys are field names and values are field values.
//   - error: An error if objValue is zero or not a struct.
func RangeFields(objValue reflect.Value) (map[string]any, error) {
	if objValue.IsZero() {
		return nil, fmt.Errorf("reflectx: unsupport zero value")
	}
	objValue = IndirectValue(objValue)
	objType := IndirectType(objValue.Type())
	if objType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("reflectx: %s is not struct", objType.Kind())
	}

	numField := objType.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)
		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return res, nil
}
