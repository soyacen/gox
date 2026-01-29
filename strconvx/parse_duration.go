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

func ParseTime(s string, layout string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, s, loc)
}

func ParseWrapperTimestamp(s string, layout string, loc *time.Location) (*timestamppb.Timestamp, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(v), nil
}

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

func ParseWrapperMonth(s string, layout string, loc *time.Location) (month.Month, error) {
	v, err := ParseTime(s, layout, loc)
	if err != nil {
		return month.Month_MONTH_UNSPECIFIED, err
	}
	switch v.Month() {
	case time.January:
		return month.Month_JANUARY
	case time.February:
		return month.Month_FEBRUARY
	case time.March:
		return month.Month_MARCH
	case time.April:
		return month.Month_APRIL
	case time.May:
		return month.Month_MAY
	case time.June:
		return month.Month_JUNE
	case time.July:
		return month.Month_JULY
	case time.August:
		return month.Month_AUGUST
	case time.September:
		return month.Month_SEPTEMBER
	case time.October:
		return month.Month_OCTOBER
	case time.November:
		return month.Month_NOVEMBER
	case time.December:
		return month.Month_DECEMBER
	}
	return month.Month_MONTH_UNSPECIFIED, errors.New("invalid month")
}

func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func ParseWrapperDuration(s string) (*durationpb.Duration, error) {
	v, err := ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return durationpb.New(v), nil
}
