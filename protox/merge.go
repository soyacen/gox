package protox

import (
	"github.com/soyacen/gox/errorx"
	"google.golang.org/protobuf/types/known/structpb"
)

// MergeStruct combines multiple structpb.Struct values into a single struct
// It creates a new target struct and merges all source structs into it
func MergeStruct(values ...*structpb.Struct) *structpb.Struct {
	target := errorx.Ignore(structpb.NewStruct(map[string]any{}))
	for _, value := range values {
		mergeStruct(target, value)
	}
	return target
}

// mergeStruct merges fields from source struct into target struct
// For each field in source, it makes a deep copy and adds to target
func mergeStruct(target *structpb.Struct, source *structpb.Struct) {
	for key, value := range source.GetFields() {
		target.Fields[key] = copyValue(value)
	}
}

// mergeList merges values from source ListValue into target ListValue
// Appends each item from source to target after making a copy
func mergeList(target *structpb.ListValue, source *structpb.ListValue) {
	for _, item := range source.GetValues() {
		target.Values = append(target.Values, copyValue(item))
	}
}

// copyValue creates a deep copy of a protobuf Value
// Handles all types of Value including structs, lists and primitive types
// Returns NullValue for nil or unknown types
func copyValue(value *structpb.Value) *structpb.Value {
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
		subValue := errorx.Ignore(structpb.NewStruct(map[string]any{}))
		mergeStruct(subValue, v.StructValue)
		return structpb.NewStructValue(subValue)
	case *structpb.Value_ListValue:
		subList := errorx.Ignore(structpb.NewList([]any{}))
		mergeList(subList, v.ListValue)
		return structpb.NewListValue(subList)
	default:
		return structpb.NewNullValue()
	}
}
