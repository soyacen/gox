package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// TransportBuilder is a builder for creating custom http.Transport instances.
type TransportBuilder struct {
	proxy                  func(*http.Request) (*url.URL, error)
	dial                   func(ctx context.Context, network string, addr string) (net.Conn, error)
	dialTLS                func(ctx context.Context, network string, addr string) (net.Conn, error)
	tlsConfig              *tls.Config
	tlsHandshakeTimeout    time.Duration
	disableKeepAlives      bool
	disableCompression     bool
	maxIdleConns           int
	maxIdleConnsPerHost    int
	maxConnsPerHost        int
	idleConnTimeout        time.Duration
	responseHeaderTimeout  time.Duration
	expectContinueTimeout  time.Duration
	tlsNextProto           map[string]func(authority string, c *tls.Conn) http.RoundTripper
	proxyConnectHeader     http.Header
	getProxyConnectHeader  func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error)
	maxResponseHeaderBytes int64
	writeBufferSize        int
	readBufferSize         int
	forceAttemptHTTP2      bool
}

// Proxy sets the proxy function for the transport.
//
// Parameters:
//   - proxy: The proxy function
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) Proxy(proxy func(*http.Request) (*url.URL, error)) *TransportBuilder {
	builder.proxy = proxy
	return builder
}

// Dial sets the dial function for the transport.
//
// Parameters:
//   - dial: The dial function
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) Dial(dial func(ctx context.Context, network string, addr string) (net.Conn, error)) *TransportBuilder {
	builder.dial = dial
	return builder
}

// DialTLS sets the TLS dial function for the transport.
//
// Parameters:
//   - dialTLS: The TLS dial function
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) DialTLS(dialTLS func(ctx context.Context, network string, addr string) (net.Conn, error)) *TransportBuilder {
	builder.dialTLS = dialTLS
	return builder
}

// TLSConfig sets the TLS configuration for the transport.
//
// Parameters:
//   - tlsConfig: The TLS configuration
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) TLSConfig(tlsConfig *tls.Config) *TransportBuilder {
	builder.tlsConfig = tlsConfig
	return builder
}

// TLSHandshakeTimeout sets the TLS handshake timeout for the transport.
//
// Parameters:
//   - timeout: The TLS handshake timeout
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) TLSHandshakeTimeout(timeout time.Duration) *TransportBuilder {
	builder.tlsHandshakeTimeout = timeout
	return builder
}

// DisableKeepAlives enables or disables HTTP keep-alives.
//
// Parameters:
//   - disable: Whether to disable keep-alives
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) DisableKeepAlives(disable bool) *TransportBuilder {
	builder.disableKeepAlives = disable
	return builder
}

// DisableCompression enables or disables compression.
//
// Parameters:
//   - disable: Whether to disable compression
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) DisableCompression(disable bool) *TransportBuilder {
	builder.disableCompression = disable
	return builder
}

// MaxIdleConns sets the maximum number of idle connections.
//
// Parameters:
//   - n: The maximum number of idle connections
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) MaxIdleConns(n int) *TransportBuilder {
	builder.maxIdleConns = n
	return builder
}

// MaxIdleConnsPerHost sets the maximum number of idle connections per host.
//
// Parameters:
//   - n: The maximum number of idle connections per host
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) MaxIdleConnsPerHost(n int) *TransportBuilder {
	builder.maxIdleConnsPerHost = n
	return builder
}

// MaxConnsPerHost sets the maximum number of connections per host.
//
// Parameters:
//   - n: The maximum number of connections per host
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) MaxConnsPerHost(n int) *TransportBuilder {
	builder.maxConnsPerHost = n
	return builder
}

// IdleConnTimeout sets the idle connection timeout.
//
// Parameters:
//   - timeout: The idle connection timeout
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) IdleConnTimeout(timeout time.Duration) *TransportBuilder {
	builder.idleConnTimeout = timeout
	return builder
}

// ResponseHeaderTimeout sets the response header timeout.
//
// Parameters:
//   - timeout: The response header timeout
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) ResponseHeaderTimeout(timeout time.Duration) *TransportBuilder {
	builder.responseHeaderTimeout = timeout
	return builder
}

// ExpectContinueTimeout sets the expect continue timeout.
//
// Parameters:
//   - timeout: The expect continue timeout
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) ExpectContinueTimeout(timeout time.Duration) *TransportBuilder {
	builder.expectContinueTimeout = timeout
	return builder
}

// TLSNextProto sets the TLS next protocol negotiation.
//
// Parameters:
//   - f: The next protocol map
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) TLSNextProto(f map[string]func(authority string, c *tls.Conn) http.RoundTripper) *TransportBuilder {
	builder.tlsNextProto = f
	return builder
}

// ProxyConnectHeader sets the proxy connect header.
//
// Parameters:
//   - h: The proxy connect header
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) ProxyConnectHeader(h http.Header) *TransportBuilder {
	builder.proxyConnectHeader = h
	return builder
}

// GetProxyConnectHeader sets the function to get proxy connect headers.
//
// Parameters:
//   - f: The function to get proxy connect headers
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) GetProxyConnectHeader(f func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error)) *TransportBuilder {
	builder.getProxyConnectHeader = f
	return builder
}

// MaxResponseHeaderBytes sets the maximum response header bytes.
//
// Parameters:
//   - n: The maximum response header bytes
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) MaxResponseHeaderBytes(n int64) *TransportBuilder {
	builder.maxResponseHeaderBytes = n
	return builder
}

// WriteBufferSize sets the write buffer size.
//
// Parameters:
//   - n: The write buffer size
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) WriteBufferSize(n int) *TransportBuilder {
	builder.writeBufferSize = n
	return builder
}

// ReadBufferSize sets the read buffer size.
//
// Parameters:
//   - n: The read buffer size
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) ReadBufferSize(n int) *TransportBuilder {
	builder.readBufferSize = n
	return builder
}

// ForceAttemptHTTP2 enables or disables HTTP/2 support.
//
// Parameters:
//   - enable: Whether to enable HTTP/2
//
// Returns:
//   - *TransportBuilder: The builder instance for method chaining
func (builder *TransportBuilder) ForceAttemptHTTP2(enable bool) *TransportBuilder {
	builder.forceAttemptHTTP2 = enable
	return builder
}

// Build creates and returns a new http.Transport with the configured settings.
//
// Returns:
//   - *http.Transport: A new HTTP transport instance
func (builder *TransportBuilder) Build() *http.Transport {
	return &http.Transport{
		Proxy:                  builder.proxy,
		DialContext:            builder.dial,
		DialTLSContext:         builder.dialTLS,
		TLSClientConfig:        builder.tlsConfig,
		TLSHandshakeTimeout:    builder.tlsHandshakeTimeout,
		DisableKeepAlives:      builder.disableKeepAlives,
		DisableCompression:     builder.disableCompression,
		MaxIdleConns:           builder.maxIdleConns,
		MaxIdleConnsPerHost:    builder.maxIdleConnsPerHost,
		MaxConnsPerHost:        builder.maxConnsPerHost,
		IdleConnTimeout:        builder.idleConnTimeout,
		ResponseHeaderTimeout:  builder.responseHeaderTimeout,
		ExpectContinueTimeout:  builder.expectContinueTimeout,
		TLSNextProto:           builder.tlsNextProto,
		ProxyConnectHeader:     builder.proxyConnectHeader,
		GetProxyConnectHeader:  builder.getProxyConnectHeader,
		MaxResponseHeaderBytes: builder.maxResponseHeaderBytes,
		WriteBufferSize:        builder.writeBufferSize,
		ReadBufferSize:         builder.readBufferSize,
		ForceAttemptHTTP2:      builder.forceAttemptHTTP2,
	}
}

// DisableKeepAlivesTransport returns a new http.Transport with similar default values to
// http.DefaultTransport, but with idle connections and keepalives disabled.
//
// Returns:
//   - *http.Transport: A transport with keep-alives disabled
func DisableKeepAlivesTransport() *http.Transport {
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	return new(TransportBuilder).
		Proxy(http.ProxyFromEnvironment).
		Dial(dialer.DialContext).
		ForceAttemptHTTP2(true).
		MaxIdleConns(100).
		IdleConnTimeout(90 * time.Second).
		TLSHandshakeTimeout(10 * time.Second).
		ExpectContinueTimeout(time.Second).
		DisableKeepAlives(true).
		MaxIdleConns(-1).Build()
}

// PooledTransport returns a new http.Transport with similar default
// values to http.DefaultTransport. Do not use this for transient transports as
// it can leak file descriptors over time. Only use this for transports that
// will be re-used for the same host(s).
//
// Returns:
//   - *http.Transport: A transport with connection pooling enabled
func PooledTransport() *http.Transport {
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	return new(TransportBuilder).
		Proxy(http.ProxyFromEnvironment).
		Dial(dialer.DialContext).
		ForceAttemptHTTP2(true).
		MaxIdleConns(100).
		IdleConnTimeout(90 * time.Second).
		TLSHandshakeTimeout(10 * time.Second).
		ExpectContinueTimeout(time.Second).
		MaxIdleConnsPerHost(runtime.GOMAXPROCS(0) + 1).Build()
}
