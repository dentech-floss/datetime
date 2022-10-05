# datetime

Contains various date/time related utility functionality, like parsing and formatting ISO8601 strings to/from UTC as well as different timezones/locations. It also contains a "TimeProvider" which is useful for controlling the perception of time in a service, by injecting this we get a single and mockable source of the current time throughout the code base.

## Install

```
go get github.com/dentech-floss/datetime@v0.1.1
```

## Usage

### datetime

[datetime.go](https://github.com/dentech-floss/datetime/blob/main/pkg/datetime/datetime.go) contains a bunch of reusable utility func's for dealing with date/time to/from UTC, it uses [relvacode/iso8601](https://github.com/relvacode/iso8601) for dealing with ISO8601 formatted strings as well as proto DATA/DATETIME to/from UTC.

Here follows an example or it's usage, check out the [datetime_test.go](https://github.com/dentech-floss/datetime/blob/main/pkg/datetime/datetime_test.go) for the full monty.

```go
package example

import (
    "testing"
    "github.com/dentech-floss/datetime/pkg/datetime"
)

func Test_ISO8601StringToTime(t *testing.T) {

    from = "2006-01-02T15:04:05-07:00"

    localTime, err = datetime.ISO8601StringToTime(from)
    if err != nil {
        t.Fatal(err)
    }
    utcTime, err = datetime.ISO8601StringToUTCTime(from)
    if err != nil {
        t.Fatal(err)
    }

    _, offsetInSeconds := localTime.Zone()
    offsetInHours := (offsetInSeconds / 60) / 60
    require.Equal(-7, offsetInHours)

    require.Equal(from, datetime.TimeToISO8601DateTimeString(localTime))
    require.Equal("2006-01-02T22:04:05Z", datetime.TimeToISO8601DateTimeString(utcTime))

    require.Equal(utcTime.In(localTime.Location()), localTime)
    require.Equal(localTime.UTC(), utcTime)
}
```

### TimeProvider

Disadvantages of using "time.Now()" in the code? Well... are we using UTC or not? What if we use "time.Now()" somewhere when we were supposed to use "time.Now().UTC()"? What if we have time-sensitive code and want to write tests for certain times? 

Inject the mockable [TimeProvider](https://github.com/dentech-floss/datetime/blob/main/pkg/datetime/time_provider.go) and get a single source of the current time.

```go
package example

import (
    "github.com/dentech-floss/datetime/pkg/datetime"
)

func main() {
    timeProvider := datetime.NewUTCTimeProvider() // provider of current time in UTC
    patientGatewayServiceV1 := service.NewPatientGatewayServiceV1(timeProvider) // inject it
}
```

```go
package example

func (s *PatientGatewayServiceV1) FindAppointments(
    ctx context.Context,
    request *patient_gateway_service_v1.FindAppointmentsRequest,
) (*patient_gateway_service_v1.FindAppointmentsResponse, error) {

    now := s.timeProvider.Now() // Get the current time in UTC
}
```

Control the perception of time when testing:

```go
package example

import (
    "testing"
    "github.com/dentech-floss/datetime/pkg/datetime"
)

func Test_FindAppointments(t *testing.T) {

    now, err := datetime.ISO8601StringToUTCTime("2022-01-02T22:04:05+07:00")
    if err != nil {
        t.Fatal(err)
    }

    timeProvider := datetime.NewFakeTimeProvider(now) // provider of current time of choice
    patientGatewayServiceV1 := service.NewPatientGatewayServiceV1(timeProvider) // inject it

    patientGatewayServiceV1.FindAppointments(...)
}
```
