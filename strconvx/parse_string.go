package strconvx

// ParseBytesSlice converts a slice of strings to a slice of byte slices.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to convert.
//
// Returns:
//   - [][]byte: the converted byte slices.
func ParseBytesSlice(s []string) [][]byte {
	if s == nil {
		return nil
	}
	r := make([][]byte, 0, len(s))
	for _, str := range s {
		r = append(r, []byte(str))
	}
	return r
}
