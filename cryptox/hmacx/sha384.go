package hmacx

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// HmacSha384 computes the HMAC-SHA384 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-SHA384 signature as raw bytes
func HmacSha384(key []byte, data []byte) []byte {
	h := hmac.New(sha512.New384, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacSha384Hex computes the HMAC-SHA384 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-SHA384 signature as a lowercase hex string
func HmacSha384Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacSha384(key, data))
}
