package httpx

import (
	"context"
	"net/http"
)

// clientKey is the context key for storing an *http.Client.
type clientKey struct{}

// ClientFromContext retrieves an *http.Client from the context.
//
// Parameters:
//   - ctx: The context to retrieve the client from
//
// Returns:
//   - *http.Client: The HTTP client, or nil if not found
//   - bool: True if the client was found in the context
func ClientFromContext(ctx context.Context) (*http.Client, bool) {
	cli, ok := ctx.Value(clientKey{}).(*http.Client)
	return cli, ok
}

// NewContextWithClient stores an *http.Client in the context.
//
// Parameters:
//   - ctx: The parent context
//   - cli: The HTTP client to store
//
// Returns:
//   - context.Context: A new context containing the HTTP client
func NewContextWithClient(ctx context.Context, cli *http.Client) context.Context {
	return context.WithValue(ctx, clientKey{}, cli)
}

// requestKey is the context key for storing an *http.Request.
type requestKey struct{}

// RequestFromContext retrieves an *http.Request from the context.
//
// Parameters:
//   - ctx: The context to retrieve the request from
//
// Returns:
//   - *http.Request: The HTTP request, or nil if not found
//   - bool: True if the request was found in the context
func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	req, ok := ctx.Value(requestKey{}).(*http.Request)
	return req, ok
}

// NewContextWithRequest stores an *http.Request in the context.
//
// Parameters:
//   - ctx: The parent context
//   - req: The HTTP request to store
//
// Returns:
//   - context.Context: A new context containing the HTTP request
func NewContextWithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

// responseKey is the context key for storing an *http.Response.
type responseKey struct{}

// ResponseFromContext retrieves an *http.Response from the context.
//
// Parameters:
//   - ctx: The context to retrieve the response from
//
// Returns:
//   - *http.Response: The HTTP response, or nil if not found
//   - bool: True if the response was found in the context
func ResponseFromContext(ctx context.Context) (*http.Response, bool) {
	resp, ok := ctx.Value(responseKey{}).(*http.Response)
	return resp, ok
}

// NewContextWithResponse stores an *http.Response in the context.
//
// Parameters:
//   - ctx: The parent context
//   - resp: The HTTP response to store
//
// Returns:
//   - context.Context: A new context containing the HTTP response
func NewContextWithResponse(ctx context.Context, resp *http.Response) context.Context {
	return context.WithValue(ctx, responseKey{}, resp)
}
