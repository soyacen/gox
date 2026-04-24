package slogx

import (
	"encoding/json"
	"log/slog"
	"time"

	"golang.org/x/exp/constraints"
)

// Bool creates a slog.Attr with a bool value.
//
// Parameters:
//   - key: the attribute key
//   - value: the bool value
//
// Returns:
//   - slog.Attr: the created attribute
func Bool[Bool ~bool](key string, value Bool) slog.Attr {
	return slog.Attr{Key: key, Value: slog.BoolValue(bool(value))}
}

// Int creates a slog.Attr with a signed integer value.
//
// Parameters:
//   - key: the attribute key
//   - value: the signed integer value
//
// Returns:
//   - slog.Attr: the created attribute
func Int[Int constraints.Signed](key string, value Int) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Int64Value(int64(value))}
}

// Uint creates a slog.Attr with an unsigned integer value.
//
// Parameters:
//   - key: the attribute key
//   - value: the unsigned integer value
//
// Returns:
//   - slog.Attr: the created attribute
func Uint[Uint constraints.Unsigned](key string, value Uint) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Uint64Value(uint64(value))}
}

// Duration creates a slog.Attr with a time.Duration value.
//
// Parameters:
//   - key: the attribute key
//   - value: the duration value
//
// Returns:
//   - slog.Attr: the created attribute
func Duration[Duration time.Duration](key string, value Duration) slog.Attr {
	return slog.Attr{Key: key, Value: slog.DurationValue(time.Duration(value))}
}

// Float creates a slog.Attr with a float value.
//
// Parameters:
//   - key: the attribute key
//   - value: the float value
//
// Returns:
//   - slog.Attr: the created attribute
func Float[Float constraints.Float](key string, value Float) slog.Attr {
	return slog.Attr{Key: key, Value: slog.Float64Value(float64(value))}
}

// String creates a slog.Attr with a string value.
//
// Parameters:
//   - key: the attribute key
//   - value: the string value
//
// Returns:
//   - slog.Attr: the created attribute
func String[String ~string](key string, value String) slog.Attr {
	return slog.Attr{Key: key, Value: slog.StringValue(string(value))}
}

// Time creates a slog.Attr with a time.Time value.
//
// Parameters:
//   - key: the attribute key
//   - value: the time value
//
// Returns:
//   - slog.Attr: the created attribute
func Time[Time time.Time](key string, value Time) slog.Attr {
	return slog.Attr{Key: key, Value: slog.TimeValue(time.Time(value))}
}

// Json creates a slog.Attr with a JSON-encoded value.
// If the value is nil, it returns a string attribute with "<nil>".
// If marshaling fails, it panics.
//
// Parameters:
//   - key: the attribute key
//   - value: the value to JSON encode
//
// Returns:
//   - slog.Attr: the created attribute with JSON string value
func Json(key string, value any) slog.Attr {
	if value == nil {
		return String(key, "<nil>")
	}
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return slog.Attr{Key: key, Value: slog.StringValue(string(data))}
}

// Error creates a slog.Attr with an error value.
// If the error is nil, it returns a string attribute with "<nil>".
//
// Parameters:
//   - key: the attribute key
//   - value: the error value
//
// Returns:
//   - slog.Attr: the created attribute with the error message
func Error(key string, value error) slog.Attr {
	if value == nil {
		return String(key, "<nil>")
	}
	return String(key, value.Error())
}

// func Valuer(key string, value slog.LogValuer) slog.Attr {
// 	return slog.Attr{Key: key, Value: value.LogValue()}
// }

// KindGroup
// KindLogValuer
