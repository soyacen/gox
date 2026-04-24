package httpx

import (
	"net/http"
	"time"
)

// ClientBuilder is a builder for creating custom http.Client instances.
// It allows configuring transport, redirect handling, cookies, and timeout settings.
type ClientBuilder struct {
	transport     http.RoundTripper
	checkRedirect func(req *http.Request, via []*http.Request) error
	jar           http.CookieJar
	timeout       time.Duration
}

// Transport sets the HTTP transport for the client.
//
// Parameters:
//   - transport: The RoundTripper to use for HTTP requests
//
// Returns:
//   - *ClientBuilder: The builder instance for method chaining
func (builder *ClientBuilder) Transport(transport http.RoundTripper) *ClientBuilder {
	builder.transport = transport
	return builder
}

// CheckRedirect sets the redirect policy function for the client.
//
// Parameters:
//   - f: Function that determines whether to follow redirects
//
// Returns:
//   - *ClientBuilder: The builder instance for method chaining
func (builder *ClientBuilder) CheckRedirect(f func(req *http.Request, via []*http.Request) error) *ClientBuilder {
	builder.checkRedirect = f
	return builder
}

// Jar sets the cookie jar for the client.
//
// Parameters:
//   - jar: The CookieJar to use for storing cookies
//
// Returns:
//   - *ClientBuilder: The builder instance for method chaining
func (builder *ClientBuilder) Jar(jar http.CookieJar) *ClientBuilder {
	builder.jar = jar
	return builder
}

// Timeout sets the request timeout for the client.
//
// Parameters:
//   - timeout: The maximum time to wait for a request to complete
//
// Returns:
//   - *ClientBuilder: The builder instance for method chaining
func (builder *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder {
	builder.timeout = timeout
	return builder
}

// Build creates and returns a new http.Client with the configured settings.
//
// Returns:
//   - *http.Client: A new HTTP client instance
func (builder *ClientBuilder) Build() *http.Client {
	return &http.Client{
		Transport:     builder.transport,
		CheckRedirect: builder.checkRedirect,
		Jar:           builder.jar,
		Timeout:       builder.timeout,
	}
}

// DisableKeepAlivesClient returns a new http.Client with similar default values to
// http.Client, but with a non-shared Transport, idle connections disabled, and
// keepalives disabled.
//
// Returns:
//   - *http.Client: An HTTP client with keep-alives disabled
func DisableKeepAlivesClient() *http.Client {
	return new(ClientBuilder).Transport(DisableKeepAlivesTransport()).Build()
}

// PooledClient returns a new http.Client with similar default values to
// http.Client, but with a shared Transport. Do not use this function for
// transient clients as it can leak file descriptors over time. Only use this
// for clients that will be re-used for the same host(s).
//
// Returns:
//   - *http.Client: An HTTP client with connection pooling enabled
func PooledClient() *http.Client {
	return new(ClientBuilder).Transport(PooledTransport()).Build()
}
