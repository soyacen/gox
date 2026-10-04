package mysqlx

import (
	"errors"
	"testing"
)

func TestYearScan(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		want    int16
		wantErr bool
	}{
		{name: "zero", src: int64(0), want: 0},
		{name: "min", src: int64(yearMin), want: yearMin},
		{name: "max", src: int64(yearMax), want: yearMax},
		{name: "string", src: "2024", want: 2024},
		{name: "bytes", src: []byte("1999"), want: 1999},
		{name: "too small", src: int64(yearMin - 1), wantErr: true},
		{name: "too large", src: int64(yearMax + 1), wantErr: true},
		{name: "negative", src: int64(-1), wantErr: true},
		{name: "invalid", src: "abcd", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var v Year
			err := v.Scan(tc.src)
			if tc.wantErr {
				if !errors.Is(err, ErrInvalidYear) && !errors.Is(err, ErrCannotScan) {
					t.Fatalf("Year.Scan(%v) error = %v, want ErrInvalidYear or ErrCannotScan", tc.src, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Year.Scan(%v) error: %v", tc.src, err)
			}
			if !v.Valid || v.Int16 != tc.want {
				t.Fatalf("Year.Scan(%v) = %#v, want %d", tc.src, v, tc.want)
			}
		})
	}

	var v Year
	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Year.Scan(nil) = %#v, %v", v, err)
	}
}

func TestYearValueAndJSON(t *testing.T) {
	v := Year{Int16: 2024, Valid: true}

	got, err := v.Value()
	if err != nil || got != int64(2024) {
		t.Fatalf("Value() = %#v, %v", got, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != "2024" {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed Year
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.Int16 != 2024 || !parsed.Valid {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte("1800")); !errors.Is(err, ErrInvalidYear) {
		t.Fatalf("UnmarshalJSON(1800) error = %v, want ErrInvalidYear", err)
	}

	outOfRange := Year{Int16: 1800, Valid: true}
	if _, err := outOfRange.Value(); !errors.Is(err, ErrInvalidYear) {
		t.Fatalf("Value() error = %v, want ErrInvalidYear", err)
	}
	if _, err := outOfRange.MarshalJSON(); !errors.Is(err, ErrInvalidYear) {
		t.Fatalf("MarshalJSON() error = %v, want ErrInvalidYear", err)
	}
}
