package main

import (
	"time"

	"github.com/getlantern/systray"
)

func (a *app) handleIdleItemSelected(mIdleTimes []*systray.MenuItem, index int) {
	prevIndex := getMinutesSliceIndexFromDuration(a.maxAllowedIdleTime, idleTimeOptionsInSettings[:])
	if prevIndex >= 0 && prevIndex < len(mIdleTimes) {
		mIdleTimes[prevIndex].Uncheck()
	}
	mIdleTimes[index].Check()

	a.setMaxAllowedIdleTime(int(idleTimeOptionsInSettings[index]))
}

func (a *app) handleKeepTimeLogsForOptionSelected(mKeepTimeLogsForOptions []*systray.MenuItem, index int) {
	prevIndex := getMinutesSliceIndexFromDuration(a.keepTimeLogsFor, keepTimeLogsForOptionsInSettings[:])
	if prevIndex >= 0 && prevIndex < len(mKeepTimeLogsForOptions) {
		mKeepTimeLogsForOptions[prevIndex].Uncheck()
	}
	mKeepTimeLogsForOptions[index].Check()

	a.setKeepTimeLogsFor(int(keepTimeLogsForOptionsInSettings[index]))

	a.checkAndCleanExpiredTimeLogs()
}

func (a *app) handleOpenAtLoginClicked(item *systray.MenuItem) {
	if a.startup.RunningAtStartup() {
		a.startup.RemoveStartupItem()
		item.Uncheck()
	} else {
		a.startup.AddStartupItem()
		item.Check()
	}
}

func (a *app) handleQuitClicked() {
	a.addActivePeriodEntry(a.lastIdleTime, time.Now())

	systray.Quit()
}
