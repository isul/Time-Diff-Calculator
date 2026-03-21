package main

import (
	"context"
	"timediff/backend"
)

// App binds to the Wails frontend.
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a App) domReady(ctx context.Context) {}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (a *App) shutdown(ctx context.Context) {}

// Calculate runs parse + duration + formatting.
func (a *App) Calculate(input string, settings backend.Settings) backend.CalculateResponse {
	return backend.Calculate(input, settings)
}

// ValidateInput provides real-time validation messages.
func (a *App) ValidateInput(input string, settings backend.Settings) backend.ValidateResponse {
	return backend.ValidateInput(input, settings)
}

// LoadSettings reads persisted settings from disk.
func (a *App) LoadSettings() backend.Settings {
	return backend.LoadSettings()
}

// SaveSettings persists settings to disk.
func (a *App) SaveSettings(settings backend.Settings) error {
	return backend.SaveSettings(settings)
}

// SystemLocale returns resolved locale ("ko" or "en") from OS environment.
func (a *App) SystemLocale() string {
	return backend.ResolveLocale("auto")
}
