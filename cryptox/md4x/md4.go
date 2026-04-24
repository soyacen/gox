package md4x

import (
	"encoding/hex"

	"golang.org/x/crypto/md4" // nolint
)

// MD4 computes the MD4 hash of the given data.
//
// Parameters:
//   - data: the input bytes to hash.
//
// Returns:
//   - []byte: the raw MD4 digest.
func MD4(data []byte) []byte {
	hash := md4.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// MD4Hex computes the MD4 hash and returns it as a hexadecimal string.
//
// Parameters:
//   - data: the input bytes to hash.
//
// Returns:
//   - string: the MD4 digest encoded in lowercase hexadecimal.
func MD4Hex(data []byte) string {
	return hex.EncodeToString(MD4(data))
}
