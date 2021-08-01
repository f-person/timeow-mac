package main

import (
	"fmt"
	"time"

	"github.com/f-person/usage_time_menubar_app/pkg/startup"
	"github.com/f-person/usage_time_menubar_app/pkg/userdefaults"
	"github.com/getlantern/systray"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

type app struct {
	isPro bool

	maxAllowedIdleTime time.Duration
	keepTimeLogsFor    time.Duration

	startup  startup.Startup
	defaults userdefaults.UserDefaults

	isIdle     bool
	isSleeping bool

	lastIdleTime   time.Time
	lastActiveTime time.Time

	notifier *notifier.Notifier

	isSystrayReady bool
	idleTimeCh     chan time.Duration
	notifierCh     chan *notifier.Activity

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
		isPro: false,

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

	// Listeners
	go app.activityListener()

	idleTimeListenerTicker := time.NewTicker(idleListenerInterval)
	go app.idleTimeListener(idleTimeListenerTicker)

	timeLogCleanerTicker := time.NewTicker(timeLogCleanerInterval)
	go app.timeLogCleaner(timeLogCleanerTicker)

	app.checkAndCleanExpiredTimeLogs()

	// Menu bar
	systray.Run(func() {
		app.onSystrayReady()
	}, func() {
		notifierInstance.Quit()
	})
}
