package hmacx

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// HmacSha512_256 computes the HMAC-SHA512/256 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA512/256 signature as raw bytes
func HmacSha512_256(key []byte, data []byte) []byte {
	h := hmac.New(sha512.New512_256, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha512_256Hex computes the HMAC-SHA512/256 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA512/256 signature as a lowercase hex string
func HmacSha512_256Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha512_256(key, data))
}
