package timex

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash"
	"time"
)

// UTCTime is a wrapper around time.Time that always serializes in UTC.
type UTCTime time.Time

// MarshalJSON implements json.Marshaler by formatting the time in UTC
// using the UTCLayout format.
//
// Returns:
//   - []byte: the JSON-encoded UTC time string.
//   - error: any marshaling error.
func (t UTCTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format(UTCLayout))
}

// UnmarshalJSON implements json.Unmarshaler by parsing a quoted UTC time
// string using the UTCLayout format.
//
// Parameters:
//   - data: the JSON data containing a quoted time string.
//
// Returns:
//   - error: an error if the data is not quoted or parsing fails.
func (t *UTCTime) UnmarshalJSON(data []byte) (err error) {
	if data[0] != []byte(`"`)[0] || data[len(data)-1] != []byte(`"`)[0] {
		return errors.New("Not quoted")
	}
	*t, err = ParseUTCTime(string(data[1 : len(data)-1]))
	return
}

// Hash32 writes the Unix nanosecond timestamp into the provided hash.
//
// Parameters:
//   - h: the hash.Hash32 to write into.
//
// Returns:
//   - error: any error from writing to the hash.
func (t UTCTime) Hash32(h hash.Hash32) error {
	err := binary.Write(h, binary.LittleEndian, time.Time(t).UnixNano())
	return err
}

// String returns the UTC time formatted using UTCLayout.
//
// Returns:
//   - string: the formatted UTC time string.
func (t UTCTime) String() string {
	return time.Time(t).Format(UTCLayout)
}

// ParseUTCTime parses a UTC time string using the UTCLayout format.
//
// Parameters:
//   - timespec: the time string to parse.
//
// Returns:
//   - UTCTime: the parsed UTC time.
//   - error: any parsing error.
func ParseUTCTime(timespec string) (UTCTime, error) {
	t, err := time.Parse(UTCLayout, timespec)
	return UTCTime(t), err
}

// MustParseUTCTime parses a UTC time string and panics on error.
//
// Parameters:
//   - timespec: the time string to parse.
//
// Returns:
//   - UTCTime: the parsed UTC time.
func MustParseUTCTime(timespec string) UTCTime {
	ts, err := ParseUTCTime(timespec)
	if err != nil {
		panic(err)
	}
	return ts
}
