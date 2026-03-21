package backend

import "testing"

func TestMergeSettingsAddsMissingKeys(t *testing.T) {
	s := Settings{
		Formats: map[string]bool{
			"seconds": false,
		},
		Language: "auto",
	}
	s = MergeSettings(s)
	if !s.Formats["milliseconds"] {
		t.Fatal("expected default true for missing key")
	}
	if s.Formats["seconds"] != false {
		t.Fatal("expected preserved false")
	}
	if !EffectiveZeroPadding(s) {
		t.Fatal("expected default zero padding true when unset")
	}
}
