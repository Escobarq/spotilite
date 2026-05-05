//go:build !windows

// Package systray is a stub on non-Windows platforms.
package systray

import "spotilite/internal/i18n"

// Manager is a no-op stub.
type Manager struct{}

// NewManager creates a no-op manager.
func NewManager(
	_ *i18n.Translator,
	_ string,
	_ bool,
	_, _, _ func(),
	_ func(bool),
	_ func(string),
) *Manager {
	return &Manager{}
}

// Start does nothing.
func (m *Manager) Start() {}

// Refresh does nothing.
func (m *Manager) Refresh() {}

// SetBackgroundState does nothing.
func (m *Manager) SetBackgroundState(_ bool) {}
