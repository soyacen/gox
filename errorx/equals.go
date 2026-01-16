package errorx

import (
	"errors"
	"reflect"

	"github.com/soyacen/gox/reflectx"
)

// Equals 判断两个错误是否相等
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
