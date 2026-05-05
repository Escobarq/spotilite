package main

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"spotilite/internal/app"
	"spotilite/internal/i18n"
	apptray "spotilite/internal/systray"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	translator := i18n.New()

	// Default: run in background when closing the window.
	runInBackground := true

	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	iconPath := filepath.Join(filepath.Dir(exe), "build", "windows", "icon.ico")
	// Fallback for development mode.
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		iconPath = filepath.Join(".", "build", "windows", "icon.ico")
	}

	var application *app.App

	trayManager := apptray.NewManager(
		translator,
		iconPath,
		runInBackground,
		func() {
			if application != nil {
				application.ToggleWindowVisibility()
			}
		},
		func() {
			if application != nil {
				application.ToggleWindowVisibility()
			}
		},
		func() {
			// Allow quit even if background mode is on.
			os.Exit(0)
		},
		func(enabled bool) {
			if application != nil {
				application.SetBackgroundMode(enabled)
			}
		},
		func(lang string) {
			if application != nil {
				application.SetLanguage(lang)
			}
		},
	)

	application = app.NewApp(translator, trayManager, runInBackground)

	err = wails.Run(&options.App{
		Title:     translator.T("app.title"),
		Width:     1024,
		Height:    768,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Mac: &mac.Options{
			Appearance: mac.NSAppearanceNameDarkAqua,
		},
		OnStartup:     application.Startup,
		OnShutdown:    application.Shutdown,
		OnBeforeClose: application.OnBeforeClose,
		Bind: []interface{}{
			application,
		},
	})
	if err != nil {
		slog.Error("failed to run wails application", "error", err)
		os.Exit(1)
	}
}
