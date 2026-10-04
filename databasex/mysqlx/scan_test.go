package mysqlx

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestCoerceInt64(t *testing.T) {
	testCases := []struct {
		name     string
		src      any
		want     int64
		wantNull bool
		wantErr  bool
	}{
		{name: "nil", src: nil, wantNull: true},
		{name: "int", src: 42, want: 42},
		{name: "int8", src: int8(-8), want: -8},
		{name: "int16", src: int16(-16), want: -16},
		{name: "int32", src: int32(-32), want: -32},
		{name: "int64", src: int64(-64), want: -64},
		{name: "uint", src: uint(7), want: 7},
		{name: "uint8", src: uint8(8), want: 8},
		{name: "uint16", src: uint16(16), want: 16},
		{name: "uint32", src: uint32(32), want: 32},
		{name: "uint64", src: uint64(math.MaxInt64), want: math.MaxInt64},
		{name: "uint64 overflow", src: uint64(math.MaxInt64) + 1, wantErr: true},
		{name: "bool true", src: true, want: 1},
		{name: "bool false", src: false, want: 0},
		{name: "string", src: "123", want: 123},
		{name: "bytes", src: []byte("-5"), want: -5},
		{name: "invalid string", src: "abc", wantErr: true},
		{name: "float unsupported", src: 1.5, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, isNull, err := coerceInt64(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("coerceInt64(%v) expected error", tc.src)
				}
				if !errors.Is(err, ErrOutOfRange) && !errors.Is(err, ErrCannotScan) {
					t.Fatalf("coerceInt64(%v) error = %v, want ErrOutOfRange or ErrCannotScan", tc.src, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("coerceInt64(%v) unexpected error: %v", tc.src, err)
			}
			if isNull != tc.wantNull {
				t.Fatalf("coerceInt64(%v) isNull = %v, want %v", tc.src, isNull, tc.wantNull)
			}
			if got != tc.want {
				t.Fatalf("coerceInt64(%v) = %d, want %d", tc.src, got, tc.want)
			}
		})
	}
}

func TestCoerceUint64(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		want    uint64
		wantErr bool
	}{
		{name: "nil", src: nil},
		{name: "uint64 max", src: uint64(math.MaxUint64), want: math.MaxUint64},
		{name: "int positive", src: int64(9), want: 9},
		{name: "int negative", src: int64(-1), wantErr: true},
		{name: "int8 negative", src: int8(-1), wantErr: true},
		{name: "bool true", src: true, want: 1},
		{name: "string", src: "18446744073709551615", want: math.MaxUint64},
		{name: "bytes", src: []byte("5"), want: 5},
		{name: "invalid", src: "x", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := coerceUint64(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("coerceUint64(%v) expected error", tc.src)
				}
				return
			}
			if err != nil {
				t.Fatalf("coerceUint64(%v) unexpected error: %v", tc.src, err)
			}
			if got != tc.want {
				t.Fatalf("coerceUint64(%v) = %d, want %d", tc.src, got, tc.want)
			}
		})
	}
}

func TestCoerceFloat64(t *testing.T) {
	if got, _, err := coerceFloat64(float32(1.5)); err != nil || got != 1.5 {
		t.Fatalf("coerceFloat64(float32) = %v, %v", got, err)
	}
	if got, _, err := coerceFloat64([]byte("2.25")); err != nil || got != 2.25 {
		t.Fatalf("coerceFloat64([]byte) = %v, %v", got, err)
	}
	if got, _, err := coerceFloat64(int64(3)); err != nil || got != 3 {
		t.Fatalf("coerceFloat64(int64) = %v, %v", got, err)
	}
	if _, _, err := coerceFloat64("nope"); err == nil {
		t.Fatal("coerceFloat64(string) expected error")
	}
}

func TestCoerceBool(t *testing.T) {
	testCases := []struct {
		name    string
		src     any
		want    bool
		wantErr bool
	}{
		{name: "nil", src: nil},
		{name: "bool", src: true, want: true},
		{name: "zero int", src: int64(0), want: false},
		{name: "non zero int", src: int64(2), want: true},
		{name: "negative", src: int64(-1), want: true},
		{name: "string one", src: "1", want: true},
		{name: "string false", src: []byte("false"), want: false},
		{name: "float", src: 0.5, want: true},
		{name: "invalid", src: "yes", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := coerceBool(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("coerceBool(%v) expected error", tc.src)
				}
				return
			}
			if err != nil {
				t.Fatalf("coerceBool(%v) unexpected error: %v", tc.src, err)
			}
			if got != tc.want {
				t.Fatalf("coerceBool(%v) = %v, want %v", tc.src, got, tc.want)
			}
		})
	}
}

func TestCoerceStringAndBytes(t *testing.T) {
	if v, isNull, err := coerceString("a"); err != nil || isNull || v != "a" {
		t.Fatalf("coerceString(string) = %q, %v, %v", v, isNull, err)
	}
	if v, isNull, err := coerceString([]byte("a")); err != nil || isNull || v != "a" {
		t.Fatalf("coerceString([]byte) = %q, %v, %v", v, isNull, err)
	}
	if _, isNull, err := coerceString(nil); err != nil || !isNull {
		t.Fatalf("coerceString(nil) = isNull %v, %v", isNull, err)
	}
	if _, _, err := coerceString(1); err == nil {
		t.Fatal("coerceString(int) expected error")
	}

	src := []byte("xyz")
	got, isNull, err := coerceBytes(src)
	if err != nil || isNull {
		t.Fatalf("coerceBytes([]byte) = %v, %v", isNull, err)
	}
	src[0] = 'X'
	if string(got) != "xyz" {
		t.Fatalf("coerceBytes did not copy the input: %q", got)
	}
	if _, _, err := coerceBytes(3); err == nil {
		t.Fatal("coerceBytes(int) expected error")
	}
}

func TestCloneBytes(t *testing.T) {
	if cloneBytes(nil) != nil {
		t.Fatal("cloneBytes(nil) should be nil")
	}
	src := []byte{1, 2}
	dst := cloneBytes(src)
	src[0] = 9
	if dst[0] != 1 {
		t.Fatalf("cloneBytes shares memory: %v", dst)
	}
}

func TestIsZeroDate(t *testing.T) {
	for _, s := range []string{"0000-00-00", "0000-00-00 00:00:00", "0000-00-00 00:00:00.000000"} {
		if !isZeroDate(s) {
			t.Fatalf("isZeroDate(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"2024-01-01", "0000-01-01", ""} {
		if isZeroDate(s) {
			t.Fatalf("isZeroDate(%q) = true, want false", s)
		}
	}
}

func TestCoerceTime(t *testing.T) {
	want := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	got, isNull, err := coerceTime("2024-01-02 03:04:05", dateTimeLayouts...)
	if err != nil || isNull || !got.Equal(want) {
		t.Fatalf("coerceTime(string) = %v, %v, %v", got, isNull, err)
	}

	got, isNull, err = coerceTime([]byte("2024-01-02 03:04:05.123456"), dateTimeLayouts...)
	if err != nil || isNull || got.Nanosecond() != 123456000 {
		t.Fatalf("coerceTime([]byte with fraction) = %v, %v, %v", got, isNull, err)
	}

	loc := time.FixedZone("CST", 8*3600)
	given := time.Date(2024, 1, 2, 3, 4, 5, 0, loc)
	got, _, err = coerceTime(given, dateTimeLayouts...)
	if err != nil || !got.Equal(given) || got.Location() != loc {
		t.Fatalf("coerceTime(time.Time) = %v, %v", got, err)
	}

	if _, isNull, err = coerceTime("0000-00-00 00:00:00", dateTimeLayouts...); err != nil || !isNull {
		t.Fatalf("coerceTime(zero date) = %v, %v", isNull, err)
	}

	if _, _, err = coerceTime("not a time", dateTimeLayouts...); err == nil {
		t.Fatal("coerceTime(bad string) expected error")
	}
}

func TestSignedAndUnsignedRange(t *testing.T) {
	if lo, hi := signedRange[int8](); lo != math.MinInt8 || hi != math.MaxInt8 {
		t.Fatalf("signedRange[int8]() = %d, %d", lo, hi)
	}
	if lo, hi := signedRange[int16](); lo != math.MinInt16 || hi != math.MaxInt16 {
		t.Fatalf("signedRange[int16]() = %d, %d", lo, hi)
	}
	if lo, hi := signedRange[int32](); lo != math.MinInt32 || hi != math.MaxInt32 {
		t.Fatalf("signedRange[int32]() = %d, %d", lo, hi)
	}
	if lo, hi := signedRange[int64](); lo != math.MinInt64 || hi != math.MaxInt64 {
		t.Fatalf("signedRange[int64]() = %d, %d", lo, hi)
	}
	if lo, hi := unsignedRange[uint8](); lo != 0 || hi != math.MaxUint8 {
		t.Fatalf("unsignedRange[uint8]() = %d, %d", lo, hi)
	}
	if lo, hi := unsignedRange[uint16](); lo != 0 || hi != math.MaxUint16 {
		t.Fatalf("unsignedRange[uint16]() = %d, %d", lo, hi)
	}
	if lo, hi := unsignedRange[uint32](); lo != 0 || hi != math.MaxUint32 {
		t.Fatalf("unsignedRange[uint32]() = %d, %d", lo, hi)
	}
	if lo, hi := unsignedRange[uint64](); lo != 0 || hi != math.MaxUint64 {
		t.Fatalf("unsignedRange[uint64]() = %d, %d", lo, hi)
	}
}

func TestIsJSONNull(t *testing.T) {
	for _, data := range []string{"null", " null ", "\nnull"} {
		if !isJSONNull([]byte(data)) {
			t.Fatalf("isJSONNull(%q) = false, want true", data)
		}
	}
	for _, data := range []string{"nullx", "0", `"null"`, ""} {
		if isJSONNull([]byte(data)) {
			t.Fatalf("isJSONNull(%q) = true, want false", data)
		}
	}
}

func TestCannotScanWrapsSentinel(t *testing.T) {
	err := cannotScan(struct{}{}, (*Text)(nil))
	if !errors.Is(err, ErrCannotScan) {
		t.Fatalf("cannotScan error = %v, want ErrCannotScan", err)
	}
}

func TestCoerceNumericTypeBranches(t *testing.T) {
	uintSources := []any{uint(1), uint8(2), uint16(3), uint32(4), uint64(5), int(6), int8(7), int16(8), int32(9), int64(10)}
	for _, src := range uintSources {
		got, isNull, err := coerceUint64(src)
		if err != nil || isNull || got == 0 {
			t.Fatalf("coerceUint64(%T) = %d, %v, %v", src, got, isNull, err)
		}
	}

	floatSources := []any{
		float64(1), float32(2), int(3), int8(4), int16(5), int32(6), int64(7),
		uint(8), uint8(9), uint16(10), uint32(11), uint64(12), "13", []byte("14"),
	}
	for _, src := range floatSources {
		got, isNull, err := coerceFloat64(src)
		if err != nil || isNull || got == 0 {
			t.Fatalf("coerceFloat64(%T) = %v, %v, %v", src, got, isNull, err)
		}
	}

	boolSources := []any{
		uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
		int(1), int8(1), int16(1), int32(1), int64(1),
		float32(1), float64(1), "true", []byte("1"),
	}
	for _, src := range boolSources {
		got, isNull, err := coerceBool(src)
		if err != nil || isNull || !got {
			t.Fatalf("coerceBool(%T) = %v, %v, %v", src, got, isNull, err)
		}
	}
}
