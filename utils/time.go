package utils

import "time"

func HourToSecond(hour int) int {
	return hour * 3600
}

func TimeNowToMillisSecond() int64 {
	return time.Now().UnixNano() / 1000000
}