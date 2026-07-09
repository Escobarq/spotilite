// Package themes provides Spicetify-compatible theme loading for the
// https://open.spotify.com web player. It parses a theme's color.ini into a
// set of named schemes, exposes CSS variables using the canonical --spice-*
// naming convention, and translates those into the web player's own design
// tokens (--background-base, --text-base, --essential-bright-accent, ...) so
// that themes authored for the desktop client apply meaningful color changes
// without rewriting Spotify's bundled CSS.
//
// On the desktop Spicetify client the preprocessing step ([Patch] + the
// colorVariableReplace regexes in preprocess.go) rewrites Spotify's hardcoded
// hex color literals to var(--spice-*) references. The web player's bundle is
// served by Spotify and cannot be rewritten, so this package emits a small
// CSS shim that overrides the web player's own tokens at the :root level.
//
// The token mapping grows as needed; unmapped --spice-* variables are still
// emitted (extensions/themes can target them directly via user.css).
package themes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spotilite/internal/ini"
)

// Scheme holds the color values for a named scheme in a color.ini file.
// Variable names use the canonical Spicetify --spice-* keys so themes authored
// against spicetify-cli map 1:1.
type Scheme struct {
	Name   string
	Colors map[string]string // variable name without "spice-" prefix (e.g. "button" -> "#1db954")
}

// Theme is a Spicetify theme loaded from disk. A theme directory contains
// color.ini (required) and optionally user.css, theme.js and an assets/
// subfolder. The Name field is the folder name; schemes are keyed by name.
type Theme struct {
	Name     string    // folder name
	Dir      string    // absolute path on disk
	Schemes  map[string]*Scheme
	UserCSS  string   // contents of user.css if present, else ""
	ThemeJS  string   // contents of theme.js if present (used as an extension)
	colorIni ini.Sections
}

// Manager resolves and caches themes from one or more search directories.
// It first looks in the user themes folder (~/spicetify/Themes by default;
// matches the spicetify-cli convention), then falls back to a built-in
// themes blob bundled with the app.
type Manager struct {
	searchDirs []string
	cache      map[string]*Theme
}

// NewManager builds a Manager that resolves themes from the given directories,
// tried in order (first hit wins).
func NewManager(searchDirs []string) *Manager {
	return &Manager{searchDirs: searchDirs, cache: make(map[string]*Theme)}
}

// Load returns a Theme by folder name. Returns an error if no theme directory
// exists in any search dir, or if its color.ini is missing or unparseable.
// Successful loads are cached by theme name.
func (m *Manager) Load(name string) (*Theme, error) {
	if name == "" {
		return nil, fmt.Errorf("themes: empty name")
	}
	if t, ok := m.cache[name]; ok {
		return t, nil
	}
	for _, dir := range m.searchDirs {
		t, err := loadFromDir(filepath.Join(dir, name))
		if err == nil {
			t.Name = name
			m.cache[name] = t
			return t, nil
		}
		// Theme not found in this dir; try the next. If dir doesn't exist,
		// silently skip.
	}
	return nil, fmt.Errorf("themes: %q not found in %v", name, m.searchDirs)
}

// ListNames returns the set of theme folder names available across all
// search directories (deduplicated, preserving first-seen order).
func (m *Manager) ListNames() []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, dir := range m.searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			n := e.Name()
			if seen[n] {
				continue
			}
			// Cheap presence ping: must contain a color.ini.
			if _, err := os.Stat(filepath.Join(dir, n, "color.ini")); err != nil {
				continue
			}
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

func loadFromDir(dir string) (*Theme, error) {
	colorPath := filepath.Join(dir, "color.ini")
	colorData, err := os.ReadFile(colorPath)
	if err != nil {
		return nil, err
	}
	secs := ini.Sections{}
	if err := ini.Parse(string(colorData), secs); err != nil {
		return nil, fmt.Errorf("themes: parse color.ini %s: %w", dir, err)
	}

	t := &Theme{
		Dir:      dir,
		Schemes:  make(map[string]*Scheme),
		colorIni: secs,
	}

	for schemeName, kv := range secs {
		colors := make(map[string]string, len(kv))
		for k, v := range kv {
			colors[strings.TrimSpace(k)] = strings.TrimSpace(stripComment(v))
		}
		t.Schemes[schemeName] = &Scheme{Name: schemeName, Colors: colors}
	}

	if css, err := os.ReadFile(filepath.Join(dir, "user.css")); err == nil {
		t.UserCSS = string(css)
	}
	if js, err := os.ReadFile(filepath.Join(dir, "theme.js")); err == nil {
		t.ThemeJS = string(js)
	}

	return t, nil
}

// stripComment removes trailing ; or # comments after a color value in
// color.ini (spicetify-cli colorVariableReplace supports both). Stops at the
// first comment marker because color values themselves are hex literals.
func stripComment(v string) string {
	for i := 0; i < len(v); i++ {
		c := v[i]
		if c == ';' || c == '#' && i > 0 {
			// Keep '#' when it's the very first char (it's the hex color).
			// Actually # at position 0 is the color itself; we never strip it.
			return strings.TrimSpace(v[:i])
		}
	}
	return strings.TrimSpace(v)
}

// Scheme returns the named scheme, falling back to the first scheme in the
// file if the requested name is empty or unknown.
func (t *Theme) Scheme(name string) *Scheme {
	if s, ok := t.Schemes[name]; ok {
		return s
	}
	for _, s := range t.Schemes {
		return s
	}
	return nil
}

// SchemeNames returns the sorted scheme names declared in color.ini.
func (t *Theme) SchemeNames() []string {
	out := make([]string, 0, len(t.Schemes))
	for n := range t.Schemes {
		out = append(out, n)
	}
	// stable order: alphabetical
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}

// SpiceCSS returns the CSS string that, when injected into the page, exposes:
//
//   1. :root { --spice-*: <color> !important; ... }  (Spicetify canonical variables)
//   2. :root { --<web-player-token>: var(--spice-*) !important; ... }  (token mapping)
//
// The !important flags are required because the Spotify web player uses
// inline styles with !important on many elements, which otherwise override
// our :root CSS variable declarations.
func (t *Theme) SpiceCSS(schemeName string) string {
	scheme := t.Scheme(schemeName)
	if scheme == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(":root {\n")
	for k, v := range scheme.Colors {
		b.WriteString("  --spice-")
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString(" !important;\n")
	}
	b.WriteString("}\n")

	b.WriteString(":root {\n")
	for spiceVar, webTokens := range TokenMap {
		if color, ok := scheme.Colors[spiceVar]; ok && color != "" {
			for _, webToken := range webTokens {
				b.WriteString("  ")
				b.WriteString(webToken)
				b.WriteString(": var(--spice-")
				b.WriteString(spiceVar)
				b.WriteString(") !important;\n")
			}
		}
	}
	b.WriteString("}\n")
	return b.String()
}

// FullCSS returns SpiceCSS(schemeName) followed by the theme's user.css, if
// any. The user.css is appended after so author rules win.
func (t *Theme) FullCSS(schemeName string) string {
	out := t.SpiceCSS(schemeName)
	if t.UserCSS != "" {
		out += "\n/* user.css */\n" + t.UserCSS + "\n"
	}
	return out
}
