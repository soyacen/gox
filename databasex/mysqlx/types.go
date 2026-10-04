package mysqlx

// MySQL data type names. The constants can be used when building DDL statements
// or when mapping driver metadata to Go types.
const (
	// TypeBool is the MySQL BOOL type name.
	TypeBool = "bool"
	// TypeBoolean is the MySQL BOOLEAN type name, a synonym for BOOL.
	TypeBoolean = "boolean"

	// TypeTinyInt is the MySQL TINYINT type name.
	TypeTinyInt = "tinyint"
	// TypeSmallInt is the MySQL SMALLINT type name.
	TypeSmallInt = "smallint"
	// TypeMediumInt is the MySQL MEDIUMINT type name.
	TypeMediumInt = "mediumint"
	// TypeInt is the MySQL INT type name.
	TypeInt = "int"
	// TypeInteger is the MySQL INTEGER type name, a synonym for INT.
	TypeInteger = "integer"
	// TypeBigInt is the MySQL BIGINT type name.
	TypeBigInt = "bigint"

	// TypeFloat is the MySQL FLOAT type name.
	TypeFloat = "float"
	// TypeDouble is the MySQL DOUBLE type name.
	TypeDouble = "double"
	// TypeReal is the MySQL REAL type name, a synonym for DOUBLE.
	TypeReal = "real"
	// TypeDoublePrecision is the MySQL DOUBLE PRECISION type name, a synonym for DOUBLE.
	TypeDoublePrecision = "double precision"
	// TypeDecimal is the MySQL DECIMAL type name.
	TypeDecimal = "decimal"
	// TypeNumeric is the MySQL NUMERIC type name, a synonym for DECIMAL.
	TypeNumeric = "numeric"
	// TypeDec is the MySQL DEC type name, a synonym for DECIMAL.
	TypeDec = "dec"
	// TypeFixed is the MySQL FIXED type name, a synonym for DECIMAL.
	TypeFixed = "fixed"
	// TypeBit is the MySQL BIT type name.
	TypeBit = "bit"

	// TypeDate is the MySQL DATE type name.
	TypeDate = "date"
	// TypeDateTime is the MySQL DATETIME type name.
	TypeDateTime = "datetime"
	// TypeTimestamp is the MySQL TIMESTAMP type name.
	TypeTimestamp = "timestamp"
	// TypeTime is the MySQL TIME type name.
	TypeTime = "time"
	// TypeYear is the MySQL YEAR type name.
	TypeYear = "year"

	// TypeChar is the MySQL CHAR type name.
	TypeChar = "char"
	// TypeVarChar is the MySQL VARCHAR type name.
	TypeVarChar = "varchar"
	// TypeTinyText is the MySQL TINYTEXT type name.
	TypeTinyText = "tinytext"
	// TypeText is the MySQL TEXT type name.
	TypeText = "text"
	// TypeMediumText is the MySQL MEDIUMTEXT type name.
	TypeMediumText = "mediumtext"
	// TypeLongText is the MySQL LONGTEXT type name.
	TypeLongText = "longtext"
	// TypeEnum is the MySQL ENUM type name.
	TypeEnum = "enum"
	// TypeSet is the MySQL SET type name.
	TypeSet = "set"

	// TypeBinary is the MySQL BINARY type name.
	TypeBinary = "binary"
	// TypeVarBinary is the MySQL VARBINARY type name.
	TypeVarBinary = "varbinary"
	// TypeTinyBlob is the MySQL TINYBLOB type name.
	TypeTinyBlob = "tinyblob"
	// TypeBlob is the MySQL BLOB type name.
	TypeBlob = "blob"
	// TypeMediumBlob is the MySQL MEDIUMBLOB type name.
	TypeMediumBlob = "mediumblob"
	// TypeLongBlob is the MySQL LONGBLOB type name.
	TypeLongBlob = "longblob"

	// TypeJSON is the MySQL JSON type name.
	TypeJSON = "json"

	// TypeUUID is the MariaDB UUID type name. MySQL itself has no native UUID
	// type and stores identifiers in CHAR(36) or BINARY(16) columns instead.
	TypeUUID = "uuid"

	// TypeGeometry is the MySQL GEOMETRY type name.
	TypeGeometry = "geometry"
	// TypePoint is the MySQL POINT type name.
	TypePoint = "point"
	// TypeLineString is the MySQL LINESTRING type name.
	TypeLineString = "linestring"
	// TypePolygon is the MySQL POLYGON type name.
	TypePolygon = "polygon"
	// TypeMultiPoint is the MySQL MULTIPOINT type name.
	TypeMultiPoint = "multipoint"
	// TypeMultiLineString is the MySQL MULTILINESTRING type name.
	TypeMultiLineString = "multilinestring"
	// TypeMultiPolygon is the MySQL MULTIPOLYGON type name.
	TypeMultiPolygon = "multipolygon"
	// TypeGeometryCollection is the MySQL GEOMETRYCOLLECTION type name.
	TypeGeometryCollection = "geometrycollection"
)
