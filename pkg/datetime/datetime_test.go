package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	dpb "google.golang.org/genproto/googleapis/type/date"
	dtpb "google.golang.org/genproto/googleapis/type/datetime"
	tod "google.golang.org/genproto/googleapis/type/timeofday"
	durpb "google.golang.org/protobuf/types/known/durationpb"
)

func Test_ISO8601StringToTime(t *testing.T) {

	// Verifies that various strings with offset/timezone are correctly parsed
	// to both local time and to UTC time + that we can keep track of the offset

	require := require.New(t)

	// ...start with just a date:

	from := "2006-01-02"

	utcTime, err := ISO8601StringToUTCTime(from)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(from, TimeToISO8601DateString(utcTime))
	require.Equal("2006-01-02T00:00:00Z", TimeToISO8601DateTimeString(utcTime))
	require.True(IsStartOfDay(utcTime))

	// ...then an UTC date time:

	from = "2006-01-02T01:04:05Z"

	localTime, err := ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}
	utcTime, err = ISO8601StringToUTCTime(from)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(from, TimeToISO8601DateTimeString(localTime))
	require.Equal(from, TimeToISO8601DateTimeString(utcTime))
	require.Equal(utcTime, localTime)

	// ...a negative UTC offset/timezone:

	from = "2006-01-02T15:04:05-07:00"

	localTime, err = ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}
	utcTime, err = ISO8601StringToUTCTime(from)
	if err != nil {
		t.Fatal(err)
	}

	_, offsetInSeconds := localTime.Zone()
	offsetInHours := (offsetInSeconds / 60) / 60
	require.Equal(-7, offsetInHours)

	require.Equal(from, TimeToISO8601DateTimeString(localTime))
	require.Equal("2006-01-02T22:04:05Z", TimeToISO8601DateTimeString(utcTime))

	require.Equal(utcTime.In(localTime.Location()), localTime)
	require.Equal(localTime.UTC(), utcTime)

	// ...a positive UTC offset/timezone:

	from = "2006-01-02T22:04:05+07:00"

	localTime, err = ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}
	utcTime, err = ISO8601StringToUTCTime(from)
	if err != nil {
		t.Fatal(err)
	}

	_, offsetInSeconds = localTime.Zone()
	offsetInHours = (offsetInSeconds / 60) / 60
	require.Equal(7, offsetInHours)

	require.Equal(from, TimeToISO8601DateTimeString(localTime))
	require.Equal("2006-01-02T15:04:05Z", TimeToISO8601DateTimeString(utcTime))

	require.Equal(utcTime.In(localTime.Location()), localTime)
	require.Equal(localTime.UTC(), utcTime)

	// ...and finally a positive UTC offset that shall "flip the date" when converted to UTC:

	from = "2006-01-02T01:04:05+04:00"

	localTime, err = ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}
	utcTime, err = ISO8601StringToUTCTime(from)
	if err != nil {
		t.Fatal(err)
	}

	_, offsetInSeconds = localTime.Zone()
	offsetInHours = (offsetInSeconds / 60) / 60
	require.Equal(4, offsetInHours)

	require.Equal(from, TimeToISO8601DateTimeString(localTime))
	require.Equal("2006-01-01T21:04:05Z", TimeToISO8601DateTimeString(utcTime))

	require.Equal(utcTime.In(localTime.Location()), localTime)
	require.Equal(localTime.UTC(), utcTime)
}

func Test_DateTimeZone(t *testing.T) {

	t.Run("Empty timeZone", func(t *testing.T) {
		tm := time.Date(2012, 4, 21, 11, 30, 0, 0, time.UTC)
		dt := &dtpb.DateTime{
			Year:       int32(tm.Year()),
			Month:      int32(tm.Month()),
			Day:        int32(tm.Day()),
			Hours:      int32(tm.Hour()),
			Minutes:    int32(tm.Minute()),
			Seconds:    int32(tm.Second()),
			TimeOffset: &dtpb.DateTime_TimeZone{TimeZone: &dtpb.TimeZone{Id: ""}},
		}

		pt, err := ProtoDateTimeToTime(dt)
		require.Nil(t, err)
		require.Equalf(t, tm, pt, "empty TimeOffset should result in UTC")

		dt.TimeOffset = &dtpb.DateTime_TimeZone{TimeZone: nil}
		pt, err = ProtoDateTimeToTime(dt)
		require.Nil(t, err)
		require.Equalf(t, tm, pt, "nil TimeZone should result in UTC")

		dt.TimeOffset = nil
		pt, err = ProtoDateTimeToTime(dt)
		require.Nil(t, err)
		require.Equalf(t, tm, pt, "nil TimeOffset should result in UTC")

	})

	t.Run("Unknown timeZone", func(t *testing.T) {
		dt := &dtpb.DateTime{
			Year:    int32(2012),
			Month:   int32(4),
			Day:     int32(21),
			Hours:   int32(11),
			Minutes: int32(30),
			Seconds: int32(0),
		}

		dt.TimeOffset = &dtpb.DateTime_TimeZone{TimeZone: &dtpb.TimeZone{Id: "fake/ness"}}
		pt, err := ProtoDateTimeToTime(dt)
		require.NotNilf(t, err, "should result in error")
		require.True(t, pt.IsZero(), "should result in zero time")
	})

	// supported happy path
	for _, test := range []struct {
		name               string
		y, mo, d, h, mi, s int
		tz                 *dtpb.TimeZone
		offset             *durpb.Duration
	}{
		{"DateTimeTZ", 2012, 4, 21, 11, 30, 0, &dtpb.TimeZone{Id: "America/New_York"}, nil},
		{"DateTimeTZ", 2012, 4, 21, 11, 30, 0, &dtpb.TimeZone{Id: "Europe/Berlin"}, nil},
		{"DateTimeTZ", 2012, 4, 21, 11, 30, 0, nil, &durpb.Duration{Seconds: 3600 * 5}},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Get the starting object.
			dt := &dtpb.DateTime{
				Year:    int32(test.y),
				Month:   int32(test.mo),
				Day:     int32(test.d),
				Hours:   int32(test.h),
				Minutes: int32(test.mi),
				Seconds: int32(test.s),
			}
			if test.tz != nil {
				dt.TimeOffset = &dtpb.DateTime_TimeZone{TimeZone: test.tz}
			}
			if test.offset != nil {
				dt.TimeOffset = &dtpb.DateTime_UtcOffset{UtcOffset: test.offset}
			}

			// Convert to a time.Time.
			tm, err := ProtoDateTimeToTime(dt)
			require.Nil(t, err)
			t.Run("ToTime", func(t *testing.T) {
				require.Equal(t, tm.Year(), test.y)
				require.Equal(t, tm.Month(), time.Month(test.mo))
				require.Equal(t, tm.Day(), test.d)
				require.Equal(t, tm.Hour(), test.h)
				require.Equal(t, tm.Minute(), test.mi)
				require.Equal(t, tm.Second(), test.s)
				if test.tz != nil {
					require.Equal(t, tm.Location().String(), test.tz.GetId())
				}
				if test.offset != nil {
					require.Equal(
						t,
						tm.Location().String(),
						fmt.Sprintf("UTC+%d", test.offset.GetSeconds()/3600),
					)
				}
			})

			// Convert back to a duration.
			t.Run("ToDateTime", func(t *testing.T) {
				durPb := TimeToProtoDateTime(tm)
				require.Equal(t, durPb.GetYear(), int32(test.y))
				require.Equal(t, durPb.GetMonth(), int32(test.mo))
				require.Equal(t, durPb.GetDay(), int32(test.d))
				require.Equal(t, durPb.GetHours(), int32(test.h))
				require.Equal(t, durPb.GetMinutes(), int32(test.mi))
				require.Equal(t, durPb.GetSeconds(), int32(test.s))
				if test.tz != nil {
					require.Equal(t, durPb.GetTimeZone().GetId(), test.tz.GetId())
				}
				if test.offset != nil {
					require.Equal(t, durPb.GetUtcOffset().GetSeconds(), test.offset.GetSeconds())
				}
			})
		})
	}
}

func Test_Date(t *testing.T) {
	for _, test := range []struct {
		name    string
		y, m, d int
	}{
		{"NormalDate", 2012, 4, 21},
		{"LongAgo", 1776, 7, 4},
		{"Future", 2032, 4, 21},
		{"FarFuture", 2062, 4, 21},
		{"MinimumDate", 1, 1, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			datePb := &dpb.Date{Year: int32(test.y), Month: int32(test.m), Day: int32(test.d)}
			local, err := ProtoDateToLocalTime(datePb)
			require.Nil(t, err)

			utc, err := ProtoDateToUTCTime(datePb)
			require.Nil(t, err)

			times := map[string]time.Time{
				"local": local,
				"utc":   utc,
			}
			for k, time := range times {
				t.Run(k, func(t *testing.T) {
					require.Equalf(t, time.Year(), test.y, "year")
					require.Equalf(t, int(time.Month()), test.m, "month")
					require.Equalf(t, time.Day(), test.d, "day")
					dtPb := TimeToProtoDate(time)
					t.Run("ToProto", func(t *testing.T) {
						require.Equalf(t, dtPb.GetYear(), int32(test.y), "year")
						require.Equalf(t, dtPb.GetMonth(), int32(test.m), "month")
						require.Equalf(t, dtPb.GetDay(), int32(test.d), "day")
					})
				})
			}
		})
	}

	// Invalid proto values
	datePb := &dpb.Date{Year: int32(0), Month: int32(0), Day: int32(0)}
	_, err := ProtoDateToLocalTime(datePb)
	require.Error(t, ErrInvalidValue, err)

	// Invalid proto values
	datePb = &dpb.Date{}
	_, err = ProtoDateToLocalTime(datePb)
	require.Error(t, ErrInvalidValue, err)

	// location panic error
	proto := TimeToProtoDate(time.Time{})
	require.Panics(t, func() { ProtoDateToTime(proto, nil) })

	// proto error
	_, err = ProtoDateToTime(nil, time.UTC)
	require.Error(t, err, ErrInvalidValue)

	// "empty" date
	proto = TimeToProtoDate(time.Time{})
	date, err := ProtoDateToTime(proto, time.UTC)
	require.Nil(t, err)
	require.Equal(t, time.Time{}, date)
	require.Equal(t, true, date.IsZero())

}

func Test_DateTime(t *testing.T) {
	from := "2006-01-02T01:04:05+04:00"
	localTime, err := ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}

	proto := TimeToProtoDateTime(localTime)
	require.Nilf(t, err, "converting from time to datetime proto")

	apTime, err := ProtoDateTimeToTime(proto)
	require.Nilf(t, err, "converting from datetime proto to time")

	// localTime have no "name" set for loc
	// should be equal if set to the apTime loc (UTC as default)
	require.Equal(t, localTime.In(apTime.Location()), apTime)

	require.Equal(t, localTime.UTC(), apTime.UTC())
	require.Equal(t, localTime.Local(), apTime.Local())

	proto = TimeToProtoDateTime(time.Time{})
	require.Nilf(t, err, "converting from time to datetime proto")

	apTime, err = ProtoDateTimeToTime(proto)
	require.Nilf(t, err, "converting from datetime proto to time")

	require.Equal(t, time.Time{}, apTime)
	require.Equal(t, true, apTime.IsZero())

	_, err = ProtoDateTimeToTime(nil)
	require.NotNil(t, err)
}

func Test_ISO8601StringToTime_date(t *testing.T) {

	require := require.New(t)

	utcTime, err := ISO8601StringToUTCTime("2006-01-02T22:04:05+07:00")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal("2006-01-02", TimeToISO8601DateString(utcTime))

	utcTime, err = ISO8601StringToUTCTime("2006-01-02T01:04:05+04:00")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal("2006-01-01", TimeToISO8601DateString(utcTime)) // flip the date...
}

func Test_TimeToLocalISO8601DateTimeString(t *testing.T) {

	require := require.New(t)

	utcTime, err := ISO8601StringToUTCTime("2006-01-05T11:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	localTime, err := ISO8601StringToTime("2006-01-02T22:04:05+07:00")
	if err != nil {
		t.Fatal(err)
	}

	_, offsetInSeconds := localTime.Zone()
	offsetInHours := (offsetInSeconds / 60) / 60
	require.Equal(7, offsetInHours)

	require.Equal(
		"2006-01-05T18:04:05+07:00",
		TimeToLocalISO8601DateTimeString(utcTime, localTime.Location()),
	)
}

func Test_IsSameDate(t *testing.T) {

	require := require.New(t)

	utcTime1, err := ISO8601StringToUTCTime("2006-01-02")
	if err != nil {
		t.Fatal(err)
	}

	utcTime2, err := ISO8601StringToUTCTime("2006-01-02T22:04:05+07:00")
	if err != nil {
		t.Fatal(err)
	}

	utcTime3, err := ISO8601StringToUTCTime("2006-01-03")
	if err != nil {
		t.Fatal(err)
	}

	require.True(IsSameDate(utcTime1, utcTime2))
	require.False(IsSameDate(utcTime1, utcTime3))
}

func Test_UTCTimeAdjustedToStartOfDay(t *testing.T) {

	require := require.New(t)

	utcTime, err := ISO8601StringToUTCTime("2006-01-02T11:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	require.False(IsStartOfDay(utcTime))
	require.False(IsEndOfDay(utcTime))

	utcTime = UTCTimeAdjustedToStartOfDay(utcTime)
	require.Equal("2006-01-02T00:00:00Z", TimeToISO8601DateTimeString(utcTime))

	require.True(IsStartOfDay(utcTime))
	require.False(IsEndOfDay(utcTime))
}

func Test_UTCTimeAdjustedToEndOfDay(t *testing.T) {

	require := require.New(t)

	utcTime, err := ISO8601StringToUTCTime("2006-01-02T11:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	require.False(IsStartOfDay(utcTime))
	require.False(IsEndOfDay(utcTime))

	utcTime = UTCTimeAdjustedToEndOfDay(utcTime)
	require.Equal("2006-01-02T23:59:59Z", TimeToISO8601DateTimeString(utcTime))

	require.False(IsStartOfDay(utcTime))
	require.True(IsEndOfDay(utcTime))
}

func Test_TimeToTimeOfDay(t *testing.T) {
	require := require.New(t)
	utcTime, err := ISO8601StringToUTCTime("2006-01-02T11:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	timeOfDay := TimeToTimeOfDay(utcTime)
	require.Equal(int32(11), timeOfDay.GetHours())
	require.Equal(int32(4), timeOfDay.GetMinutes())
	require.Equal(int32(5), timeOfDay.GetSeconds())
	require.Equal(int32(0), timeOfDay.GetNanos())
}

func Test_TimeOfDayToTime(t *testing.T) {
	require := require.New(t)
	timeOfDay := &tod.TimeOfDay{
		Hours:   11,
		Minutes: 4,
		Seconds: 5,
	}
	utcTime := TimeOfDayToTime(timeOfDay)
	require.Equal("0000-01-01T11:04:05Z", TimeToISO8601DateTimeString(utcTime))
}
