package md5x

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// MD5 computes the MD5 hash of the given data.
//
// Parameters:
//   - data: the input bytes to hash.
//
// Returns:
//   - []byte: the raw MD5 digest.
func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// MD5Hex computes the MD5 hash and returns it as a hexadecimal string.
//
// Parameters:
//   - data: the input bytes to hash.
//
// Returns:
//   - string: the MD5 digest encoded in lowercase hexadecimal.
func MD5Hex(data []byte) string {
	return hex.EncodeToString(MD5(data))
}

// TextMD5 computes the MD5 hash of the given text string.
//
// Parameters:
//   - text: the input string to hash.
//
// Returns:
//   - []byte: the raw MD5 digest.
func TextMD5(text string) []byte {
	hash := md5.New()
	_, _ = io.WriteString(hash, text)
	return hash.Sum(nil)
}

// TextMD5Hex computes the MD5 hash of the given text string and returns it as a hexadecimal string.
//
// Parameters:
//   - text: the input string to hash.
//
// Returns:
//   - string: the MD5 digest encoded in lowercase hexadecimal.
func TextMD5Hex(text string) string {
	return hex.EncodeToString(TextMD5(text))
}
