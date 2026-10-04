package mysqlx

import (
	"bytes"
	"errors"
	"testing"
)

func TestJSONScan(t *testing.T) {
	var v JSON
	if err := v.Scan([]byte(`{"a":1}`)); err != nil || !v.Valid {
		t.Fatalf("JSON.Scan(object) = %#v, %v", v, err)
	}
	if string(v.Bytes) != `{"a":1}` {
		t.Fatalf("JSON.Scan(object) bytes = %s", v.Bytes)
	}

	if err := v.Scan("[1,2,3]"); err != nil || !v.Valid {
		t.Fatalf("JSON.Scan(string array) = %#v, %v", v, err)
	}

	if err := v.Scan("1"); err != nil || !v.Valid {
		t.Fatalf("JSON.Scan(number) = %#v, %v", v, err)
	}

	if err := v.Scan([]byte("{oops")); !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("JSON.Scan(invalid) error = %v, want ErrInvalidJSON", err)
	}

	if err := v.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("JSON.Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid || v.Bytes != nil {
		t.Fatalf("JSON.Scan(nil) = %#v, %v", v, err)
	}
}

func TestJSONScanCopiesInput(t *testing.T) {
	src := []byte(`{"a":1}`)

	var v JSON
	if err := v.Scan(src); err != nil {
		t.Fatalf("JSON.Scan() error: %v", err)
	}

	src[2] = 'X'
	if string(v.Bytes) != `{"a":1}` {
		t.Fatalf("JSON.Scan() retained driver memory: %s", v.Bytes)
	}
}

func TestJSONValueAndMarshal(t *testing.T) {
	v := JSON{Bytes: []byte(`{"a":1}`), Valid: true}

	got, err := v.Value()
	if err != nil || !bytes.Equal(got.([]byte), v.Bytes) {
		t.Fatalf("Value() = %#v, %v", got, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `{"a":1}` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	broken := JSON{Bytes: []byte("{oops"), Valid: true}
	if _, err := broken.Value(); !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("Value() error = %v, want ErrInvalidJSON", err)
	}
	if _, err := broken.MarshalJSON(); !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("MarshalJSON() error = %v, want ErrInvalidJSON", err)
	}

	if nullValue, err := (JSON{}).Value(); err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %#v, %v", nullValue, err)
	}
	if nullJSON, err := (JSON{}).MarshalJSON(); err != nil || string(nullJSON) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", nullJSON, err)
	}
}

func TestJSONUnmarshal(t *testing.T) {
	var v JSON
	if err := v.UnmarshalJSON([]byte(`{"a":1}`)); err != nil || !v.Valid || string(v.Bytes) != `{"a":1}` {
		t.Fatalf("UnmarshalJSON(object) = %#v, %v", v, err)
	}

	if err := v.UnmarshalJSON([]byte("null")); err != nil || v.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", v, err)
	}

	if err := v.UnmarshalJSON([]byte("{oops")); !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("UnmarshalJSON(invalid) error = %v, want ErrInvalidJSON", err)
	}
}
