package helper

import (
	"fmt"
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

func FormatDuration(d time.Duration) string {
	if d.Hours() >= 24 {
		days := int(d.Hours()) / 24
		hours := int(d.Hours()) % 24
		if hours > 0 {
			return fmt.Sprintf("%d days %d hours", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	} else if d.Hours() >= 1 {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes > 0 {
			return fmt.Sprintf("%d hours %d minutes", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	} else if d.Minutes() >= 1 {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds > 0 {
			return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
		}
		return fmt.Sprintf("%d minutes", minutes)
	} else {
		seconds := int(d.Seconds())
		return fmt.Sprintf("%d seconds", seconds)
	}
}
