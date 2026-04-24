package textx

import (
	"bytes"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"io"
)

// EUCKRToUtf8 converts EUC-KR encoded bytes to UTF-8.
//
// Parameters:
//   - s: the EUC-KR encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func EUCKRToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), korean.EUCKR.NewDecoder()))
}

// Utf8ToEUCKR converts UTF-8 encoded bytes to EUC-KR.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the EUC-KR encoded result.
//   - error: an error if the conversion fails.
func Utf8ToEUCKR(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), korean.EUCKR.NewEncoder()))
}
