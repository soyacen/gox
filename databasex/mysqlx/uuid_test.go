package mysqlx

import (
	"bytes"
	"errors"
	"testing"
	"uuid"
)

// testUUID is a fixed identifier used across the UUID tests.
const testUUID = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"

func TestUUIDConstructors(t *testing.T) {
	parsed := MustParseUUID(testUUID)
	if !parsed.Valid {
		t.Fatal("MustParseUUID returned an invalid value")
	}
	if got := parsed.UUID.String(); got != testUUID {
		t.Fatalf("MustParseUUID().String() = %q, want %q", got, testUUID)
	}

	v4 := NewUUID()
	if !v4.Valid {
		t.Fatal("NewUUID returned an invalid value")
	}
	if version := v4.UUID[6] >> 4; version != 4 {
		t.Fatalf("NewUUID version = %d, want 4", version)
	}

	if other := NewUUIDV4(); !other.Valid || other.UUID == v4.UUID {
		t.Fatalf("NewUUIDV4() = %v, want a distinct valid UUID", other)
	}

	v7 := NewUUIDV7()
	if !v7.Valid {
		t.Fatal("NewUUIDV7 returned an invalid value")
	}
	if version := v7.UUID[6] >> 4; version != 7 {
		t.Fatalf("NewUUIDV7 version = %d, want 7", version)
	}
	if variant := v7.UUID[8] >> 6; variant != 0b10 {
		t.Fatalf("NewUUIDV7 variant = %#b, want %#b", variant, 0b10)
	}
}

func TestUUIDV7SortsInOrder(t *testing.T) {
	first := NewUUIDV7()
	second := NewUUIDV7()

	if first.UUID.Compare(second.UUID) >= 0 {
		t.Fatalf("NewUUIDV7 is not increasing: %v then %v", first.UUID, second.UUID)
	}
}

func TestParseUUIDForms(t *testing.T) {
	testCases := []struct {
		name string
		in   string
	}{
		{name: "canonical", in: testUUID},
		{name: "uppercase", in: "F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6"},
		{name: "bare hex", in: "f81d4fae7dec11d0a76500a0c91e6bf6"},
		{name: "braces", in: "{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}"},
		{name: "urn", in: "urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseUUID(tc.in)
			if err != nil {
				t.Fatalf("ParseUUID(%q) error: %v", tc.in, err)
			}
			if !got.Valid || got.String() != testUUID {
				t.Fatalf("ParseUUID(%q) = %#v, want %s", tc.in, got, testUUID)
			}
		})
	}
}

func TestParseUUIDErrors(t *testing.T) {
	for _, in := range []string{"", "not-a-uuid", "f81d4fae-7dec-11d0-a765-00a0c91e6bf", "f81d4fae-7dec-11d0-a765-00a0c91e6bfg"} {
		if _, err := ParseUUID(in); !errors.Is(err, ErrInvalidUUID) {
			t.Fatalf("ParseUUID(%q) error = %v, want ErrInvalidUUID", in, err)
		}
	}
}

func TestMustParseUUIDPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MustParseUUID did not panic on an invalid UUID")
		}
	}()

	MustParseUUID("not-a-uuid")
}

func TestUUIDScan(t *testing.T) {
	var packed [16]byte
	packedUUID := MustParseUUID(testUUID)
	copy(packed[:], packedUUID.UUID[:])

	testCases := []struct {
		name string
		src  any
	}{
		{name: "string", src: testUUID},
		{name: "text bytes", src: []byte(testUUID)},
		{name: "bare hex", src: []byte("f81d4fae7dec11d0a76500a0c91e6bf6")},
		{name: "packed bytes", src: MustParseUUID(testUUID).Bin()},
		{name: "array", src: packed},
		{name: "uuid", src: uuid.MustParse(testUUID)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v UUID
			if err := v.Scan(tc.src); err != nil {
				t.Fatalf("Scan(%T) error: %v", tc.src, err)
			}
			if !v.Valid || v.String() != testUUID {
				t.Fatalf("Scan(%T) = %#v, want %s", tc.src, v, testUUID)
			}
		})
	}

	var v UUID
	if err := v.Scan("not-a-uuid"); !errors.Is(err, ErrInvalidUUID) {
		t.Fatalf("Scan(bad string) error = %v, want ErrInvalidUUID", err)
	}

	if err := v.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Scan(nil) = %#v, %v", v, err)
	}
}

func TestUUIDValue(t *testing.T) {
	v := MustParseUUID(testUUID)

	got, err := v.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if got != testUUID {
		t.Fatalf("Value() = %#v, want %s", got, testUUID)
	}

	nullValue, err := (UUID{}).Value()
	if err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %#v, %v", nullValue, err)
	}
}

func TestUUIDJSON(t *testing.T) {
	v := MustParseUUID(testUUID)

	jsonBytes, err := v.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}
	if want := `"` + testUUID + `"`; string(jsonBytes) != want {
		t.Fatalf("MarshalJSON() = %s, want %s", jsonBytes, want)
	}

	var parsed UUID
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || !parsed.Valid || parsed.String() != testUUID {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"not-a-uuid"`)); !errors.Is(err, ErrInvalidUUID) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrInvalidUUID", err)
	}

	if err := parsed.UnmarshalJSON([]byte("1")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("null")); err != nil || parsed.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", parsed, err)
	}

	nullJSON, err := (UUID{}).MarshalJSON()
	if err != nil || string(nullJSON) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", nullJSON, err)
	}
}

func TestUUIDHelpers(t *testing.T) {
	v := MustParseUUID(testUUID)

	if got := v.String(); got != testUUID {
		t.Fatalf("String() = %q, want %q", got, testUUID)
	}

	bytesArray := v.Bytes()
	if !bytes.Equal(bytesArray[:], v.UUID[:]) {
		t.Fatalf("Bytes() = %v, want %v", bytesArray, v.UUID)
	}

	bin := v.Bin()
	if !bytes.Equal(bin, v.UUID[:]) {
		t.Fatalf("Bin() = %v, want %v", bin, v.UUID)
	}

	// Bin must allocate a fresh slice so callers cannot mutate the value.
	bin[0] = 'X'
	if v.UUID[0] == 'X' {
		t.Fatal("Bin() exposed the underlying array")
	}

	var null UUID
	if got := null.String(); got != "" {
		t.Fatalf("String() on NULL = %q, want an empty string", got)
	}
	if got := null.Bytes(); got != ([16]byte{}) {
		t.Fatalf("Bytes() on NULL = %v, want the zero array", got)
	}
	if got := null.Bin(); got != nil {
		t.Fatalf("Bin() on NULL = %v, want nil", got)
	}
}
