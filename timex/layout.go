package timex

import (
	"net/http"
)

// YearMonth is the layout string for "2006-01" format.
//
// It represents a year and month layout used for time formatting and parsing.
const YearMonth = "2006-01"

// DateMinute is the layout string for "2006-01-02 15:04" format.
//
// It represents a layout including year, month, day, hour, and minute.
const DateMinute = "2006-01-02 15:04"

// DateHour is the layout string for "2006-01-02 15" format.
//
// It represents a layout including year, month, day, and hour.
const DateHour = "2006-01-02 15"

// DateTimeNano is the layout string for "2006-01-02 15:04:05.999999999" format.
//
// It represents a layout including full date, time, and nanosecond precision.
const DateTimeNano = "2006-01-02 15:04:05.999999999"

// UTCLayout is the layout string for "2006-01-02T15:04:05.000Z" UTC format.
//
// It represents the ISO 8601 UTC timestamp layout.
const UTCLayout = "2006-01-02T15:04:05.000Z"

// HttpTimeLayout is the layout string for HTTP time format.
//
// It represents the standard HTTP time format as defined by RFC 7231.
const HttpTimeLayout = http.TimeFormat
