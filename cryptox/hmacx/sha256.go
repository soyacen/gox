package hmacx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSha256 computes the HMAC-SHA256 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA256 signature as raw bytes
func HmacSha256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha256Hex computes the HMAC-SHA256 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA256 signature as a lowercase hex string
func HmacSha256Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha256(key, data))
}
