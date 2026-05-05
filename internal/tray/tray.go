// Package tray builds and manages the system tray context menu for spotilite.
package tray

import (
	"context"
	"log/slog"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/i18n"
)

// Tray manages the system tray menu and its lifecycle.
type Tray struct {
	ctx                context.Context
	i18n               *i18n.Translator
	runInBackground    bool
	onToggleBackground func(bool)
}

// New creates a new Tray instance.
func New(i18n *i18n.Translator, runInBackground bool, onToggleBackground func(bool)) *Tray {
	return &Tray{
		i18n:               i18n,
		runInBackground:    runInBackground,
		onToggleBackground: onToggleBackground,
	}
}

// SetContext stores the Wails runtime context so tray callbacks can call
// runtime helpers (Show, Hide, Quit, etc.).
func (t *Tray) SetContext(ctx context.Context) {
	t.ctx = ctx
}

// Build constructs the application menu based on the current translations.
func (t *Tray) Build() *menu.Menu {
	appMenu := menu.NewMenu()

	// Show Window
	appMenu.Append(menu.Text(t.i18n.T("tray.show"), nil, func(_ *menu.CallbackData) {
		runtime.Show(t.ctx)
	}))

	// Hide Window
	appMenu.Append(menu.Text(t.i18n.T("tray.hide"), nil, func(_ *menu.CallbackData) {
		runtime.Hide(t.ctx)
	}))

	appMenu.Append(menu.Separator())

	// Run in Background toggle (using visual checkmark for reliability on Windows)
	bgLabel := t.formatToggleLabel("tray.runInBackground", t.runInBackground)
	appMenu.Append(menu.Text(bgLabel, nil, func(_ *menu.CallbackData) {
		t.runInBackground = !t.runInBackground
		if t.onToggleBackground != nil {
			t.onToggleBackground(t.runInBackground)
		}
		// Refresh menu to update the checkmark visual
		t.Refresh()
	}))

	appMenu.Append(menu.Separator())

	// Language submenu
	langMenu := appMenu.AddSubmenu(t.i18n.T("tray.language"))
	langMenu.Append(menu.Text(t.i18n.T("tray.lang.en"), nil, func(_ *menu.CallbackData) {
		t.changeLanguage(i18n.LangEnglish)
	}))
	langMenu.Append(menu.Text(t.i18n.T("tray.lang.es"), nil, func(_ *menu.CallbackData) {
		t.changeLanguage(i18n.LangSpanish)
	}))

	appMenu.Append(menu.Separator())

	// Quit
	appMenu.Append(menu.Text(t.i18n.T("tray.quit"), nil, func(_ *menu.CallbackData) {
		runtime.Quit(t.ctx)
	}))

	return appMenu
}

// formatToggleLabel prefixes the label with a checkmark when active,
// or an indentation when inactive, so the state is visually obvious.
func (t *Tray) formatToggleLabel(key string, active bool) string {
	if active {
		return "[x] " + t.i18n.T(key)
	}
	return "[ ] " + t.i18n.T(key)
}

// Apply registers the built menu with the Wails runtime.
func (t *Tray) Apply() {
	if t.ctx == nil {
		slog.Warn("tray context not set, skipping menu apply")
		return
	}
	runtime.MenuSetApplicationMenu(t.ctx, t.Build())
	slog.Info("application menu applied", "lang", t.i18n.Language(), "background", t.runInBackground)
}

// Refresh rebuilds and re-applies the menu. Useful after language changes.
func (t *Tray) Refresh() {
	t.Apply()
}

func (t *Tray) changeLanguage(lang string) {
	t.i18n.SetLanguage(lang)
	slog.Info("language changed", "lang", lang)
	t.Refresh()
}
