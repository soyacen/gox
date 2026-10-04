package mysqlx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"testing"
)

// bigEndianWKBPoint builds a big-endian WKB POINT for tests.
//
// Parameters:
//   - x: the X coordinate.
//   - y: the Y coordinate.
//
// Returns:
//   - []byte: the WKB payload.
func bigEndianWKBPoint(x float64, y float64) []byte {
	b := make([]byte, wkbPointSize)
	b[0] = wkbByteOrderBigEndian
	binary.BigEndian.PutUint32(b[1:5], wkbPointType)
	binary.BigEndian.PutUint64(b[5:13], math.Float64bits(x))
	binary.BigEndian.PutUint64(b[13:21], math.Float64bits(y))
	return b
}

func TestPointBareWKBRoundTrip(t *testing.T) {
	src := encodeWKBPoint(1.5, -2.25)

	var p Point
	if err := p.Scan(src); err != nil {
		t.Fatalf("Point.Scan(bare WKB) error: %v", err)
	}
	if !p.Valid || p.X != 1.5 || p.Y != -2.25 || p.SRID != nil {
		t.Fatalf("Point.Scan(bare WKB) = %#v", p)
	}

	got, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if !bytes.Equal(got.([]byte), src) {
		t.Fatalf("Value() = %v, want %v", got, src)
	}
}

func TestPointWithSRIDRoundTrip(t *testing.T) {
	srid := uint32(4326)
	src := prependSRID(srid, encodeWKBPoint(3, 4))

	var p Point
	if err := p.Scan(src); err != nil {
		t.Fatalf("Point.Scan() error: %v", err)
	}
	if !p.Valid || p.X != 3 || p.Y != 4 || p.SRID == nil || *p.SRID != srid {
		t.Fatalf("Point.Scan() = %#v", p)
	}

	got, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if !bytes.Equal(got.([]byte), src) {
		t.Fatalf("Value() = %v, want %v", got, src)
	}
}

func TestPointBigEndianScan(t *testing.T) {
	var p Point
	if err := p.Scan(bigEndianWKBPoint(7, 8)); err != nil {
		t.Fatalf("Point.Scan(big-endian) error: %v", err)
	}
	if p.X != 7 || p.Y != 8 {
		t.Fatalf("Point.Scan(big-endian) = %#v", p)
	}

	// Encoding normalizes to little-endian.
	got, err := p.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if !bytes.Equal(got.([]byte), encodeWKBPoint(7, 8)) {
		t.Fatalf("Value() = %v, want little-endian encoding", got)
	}
}

func TestPointScanErrors(t *testing.T) {
	var p Point
	if err := p.Scan([]byte{1, 2, 3}); !errors.Is(err, ErrInvalidGeometry) {
		t.Fatalf("Point.Scan(short) error = %v, want ErrInvalidGeometry", err)
	}

	notAPoint := encodeWKBPoint(0, 0)
	binary.LittleEndian.PutUint32(notAPoint[1:5], 2)
	if err := p.Scan(notAPoint); !errors.Is(err, ErrInvalidGeometry) {
		t.Fatalf("Point.Scan(linestring WKB) error = %v, want ErrInvalidGeometry", err)
	}

	badOrder := encodeWKBPoint(0, 0)
	badOrder[0] = 9
	if err := p.Scan(badOrder); !errors.Is(err, ErrInvalidGeometry) {
		t.Fatalf("Point.Scan(bad byte order) error = %v, want ErrInvalidGeometry", err)
	}

	if err := p.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Point.Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := p.Scan(nil); err != nil || p.Valid {
		t.Fatalf("Point.Scan(nil) = %#v, %v", p, err)
	}
}

func TestPointJSON(t *testing.T) {
	p := Point{X: 1, Y: 2, Valid: true}
	jsonBytes, err := p.MarshalJSON()
	if err != nil || string(jsonBytes) != `{"x":1,"y":2}` {
		t.Fatalf("MarshalJSON() = %s, %v", jsonBytes, err)
	}

	var parsed Point
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.X != 1 || parsed.Y != 2 || parsed.SRID != nil {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}

	srid := uint32(4326)
	withSRID := Point{X: 1, Y: 2, SRID: &srid, Valid: true}
	jsonBytes, err = withSRID.MarshalJSON()
	if err != nil || string(jsonBytes) != `{"x":1,"y":2,"srid":4326}` {
		t.Fatalf("MarshalJSON() with SRID = %s, %v", jsonBytes, err)
	}

	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || parsed.SRID == nil || *parsed.SRID != srid {
		t.Fatalf("UnmarshalJSON() with SRID = %#v, %v", parsed, err)
	}

	if err := parsed.UnmarshalJSON([]byte("5")); err == nil {
		t.Fatal("UnmarshalJSON(number) expected error")
	}
}

func TestGeometryScan(t *testing.T) {
	srid := uint32(4326)
	raw := prependSRID(srid, encodeWKBPoint(1, 2))

	var g Geometry
	if err := g.Scan(raw); err != nil || !g.Valid {
		t.Fatalf("Geometry.Scan() = %#v, %v", g, err)
	}
	if !bytes.Equal(g.Bytes, raw) {
		t.Fatalf("Geometry.Scan() bytes = %v, want %v", g.Bytes, raw)
	}
	if g.SRID == nil || *g.SRID != srid {
		t.Fatalf("Geometry.Scan() SRID = %v, want %d", g.SRID, srid)
	}

	if err := g.Scan([]byte{1, 2, 3}); err != nil || g.SRID != nil {
		t.Fatalf("Geometry.Scan(short) = %#v, %v", g, err)
	}

	if err := g.Scan([]byte{}); !errors.Is(err, ErrInvalidGeometry) {
		t.Fatalf("Geometry.Scan(empty) error = %v, want ErrInvalidGeometry", err)
	}

	if err := g.Scan(int64(1)); !errors.Is(err, ErrCannotScan) {
		t.Fatalf("Geometry.Scan(int64) error = %v, want ErrCannotScan", err)
	}

	if err := g.Scan(nil); err != nil || g.Valid {
		t.Fatalf("Geometry.Scan(nil) = %#v, %v", g, err)
	}
}

func TestGeometryValueAndJSON(t *testing.T) {
	raw := prependSRID(42, encodeWKBPoint(1, 2))
	g := Geometry{Bytes: raw, Valid: true}

	got, err := g.Value()
	if err != nil || !bytes.Equal(got.([]byte), raw) {
		t.Fatalf("Value() = %v, %v", got, err)
	}

	jsonBytes, err := g.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}

	var parsed Geometry
	if err := parsed.UnmarshalJSON(jsonBytes); err != nil || !bytes.Equal(parsed.Bytes, raw) {
		t.Fatalf("UnmarshalJSON() = %#v, %v", parsed, err)
	}
	if parsed.SRID == nil || *parsed.SRID != 42 {
		t.Fatalf("UnmarshalJSON() SRID = %v, want 42", parsed.SRID)
	}

	if nullValue, err := (Geometry{}).Value(); err != nil || nullValue != nil {
		t.Fatalf("Value() on NULL = %#v, %v", nullValue, err)
	}
}

func TestSplitGeometrySRID(t *testing.T) {
	wkb := encodeWKBPoint(1, 2)

	if _, srid := splitGeometrySRID(wkb); srid != nil {
		t.Fatalf("splitGeometrySRID(bare WKB) SRID = %v, want nil", srid)
	}

	srid := uint32(7)
	raw := prependSRID(srid, wkb)
	payload, parsed := splitGeometrySRID(raw)
	if parsed == nil || *parsed != srid {
		t.Fatalf("splitGeometrySRID() SRID = %v, want %d", parsed, srid)
	}
	if !bytes.Equal(payload, wkb) {
		t.Fatalf("splitGeometrySRID() payload = %v, want %v", payload, wkb)
	}
}
