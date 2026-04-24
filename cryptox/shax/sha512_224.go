package shax

import (
	"crypto/sha512"
	"encoding/hex"
)

// Sha512_224Hex returns the hexadecimal encoding of the SHA-512/224 hash of data.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - string: the hex-encoded SHA-512/224 digest
func Sha512_224Hex(data []byte) string {
	return hex.EncodeToString(Sha512_224(data))
}

// Sha512_224 returns the SHA-512/224 hash of data as a byte slice.
//
// Parameters:
//   - data: the input bytes to hash
//
// Returns:
//   - []byte: the raw SHA-512/224 digest
func Sha512_224(data []byte) []byte {
	digest := sha512.New512_224()
	digest.Write(data)
	return digest.Sum(nil)
}
