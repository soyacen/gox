package mysqlx

import (
	"bytes"
	"errors"
	"testing"
)

func TestBlobScanCopiesInput(t *testing.T) {
	src := []byte("payload")

	var v Blob
	if err := v.Scan(src); err != nil || !v.Valid {
		t.Fatalf("Blob.Scan() = %#v, %v", v, err)
	}

	src[0] = 'X'
	if string(v.Bytes) != "payload" {
		t.Fatalf("Blob.Scan() retained driver memory: %q", v.Bytes)
	}

	if err := v.Scan(nil); err != nil || v.Valid || v.Bytes != nil {
		t.Fatalf("Blob.Scan(nil) = %#v, %v", v, err)
	}

	if err := v.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Blob.Scan(int64) error = %v, want ErrCannotScan", err)
	}
}

func TestBlobValueAndJSON(t *testing.T) {
	v := Blob{Bytes: []byte{1, 2, 3}, Valid: true}

	got, err := v.Value()
	if err != nil || !bytes.Equal(got.([]byte), []byte{1, 2, 3}) {
		t.Fatalf("Value() = %#v, %v", got, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `"AQID"` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed Blob
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || !bytes.Equal(parsed.Bytes, v.Bytes) {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"!!!"`)); err == nil {
		t.Fatal("UnmarshalJSON(invalid base64) expected error")
	}

	if nullValue, err := (Blob{}).Value(); err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %#v, %v", nullValue, err)
	}
}
