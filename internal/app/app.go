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
	"spotilite/internal/spotify/spicetify"
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

	// Build the full injection payload that runs *after* the bundled spicetify
	// shims. spicetify-js lives in SetExtraCSS/SetExtraJS separately; extra
	// here carries the runtime overlay (DevTools shortcut, etc.).
	var extra strings.Builder

	// Build Spicetify wrapper bundle (matches spicetify-cli apply.go htmlMod).
	spiceJS, spiceCSS := a.buildSpicetifyPayload()
	extra.WriteString(spiceJS)
	if spiceCSS != "" {
		a.injector.SetExtraCSS(spiceCSS)
	}

	// DevTools hotkey (Ctrl+Shift+D)
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

	slog.Info("[Spotilite] ExtraJS prepared", "len", extra.Len())
	a.injector.SetExtraJS(extra.String())
	a.injector.Start(ctx)
}

// buildSpicetifyPayload assembles the runtime Spicetify bundle for injection
// into the Spotify webview. It returns (js, css). Order matches the
// spicetify-cli apply package: shims first, then Spicetify.Config, then
// theme.js, helpers, user extensions, custom apps.
//
// If no config is loaded we still emit the bare shims so extensions that
// install only require Spicetify events / URI etc. continue to work.
func (a *App) buildSpicetifyPayload() (js string, css string) {
	ctx := spicetify.BuildContext{
		Config:   a.cfg,
		LocalAPI: spotify.APIBaseURL,
	}

	// Theme: best-effort; missing theme leaves css "" (no-op in injector).
	if a.cfg != nil && a.cfg.CurrentTheme != "" && a.themeManager != nil {
		if t, err := a.themeManager.Load(a.cfg.CurrentTheme); err == nil {
			ctx.Theme = t
		} else {
			slog.Warn("theme not loaded; skipping theme injection", "theme", a.cfg.CurrentTheme, "error", err)
		}
	}

	// Extensions: load files declared in AdditionalOptions; missing ones are
	// logged and skipped (matches spicetify-cli error visibility).
	if a.cfg != nil && a.extLoader != nil {
		names := a.cfg.EnabledExtensions()
		if len(names) > 0 {
			loaded, skipped, err := a.extLoader.Load(names)
			if err != nil {
				slog.Warn("extension loader error", "error", err)
			}
			if len(skipped) > 0 {
				slog.Warn("skipped extensions not found on disk", "names", skipped)
			}
			ctx.Extensions = loaded
		}
	}

	// Custom apps: same shape.
	if a.cfg != nil && a.appManager != nil {
		apps := make([]*customapps.App, 0)
		for _, name := range a.cfg.CustomApps {
			if strings.HasSuffix(name, "-") {
				continue
			}
			app, err := a.appManager.Load(strings.TrimSuffix(name, "-"))
			if err != nil {
				slog.Warn("custom app load failed", "name", name, "error", err)
				continue
			}
			apps = append(apps, app)
		}
		ctx.CustomApps = apps
	}

	js = spicetify.Bundle(ctx)
	css = ctx.ThemeCSS()
	return
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

// --- spicetify bridge ---------------------------------------------------------
//
// Implements api.SpicetifyHandler. These methods extend the API surface that
// webview-side extensions consult when they need to know the active
// configuration or to ask Go to switch themes / toggle extensions / reload
// the injection payload.

func (a *App) GetSpicetifyConfig() api.SpicetifyConfigDTO {
	dto := api.SpicetifyConfigDTO{
		Version:              "spotilite-1.0",
		LocalAPI:             "http://localhost:" + api.DefaultPort,
		InjectCSS:            true,
		InjectThemeJS:        true,
		ReplaceColors:        true,
		SidebarConfig:        true,
		HomeConfig:           true,
		ExperimentalFeatures: true,
	}
	if a.cfg != nil {
		dto.CurrentTheme = a.cfg.CurrentTheme
		dto.ColorScheme = a.cfg.ColorScheme
		dto.Extensions = append([]string(nil), a.cfg.Extensions...)
		dto.CustomApps = append([]string(nil), a.cfg.CustomApps...)
		dto.InjectCSS = a.cfg.InjectCSS
		dto.InjectThemeJS = a.cfg.InjectThemeJS
		dto.ReplaceColors = a.cfg.ReplaceColors
		dto.SidebarConfig = a.cfg.SidebarConfig
		dto.HomeConfig = a.cfg.HomeConfig
		dto.ExperimentalFeatures = a.cfg.ExpFeatures
	}
	return dto
}

func (a *App) GetSpicetifyExtensions() []api.ExtensionDTO {
	if a.cfg == nil || a.extLoader == nil {
		return nil
	}
	names := a.cfg.Extensions
	enabled := map[string]bool{}
	for _, n := range a.cfg.EnabledExtensions() {
		enabled[n] = true
	}
	out := make([]api.ExtensionDTO, 0, len(names))
	for _, raw := range names {
		clean := strings.TrimSuffix(raw, "-")
		isMJS := strings.HasSuffix(clean, ".mjs")
		display := clean
		if isMJS {
			display = strings.TrimSuffix(clean, ".mjs")
		} else if strings.HasSuffix(clean, ".js") {
			display = strings.TrimSuffix(clean, ".js")
		}
		out = append(out, api.ExtensionDTO{
			Name:    display,
			File:    clean,
			Enabled: enabled[raw],
			IsMJS:   isMJS,
		})
	}
	return out
}

func (a *App) SetSpicetifyExtension(name string, enabled bool) error {
	if a.cfg == nil {
		return fmt.Errorf("app: no config loaded")
	}
	name = strings.TrimSuffix(name, ".js")
	name = strings.TrimSuffix(name, ".mjs")
	a.cfg.SetExtension(name, enabled)
	if err := a.cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a *App) GetSpicetifyThemes() []string {
	if a.themeManager == nil {
		return nil
	}
	return a.themeManager.ListNames()
}

func (a *App) SetSpicetifyTheme(name, colorScheme string) error {
	if a.cfg == nil {
		return fmt.Errorf("app: no config loaded")
	}
	if name == "" {
		return fmt.Errorf("theme name empty")
	}
	if a.themeManager != nil {
		if _, err := a.themeManager.Load(name); err != nil {
			return fmt.Errorf("theme not found: %w", err)
		}
	}
	a.cfg.SetTheme(name, colorScheme)
	if err := a.cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (a *App) GetSpicetifyCustomApps() []string {
	out := []string{}
	if a.cfg != nil {
		out = append(out, a.cfg.CustomApps...)
	}
	return out
}

func (a *App) SetSpicetifyCustomApp(name string, enabled bool) error {
	if a.cfg == nil {
		return fmt.Errorf("app: no config loaded")
	}
	a.cfg.SetCustomApp(name, enabled)
	if err := a.cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

// ReloadInjection rebuilds the spicetify bundle from the current on-disk
// config and triggers an immediate re-inject. Called by the
// /api/spicetify/reload endpoint and any UI flow that mutates the config.
func (a *App) ReloadInjection(ctx context.Context) error {
	if a.ctx == nil {
		return fmt.Errorf("app: not running")
	}
	slog.Info("reload injection requested from API")
	var extra strings.Builder
	spiceJS, spiceCSS := a.buildSpicetifyPayload()
	extra.WriteString(spiceJS)
	a.injector.SetExtraJS(extra.String())
	if spiceCSS != "" {
		a.injector.SetExtraCSS(spiceCSS)
	}
	return nil
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
