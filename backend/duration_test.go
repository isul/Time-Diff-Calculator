package backend

import (
	"testing"
	"time"
)

func TestCalculateDurationBreakdown(t *testing.T) {
	t.Setenv("TZ", "UTC")
	loc := time.UTC
	start := time.Date(2026, 3, 21, 0, 0, 0, 0, loc)
	end := time.Date(2026, 3, 22, 1, 2, 3, 4_000_000, loc)
	r := CalculateDuration(start, end)
	if r.Days != 1 || r.Hours != 1 || r.Minutes != 2 || r.Seconds != 3 || r.Ms != 4 {
		t.Fatalf("%+v", r)
	}
	if r.TotalMs != end.Sub(start).Milliseconds() {
		t.Fatalf("total ms")
	}
}

func TestCalculateDurationZero(t *testing.T) {
	t.Setenv("TZ", "UTC")
	loc := time.UTC
	tm := time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
	r := CalculateDuration(tm, tm)
	if r.TotalMs != 0 || r.Days != 0 {
		t.Fatalf("%+v", r)
	}
}
