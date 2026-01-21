// Package randx provides utilities for generating random numbers and strings.
package randx

import (
	"crypto/rand"         // For generating cryptographically secure random bytes.
	"encoding/binary"     // To convert byte slices to numeric types.
	randv2 "math/rand/v2" // The version 2 of Go's math/rand package, offering better PRNGs.
	"sync"                // For synchronization primitives like sync.Pool.

	"github.com/soyacen/goconc/poolx"
	"github.com/soyacen/gox/errorx"
	// A custom package for managing byte buffers.
)

// stringPool is a pool of byte buffers used to minimize allocations during string generation.
var stringPool = errorx.Ignore(poolx.NewBucketBufferPool(16, 16*1024))

// Constants define common character sets for generating random strings.
const (
	// Lowercase alphabetic characters.
	Lowercase = "abcdefghijklmnopqrstuvwxyz"

	// Uppercase alphabetic characters.
	Uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Decimal numeric digits.
	Numeric = "0123456789"

	// Combination of lowercase, uppercase letters and digits.
	Alphanumeric = Lowercase + Uppercase + Numeric

	// Hexadecimal digits including lowercase a-f.
	Hex = Numeric + "abcdef"

	// Standard Base64 encoding character set.
	Base64 = Uppercase + Lowercase + Numeric + "+/"

	// URL-safe Base64 encoding character set (replaces '+' with '-', '/' with '_').
	URLSafeBase64 = Uppercase + Lowercase + Numeric + "-_"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// chacha8Pool manages a pool of ChaCha8 pseudo-random number generators.
var chacha8Pool = &sync.Pool{
	New: func() any {
		r, err := NewChaCha8() // Create a new ChaCha8 generator if none available.
		if err != nil {
			panic(err) // Panic on failure to initialize due to crypto/rand issues.
		}
		return r
	},
}

// pcgPool manages a pool of PCG pseudo-random number generators.
var pcgPool = &sync.Pool{
	New: func() any {
		r, err := NewPCG() // Create a new PCG generator if none available.
		if err != nil {
			panic(err) // Panic on failure to initialize due to crypto/rand issues.
		}
		return r
	},
}

// GetChaCha8 retrieves a ChaCha8 PRNG from the pool or creates one if needed.
func GetChaCha8() *randv2.Rand {
	return chacha8Pool.Get().(*randv2.Rand)
}

// PutChaCha8 returns a ChaCha8 PRNG to the pool for reuse.
func PutChaCha8(r *randv2.Rand) {
	chacha8Pool.Put(r)
}

// GetPCG retrieves a PCG PRNG from the pool or creates one if needed.
func GetPCG() *randv2.Rand {
	return pcgPool.Get().(*randv2.Rand)
}

// PutPCG returns a PCG PRNG to the pool for reuse.
func PutPCG(r *randv2.Rand) {
	pcgPool.Put(r)
}

// NewChaCha8 initializes a new ChaCha8-based PRNG with a securely randomized seed.
func NewChaCha8() (*randv2.Rand, error) {
	var seed [32]byte
	_, err := rand.Read(seed[:]) // Fill seed with cryptographically secure random data.
	if err != nil {
		return nil, err
	}
	return randv2.New(randv2.NewChaCha8(seed)), nil
}

// WithChaCha8 executes a function using a pooled ChaCha8 PRNG and ensures cleanup.
func WithChaCha8[T any](fn func(rng *randv2.Rand) T) T {
	rng := GetChaCha8()
	res := fn(rng)
	PutChaCha8(rng)
	return res
}

// NewPCG initializes a new PCG-based PRNG with two securely randomized seeds.
func NewPCG() (*randv2.Rand, error) {
	var b [16]byte
	_, err := rand.Read(b[:]) // Fill byte slice with cryptographically secure random data.
	if err != nil {
		return nil, err
	}
	seed1 := binary.BigEndian.Uint64(b[:8]) // Convert first 8 bytes to uint64.
	seed2 := binary.BigEndian.Uint64(b[8:]) // Convert next 8 bytes to uint64.
	return randv2.New(randv2.NewPCG(seed1, seed2)), nil
}

// WithPCG executes a function using a pooled PCG PRNG and ensures cleanup.
func WithPCG[T any](fn func(rng *randv2.Rand) T) T {
	rng := GetPCG()
	res := fn(rng)
	PutPCG(rng)
	return res
}

// Bool returns a randomly chosen boolean value.
func Bool() bool {
	return WithPCG(func(rng *randv2.Rand) bool {
		return rng.Uint64()%2 == 0 // Even numbers are true, odd numbers are false.
	})
}

// Int returns a non-negative pseudo-random int.
func Int() int {
	return WithPCG(func(rng *randv2.Rand) int {
		return rng.Int()
	})
}

// IntN returns a pseudo-random integer in [0,n).
func IntN(n int) int {
	return WithPCG(func(rng *randv2.Rand) int {
		return rng.IntN(n)
	})
}

// IntRange returns a pseudo-random integer in [min,max).
func IntRange(min, max int) int {
	return IntN(max-min) + min
}

// Int32 returns a non-negative pseudo-random 32-bit integer as an int32.
func Int32() int32 {
	return WithPCG(func(rng *randv2.Rand) int32 {
		return rng.Int32()
	})
}

// Int32N returns a pseudo-random 32-bit integer in [0,n).
func Int32N(n int32) int32 {
	return WithPCG(func(rng *randv2.Rand) int32 {
		return rng.Int32N(n)
	})
}

// Int32Range returns a pseudo-random 32-bit integer in [min,max).
func Int32Range(min, max int32) int32 {
	return Int32N(max-min) + min
}

// Int64 returns a non-negative pseudo-random 64-bit integer as an int64.
func Int64() int64 {
	return WithPCG(func(rng *randv2.Rand) int64 {
		return rng.Int64()
	})
}

// Int64N returns a pseudo-random 64-bit integer in [0,n).
func Int64N(n int64) int64 {
	return WithPCG(func(rng *randv2.Rand) int64 {
		return rng.Int64N(n)
	})
}

// Int64Range returns a pseudo-random 64-bit integer in [min,max).
func Int64Range(min, max int64) int64 {
	return Int64N(max-min) + min
}

// Uint returns a pseudo-random unsigned integer.
func Uint() uint {
	return WithPCG(func(rng *randv2.Rand) uint {
		return rng.Uint()
	})
}

// UintN returns a pseudo-random unsigned integer in [0,n).
func UintN(n uint) uint {
	return WithPCG(func(rng *randv2.Rand) uint {
		return rng.UintN(n)
	})
}

// UintRange returns a pseudo-random unsigned integer in [min,max).
func UintRange(min, max uint) uint {
	return UintN(max-min) + min
}

// Uint32 returns a pseudo-random 32-bit unsigned integer as a uint32.
func Uint32() uint32 {
	return WithPCG(func(rng *randv2.Rand) uint32 {
		return rng.Uint32()
	})
}

// Uint32N returns a pseudo-random 32-bit unsigned integer in [0,n).
func Uint32N(n uint32) uint32 {
	return WithPCG(func(rng *randv2.Rand) uint32 {
		return rng.Uint32N(n)
	})
}

// Uint32Range returns a pseudo-random 32-bit unsigned integer in [min,max).
func Uint32Range(min, max uint32) uint32 {
	return Uint32N(max-min) + min
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func Uint64() uint64 {
	return WithPCG(func(rng *randv2.Rand) uint64 {
		return rng.Uint64()
	})
}

// Uint64N returns a pseudo-random 64-bit unsigned integer in [0,n).
func Uint64N(n uint64) uint64 {
	return WithPCG(func(rng *randv2.Rand) uint64 {
		return rng.Uint64N(n)
	})
}

// Uint64Range returns a pseudo-random 64-bit unsigned integer in [min,max).
func Uint64Range(min, max uint64) uint64 {
	return Uint64N(max-min) + min
}

// Float32 returns a pseudo-random float32 in [0.0,1.0).
func Float32() float32 {
	return WithPCG(func(rng *randv2.Rand) float32 {
		return rng.Float32()
	})
}

// Float64 returns a pseudo-random float64 in [0.0,1.0).
func Float64() float64 {
	return WithPCG(func(rng *randv2.Rand) float64 {
		return rng.Float64()
	})
}

// String generates a random string of specified length using the provided character set.
func String(length int, charset string) string {
	if length <= 0 {
		return "" // Return empty string for invalid lengths.
	}
	return WithPCG(func(rng *randv2.Rand) string {
		buf := stringPool.Get(length) // Acquire a buffer from the pool.
		defer stringPool.Put(buf)     // Return the buffer to the pool after use.
		for i, cache, remain := length-1, rng.Int64(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = rng.Int64(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(charset) {
				buf.WriteByte(charset[idx])
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
		return buf.String() // Convert buffer to string before returning.
	})
}
