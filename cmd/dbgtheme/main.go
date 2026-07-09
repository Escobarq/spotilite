//go:build debugtheme

package main

import (
	"fmt"
	"os"

	"spotilite/internal/config"
	"spotilite/internal/themes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}
	fmt.Println("config path:", cfg.Path())
	fmt.Println("theme:", cfg.CurrentTheme, "scheme:", cfg.ColorScheme)
	fmt.Println("extensions:", cfg.EnabledExtensions())
	fmt.Println("extensionsDir:", cfg.UserExtensionsDir)
	fmt.Println("themesDir:", cfg.UserThemesDir)

	m := themes.NewManager([]string{cfg.UserThemesDir})
	names := m.ListNames()
	fmt.Println("available themes:", names)

	if cfg.CurrentTheme == "" {
		fmt.Println("no theme configured")
		os.Exit(0)
	}
	t, err := m.Load(cfg.CurrentTheme)
	if err != nil {
		fmt.Println("theme load error:", err)
		os.Exit(1)
	}
	fmt.Println("schemes:", t.SchemeNames())
	css := t.SpiceCSS(cfg.ColorScheme)
	fmt.Printf("\n=== Generated SpiceCSS (first 600 chars) ===\n%s\n...===\n", truncate(css, 600))
	fmt.Println("Full CSS length:", len(t.FullCSS(cfg.ColorScheme)), "chars")
	if t.UserCSS != "" {
		fmt.Println("user.css loaded:", len(t.UserCSS), "chars")
	}
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
