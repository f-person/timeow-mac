package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
)

type period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (p *period) duration() time.Duration {
	return calculateDuration(p.Start, p.End)
}

func (p *period) string() string {
	var format string

	if p.Start.Day() == time.Now().Day() {
		format = "15:04"
	} else {
		format = "2 Jan 15:04"
	}
	duration := p.duration()
	limit := 1
	if duration > time.Hour {
		limit = 2
	}

	return fmt.Sprintf(
		"%s - %s (%s)",
		p.Start.Format(format),
		p.End.Format(format),
		durafmt.Parse(duration).LimitFirstN(limit).String(),
	)
}

func (a *app) readPeriodsFromStorage(key string) ([]period, error) {
	var periods []period
	err := a.defaults.Unmarshal(key, &periods)

	return periods, err
}

func (a *app) savePeriodsToStorage(key string, periods []period) error {
	return a.defaults.Marshal(key, periods)
}

func updatePeriodMenuItems(
	periods []period,
	periodsMenuItem *systray.MenuItem,
	currentMenuItems []*systray.MenuItem,
) []*systray.MenuItem {
	menuItems := currentMenuItems

	totalNewMenuItems := len(periods) - len(menuItems)
	fmt.Printf("totalNewMenuItems: %v", totalNewMenuItems)

	if totalNewMenuItems < 0 {
		// Hide redundant menu items.
		for i := len(menuItems) - 1; i >= len(menuItems)-(-totalNewMenuItems); i-- {
			menuItems[i].Hide()
		}
	} else {
		// Add missing menu items.
		for i := 0; i < totalNewMenuItems; i++ {
			item := periodsMenuItem.AddSubMenuItem("", "")
			item.Disable()
			item.Show()
			menuItems = append(menuItems, item)
		}
	}

	length := len(periods)
	for index, entry := range periods {
		menuItems[length-index-1].SetTitle(entry.string())
	}

	return menuItems
}
