package hmacx

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

// HmacMD5 computes the HMAC-MD5 signature of data using the provided key.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - []byte: the HMAC-MD5 signature as raw bytes
func HmacMD5(key []byte, data []byte) []byte {
	h := hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HmacMD5Hex computes the HMAC-MD5 signature of data using the provided key and returns it as a lowercase hex string.
//
// Parameters:
//   - key: the HMAC key
//   - data: the data to sign
//
// Returns:
//   - string: the HMAC-MD5 signature as a lowercase hex string
func HmacMD5Hex(key []byte, data []byte) string {
	return hex.EncodeToString(HmacMD5(key, data))
}
