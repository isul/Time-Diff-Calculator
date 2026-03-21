package backend

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
)

//go:embed locales/ko.json
var koJSON []byte

//go:embed locales/en.json
var enJSON []byte

type messageMap map[string]string

func loadMessages(raw []byte) messageMap {
	var m messageMap
	_ = json.Unmarshal(raw, &m)
	return m
}

var (
	koMsgs = loadMessages(koJSON)
	enMsgs = loadMessages(enJSON)
)

// ResolveLocale returns "ko" or "en" from override or environment.
func ResolveLocale(override string) string {
	o := strings.TrimSpace(strings.ToLower(override))
	if o == "ko" || o == "en" {
		return o
	}
	for _, k := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		v := os.Getenv(k)
		if v == "" {
			continue
		}
		lv := strings.ToLower(v)
		if strings.HasPrefix(lv, "ko") {
			return "ko"
		}
		if strings.HasPrefix(lv, "en") {
			return "en"
		}
	}
	return "en"
}

func messagesFor(locale string) messageMap {
	if locale == "ko" {
		return koMsgs
	}
	return enMsgs
}

func tr(locale, key string) string {
	m := messagesFor(locale)
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	if locale != "en" {
		if v, ok := enMsgs[key]; ok {
			return v
		}
	}
	return key
}

func enPlural(locale string, n int, singularKey, pluralKey string) string {
	if locale != "en" {
		return tr(locale, singularKey)
	}
	if n == 1 {
		return tr("en", singularKey)
	}
	return tr("en", pluralKey)
}
