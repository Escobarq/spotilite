//go:build windows

// Package systray wraps getlantern/systray to provide a native Windows system
// tray icon and context menu for spotilite.
package systray

import (
	"log/slog"
	"os"

	"github.com/getlantern/systray"

	"spotilite/internal/i18n"
)

// Manager handles the system tray lifecycle, icon and context menu.
type Manager struct {
	i18n               *i18n.Translator
	runInBackground    bool
	iconPath           string
	onShow             func()
	onHide             func()
	onQuit             func()
	onToggleBackground func(enabled bool)
	onSetLanguage      func(lang string)
	items              menuItems
}

type menuItems struct {
	show             *systray.MenuItem
	hide             *systray.MenuItem
	bg               *systray.MenuItem
	langEN           *systray.MenuItem
	langES           *systray.MenuItem
	quit             *systray.MenuItem
}

// NewManager creates a system tray manager. iconPath should point to a .ico
// file (e.g. "build/windows/icon.ico").
func NewManager(
	i18n *i18n.Translator,
	iconPath string,
	runInBackground bool,
	onShow, onHide, onQuit func(),
	onToggleBackground func(enabled bool),
	onSetLanguage func(lang string),
) *Manager {
	return &Manager{
		i18n:               i18n,
		iconPath:           iconPath,
		runInBackground:    runInBackground,
		onShow:             onShow,
		onHide:             onHide,
		onQuit:             onQuit,
		onToggleBackground: onToggleBackground,
		onSetLanguage:      onSetLanguage,
	}
}

// Start launches the system tray in its own goroutine. This method is
// non-blocking.
func (m *Manager) Start() {
	go systray.Run(m.onReady, m.onExit)
}

// Refresh updates all menu item labels to reflect the current language and
// background-mode state. Must be called after language changes.
func (m *Manager) Refresh() {
	if m.items.show == nil {
		return
	}
	m.items.show.SetTitle(m.i18n.T("tray.show"))
	m.items.hide.SetTitle(m.i18n.T("tray.hide"))
	m.items.bg.SetTitle(m.formatToggleLabel())
	m.items.langEN.SetTitle(m.i18n.T("tray.lang.en"))
	m.items.langES.SetTitle(m.i18n.T("tray.lang.es"))
	m.items.quit.SetTitle(m.i18n.T("tray.quit"))
}

// SetBackgroundState updates the visual state of the background-mode item.
func (m *Manager) SetBackgroundState(enabled bool) {
	m.runInBackground = enabled
	if m.items.bg != nil {
		m.items.bg.SetTitle(m.formatToggleLabel())
	}
}

func (m *Manager) onReady() {
	iconBytes, err := os.ReadFile(m.iconPath)
	if err != nil {
		slog.Warn("failed to read tray icon", "path", m.iconPath, "error", err)
	} else {
		systray.SetIcon(iconBytes)
	}
	systray.SetTooltip("Spotilite")

	m.items.show = systray.AddMenuItem(m.i18n.T("tray.show"), "Show window")
	m.items.hide = systray.AddMenuItem(m.i18n.T("tray.hide"), "Hide window")
	systray.AddSeparator()
	m.items.bg = systray.AddMenuItem(m.formatToggleLabel(), "Toggle background mode")
	systray.AddSeparator()

	langMenu := systray.AddMenuItem(m.i18n.T("tray.language"), "Change language")
	m.items.langEN = langMenu.AddSubMenuItem(m.i18n.T("tray.lang.en"), "Switch to English")
	m.items.langES = langMenu.AddSubMenuItem(m.i18n.T("tray.lang.es"), "Switch to Spanish")

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
		case <-m.items.hide.ClickedCh:
			if m.onHide != nil {
				m.onHide()
			}
		case <-m.items.bg.ClickedCh:
			m.runInBackground = !m.runInBackground
			m.items.bg.SetTitle(m.formatToggleLabel())
			if m.onToggleBackground != nil {
				m.onToggleBackground(m.runInBackground)
			}
		case <-m.items.langEN.ClickedCh:
			if m.onSetLanguage != nil {
				m.onSetLanguage(i18n.LangEnglish)
			}
		case <-m.items.langES.ClickedCh:
			if m.onSetLanguage != nil {
				m.onSetLanguage(i18n.LangSpanish)
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

func (m *Manager) formatToggleLabel() string {
	if m.runInBackground {
		return "[x] " + m.i18n.T("tray.runInBackground")
	}
	return "[ ] " + m.i18n.T("tray.runInBackground")
}
