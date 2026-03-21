package backend

import "testing"

func TestCalculateHappyPath(t *testing.T) {
	t.Setenv("TZ", "UTC")
	s := DefaultSettings()
	s.Language = "ko"
	resp := Calculate("2026-03-21 00:00:00 2026-03-21 01:00:00", s)
	if !resp.Ok {
		t.Fatalf("%+v", resp)
	}
	if len(resp.Results) < 6 {
		t.Fatalf("results %d", len(resp.Results))
	}
}

func TestCalculateCustomFormat(t *testing.T) {
	t.Setenv("TZ", "UTC")
	s := DefaultSettings()
	s.CustomFormat = "{hh}:{mm}:{ss}"
	s.Language = "en"
	resp := Calculate("2026-03-21 00:00:00 2026-03-21 00:30:00", s)
	if !resp.Ok {
		t.Fatal(resp.Error)
	}
	var found bool
	for _, r := range resp.Results {
		if r.Key == "custom" && r.Value == "00:30:00" {
			found = true
		}
	}
	if !found {
		t.Fatalf("%+v", resp.Results)
	}
}

func TestCalculateCustomFormatNoZeroPad(t *testing.T) {
	t.Setenv("TZ", "UTC")
	s := DefaultSettings()
	f := false
	s.ZeroPadding = &f
	s.CustomFormat = "{hh}:{mm}:{ss}"
	s.Language = "en"
	resp := Calculate("2026-03-21 00:00:00 2026-03-21 00:30:00", s)
	if !resp.Ok {
		t.Fatal(resp.Error)
	}
	var found bool
	for _, r := range resp.Results {
		if r.Key == "custom" {
			found = true
			if r.Value != "0:30:0" {
				t.Fatalf("got %q", r.Value)
			}
		}
	}
	if !found {
		t.Fatal("missing custom result")
	}
}

func TestValidateInputEmpty(t *testing.T) {
	s := DefaultSettings()
	s.Language = "en"
	v := ValidateInput("   ", s)
	if v.Ready || v.Level != "warn" {
		t.Fatalf("%+v", v)
	}
}
