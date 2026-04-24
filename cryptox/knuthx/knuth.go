package knuthx

import (
	"unsafe"
)

// Knuth implements Knuth's multiplicative method for hash index calculation.
// It optimizes modulo operations to reduce bias.
//
// The algorithm maps a hash value x to an index in the range [0, m-1] using
// a carefully chosen multiplier constant:
//   1. Select a multiplier a based on architecture (32-bit or 64-bit).
//   2. Compute x * a.
//   3. Right-shift the result by (w - p) bits, where w is the word size
//      and p is the number of bits needed for the target range.
//   4. Modulo by the target range size m: y = ((x * a) >> (w - p)) % m.
//
// Reference: https://en.wikipedia.org/wiki/Knuth_multiplicative_method
//
// Parameters:
//   - hashSum: the hash value to map
//
// Returns:
//   - uint: the computed index
const (
	knuthMultiplier32 uint64 = 2654435761
	knuthMultiplier64 uint64 = 6364136223846793005
)

func Knuth(hashSum uint) uint {
	var multiplier uint64
	switch unsafe.Sizeof(uintptr(0)) {
	case 4:
		multiplier = knuthMultiplier32
	case 8:
		multiplier = knuthMultiplier64
	default:
		panic("hashx: unsupported architecture")
	}
	return uint((multiplier * uint64(hashSum)) >> 32)
}
