package strconvx

import (
	"errors"
	"time"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/genproto/googleapis/type/month"
	"google.golang.org/genproto/googleapis/type/timeofday"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseTime parses a string into a time.Time value in the given location.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string (e.g., time.RFC3339).
//   - loc: the location for time zone offset.
//
// Returns:
//   - time.Time: the parsed time value.
//   - error: an error if parsing fails.
func ParseTime(s string, layout string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, s, loc)
}

// ParseWrapperTimestamp parses a string into a protobuf Timestamp wrapper.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - *timestamppb.Timestamp: the protobuf wrapper containing the parsed time.
//   - error: an error if parsing fails.
func ParseWrapperTimestamp(s string, layout string, loc *time.Location) (*timestamppb.Timestamp, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(v), nil
}

// ParseWrapperDateTime parses a string into a protobuf DateTime wrapper.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - *datetime.DateTime: the protobuf wrapper containing the parsed datetime.
//   - error: an error if parsing fails.
func ParseWrapperDateTime(s string, layout string, loc *time.Location) (*datetime.DateTime, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return nil, err
	}
	_, offset := v.Zone()
	dt := &datetime.DateTime{
		Year:    int32(v.Year()),
		Month:   int32(v.Month()),
		Day:     int32(v.Day()),
		Hours:   int32(v.Hour()),
		Minutes: int32(v.Minute()),
		Seconds: int32(v.Second()),
		Nanos:   int32(v.Nanosecond()),
		TimeOffset: &datetime.DateTime_UtcOffset{
			UtcOffset: durationpb.New(time.Duration(offset) * time.Second),
		},
	}
	return dt, nil
}

// ParseWrapperTime parses a string into a protobuf Date wrapper.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - *date.Date: the protobuf wrapper containing the parsed date.
//   - error: an error if parsing fails.
func ParseWrapperTime(s string, layout string, loc *time.Location) (*date.Date, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return nil, err
	}
	dt := &date.Date{
		Year:  int32(v.Year()),
		Month: int32(v.Month()),
		Day:   int32(v.Day()),
	}
	return dt, nil
}

// ParseWrapperTimeOfDay parses a string into a protobuf TimeOfDay wrapper.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - *timeofday.TimeOfDay: the protobuf wrapper containing the parsed time of day.
//   - error: an error if parsing fails.
func ParseWrapperTimeOfDay(s string, layout string, loc *time.Location) (*timeofday.TimeOfDay, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return nil, err
	}
	dt := &timeofday.TimeOfDay{
		Hours:   int32(v.Hour()),
		Minutes: int32(v.Minute()),
		Seconds: int32(v.Second()),
		Nanos:   int32(v.Nanosecond()),
	}
	return dt, nil
}

// ParseWrapperDayOfWeek parses a string into a DayOfWeek enum value.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - dayofweek.DayOfWeek: the parsed day of week enum value.
//   - error: an error if parsing fails.
func ParseWrapperDayOfWeek(s string, layout string, loc *time.Location) (dayofweek.DayOfWeek, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return dayofweek.DayOfWeek_DAY_OF_WEEK_UNSPECIFIED, err
	}
	switch v.Weekday() {
	case time.Sunday:
		return dayofweek.DayOfWeek_SUNDAY, nil
	case time.Monday:
		return dayofweek.DayOfWeek_MONDAY, nil
	case time.Tuesday:
		return dayofweek.DayOfWeek_TUESDAY, nil
	case time.Wednesday:
		return dayofweek.DayOfWeek_WEDNESDAY, nil
	case time.Thursday:
		return dayofweek.DayOfWeek_THURSDAY, nil
	case time.Friday:
		return dayofweek.DayOfWeek_FRIDAY, nil
	case time.Saturday:
		return dayofweek.DayOfWeek_SATURDAY, nil
	}
	return dayofweek.DayOfWeek_DAY_OF_WEEK_UNSPECIFIED, errors.New("invalid day of week")
}

// ParseWrapperMonth parses a string into a Month enum value.
//
// Parameters:
//   - s: the string to parse.
//   - layout: the time layout string.
//   - loc: the location for time zone offset.
//
// Returns:
//   - month.Month: the parsed month enum value.
//   - error: an error if parsing fails.
func ParseWrapperMonth(s string, layout string, loc *time.Location) (month.Month, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return month.Month_MONTH_UNSPECIFIED, err
	}
	switch v.Month() {
	case time.January:
		return month.Month_JANUARY, nil
	case time.February:
		return month.Month_FEBRUARY, nil
	case time.March:
		return month.Month_MARCH, nil
	case time.April:
		return month.Month_APRIL, nil
	case time.May:
		return month.Month_MAY, nil
	case time.June:
		return month.Month_JUNE, nil
	case time.July:
		return month.Month_JULY, nil
	case time.August:
		return month.Month_AUGUST, nil
	case time.September:
		return month.Month_SEPTEMBER, nil
	case time.October:
		return month.Month_OCTOBER, nil
	case time.November:
		return month.Month_NOVEMBER, nil
	case time.December:
		return month.Month_DECEMBER, nil
	}
	return month.Month_MONTH_UNSPECIFIED, errors.New("invalid month")
}

// ParseDuration parses a duration string.
//
// Parameters:
//   - s: the duration string to parse (e.g., "1h30m").
//
// Returns:
//   - time.Duration: the parsed duration.
//   - error: an error if parsing fails.
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// ParseWrapperDuration parses a duration string into a protobuf Duration wrapper.
//
// Parameters:
//   - s: the duration string to parse.
//
// Returns:
//   - *durationpb.Duration: the protobuf wrapper containing the parsed duration.
//   - error: an error if parsing fails.
func ParseWrapperDuration(s string) (*durationpb.Duration, error) {
	v, err := ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return durationpb.New(v), nil
}
