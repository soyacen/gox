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
