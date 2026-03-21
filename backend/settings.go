package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func boolPtr(b bool) *bool {
	return &b
}

// DefaultSettings matches PRD: all predefined formats on; custom available when checked + non-empty template.
func DefaultSettings() Settings {
	return Settings{
		Formats: map[string]bool{
			"seconds":      true,
			"milliseconds": true,
			"mmss":         true,
			"hhmmss":       true,
			"ddhhmmss":     true,
			"full":         true,
			"custom":       true,
		},
		CustomFormat: "",
		Language:     "auto",
		ZeroPadding:  boolPtr(true),
	}
}

// EffectiveZeroPadding reports whether numeric time parts use leading zeros (default true when unset).
func EffectiveZeroPadding(s Settings) bool {
	if s.ZeroPadding == nil {
		return true
	}
	return *s.ZeroPadding
}

// SettingsDir returns the directory for app config.
func SettingsDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "timediff")
	return dir, nil
}

// SettingsFilePath returns absolute path to settings.json.
func SettingsFilePath() (string, error) {
	dir, err := SettingsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

// MergeSettings fills missing format keys from defaults.
func MergeSettings(s Settings) Settings {
	def := DefaultSettings()
	if s.Formats == nil {
		s.Formats = make(map[string]bool)
	}
	for k, v := range def.Formats {
		if _, ok := s.Formats[k]; !ok {
			s.Formats[k] = v
		}
	}
	if strings.TrimSpace(s.Language) == "" {
		s.Language = "auto"
	}
	if s.ZeroPadding == nil {
		s.ZeroPadding = boolPtr(true)
	}
	return s
}

// LoadSettings reads JSON from disk; on error returns defaults.
func LoadSettings() Settings {
	path, err := SettingsFilePath()
	if err != nil {
		return DefaultSettings()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultSettings()
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return DefaultSettings()
	}
	return MergeSettings(s)
}

// SaveSettings writes JSON (creates directory if needed).
func SaveSettings(s Settings) error {
	path, err := SettingsFilePath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	s = MergeSettings(s)
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
