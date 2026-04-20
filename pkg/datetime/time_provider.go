package datetime

import (
	"time"
)

type TimeProvider interface {
	// Now - Get the current time
	Now() time.Time
	// Since - Get the duration since start
	Since(start time.Time) time.Duration
}

type utcTimeProvider struct{}

func (*utcTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

func (*utcTimeProvider) Since(start time.Time) time.Duration {
	return time.Since(start)
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

func (f *fakeTimeProvider) Since(start time.Time) time.Duration {
	return f.now.Sub(start)
}

func NewFakeTimeProvider(now time.Time) TimeProvider {
	return &fakeTimeProvider{now: now}
}
