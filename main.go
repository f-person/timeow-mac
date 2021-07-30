package main

import (
	"fmt"
	"time"

	"github.com/f-person/usage_time_menubar_app/pkg/startup"
	"github.com/f-person/usage_time_menubar_app/pkg/userdefaults"
	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
	"github.com/lextoumbourou/idle"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

type app struct {
	maxAllowedIdleTime time.Duration
	keepTimeLogsFor    time.Duration

	startup  startup.Startup
	defaults userdefaults.UserDefaults

	isIdle     bool
	isSleeping bool

	lastIdleTime   time.Time
	lastActiveTime time.Time

	notifier *notifier.Notifier
	ticker   *time.Ticker

	idleTimeCh chan time.Duration
	notifierCh chan *notifier.Activity

	breaks         []period
	mBreaks        *systray.MenuItem
	breakMenuItems []*systray.MenuItem

	activePeriods         []period
	mActivePeriods        *systray.MenuItem
	activePeriodMenuItems []*systray.MenuItem
}

func main() {
	fmt.Println("Started the app at", time.Now())

	notifierInstance := notifier.GetInstance()
	defaults := *userdefaults.Defaults()
	app := app{
		maxAllowedIdleTime: time.Minute * time.Duration(defaults.Integer(maxAllowedIdleTimeKey)),
		keepTimeLogsFor:    time.Minute * time.Duration(defaults.Integer(keepTimeLogsForKey)),

		defaults: defaults,
		startup: startup.Startup{
			AppLabel: appLabel,
			AppName:  appName,
		},

		isIdle:         false,
		isSleeping:     false,
		lastIdleTime:   time.Now(),
		lastActiveTime: time.Now(),
		notifier:       notifierInstance,
		ticker:         time.NewTicker(idleListenerInterval),

		idleTimeCh: make(chan time.Duration),
		notifierCh: notifierInstance.Start(),
	}

	// [maxAllowedIdleTime] have never been set before,
	// use the default value and save it to user defaults.
	if app.maxAllowedIdleTime == 0 {
		app.setMaxAllowedIdleTime(int(defaultMaxAllowedIdleTime.Minutes()))
	}

	if app.keepTimeLogsFor == 0 {
		app.setKeepTimeLogsFor(int(defaultKeepTimeLogsFor.Minutes()))
	}

	breaks, err := app.readPeriodsFromStorage(breaksKey)
	if err == nil {
		app.breaks = breaks
	} else {
		fmt.Printf("An error occurred while reading breaks: %v", err)
	}

	activePeriods, err := app.readPeriodsFromStorage(activePeriodsKey)
	if err == nil {
		app.activePeriods = activePeriods
	} else {
		fmt.Printf("An error occurred while reading active times: %v", err)
	}

	go app.activityListener()
	go app.idleTimeListener(app.ticker)

	systray.Run(func() {
		app.onSystrayReady()
	}, func() {
		notifierInstance.Quit()
	})
}

func (a *app) onSystrayReady() {
	systray.SetTitle("0m")

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

	mQuit := systray.AddMenuItem("Quit", "")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				a.handleQuitClicked()
			case <-mOpenAtLogin.ClickedCh:
				a.handleOpenAtLoginClicked(mOpenAtLogin)
			case index := <-idleTimeSelected:
				a.handleIdleItemSelected(mIdleTimes, index)
			case index := <-keepTimeLogsForSelected:
				a.handleKeepTimeLogsForOptionSelected(mKeepTimeLogsForOptions, index)
			}
		}
	}()
}

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

			a.isSleeping = false
		}
	}
}

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
