package main

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
)

type breakEntry struct {
	start time.Time
	end   time.Time
}

// Break times are converted to Unix and the difference in seconds is calculated
// because time.Sub behaves strangely after the machine wakes up from sleep.
func (b *breakEntry) duration() time.Duration {
	startSeconds := b.start.Unix()
	endSeconds := b.end.Local().Unix()

	return time.Second * time.Duration(endSeconds-startSeconds)
}

func (b *breakEntry) string() string {
	var format string

	if b.start.Day() == time.Now().Day() {
		format = "15:04"
	} else {
		format = "2 Jan 15:04"
	}
	duration := b.duration()
	limit := 1
	if duration > time.Hour {
		limit = 2
	}

	return fmt.Sprintf(
		"%s - %s (%s)",
		b.start.Format(format),
		b.end.Format(format),
		durafmt.Parse(duration).LimitFirstN(limit).String(),
	)
}
