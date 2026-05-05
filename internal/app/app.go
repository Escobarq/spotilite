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
	apptray "spotilite/internal/systray"
)

// App is the Wails application struct that exposes bound methods to the
// frontend and coordinates backend services.
type App struct {
	ctx             context.Context
	i18n            *i18n.Translator
	tray            *apptray.Manager
	injector        *spotify.Injector
	runInBackground bool
	windowVisible   bool
	maximized       bool
}

// NewApp creates a new App application struct with its dependencies injected.
func NewApp(i18n *i18n.Translator, tray *apptray.Manager, runInBackground bool) *App {
	return &App{
		i18n:            i18n,
		tray:            tray,
		injector:        spotify.NewInjector(),
		runInBackground: runInBackground,
		windowVisible:   true,
	}
}

// Startup is called when the app starts. It wires runtime context, starts the
// Spotify injector, applies the system tray, and registers global shortcuts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.tray.Start()

	if err := shortcut.Register(a.ToggleWindowVisibility); err != nil {
		slog.Warn("failed to register global shortcut", "error", err)
	} else {
		slog.Info("global shortcut registered", "shortcut", "Ctrl+Shift+S")
	}

	a.injector.Start(ctx)
}

// Shutdown is called when the app is about to quit.
func (a *App) Shutdown(_ context.Context) {
	slog.Info("shutting down, stopping global shortcuts")
	shortcut.Unregister()
}

// SetBackgroundMode updates whether the app should hide to tray on close.
func (a *App) SetBackgroundMode(enabled bool) {
	a.runInBackground = enabled
	a.tray.SetBackgroundState(enabled)
	slog.Info("background mode changed", "enabled", enabled)
}

// IsBackgroundMode reports whether the app is configured to run in background.
func (a *App) IsBackgroundMode() bool {
	return a.runInBackground
}

// SetLanguage changes the app language and refreshes the tray menu labels.
func (a *App) SetLanguage(lang string) {
	a.i18n.SetLanguage(lang)
	a.tray.Refresh()
	slog.Info("language changed", "lang", lang)
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

// ---------------------------------------------------------------------------
// Window controls (bound to frontend)
// ---------------------------------------------------------------------------

// Minimize minimises the main window.
func (a *App) Minimize() {
	runtime.WindowMinimise(a.ctx)
}

// Maximize maximises the main window.
func (a *App) Maximize() {
	runtime.WindowMaximise(a.ctx)
	a.maximized = true
}

// UnMaximize restores the main window from maximised state.
func (a *App) UnMaximize() {
	runtime.WindowUnmaximise(a.ctx)
	a.maximized = false
}

// ToggleMaximize toggles between maximised and normal window state.
func (a *App) ToggleMaximize() {
	if a.maximized {
		a.UnMaximize()
	} else {
		a.Maximize()
	}
}

// Close closes the main window. Respects background mode.
func (a *App) Close() {
	if a.runInBackground {
		runtime.Hide(a.ctx)
		a.windowVisible = false
		beeep.Notify(a.i18n.T("app.title"), a.i18n.T("notif.minimizedToTray"), "")
	} else {
		runtime.Quit(a.ctx)
	}
}

// IsMaximized reports whether the window is currently maximised.
func (a *App) IsMaximized() bool {
	return runtime.WindowIsMaximised(a.ctx)
}

// ToggleWindowVisibility shows or hides the main window.
func (a *App) ToggleWindowVisibility() {
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
