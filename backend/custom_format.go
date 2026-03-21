package backend

import (
	"fmt"
	"strings"
)

// ApplyCustomFormat replaces tokens in template with values from result.
func ApplyCustomFormat(template string, r DurationResult, zeroPad bool) string {
	if template == "" {
		return ""
	}
	dd := fmtIntWidth(r.Days, zeroPad, 2)
	hh := fmtIntWidth(r.Hours, zeroPad, 2)
	mm := fmtIntWidth(r.Minutes, zeroPad, 2)
	ss := fmtIntWidth(r.Seconds, zeroPad, 2)
	mss := fmtIntWidth(r.Ms, zeroPad, 3)
	replacements := []struct {
		token string
		val   string
	}{
		{"{total_seconds}", fmt.Sprintf("%d", r.TotalMs/1000)},
		{"{total_ms}", fmt.Sprintf("%d", r.TotalMs)},
		{"{dd}", dd},
		{"{hh}", hh},
		{"{mm}", mm},
		{"{ss}", ss},
		{"{ms}", mss},
	}
	out := template
	for _, p := range replacements {
		out = strings.ReplaceAll(out, p.token, p.val)
	}
	return out
}
