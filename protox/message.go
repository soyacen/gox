package protox

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// MessageToStruct converts a protobuf message to a structpb.Struct.
// It marshals the message to JSON using proto names, then unmarshals into
// a map and creates a struct from it.
//
// Parameters:
//   - msg: the protobuf message to convert.
//
// Returns:
//   - *structpb.Struct: the resulting struct representation.
//   - error: an error if marshaling or unmarshaling fails.
func MessageToStruct(msg proto.Message) (*structpb.Struct, error) {
	data, err := protojson.MarshalOptions{
		UseProtoNames: true,
		AllowPartial:  true,
	}.Marshal(msg)
	if err != nil {
		return nil, err
	}
	m := map[string]any{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}

// StructToMessage converts a structpb.Struct to a protobuf message.
// It marshals the struct to JSON, and unmarshals into the provided proto
// message.
//
// Parameters:
//   - s: the structpb.Struct to convert.
//   - msg: the target protobuf message to populate.
//
// Returns:
//   - error: an error if marshaling or unmarshaling fails.
func StructToMessage(s *structpb.Struct, msg proto.Message) error {
	if s == nil || msg == nil {
		return nil
	}

	data, err := protojson.MarshalOptions{
		UseProtoNames: true,
		AllowPartial:  true,
	}.Marshal(s)
	if err != nil {
		return err
	}
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(data, msg)
}
