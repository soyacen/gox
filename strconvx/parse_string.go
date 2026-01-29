package strconvx

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

// func ParseWrapperString(s string) (*wrapperspb.StringValue, error) {
// 	wrapperspb.Bytes()
// 	return wrapperspb.String(s), nil
// }
