package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// maxMySQLTime is the largest MySQL TIME value, 838:59:59.
const maxMySQLTime = 838*time.Hour + 59*time.Minute + 59*time.Second

// maxTimeDigits is the number of fractional second digits MySQL stores.
const maxTimeDigits = 6

// Time represents a MySQL TIME value. MySQL TIME is a signed duration between
// -838:59:59 and 838:59:59, so it is represented as a time.Duration.
type Time struct {
	// Duration is the signed duration.
	Duration time.Duration
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Time.
func (dst *Time) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = Time{}
		return nil
	case time.Duration:
		if v > maxMySQLTime || v < -maxMySQLTime {
			return fmt.Errorf("%w: %s does not fit in %s", ErrOutOfRange, v, TypeTime)
		}
		*dst = Time{Duration: v, Valid: true}
		return nil
	case string:
		return dst.scanString(v)
	case []byte:
		return dst.scanString(string(v))
	default:
		return cannotScan(src, (*Time)(nil))
	}
}

// scanString parses and stores a MySQL TIME literal.
//
// Parameters:
//   - s: the TIME literal, for example "12:34:56.789".
//
// Returns:
//   - error: non-nil when s is not a valid TIME literal.
func (dst *Time) scanString(s string) error {
	d, err := parseMySQLTime(s)
	if err != nil {
		return err
	}
	*dst = Time{Duration: d, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface. The value is
// encoded as a MySQL TIME literal so that it round-trips exactly.
//
// Returns:
//   - driver.Value: the TIME literal, or nil when the value is NULL.
//   - error: non-nil when the duration is outside the MySQL TIME range.
func (src Time) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	if src.Duration > maxMySQLTime || src.Duration < -maxMySQLTime {
		return nil, fmt.Errorf("%w: %s does not fit in %s", ErrOutOfRange, src.Duration, TypeTime)
	}
	return formatMySQLTime(src.Duration), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The value is
// encoded as a MySQL TIME literal string.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the duration is outside the MySQL TIME range.
func (src Time) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	if src.Duration > maxMySQLTime || src.Duration < -maxMySQLTime {
		return nil, fmt.Errorf("%w: %s does not fit in %s", ErrOutOfRange, src.Duration, TypeTime)
	}
	return json.Marshal(formatMySQLTime(src.Duration))
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The value
// is expected to be a MySQL TIME literal string.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a TIME literal string or null.
func (dst *Time) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Time{}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return dst.scanString(s)
}

// parseMySQLTime parses a MySQL TIME literal into a time.Duration.
//
// Parameters:
//   - s: the TIME literal, for example "-12:34:56.789".
//
// Returns:
//   - time.Duration: the parsed signed duration.
//   - error: non-nil when s is not a valid TIME literal.
func parseMySQLTime(s string) (time.Duration, error) {
	raw := s
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("%w: empty literal", ErrInvalidTime)
	}

	neg := false
	switch s[0] {
	case '-':
		neg = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	frac := ""
	hasFrac := false
	if i := strings.IndexByte(s, '.'); i >= 0 {
		frac = s[i+1:]
		hasFrac = true
		s = s[:i]
	}

	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("%w: %q", ErrInvalidTime, raw)
	}

	hours, err := parseTimeField(parts[0], 0, 838)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidTime, raw)
	}
	minutes, err := parseTimeField(parts[1], 0, 59)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidTime, raw)
	}
	seconds, err := parseTimeField(parts[2], 0, 59)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidTime, raw)
	}

	d := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second

	if hasFrac {
		micros, err := parseMicroseconds(frac)
		if err != nil {
			return 0, fmt.Errorf("%w: %q", ErrInvalidTime, raw)
		}
		d += time.Duration(micros) * time.Microsecond
	}

	if d > maxMySQLTime {
		return 0, fmt.Errorf("%w: %q is outside the %s range", ErrInvalidTime, raw, TypeTime)
	}
	if neg {
		d = -d
	}

	return d, nil
}

// parseTimeField parses a single mandatory-digit time field and validates its
// range.
//
// Parameters:
//   - s: the field text.
//   - min: the inclusive minimum value.
//   - max: the inclusive maximum value.
//
// Returns:
//   - int: the parsed value.
//   - error: non-nil when s is not a decimal number within the range.
func parseTimeField(s string, min int, max int) (int, error) {
	if s == "" {
		return 0, ErrInvalidTime
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if v < min || v > max {
		return 0, ErrInvalidTime
	}
	return v, nil
}

// parseMicroseconds parses the fractional seconds of a TIME literal.
//
// Parameters:
//   - frac: the digits following the decimal point.
//
// Returns:
//   - int: the microsecond value.
//   - error: non-nil when frac is empty or not made of decimal digits.
func parseMicroseconds(frac string) (int, error) {
	if frac == "" {
		return 0, ErrInvalidTime
	}
	if len(frac) > maxTimeDigits {
		frac = frac[:maxTimeDigits]
	}
	for len(frac) < maxTimeDigits {
		frac += "0"
	}
	return strconv.Atoi(frac)
}

// formatMySQLTime renders a time.Duration as a MySQL TIME literal. The
// fractional part is omitted when it is zero.
//
// Parameters:
//   - d: the duration to format.
//
// Returns:
//   - string: the TIME literal.
func formatMySQLTime(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}

	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	micros := (d - seconds*time.Second) / time.Microsecond

	if micros == 0 {
		return fmt.Sprintf("%s%02d:%02d:%02d", sign, hours, minutes, seconds)
	}
	return fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hours, minutes, seconds, micros)
}
