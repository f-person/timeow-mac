package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
)

func (a *app) onSystrayReady() {
	systray.SetTitle("0m")

	var getProClickedCh chan struct{}

	if !a.isPro {
		mGetPro := systray.AddMenuItem(fmt.Sprintf("⭐️ Get %s Pro", appName), "")
		getProClickedCh = mGetPro.ClickedCh
		systray.AddSeparator()
	}

	// Setup break entries
	mBreaks := systray.AddMenuItem("Breaks", "")
	mNoBreaks := mBreaks.AddSubMenuItem("No breaks yet", "")
	mNoBreaks.Disable()

	a.mBreaks = mBreaks
	a.breakMenuItems = append(a.breakMenuItems, mNoBreaks)
	if len(a.breaks) > 0 {
		a.updateBreakMenuItems()
	}

	// Setup active time entries
	mActivePeriods := systray.AddMenuItem("Active periods", "")
	mNoActivePeriods := mActivePeriods.AddSubMenuItem("No active periods yet", "")
	mNoActivePeriods.Disable()

	a.mActivePeriods = mActivePeriods
	a.activePeriodMenuItems = append(a.activePeriodMenuItems, mNoActivePeriods)
	if len(a.activePeriods) > 0 {
		a.updateActivePeriodMenuItems()
	}

	mPreferences := systray.AddMenuItem("Preferences", "")
	mGoIdleAfter := mPreferences.AddSubMenuItem("Reset after inactivity for", "")
	mKeepTimeLogsFor := mPreferences.AddSubMenuItem("Keep time logs for", "")
	mOpenAtLogin := mPreferences.AddSubMenuItemCheckbox("Start at Login", "", a.startup.RunningAtStartup())

	createMinutesSelectionItems := func(
		menuItem *systray.MenuItem,
		selectedItem time.Duration,
		optionsInSettings []uint32,
	) (
		itemSelected chan int, menuItems []*systray.MenuItem,
	) {
		selectedItemIndex := getMinutesSliceIndexFromDuration(selectedItem, optionsInSettings[:])
		for index, minutes := range optionsInSettings {
			durationString := durafmt.Parse(time.Duration(minutes) * time.Minute).LimitFirstN(1).String()
			menuItems = append(
				menuItems,
				menuItem.AddSubMenuItemCheckbox(durationString, "", index == selectedItemIndex),
			)
		}

		itemSelected = make(chan int)
		for i, mItem := range menuItems {
			go func(c chan struct{}, index int) {
				for range c {
					itemSelected <- index
				}
			}(mItem.ClickedCh, i)
		}

		return itemSelected, menuItems
	}

	idleTimeSelected, mIdleTimes := createMinutesSelectionItems(
		mGoIdleAfter,
		a.maxAllowedIdleTime,
		idleTimeOptionsInSettings[:],
	)
	keepTimeLogsForSelected, mKeepTimeLogsForOptions := createMinutesSelectionItems(
		mKeepTimeLogsFor,
		a.keepTimeLogsFor,
		keepTimeLogsForOptionsInSettings[:],
	)

	systray.AddSeparator()

	mAbout := systray.AddMenuItem(fmt.Sprintf("About %s", appName), "")
	mQuit := systray.AddMenuItem(fmt.Sprintf("Quit %s", appName), "")

	go func() {
		for {
			select {
			case <-getProClickedCh:
				openURL(getProURL)
			case index := <-idleTimeSelected:
				a.handleIdleItemSelected(mIdleTimes, index)
			case index := <-keepTimeLogsForSelected:
				a.handleKeepTimeLogsForOptionSelected(mKeepTimeLogsForOptions, index)
			case <-mOpenAtLogin.ClickedCh:
				a.handleOpenAtLoginClicked(mOpenAtLogin)
			case <-mAbout.ClickedCh:
				a.addAnalyticsEvent("aboutClicked")
				a.handleAboutClicked()
			case <-mQuit.ClickedCh:
				a.handleQuitClicked()
			}
		}
	}()

	a.isSystrayReady = true
}
