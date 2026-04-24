package hmacx

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// HmacSha512 computes the HMAC-SHA512 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA512 signature as raw bytes
func HmacSha512(key []byte, data []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha512Hex computes the HMAC-SHA512 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA512 signature as a lowercase hex string
func HmacSha512Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha512(key, data))
}
