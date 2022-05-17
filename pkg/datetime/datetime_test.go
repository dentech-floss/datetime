package datetime

import (
	"testing"

	"github.com/stretchr/testify/require"
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
