package main

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
)

type breakEntry struct {
	start time.Time
	end   time.Time
}

func (b *breakEntry) duration() time.Duration {
	return calculateDuration(b.start, b.end)
}

func (b *breakEntry) string() string {
	var format string

	if b.start.Day() == time.Now().Day() {
		format = "15:04"
	} else {
		format = "2 Jan 15:04"
	}
	duration := b.duration()
	limit := 1
	if duration > time.Hour {
		limit = 2
	}

	return fmt.Sprintf(
		"%s - %s (%s)",
		b.start.Format(format),
		b.end.Format(format),
		durafmt.Parse(duration).LimitFirstN(limit).String(),
	)
}

func (a *app) addBreakEntry(start, end time.Time) {
	entry := breakEntry{
		start: start,
		end:   end,
	}

	// Avoid duplicate breaks.
	length := len(a.breaks)
	if length > 0 {
		lastBreak := a.breaks[length-1]
		if entry.start.Truncate(time.Minute) == lastBreak.start.Truncate(time.Minute) {
			if entry.end.After(lastBreak.end) {
				a.breaks[length-1] = entry
			}
		} else {
			a.breaks = append(a.breaks, entry)
		}
	} else {
		a.breaks = append(a.breaks, entry)
	}

	fmt.Println("----------")
	fmt.Printf("New break added: start = %v, end =%v , duration = %v.\n", entry.start, entry.end, entry.duration())
	fmt.Printf("string: %v\n", entry.string())
	fmt.Println("current time:", time.Now())
	fmt.Println("----------")

	if len(a.breakMenuItems) == len(a.breaks)-1 {
		item := a.mBreaks.AddSubMenuItem("", "")
		item.Disable()
		a.breakMenuItems = append(a.breakMenuItems, item)
	}

	length = len(a.breaks)
	for index, entry := range a.breaks {
		a.breakMenuItems[length-index-1].SetTitle(entry.string())
	}
}
