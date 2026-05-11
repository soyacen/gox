package stringx

import (
	"strconv"
	"strings"
)

// Builder wraps strings.Builder. The zero value is ready to use.
//
// Builder embeds strings.Builder by value so that copying a Builder
// triggers the underlying copyCheck on the next mutating call, matching
// the semantics of strings.Builder itself.
type Builder struct {
	b strings.Builder
}

// String returns the accumulated string.
//
// Returns:
//   - string: the accumulated string.
func (b *Builder) String() string {
	return b.b.String()
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
//
// Returns:
//   - int: the number of accumulated bytes.
func (b *Builder) Len() int {
	return b.b.Len()
}

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
//
// Returns:
//   - int: the capacity of the underlying byte slice.
func (b *Builder) Cap() int {
	return b.b.Cap()
}

// Reset resets the Builder to be empty.
func (b *Builder) Reset() {
	b.b.Reset()
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
//
// Parameters:
//   - n: the number of bytes to grow.
func (b *Builder) Grow(n int) {
	b.b.Grow(n)
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
//
// Parameters:
//   - p: the bytes to write.
//
// Returns:
//   - int: the number of bytes written.
//   - error: always nil.
func (b *Builder) Write(p []byte) (int, error) {
	return b.b.Write(p)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
//
// Parameters:
//   - c: the byte to write.
//
// Returns:
//   - error: always nil.
func (b *Builder) WriteByte(c byte) error {
	return b.b.WriteByte(c)
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
//
// Parameters:
//   - r: the rune to write.
//
// Returns:
//   - int: the length of the rune.
//   - error: always nil.
func (b *Builder) WriteRune(r rune) (int, error) {
	return b.b.WriteRune(r)
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
//
// Parameters:
//   - s: the string to write.
//
// Returns:
//   - int: the length of the string.
//   - error: always nil.
func (b *Builder) WriteString(s string) (int, error) {
	return b.b.WriteString(s)
}

// WriteInt appends the string form of the integer i.
//
// Parameters:
//   - i: the integer to write.
//   - base: the base to use.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteInt(i int64, base int) error {
	_, err := b.b.Write(strconv.AppendInt(nil, i, base))
	return err
}

// WriteUint appends the string form of the unsigned integer i.
//
// Parameters:
//   - i: the unsigned integer to write.
//   - base: the base to use.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteUint(i uint64, base int) error {
	_, err := b.b.Write(strconv.AppendUint(nil, i, base))
	return err
}

// WriteBool appends "true" or "false".
//
// Parameters:
//   - bl: the boolean to write.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteBool(bl bool) error {
	_, err := b.b.Write(strconv.AppendBool(nil, bl))
	return err
}

// WriteFloat appends the string form of the floating-point number f.
//
// Parameters:
//   - f: the float to write.
//   - fmt: the format.
//   - prec: the precision.
//   - bitSize: the bit size.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteFloat(f float64, fmt byte, prec, bitSize int) error {
	_, err := b.b.Write(strconv.AppendFloat(nil, f, fmt, prec, bitSize))
	return err
}

// WriteQuote appends a double-quoted Go string literal representing s.
//
// Parameters:
//   - s: the string to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuote(s string) error {
	_, err := b.b.Write(strconv.AppendQuote(nil, s))
	return err
}

// WriteQuoteRune appends a single-quoted Go character literal representing the rune.
//
// Parameters:
//   - r: the rune to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuoteRune(r rune) error {
	_, err := b.b.Write(strconv.AppendQuoteRune(nil, r))
	return err
}

// WriteQuoteRuneToASCII appends a single-quoted Go character literal representing the rune.
//
// Parameters:
//   - r: the rune to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuoteRuneToASCII(r rune) error {
	_, err := b.b.Write(strconv.AppendQuoteRuneToASCII(nil, r))
	return err
}

// WriteQuoteRuneToGraphic appends a single-quoted Go character literal representing the rune.
//
// Parameters:
//   - r: the rune to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuoteRuneToGraphic(r rune) error {
	_, err := b.b.Write(strconv.AppendQuoteRuneToGraphic(nil, r))
	return err
}

// WriteQuoteToASCII appends a double-quoted Go string literal representing s.
//
// Parameters:
//   - s: the string to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuoteToASCII(s string) error {
	_, err := b.b.Write(strconv.AppendQuoteToASCII(nil, s))
	return err
}

// WriteQuoteToGraphic appends a double-quoted Go string literal representing s.
//
// Parameters:
//   - s: the string to quote.
//
// Returns:
//   - error: any error encountered.
func (b *Builder) WriteQuoteToGraphic(s string) error {
	_, err := b.b.Write(strconv.AppendQuoteToGraphic(nil, s))
	return err
}

// NewBuilder creates a new Builder.
//
// Returns:
//   - *Builder: a new Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// NewBuilderBuilder creates a new Builder seeded with the contents of the
// provided strings.Builder. The source builder is not retained; only its
// current content is copied so that subsequent writes to either builder do
// not affect the other.
//
// Parameters:
//   - sb: the strings.Builder to copy content from. nil is treated as empty.
//
// Returns:
//   - *Builder: a new Builder.
func NewBuilderBuilder(sb *strings.Builder) *Builder {
	nb := &Builder{}
	if sb != nil {
		_, _ = nb.b.WriteString(sb.String())
	}
	return nb
}
