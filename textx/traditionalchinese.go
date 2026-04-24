package textx

import (
	"bytes"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io"
)

// Big5ToUtf8 converts Big5 encoded bytes to UTF-8.
//
// Parameters:
//   - s: the Big5 encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func Big5ToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), traditionalchinese.Big5.NewDecoder()))
}

// Utf8ToBig5 converts UTF-8 encoded bytes to Big5.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the Big5 encoded result.
//   - error: an error if the conversion fails.
func Utf8ToBig5(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), traditionalchinese.Big5.NewEncoder()))
}
