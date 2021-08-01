package main

import (
	"os/exec"
	"time"
)

func getMinutesSliceIndexFromDuration(d time.Duration, minutesSlice []uint32) int {
	minutes := uint32(d.Minutes())
	for index, value := range minutesSlice {
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

func (a *app) setKeepTimeLogsFor(minutes int) {
	a.keepTimeLogsFor = time.Duration(minutes) * time.Minute
	a.defaults.SetInteger(keepTimeLogsForKey, minutes)
}

// This function exists because time.Sub behaves strangely after the machine wakes up from sleep.
func calculateDuration(start, end time.Time) time.Duration {
	startSeconds := start.Unix()
	endSeconds := end.Local().Unix()

	return time.Second * time.Duration(endSeconds-startSeconds)
}

// opens [url] using the `open` command
func openURL(url string) error {
	return exec.Command("open", url).Start()
}
