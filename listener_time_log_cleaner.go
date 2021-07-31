package main

import (
	"fmt"
	"time"
)

func (a *app) timeLogCleaner(ticker *time.Ticker) {
	for range ticker.C {
		a.checkAndCleanExpiredTimeLogs()
	}
}

func (a *app) checkAndCleanExpiredTimeLogs() {
	cleanExpiredPeriods := func(periods []period) []period {
		minimalAllowedDate := time.Now().Add(-a.keepTimeLogsFor)

		index := 0
		for ; index < len(periods); index++ {
			entry := periods[index]
			if entry.End.After(minimalAllowedDate) {
				fmt.Printf(
					"Found a date that is after minimalAllowedDate (%v) at index %v – %v; Breaking\n",
					minimalAllowedDate,
					index,
					entry.string(),
				)
				break
			}
		}

		return periods[index:]
	}

	a.breaks = cleanExpiredPeriods(a.breaks)
	a.activePeriods = cleanExpiredPeriods(a.activePeriods)

	a.savePeriodsToStorage(breaksKey, a.breaks)
	a.savePeriodsToStorage(activePeriodsKey, a.activePeriods)

	if a.isSystrayReady {
		a.updateBreakMenuItems()
		a.updateActivePeriodMenuItems()
	}
}
