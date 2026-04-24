package protox

import (
	"github.com/soyacen/gox/errorx"
	"google.golang.org/protobuf/types/known/structpb"
)

// MergeStructs combines multiple structpb.Struct values into a single struct.
// It creates a new target struct and merges all source structs into it.
//
// Parameters:
//   - values: the source structs to merge.
//
// Returns:
//   - *structpb.Struct: the merged struct containing all fields.
func MergeStructs(values ...*structpb.Struct) *structpb.Struct {
	target := errorx.Ignore(structpb.NewStruct(map[string]any{}))
	for _, value := range values {
		MergeStruct(target, value)
	}
	return target
}

// MergeStruct merges fields from source struct into target struct.
// For each field in source, it makes a deep copy and adds to target.
//
// Parameters:
//   - target: the struct to merge fields into.
//   - source: the struct to copy fields from.
func MergeStruct(target *structpb.Struct, source *structpb.Struct) {
	for key, value := range source.GetFields() {
		target.Fields[key] = CopyValue(value)
	}
}

// MergeList merges values from source ListValue into target ListValue.
// Appends each item from source to target after making a copy.
//
// Parameters:
//   - target: the list to append values to.
//   - source: the list to copy values from.
func MergeList(target *structpb.ListValue, source *structpb.ListValue) {
	for _, item := range source.GetValues() {
		target.Values = append(target.Values, CopyValue(item))
	}
}

// CopyValue creates a deep copy of a protobuf Value.
// Handles all types of Value including structs, lists and primitive types.
// Returns NullValue for nil or unknown types.
//
// Parameters:
//   - value: the protobuf Value to copy.
//
// Returns:
//   - *structpb.Value: a deep copy of the input value, or nil if input is nil.
func CopyValue(value *structpb.Value) *structpb.Value {
	if value == nil {
		return nil
	}
	switch v := value.GetKind().(type) {
	case *structpb.Value_NumberValue:
		return structpb.NewNumberValue(v.NumberValue)
	case *structpb.Value_StringValue:
		return structpb.NewStringValue(v.StringValue)
	case *structpb.Value_BoolValue:
		return structpb.NewBoolValue(v.BoolValue)
	case *structpb.Value_NullValue:
		return structpb.NewNullValue()
	case *structpb.Value_StructValue:
		newStruct := errorx.Ignore(structpb.NewStruct(map[string]any{}))
		MergeStruct(newStruct, v.StructValue)
		return structpb.NewStructValue(newStruct)
	case *structpb.Value_ListValue:
		newList := errorx.Ignore(structpb.NewList([]any{}))
		MergeList(newList, v.ListValue)
		return structpb.NewListValue(newList)
	default:
		return structpb.NewNullValue()
	}
}
