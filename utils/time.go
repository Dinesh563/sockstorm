package utils

import (
	"time"
)

func TimeDifference(t1 time.Time, t2 time.Time) int {
	s1 := (t1.Hour() - t2.Hour()) * 60 * 60
	s2 := (t1.Minute() - t2.Minute()) * 60
	s3 := t1.Second() - t2.Second()
	return s1 + s2 + s3
}

// formatTimestamp converts UNIX timestamp to human-readable IST format
func FormatTimestamp(timestamp int32) time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	if timestamp > 0 {
		return time.Unix(int64(timestamp), 0).In(loc)
	}
	return time.Now().Add(-12 * time.Hour)
}
