package main

import (
	"fmt"

	ga "github.com/f-person/timeow-mac/pkg/ga"
)

func (a *app) addAnalyticsEvent(action string) {
	event := ga.NewEvent(analyticsCategory, action).Label("menuBarInteraction")
	err := a.analytics.Send(event)
	if err != nil {
		fmt.Printf("Error occurred when sending Analytics event: %v\n", err)
	} else {
		fmt.Printf("Successfully sent event %v\n", action)
	}
}
