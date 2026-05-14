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

//go:embed icon_tray.ico
var trayIconData []byte

type menuItems struct {
	show         *systray.MenuItem
	background   *systray.MenuItem
	adblock      *systray.MenuItem
	proxyblock   *systray.MenuItem
	sectionblock *systray.MenuItem
	premiumspoof *systray.MenuItem
	experiments  *systray.MenuItem
	history      *systray.MenuItem
	quit         *systray.MenuItem
}

type Manager struct {
	i18n           *i18n.Translator
	onShow         func()
	onQuit         func()
	onToggle       func(module string, enabled bool)
	onToggleProxy   func(enabled bool)
	onToggleBg     func(enabled bool)
	getBgMode      func() bool
	items          menuItems
	state          *TrayState
}

type TrayState struct {
	Background   bool
	AdBlock      bool
	SectionBlock bool
	PremiumSpoof bool
	Experiments  bool
	History      bool
}

func NewManager(
	i18n *i18n.Translator,
	iconPath string,
	onShow, onQuit func(),
	onToggle func(module string, enabled bool),
	onToggleProxy func(enabled bool),
	onToggleBg func(enabled bool),
	getBgMode func() bool,
) *Manager {
	return &Manager{
		i18n:          i18n,
		onShow:        onShow,
		onQuit:        onQuit,
		onToggle:      onToggle,
		onToggleProxy: onToggleProxy,
		onToggleBg:    onToggleBg,
		getBgMode:     getBgMode,
		state: &TrayState{
			AdBlock:      true,
			SectionBlock: true,
			PremiumSpoof: true,
			Experiments:  true,
			History:      true,
		},
	}
}

func (m *Manager) Start() {
	go systray.Run(m.onReady, m.onExit)
}

func (m *Manager) UpdateState(state *TrayState) {
	m.state = state
	m.Refresh()
}

func (m *Manager) Refresh() {
	if m.items.adblock == nil {
		return
	}
	if m.items.background != nil {
		bgMode := false
		if m.getBgMode != nil {
			bgMode = m.getBgMode()
		}
		m.updateCheckbox(m.items.background, bgMode)
	}
	m.updateCheckbox(m.items.adblock, m.state.AdBlock)
	m.updateCheckbox(m.items.sectionblock, m.state.SectionBlock)
	m.updateCheckbox(m.items.premiumspoof, m.state.PremiumSpoof)
	m.updateCheckbox(m.items.experiments, m.state.Experiments)
	m.updateCheckbox(m.items.history, m.state.History)
}

func (m *Manager) updateCheckbox(item *systray.MenuItem, checked bool) {
	if checked {
		item.Check()
	} else {
		item.Uncheck()
	}
}

func (m *Manager) onReady() {
	systray.SetIcon(trayIconData)
	systray.SetTooltip("Spotilite")

	m.items.show = systray.AddMenuItem(m.i18n.T("tray.show"), m.i18n.T("tray.show"))

	systray.AddSeparator()

	bgMode := false
	if m.getBgMode != nil {
		bgMode = m.getBgMode()
	}
	m.items.background = systray.AddMenuItemCheckbox(m.i18n.T("tray.runInBackground"), m.i18n.T("tray.runInBackground"), bgMode)

	systray.AddSeparator()

	m.items.adblock = systray.AddMenuItemCheckbox(m.i18n.T("spotx.adblock"), m.i18n.T("spotx.adblock"), true)
	m.items.proxyblock = systray.AddMenuItemCheckbox(m.i18n.T("spotx.proxyblock"), m.i18n.T("spotx.proxyblock"), true)
	m.items.sectionblock = systray.AddMenuItemCheckbox(m.i18n.T("spotx.sections"), m.i18n.T("spotx.sections"), true)
	m.items.premiumspoof = systray.AddMenuItemCheckbox(m.i18n.T("spotx.premium"), m.i18n.T("spotx.premium"), true)
	m.items.experiments = systray.AddMenuItemCheckbox(m.i18n.T("spotx.experiments"), m.i18n.T("spotx.experiments"), true)
	m.items.history = systray.AddMenuItemCheckbox(m.i18n.T("spotx.history"), m.i18n.T("spotx.history"), true)

	systray.AddSeparator()

	m.items.quit = systray.AddMenuItem(m.i18n.T("tray.quit"), m.i18n.T("tray.quit"))

	go m.handleClicks()
}

func (m *Manager) handleClicks() {
	for {
		select {
		case <-m.items.show.ClickedCh:
			if m.onShow != nil {
				m.onShow()
			}
		case <-m.items.background.ClickedCh:
			m.state.Background = !m.state.Background
			m.updateCheckbox(m.items.background, m.state.Background)
			if m.onToggleBg != nil {
				m.onToggleBg(m.state.Background)
			}
		case <-m.items.adblock.ClickedCh:
			m.state.AdBlock = !m.state.AdBlock
			m.updateCheckbox(m.items.adblock, m.state.AdBlock)
			if m.onToggle != nil {
				m.onToggle("adblock", m.state.AdBlock)
			}
		case <-m.items.proxyblock.ClickedCh:
			m.state.AdBlock = !m.state.AdBlock
			m.updateCheckbox(m.items.proxyblock, m.state.AdBlock)
			if m.onToggleProxy != nil {
				m.onToggleProxy(m.state.AdBlock)
			}
		case <-m.items.sectionblock.ClickedCh:
			m.state.SectionBlock = !m.state.SectionBlock
			m.updateCheckbox(m.items.sectionblock, m.state.SectionBlock)
			if m.onToggle != nil {
				m.onToggle("sectionblock", m.state.SectionBlock)
			}
		case <-m.items.premiumspoof.ClickedCh:
			m.state.PremiumSpoof = !m.state.PremiumSpoof
			m.updateCheckbox(m.items.premiumspoof, m.state.PremiumSpoof)
			if m.onToggle != nil {
				m.onToggle("premium_spoof", m.state.PremiumSpoof)
			}
		case <-m.items.experiments.ClickedCh:
			m.state.Experiments = !m.state.Experiments
			m.updateCheckbox(m.items.experiments, m.state.Experiments)
			if m.onToggle != nil {
				m.onToggle("experiments", m.state.Experiments)
			}
		case <-m.items.history.ClickedCh:
			m.state.History = !m.state.History
			m.updateCheckbox(m.items.history, m.state.History)
			if m.onToggle != nil {
				m.onToggle("history", m.state.History)
			}
		case <-m.items.quit.ClickedCh:
			if m.onQuit != nil {
				m.onQuit()
			}
			return
		}
	}
}

func (m *Manager) onExit() {
	slog.Info("system tray exited")
}

func (m *Manager) Quit() {
	systray.Quit()
}
