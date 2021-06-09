package main

import "time"

const idleListenerInterval = time.Second
const appLabel = "com.github.f-person.usage_time_menubar_app"
const appName = "Usage Time Menu Bar App"
const defaultMaxAllowedIdleTime time.Duration = time.Minute * 3

const maxAllowedIdleTimeKey = "maxAllowedIdleTime"
