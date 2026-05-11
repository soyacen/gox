package timex

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"testing"
	"time"
)

var loc = time.UTC

func mustTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	v, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		t.Fatalf("parse %q: %v", value, err)
	}
	return v
}

func TestMax(t *testing.T) {
	a := mustTime(t, time.RFC3339, "2025-01-01T00:00:00Z")
	b := mustTime(t, time.RFC3339, "2025-06-01T00:00:00Z")
	c := mustTime(t, time.RFC3339, "2024-12-31T23:59:59Z")

	if got := Max(a); !got.Equal(a) {
		t.Errorf("Max(a) = %v, want %v", got, a)
	}
	if got := Max(a, b, c); !got.Equal(b) {
		t.Errorf("Max(a,b,c) = %v, want %v", got, b)
	}
	if got := Max(b, a, c); !got.Equal(b) {
		t.Errorf("Max(b,a,c) = %v, want %v", got, b)
	}
}

func TestMin(t *testing.T) {
	a := mustTime(t, time.RFC3339, "2025-01-01T00:00:00Z")
	b := mustTime(t, time.RFC3339, "2025-06-01T00:00:00Z")
	c := mustTime(t, time.RFC3339, "2024-12-31T23:59:59Z")

	if got := Min(a); !got.Equal(a) {
		t.Errorf("Min(a) = %v, want %v", got, a)
	}
	if got := Min(a, b, c); !got.Equal(c) {
		t.Errorf("Min(a,b,c) = %v, want %v", got, c)
	}
}

func TestDate(t *testing.T) {
	in := time.Date(2025, 5, 20, 13, 45, 30, 1234, loc)
	want := time.Date(2025, 5, 20, 0, 0, 0, 0, loc)
	got := Date(in)
	if !got.Equal(want) {
		t.Errorf("Date(%v) = %v, want %v", in, got, want)
	}
}

func TestToday(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		got := Today()
		now := time.Now()
		want := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("Today() = %v, want %v", got, want)
		}
	})

	t.Run("with arg", func(t *testing.T) {
		in := time.Date(2025, 5, 20, 9, 0, 0, 0, loc)
		want := time.Date(2025, 5, 20, 0, 0, 0, 0, loc)
		if got := Today(in); !got.Equal(want) {
			t.Errorf("Today(in) = %v, want %v", got, want)
		}
	})

	t.Run("zero arg falls back to now", func(t *testing.T) {
		got := Today(time.Time{})
		now := time.Now()
		want := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("Today(zero) = %v, want %v", got, want)
		}
	})
}

func TestTomorrow(t *testing.T) {
	in := time.Date(2025, 5, 20, 23, 59, 59, 0, loc)
	want := time.Date(2025, 5, 21, 0, 0, 0, 0, loc)
	if got := Tomorrow(in); !got.Equal(want) {
		t.Errorf("Tomorrow(in) = %v, want %v", got, want)
	}
}

func TestYesterday(t *testing.T) {
	in := time.Date(2025, 5, 20, 0, 0, 1, 0, loc)
	want := time.Date(2025, 5, 19, 0, 0, 0, 0, loc)
	if got := Yesterday(in); !got.Equal(want) {
		t.Errorf("Yesterday(in) = %v, want %v", got, want)
	}
}

func TestHour(t *testing.T) {
	in := time.Date(2025, 5, 20, 13, 45, 30, 1234, loc)
	want := time.Date(2025, 5, 20, 13, 0, 0, 0, loc)
	if got := Hour(in); !got.Equal(want) {
		t.Errorf("Hour(in) = %v, want %v", got, want)
	}
}

func TestThisHour(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		got := ThisHour()
		now := time.Now()
		want := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("ThisHour() = %v, want %v", got, want)
		}
	})

	t.Run("with arg", func(t *testing.T) {
		in := time.Date(2025, 5, 20, 13, 45, 30, 0, loc)
		want := time.Date(2025, 5, 20, 13, 0, 0, 0, loc)
		if got := ThisHour(in); !got.Equal(want) {
			t.Errorf("ThisHour(in) = %v, want %v", got, want)
		}
	})
}

func TestLastHour(t *testing.T) {
	in := time.Date(2025, 5, 20, 13, 45, 30, 0, loc)
	want := time.Date(2025, 5, 20, 12, 0, 0, 0, loc)
	if got := LastHour(in); !got.Equal(want) {
		t.Errorf("LastHour(in) = %v, want %v", got, want)
	}
}

func TestNextHour(t *testing.T) {
	in := time.Date(2025, 5, 20, 13, 45, 30, 0, loc)
	want := time.Date(2025, 5, 20, 14, 0, 0, 0, loc)
	if got := NextHour(in); !got.Equal(want) {
		t.Errorf("NextHour(in) = %v, want %v", got, want)
	}
}

func TestThisMonth(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		got := ThisMonth()
		now := time.Now()
		want := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("ThisMonth() = %v, want %v", got, want)
		}
	})

	t.Run("with arg", func(t *testing.T) {
		in := time.Date(2025, 5, 20, 13, 45, 0, 0, loc)
		want := time.Date(2025, 5, 1, 0, 0, 0, 0, loc)
		if got := ThisMonth(in); !got.Equal(want) {
			t.Errorf("ThisMonth(in) = %v, want %v", got, want)
		}
	})
}

func TestLastMonth(t *testing.T) {
	in := time.Date(2025, 5, 20, 0, 0, 0, 0, loc)
	want := time.Date(2025, 4, 1, 0, 0, 0, 0, loc)
	if got := LastMonth(in); !got.Equal(want) {
		t.Errorf("LastMonth(in) = %v, want %v", got, want)
	}
}

func TestNextMonth(t *testing.T) {
	in := time.Date(2025, 12, 20, 0, 0, 0, 0, loc)
	want := time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
	if got := NextMonth(in); !got.Equal(want) {
		t.Errorf("NextMonth(in) = %v, want %v", got, want)
	}
}

func TestMonthOfLastDay(t *testing.T) {
	t.Run("January 31", func(t *testing.T) {
		in := time.Date(2025, 1, 15, 0, 0, 0, 0, loc)
		want := time.Date(2025, 1, 31, 0, 0, 0, 0, loc)
		if got := MonthOfLastDay(in); !got.Equal(want) {
			t.Errorf("MonthOfLastDay(Jan) = %v, want %v", got, want)
		}
	})

	t.Run("Leap February", func(t *testing.T) {
		in := time.Date(2024, 2, 10, 0, 0, 0, 0, loc)
		want := time.Date(2024, 2, 29, 0, 0, 0, 0, loc)
		if got := MonthOfLastDay(in); !got.Equal(want) {
			t.Errorf("MonthOfLastDay(LeapFeb) = %v, want %v", got, want)
		}
	})

	t.Run("Non-leap February", func(t *testing.T) {
		in := time.Date(2025, 2, 10, 0, 0, 0, 0, loc)
		want := time.Date(2025, 2, 28, 0, 0, 0, 0, loc)
		if got := MonthOfLastDay(in); !got.Equal(want) {
			t.Errorf("MonthOfLastDay(Feb) = %v, want %v", got, want)
		}
	})
}

func TestYear(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		got := Year()
		now := time.Now()
		want := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("Year() = %v, want %v", got, want)
		}
	})

	t.Run("with arg", func(t *testing.T) {
		in := time.Date(2025, 7, 15, 13, 45, 0, 0, loc)
		want := time.Date(2025, 1, 1, 0, 0, 0, 0, loc)
		if got := Year(in); !got.Equal(want) {
			t.Errorf("Year(in) = %v, want %v", got, want)
		}
	})

	t.Run("zero falls back to now", func(t *testing.T) {
		got := Year(time.Time{})
		now := time.Now()
		want := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		if !got.Equal(want) {
			t.Errorf("Year(zero) = %v, want %v", got, want)
		}
	})
}

func TestIsUnixZero(t *testing.T) {
	if !IsUnixZero(UnixZero) {
		t.Error("IsUnixZero(UnixZero) should be true")
	}
	if IsUnixZero(time.Time{}) {
		t.Error("IsUnixZero(zero time.Time) should be false")
	}
	if IsUnixZero(time.Date(2025, 1, 1, 0, 0, 0, 0, loc)) {
		t.Error("IsUnixZero(2025) should be false")
	}
}

func TestIsZero(t *testing.T) {
	if !IsZero(time.Time{}) {
		t.Error("IsZero(zero) should be true")
	}
	if !IsZero(UnixZero) {
		t.Error("IsZero(UnixZero) should be true")
	}
	if IsZero(time.Now()) {
		t.Error("IsZero(now) should be false")
	}
}

func TestIsNotZero(t *testing.T) {
	if IsNotZero(time.Time{}) {
		t.Error("IsNotZero(zero) should be false")
	}
	if IsNotZero(UnixZero) {
		t.Error("IsNotZero(UnixZero) should be false")
	}
	if !IsNotZero(time.Now()) {
		t.Error("IsNotZero(now) should be true")
	}
}

func TestUTCTimeMarshalUnmarshal(t *testing.T) {
	orig := UTCTime(time.Date(2025, 5, 20, 13, 45, 30, int(123*time.Millisecond), time.UTC))
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want := `"2025-05-20T13:45:30.123Z"`
	if string(data) != want {
		t.Fatalf("Marshal got %s, want %s", data, want)
	}

	var got UTCTime
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !time.Time(got).Equal(time.Time(orig)) {
		t.Fatalf("Unmarshal got %v, want %v", time.Time(got), time.Time(orig))
	}
}

func TestUTCTimeUnmarshalErrors(t *testing.T) {
	var v UTCTime
	err := v.UnmarshalJSON([]byte("12345"))
	if err == nil {
		t.Fatal("expected error for non-quoted input")
	}

	err = v.UnmarshalJSON([]byte(`"not-a-time"`))
	if err == nil {
		t.Fatal("expected parse error for invalid timestamp")
	}
}

func TestUTCTimeString(t *testing.T) {
	in := UTCTime(time.Date(2025, 5, 20, 13, 45, 30, int(456*time.Millisecond), time.UTC))
	want := "2025-05-20T13:45:30.456Z"
	if got := in.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestUTCTimeHash32(t *testing.T) {
	in := UTCTime(time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC))
	h := fnv.New32a()
	if err := in.Hash32(h); err != nil {
		t.Fatalf("Hash32: %v", err)
	}
	if h.Sum32() == 0 {
		t.Error("Hash32 sum should be non-zero")
	}

	h2 := fnv.New32a()
	if err := in.Hash32(h2); err != nil {
		t.Fatalf("Hash32 (2nd): %v", err)
	}
	if h.Sum32() != h2.Sum32() {
		t.Errorf("Hash32 differs for equal inputs: %d vs %d", h.Sum32(), h2.Sum32())
	}

	other := UTCTime(time.Date(2025, 5, 20, 0, 0, 1, 0, time.UTC))
	h3 := fnv.New32a()
	if err := other.Hash32(h3); err != nil {
		t.Fatalf("Hash32 (other): %v", err)
	}
	if h.Sum32() == h3.Sum32() {
		t.Error("Hash32 should differ for distinct inputs")
	}
}

func TestParseUTCTime(t *testing.T) {
	got, err := ParseUTCTime("2025-05-20T13:45:30.000Z")
	if err != nil {
		t.Fatalf("ParseUTCTime: %v", err)
	}
	want := time.Date(2025, 5, 20, 13, 45, 30, 0, time.UTC)
	if !time.Time(got).Equal(want) {
		t.Errorf("ParseUTCTime = %v, want %v", time.Time(got), want)
	}

	if _, err := ParseUTCTime("not-a-time"); err == nil {
		t.Error("ParseUTCTime invalid should error")
	}
}

func TestMustParseUTCTime(t *testing.T) {
	got := MustParseUTCTime("2025-05-20T13:45:30.000Z")
	want := time.Date(2025, 5, 20, 13, 45, 30, 0, time.UTC)
	if !time.Time(got).Equal(want) {
		t.Errorf("MustParseUTCTime = %v, want %v", time.Time(got), want)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseUTCTime should panic on invalid input")
		}
	}()
	MustParseUTCTime("not-a-time")
}

func TestLayoutConstants(t *testing.T) {
	in := time.Date(2025, 5, 20, 13, 45, 30, int(123456789), time.UTC)
	cases := []struct {
		name, layout, want string
	}{
		{"YearMonth", YearMonth, "2025-05"},
		{"DateMinute", DateMinute, "2025-05-20 13:45"},
		{"DateHour", DateHour, "2025-05-20 13"},
		{"DateTimeNano", DateTimeNano, "2025-05-20 13:45:30.123456789"},
		{"UTCLayout", UTCLayout, "2025-05-20T13:45:30.123Z"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := in.Format(c.layout); got != c.want {
				t.Errorf("Format(%s) = %q, want %q", c.name, got, c.want)
			}
		})
	}
	if HttpTimeLayout == "" {
		t.Error("HttpTimeLayout should not be empty")
	}
}

func TestUTCTimeRoundTripBytes(t *testing.T) {
	orig := UTCTime(time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC))
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !bytes.HasPrefix(data, []byte(`"`)) || !bytes.HasSuffix(data, []byte(`"`)) {
		t.Fatalf("marshaled value should be quoted JSON string: %s", data)
	}
}
