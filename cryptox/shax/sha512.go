package shax

import (
	"crypto/sha512"
	"encoding/hex"
)

// Sha512Hex returns the hexadecimal encoding of the SHA-512 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-512 digest
func Sha512Hex(data []byte) string {
	return hex.EncodeToString(Sha512(data))
}

// Sha512 returns the SHA-512 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-512 digest
func Sha512(data []byte) []byte {
	digest := sha512.New()
	digest.Write(data)
	return digest.Sum(nil)
}
