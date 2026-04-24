package protox

import "google.golang.org/protobuf/types/known/wrapperspb"

// WrapBoolSlice wraps a slice of bool values into a slice of BoolValue wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of bool values.
//
// Returns:
//   - []*wrapperspb.BoolValue: the wrapped slice.
func WrapBoolSlice(s []bool) []*wrapperspb.BoolValue {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.BoolValue, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Bool(v))
	}
	return r
}

// UnwrapBoolSlice unwraps a slice of BoolValue wrappers into a slice of bool values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of BoolValue wrappers.
//
// Returns:
//   - []bool: the unwrapped slice.
func UnwrapBoolSlice(s []*wrapperspb.BoolValue) []bool {
	if s == nil {
		return nil
	}
	r := make([]bool, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapInt32Slice wraps a slice of int32 values into a slice of Int32Value wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of int32 values.
//
// Returns:
//   - []*wrapperspb.Int32Value: the wrapped slice.
func WrapInt32Slice(s []int32) []*wrapperspb.Int32Value {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.Int32Value, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Int32(v))
	}
	return r
}

// UnwrapInt32Slice unwraps a slice of Int32Value wrappers into a slice of int32 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of Int32Value wrappers.
//
// Returns:
//   - []int32: the unwrapped slice.
func UnwrapInt32Slice(s []*wrapperspb.Int32Value) []int32 {
	if s == nil {
		return nil
	}
	r := make([]int32, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapInt64Slice wraps a slice of int64 values into a slice of Int64Value wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of int64 values.
//
// Returns:
//   - []*wrapperspb.Int64Value: the wrapped slice.
func WrapInt64Slice(s []int64) []*wrapperspb.Int64Value {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.Int64Value, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Int64(v))
	}
	return r
}

// UnwrapInt64Slice unwraps a slice of Int64Value wrappers into a slice of int64 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of Int64Value wrappers.
//
// Returns:
//   - []int64: the unwrapped slice.
func UnwrapInt64Slice(s []*wrapperspb.Int64Value) []int64 {
	if s == nil {
		return nil
	}
	r := make([]int64, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapFloat32Slice wraps a slice of float32 values into a slice of FloatValue wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of float32 values.
//
// Returns:
//   - []*wrapperspb.FloatValue: the wrapped slice.
func WrapFloat32Slice(s []float32) []*wrapperspb.FloatValue {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.FloatValue, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Float(v))
	}
	return r
}

// UnwrapFloat32Slice unwraps a slice of FloatValue wrappers into a slice of float32 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of FloatValue wrappers.
//
// Returns:
//   - []float32: the unwrapped slice.
func UnwrapFloat32Slice(s []*wrapperspb.FloatValue) []float32 {
	if s == nil {
		return nil
	}
	r := make([]float32, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapFloat64Slice wraps a slice of float64 values into a slice of DoubleValue wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of float64 values.
//
// Returns:
//   - []*wrapperspb.DoubleValue: the wrapped slice.
func WrapFloat64Slice(s []float64) []*wrapperspb.DoubleValue {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.DoubleValue, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Double(v))
	}
	return r
}

// UnwrapFloat64Slice unwraps a slice of DoubleValue wrappers into a slice of float64 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of DoubleValue wrappers.
//
// Returns:
//   - []float64: the unwrapped slice.
func UnwrapFloat64Slice(s []*wrapperspb.DoubleValue) []float64 {
	if s == nil {
		return nil
	}
	r := make([]float64, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapUint32Slice wraps a slice of uint32 values into a slice of UInt32Value wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of uint32 values.
//
// Returns:
//   - []*wrapperspb.UInt32Value: the wrapped slice.
func WrapUint32Slice(s []uint32) []*wrapperspb.UInt32Value {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.UInt32Value, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.UInt32(v))
	}
	return r
}

// UnwrapUint32Slice unwraps a slice of UInt32Value wrappers into a slice of uint32 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of UInt32Value wrappers.
//
// Returns:
//   - []uint32: the unwrapped slice.
func UnwrapUint32Slice(s []*wrapperspb.UInt32Value) []uint32 {
	if s == nil {
		return nil
	}
	r := make([]uint32, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapUint64Slice wraps a slice of uint64 values into a slice of UInt64Value wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of uint64 values.
//
// Returns:
//   - []*wrapperspb.UInt64Value: the wrapped slice.
func WrapUint64Slice(s []uint64) []*wrapperspb.UInt64Value {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.UInt64Value, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.UInt64(v))
	}
	return r
}

// UnwrapUint64Slice unwraps a slice of UInt64Value wrappers into a slice of uint64 values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of UInt64Value wrappers.
//
// Returns:
//   - []uint64: the unwrapped slice.
func UnwrapUint64Slice(s []*wrapperspb.UInt64Value) []uint64 {
	if s == nil {
		return nil
	}
	r := make([]uint64, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapStringSlice wraps a slice of string values into a slice of StringValue wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of string values.
//
// Returns:
//   - []*wrapperspb.StringValue: the wrapped slice.
func WrapStringSlice(s []string) []*wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.StringValue, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.String(v))
	}
	return r
}

// UnwrapStringSlice unwraps a slice of StringValue wrappers into a slice of string values.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of StringValue wrappers.
//
// Returns:
//   - []string: the unwrapped slice.
func UnwrapStringSlice(s []*wrapperspb.StringValue) []string {
	if s == nil {
		return nil
	}
	r := make([]string, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}

// WrapBytesSlice wraps a slice of byte slices into a slice of BytesValue wrappers.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of byte slices.
//
// Returns:
//   - []*wrapperspb.BytesValue: the wrapped slice.
func WrapBytesSlice(s [][]byte) []*wrapperspb.BytesValue {
	if s == nil {
		return nil
	}
	r := make([]*wrapperspb.BytesValue, 0, len(s))
	for _, v := range s {
		r = append(r, wrapperspb.Bytes(v))
	}
	return r
}

// UnwrapBytesSlice unwraps a slice of BytesValue wrappers into a slice of byte slices.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of BytesValue wrappers.
//
// Returns:
//   - [][]byte: the unwrapped slice.
func UnwrapBytesSlice(s []*wrapperspb.BytesValue) [][]byte {
	if s == nil {
		return nil
	}
	r := make([][]byte, 0, len(s))
	for _, v := range s {
		r = append(r, v.GetValue())
	}
	return r
}
