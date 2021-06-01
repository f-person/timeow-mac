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

	if days > 0 {
		b.WriteString(fmt.Sprintf("%dd ", days))
	}
	if hours > 0 || (days == 0 && minutes == 0) {
		b.WriteString(fmt.Sprintf("%dh ", hours))
	}
	if minutes > 0 || hours == 0 {
		b.WriteString(fmt.Sprintf("%dm", minutes))
	}

	return b.String()
}
