package mysqlx

import (
	"encoding/binary"
	"fmt"
	"math"
)

// WKB byte order markers.
const (
	wkbByteOrderBigEndian    = 0
	wkbByteOrderLittleEndian = 1
)

// wkbPointType is the WKB geometry type code of a POINT.
const wkbPointType = 1

// wkbPointSize is the size in bytes of a little-endian WKB POINT: one byte order
// byte, four type bytes and two float64 coordinates.
const wkbPointSize = 21

// wkbSRIDSize is the size in bytes of the SRID prefix MySQL prepends to WKB.
const wkbSRIDSize = 4

// validWKBType reports whether t is one of the WKB geometry type codes used by
// MySQL.
//
// Parameters:
//   - t: the WKB geometry type code.
//
// Returns:
//   - bool: true when t is between 1 (POINT) and 7 (GEOMETRYCOLLECTION).
func validWKBType(t uint32) bool {
	return t >= 1 && t <= 7
}

// decodeWKBPoint decodes a bare WKB POINT in either byte order.
//
// Parameters:
//   - b: the WKB payload, exactly wkbPointSize bytes long.
//
// Returns:
//   - float64: the X coordinate.
//   - float64: the Y coordinate.
//   - error: non-nil when b is not a WKB POINT.
func decodeWKBPoint(b []byte) (float64, float64, error) {
	if len(b) != wkbPointSize {
		return 0, 0, fmt.Errorf("%w: WKB %s must be %d bytes, got %d", ErrInvalidGeometry, TypePoint, wkbPointSize, len(b))
	}

	switch b[0] {
	case wkbByteOrderLittleEndian:
		if t := binary.LittleEndian.Uint32(b[1:5]); t != wkbPointType {
			return 0, 0, fmt.Errorf("%w: WKB geometry type %d is not %s", ErrInvalidGeometry, t, TypePoint)
		}
		x := math.Float64frombits(binary.LittleEndian.Uint64(b[5:13]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[13:21]))
		return x, y, nil
	case wkbByteOrderBigEndian:
		if t := binary.BigEndian.Uint32(b[1:5]); t != wkbPointType {
			return 0, 0, fmt.Errorf("%w: WKB geometry type %d is not %s", ErrInvalidGeometry, t, TypePoint)
		}
		x := math.Float64frombits(binary.BigEndian.Uint64(b[5:13]))
		y := math.Float64frombits(binary.BigEndian.Uint64(b[13:21]))
		return x, y, nil
	default:
		return 0, 0, fmt.Errorf("%w: unknown WKB byte order %d", ErrInvalidGeometry, b[0])
	}
}

// encodeWKBPoint encodes a point as a little-endian WKB POINT, the byte order
// MySQL uses.
//
// Parameters:
//   - x: the X coordinate.
//   - y: the Y coordinate.
//
// Returns:
//   - []byte: the WKB payload.
func encodeWKBPoint(x float64, y float64) []byte {
	b := make([]byte, wkbPointSize)
	b[0] = wkbByteOrderLittleEndian
	binary.LittleEndian.PutUint32(b[1:5], wkbPointType)
	binary.LittleEndian.PutUint64(b[5:13], math.Float64bits(x))
	binary.LittleEndian.PutUint64(b[13:21], math.Float64bits(y))
	return b
}

// splitGeometrySRID splits a raw MySQL geometry value into its WKB payload and
// its optional SRID prefix. MySQL stores an SRID prefix followed by WKB, but a
// bare WKB value is also accepted.
//
// Parameters:
//   - b: the raw geometry bytes.
//
// Returns:
//   - []byte: the WKB payload without the SRID prefix.
//   - *uint32: the parsed SRID, or nil when b has no SRID prefix.
func splitGeometrySRID(b []byte) ([]byte, *uint32) {
	if len(b) < wkbSRIDSize+5 {
		return b, nil
	}

	order := b[wkbSRIDSize]
	var geometryType uint32
	switch order {
	case wkbByteOrderLittleEndian:
		geometryType = binary.LittleEndian.Uint32(b[wkbSRIDSize+1 : wkbSRIDSize+5])
	case wkbByteOrderBigEndian:
		geometryType = binary.BigEndian.Uint32(b[wkbSRIDSize+1 : wkbSRIDSize+5])
	default:
		return b, nil
	}

	if !validWKBType(geometryType) {
		return b, nil
	}

	srid := binary.LittleEndian.Uint32(b[:wkbSRIDSize])
	return b[wkbSRIDSize:], &srid
}

// prependSRID prefixes a WKB payload with a little-endian SRID.
//
// Parameters:
//   - srid: the SRID to prepend.
//   - wkb: the WKB payload.
//
// Returns:
//   - []byte: the MySQL geometry value.
func prependSRID(srid uint32, wkb []byte) []byte {
	out := make([]byte, wkbSRIDSize+len(wkb))
	binary.LittleEndian.PutUint32(out[:wkbSRIDSize], srid)
	copy(out[wkbSRIDSize:], wkb)
	return out
}
