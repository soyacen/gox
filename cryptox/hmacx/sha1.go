package hmacx

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// HmacSha1 computes the HMAC-SHA1 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA1 signature as raw bytes
func HmacSha1(key []byte, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha1Hex computes the HMAC-SHA1 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA1 signature as a lowercase hex string
func HmacSha1Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha1(key, data))
}
