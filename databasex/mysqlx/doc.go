// Package mysqlx provides nullable Go representations for the MySQL data types.
//
// The package is modeled after github.com/jackc/pgx/v5/pgtype: every MySQL type
// is represented by a plain struct carrying the decoded value plus a Valid flag
// that distinguishes the SQL NULL value from the zero value. Each type
// implements the [database/sql.Scanner] and [database/sql/driver.Valuer]
// interfaces as well as [encoding/json.Marshaler] and
// [encoding/json.Unmarshaler], so the types can be used directly with
// database/sql, query builders and ORMs.
//
// # Type Mapping
//
// The following table summarizes the mapping between Go and MySQL. Length and
// semantic variants that share the exact same Go representation are exposed as
// type aliases.
//
//	Go                        MySQL
//	------------------------------------------------------------------
//	Bool                      bool
//	                          boolean
//
//	TinyInt                   tinyint
//	TinyIntUnsigned           tinyint unsigned
//	SmallInt                  smallint
//	SmallIntUnsigned          smallint unsigned
//	MediumInt                 mediumint
//	MediumIntUnsigned         mediumint unsigned
//	Int                       int
//	                          integer
//	IntUnsigned               int unsigned
//	BigInt                    bigint
//	BigIntUnsigned            bigint unsigned
//
//	Float                     float
//	Double                    double
//	                          real
//	                          double precision
//	Decimal                   decimal
//	                          numeric
//	                          dec
//	                          fixed
//	Bit                       bit
//
//	Date                      date
//	DateTime                  datetime
//	Timestamp                 timestamp
//	Time                      time
//	Year                      year
//
//	Text                      char
//	                          varchar
//	                          tinytext
//	                          text
//	                          mediumtext
//	                          longtext
//	Enum                      enum
//	Set                       set
//
//	Blob                      binary
//	                          varbinary
//	                          tinyblob
//	                          blob
//	                          mediumblob
//	                          longblob
//
//	JSON                      json
//	UUID                      uuid (MariaDB only)
//
//	Geometry                  geometry
//	                          linestring
//	                          polygon
//	                          multipoint
//	                          multilinestring
//	                          multipolygon
//	                          geometrycollection
//	Point                     point
//
// # Null Values
//
// A NULL column value is represented by Valid being false. Scanning NULL resets
// the destination to its zero value, and encoding an invalid value yields a SQL
// NULL. Marshaling an invalid value yields the JSON null literal.
//
//	// column is nullable
//	var name mysqlx.Text
//	if err := row.Scan(&name); err != nil {
//		return err
//	}
//	if name.Valid {
//		fmt.Println(name.String)
//	}
//
// When passing a nullable type as a query parameter, remember to set Valid
// explicitly, otherwise the parameter is encoded as NULL.
//
//	mysqlx.Text{String: "gopher", Valid: true}
//
// # Scanning Sources
//
// database/sql drivers are free to return different Go types for the same
// column depending on the driver configuration. The Scan methods of this
// package accept every representation commonly produced by MySQL drivers:
//
//	nil                                     NULL
//	[]byte, string                          text, numeric and binary columns
//	int64, uint64                          integer columns
//	float32, float64                       floating point columns
//	time.Time                              date/time columns (parseTime=true)
//	bool                                   boolean columns
//
// Unsupported source types produce an error wrapping [ErrCannotScan].
//
// # Zero Dates
//
// MySQL can store the all-zero date literals "0000-00-00" and
// "0000-00-00 00:00:00". Because time.Time cannot represent them, [Date],
// [DateTime] and [Timestamp] scan them as NULL (Valid stays false) without
// returning an error.
//
// # Time Zones
//
// DATETIME and DATE columns carry no time zone information. When the driver
// returns the value as text, it is parsed in UTC. When the driver returns a
// time.Time (parseTime=true) its location is preserved.
//
// # UUID Values
//
// MySQL has no dedicated UUID type. Identifiers are stored in CHAR(36) columns
// using the canonical hyphenated form, or in BINARY(16) columns using the packed
// form produced by the UUID_TO_BIN function. [UUID] scans both representations
// and writes the canonical form, which the server converts with UUID_TO_BIN when
// the target column is BINARY(16).
//
//	// ORDER BY created_at keeps v7 identifiers in insertion order.
//	id := mysqlx.NewUUIDV7()
//
//	// BINARY(16) columns can be round-tripped without a textual form.
//	var packed []byte = id.Bin()
//
// [UUID] is backed by the standard library [uuid] package, so building this
// package requires Go 1.27 or later.
package mysqlx
