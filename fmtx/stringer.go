package fmtx

import (
	"fmt"
	"reflect"
)

// StringerType is the reflect.Type of fmt.Stringer.
//
// It represents the interface type for types that implement the String() method.
var StringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
