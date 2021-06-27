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

func (b *breakEntry) duration() time.Duration {
	return b.end.Sub(b.start)
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
