package mysqlx

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

func TestNewDecimalFromString(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{in: "123.45", want: "123.45"},
		{in: "123.4500", want: "123.4500"},
		{in: "-0.001", want: "-0.001"},
		{in: "0", want: "0"},
		{in: "-0", want: "0"},
		{in: "+7", want: "7"},
		{in: ".5", want: "0.5"},
		{in: "5.", want: "5"},
		{in: "1e3", want: "1000"},
		{in: "1.5e3", want: "1500"},
		{in: "1.5E-3", want: "0.0015"},
		{in: " 42.5 ", want: "42.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			d, err := NewDecimalFromString(tc.in)
			if err != nil {
				t.Fatalf("NewDecimalFromString(%q) error: %v", tc.in, err)
			}
			if got := d.String(); got != tc.want {
				t.Fatalf("NewDecimalFromString(%q).String() = %q, want %q", tc.in, got, tc.want)
			}
			if !d.Valid {
				t.Fatalf("NewDecimalFromString(%q) is not Valid", tc.in)
			}
		})
	}
}

func TestNewDecimalFromStringErrors(t *testing.T) {
	for _, in := range []string{"", " ", ".", "abc", "1.2.3", "1e", "1e99999999999", "--1", "1-1"} {
		if _, err := NewDecimalFromString(in); !errors.Is(err, ErrInvalidDecimal) {
			t.Fatalf("NewDecimalFromString(%q) error = %v, want ErrInvalidDecimal", in, err)
		}
	}
}

func TestDecimalScan(t *testing.T) {
	var d Decimal
	if err := d.Scan([]byte("123.45")); err != nil || d.String() != "123.45" {
		t.Fatalf("Decimal.Scan([]byte) = %#v, %v", d, err)
	}
	if err := d.Scan(int64(-7)); err != nil || d.String() != "-7" {
		t.Fatalf("Decimal.Scan(int64) = %#v, %v", d, err)
	}
	if err := d.Scan(uint64(math.MaxUint64)); err != nil || d.String() != "18446744073709551615" {
		t.Fatalf("Decimal.Scan(uint64) = %#v, %v", d, err)
	}
	if err := d.Scan(0.5); err != nil || d.String() != "0.5" {
		t.Fatalf("Decimal.Scan(float64) = %#v, %v", d, err)
	}
	if err := d.Scan(math.NaN()); err == nil {
		t.Fatal("Decimal.Scan(NaN) expected error")
	}
	if err := d.Scan(struct{}{}); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Decimal.Scan(struct) error = %v, want ErrCannotScan", err)
	}
	if err := d.Scan(nil); err != nil || d.Valid {
		t.Fatalf("Decimal.Scan(nil) = %#v, %v", d, err)
	}
}

func TestDecimalValueAndFloat(t *testing.T) {
	d, err := NewDecimalFromString("0.5")
	if err != nil {
		t.Fatal(err)
	}

	v, err := d.Value()
	if err != nil || v != "0.5" {
		t.Fatalf("Value() = %#v, %v", v, err)
	}

	f, err := d.Float64()
	if err != nil || f != 0.5 {
		t.Fatalf("Float64() = %v, %v", f, err)
	}

	if _, err := (Decimal{}).Float64(); !errors.Is(err, ErrInvalidDecimal) {
		t.Fatalf("Float64() on NULL error = %v, want ErrInvalidDecimal", err)
	}

	nullValue, err := (Decimal{}).Value()
	if err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %#v, %v", nullValue, err)
	}
}

func TestDecimalJSON(t *testing.T) {
	var d Decimal
	if err := d.UnmarshalJSON([]byte("1.25")); err != nil || d.String() != "1.25" {
		t.Fatalf("UnmarshalJSON(number) = %#v, %v", d, err)
	}
	if err := d.UnmarshalJSON([]byte(`"2.5e2"`)); err != nil || d.String() != "250" {
		t.Fatalf("UnmarshalJSON(string) = %#v, %v", d, err)
	}
	if err := d.UnmarshalJSON([]byte(`"x"`)); !errors.Is(err, ErrInvalidDecimal) {
		t.Fatalf("UnmarshalJSON(bad string) error = %v, want ErrInvalidDecimal", err)
	}
	if err := d.UnmarshalJSON([]byte("{]")); err == nil {
		t.Fatal("UnmarshalJSON(object) expected error")
	}

	got, err := d.MarshalJSON()
	if err != nil || string(got) != "250" {
		t.Fatalf("MarshalJSON() = %s, %v", got, err)
	}

	if got, err := (Decimal{}).MarshalJSON(); err != nil || string(got) != "null" {
		t.Fatalf("MarshalJSON() on NULL = %s, %v", got, err)
	}
}

func TestDecimalScientificFormatting(t *testing.T) {
	huge := Decimal{Int: big.NewInt(1), Exp: maxDecimalExp + 1, Valid: true}
	if got := huge.String(); got != "1e100001" {
		t.Fatalf("String() = %q, want 1e100001", got)
	}

	tiny := Decimal{Int: big.NewInt(-2), Exp: -maxDecimalExp - 1, Valid: true}
	if got := tiny.String(); got != "-2e-100001" {
		t.Fatalf("String() = %q, want -2e-100001", got)
	}
}

func TestDecimalNullString(t *testing.T) {
	if got := (Decimal{}).String(); got != "null" {
		t.Fatalf("String() on NULL = %q, want null", got)
	}
}

func TestDecimalScanNumericBranches(t *testing.T) {
	sources := []any{
		int(1), int8(2), int16(3), int32(4), int64(5),
		uint(6), uint8(7), uint16(8), uint32(9), uint64(10),
		float32(1.5), float64(2.5),
	}

	for _, src := range sources {
		var d Decimal
		if err := d.Scan(src); err != nil || !d.Valid {
			t.Fatalf("Decimal.Scan(%T) = %#v, %v", src, d, err)
		}
		if got := d.String(); got == "null" || got == "" {
			t.Fatalf("Decimal.Scan(%T).String() = %q", src, got)
		}
	}
}
