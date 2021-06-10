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

	app.notifierCh = app.notifier.Start()

	go app.activityListener()
	go app.idleTimeListener(app.ticker)

	systray.Run(func() {
		app.onSystrayReady()
	}, func() {
		notifierInstance.Quit()
	})
}

func (a *app) setMaxAllowedIdleTime(minutes int) {
	a.maxAllowedIdleTime = time.Duration(minutes) * time.Minute
	a.defaults.SetInteger(maxAllowedIdleTimeKey, minutes)
}

func (a *app) onSystrayReady() {
	systray.SetTitle("0m")

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

func getIdleTimeIndexFromDuration(d time.Duration) int {
	minutes := uint8(d.Minutes())
	for index, value := range idleTimes {
		if value == minutes {
			return index
		}
	}

	return -1
}

func (a *app) activityListener() {
	for activity := range a.notifierCh {
		// TODO: update [a.lastActiveTime] and [a.lastIdleTime]
		switch activity.Type {
		case notifier.Sleep:
			a.isSleeping = true
		case notifier.Awake:
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
		} else {
			a.lastActiveTime = time.Now().Add(-idleTime)

			if idleTime > a.maxAllowedIdleTime {
				a.lastIdleTime = time.Now()
				a.isIdle = true
			} else if a.isIdle {
				// Reset
				a.lastIdleTime = time.Now()
				a.lastActiveTime = time.Now()
				a.isIdle = false

				// TODO: add idleDuration to durations list
			}

			d := time.Since(a.lastIdleTime)
			fmt.Println(d)
			systray.SetTitle(formatDuration(d))
		}
	}
}

func (a *app) handleIdleItemSelected(mIdleTimes []*systray.MenuItem, index int) {
	prevIndex := getIdleTimeIndexFromDuration(a.maxAllowedIdleTime)
	mIdleTimes[prevIndex].Uncheck()
	mIdleTimes[index].Check()

	a.setMaxAllowedIdleTime(int(idleTimes[index]))
}

func (a *app) handleQuitClicked() {
	systray.Quit()
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
