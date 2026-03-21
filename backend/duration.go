package backend

import "time"

// CalculateDuration converts wall-clock difference into decomposed components.
func CalculateDuration(start, end time.Time) DurationResult {
	if end.Before(start) {
		start, end = end, start
	}
	totalMs := end.Sub(start).Milliseconds()
	if totalMs < 0 {
		totalMs = 0
	}
	const (
		msPerMinute = 60 * 1000
		msPerHour   = 60 * msPerMinute
		msPerDay    = 24 * msPerHour
	)
	days := int(totalMs / msPerDay)
	rem := totalMs % msPerDay
	hours := int(rem / msPerHour)
	rem %= msPerHour
	minutes := int(rem / msPerMinute)
	rem %= msPerMinute
	seconds := int(rem / 1000)
	ms := int(rem % 1000)
	return DurationResult{
		TotalMs: totalMs,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
		Ms:      ms,
	}
}
