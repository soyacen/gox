package mysqlx

import (
	"errors"
	"math"
	"testing"
)

func TestFloatScan(t *testing.T) {
	var v Float
	if err := v.Scan(float32(1.25)); err != nil || !v.Valid || v.Float32 != 1.25 {
		t.Fatalf("Float.Scan(float32) = %#v, %v", v, err)
	}

	if err := v.Scan([]byte("2.5")); err != nil || v.Float32 != 2.5 {
		t.Fatalf("Float.Scan([]byte) = %#v, %v", v, err)
	}

	if err := v.Scan(float64(math.MaxFloat32) * 2); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("Float.Scan(overflow) error = %v, want ErrOutOfRange", err)
	}

	if err := v.Scan(struct{}{}); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Float.Scan(struct) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Float.Scan(nil) = %#v, %v", v, err)
	}
}

func TestDoubleScan(t *testing.T) {
	var v Double
	if err := v.Scan("3.5"); err != nil || !v.Valid || v.Float64 != 3.5 {
		t.Fatalf("Double.Scan(string) = %#v, %v", v, err)
	}
	if err := v.Scan(int64(4)); err != nil || v.Float64 != 4 {
		t.Fatalf("Double.Scan(int64) = %#v, %v", v, err)
	}
	if err := v.Scan("x"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Double.Scan(x) error = %v, want ErrCannotScan", err)
	}
	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Double.Scan(nil) = %#v, %v", v, err)
	}
}

func TestFloatJSON(t *testing.T) {
	var v Double
	if err := v.UnmarshalJSON([]byte("1.5")); err != nil || v.Float64 != 1.5 {
		t.Fatalf("Double.UnmarshalJSON = %#v, %v", v, err)
	}
	if err := v.UnmarshalJSON([]byte(`"1.5"`)); err == nil {
		t.Fatal("Double.UnmarshalJSON(string) expected error")
	}
	if err := v.UnmarshalJSON([]byte("null")); err != nil || v.Valid {
		t.Fatalf("Double.UnmarshalJSON(null) = %#v, %v", v, err)
	}

	var f Float
	if err := f.UnmarshalJSON([]byte("1.5")); err != nil || f.Float32 != 1.5 {
		t.Fatalf("Float.UnmarshalJSON = %#v, %v", f, err)
	}
}
