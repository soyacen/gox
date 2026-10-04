package mysqlx

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// pointJSON is the JSON representation of a Point.
type pointJSON struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	SRID *uint32 `json:"srid,omitempty"`
}

// Point represents a MySQL POINT value. The full WKB POINT encoding is
// supported, so the coordinates can be read and written directly.
//
// MySQL stores geometry values as a four byte little-endian SRID followed by
// WKB. Scanning accepts both that layout and a bare WKB POINT. Encoding always
// produces a little-endian WKB POINT, prefixed with the SRID when SRID is not
// nil. The extended WKB (EWKB) representation is not supported.
type Point struct {
	// X is the X coordinate.
	X float64
	// Y is the Y coordinate.
	Y float64
	// SRID is the spatial reference system identifier, or nil when the value
	// carries no SRID prefix.
	SRID *uint32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src is not a valid MySQL POINT value.
func (dst *Point) Scan(src any) error {
	b, isNull, err := coerceBytes(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Point{}
		return nil
	}

	switch len(b) {
	case wkbPointSize:
		x, y, err := decodeWKBPoint(b)
		if err != nil {
			return err
		}
		*dst = Point{X: x, Y: y, Valid: true}
		return nil
	case wkbSRIDSize + wkbPointSize:
		x, y, err := decodeWKBPoint(b[wkbSRIDSize:])
		if err != nil {
			return err
		}
		srid := binary.LittleEndian.Uint32(b[:wkbSRIDSize])
		*dst = Point{X: x, Y: y, SRID: &srid, Valid: true}
		return nil
	default:
		return fmt.Errorf("%w: %d bytes is not a MySQL %s value", ErrInvalidGeometry, len(b), TypePoint)
	}
}

// Value implements the database/sql/driver.Valuer interface.
//
// Returns:
//   - driver.Value: the MySQL geometry value, or nil when the value is NULL.
//   - error: always nil.
func (src Point) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}

	wkb := encodeWKBPoint(src.X, src.Y)
	if src.SRID == nil {
		return wkb, nil
	}
	return prependSRID(*src.SRID, wkb), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The value is
// encoded as an object with x and y fields, plus srid when it is present.
//
// Returns:
//   - []byte: the JSON object, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Point) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(pointJSON{X: src.X, Y: src.Y, SRID: src.SRID})
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON object or null.
func (dst *Point) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Point{}
		return nil
	}

	var v pointJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*dst = Point{X: v.X, Y: v.Y, SRID: v.SRID, Valid: true}
	return nil
}

// Geometry represents a MySQL geometry value other than POINT. Its raw MySQL
// bytes are preserved unchanged so that any geometry value round-trips without
// loss.
//
// The payload is not decoded into geometric primitives; only the SRID prefix is
// parsed when the value is recognized as a MySQL geometry value. JSON encoding
// uses the base64 encoding of the raw bytes.
type Geometry struct {
	// Bytes is the raw MySQL geometry value, including the SRID prefix.
	Bytes []byte
	// SRID is the parsed spatial reference system identifier, or nil when the
	// value has no recognizable SRID prefix.
	SRID *uint32
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Geometry.
func (dst *Geometry) Scan(src any) error {
	b, isNull, err := coerceBytes(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Geometry{}
		return nil
	}

	if len(b) == 0 {
		return fmt.Errorf("%w: empty %s value", ErrInvalidGeometry, TypeGeometry)
	}

	_, srid := splitGeometrySRID(b)
	*dst = Geometry{Bytes: b, SRID: srid, Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface. The raw bytes are
// returned unchanged.
//
// Returns:
//   - driver.Value: the raw geometry value, or nil when the value is NULL.
//   - error: always nil.
func (src Geometry) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Bytes, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The raw bytes
// are encoded as a base64 JSON string.
//
// Returns:
//   - []byte: the JSON string, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Geometry) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Bytes)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The value
// is expected to be a base64 JSON string.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a base64 JSON string or null.
func (dst *Geometry) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Geometry{}
		return nil
	}

	var b []byte
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	_, srid := splitGeometrySRID(b)
	*dst = Geometry{Bytes: b, SRID: srid, Valid: true}
	return nil
}

// LineString is a MySQL LINESTRING value. It is an alias for Geometry.
type LineString = Geometry

// Polygon is a MySQL POLYGON value. It is an alias for Geometry.
type Polygon = Geometry

// MultiPoint is a MySQL MULTIPOINT value. It is an alias for Geometry.
type MultiPoint = Geometry

// MultiLineString is a MySQL MULTILINESTRING value. It is an alias for Geometry.
type MultiLineString = Geometry

// MultiPolygon is a MySQL MULTIPOLYGON value. It is an alias for Geometry.
type MultiPolygon = Geometry

// GeometryCollection is a MySQL GEOMETRYCOLLECTION value. It is an alias for
// Geometry.
type GeometryCollection = Geometry
