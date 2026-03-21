package backend

import "time"

// ParsedInput holds parsed start/end instants (local timezone).
type ParsedInput struct {
	Start time.Time
	End   time.Time
}

// DurationResult is the decomposed duration in local wall-clock components.
type DurationResult struct {
	TotalMs int64
	Days    int
	Hours   int
	Minutes int
	Seconds int
	Ms      int
}

// Settings persisted and sent from the UI.
type Settings struct {
	Formats map[string]bool `json:"formats"`
	// CustomFormat is optional; if non-empty and enabled, rendered as extra line.
	CustomFormat string `json:"customFormat"`
	// Language: "auto" | "ko" | "en"
	Language string `json:"language"`
	// ZeroPadding pads hh/mm/ss/ms (and related tokens) with leading zeros. Omitted in JSON means default true.
	ZeroPadding *bool `json:"zeroPadding,omitempty"`
}

// FormatLine is one output row.
type FormatLine struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}

// CalculateResponse is returned to the frontend.
type CalculateResponse struct {
	Ok       bool         `json:"ok"`
	Error    string       `json:"error,omitempty"`
	Warning  string       `json:"warning,omitempty"`
	Results  []FormatLine `json:"results,omitempty"`
	Locale   string       `json:"locale"`
}

// ValidateResponse for real-time checks without full calculation.
type ValidateResponse struct {
	Ready   bool   `json:"ready"`
	Level   string `json:"level"` // "ok" | "warn" | "error"
	Message string `json:"message,omitempty"`
}
