//go:build !go1.21

package strconvx

import (
	"unsafe"
)

type (
	stringHeader struct {
		Data unsafe.Pointer
		Len  int
	}
	sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}
)

// StringToBytes returns an unsafe bytes slice reference of s.
// The caller must treat returned slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user.
func StringToBytes(s string) []byte {
	var b []byte
	src := (*stringHeader)(unsafe.Pointer(&s))
	dst := (*sliceHeader)(unsafe.Pointer(&b))
	dst.Data = src.Data
	dst.Len = src.Len
	dst.Cap = src.Len
	return b
}

// BytesToString returns an unsafe string reference of b.
// The caller must treat the input slice as immutable.
//
// WARNING: Use carefully. The returned result must not leak to the end user
// unless the input slice is provably immutable.
func BytesToString(b []byte) string {
	var s string
	src := (*sliceHeader)(unsafe.Pointer(&b))
	dst := (*stringHeader)(unsafe.Pointer(&s))
	dst.Data = src.Data
	dst.Len = src.Len
	return s
}
