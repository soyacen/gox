package httpx

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// ResponseHelper provides convenient methods for processing HTTP responses.
type ResponseHelper struct {
	err        error
	resp       *http.Response
	statusCode int
	headers    http.Header
	trailers   http.Header
	cookies    []*http.Cookie
	bodyBytes  []byte
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// NewResponseHelper creates a new ResponseHelper from an HTTP response.
// Reads and stores the response body for later processing.
//
// Parameters:
//   - resp: The HTTP response
//   - err: Any error from the request
//
// Returns:
//   - *ResponseHelper: A helper for processing the response
func NewResponseHelper(resp *http.Response, err error) *ResponseHelper {
	respHelper := &ResponseHelper{resp: resp, err: err}
	if err != nil {
		return respHelper
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respHelper.err = err
		return respHelper
	}
	err = resp.Body.Close()
	if err != nil {
		respHelper.err = err
		return respHelper
	}

	respHelper.resp = resp
	respHelper.statusCode = resp.StatusCode
	respHelper.headers = resp.Header
	respHelper.trailers = resp.Trailer
	respHelper.cookies = resp.Cookies()
	respHelper.bodyBytes = body
	return respHelper
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Err returns any error that occurred during the request or response processing.
//
// Returns:
//   - error: The error, or nil if no error occurred
func (helper *ResponseHelper) Err() error {
	return helper.err
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// StatusCode returns the HTTP status code of the response.
//
// Returns:
//   - int: The HTTP status code
//   - error: An error if the request failed
func (helper *ResponseHelper) StatusCode() (int, error) {
	if helper.err != nil {
		return 0, helper.err
	}
	return helper.statusCode, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Headers returns the HTTP response headers.
//
// Returns:
//   - http.Header: The response headers
//   - error: An error if the request failed
func (helper *ResponseHelper) Headers() (http.Header, error) {
	if helper.err != nil {
		return nil, helper.err
	}
	return helper.headers, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Trailer returns the HTTP response trailers.
//
// Returns:
//   - http.Header: The response trailers
//   - error: An error if the request failed
func (helper *ResponseHelper) Trailer() (http.Header, error) {
	if helper.err != nil {
		return nil, helper.err
	}
	return helper.trailers, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Cookies returns the HTTP response cookies.
//
// Returns:
//   - []*http.Cookie: The response cookies
//   - error: An error if the request failed
func (helper *ResponseHelper) Cookies() ([]*http.Cookie, error) {
	if helper.err != nil {
		return nil, helper.err
	}
	return helper.cookies, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// Body returns the response body as a read-closer.
//
// Returns:
//   - io.ReadCloser: The response body
//   - error: An error if the request failed
func (helper *ResponseHelper) Body() (io.ReadCloser, error) {
	if helper.err != nil {
		return nil, helper.err
	}
	return io.NopCloser(bytes.NewReader(helper.bodyBytes)), nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// BytesBody returns the response body as a byte slice.
//
// Returns:
//   - []byte: The response body bytes
//   - error: An error if the request failed
func (helper *ResponseHelper) BytesBody() ([]byte, error) {
	if helper.err != nil {
		return nil, helper.err
	}
	return helper.bodyBytes, nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// TextBody returns the response body as a string.
//
// Returns:
//   - string: The response body text
//   - error: An error if the request failed
func (helper *ResponseHelper) TextBody() (string, error) {
	if helper.err != nil {
		return "", helper.err
	}
	return string(helper.bodyBytes), nil
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// ObjectBody unmarshals the response body into the provided object.
//
// Parameters:
//   - body: The object to unmarshal into
//   - unmarshal: The unmarshal function
//
// Returns:
//   - error: An error if the request failed or unmarshal fails
func (helper *ResponseHelper) ObjectBody(body any, unmarshal func([]byte, any) error) error {
	if helper.err != nil {
		return helper.err
	}
	err := unmarshal(helper.bodyBytes, body)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal body, body is %s, %w", helper.bodyBytes, err)
	}
	return err
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// JSONBody unmarshals the response body as JSON into the provided object.
//
// Parameters:
//   - body: The object to unmarshal into
//
// Returns:
//   - error: An error if the request failed or JSON unmarshal fails
func (helper *ResponseHelper) JSONBody(body any) error {
	if helper.err != nil {
		return helper.err
	}
	return helper.ObjectBody(body, json.Unmarshal)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// XMLBody unmarshals the response body as XML into the provided object.
//
// Parameters:
//   - body: The object to unmarshal into
//
// Returns:
//   - error: An error if the request failed or XML unmarshal fails
func (helper *ResponseHelper) XMLBody(body any) error {
	if helper.err != nil {
		return helper.err
	}
	return helper.ObjectBody(body, xml.Unmarshal)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// ProtobufBody unmarshals the response body as protobuf into the provided message.
//
// Parameters:
//   - body: The protobuf message to unmarshal into
//
// Returns:
//   - error: An error if the request failed or protobuf unmarshal fails
func (helper *ResponseHelper) ProtobufBody(body proto.Message) error {
	if helper.err != nil {
		return helper.err
	}
	unmarshal := func(data []byte, v any) error {
		m := v.(proto.Message)
		return proto.Unmarshal(data, m)
	}
	return helper.ObjectBody(body, unmarshal)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// GobBody unmarshals the response body as Gob into the provided message.
//
// Parameters:
//   - body: The object to unmarshal into
//
// Returns:
//   - error: An error if the request failed or Gob unmarshal fails
func (helper *ResponseHelper) GobBody(body proto.Message) error {
	if helper.err != nil {
		return helper.err
	}
	unmarshal := func(data []byte, v any) error {
		return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
	}
	return helper.ObjectBody(body, unmarshal)
}

// Deprecated: Do not use. Use github.com/soyacen/netx/httpx/outgoing instead.
// FileBody writes the response body to the provided writer.
//
// Parameters:
//   - file: The writer to write the response body to
//
// Returns:
//   - error: An error if the request failed or write fails
func (helper *ResponseHelper) FileBody(file io.Writer) error {
	if helper.err != nil {
		return helper.err
	}
	_, err := file.Write(helper.bodyBytes)
	return err
}
