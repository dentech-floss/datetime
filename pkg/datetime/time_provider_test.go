package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_UTCTimeProvider(t *testing.T) {
	require := require.New(t)

	timeProvider := NewUTCTimeProvider()
	require.Equal(time.UTC, timeProvider.Now().Location())
}
