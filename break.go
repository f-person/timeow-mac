package main

import (
	"fmt"
	"time"
)

func (a *app) addBreakEntry(start, end time.Time) {
	entry := period{
		Start: start,
		End:   end,
	}

	// Avoid duplicate breaks.
	length := len(a.breaks)
	if length > 0 {
		lastBreak := a.breaks[length-1]
		if entry.Start.Truncate(time.Minute) == lastBreak.Start.Truncate(time.Minute) {
			if entry.End.After(lastBreak.End) {
				a.breaks[length-1] = entry
			}
		} else {
			a.breaks = append(a.breaks, entry)
		}
	} else {
		a.breaks = append(a.breaks, entry)
	}

	fmt.Println("----------")
	fmt.Printf("New break added: start = %v, end =%v , duration = %v.\n", entry.Start, entry.End, entry.duration())
	fmt.Printf("string: %v\n", entry.string())
	fmt.Println("current time:", time.Now())
	fmt.Println("----------")

	a.updateBreakMenuItems()
	a.savePeriodsToStorage(breaksKey, a.breaks)
}

func (a *app) updateBreakMenuItems() {
	a.breakMenuItems = updatePeriodMenuItems(a.breaks, a.mBreaks, a.breakMenuItems)
}
