package mysqlx

import (
	"bytes"
	"errors"
	"math"
	"testing"
)

func TestBitScan(t *testing.T) {
	var v Bit
	if err := v.Scan([]byte{0x01, 0x02}); err != nil || !v.Valid || v.Uint64 != 258 {
		t.Fatalf("Bit.Scan([]byte) = %#v, %v", v, err)
	}

	if err := v.Scan("258"); err != nil || v.Uint64 != 258 {
		t.Fatalf("Bit.Scan(string) = %#v, %v", v, err)
	}

	if err := v.Scan(int64(5)); err != nil || v.Uint64 != 5 {
		t.Fatalf("Bit.Scan(int64) = %#v, %v", v, err)
	}

	if err := v.Scan(uint64(math.MaxUint64)); err != nil || v.Uint64 != math.MaxUint64 {
		t.Fatalf("Bit.Scan(uint64) = %#v, %v", v, err)
	}

	if err := v.Scan([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("Bit.Scan(9 bytes) error = %v, want ErrOutOfRange", err)
	}

	if err := v.Scan("abc"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Bit.Scan(abc) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Bit.Scan(nil) = %#v, %v", v, err)
	}
}

func TestBitValue(t *testing.T) {
	v := Bit{Uint64: 258, Valid: true}
	got, err := v.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if !bytes.Equal(got.([]byte), []byte{0x01, 0x02}) {
		t.Fatalf("Value() = %v, want [1 2]", got)
	}

	zero, err := (Bit{Uint64: 0, Valid: true}).Value()
	if err != nil || !bytes.Equal(zero.([]byte), []byte{0x00}) {
		t.Fatalf("Value() for zero = %v, %v", zero, err)
	}

	max, err := (Bit{Uint64: math.MaxUint64, Valid: true}).Value()
	if err != nil || len(max.([]byte)) != 8 {
		t.Fatalf("Value() for max = %v, %v", max, err)
	}

	if nullValue, err := (Bit{}).Value(); err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %v, %v", nullValue, err)
	}
}

func TestBitJSON(t *testing.T) {
	var v Bit
	if err := v.UnmarshalJSON([]byte("258")); err != nil || v.Uint64 != 258 {
		t.Fatalf("UnmarshalJSON(number) = %#v, %v", v, err)
	}
	if err := v.UnmarshalJSON([]byte(`"258"`)); err != nil || v.Uint64 != 258 {
		t.Fatalf("UnmarshalJSON(string) = %#v, %v", v, err)
	}
	if err := v.UnmarshalJSON([]byte(`"x"`)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrCannotScan", err)
	}
	if err := v.UnmarshalJSON([]byte("null")); err != nil || v.Valid {
		t.Fatalf("UnmarshalJSON(null) = %#v, %v", v, err)
	}
}
