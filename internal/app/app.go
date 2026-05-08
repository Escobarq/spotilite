package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/api"
	"spotilite/internal/i18n"
	"spotilite/internal/shortcut"
	"spotilite/internal/spotify"
	"spotilite/internal/spotify/modules"
	apptray "spotilite/internal/systray"
)

type App struct {
	ctx context.Context
	i18n *i18n.Translator
	tray *apptray.Manager
	injector *spotify.Injector
	api *api.Server
	runInBackground bool
	windowVisible bool
	maximized bool
	hasNotifiedBackground bool
	iconPath string
}

func NewApp(i18n *i18n.Translator, tray *apptray.Manager, apiServer *api.Server, runInBackground bool, iconPath string) *App {
	injector := spotify.NewInjector()
	return &App{
		i18n: i18n,
		tray: tray,
		injector: injector,
		api: apiServer,
		runInBackground: runInBackground,
		windowVisible: true,
		iconPath: iconPath,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.api.Start()
	a.tray.Start()

	if err := shortcut.Register(a.ToggleWindowVisibility); err != nil {
		slog.Warn("failed to register global shortcut", "error", err)
	} else {
		slog.Info("global shortcut registered", "shortcut", "Ctrl+Shift+S")
	}

	a.injector.Start(ctx)
}

func (a *App) Shutdown(_ context.Context) {
	slog.Info("shutting down, stopping global shortcuts and api")
	shortcut.Unregister()
	if err := a.api.Stop(context.Background()); err != nil {
		slog.Warn("failed to stop api server", "error", err)
	}
}

func (a *App) SetBackgroundMode(enabled bool) {
	a.runInBackground = enabled
	slog.Info("background mode changed", "enabled", enabled)
}

func (a *App) IsBackgroundMode() bool {
	return a.runInBackground
}

func (a *App) SetLanguage(lang string) {
	a.i18n.SetLanguage(lang)
	a.tray.Refresh()
	slog.Info("language changed", "lang", lang)
}

func (a *App) GetSettings() api.Settings {
	return api.Settings{
		Language:       a.i18n.Language(),
		BackgroundMode: a.runInBackground,
	}
}

func (a *App) SetModuleEnabled(name string, enabled bool) {
	m := a.injector.GetModule(name)
	if m == nil {
		slog.Warn("module not found", "name", name)
		return
	}
	m.SetEnabled(enabled)
	slog.Info("module toggled", "name", name, "enabled", enabled)
}

func (a *App) SetLyricsTheme(css string) {
	localStorage := `localStorage.setItem('spotilite.custom_css', ` + "`" + css + "`" + `);`
	runtime.WindowExecJS(a.ctx, localStorage)
	a.injector.UpdateCustomCSS(a.ctx)
	slog.Info("custom css updated")
}

func (a *App) GetSpotXSettings() api.SpotXSettings {
	settings := api.SpotXSettings{
		AdBlock:      true,
		SectionBlock: true,
		PremiumSpoof: true,
		Experiments:  true,
		LyricsTheme:  modules.DefaultTheme,
		TrackHistory: true,
	}
	if m := a.injector.GetModule("adblock"); m != nil {
		settings.AdBlock = m.Enabled()
	}
	if m := a.injector.GetModule("sectionblock"); m != nil {
		settings.SectionBlock = m.Enabled()
	}
	if m := a.injector.GetModule("premium_spoof"); m != nil {
		settings.PremiumSpoof = m.Enabled()
	}
	if m := a.injector.GetModule("experiments"); m != nil {
		settings.Experiments = m.Enabled()
	}
	settings.LyricsTheme = "custom"
	if m := a.injector.GetModule("history"); m != nil {
		settings.TrackHistory = m.Enabled()
	}
	return settings
}

func (a *App) OnBeforeClose(_ context.Context) bool {
	if a.runInBackground && !a.hasNotifiedBackground {
		slog.Info("window close requested, hiding to system tray")
		runtime.Hide(a.ctx)
		a.windowVisible = false

		if err := beeep.Notify(
			a.i18n.T("app.title"),
			a.i18n.T("notif.minimizedToTray"),
			a.iconPath,
		); err != nil {
			slog.Warn("failed to show native notification", "error", err)
		}
		a.hasNotifiedBackground = true
		return true
	}

	slog.Info("window close requested, quitting application")
	return false
}

func (a *App) Minimize() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) Maximize() {
	runtime.WindowMaximise(a.ctx)
	a.maximized = true
}

func (a *App) UnMaximize() {
	runtime.WindowUnmaximise(a.ctx)
	a.maximized = false
}

func (a *App) Close() {
	if a.runInBackground && !a.hasNotifiedBackground {
		runtime.Hide(a.ctx)
		a.windowVisible = false
		beeep.Notify(a.i18n.T("app.title"), a.i18n.T("notif.minimizedToTray"), a.iconPath)
		a.hasNotifiedBackground = true
	} else {
		runtime.Quit(a.ctx)
	}
}

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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) OpenDevTools() {
	runtime.BrowserOpenURL(a.ctx, "http://localhost:34115")
}
