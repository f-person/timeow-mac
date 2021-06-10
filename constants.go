package main

import "time"

const idleListenerInterval = time.Second
const appLabel = "com.github.f-person.usage_time_menubar_app"
const appName = "Usage Time Menu Bar App"
const defaultMaxAllowedIdleTime time.Duration = time.Minute * 3

const maxAllowedIdleTimeKey = "maxAllowedIdleTime"

var idleTimes = [...]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
