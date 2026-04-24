//go:build go1.21

package strconvx

import (
	"unsafe"
)

// StringToBytes converts a string to a byte slice without a memory allocation.
//
// Parameters:
//   - s: the string to convert.
//
// Returns:
//   - []byte: the byte slice referencing the string's underlying data.
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts a byte slice to a string without a memory allocation.
//
// Parameters:
//   - b: the byte slice to convert.
//
// Returns:
//   - string: the string referencing the slice's underlying data.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
