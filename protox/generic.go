package protox

import "google.golang.org/protobuf/proto"

// Clone creates a deep copy of a protobuf message.
//
// Parameters:
//   - m: the message to clone.
//
// Returns:
//   - M: a deep copy of the input message.
func Clone[M proto.Message](m M) M {
	return proto.Clone(m).(M)
}

// CloneSlice creates a deep copy of a slice of protobuf messages.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of messages to clone.
//
// Returns:
//   - S: a new slice containing deep copies of the input messages.
func CloneSlice[S ~[]M, M proto.Message](s S) S {
	var zero S
	if s == nil {
		return zero
	}
	r := make(S, 0, len(s))
	for _, m := range s {
		r = append(r, Clone(m))
	}
	return r
}

// MessageSlice converts a typed slice of protobuf messages to a slice of
// the proto.Message interface.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the typed slice of messages.
//
// Returns:
//   - []proto.Message: the converted slice of interface values.
func MessageSlice[S []E, E proto.Message](s S) []proto.Message {
	if s == nil {
		return nil
	}
	r := make([]proto.Message, 0, len(s))
	for _, e := range s {
		r = append(r, e)
	}
	return r
}

// ProtoSlice converts a slice of proto.Message interface values to a typed
// slice of a specific message type.
// Elements that do not match the target type are skipped.
// If the input slice is nil, it returns nil.
//
// Parameters:
//   - s: the slice of proto.Message interface values.
//
// Returns:
//   - S: the typed slice containing only elements of type E.
func ProtoSlice[S []E, E proto.Message](s []proto.Message) S {
	if s == nil {
		return nil
	}
	r := make(S, 0, len(s))
	for _, e := range s {
		p, ok := e.(E)
		if !ok {
			continue
		}
		r = append(r, p)
	}
	return r
}
