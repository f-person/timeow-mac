package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/lextoumbourou/idle"
)

func (a *app) idleTimeListener(ticker *time.Ticker) {
	for range ticker.C {
		if a.isSleeping {
			continue
		}

		idleTime, err := idle.Get()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if idleTime > a.maxAllowedIdleTime {
			if !a.isIdle {
				a.addActivePeriodEntry(a.lastIdleTime, a.lastActiveTime)

				a.isIdle = true
			}

			a.lastIdleTime = time.Now()
		} else if a.isIdle {
			a.addBreakEntry(a.lastActiveTime, time.Now())

			// Reset
			a.lastIdleTime = time.Now()
			a.lastActiveTime = time.Now()
			a.isIdle = false
		} else {
			a.lastActiveTime = time.Now().Add(-idleTime)
		}

		d := time.Since(a.lastIdleTime)
		fmt.Println(d)
		systray.SetTitle(formatDuration(d))
	}
}
