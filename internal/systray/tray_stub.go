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
	_, _ func(),
) *Manager {
	return &Manager{}
}

// Start does nothing.
func (m *Manager) Start() {}

// Refresh does nothing.
func (m *Manager) Refresh() {}

// Quit does nothing.
func (m *Manager) Quit() {}
