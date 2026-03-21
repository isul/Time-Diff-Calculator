package backend

import "testing"

func TestFormatPredefinedKo(t *testing.T) {
	r := DurationResult{TotalMs: 2000, Days: 0, Hours: 0, Minutes: 0, Seconds: 2, Ms: 0}
	if s := FormatPredefined("seconds", r, "ko", true); s != "2초" {
		t.Fatalf("got %q", s)
	}
	if s := FormatPredefined("milliseconds", r, "ko", true); s != "2,000ms" {
		t.Fatalf("got %q", s)
	}
}

func TestFormatPredefinedEnPlural(t *testing.T) {
	r1 := DurationResult{TotalMs: 1000, Days: 0, Hours: 0, Minutes: 0, Seconds: 1, Ms: 0}
	if s := FormatPredefined("seconds", r1, "en", true); s != "1 second" {
		t.Fatalf("got %q", s)
	}
	r2 := DurationResult{TotalMs: 2000, Days: 0, Hours: 0, Minutes: 0, Seconds: 2, Ms: 0}
	if s := FormatPredefined("seconds", r2, "en", true); s != "2 seconds" {
		t.Fatalf("got %q", s)
	}
}

func TestFormatPredefinedNoZeroPad(t *testing.T) {
	r := DurationResult{TotalMs: 125000, Days: 0, Hours: 0, Minutes: 2, Seconds: 5, Ms: 7}
	if s := FormatPredefined("mmss", r, "en", false); s != "2 minutes 5 seconds" {
		t.Fatalf("got %q", s)
	}
	if s := FormatPredefined("hhmmss", r, "ko", false); s != "0시간 2분 5초" {
		t.Fatalf("got %q", s)
	}
}

func TestApplyCustomFormat(t *testing.T) {
	r := DurationResult{TotalMs: 3661001, Days: 0, Hours: 1, Minutes: 1, Seconds: 1, Ms: 1}
	out := ApplyCustomFormat("{total_ms} {total_seconds} {dd}/{hh}/{mm}/{ss}/{ms}", r, true)
	if out != "3661001 3661 00/01/01/01/001" {
		t.Fatalf("got %q", out)
	}
}

func TestApplyCustomFormatNoPad(t *testing.T) {
	r := DurationResult{TotalMs: 3661001, Days: 0, Hours: 1, Minutes: 1, Seconds: 1, Ms: 1}
	out := ApplyCustomFormat("{dd}/{hh}/{mm}/{ss}/{ms}", r, false)
	if out != "0/1/1/1/1" {
		t.Fatalf("got %q", out)
	}
}
