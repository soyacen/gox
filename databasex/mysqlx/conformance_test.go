package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"reflect"
	"testing"
	"time"
)

// mustMarshalJSON marshals v or panics, for use in test fixtures.
//
// Parameters:
//   - v: the value to marshal.
//
// Returns:
//   - []byte: the JSON encoding.
func mustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// nullable is the common surface implemented by every exported type in the
// package.
type nullable interface {
	Scan(src any) error
	Value() (driver.Value, error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// valuesEqual compares two driver values, treating time.Time specially.
//
// Parameters:
//   - a: the first value.
//   - b: the second value.
//
// Returns:
//   - bool: true when the values are equivalent.
func valuesEqual(a any, b any) bool {
	at, aok := a.(time.Time)
	bt, bok := b.(time.Time)
	if aok && bok {
		return at.Equal(bt)
	}
	return reflect.DeepEqual(a, b)
}

func TestTypeConformance(t *testing.T) {
	srid := uint32(4326)
	pointWKB := prependSRID(srid, encodeWKBPoint(1, 2))
	geometryJSON := string(mustMarshalJSON(pointWKB))

	dateTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	date := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name      string
		new       func() nullable
		src       any
		wantValue any
		wantJSON  string
	}{
		{
			name: "Bool", new: func() nullable { return &Bool{} },
			src: true, wantValue: true, wantJSON: "true",
		},
		{
			name: "TinyInt", new: func() nullable { return &TinyInt{} },
			src: int64(math.MinInt8), wantValue: int64(math.MinInt8), wantJSON: "-128",
		},
		{
			name: "TinyIntUnsigned", new: func() nullable { return &TinyIntUnsigned{} },
			src: uint64(math.MaxUint8), wantValue: uint64(math.MaxUint8), wantJSON: "255",
		},
		{
			name: "SmallInt", new: func() nullable { return &SmallInt{} },
			src: int64(math.MinInt16), wantValue: int64(math.MinInt16), wantJSON: "-32768",
		},
		{
			name: "SmallIntUnsigned", new: func() nullable { return &SmallIntUnsigned{} },
			src: uint64(math.MaxUint16), wantValue: uint64(math.MaxUint16), wantJSON: "65535",
		},
		{
			name: "MediumInt", new: func() nullable { return &MediumInt{} },
			src: int64(mediumIntMin), wantValue: int64(mediumIntMin), wantJSON: "-8388608",
		},
		{
			name: "MediumIntUnsigned", new: func() nullable { return &MediumIntUnsigned{} },
			src: uint64(mediumIntUnsignedMax), wantValue: uint64(mediumIntUnsignedMax), wantJSON: "16777215",
		},
		{
			name: "Int", new: func() nullable { return &Int{} },
			src: int64(math.MinInt32), wantValue: int64(math.MinInt32), wantJSON: "-2147483648",
		},
		{
			name: "IntUnsigned", new: func() nullable { return &IntUnsigned{} },
			src: uint64(math.MaxUint32), wantValue: uint64(math.MaxUint32), wantJSON: "4294967295",
		},
		{
			name: "BigInt", new: func() nullable { return &BigInt{} },
			src: int64(math.MinInt64), wantValue: int64(math.MinInt64), wantJSON: "-9223372036854775808",
		},
		{
			name: "BigIntUnsigned", new: func() nullable { return &BigIntUnsigned{} },
			src: uint64(math.MaxUint64), wantValue: uint64(math.MaxUint64), wantJSON: "18446744073709551615",
		},
		{
			name: "Float", new: func() nullable { return &Float{} },
			src: float32(1.5), wantValue: float64(1.5), wantJSON: "1.5",
		},
		{
			name: "Double", new: func() nullable { return &Double{} },
			src: float64(2.25), wantValue: float64(2.25), wantJSON: "2.25",
		},
		{
			name: "Decimal", new: func() nullable { return &Decimal{} },
			src: []byte("123.4500"), wantValue: "123.4500", wantJSON: "123.4500",
		},
		{
			name: "Bit", new: func() nullable { return &Bit{} },
			src: []byte{0x01, 0x02}, wantValue: []byte{0x01, 0x02}, wantJSON: "258",
		},
		{
			name: "Date", new: func() nullable { return &Date{} },
			src: []byte("2024-01-02"), wantValue: date, wantJSON: `"2024-01-02"`,
		},
		{
			name: "DateTime", new: func() nullable { return &DateTime{} },
			src: []byte("2024-01-02 03:04:05"), wantValue: dateTime, wantJSON: `"2024-01-02T03:04:05Z"`,
		},
		{
			name: "Timestamp", new: func() nullable { return &Timestamp{} },
			src: "2024-01-02 03:04:05", wantValue: dateTime, wantJSON: `"2024-01-02T03:04:05Z"`,
		},
		{
			name: "Time", new: func() nullable { return &Time{} },
			src: "12:34:56.789", wantValue: "12:34:56.789000", wantJSON: `"12:34:56.789000"`,
		},
		{
			name: "Year", new: func() nullable { return &Year{} },
			src: int64(2024), wantValue: int64(2024), wantJSON: "2024",
		},
		{
			name: "Text", new: func() nullable { return &Text{} },
			src: "hello", wantValue: "hello", wantJSON: `"hello"`,
		},
		{
			name: "Enum", new: func() nullable { return &Enum{} },
			src: []byte("a"), wantValue: "a", wantJSON: `"a"`,
		},
		{
			name: "Set", new: func() nullable { return &Set{} },
			src: "a,b", wantValue: "a,b", wantJSON: `["a","b"]`,
		},
		{
			name: "Blob", new: func() nullable { return &Blob{} },
			src: []byte{1, 2, 3}, wantValue: []byte{1, 2, 3}, wantJSON: `"AQID"`,
		},
		{
			name: "JSON", new: func() nullable { return &JSON{} },
			src: []byte(`{"a":1}`), wantValue: []byte(`{"a":1}`), wantJSON: `{"a":1}`,
		},
		{
			name: "UUID", new: func() nullable { return &UUID{} },
			src: testUUID, wantValue: testUUID, wantJSON: `"` + testUUID + `"`,
		},
		{
			name: "IP", new: func() nullable { return &IP{} },
			src: "192.168.1.1", wantValue: "192.168.1.1", wantJSON: `"192.168.1.1"`,
		},
		{
			name: "IPPrefix", new: func() nullable { return &IPPrefix{} },
			src: "192.168.1.0/24", wantValue: "192.168.1.0/24", wantJSON: `"192.168.1.0/24"`,
		},
		{
			name: "IPPort", new: func() nullable { return &IPPort{} },
			src: "192.168.1.1:3306", wantValue: "192.168.1.1:3306", wantJSON: `"192.168.1.1:3306"`,
		},
		{
			name: "Geometry", new: func() nullable { return &Geometry{} },
			src: pointWKB, wantValue: pointWKB, wantJSON: geometryJSON,
		},
		{
			name: "Point", new: func() nullable { return &Point{} },
			src: pointWKB, wantValue: pointWKB, wantJSON: `{"x":1,"y":2,"srid":4326}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"/scan-and-value", func(t *testing.T) {
			v := tc.new()
			if err := v.Scan(tc.src); err != nil {
				t.Fatalf("Scan(%v) returned error: %v", tc.src, err)
			}

			got, err := v.Value()
			if err != nil {
				t.Fatalf("Value() returned error: %v", err)
			}
			if !valuesEqual(got, tc.wantValue) {
				t.Fatalf("Value() = %#v, want %#v", got, tc.wantValue)
			}

			jsonBytes, err := v.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() returned error: %v", err)
			}
			if string(jsonBytes) != tc.wantJSON {
				t.Fatalf("MarshalJSON() = %s, want %s", jsonBytes, tc.wantJSON)
			}
		})

		t.Run(tc.name+"/null", func(t *testing.T) {
			v := tc.new()
			if err := v.Scan(nil); err != nil {
				t.Fatalf("Scan(nil) returned error: %v", err)
			}

			got, err := v.Value()
			if err != nil {
				t.Fatalf("Value() returned error: %v", err)
			}
			if got != nil {
				t.Fatalf("Value() = %#v, want nil", got)
			}

			jsonBytes, err := v.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() returned error: %v", err)
			}
			if string(jsonBytes) != "null" {
				t.Fatalf("MarshalJSON() = %s, want null", jsonBytes)
			}
		})

		t.Run(tc.name+"/json-round-trip", func(t *testing.T) {
			v := tc.new()
			if err := v.UnmarshalJSON([]byte(tc.wantJSON)); err != nil {
				t.Fatalf("UnmarshalJSON(%s) returned error: %v", tc.wantJSON, err)
			}

			got, err := v.Value()
			if err != nil {
				t.Fatalf("Value() returned error: %v", err)
			}
			if !valuesEqual(got, tc.wantValue) {
				t.Fatalf("after UnmarshalJSON(%s) Value() = %#v, want %#v", tc.wantJSON, got, tc.wantValue)
			}
		})

		t.Run(tc.name+"/json-null", func(t *testing.T) {
			v := tc.new()
			if err := v.Scan(tc.src); err != nil {
				t.Fatalf("Scan(%v) returned error: %v", tc.src, err)
			}
			if err := v.UnmarshalJSON([]byte("null")); err != nil {
				t.Fatalf("UnmarshalJSON(null) returned error: %v", err)
			}
			got, err := v.Value()
			if err != nil {
				t.Fatalf("Value() returned error: %v", err)
			}
			if got != nil {
				t.Fatalf("Value() after null = %#v, want nil", got)
			}
		})
	}
}

func TestTypeAliases(t *testing.T) {
	// The alias declarations are checked at compile time by assigning between
	// the alias and its canonical type.
	var text Text = Char{String: "a", Valid: true}
	var char Char = text
	var tinyText TinyText = char
	var varChar VarChar = tinyText
	var mediumText MediumText = varChar
	var longText LongText = mediumText
	if longText.String != "a" {
		t.Fatalf("text aliases did not round-trip: %#v", longText)
	}

	var blob Blob = Binary{Bytes: []byte{1}, Valid: true}
	var varBinary VarBinary = blob
	var tinyBlob TinyBlob = varBinary
	var mediumBlob MediumBlob = tinyBlob
	var longBlob LongBlob = mediumBlob
	if len(longBlob.Bytes) != 1 {
		t.Fatalf("blob aliases did not round-trip: %#v", longBlob)
	}

	var geometry Geometry = LineString{Bytes: []byte{1}, Valid: true}
	var polygon Polygon = geometry
	var multiPoint MultiPoint = polygon
	var multiLineString MultiLineString = multiPoint
	var multiPolygon MultiPolygon = multiLineString
	var collection GeometryCollection = multiPolygon
	if len(collection.Bytes) != 1 {
		t.Fatalf("geometry aliases did not round-trip: %#v", collection)
	}
}
