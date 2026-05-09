//go:build windows

// Package systray wraps getlantern/systray to provide a native Windows system
// tray icon and a minimal context menu for spotilite.
package systray

import (
	_ "embed"
	"log/slog"

	"github.com/getlantern/systray"

	"spotilite/internal/i18n"
)

//go:embed ../../build/windows/icon_tray.ico
var trayIconData []byte

type menuItems struct {
	show *systray.MenuItem
	quit *systray.MenuItem
}

type Manager struct {
	i18n   *i18n.Translator
	onShow func()
	onQuit func()
	items  menuItems
}

func NewManager(
	i18n *i18n.Translator,
	iconPath string,
	onShow, onQuit func(),
) *Manager {
	return &Manager{
		i18n:   i18n,
		onShow: onShow,
		onQuit: onQuit,
	}
}

func (m *Manager) Start() {
	go systray.Run(m.onReady, m.onExit)
}

func (m *Manager) Refresh() {
	if m.items.show == nil {
		return
	}
	m.items.show.SetTitle(m.i18n.T("tray.show"))
	m.items.quit.SetTitle(m.i18n.T("tray.quit"))
}

func (m *Manager) onReady() {
	systray.SetIcon(trayIconData)
	systray.SetTooltip(m.i18n.T("app.title"))

	m.items.show = systray.AddMenuItem(m.i18n.T("tray.show"), "Show window")
	systray.AddSeparator()
	m.items.quit = systray.AddMenuItem(m.i18n.T("tray.quit"), "Quit application")

	go m.handleClicks()
}

func (m *Manager) handleClicks() {
	for {
		select {
		case <-m.items.show.ClickedCh:
			if m.onShow != nil {
				m.onShow()
			}
		case <-m.items.quit.ClickedCh:
			if m.onQuit != nil {
				m.onQuit()
			}
			systray.Quit()
		}
	}
}

func (m *Manager) onExit() {
	slog.Info("system tray exited")
}
