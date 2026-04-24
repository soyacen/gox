package shax

import (
	"crypto/sha512"
	"encoding/hex"
)

// Sha384Hex returns the hexadecimal encoding of the SHA-384 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-384 digest
func Sha384Hex(data []byte) string {
	return hex.EncodeToString(Sha384(data))
}

// Sha384 returns the SHA-384 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-384 digest
func Sha384(data []byte) []byte {
	digest := sha512.New384()
	digest.Write(data)
	return digest.Sum(nil)
}
