package textx

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

// GBKToUtf8 converts GBK encoded bytes to UTF-8.
//
// Parameters:
//   - s: the GBK encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func GBKToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder()))
}

// Utf8ToGBK converts UTF-8 encoded bytes to GBK.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the GBK encoded result.
//   - error: an error if the conversion fails.
func Utf8ToGBK(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder()))
}

// GB18030ToUtf8 converts GB18030 encoded bytes to UTF-8.
//
// Parameters:
//   - s: the GB18030 encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func GB18030ToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewDecoder()))
}

// Utf8ToGB18030 converts UTF-8 encoded bytes to GB18030.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the GB18030 encoded result.
//   - error: an error if the conversion fails.
func Utf8ToGB18030(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewEncoder()))
}

// HZGB2312ToUtf8 converts HZ-GB2312 encoded bytes to UTF-8.
//
// Parameters:
//   - s: the HZ-GB2312 encoded byte slice.
//
// Returns:
//   - []byte: the UTF-8 encoded result.
//   - error: an error if the conversion fails.
func HZGB2312ToUtf8(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewDecoder()))
}

// Utf8ToHZGB2312 converts UTF-8 encoded bytes to HZ-GB2312.
//
// Parameters:
//   - s: the UTF-8 encoded byte slice.
//
// Returns:
//   - []byte: the HZ-GB2312 encoded result.
//   - error: an error if the conversion fails.
func Utf8ToHZGB2312(s []byte) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewEncoder()))
}
