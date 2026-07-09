package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/api"
	"spotilite/internal/config"
	"spotilite/internal/customapps"
	"spotilite/internal/extensions"
	"spotilite/internal/i18n"
	"spotilite/internal/proxy"
	"spotilite/internal/shortcut"
	"spotilite/internal/spotify"
	"spotilite/internal/spotify/modules"
	"spotilite/internal/themes"
	apptray "spotilite/internal/systray"
)

type App struct {
	ctx                   context.Context
	i18n                  *i18n.Translator
	tray                  *apptray.Manager
	injector              *spotify.Injector
	api                   *api.Server
	cfg                   *config.Config
	themeManager          *themes.Manager
	extLoader             *extensions.Loader
	appManager            *customapps.Manager
	runInBackground       bool
	windowVisible         bool
	maximized             bool
	hasNotifiedBackground bool
	iconPath              string
	adProxy               *proxy.AdBlockProxy
}

func NewApp(i18n *i18n.Translator, tray *apptray.Manager, apiServer *api.Server, cfg *config.Config, themeManager *themes.Manager, extLoader *extensions.Loader, appManager *customapps.Manager, runInBackground bool, iconPath string) *App {
	injector := spotify.NewInjector()
	adProxy := proxy.NewAdBlockProxy("8766")
	return &App{
		i18n:            i18n,
		tray:            tray,
		injector:        injector,
		api:             apiServer,
		cfg:             cfg,
		themeManager:    themeManager,
		extLoader:       extLoader,
		appManager:      appManager,
		runInBackground: runInBackground,
		windowVisible:   true,
		iconPath:        iconPath,
		adProxy:         adProxy,
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

	if err := a.adProxy.Start(ctx); err != nil {
		slog.Warn("failed to start ad-block proxy", "error", err)
	}

	// Prepare the full injection payload: Spicetify wrapper + theme +
	// extensions + custom apps. This payload is prepended before the builtin
	// modules and the title bar on every injection tick.
	var extra strings.Builder

	// Note: spicetify.Bundle() is injected separately in the Injector

	extra.WriteString(`(function(){
var pressed={};
window.addEventListener('keydown',function(e){
pressed[e.key]=true;
if(pressed['d']&&(pressed['Control']||pressed['Meta'])){
e.preventDefault();
try{
if(window.chrome&&window.chrome.webview&&window.chrome.webview.openDevToolsWindow){
window.chrome.webview.openDevToolsWindow();
}else{
console.warn('[Spotilite] DevTools not available via chrome.webview');
}
}catch(err){console.error('[Spotilite] DevTools error:',err);}
}
},true);
window.addEventListener('keyup',function(e){delete pressed[e.key];});
})();`)

	// Extensions
	if a.cfg != nil {
		extNames := a.cfg.EnabledExtensions()
		if len(extNames) > 0 {
			script, skipped, err := a.extLoader.LoadAndBundle(extNames)
			if err != nil {
				slog.Warn("extension loader error", "error", err)
			}
			if len(skipped) > 0 {
				slog.Warn("skipped extensions not found on disk", "names", skipped)
			}
			extra.WriteString(script)
		}
	}

	slog.Info("[Spotilite] ExtraJS prepared", "len", extra.Len())
	a.injector.SetExtraJS(extra.String())
	a.injector.Start(ctx)
}

func (a *App) Shutdown(_ context.Context) {
	slog.Info("shutting down, stopping global shortcuts, api, proxy and systray")
	shortcut.Unregister()
	a.tray.Quit()
	a.api.Stop(context.Background())
	a.adProxy.Stop()
	os.Exit(0)
}

func (a *App) SetBackgroundMode(enabled bool) {
	a.runInBackground = enabled
	slog.Info("background mode changed", "enabled", enabled)

	if !enabled && !a.windowVisible {
		slog.Info("background mode disabled while window hidden, quitting application")
		shortcut.Unregister()
		a.tray.Quit()
		os.Exit(0)
	}
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

func (a *App) SetProxyEnabled(enabled bool) {
	a.adProxy.SetEnabled(enabled)
	slog.Info("proxy ad-blocking toggled", "enabled", enabled)
}

func (a *App) SetLyricsTheme(theme string) {
	slog.Info("lyrics theme set", "theme", theme)
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
	if a.runInBackground {
		slog.Info("window close requested, hiding to system tray")
		runtime.Hide(a.ctx)
		a.windowVisible = false

		if !a.hasNotifiedBackground {
			if err := beeep.Notify(
				a.i18n.T("app.title"),
				a.i18n.T("notif.minimizedToTray"),
				a.iconPath,
			); err != nil {
				slog.Warn("failed to show native notification", "error", err)
			}
			a.hasNotifiedBackground = true
		}

		// Liberar memoria de Go cuando pasa a segundo plano
		go debug.FreeOSMemory()

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

func (a *App) ForceQuit() {
	slog.Info("force quit requested")
	shortcut.Unregister()
	a.tray.Quit()
	os.Exit(0)
}

func (a *App) Close() {
	if a.runInBackground {
		slog.Info("hiding to system tray (background mode enabled)")
		runtime.Hide(a.ctx)
		a.windowVisible = false
		beeep.Notify(a.i18n.T("app.title"), a.i18n.T("notif.minimizedToTray"), a.iconPath)
		a.hasNotifiedBackground = true
	} else {
		slog.Info("closing application (background mode disabled)")
		shortcut.Unregister()
		a.tray.Quit()
		os.Exit(0)
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
		go debug.FreeOSMemory()
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
	if a.ctx != nil {
		runtime.WindowExecJS(a.ctx, `try{if(window.chrome&&window.chrome.webview&&window.chrome.webview.openDevToolsWindow){window.chrome.webview.openDevToolsWindow()}else{console.warn('[Spotilite] DevTools not available via chrome.webview')}}catch(e){console.error('[Spotilite] DevTools error:',e)}`)
		slog.Info("devtools opened")
	}
}

func jsStringGo(s string) string {
	var b strings.Builder
	b.WriteRune('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\u0000':
			b.WriteString(`\0`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteRune('"')
	return b.String()
}
