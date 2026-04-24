package shax

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hex returns the hexadecimal encoding of the SHA-256 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-256 digest
func Sha256Hex(data []byte) string {
	return hex.EncodeToString(Sha256(data))
}

// Sha256 returns the SHA-256 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-256 digest
func Sha256(data []byte) []byte {
	digest := sha256.New()
	digest.Write(data)
	return digest.Sum(nil)
}
