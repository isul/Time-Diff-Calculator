package backend

import "testing"

func TestResolveLocaleOverride(t *testing.T) {
	if ResolveLocale("ko") != "ko" {
		t.Fail()
	}
	if ResolveLocale("en") != "en" {
		t.Fail()
	}
}

func TestResolveLocaleFromEnv(t *testing.T) {
	t.Setenv("LANG", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "ko_KR.UTF-8")
	if ResolveLocale("auto") != "ko" {
		t.Fatalf("got %s", ResolveLocale("auto"))
	}
}

func TestResolveLocaleFallbackEn(t *testing.T) {
	t.Setenv("LANG", "fr_FR.UTF-8")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	// fr -> en fallback
	if ResolveLocale("auto") != "en" {
		t.Fatalf("got %s", ResolveLocale("auto"))
	}
}
