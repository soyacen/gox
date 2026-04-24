package hmacx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSha224 computes the HMAC-SHA224 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA224 signature as raw bytes
func HmacSha224(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New224, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha224Hex computes the HMAC-SHA224 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA224 signature as a lowercase hex string
func HmacSha224Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha224(key, data))
}
