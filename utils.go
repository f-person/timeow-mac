package main

import "time"

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
