package helper

import (
	"strings"
	"time"
)

func CalculateDueBy(priority string) *time.Time {
	var duration time.Duration
	switch strings.ToLower(priority) {
	case "high":
		duration = 60 * time.Minute
	case "medium":
		duration = 90 * time.Minute
	case "low":
		duration = 120 * time.Minute
	case "very_low":
		duration = 240 * time.Minute
	default:
		duration = 90 * time.Minute
	}
	due := time.Now().Add(duration)
	return &due
}
