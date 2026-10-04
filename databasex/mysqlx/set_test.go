package mysqlx

import (
	"errors"
	"testing"
)

func TestSetScan(t *testing.T) {
	var v Set
	if err := v.Scan("a,b,c"); err != nil || !v.Valid {
		t.Fatalf("Set.Scan() = %#v, %v", v, err)
	}
	if len(v.Strings) != 3 || v.Strings[0] != "a" || v.Strings[2] != "c" {
		t.Fatalf("Set.Scan() members = %#v", v.Strings)
	}

	if err := v.Scan(""); err != nil || !v.Valid || v.Strings == nil || len(v.Strings) != 0 {
		t.Fatalf("Set.Scan(empty) = %#v, %v", v, err)
	}

	if err := v.Scan([]byte("solo")); err != nil || len(v.Strings) != 1 || v.Strings[0] != "solo" {
		t.Fatalf("Set.Scan([]byte) = %#v, %v", v, err)
	}

	if err := v.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Set.Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid || v.Strings != nil {
		t.Fatalf("Set.Scan(nil) = %#v, %v", v, err)
	}
}

func TestSetValue(t *testing.T) {
	v := Set{Strings: []string{"a", "b"}, Valid: true}
	got, err := v.Value()
	if err != nil || got != "a,b" {
		t.Fatalf("Value() = %#v, %v", got, err)
	}

	empty := Set{Strings: []string{}, Valid: true}
	got, err = empty.Value()
	if err != nil || got != "" {
		t.Fatalf("Value() for empty = %#v, %v", got, err)
	}

	null, err := (Set{}).Value()
	if err != nil || null != nil {
		t.Fatalf("Value() for NULL = %#v, %v", null, err)
	}
}

func TestSetJSON(t *testing.T) {
	v := Set{Strings: []string{"a", "b"}, Valid: true}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `["a","b"]` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed Set
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || len(parsed.Strings) != 2 {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	empty := Set{Strings: []string{}, Valid: true}
	emptyJSON, err := empty.MarshalJSON()
	if err != nil || string(emptyJSON) != "[]" {
		t.Fatalf("MarshalJSON() for empty = %s, %v", emptyJSON, err)
	}

	if err := parsed.UnmarshalJSON([]byte("[]")); err != nil || parsed.Strings == nil || len(parsed.Strings) != 0 {
		t.Fatalf("UnmarshalJSON([]) = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"a,b"`)); err == nil {
		t.Fatal("UnmarshalJSON(string) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("null")); err != nil || parsed.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", parsed, err)
	}
}
