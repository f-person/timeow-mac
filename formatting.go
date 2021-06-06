package main

import (
	"fmt"
	"strings"
	"time"
)

const day = time.Hour * 24

func formatDuration(d time.Duration) string {
	var b strings.Builder

	days := d / day
	hours := (d - (days * day)) / time.Hour
	minutes := (d - (days * day) - (hours * time.Hour)) / time.Minute

	// Maybe days are redundant?
	if days > 0 {
		b.WriteString(fmt.Sprintf("%dd ", days))
	}
	if hours > 0 {
		b.WriteString(fmt.Sprintf("%dh ", hours))
	}
	b.WriteString(fmt.Sprintf("%dm", minutes))

	return b.String()
}
