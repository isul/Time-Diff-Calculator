package backend

import (
	"strings"
)

var formatOrder = []struct {
	key     string
	labelTr string
}{
	{"seconds", "format_seconds"},
	{"milliseconds", "format_milliseconds"},
	{"mmss", "format_mmss"},
	{"hhmmss", "format_hhmmss"},
	{"ddhhmmss", "format_ddhhmmss"},
	{"full", "format_full"},
}

// ValidateInput checks input for real-time UI feedback.
func ValidateInput(input string, settings Settings) ValidateResponse {
	locale := ResolveLocale(settings.Language)
	input = strings.TrimSpace(input)
	if input == "" {
		return ValidateResponse{Ready: false, Level: "warn", Message: tr(locale, "warn_empty")}
	}
	n := CountDatetimeTokens(input)
	if n < 2 {
		return ValidateResponse{Ready: false, Level: "error", Message: tr(locale, "err_need_two")}
	}
	return ValidateResponse{Ready: true, Level: "ok", Message: tr(locale, "ok_ready")}
}

// Calculate parses input, computes duration, and formats according to settings.
func Calculate(input string, settings Settings) CalculateResponse {
	settings = MergeSettings(settings)
	locale := ResolveLocale(settings.Language)
	resp := CalculateResponse{Locale: locale}

	in := strings.TrimSpace(input)
	if in == "" {
		resp.Ok = false
		resp.Warning = tr(locale, "warn_empty")
		return resp
	}

	parsed, err := ParseInput(in)
	if err != nil {
		resp.Ok = false
		resp.Error = err.Error()
		return resp
	}

	dur := CalculateDuration(parsed.Start, parsed.End)
	resp.Ok = true
	pad := EffectiveZeroPadding(settings)

	for _, fo := range formatOrder {
		if !settings.Formats[fo.key] {
			continue
		}
		val := FormatPredefined(fo.key, dur, locale, pad)
		resp.Results = append(resp.Results, FormatLine{
			Key:   fo.key,
			Label: tr(locale, fo.labelTr),
			Value: val,
		})
	}

	if settings.Formats["custom"] && strings.TrimSpace(settings.CustomFormat) != "" {
		resp.Results = append(resp.Results, FormatLine{
			Key:   "custom",
			Label: tr(locale, "format_custom"),
			Value: ApplyCustomFormat(settings.CustomFormat, dur, pad),
		})
	}

	return resp
}
