package strconvx

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

var quotePool = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, 16)) }}

// Quote quotes a string with the given quote characters.
//
// Parameters:
//   - e: the string to quote.
//   - quote: the quote string to wrap around e.
//
// Returns:
//   - E: the quoted string.
func Quote[E ~string](e E, quote string) E {
	buffer := quotePool.Get().(*bytes.Buffer)
	defer quotePool.Put(buffer)
	buffer.Reset()
	buffer.WriteString(quote)
	buffer.WriteString(string(e))
	buffer.WriteString(quote)
	return E(buffer.String())
}

func quoteV2[E ~string](e E, quote string) E {
	buffer := quotePool.Get().(*bytes.Buffer)
	defer quotePool.Put(buffer)
	buffer.Reset()
	_, _ = buffer.WriteString(fmt.Sprintf("%s%s%s", quote, e, quote))
	return E(buffer.String())
}

func quoteV3[E ~string](e E, quote string) E {
	return E(strings.Join([]string{quote, string(e), quote}, ""))
}

// QuoteSlice quotes each string element in the slice with the given quote characters.
// Returns nil if the input slice is nil.
//
// Parameters:
//   - s: the slice of strings to quote.
//   - quote: the quote string to wrap around each element.
//
// Returns:
//   - S: the slice with each element quoted.
func QuoteSlice[S ~[]E, E ~string](s S, quote string) S {
	if s == nil {
		return s
	}
	r := make(S, 0, len(s))
	for _, e := range s {
		r = append(r, Quote(e, quote))
	}
	return r
}
