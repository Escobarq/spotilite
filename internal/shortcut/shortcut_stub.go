//go:build !windows

// Package shortcut provides global keyboard shortcuts without CGO.
// Stub implementation for non-Windows platforms.
package shortcut

import "log/slog"

// Register is a no-op on non-Windows platforms.
func Register(callback func()) error {
	slog.Warn("global shortcuts are only supported on Windows")
	return nil
}

// Unregister is a no-op on non-Windows platforms.
func Unregister() {}
