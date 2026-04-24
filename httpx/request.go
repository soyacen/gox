package httpx

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"

	"github.com/soyacen/gox/stringx"
)

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// RequestBuilder is a builder for constructing HTTP requests.
type RequestBuilder struct {
	err     error
	method  string
	uri     *url.URL
	queries url.Values
	headers http.Header
	body    io.Reader
	cookies []*http.Cookie
	req     *http.Request
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// NewRequestBuilder creates a new RequestBuilder.
//
// Returns:
//   - *RequestBuilder: A new request builder instance
func NewRequestBuilder() *RequestBuilder {
	return new(RequestBuilder)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Method sets the HTTP method for the request.
//
// Parameters:
//   - method: The HTTP method to use
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Method(method string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	if stringx.IsBlank(method) {
		builder.err = errors.New("method is blank")
		return builder
	}
	builder.method = method
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Get sets the HTTP method to GET.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Get() *RequestBuilder {
	return builder.Method(http.MethodGet)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Head sets the HTTP method to HEAD.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Head() *RequestBuilder {
	return builder.Method(http.MethodHead)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Post sets the HTTP method to POST.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Post() *RequestBuilder {
	return builder.Method(http.MethodPost)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Put sets the HTTP method to PUT.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Put() *RequestBuilder {
	return builder.Method(http.MethodPut)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Patch sets the HTTP method to PATCH.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Patch() *RequestBuilder {
	return builder.Method(http.MethodPatch)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Delete sets the HTTP method to DELETE.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Delete() *RequestBuilder {
	return builder.Method(http.MethodDelete)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Connect sets the HTTP method to CONNECT.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Connect() *RequestBuilder {
	return builder.Method(http.MethodConnect)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Options sets the HTTP method to OPTIONS.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Options() *RequestBuilder {
	return builder.Method(http.MethodOptions)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Trace sets the HTTP method to TRACE.
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Trace() *RequestBuilder {
	return builder.Method(http.MethodTrace)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// URL sets the request URL.
//
// Parameters:
//   - uri: The URL to set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) URL(uri *url.URL) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.uri = uri
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// URLString parses and sets the request URL from a string.
// Automatically converts ws:// and wss:// schemes to http:// and https://.
//
// Parameters:
//   - urlString: The URL string to parse and set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) URLString(urlString string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	if strings.HasPrefix(urlString, "ws:") {
		urlString = "http:" + strings.TrimPrefix(urlString, "ws:")
	} else if strings.HasPrefix(urlString, "wss") {
		urlString = "http:" + strings.TrimPrefix(urlString, "https:")
	}
	uri, err := url.Parse(urlString)
	if err != nil {
		builder.err = err
		return builder
	}
	return builder.URL(uri)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// query returns the query values, initializing them if necessary.
//
// Returns:
//   - url.Values: The query values
func (builder *RequestBuilder) query() url.Values {
	if builder.queries == nil {
		builder.queries = make(url.Values)
	}
	return builder.queries
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Query sets a query parameter.
//
// Parameters:
//   - name: The query parameter name
//   - value: The query parameter value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Query(name, value string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.query().Set(name, value)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// AddQuery adds a query parameter (allows multiple values for the same key).
//
// Parameters:
//   - key: The query parameter name
//   - value: The query parameter value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) AddQuery(key, value string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.query().Add(key, value)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// RemoveQuery removes a query parameter.
//
// Parameters:
//   - name: The query parameter name to remove
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) RemoveQuery(name string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.query().Del(name)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// QueryString parses and sets query parameters from a query string.
//
// Parameters:
//   - q: The query string to parse
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) QueryString(q string) *RequestBuilder {
	queries, err := url.ParseQuery(q)
	if err != nil {
		builder.err = err
		return builder
	}
	return builder.Queries(queries)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Queries sets multiple query parameters from url.Values.
//
// Parameters:
//   - queries: The query values to set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Queries(queries url.Values) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	for key, values := range queries {
		for _, value := range values {
			builder.query().Add(key, value)
		}
	}
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// header returns the HTTP headers, initializing them if necessary.
//
// Returns:
//   - http.Header: The HTTP headers
func (builder *RequestBuilder) header() http.Header {
	if builder.headers == nil {
		builder.headers = make(http.Header)
	}
	return builder.headers
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Header sets an HTTP header.
//
// Parameters:
//   - name: The header name
//   - value: The header value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Header(name, value string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.header().Set(name, value)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// AddHeader adds an HTTP header (allows multiple values for the same key).
//
// Parameters:
//   - name: The header name
//   - value: The header value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) AddHeader(name, value string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.header().Add(name, value)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// RemoveHeader removes an HTTP header.
//
// Parameters:
//   - name: The header name to remove
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) RemoveHeader(name string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.header().Del(name)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Headers sets multiple HTTP headers from an http.Header.
//
// Parameters:
//   - header: The headers to set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Headers(header http.Header) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	for key, values := range header {
		for _, value := range values {
			builder.header().Add(key, value)
		}
	}
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// UserAgent sets the User-Agent header.
//
// Parameters:
//   - ua: The User-Agent string
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) UserAgent(ua string) *RequestBuilder {
	return builder.Header("User-Agent", ua)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// IfModifiedSince sets the If-Modified-Since header.
//
// Parameters:
//   - time: The time value for the header
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) IfModifiedSince(time string) *RequestBuilder {
	return builder.Header("If-Modified-Since", time)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// IfUnmodifiedSince sets the If-Unmodified-Since header.
//
// Parameters:
//   - time: The time value for the header
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) IfUnmodifiedSince(time string) *RequestBuilder {
	return builder.Header("If-Unmodified-Since", time)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// IfNoneMatch sets the If-None-Match header.
//
// Parameters:
//   - etag: The ETag value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) IfNoneMatch(etag string) *RequestBuilder {
	return builder.Header("If-None-Match", etag)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// IfMatch sets the If-Match header.
//
// Parameters:
//   - etags: The ETag values
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) IfMatch(etags ...string) *RequestBuilder {
	return builder.Header("If-Match", strings.Join(etags, ", "))
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// CacheControl sets the Cache-Control header.
//
// Parameters:
//   - directives: The cache control directives
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) CacheControl(directives ...string) *RequestBuilder {
	return builder.Header("Cache-Control", strings.Join(directives, ", "))
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Body sets the request body with the specified content type.
// Automatically sets the Content-Length header if the body implements Len().
//
// Parameters:
//   - body: The request body
//   - contentType: The content type of the body
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Body(body io.Reader, contentType string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.body = body
	builder.Header("Content-Type", contentType)
	if lenGetter, ok := body.(interface{ Len() int }); ok {
		builder.Header("Content-Length", strconv.Itoa(lenGetter.Len()))
	}
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// BytesBody sets the request body from a byte slice.
//
// Parameters:
//   - body: The request body bytes
//   - contentType: The content type of the body
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) BytesBody(body []byte, contentType string) *RequestBuilder {
	reader := bytes.NewReader(body)
	return builder.Body(reader, contentType)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// TextBody sets the request body from a string.
//
// Parameters:
//   - body: The request body text
//   - contentType: The content type of the body
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) TextBody(body string, contentType string) *RequestBuilder {
	reader := strings.NewReader(body)
	return builder.Body(reader, contentType)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// FormBody sets the request body from form values with application/x-www-form-urlencoded.
//
// Parameters:
//   - form: The form values
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) FormBody(form url.Values) *RequestBuilder {
	return builder.TextBody(form.Encode(), "application/x-www-form-urlencoded")
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// ObjectBody sets the request body by marshaling an object.
//
// Parameters:
//   - body: The object to marshal
//   - marshal: The marshal function
//   - contentType: The content type of the body
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) ObjectBody(body any, marshal func(any) ([]byte, error), contentType string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	data, err := marshal(body)
	if err != nil {
		builder.err = err
		return builder
	}
	return builder.BytesBody(data, contentType)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// JSONBody sets the request body by marshaling an object as JSON.
//
// Parameters:
//   - body: The object to marshal as JSON
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) JSONBody(body any) *RequestBuilder {
	return builder.ObjectBody(body, func(v any) ([]byte, error) {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(v)
		return buffer.Bytes(), err
	}, "application/json")
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// XMLBody sets the request body by marshaling an object as XML.
//
// Parameters:
//   - body: The object to marshal as XML
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) XMLBody(body any) *RequestBuilder {
	return builder.ObjectBody(body, xml.Marshal, "application/xml")
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// ProtobufBody sets the request body by marshaling a protobuf message.
//
// Parameters:
//   - body: The protobuf message to marshal
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) ProtobufBody(body proto.Message) *RequestBuilder {
	marshal := func(v any) ([]byte, error) {
		message, _ := v.(proto.Message)
		return proto.Marshal(message)
	}
	return builder.ObjectBody(body, marshal, "application/x-protobuf")
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// GobBody sets the request body by marshaling an object as Gob.
//
// Parameters:
//   - body: The object to marshal as Gob
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) GobBody(body any) *RequestBuilder {
	marshal := func(v any) ([]byte, error) {
		var b bytes.Buffer
		if err := gob.NewEncoder(&b).Encode(v); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	}
	return builder.ObjectBody(body, marshal, "application/x-gob")
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// FormData represents a single field in a multipart form.
type FormData struct {
	FieldName string
	Value     string
	File      io.Reader
	Filename  string
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// MultipartBody sets the request body as multipart form data.
//
// Parameters:
//   - formData: The form data fields to include
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) MultipartBody(formData ...*FormData) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	payload := new(bytes.Buffer)
	writer := multipart.NewWriter(payload)
	for _, form := range formData {
		if form.File != nil {
			mf, err := writer.CreateFormFile(form.FieldName, filepath.Base(form.Filename))
			if err != nil {
				builder.err = err
				return builder
			}
			if _, err = io.Copy(mf, form.File); err != nil {
				builder.err = err
				return builder
			}
		} else {
			_ = writer.WriteField(form.FieldName, form.Value)
		}
	}
	if err := writer.Close(); err != nil {
		builder.err = err
		return builder
	}
	return builder.Body(payload, writer.FormDataContentType())
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// BasicAuth sets the Authorization header with Basic authentication.
//
// Parameters:
//   - username: The username
//   - password: The password
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) BasicAuth(username, password string) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	token := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return builder.CustomAuth("Basic", token)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// BearerAuth sets the Authorization header with Bearer authentication.
//
// Parameters:
//   - token: The bearer token
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) BearerAuth(token string) *RequestBuilder {
	return builder.CustomAuth("Bearer", token)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// CustomAuth sets the Authorization header with a custom scheme.
//
// Parameters:
//   - scheme: The authentication scheme
//   - token: The authentication token
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) CustomAuth(scheme, token string) *RequestBuilder {
	return builder.APIKey("Authorization", scheme+" "+token)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// APIKey sets a header with the specified key and value.
//
// Parameters:
//   - key: The header name
//   - value: The header value
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) APIKey(key string, value string) *RequestBuilder {
	return builder.Header(key, value)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Cookie sets a cookie, replacing any existing cookie with the same name.
//
// Parameters:
//   - cookie: The cookie to set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Cookie(cookie *http.Cookie) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	index := slices.IndexFunc(builder.cookies, func(c *http.Cookie) bool {
		return c.Name == cookie.Name
	})
	if index >= 0 {
		builder.cookies[index] = cookie
		return builder
	}
	return builder.AddCookie(cookie)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// AddCookie adds a cookie to the request.
//
// Parameters:
//   - cookie: The cookie to add
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) AddCookie(cookie *http.Cookie) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.cookies = append(builder.cookies, cookie)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// RemoveCookie removes a cookie from the request by name.
//
// Parameters:
//   - cookie: The cookie to remove
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) RemoveCookie(cookie *http.Cookie) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	index := slices.IndexFunc(builder.cookies, func(c *http.Cookie) bool {
		return c.Name == cookie.Name
	})
	if index == -1 {
		return builder
	}
	builder.cookies = slices.Delete(builder.cookies, index, index+1)
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Cookies sets all cookies for the request.
//
// Parameters:
//   - cookies: The cookies to set
//
// Returns:
//   - *RequestBuilder: The builder instance for method chaining
func (builder *RequestBuilder) Cookies(cookies ...*http.Cookie) *RequestBuilder {
	if builder.err != nil {
		return builder
	}
	builder.cookies = make([]*http.Cookie, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}
		builder.cookies = append(builder.cookies, cookie)
	}
	return builder
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// build constructs the HTTP request.
//
// Parameters:
//   - ctx: The context for the request
//
// Returns:
//   - *http.Request: The constructed HTTP request
//   - error: An error if the request cannot be built
func (builder *RequestBuilder) build(ctx context.Context) (*http.Request, error) {
	if stringx.IsBlank(builder.method) {
		return nil, errors.New("method is blank")
	}
	if builder.uri == nil {
		return nil, errors.New("url is nil")
	}
	query := builder.uri.Query()
	for name, values := range builder.query() {
		if query.Has(name) {
			query.Del(name)
		}
		for _, value := range values {
			query.Add(name, value)
		}
	}
	builder.uri.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, builder.method, builder.uri.String(), builder.body)
	if err != nil {
		return nil, err
	}
	for key, values := range builder.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	for _, cookie := range builder.cookies {
		req.AddCookie(cookie)
	}
	builder.req = req
	return req, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Build constructs and returns the HTTP request.
//
// Parameters:
//   - ctx: The context for the request
//
// Returns:
//   - *http.Request: The constructed HTTP request
//   - error: An error if the request cannot be built
func (builder *RequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	if builder.err != nil {
		return nil, builder.err
	}
	return builder.build(ctx)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Execute builds and sends the HTTP request using the provided client.
//
// Parameters:
//   - ctx: The context for the request
//   - cli: The HTTP client to use
//
// Returns:
//   - *ResponseHelper: A helper for processing the response
func (builder *RequestBuilder) Execute(ctx context.Context, cli *http.Client) *ResponseHelper {
	if builder.err != nil {
		return NewResponseHelper(nil, builder.err)
	}
	req, err := builder.build(ctx)
	if err != nil {
		return NewResponseHelper(nil, err)
	}
	return NewResponseHelper(cli.Do(req))
}
