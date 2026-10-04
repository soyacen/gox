package mysqlx

import (
	"errors"
	"testing"
)

func TestTextScan(t *testing.T) {
	var v Text
	if err := v.Scan([]byte("hello")); err != nil || !v.Valid || v.String != "hello" {
		t.Fatalf("Text.Scan([]byte) = %#v, %v", v, err)
	}

	if err := v.Scan(""); err != nil || !v.Valid || v.String != "" {
		t.Fatalf("Text.Scan(empty) = %#v, %v, empty string must not be NULL", v, err)
	}

	if err := v.Scan("NULL"); err != nil || !v.Valid || v.String != "NULL" {
		t.Fatalf("Text.Scan(NULL) = %#v, %v, the literal %q must not be NULL", v, err, "NULL")
	}

	if err := v.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Text.Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid || v.String != "" {
		t.Fatalf("Text.Scan(nil) = %#v, %v", v, err)
	}
}

func TestTextJSON(t *testing.T) {
	v := Text{String: "a\"b", Valid: true}

	jsonBytes, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}

	var parsed Text
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.String != v.String {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte("1")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}
}

func TestEnumScan(t *testing.T) {
	var v Enum
	if err := v.Scan("pending"); err != nil || !v.Valid || v.String != "pending" {
		t.Fatalf("Enum.Scan() = %#v, %v", v, err)
	}

	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Enum.Scan(nil) = %#v, %v", v, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != "null" {
		t.Fatalf("Enum.MarshalJSON() = %s, %v", jsonBytes, err)
	}

	valid := Enum{String: "pending", Valid: true}
	if err := valid.UnmarshalJSON([]byte(`"pending"`)); err != nil || valid.String != "pending" {
		t.Fatalf("Enum.UnmarshalJSON() = %#v, %v", valid, err)
	}
}
