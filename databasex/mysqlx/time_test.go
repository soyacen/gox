package mysqlx

import (
	"errors"
	"testing"
	"time"
)

func TestMySQLTimeParseFormat(t *testing.T) {
	testCases := []struct {
		in   string
		want string
	}{
		{in: "00:00:00", want: "00:00:00"},
		{in: "12:34:56", want: "12:34:56"},
		{in: "12:34:56.789", want: "12:34:56.789000"},
		{in: "-12:34:56", want: "-12:34:56"},
		{in: "+01:02:03", want: "01:02:03"},
		{in: "838:59:59", want: "838:59:59"},
		{in: "-838:59:59", want: "-838:59:59"},
		{in: "00:00:00.000001", want: "00:00:00.000001"},
		{in: " 01:02:03 ", want: "01:02:03"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			d, err := parseMySQLTime(tc.in)
			if err != nil {
				t.Fatalf("parseMySQLTime(%q) error: %v", tc.in, err)
			}
			if got := formatMySQLTime(d); got != tc.want {
				t.Fatalf("formatMySQLTime(parse(%q)) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestMySQLTimeParseErrors(t *testing.T) {
	for _, in := range []string{"", "12:34", "839:00:00", "12:60:00", "12:34:60", "ab:cd:ef", "12:34:56.", "-", "1:2:3:4"} {
		if _, err := parseMySQLTime(in); !errors.Is(err, ErrInvalidTime) {
			t.Fatalf("parseMySQLTime(%q) error = %v, want ErrInvalidTime", in, err)
		}
	}
}

func TestTimeScan(t *testing.T) {
	var v Time
	if err := v.Scan("12:34:56"); err != nil || !v.Valid || v.Duration != 12*time.Hour+34*time.Minute+56*time.Second {
		t.Fatalf("Time.Scan(string) = %#v, %v", v, err)
	}

	if err := v.Scan(90 * time.Minute); err != nil || v.Duration != 90*time.Minute {
		t.Fatalf("Time.Scan(duration) = %#v, %v", v, err)
	}

	if err := v.Scan(839 * time.Hour); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("Time.Scan(out of range) error = %v, want ErrOutOfRange", err)
	}

	if err := v.Scan(time.Time{}); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Time.Scan(time.Time) error = %v, want ErrCannotScan", err)
	}

	if err := v.Scan(nil); err != nil || v.Valid {
		t.Fatalf("Time.Scan(nil) = %#v, %v", v, err)
	}
}

func TestTimeValueAndJSON(t *testing.T) {
	v := Time{Duration: 12*time.Hour + 34*time.Minute + 56*time.Second, Valid: true}

	got, err := v.Value()
	if err != nil || got != "12:34:56" {
		t.Fatalf("Value() = %#v, %v", got, err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil || string(jsonBytes) != `"12:34:56"` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed Time
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.Duration != v.Duration {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte(`"bad"`)); !errors.Is(err, ErrInvalidTime) {
		t.Fatalf("UnmarshalJSON(bad) error = %v, want ErrInvalidTime", err)
	}

	if err := parsed.UnmarshalJSON([]byte("5")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}

	outOfRange := Time{Duration: maxMySQLTime + time.Second, Valid: true}
	if _, err := outOfRange.Value(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("Value() error = %v, want ErrOutOfRange", err)
	}
	if _, err := outOfRange.MarshalJSON(); !errors.Is(err, ErrOutOfRange) {
		t.Fatalf("MarshalJSON() error = %v, want ErrOutOfRange", err)
	}
}
