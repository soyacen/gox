package mysqlx

import (
	"errors"
	"math"
	"testing"
)

func TestIntegerRangeChecks(t *testing.T) {
	testCases := []struct {
		name    string
		new     func() nullable
		src     any
		wantErr bool
	}{
		{name: "TinyInt max ok", new: func() nullable { return &TinyInt{} }, src: int64(math.MaxInt8)},
		{name: "TinyInt max+1", new: func() nullable { return &TinyInt{} }, src: int64(math.MaxInt8) + 1, wantErr: true},
		{name: "TinyInt min-1", new: func() nullable { return &TinyInt{} }, src: int64(math.MinInt8) - 1, wantErr: true},
		{name: "TinyIntUnsigned ok", new: func() nullable { return &TinyIntUnsigned{} }, src: uint64(math.MaxUint8)},
		{name: "TinyIntUnsigned negative", new: func() nullable { return &TinyIntUnsigned{} }, src: int64(-1), wantErr: true},
		{name: "TinyIntUnsigned max+1", new: func() nullable { return &TinyIntUnsigned{} }, src: uint64(math.MaxUint8) + 1, wantErr: true},
		{name: "SmallInt max ok", new: func() nullable { return &SmallInt{} }, src: int64(math.MaxInt16)},
		{name: "SmallInt max+1", new: func() nullable { return &SmallInt{} }, src: int64(math.MaxInt16) + 1, wantErr: true},
		{name: "SmallIntUnsigned max+1", new: func() nullable { return &SmallIntUnsigned{} }, src: uint64(math.MaxUint16) + 1, wantErr: true},
		{name: "MediumInt max ok", new: func() nullable { return &MediumInt{} }, src: int64(mediumIntMax)},
		{name: "MediumInt max+1", new: func() nullable { return &MediumInt{} }, src: int64(mediumIntMax) + 1, wantErr: true},
		{name: "MediumInt min-1", new: func() nullable { return &MediumInt{} }, src: int64(mediumIntMin) - 1, wantErr: true},
		{name: "MediumIntUnsigned max ok", new: func() nullable { return &MediumIntUnsigned{} }, src: uint64(mediumIntUnsignedMax)},
		{name: "MediumIntUnsigned max+1", new: func() nullable { return &MediumIntUnsigned{} }, src: uint64(mediumIntUnsignedMax) + 1, wantErr: true},
		{name: "Int max ok", new: func() nullable { return &Int{} }, src: int64(math.MaxInt32)},
		{name: "Int max+1", new: func() nullable { return &Int{} }, src: int64(math.MaxInt32) + 1, wantErr: true},
		{name: "IntUnsigned max+1", new: func() nullable { return &IntUnsigned{} }, src: uint64(math.MaxUint32) + 1, wantErr: true},
		{name: "BigInt max+1", new: func() nullable { return &BigInt{} }, src: uint64(math.MaxInt64) + 1, wantErr: true},
		{name: "BigInt from string", new: func() nullable { return &BigInt{} }, src: []byte("-9223372036854775808")},
		{name: "unsupported source", new: func() nullable { return &BigInt{} }, src: struct{}{}, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.new()
			err := v.Scan(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Scan(%v) expected error", tc.src)
				}
				if !errors.Is(err, ErrOutOfRange) && !errors.Is(err, ErrCannotScan) {
					t.Fatalf("Scan(%v) error = %v, want ErrOutOfRange or ErrCannotScan", tc.src, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan(%v) unexpected error: %v", tc.src, err)
			}
		})
	}
}

func TestMediumIntConstructedOutOfRange(t *testing.T) {
	v := MediumInt{Int32: mediumIntMax + 1, Valid: true}
	if _, err := v.Value(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("Value() error = %v, want ErrOutOfRange", err)
	}
	if _, err := v.MarshalJSON(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("MarshalJSON() error = %v, want ErrOutOfRange", err)
	}

	u := MediumIntUnsigned{Uint32: mediumIntUnsignedMax + 1, Valid: true}
	if _, err := u.Value(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("unsigned Value() error = %v, want ErrOutOfRange", err)
	}
	if _, err := u.MarshalJSON(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("unsigned MarshalJSON() error = %v, want ErrOutOfRange", err)
	}
}

func TestIntegerJSONRangeChecks(t *testing.T) {
	var tiny TinyInt
	if err := tiny.UnmarshalJSON([]byte("128")); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("TinyInt.UnmarshalJSON(128) error = %v, want ErrOutOfRange", err)
	}

	var u TinyIntUnsigned
	if err := u.UnmarshalJSON([]byte("-1")); err == nil {
		t.Fatal("TinyIntUnsigned.UnmarshalJSON(-1) expected error")
	}

	var medium MediumInt
	if err := medium.UnmarshalJSON([]byte("8388608")); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("MediumInt.UnmarshalJSON error = %v, want ErrOutOfRange", err)
	}

	var mediumU MediumIntUnsigned
	if err := mediumU.UnmarshalJSON([]byte("16777216")); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("MediumIntUnsigned.UnmarshalJSON error = %v, want ErrOutOfRange", err)
	}

	var big BigInt
	if err := big.UnmarshalJSON([]byte(`"5"`)); err == nil {
		t.Fatal("BigInt.UnmarshalJSON(string) expected error")
	}
}

func TestIntegerInvalidSources(t *testing.T) {
	var v Int
	if err := v.Scan("not a number"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Int.Scan(bad string) error = %v, want ErrCannotScan", err)
	}
	if err := v.Scan(1.5); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Int.Scan(float) error = %v, want ErrCannotScan", err)
	}
}

func TestBoolScanSources(t *testing.T) {
	testCases := []struct {
		src  any
		want bool
	}{
		{src: true, want: true},
		{src: int64(0), want: false},
		{src: int64(3), want: true},
		{src: []byte("1"), want: true},
		{src: "false", want: false},
	}

	for _, tc := range testCases {
		var v Bool
		if err := v.Scan(tc.src); err != nil {
			t.Fatalf("Bool.Scan(%v) error: %v", tc.src, err)
		}
		if !v.Valid || v.Bool != tc.want {
			t.Fatalf("Bool.Scan(%v) = %#v, want %v", tc.src, v, tc.want)
		}
	}

	var v Bool
	if err := v.Scan("maybe"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Bool.Scan(maybe) error = %v, want ErrCannotScan", err)
	}
}
