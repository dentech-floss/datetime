package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	dpb "google.golang.org/genproto/googleapis/type/date"
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

func Test_Date(t *testing.T) {
	for _, test := range []struct {
		name    string
		y, m, d int
	}{
		{"NormalDate", 2012, 4, 21},
		{"LongAgo", 1776, 7, 4},
		{"Future", 2032, 4, 21},
		{"FarFuture", 2062, 4, 21},
	} {
		t.Run(test.name, func(t *testing.T) {
			datePb := &dpb.Date{Year: int32(test.y), Month: int32(test.m), Day: int32(test.d)}
			times := map[string]time.Time{
				"local": ProtoDateToLocalTime(datePb),
				"utc":   ProtoDateToUTCTime(datePb),
			}
			for k, time := range times {
				t.Run(k, func(t *testing.T) {
					require.Equalf(t, time.Year(), test.y, "year")
					require.Equalf(t, int(time.Month()), test.m, "month")
					require.Equalf(t, time.Day(), test.d, "day")
					dtPb := TimeToProtoDate(&time)
					t.Run("ToProto", func(t *testing.T) {
						require.Equalf(t, dtPb.GetYear(), int32(test.y), "year")
						require.Equalf(t, dtPb.GetMonth(), int32(test.m), "month")
						require.Equalf(t, dtPb.GetDay(), int32(test.d), "day")
					})
				})
			}
		})
	}
}

func Test_DateTime(t *testing.T) {
	from := "2006-01-02T01:04:05+04:00"
	localTime, err := ISO8601StringToTime(from)
	if err != nil {
		t.Fatal(err)
	}

	proto, err := TimeToProtoDateTime(localTime)
	require.Nilf(t, err, "converting from time to datetime proto")

	apTime, err := ProtoDateTimeToTime(proto)
	require.Nilf(t, err, "converting from datetime proto to time")

	// localTime have no "name" set for loc
	// should be equal if set to the apTime loc (UTC as default)
	require.Equal(t, localTime.In(apTime.Location()), apTime)

	require.Equal(t, localTime.UTC(), apTime.UTC())
	require.Equal(t, localTime.Local(), apTime.Local())
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

	require.Equal("2006-01-05T18:04:05+07:00", TimeToLocalISO8601DateTimeString(utcTime, localTime.Location()))
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
