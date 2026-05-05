//go:build windows

// Package systray wraps getlantern/systray to provide a native Windows system
// tray icon and a minimal context menu for spotilite.
package systray

import (
	"log/slog"
	"os"

	"github.com/getlantern/systray"

	"spotilite/internal/i18n"
)

// Manager handles the system tray lifecycle, icon and minimal context menu.
type Manager struct {
	i18n    *i18n.Translator
	iconPath string
	onShow  func()
	onQuit  func()
	items   menuItems
}

type menuItems struct {
	show *systray.MenuItem
	quit *systray.MenuItem
}

// NewManager creates a system tray manager. iconPath should point to a .ico
// file (e.g. "build/windows/icon.ico").
func NewManager(
	i18n *i18n.Translator,
	iconPath string,
	onShow, onQuit func(),
) *Manager {
	return &Manager{
		i18n:     i18n,
		iconPath: iconPath,
		onShow:   onShow,
		onQuit:   onQuit,
	}
}

// Start launches the system tray in its own goroutine. This method is
// non-blocking.
func (m *Manager) Start() {
	go systray.Run(m.onReady, m.onExit)
}

// Refresh updates menu item labels to reflect the current language.
func (m *Manager) Refresh() {
	if m.items.show == nil {
		return
	}
	m.items.show.SetTitle(m.i18n.T("tray.show"))
	m.items.quit.SetTitle(m.i18n.T("tray.quit"))
}

func (m *Manager) onReady() {
	iconBytes, err := os.ReadFile(m.iconPath)
	if err != nil {
		slog.Warn("failed to read tray icon", "path", m.iconPath, "error", err)
	} else {
		systray.SetIcon(iconBytes)
	}
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
