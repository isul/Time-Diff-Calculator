package backend

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// datetimeTokenRe matches one date+time token (local, no timezone).
var datetimeTokenRe = regexp.MustCompile(
	`(?:\d{4}[-/.]\d{1,2}[-/.]\d{1,2}|\d{1,2}[-/.]\d{1,2})\s+\d{1,2}:\d{1,2}:\d{1,2}(?:\.\d{1,3})?`,
)

var reTimeSuffix = regexp.MustCompile(
	`^(.+)\s+(\d{1,2}):(\d{1,2}):(\d{1,2})(?:\.(\d{1,3}))?$`,
)

func padMillis(msStr string) int {
	if msStr == "" {
		return 0
	}
	for len(msStr) < 3 {
		msStr += "0"
	}
	if len(msStr) > 3 {
		msStr = msStr[:3]
	}
	n, _ := strconv.Atoi(msStr)
	if n < 0 {
		return 0
	}
	if n > 999 {
		return 999
	}
	return n
}

func parseDatePart(datePart string, defaultYear int) (y, mo, d int, err error) {
	datePart = strings.TrimSpace(datePart)
	if datePart == "" {
		return 0, 0, 0, fmt.Errorf("empty date")
	}
	// YYYY-MM-DD (exactly two dashes, year first)
	if c := strings.Count(datePart, "-"); c == 2 {
		p := strings.Split(datePart, "-")
		if len(p) == 3 && len(p[0]) == 4 {
			y, _ = strconv.Atoi(p[0])
			mo, _ = strconv.Atoi(p[1])
			d, _ = strconv.Atoi(p[2])
			if y >= 1000 && mo >= 1 && mo <= 12 && d >= 1 && d <= 31 {
				return y, mo, d, nil
			}
		}
	}
	// YYYY/MM/DD
	if c := strings.Count(datePart, "/"); c == 2 {
		p := strings.Split(datePart, "/")
		if len(p) == 3 && len(p[0]) == 4 {
			y, _ = strconv.Atoi(p[0])
			mo, _ = strconv.Atoi(p[1])
			d, _ = strconv.Atoi(p[2])
			if y >= 1000 && mo >= 1 && mo <= 12 && d >= 1 && d <= 31 {
				return y, mo, d, nil
			}
		}
	}
	// YYYY.MM.DD
	if c := strings.Count(datePart, "."); c == 2 {
		p := strings.Split(datePart, ".")
		if len(p) == 3 && len(p[0]) == 4 {
			y, _ = strconv.Atoi(p[0])
			mo, _ = strconv.Atoi(p[1])
			d, _ = strconv.Atoi(p[2])
			if y >= 1000 && mo >= 1 && mo <= 12 && d >= 1 && d <= 31 {
				return y, mo, d, nil
			}
		}
	}
	// MM-DD / MM/DD / MM.DD (no year)
	seps := []struct {
		parts []string
	}{
		{strings.Split(datePart, "-")},
		{strings.Split(datePart, "/")},
		{strings.Split(datePart, ".")},
	}
	for _, sp := range seps {
		if len(sp.parts) == 2 {
			mo, _ = strconv.Atoi(sp.parts[0])
			d, _ = strconv.Atoi(sp.parts[1])
			if mo >= 1 && mo <= 12 && d >= 1 && d <= 31 {
				return defaultYear, mo, d, nil
			}
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid date: %q", datePart)
}

func parseOneDateTime(s string, defaultYear int, loc *time.Location) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty datetime")
	}
	m := reTimeSuffix.FindStringSubmatch(s)
	if m == nil {
		return time.Time{}, fmt.Errorf("invalid datetime: %q", s)
	}
	datePart := m[1]
	h, _ := strconv.Atoi(m[2])
	mi, _ := strconv.Atoi(m[3])
	sec, _ := strconv.Atoi(m[4])
	ms := padMillis(m[5])

	y, mo, d, err := parseDatePart(datePart, defaultYear)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(y, time.Month(mo), d, h, mi, sec, ms*1_000_000, loc), nil
}

// ParseInput extracts two datetimes from free text (whitespace-separated tokens).
func ParseInput(input string) (ParsedInput, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return ParsedInput{}, fmt.Errorf("missing input")
	}
	matches := datetimeTokenRe.FindAllString(input, -1)
	if len(matches) < 2 {
		return ParsedInput{}, fmt.Errorf("need two date-times; found %d", len(matches))
	}
	aStr, bStr := matches[0], matches[1]
	loc := time.Local
	now := time.Now().In(loc)
	year := now.Year()

	start, err := parseOneDateTime(aStr, year, loc)
	if err != nil {
		return ParsedInput{}, fmt.Errorf("start: %w", err)
	}
	end, err := parseOneDateTime(bStr, year, loc)
	if err != nil {
		return ParsedInput{}, fmt.Errorf("end: %w", err)
	}
	if start.After(end) {
		start, end = end, start
	}
	return ParsedInput{Start: start, End: end}, nil
}

// CountDatetimeTokens returns how many datetime-like tokens are in input.
func CountDatetimeTokens(input string) int {
	return len(datetimeTokenRe.FindAllString(strings.TrimSpace(input), -1))
}
