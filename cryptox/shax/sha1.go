package shax

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1Hex returns the hexadecimal encoding of the SHA-1 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-1 digest
func Sha1Hex(data []byte) string {
	return hex.EncodeToString(Sha1(data))
}

// Sha1 returns the SHA-1 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-1 digest
func Sha1(data []byte) []byte {
	digest := sha1.New()
	digest.Write(data)
	return digest.Sum(nil)
}
