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

	"spotilite/internal/api"
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
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		iconPath = filepath.Join(".", "build", "windows", "icon.ico")
	}

	trayIconPath := filepath.Join(filepath.Dir(exe), "build", "windows", "icon_tray.ico")
	if _, err := os.Stat(trayIconPath); os.IsNotExist(err) {
		trayIconPath = filepath.Join(".", "build", "windows", "icon_tray.ico")
	}

	notificationIconPath := filepath.Join(filepath.Dir(exe), "build", "windows", "icon_notification.png")
	if _, err := os.Stat(notificationIconPath); os.IsNotExist(err) {
		notificationIconPath = iconPath
	}

	var application *app.App

	apiServer := api.NewServer(nil) // Handler will be set after App is created

	trayManager := apptray.NewManager(
		translator,
		trayIconPath,
		func() {
			if application != nil {
				application.ToggleWindowVisibility()
			}
		},
		func() {
			if application != nil {
				application.ForceQuit()
			} else {
				os.Exit(0)
			}
		},
		func(module string, enabled bool) {
			if application != nil {
				application.SetModuleEnabled(module, enabled)
			}
		},
		func(enabled bool) {
			if application != nil {
				application.SetProxyEnabled(enabled)
			}
		},
		func(enabled bool) {
			if application != nil {
				application.SetBackgroundMode(enabled)
			}
		},
		func() bool {
			if application != nil {
				return application.IsBackgroundMode()
			}
			return true
		},
	)

	application = app.NewApp(translator, trayManager, apiServer, runInBackground, notificationIconPath)
	apiServer.SetHandler(application)

	err = wails.Run(&options.App{
		Title:  translator.T("app.title"),
		Width:  960,
		Height: 640,
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
		OnStartup:    application.Startup,
		OnShutdown:   application.Shutdown,
		OnBeforeClose: application.OnBeforeClose,
		Bind: []interface{}{
			application,
		},
	})
	if err != nil {
		slog.Error("failed to run wails application", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}
