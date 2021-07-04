package main

import (
	"fmt"
	"time"
)

func (a *app) addActivePeriodEntry(start, end time.Time) {
	entry := period{
		Start: start,
		End:   end,
	}

	a.activePeriods = append(a.activePeriods, entry)

	fmt.Println("----------")
	fmt.Printf("New activePeriod added: start = %v, end =%v , duration = %v.\n", entry.Start, entry.End, entry.duration())
	fmt.Printf("string: %v\n", entry.string())
	fmt.Println("current time:", time.Now())
	fmt.Println("----------")

	a.updateActivePeriodMenuItems()
	a.savePeriodsToStorage(activePeriodsKey, a.activePeriods)
}

func (a *app) updateActivePeriodMenuItems() {
	a.activePeriodMenuItems = updatePeriodMenuItems(a.activePeriods, a.mActivePeriods, a.activePeriodMenuItems)
}
