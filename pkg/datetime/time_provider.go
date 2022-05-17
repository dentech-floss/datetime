package datetime

import (
	"time"
)

type TimeProvider interface {
	// Get the current time
	Now() time.Time
}

type utcTimeProvider struct{}

func (*utcTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

func NewUTCTimeProvider() TimeProvider {
	return &utcTimeProvider{}
}

type fakeTimeProvider struct {
	now time.Time
}

func (f *fakeTimeProvider) Now() time.Time {
	return f.now
}

func NewFakeTimeProvider(now time.Time) TimeProvider {
	return &fakeTimeProvider{now: now}
}
