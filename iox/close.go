package iox

import "io"

// QuiteClose closes the closer quietly, ignoring any errors.
// It checks if the closer is not nil before calling Close.
//
// Parameters:
//   - closer: the io.Closer to close
func QuiteClose(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}
