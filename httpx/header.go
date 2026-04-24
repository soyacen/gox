package httpx

import "net/http"

// CopyHeader copies all header values from src to tgt.
// Existing values in tgt for the same keys are preserved.
//
// Parameters:
//   - tgt: The target header to copy values to
//   - src: The source header to copy values from
func CopyHeader(tgt http.Header, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			tgt.Add(key, value)
		}
	}
}
