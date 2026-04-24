package shax

import (
	"crypto/sha512"
	"encoding/hex"
)

// Sha512_256Hex returns the hexadecimal encoding of the SHA-512/256 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-512/256 digest
func Sha512_256Hex(data []byte) string {
	return hex.EncodeToString(Sha512_256(data))
}

// Sha512_256 returns the SHA-512/256 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-512/256 digest
func Sha512_256(data []byte) []byte {
	digest := sha512.New512_256()
	digest.Write(data)
	return digest.Sum(nil)
}
