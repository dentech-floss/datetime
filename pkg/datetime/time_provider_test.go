package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_UTCTimeProvider(t *testing.T) {
	rq := require.New(t)

	timeProvider := NewUTCTimeProvider()
	rq.Equal(time.UTC, timeProvider.Now().Location())
}

func Test_UTCTimeProvider_Since(t *testing.T) {
	rq := require.New(t)

	timeProvider := NewUTCTimeProvider()
	start := time.Now().Add(-50 * time.Millisecond)
	elapsed := timeProvider.Since(start)
	rq.GreaterOrEqual(elapsed, 50*time.Millisecond)
}

func Test_FakeTimeProvider_Since(t *testing.T) {
	rq := require.New(t)

	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	start := now.Add(-2 * time.Hour)

	timeProvider := NewFakeTimeProvider(now)
	rq.Equal(2*time.Hour, timeProvider.Since(start))
}
