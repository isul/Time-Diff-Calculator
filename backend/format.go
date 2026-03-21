package backend

import (
	"fmt"
	"strconv"
	"strings"
)

func commaInt64(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	lead := len(s) % 3
	if lead == 0 {
		lead = 3
	}
	b.WriteString(s[:lead])
	for i := lead; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

func fmtIntWidth(v int, zeroPad bool, width int) string {
	if !zeroPad {
		return strconv.Itoa(v)
	}
	switch width {
	case 2:
		return fmt.Sprintf("%02d", v)
	case 3:
		return fmt.Sprintf("%03d", v)
	default:
		return fmt.Sprintf("%0*d", width, v)
	}
}

// FormatPredefined renders one built-in format key.
func FormatPredefined(key string, r DurationResult, locale string, zeroPad bool) string {
	switch key {
	case "seconds":
		return formatSecondsStr(r, locale)
	case "milliseconds":
		return formatMilliseconds(r)
	case "mmss":
		return formatMMSS(r, locale, zeroPad)
	case "hhmmss":
		return formatHHMMSS(r, locale, zeroPad)
	case "ddhhmmss":
		return formatDDHHMMSS(r, locale, zeroPad)
	case "full":
		return formatFull(r, locale, zeroPad)
	default:
		return ""
	}
}

func formatSecondsStr(r DurationResult, locale string) string {
	sec := r.TotalMs / 1000
	if locale == "ko" {
		return fmt.Sprintf("%d%s", sec, tr(locale, "unit_second"))
	}
	u := enPlural(locale, int(sec), "unit_second", "unit_seconds")
	return fmt.Sprintf("%d %s", sec, u)
}

func formatMilliseconds(r DurationResult) string {
	return fmt.Sprintf("%sms", commaInt64(r.TotalMs))
}

func formatMMSS(r DurationResult, locale string, zeroPad bool) string {
	totalSec := r.TotalMs / 1000
	mm := totalSec / 60
	ss := totalSec % 60
	ms := fmtIntWidth(int(mm), zeroPad, 2)
	sss := fmtIntWidth(int(ss), zeroPad, 2)
	if locale == "ko" {
		return fmt.Sprintf("%s%s %s%s", ms, tr(locale, "unit_minute"), sss, tr(locale, "unit_second"))
	}
	return fmt.Sprintf("%s %s %s %s", ms, enPlural(locale, int(mm), "unit_minute", "unit_minutes"), sss, enPlural(locale, int(ss), "unit_second", "unit_seconds"))
}

func formatHHMMSS(r DurationResult, locale string, zeroPad bool) string {
	totalSec := r.TotalMs / 1000
	h := totalSec / 3600
	rem := totalSec % 3600
	m := rem / 60
	s := rem % 60
	hs := fmtIntWidth(int(h), zeroPad, 2)
	ms := fmtIntWidth(int(m), zeroPad, 2)
	ss := fmtIntWidth(int(s), zeroPad, 2)
	if locale == "ko" {
		return fmt.Sprintf("%s%s %s%s %s%s", hs, tr(locale, "unit_hour"), ms, tr(locale, "unit_minute"), ss, tr(locale, "unit_second"))
	}
	return fmt.Sprintf("%s %s %s %s %s %s",
		hs, enPlural(locale, int(h), "unit_hour", "unit_hours"),
		ms, enPlural(locale, int(m), "unit_minute", "unit_minutes"),
		ss, enPlural(locale, int(s), "unit_second", "unit_seconds"),
	)
}

func formatDDHHMMSS(r DurationResult, locale string, zeroPad bool) string {
	d := r.Days
	h := r.Hours
	m := r.Minutes
	s := r.Seconds
	ds := fmtIntWidth(d, zeroPad, 2)
	hs := fmtIntWidth(h, zeroPad, 2)
	ms := fmtIntWidth(m, zeroPad, 2)
	ss := fmtIntWidth(s, zeroPad, 2)
	if locale == "ko" {
		return fmt.Sprintf("%s%s %s%s %s%s %s%s", ds, tr(locale, "unit_day"), hs, tr(locale, "unit_hour"), ms, tr(locale, "unit_minute"), ss, tr(locale, "unit_second"))
	}
	return fmt.Sprintf("%s %s %s %s %s %s %s %s",
		ds, enPlural(locale, d, "unit_day", "unit_days"),
		hs, enPlural(locale, h, "unit_hour", "unit_hours"),
		ms, enPlural(locale, m, "unit_minute", "unit_minutes"),
		ss, enPlural(locale, s, "unit_second", "unit_seconds"),
	)
}

func formatFull(r DurationResult, locale string, zeroPad bool) string {
	base := formatDDHHMMSS(r, locale, zeroPad)
	mss := fmtIntWidth(r.Ms, zeroPad, 3)
	if locale == "ko" {
		return fmt.Sprintf("%s %sms", base, mss)
	}
	return fmt.Sprintf("%s %s ms", base, mss)
}
