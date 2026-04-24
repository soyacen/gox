package reflectx

import (
	"fmt"
	"reflect"
)

// RangeArrayOrSlice converts an array or slice into a slice of any.
//
// Parameters:
//   - obj: The array or slice to convert.
//
// Returns:
//   - []any: A slice containing all elements of the array or slice.
//   - error: An error if obj is not an array or slice.
func RangeArrayOrSlice(obj any) ([]any, error) {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Array && objValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("reflectx: %T is not array or slice", obj)
	}
	res := make([]any, 0, objValue.Len())
	for i := 0; i < objValue.Len(); i++ {
		elem := objValue.Index(i)
		res = append(res, elem.Interface())
	}
	return res, nil
}

// RangeMap converts a map into two slices: one for keys and one for values.
//
// Parameters:
//   - obj: The map to convert.
//
// Returns:
//   - []any: A slice containing all keys of the map.
//   - []any: A slice containing all values of the map.
//   - error: An error if obj is not a map.
func RangeMap(obj any) ([]any, []any, error) {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Map {
		return nil, nil, fmt.Errorf("reflectx: %T is not map", obj)
	}
	resKeys := make([]any, 0, objValue.Len())
	resValues := make([]any, 0, objValue.Len())
	itr := objValue.MapRange()
	for itr.Next() {
		resKeys = append(resKeys, itr.Key().Interface())
		resValues = append(resValues, itr.Value().Interface())
	}
	return resKeys, resValues, nil
}
