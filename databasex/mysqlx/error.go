package mysqlx

import "errors"

// ErrCannotScan is returned when a value produced by a database driver cannot be
// scanned into the destination type.
var ErrCannotScan = errors.New("mysqlx: cannot scan value into destination")

// ErrOutOfRange is returned when a value does not fit into the range of the
// target MySQL type.
var ErrOutOfRange = errors.New("mysqlx: value out of range")

// ErrInvalidDecimal is returned when a value cannot be parsed as a MySQL
// DECIMAL literal.
var ErrInvalidDecimal = errors.New("mysqlx: invalid decimal")

// ErrInvalidTime is returned when a value cannot be parsed as a MySQL TIME
// literal.
var ErrInvalidTime = errors.New("mysqlx: invalid time")

// ErrInvalidYear is returned when a value cannot be parsed as a MySQL YEAR
// value.
var ErrInvalidYear = errors.New("mysqlx: invalid year")

// ErrInvalidGeometry is returned when a value cannot be parsed as a MySQL
// geometry value.
var ErrInvalidGeometry = errors.New("mysqlx: invalid geometry")

// ErrInvalidJSON is returned when a value cannot be parsed as a JSON document.
var ErrInvalidJSON = errors.New("mysqlx: invalid json")
