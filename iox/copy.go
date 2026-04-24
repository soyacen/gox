package iox

import "io"

// Copy copies data from src to dst.
// It wraps io.Copy and returns any error encountered.
//
// Parameters:
//   - dst: the destination writer
//   - src: the source reader
//
// Returns:
//   - error: any error encountered during the copy
func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}
