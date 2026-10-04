package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// maxDecimalExp bounds the absolute decimal exponent accepted by the parser and
// formatted positionally. Larger exponents are formatted in scientific
// notation.
const maxDecimalExp = 100000

// Decimal represents a MySQL DECIMAL value. DECIMAL, NUMERIC, DEC and FIXED are
// synonyms.
//
// The exact value is Int × 10^Exp. Unlike float64 the representation is lossless
// and preserves the scale reported by MySQL, so "123.4500" round-trips
// unchanged.
type Decimal struct {
	// Int is the unscaled integer value.
	Int *big.Int
	// Exp is the base ten exponent applied to Int.
	Exp int32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// NewDecimalFromString parses a MySQL DECIMAL literal. Both plain and
// scientific notation are accepted.
//
// Parameters:
//   - s: the decimal literal, for example "123.45" or "1.5e3".
//
// Returns:
//   - Decimal: the parsed value with Valid set to true.
//   - error: non-nil when s is not a valid decimal literal.
func NewDecimalFromString(s string) (Decimal, error) {
	n, exp, err := parseDecimal(s)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal{Int: n, Exp: exp, Valid: true}, nil
}

// String returns the decimal literal of the value, or "null" when the value is
// NULL.
//
// Returns:
//   - string: the decimal literal.
func (d Decimal) String() string {
	if !d.Valid {
		return "null"
	}
	return formatDecimal(d.Int, d.Exp)
}

// Float64 converts the value to the nearest float64. Precision may be lost for
// values that are not exactly representable in binary floating point.
//
// Returns:
//   - float64: the converted value.
//   - error: non-nil when the value is NULL, not a number, or out of the float64
//     range.
func (d Decimal) Float64() (float64, error) {
	if !d.Valid {
		return 0, fmt.Errorf("%w: cannot convert NULL to float64", ErrInvalidDecimal)
	}

	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidDecimal, err)
	}
	return f, nil
}

// Scan implements the database/sql.Scanner interface. Floating point sources
// are converted through their shortest decimal representation.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be parsed as a decimal.
func (dst *Decimal) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*dst = Decimal{}
		return nil
	case string:
		return dst.scanString(v)
	case []byte:
		return dst.scanString(string(v))
	case int:
		*dst = Decimal{Int: big.NewInt(int64(v)), Valid: true}
		return nil
	case int8:
		*dst = Decimal{Int: big.NewInt(int64(v)), Valid: true}
		return nil
	case int16:
		*dst = Decimal{Int: big.NewInt(int64(v)), Valid: true}
		return nil
	case int32:
		*dst = Decimal{Int: big.NewInt(int64(v)), Valid: true}
		return nil
	case int64:
		*dst = Decimal{Int: big.NewInt(v), Valid: true}
		return nil
	case uint:
		return dst.scanUint(uint64(v))
	case uint8:
		return dst.scanUint(uint64(v))
	case uint16:
		return dst.scanUint(uint64(v))
	case uint32:
		return dst.scanUint(uint64(v))
	case uint64:
		return dst.scanUint(v)
	case float32:
		return dst.scanString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		return dst.scanString(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return cannotScan(src, (*Decimal)(nil))
	}
}

// scanUint stores an unsigned integer source.
//
// Parameters:
//   - v: the unsigned value.
//
// Returns:
//   - error: always nil.
func (dst *Decimal) scanUint(v uint64) error {
	*dst = Decimal{Int: new(big.Int).SetUint64(v), Valid: true}
	return nil
}

// scanString parses and stores a decimal literal.
//
// Parameters:
//   - s: the decimal literal.
//
// Returns:
//   - error: non-nil when s is not a valid decimal literal.
func (dst *Decimal) scanString(s string) error {
	n, exp, err := parseDecimal(s)
	if err != nil {
		return err
	}
	*dst = Decimal{Int: n, Exp: exp, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface. The value is
// encoded as a decimal string to avoid precision loss.
//
// Returns:
//   - driver.Value: the decimal literal, or nil when the value is NULL.
//   - error: always nil.
func (src Decimal) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.String(), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The value is
// encoded as a JSON number, or null when the value is NULL.
//
// Returns:
//   - []byte: the JSON number.
//   - error: always nil.
func (src Decimal) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return []byte(src.String()), nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. Both JSON
// numbers and JSON strings containing a decimal literal are accepted.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a decimal number, a decimal string, or
//     null.
func (dst *Decimal) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Decimal{}
		return nil
	}

	raw := strings.TrimSpace(string(data))
	if len(raw) > 1 && raw[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		raw = s
	}

	n, exp, err := parseDecimal(raw)
	if err != nil {
		return err
	}
	*dst = Decimal{Int: n, Exp: exp, Valid: true}
	return nil
}

// parseDecimal parses a decimal literal into an unscaled integer and a base ten
// exponent.
//
// Parameters:
//   - s: the decimal literal.
//
// Returns:
//   - *big.Int: the unscaled integer value.
//   - int32: the base ten exponent.
//   - error: non-nil when s is not a valid decimal literal.
func parseDecimal(s string) (*big.Int, int32, error) {
	raw := s
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, 0, fmt.Errorf("%w: empty literal", ErrInvalidDecimal)
	}

	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		neg = true
		s = s[1:]
	}
	if s == "" {
		return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
	}

	exp := int64(0)
	if i := strings.IndexAny(s, "eE"); i >= 0 {
		e := s[i+1:]
		s = s[:i]
		parsed, err := strconv.ParseInt(e, 10, 32)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
		}
		exp = parsed
	}

	intPart := s
	fracPart := ""
	if i := strings.IndexByte(s, '.'); i >= 0 {
		intPart, fracPart = s[:i], s[i+1:]
		if strings.IndexByte(fracPart, '.') >= 0 {
			return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
		}
	}

	digits := intPart + fracPart
	if digits == "" {
		return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
		}
	}

	exp -= int64(len(fracPart))
	if exp > maxDecimalExp || exp < -maxDecimalExp {
		return nil, 0, fmt.Errorf("%w: exponent out of range in %q", ErrInvalidDecimal, raw)
	}

	n, ok := new(big.Int).SetString(digits, 10)
	if !ok {
		return nil, 0, fmt.Errorf("%w: %q", ErrInvalidDecimal, raw)
	}
	if neg {
		n.Neg(n)
	}

	return n, int32(exp), nil
}

// formatDecimal renders an unscaled integer and exponent as a decimal literal.
//
// Parameters:
//   - n: the unscaled integer value.
//   - exp: the base ten exponent.
//
// Returns:
//   - string: the decimal literal.
func formatDecimal(n *big.Int, exp int32) string {
	if n == nil || n.Sign() == 0 {
		return "0"
	}

	neg := n.Sign() < 0
	digits := new(big.Int).Abs(n).String()

	switch {
	case exp == 0:
	case exp > 0:
		if int64(exp) > maxDecimalExp {
			return scientificDecimal(digits, neg, exp)
		}
		digits += strings.Repeat("0", int(exp))
	default:
		if -int64(exp) > maxDecimalExp {
			return scientificDecimal(digits, neg, exp)
		}
		pointPos := len(digits) + int(exp)
		if pointPos > 0 {
			digits = digits[:pointPos] + "." + digits[pointPos:]
		} else {
			digits = "0." + strings.Repeat("0", -pointPos) + digits
		}
	}

	if neg {
		return "-" + digits
	}
	return digits
}

// scientificDecimal renders a decimal value in scientific notation.
//
// Parameters:
//   - digits: the magnitude digits without a sign.
//   - neg: whether the value is negative.
//   - exp: the base ten exponent.
//
// Returns:
//   - string: the scientific notation literal.
func scientificDecimal(digits string, neg bool, exp int32) string {
	out := digits + "e" + strconv.FormatInt(int64(exp), 10)
	if neg {
		return "-" + out
	}
	return out
}
