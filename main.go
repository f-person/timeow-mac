package main

import (
	"fmt"
	"time"

	"github.com/f-person/usage_time_menubar_app/pkg/startup"
	"github.com/getlantern/systray"
	"github.com/lextoumbourou/idle"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
)

var appStartup = startup.Startup{
	AppLabel: appLabel,
	AppName:  appName,
}

var isIdle = false
var lastIdleTime = time.Now()
var lastActiveTime = time.Now()

func main() {
	_ = lastActiveTime // TODO: delete the line after using [lastActiveTime]

	notifierInstance := notifier.GetInstance()

	notifierCh := notifierInstance.Start()

	go func() {
		for activity := range notifierCh {
			fmt.Println(activity.Type, time.Now())
		}
	}()

	systray.Run(func() {
		onSystrayReady()
	}, func() {
		notifierInstance.Quit()
	})
}

func onSystrayReady() {
	systray.SetTitle("0m")

	idleTimeCh := make(chan time.Duration)
	// TODO: stop the listener when going to sleep

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	go idleTimeListener(ticker, done, idleTimeCh)

	mPreferences := systray.AddMenuItem("Preferences", "")
	mOpenAtLogin := mPreferences.AddSubMenuItemCheckbox("Start at Login", "", appStartup.RunningAtStartup())

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				handleQuitClicked()
			case <-mOpenAtLogin.ClickedCh:
				handleOpenAtLoginClicked(mOpenAtLogin)
			}
		}
	}()
}

func idleTimeListener(ticker *time.Ticker, done chan bool, idleTimeCh chan time.Duration) {
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			idleTime, err := idle.Get()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				// need to keep [lastIdleTime] to subtract it from [lastActiveTime] instead of [idleTime]
				lastActiveTime = time.Now().Add(-idleTime)

				if idleTime > maxAllowedIdleTime {
					lastIdleTime = time.Now()
					isIdle = true
				} else if isIdle {
					// Reset
					lastIdleTime = time.Now()
					lastActiveTime = time.Now()
					isIdle = false
					// need to handle going to sleep

					// TODO: add idleDuration to durations list
					// IDEA: timer every hour check and delete old idle breaks
				}

				d := time.Since(lastIdleTime)
				fmt.Println(d)
				systray.SetTitle(formatDuration(d))
			}

		}
	}
}

func handleQuitClicked() {
	systray.Quit()
}

func handleOpenAtLoginClicked(item *systray.MenuItem) {
	if appStartup.RunningAtStartup() {
		appStartup.RemoveStartupItem()
		item.Uncheck()
	} else {
		appStartup.AddStartupItem()
		item.Check()
	}
}
