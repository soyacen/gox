package mysqlx

import (
	"errors"
	"testing"
	"time"
)

func TestDateTimeScan(t *testing.T) {
	want := time.Date(2024, 1, 2, 3, 4, 5, 123456000, time.UTC)

	var dt DateTime
	if err := dt.Scan("2024-01-02 03:04:05.123456"); err != nil || !dt.Valid || !dt.Time.Equal(want) {
		t.Fatalf("DateTime.Scan(with fraction) = %#v, %v", dt, err)
	}

	if err := dt.Scan([]byte("2024-01-02 03:04:05")); err != nil || dt.Time.Nanosecond() != 0 {
		t.Fatalf("DateTime.Scan(without fraction) = %#v, %v", dt, err)
	}

	if err := dt.Scan("2024-01-02T03:04:05Z"); err != nil || !dt.Valid {
		t.Fatalf("DateTime.Scan(RFC3339) = %#v, %v", dt, err)
	}

	if err := dt.Scan("0000-00-00 00:00:00"); err != nil || dt.Valid {
		t.Fatalf("DateTime.Scan(zero date) = %#v, %v", dt, err)
	}

	if err := dt.Scan("nope"); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("DateTime.Scan(bad) error = %v, want ErrCannotScan", err)
	}

	if err := dt.Scan(nil); err != nil || dt.Valid {
		t.Fatalf("DateTime.Scan(nil) = %#v, %v", dt, err)
	}
}

func TestTimestampScan(t *testing.T) {
	var ts Timestamp
	if err := ts.Scan([]byte("2024-01-02 03:04:05")); err != nil {
		t.Fatalf("Timestamp.Scan error: %v", err)
	}
	if want := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC); !ts.Valid || !ts.Time.Equal(want) {
		t.Fatalf("Timestamp.Scan = %#v, want %v", ts, want)
	}

	if err := ts.Scan("0000-00-00 00:00:00"); err != nil || ts.Valid {
		t.Fatalf("Timestamp.Scan(zero date) = %#v, %v", ts, err)
	}

	if err := ts.Scan(nil); err != nil || ts.Valid {
		t.Fatalf("Timestamp.Scan(nil) = %#v, %v", ts, err)
	}
}

func TestDateTimeValueAndJSON(t *testing.T) {
	dt := DateTime{Time: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), Valid: true}

	got, err := dt.MarshalJSON()
	if err != nil || string(got) != `"2024-01-02T03:04:05Z"` {
		t.Fatalf("DateTime.MarshalJSON() = %s, %v", got, err)
	}

	var parsed DateTime
	if err := parsed.UnmarshalJSON(got); err != nil || !parsed.Valid || !parsed.Time.Equal(dt.Time) {
		t.Fatalf("DateTime.UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"2024-01-02 03:04:05"`)); err != nil || !parsed.Valid {
		t.Fatalf("DateTime.UnmarshalJSON(plain) = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"bad"`)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("DateTime.UnmarshalJSON(bad) error = %v, want ErrCannotScan", err)
	}

	ts := Timestamp{Time: dt.Time, Valid: true}
	tsJSON, err := ts.MarshalJSON()
	if err != nil || string(tsJSON) != `"2024-01-02T03:04:05Z"` {
		t.Fatalf("Timestamp.MarshalJSON() = %s, %v", tsJSON, err)
	}

	if v, err := ts.Value(); err != nil || v == nil {
		t.Fatalf("Timestamp.Value() = %v, %v", v, err)
	}
}
