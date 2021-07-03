package main

import (
	"time"
)

func getIdleTimeIndexFromDuration(d time.Duration) int {
	minutes := uint8(d.Minutes())
	for index, value := range idleTimes {
		if value == minutes {
			return index
		}
	}

	return -1
}

func (a *app) setMaxAllowedIdleTime(minutes int) {
	a.maxAllowedIdleTime = time.Duration(minutes) * time.Minute
	a.defaults.SetInteger(maxAllowedIdleTimeKey, minutes)
}

// This function exists because time.Sub behaves strangely after the machine wakes up from sleep.
func calculateDuration(start, end time.Time) time.Duration {
	startSeconds := start.Unix()
	endSeconds := end.Local().Unix()

	return time.Second * time.Duration(endSeconds-startSeconds)
}
