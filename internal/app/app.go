// Package app defines the Wails application structure and lifecycle hooks.
package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/i18n"
	"spotilite/internal/shortcut"
	"spotilite/internal/spotify"
	"spotilite/internal/tray"
)

// App is the Wails application struct that exposes bound methods to the
// frontend and coordinates backend services.
type App struct {
	ctx             context.Context
	i18n            *i18n.Translator
	tray            *tray.Tray
	injector        *spotify.Injector
	runInBackground bool
	windowVisible   bool
}

// NewApp creates a new App application struct with its dependencies injected.
func NewApp(i18n *i18n.Translator, tray *tray.Tray, runInBackground bool) *App {
	return &App{
		i18n:            i18n,
		tray:            tray,
		injector:        spotify.NewInjector(),
		runInBackground: runInBackground,
		windowVisible:   true,
	}
}

// Startup is called when the app starts. It wires runtime context, starts the
// Spotify injector, applies the system tray menu, and registers global shortcuts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.tray.SetContext(ctx)
	a.tray.Apply()

	if err := shortcut.Register(a.toggleWindowVisibility); err != nil {
		slog.Warn("failed to register global shortcut", "error", err)
	} else {
		slog.Info("global shortcut registered", "shortcut", "Ctrl+Shift+S")
	}

	a.injector.Start(ctx)
}

// Shutdown is called when the app is about to quit. It cleans up resources
// such as the global keyboard hook.
func (a *App) Shutdown(_ context.Context) {
	slog.Info("shutting down, stopping global shortcuts")
	shortcut.Unregister()
}

// SetBackgroundMode updates whether the app should hide to tray on close.
func (a *App) SetBackgroundMode(enabled bool) {
	a.runInBackground = enabled
	slog.Info("background mode changed", "enabled", enabled)
}

// IsBackgroundMode reports whether the app is configured to run in background.
func (a *App) IsBackgroundMode() bool {
	return a.runInBackground
}

// OnBeforeClose is invoked when the user attempts to close the main window.
// If background mode is enabled, the window is hidden instead of destroyed.
func (a *App) OnBeforeClose(_ context.Context) bool {
	if a.runInBackground {
		slog.Info("window close requested, hiding to system tray")
		runtime.Hide(a.ctx)
		a.windowVisible = false

		if err := beeep.Notify(
			a.i18n.T("app.title"),
			a.i18n.T("notif.minimizedToTray"),
			"",
		); err != nil {
			slog.Warn("failed to show native notification", "error", err)
		}
		return true // prevent actual window destruction
	}

	slog.Info("window close requested, quitting application")
	return false // allow Wails to close the app
}

// toggleWindowVisibility shows or hides the main window.
func (a *App) toggleWindowVisibility() {
	if a.ctx == nil {
		return
	}
	if a.windowVisible {
		runtime.Hide(a.ctx)
		a.windowVisible = false
		slog.Info("window hidden via global shortcut")
	} else {
		runtime.Show(a.ctx)
		a.windowVisible = true
		slog.Info("window shown via global shortcut")
	}
}

// Greet returns a greeting for the given name. Exposed to the frontend.
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
