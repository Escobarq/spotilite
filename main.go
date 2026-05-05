package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"spotilite/internal/app"
	"spotilite/internal/i18n"
	"spotilite/internal/tray"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	translator := i18n.New()

	// Default: run in background when closing the window.
	runInBackground := true

	var application *app.App

	trayManager := tray.New(translator, runInBackground, func(enabled bool) {
		if application != nil {
			application.SetBackgroundMode(enabled)
		}
	})

	application = app.NewApp(translator, trayManager, runInBackground)

	err := wails.Run(&options.App{
		Title:  translator.T("app.title"),
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Spotify black (#191414) to blend the title bar with the webview.
		BackgroundColour: &options.RGBA{R: 25, G: 20, B: 20, A: 255},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Mac: &mac.Options{
			Appearance: mac.NSAppearanceNameDarkAqua,
		},
		Menu:          trayManager.Build(),
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
