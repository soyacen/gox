package mysqlx

import (
	"errors"
	"testing"
	"time"
)

func TestDateScan(t *testing.T) {
	var d Date
	if err := d.Scan([]byte("2024-01-02")); err != nil {
		t.Fatalf("Date.Scan([]byte) error: %v", err)
	}
	if want := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC); !d.Valid || !d.Time.Equal(want) {
		t.Fatalf("Date.Scan([]byte) = %#v, want %v", d, want)
	}

	loc := time.FixedZone("CST", 8*3600)
	given := time.Date(2024, 1, 2, 0, 0, 0, 0, loc)
	if err := d.Scan(given); err != nil || !d.Time.Equal(given) || d.Time.Location() != loc {
		t.Fatalf("Date.Scan(time.Time) = %#v, %v", d, err)
	}

	if err := d.Scan("0000-00-00"); err != nil || d.Valid {
		t.Fatalf("Date.Scan(zero date) = %#v, %v", d, err)
	}

	if err := d.Scan("not-a-date"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Date.Scan(bad) error = %v, want ErrCannotScan", err)
	}

	if err := d.Scan("2024-01-02 03:04:05"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Date.Scan(datetime literal) error = %v, want ErrCannotScan", err)
	}

	if err := d.Scan(nil); err != nil || d.Valid {
		t.Fatalf("Date.Scan(nil) = %#v, %v", d, err)
	}
}

func TestDateValueAndJSON(t *testing.T) {
	d := Date{Time: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), Valid: true}

	v, err := d.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if _, ok := v.(time.Time); !ok {
		t.Fatalf("Value() = %T, want time.Time", v)
	}

	got, err := d.MarshalJSON()
	if err != nil || string(got) != `"2024-01-02"` {
		t.Fatalf("MarshalJSON() = %s, %v", got, err)
	}

	var parsed Date
	if err := parsed.UnmarshalJSON([]byte(`"2024-01-02"`)); err != nil || !parsed.Valid {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"bad"`)); err == nil {
		t.Fatal("UnmarshalJSON(bad) expected error")
	}

	if err := parsed.UnmarshalJSON([]byte("2024")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}
}
