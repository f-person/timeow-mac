package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
	"github.com/lextoumbourou/idle"
)

const maxAllowedIdleTime time.Duration = time.Second * 10

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("0m")

	mIdleTime := systray.AddMenuItem("Haven't been idle yet", "")
	mIdleTime.Disable()

	mOpenAtLogin := systray.AddMenuItemCheckbox("Open at Login", "", true)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	go idleTimeListener(mIdleTime)

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				handleQuit()
			case <-mOpenAtLogin.ClickedCh:
				handleOpenAtLoginClicked(mOpenAtLogin)
			}
		}
	}()
}

func idleTimeListener(mIdleTime *systray.MenuItem) {
	for range time.Tick(time.Second * 1) {
		idleTime, err := idle.Get()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			mIdleTime.SetTitle(durafmt.Parse(idleTime).LimitFirstN(2).String())

			if idleTime > maxAllowedIdleTime {
				fmt.Println("is taking a break")
				// TODO: handle going idle
			}
		}
	}
}

func onExit() {
	fmt.Println("onExit()")
}

func handleQuit() {
	systray.Quit()
}

func handleOpenAtLoginClicked(item *systray.MenuItem) {
	if item.Checked() {
		item.Uncheck()
	} else {
		item.Check()
	}
}
