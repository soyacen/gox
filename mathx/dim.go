package mathx

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// Dim returns the maximum of x-y or 0.
//
// Parameters:
//   - x: The first value
//   - y: The second value
//
// Returns:
//   - N: The maximum of x-y or 0
func Dim[N constraints.Integer | constraints.Float](x, y N) N {
	v := x - y
	if v <= 0 {
		return 0
	}
	return v
	//return N(math.Dim(float64(x), float64(y)))
}

// Max returns the maximum of two values.
//
// Parameters:
//   - x: The first value
//   - y: The second value
//
// Returns:
//   - N: The maximum value
func Max[N constraints.Integer | constraints.Float](x, y N) N {
	var r N
	max(x, y, &r)
	return r
}

// Min returns the minimum of two values.
//
// Parameters:
//   - x: The first value
//   - y: The second value
//
// Returns:
//   - N: The minimum value
func Min[N constraints.Integer | constraints.Float](x, y N) N {
	var r N
	min(x, y, &r)
	return r
}

func max(x, y any, r any) {
	switch x.(type) {
	case float32:
		ptr := r.(*float32)
		*ptr = float32(math.Max(float64(x.(float32)), float64(y.(float32))))
		return
	case float64:
		ptr := r.(*float64)
		*ptr = math.Max(x.(float64), y.(float64))
		return
	case int:
		ptr := r.(*int)
		*ptr = integerMax(x.(int), y.(int))
		return
	case int8:
		ptr := r.(*int8)
		*ptr = integerMax(x.(int8), y.(int8))
		return
	case int16:
		ptr := r.(*int16)
		*ptr = integerMax(x.(int16), y.(int16))
		return
	case int32:
		ptr := r.(*int32)
		*ptr = integerMax(x.(int32), y.(int32))
		return
	case int64:
		ptr := r.(*int64)
		*ptr = integerMax(x.(int64), y.(int64))
		return
	case uint:
		ptr := r.(*uint)
		*ptr = integerMax(x.(uint), y.(uint))
		return
	case uint8:
		ptr := r.(*uint8)
		*ptr = integerMax(x.(uint8), y.(uint8))
		return
	case uint16:
		ptr := r.(*uint16)
		*ptr = integerMax(x.(uint16), y.(uint16))
		return
	case uint32:
		ptr := r.(*uint32)
		*ptr = integerMax(x.(uint32), y.(uint32))
		return
	case uint64:
		ptr := r.(*uint64)
		*ptr = integerMax(x.(uint64), y.(uint64))
		return
	case uintptr:
		ptr := r.(*uintptr)
		*ptr = integerMax(x.(uintptr), y.(uintptr))
		return
	default:
		panic(fmt.Errorf("unknown type %T", x))
	}
}

func integerMax[N constraints.Integer](x, y N) N {
	if x > y {
		return x
	}
	return y
}

func min(x, y any, r any) {
	switch x.(type) {
	case float32:
		ptr := r.(*float32)
		*ptr = float32(math.Min(float64(x.(float32)), float64(y.(float32))))
		return
	case float64:
		ptr := r.(*float64)
		*ptr = math.Min(x.(float64), y.(float64))
		return
	case int:
		ptr := r.(*int)
		*ptr = integerMin(x.(int), y.(int))
		return
	case int8:
		ptr := r.(*int8)
		*ptr = integerMin(x.(int8), y.(int8))
		return
	case int16:
		ptr := r.(*int16)
		*ptr = integerMin(x.(int16), y.(int16))
		return
	case int32:
		ptr := r.(*int32)
		*ptr = integerMin(x.(int32), y.(int32))
		return
	case int64:
		ptr := r.(*int64)
		*ptr = integerMin(x.(int64), y.(int64))
		return
	case uint:
		ptr := r.(*uint)
		*ptr = integerMin(x.(uint), y.(uint))
		return
	case uint8:
		ptr := r.(*uint8)
		*ptr = integerMin(x.(uint8), y.(uint8))
		return
	case uint16:
		ptr := r.(*uint16)
		*ptr = integerMin(x.(uint16), y.(uint16))
		return
	case uint32:
		ptr := r.(*uint32)
		*ptr = integerMin(x.(uint32), y.(uint32))
		return
	case uint64:
		ptr := r.(*uint64)
		*ptr = integerMin(x.(uint64), y.(uint64))
		return
	case uintptr:
		ptr := r.(*uintptr)
		*ptr = integerMin(x.(uintptr), y.(uintptr))
		return
	default:
		panic(fmt.Errorf("unknown type %T", x))
	}
}

func integerMin[N constraints.Integer](x, y N) N {
	if x < y {
		return x
	}
	return y
}
