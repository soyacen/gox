package protox

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// MessageToStruct
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
