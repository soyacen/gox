package textx

import (
	"bytes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
)

// EUCJPToUtf8 converts EUC-JP encoded bytes to UTF-8.
//
// Parameters:
//   - s: the EUC-JP encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func EUCJPToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.EUCJP.NewDecoder()))
}

// Utf8ToEUCJP converts UTF-8 encoded bytes to EUC-JP.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the EUC-JP encoded result.
//   - error: an error if the conversion fails.
func Utf8ToEUCJP(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.EUCJP.NewEncoder()))
}

// ISO2022JPToUtf8 converts ISO-2022-JP encoded bytes to UTF-8.
//
// Parameters:
//   - s: the ISO-2022-JP encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func ISO2022JPToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.ISO2022JP.NewDecoder()))
}

// Utf8ToISO2022JP converts UTF-8 encoded bytes to ISO-2022-JP.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the ISO-2022-JP encoded result.
//   - error: an error if the conversion fails.
func Utf8ToISO2022JP(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.ISO2022JP.NewEncoder()))
}

// ShiftJISToUtf8 converts Shift_JIS encoded bytes to UTF-8.
//
// Parameters:
//   - s: the Shift_JIS encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func ShiftJISToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.ShiftJIS.NewDecoder()))
}

// Utf8ToShiftJIS converts UTF-8 encoded bytes to Shift_JIS.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the Shift_JIS encoded result.
//   - error: an error if the conversion fails.
func Utf8ToShiftJIS(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), japanese.ShiftJIS.NewEncoder()))
}
