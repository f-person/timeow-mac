package main

import "time"

const idleListenerInterval = time.Second
const timeLogCleanerInterval = time.Minute * 30
const appLabel = "com.timeow.timeow-mac"
const appName = "Timeow"
const defaultMaxAllowedIdleTime time.Duration = time.Minute * 3
const defaultKeepTimeLogsFor time.Duration = time.Hour * 24
const minAllowedActiveTime time.Duration = time.Minute

// User Default keys
const maxAllowedIdleTimeKey = "maxAllowedIdleTime"
const keepTimeLogsForKey = "keepTimeLogsFor"
const breaksKey = "breaks"
const activePeriodsKey = "activePeriods"

const hourInMinutes = 60
const dayInMinutes = hourInMinutes * 24

var idleTimeOptionsInSettings = [...]uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var keepTimeLogsForOptionsInSettings = [...]uint32{
	hourInMinutes,
	2 * hourInMinutes,
	3 * hourInMinutes,
	4 * hourInMinutes,
	5 * hourInMinutes,
	10 * hourInMinutes,
	12 * hourInMinutes,
	dayInMinutes,
	2 * dayInMinutes,
	3 * dayInMinutes,
	7 * dayInMinutes,
	30 * dayInMinutes,
}

const getProURL = "https://github.com/f-person"
