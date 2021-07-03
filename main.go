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
	startup            startup.Startup

	defaults userdefaults.UserDefaults

	isIdle         bool
	isSleeping     bool
	lastIdleTime   time.Time
	lastActiveTime time.Time

	notifier *notifier.Notifier
	ticker   *time.Ticker

	idleTimeCh chan time.Duration
	notifierCh chan *notifier.Activity

	breaks         []breakEntry
	mBreaks        *systray.MenuItem
	breakMenuItems []*systray.MenuItem
}

func main() {
	notifierInstance := notifier.GetInstance()
	defaults := *userdefaults.Defaults()
	app := app{
		defaults: defaults,
		startup: startup.Startup{
			AppLabel: appLabel,
			AppName:  appName,
		},

		maxAllowedIdleTime: time.Minute * time.Duration(defaults.Integer(maxAllowedIdleTimeKey)),

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
		app.setMaxAllowedIdleTime(int(defaultMaxAllowedIdleTime))
	}

	breaks, err := app.readBreaksFromStorage()
	if err == nil {
		app.breaks = breaks
	} else {
		fmt.Printf("An error occurred while reading breaks: %v", err)
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

	mBreaks := systray.AddMenuItem("Breaks", "")
	mNoBreaks := mBreaks.AddSubMenuItem("No breaks yet", "")
	mNoBreaks.Disable()

	a.mBreaks = mBreaks
	a.breakMenuItems = append(a.breakMenuItems, mNoBreaks)

	if len(a.breaks) > 0 {
		a.updateBreakMenuItems()
	}

	mPreferences := systray.AddMenuItem("Preferences", "")
	mGoIdleAfter := mPreferences.AddSubMenuItem("Reset after inactivity for", "")
	mOpenAtLogin := mPreferences.AddSubMenuItemCheckbox("Start at Login", "", a.startup.RunningAtStartup())

	var mIdleTimes [len(idleTimes)]*systray.MenuItem
	selectedIdleTimeIndex := getIdleTimeIndexFromDuration(a.maxAllowedIdleTime)
	for index, minutes := range idleTimes {
		durationString := durafmt.Parse(time.Duration(minutes) * time.Minute).LimitFirstN(1).String()
		mIdleTimes[index] = mGoIdleAfter.AddSubMenuItemCheckbox(durationString, "", index == selectedIdleTimeIndex)
	}

	idleTimeSelected := make(chan int)
	for i, mIdleTime := range mIdleTimes {
		go func(c chan struct{}, index int) {
			for range c {
				idleTimeSelected <- index
			}
		}(mIdleTime.ClickedCh, i)
	}

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
				a.handleIdleItemSelected(mIdleTimes[:], index)
			}
		}
	}()
}

func (a *app) activityListener() {
	for activity := range a.notifierCh {
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
			a.lastIdleTime = time.Now()
			a.isIdle = true
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
