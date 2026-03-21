package backend

import (
	"testing"
	"time"
)

func TestParseInputShortDates(t *testing.T) {
	t.Setenv("TZ", "UTC")
	y := time.Now().UTC().Year()
	pi, err := ParseInput("03-21 00:59:24 03-21 01:06:28")
	if err != nil {
		t.Fatal(err)
	}
	if pi.Start.Year() != y || pi.Start.Month() != 3 || pi.Start.Day() != 21 {
		t.Fatalf("start date: %v", pi.Start)
	}
	if pi.End.Sub(pi.Start) != 7*time.Minute+4*time.Second {
		t.Fatalf("delta: %v", pi.End.Sub(pi.Start))
	}
}

func TestParseInputReverseOrder(t *testing.T) {
	t.Setenv("TZ", "UTC")
	pi, err := ParseInput("2026-03-21 10:00:00 2026-03-21 08:00:00")
	if err != nil {
		t.Fatal(err)
	}
	if !pi.Start.Before(pi.End) {
		t.Fatal("expected swap")
	}
	if pi.End.Sub(pi.Start) != 2*time.Hour {
		t.Fatalf("delta %v", pi.End.Sub(pi.Start))
	}
}

func TestParseInputMillis(t *testing.T) {
	t.Setenv("TZ", "UTC")
	pi, err := ParseInput("2026/03/21 00:00:00.001 2026/03/21 00:00:00.101")
	if err != nil {
		t.Fatal(err)
	}
	if pi.End.Sub(pi.Start) != 100*time.Millisecond {
		t.Fatalf("delta %v", pi.End.Sub(pi.Start))
	}
}

func TestParseInputNeedTwo(t *testing.T) {
	_, err := ParseInput("2026-03-21 12:00:00")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCountDatetimeTokens(t *testing.T) {
	if n := CountDatetimeTokens("a 2026-03-21 12:00:00 b 2026-03-21 13:00:00"); n != 2 {
		t.Fatalf("got %d", n)
	}
}
