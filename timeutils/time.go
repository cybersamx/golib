package timeutils

import (
	"time"
)

func NowInMilli() int64 {
	return time.Now().UnixMilli()
}

func ToSeconds(duration time.Duration) int {
	const nanosecondMultiplier = time.Nanosecond * time.Microsecond * time.Millisecond
	return int(duration / nanosecondMultiplier)
}

// CompareTimeAndTimeString parses a RFC1123 string `compare` in timezone `loc` and then compares against
// a time.Time object `timestamp`. Returns true if they are equal and false otherwise.
func CompareTimeAndTimeString(timestamp time.Time, compare string, loc *time.Location) bool {
	compareTime, err := time.ParseInLocation(time.RFC1123, compare, loc)
	if err != nil {
		return false
	}

	return timestamp.Equal(compareTime)
}
