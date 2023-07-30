package timeutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/cybersamx/golib/timeutils"
)

func TestTime_DurationSeconds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		timestamp time.Duration
		want      int
	}{
		{time.Nanosecond, 0},
		{time.Microsecond, 0},
		{time.Millisecond, 0},
		{time.Second, 1},
		{time.Minute, 60},
		{time.Hour, 3600},
	}

	for _, tc := range tests {
		seconds := ToSeconds(tc.timestamp)
		assert.Equal(t, tc.want, seconds)
	}
}

func TestCompareTimeAndTimeString(t *testing.T) {
	timestamp := time.Unix(0, 0)

	est, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	pst, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	tests := []struct {
		compareTime string
		loc         *time.Location
		want        bool
	}{
		{"Wed, 31 Dec 1969 19:00:00 EST", est, true},
		{"Wed, 31 Dec 1969 16:00:00 PST", pst, true},
		{"Thu, 01 Jan 1970 00:00:00 UTC", time.UTC, true},
		{"Thu, 1 Jan 1970 00:00:00 UTC", time.UTC, false},
		{"1970-01-01T00:00:00Z", time.UTC, false},
		{"abc123", time.UTC, false},
	}

	for _, tc := range tests {
		if tc.want {
			assert.True(t, CompareTimeAndTimeString(timestamp, tc.compareTime, tc.loc))
		} else {
			assert.False(t, CompareTimeAndTimeString(timestamp, tc.compareTime, tc.loc))
		}
	}
}
