package outgoing

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/go-leo/gonv"
	"github.com/soyacen/gox/errorx"
	"google.golang.org/protobuf/proto"
)

type UnmarshalError struct {
	body []byte
	err  error
}

func (e UnmarshalError) Error() string {
	return "failed to unmarshal body"
}

func (e UnmarshalError) Unwrap() error {
	return e.err
}

func (e UnmarshalError) Body() []byte {
	return e.body
}

type Receiver interface {
	Response() *http.Response
	Status() string
	StatusCode() int
	Proto() string
	ProtoMajor() int
	ProtoMinor() int
	ContentLength() int64
	TransferEncoding() []string
	Headers() http.Header
	Trailers() http.Header
	Cookies() []*http.Cookie
	BytesBody() ([]byte, error)
	TextBody() (string, error)
	ObjectBody(body any, unmarshal func([]byte, any) error) error
	JSONBody(body any) error
	XMLBody(body any) error
	ProtobufBody(body proto.Message) error
	GobBody(body any) error
	WriterBody(file io.Writer) error
}

type receiver struct {
	resp *http.Response
}

func (r *receiver) Response() *http.Response {
	return r.resp
}

func (r *receiver) Status() string {
	return r.resp.Status
}

func (r *receiver) StatusCode() int {
	return r.resp.StatusCode
}

func (r *receiver) Proto() string {
	return r.resp.Proto
}

func (r *receiver) ProtoMajor() int {
	return r.resp.ProtoMajor
}

func (r *receiver) ProtoMinor() int {
	return r.resp.ProtoMinor
}

func (r *receiver) ContentLength() int64 {
	return r.resp.ContentLength
}

func (r *receiver) TransferEncoding() []string {
	return r.resp.TransferEncoding
}

func (r *receiver) Headers() http.Header {
	return r.resp.Header
}

func (r *receiver) Trailers() http.Header {
	return r.resp.Trailer
}

func (r *receiver) Cookies() []*http.Cookie {
	return r.resp.Cookies()
}

func (r *receiver) BytesBody() ([]byte, error) {
	body, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return nil, err
	}
	defer errorx.Silence(r.resp.Body.Close())
	return body, nil
}

func (r *receiver) TextBody() (string, error) {
	bodyBytes, err := r.BytesBody()
	if err != nil {
		return "", err
	}
	return gonv.BytesToString(bodyBytes), nil
}

func (r *receiver) ObjectBody(body any, unmarshal func([]byte, any) error) error {
	bodyBytes, err := r.BytesBody()
	if err != nil {
		return err
	}
	if err := unmarshal(bodyBytes, body); err != nil {
		return UnmarshalError{body: bodyBytes, err: err}
	}
	return nil
}

func (r *receiver) JSONBody(body any) error {
	return r.ObjectBody(body, json.Unmarshal)
}

func (r *receiver) XMLBody(body any) error {
	return r.ObjectBody(body, xml.Unmarshal)
}

func (r *receiver) ProtobufBody(body proto.Message) error {
	unmarshal := func(data []byte, v any) error { return proto.Unmarshal(data, v.(proto.Message)) }
	return r.ObjectBody(body, unmarshal)
}

func (r *receiver) GobBody(body any) error {
	unmarshal := func(data []byte, v any) error { return gob.NewDecoder(bytes.NewReader(data)).Decode(v) }
	return r.ObjectBody(body, unmarshal)
}

func (r *receiver) WriterBody(file io.Writer) error {
	_, err := io.Copy(file, r.resp.Body)
	return err
}
