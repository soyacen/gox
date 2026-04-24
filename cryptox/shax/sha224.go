package shax

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha224Hex returns the hexadecimal encoding of the SHA-224 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-224 digest
func Sha224Hex(data []byte) string {
	return hex.EncodeToString(Sha224(data))
}

// Sha224 returns the SHA-224 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-224 digest
func Sha224(data []byte) []byte {
	digest := sha256.New224()
	digest.Write(data)
	return digest.Sum(nil)
}
