package datetime

import (
	"time"

	"github.com/relvacode/iso8601"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	ISO8601Date     = "2006-01-02"
	ISO8601DateTime = "2006-01-02T15:04:05Z07:00"
)

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
