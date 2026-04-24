package contextx

import (
	"context"
	"reflect"
)

// ContextType is the reflect.Type of context.Context.
//
// It represents the interface type for context values used to carry deadlines,
// cancellation signals, and other request-scoped values across API boundaries.
var ContextType = reflect.TypeOf((*context.Context)(nil)).Elem()
