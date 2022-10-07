package datetime

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/relvacode/iso8601"
	"google.golang.org/protobuf/types/known/wrapperspb"

	dpb "google.golang.org/genproto/googleapis/type/date"
	dtpb "google.golang.org/genproto/googleapis/type/datetime"

	durpb "google.golang.org/protobuf/types/known/durationpb"
)

const (
	ISO8601Date     = "2006-01-02"
	ISO8601DateTime = "2006-01-02T15:04:05Z07:00"
)

var ErrInvalidValue = errors.New("invalid value")

func TimeToISO8601DateStringWrapper(t *time.Time) *wrapperspb.StringValue {
	if t != nil {
		return wrapperspb.String(TimeToISO8601DateString(*t))
	} else {
		return nil
	}
}

func TimeToISO8601DateTimeStringWrapper(t *time.Time) *wrapperspb.StringValue {
	if t != nil {
		return wrapperspb.String(TimeToISO8601DateTimeString(*t))
	} else {
		return nil
	}
}

func TimeToLocalISO8601DateTimeStringWrapper(t *time.Time, location *time.Location) *wrapperspb.StringValue {
	if t != nil {
		localTime := (*t).In(location)
		return TimeToISO8601DateTimeStringWrapper(&localTime)
	} else {
		return nil
	}
}

func TimeToISO8601DateString(t time.Time) string {
	return t.Format(ISO8601Date)
}

func TimeToISO8601DateTimeString(t time.Time) string {
	return t.Format(ISO8601DateTime)
}

func TimeToLocalISO8601DateTimeString(t time.Time, location *time.Location) string {
	return TimeToISO8601DateTimeString(t.In(location))
}

// Converts to "local" time, meaning that the offset/timezone does not have to be UTC
func ISO8601StringToTime(dateTime string) (time.Time, error) {
	return iso8601.ParseString(dateTime)
}

// Converts to UTC time
func ISO8601StringToUTCTime(dateTime string) (time.Time, error) {
	time, err := ISO8601StringToTime(dateTime)
	if err != nil {
		return time, err
	}
	return time.UTC(), err
}

// ProtoDateToLocalTime returns a new Time based on the google.type.Date, in
// the system's time zone.
//
// Hours, minues, seconds, and nanoseconds are set to 0.
func ProtoDateToLocalTime(d *dpb.Date) (time.Time, error) {
	return ProtoDateToTime(d, time.Local)
}

// ProtoDateToUTCTime returns a new Time based on the google.type.Date, in UTC.
//
// Hours, minutes, seconds, and nanoseconds are set to 0.
func ProtoDateToUTCTime(d *dpb.Date) (time.Time, error) {
	return ProtoDateToTime(d, time.UTC)
}

// ProtoDateToTime returns a new Time based on the google.type.Date and provided
// *time.Location.
//
// Hours, minutes, seconds, and nanoseconds are set to 0.
func ProtoDateToTime(d *dpb.Date, l *time.Location) (time.Time, error) {
	if d == nil {
		return time.Time{}, fmt.Errorf("%w: date parameter not set", ErrInvalidValue)
	}

	if d.GetYear() < 1 || d.GetMonth() < 1 || d.GetDay() < 1 {
		return time.Time{}, fmt.Errorf("%w: year, month, day not set", ErrInvalidValue)
	}

	return time.Date(
		int(d.GetYear()),
		time.Month(d.GetMonth()),
		int(d.GetDay()), 0, 0, 0, 0, l), nil
}

// TimeToProtoDate returns a new google.type.Date based on the provided time.Time.
// The location is ignored, as is anything more precise than the day.
func TimeToProtoDate(t time.Time) *dpb.Date {
	return &dpb.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

func TimeToProtoDateTime(t time.Time) *dtpb.DateTime {
	dt := &dtpb.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}

	// If the location is a UTC offset, encode it as such in the proto.
	zone, offset := t.Zone()

	// distinguish between time zone and utc offset
	var offsetRegexp = regexp.MustCompile(`^UTC([+-][\d]{1,2})$`)

	// Use utc offset if match or empty
	match := offsetRegexp.FindStringSubmatch(zone)
	if len(zone) == 0 || len(match) > 0 {
		if offset > 0 {
			dt.TimeOffset = &dtpb.DateTime_UtcOffset{
				UtcOffset: &durpb.Duration{Seconds: int64(offset)},
			}
		}
	} else {
		dt.TimeOffset = &dtpb.DateTime_TimeZone{
			TimeZone: &dtpb.TimeZone{Id: zone},
		}
	}

	return dt
}

func ProtoDateTimeToTime(d *dtpb.DateTime) (time.Time, error) {
	if d == nil {
		return time.Time{}, fmt.Errorf("%w: date parameter not set", ErrInvalidValue)
	}

	if d.GetYear() < 1 || d.GetMonth() < 1 || d.GetDay() < 1 {
		return time.Time{}, fmt.Errorf("%w: year, month, day not set", ErrInvalidValue)
	}

	var err error

	// Determine the location.
	loc := time.UTC
	if tz := d.GetTimeZone(); tz != nil {
		loc, err = time.LoadLocation(tz.GetId())
		if err != nil {
			return time.Time{}, err
		}
	}
	if offset := d.GetUtcOffset(); offset != nil {
		hours := int(offset.GetSeconds()) / 3600
		loc = time.FixedZone(fmt.Sprintf("UTC%+d", hours), int(offset.GetSeconds()))
	}

	// Return the Time.
	return time.Date(
		int(d.GetYear()),
		time.Month(d.GetMonth()),
		int(d.GetDay()),
		int(d.GetHours()),
		int(d.GetMinutes()),
		int(d.GetSeconds()),
		int(d.GetNanos()),
		loc,
	), nil
}

func UTCTimeAdjustedToStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func UTCTimeAdjustedToEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC)
}

func IsSameDate(t1, t2 time.Time) bool {
	if t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year() {
		return true
	} else {
		return false
	}
}

func IsStartOfDay(t time.Time) bool {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return true
	} else {
		return false
	}
}

func IsEndOfDay(t time.Time) bool {
	if t.Hour() == 23 && t.Minute() == 59 && t.Second() == 59 {
		return true
	} else {
		return false
	}
}
