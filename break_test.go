package main

import (
	"fmt"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	now := time.Now()
	loc := now.Location()
	tests := []struct {
		start time.Time
		end   time.Time
		want  string
	}{
		{
			time.Date(2020, time.December, now.Day(), 15, 2, 0, 0, loc),
			time.Date(2020, time.December, now.Day(), 15, 8, 0, 0, loc),
			"15:02 - 15:08 (6 minutes)",
		},
		{
			time.Date(2020, time.December, now.Day()-1, 23, 2, 0, 0, loc),
			time.Date(2020, time.December, now.Day(), 3, 8, 0, 0, loc),
			fmt.Sprintf("%v Dec 23:02 - %v Dec 03:08 (4 hours 6 minutes)", now.Day()-1, now.Day()),
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v, %v", tt.start, tt.end)
		t.Run(testname, func(t *testing.T) {
			entry := breakEntry{tt.start, tt.end}
			answer := entry.string()
			if answer != tt.want {
				t.Fatalf("got %v, want %v", answer, tt.want)
			}
		})
	}
}
