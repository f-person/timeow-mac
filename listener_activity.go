package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

func (a *app) activityListener() {
	for {
		activity := <-a.notifierCh
		fmt.Printf("(%v) new event: %v\n", time.Now(), activity.Type)
		switch activity.Type {
		case notifier.Sleep:
			fmt.Println("Sleeping")
			a.isSleeping = true
		case notifier.Awake:
			now := time.Now()
			totalIdleTime := calculateDuration(a.lastActiveTime, now)

			fmt.Println("##########")
			fmt.Printf("totalIdleTime = %v, now = %v, lastActiveTime = %v.\n", totalIdleTime, now, a.lastActiveTime)
			fmt.Println("##########")

			if totalIdleTime > a.maxAllowedIdleTime {
				a.addActivePeriodEntry(a.lastIdleTime, a.lastActiveTime)

				a.addBreakEntry(a.lastActiveTime, now)
				a.lastIdleTime = time.Now()
			}

			a.checkAndCleanExpiredTimeLogs()

			a.isSleeping = false
		}
	}
}
