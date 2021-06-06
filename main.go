package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/lextoumbourou/idle"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
	"github.com/sqweek/dialog"
)

func main() {
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

	mQuit := systray.AddMenuItem("Quit", "")

	idleTimeCh := make(chan time.Duration)
	// TODO: stop the listener when going to sleep
	go idleTimeListener(idleTimeCh)

	mPreferences := systray.AddMenuItem("Preferences", "")
	mOpenAtLogin := mPreferences.AddSubMenuItemCheckbox("Open at Login", "", true)

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

func idleTimeListener(idleTimeCh chan time.Duration) {
	isIdle := false
	lastIdleTime := time.Now()
	lastActiveTime := time.Now()
	_ = lastActiveTime

	for range time.Tick(time.Second * 1) {
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

func handleQuitClicked() {
	answer := dialog.Message("Are you sure you want to quit the app?").YesNo()
	if answer {
		systray.Quit()
	}

}

func handleOpenAtLoginClicked(item *systray.MenuItem) {
	if item.Checked() {
		item.Uncheck()
	} else {
		item.Check()
	}
}
