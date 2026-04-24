package hmacx

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// HmacSha512_224 computes the HMAC-SHA512/224 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA512/224 signature as raw bytes
func HmacSha512_224(key []byte, data []byte) []byte {
	h := hmac.New(sha512.New512_224, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha512_224Hex computes the HMAC-SHA512/224 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA512/224 signature as a lowercase hex string
func HmacSha512_224Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha512_224(key, data))
}
