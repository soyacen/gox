package errorx

import (
	"errors"
	"reflect"

	"github.com/soyacen/gox/reflectx"
)

// Equals checks whether two errors are equal.
// It first checks for nil equality, then uses errors.Is, reflection-based type comparison,
// and deep equality to determine if the errors represent the same value.
//
// Parameters:
//   - err: The error to compare
//   - target: The target error to compare against
//
// Returns:
//   - bool: true if the errors are equal, false otherwise
func Equals(err, target error) bool {
	if err == nil && target == nil {
		return true
	}
	if err != nil && target != nil {
		return equals(err, target)
	}
	return false
}

func equals(err error, target error) bool {
	if errors.Is(err, target) {
		return true
	}

	errVal := reflectx.IndirectValue(reflect.ValueOf(err))
	errType := errVal.Type()
	targetVal := reflectx.IndirectValue(reflect.ValueOf(target))
	targetType := targetVal.Type()
	if errType != targetType {
		return false
	}
	if errType.Comparable() && targetType.Comparable() && errVal.Equal(targetVal) {
		return true
	}

	if reflect.DeepEqual(errVal.Interface(), targetVal.Interface()) {
		return true
	}

	return false
}
