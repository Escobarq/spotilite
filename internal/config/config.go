package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	configFilename = "config-xpui.ini"
	prefsFilename  = "prefs"
)

// DefaultExtensionSuffix is appended to extension names lacking a .js/.mjs suffix
// when resolving them on disk.
const DefaultExtensionSuffix = ".js"

// Config mirrors the spicetify-cli config-xpui.ini shape. Only the keys Spotilite
// actually consumes are typed; unknown sections/keys are preserved verbatim so
// the file can round-trip without data loss when Save() is called.
type Config struct {
	// Setting
	SpotifyPath        string
	PrefsPath          string
	SpotifyLaunchFlags string
	CheckSpicetifyUpg  bool

	// Preprocesses
	DisableSentry     bool
	DisableUILogging  bool
	RemoveRTLRule     bool
	ExposeAPIs        bool
	DisableWinCasts   bool // legacy, kept for round-trip

	// AdditionalOptions
	CurrentTheme    string
	ColorScheme     string
	Extensions      []string // pipe-separated in ini, slice here
	CustomApps      []string // pipe-separated in ini, slice here
	InjectCSS       bool
	InjectThemeJS   bool
	ReplaceColors   bool
	OverwriteAssets bool
	SidebarConfig   bool
	HomeConfig      bool
	ExpFeatures     bool

	// CustomApps/Extensions user folder (resolved at Load time).
	// If empty, defaults to the Spicetify CLI folder when present.
	UserExtensionsDir string
	UserCustomAppsDir string
	UserThemesDir     string

	// Path the config was loaded from. Save() writes back here.
	path string

	// raw holds non-typed sections for round-trip preservation.
	raw sections
}

type sections map[string]map[string]string

// Lookup retrieves a raw string value from a section; ok=false if absent.
func (c *Config) Lookup(section, key string) (string, bool) {
	if c.raw == nil {
		return "", false
	}
	if sec, ok := c.raw[section]; ok {
		if v, ok := sec[key]; ok {
			return v, true
		}
	}
	return "", false
}

// EnabledExtensions returns the extension entries from [AdditionalOptions].Extensions
// whose trailing character is not "-". Spicetify uses "name-" to disable.
func (c *Config) EnabledExtensions() []string {
	out := make([]string, 0, len(c.Extensions))
	for _, e := range c.Extensions {
		if strings.HasSuffix(e, "-") {
			continue
		}
		out = append(out, e)
	}
	return out
}

// ExtensionEnabled reports whether the given name appears in the extensions list
// without the trailing "-" disable marker.
func (c *Config) ExtensionEnabled(name string) bool {
	for _, e := range c.Extensions {
		if strings.TrimSuffix(e, "-") == name {
			return !strings.HasSuffix(e, "-")
		}
	}
	return false
}

// SetExtension flips a single extension entry in the list, persisting on Save().
func (c *Config) SetExtension(name string, enabled bool) {
	// Strip any existing entry for this name (in both forms).
	filtered := make([]string, 0, len(c.Extensions))
	for _, e := range c.Extensions {
		if strings.TrimSuffix(e, "-") == name {
			continue
		}
		filtered = append(filtered, e)
	}
	if enabled {
		filtered = append(filtered, name)
	} else {
		filtered = append(filtered, name+"-")
	}
	c.Extensions = filtered
	if c.raw == nil {
		c.raw = sections{}
	}
	if c.raw["AdditionalOptions"] == nil {
		c.raw["AdditionalOptions"] = map[string]string{}
	}
	c.raw["AdditionalOptions"]["extensions"] = joinPipe(c.Extensions)
}

// SetCustomApp mirrors SetExtension for the custom_apps list.
func (c *Config) SetCustomApp(name string, enabled bool) {
	filtered := make([]string, 0, len(c.CustomApps))
	for _, e := range c.CustomApps {
		if strings.TrimSuffix(e, "-") == name {
			continue
		}
		filtered = append(filtered, e)
	}
	if enabled {
		filtered = append(filtered, name)
	} else {
		filtered = append(filtered, name+"-")
	}
	c.CustomApps = filtered
	if c.raw == nil {
		c.raw = sections{}
	}
	if c.raw["AdditionalOptions"] == nil {
		c.raw["AdditionalOptions"] = map[string]string{}
	}
	c.raw["AdditionalOptions"]["custom_apps"] = joinPipe(c.CustomApps)
}

// SetTheme updates current_theme + color_scheme in the raw map for Save().
func (c *Config) SetTheme(theme, scheme string) {
	c.CurrentTheme = theme
	c.ColorScheme = scheme
	if c.raw == nil {
		c.raw = sections{}
	}
	if c.raw["AdditionalOptions"] == nil {
		c.raw["AdditionalOptions"] = map[string]string{}
	}
	c.raw["AdditionalOptions"]["current_theme"] = theme
	c.raw["AdditionalOptions"]["color_scheme"] = scheme
}

// Path returns the file path the configuration was loaded from.
func (c *Config) Path() string { return c.path }

// Load reads the Spicetify config file. It tries, in order:
//  1. SPICETIFY_CONFIG env var (explicit path)
//  2. The Spicetify CLI user folder (%APPDATA%/spicetify on Windows,
//     $XDG_CONFIG_HOME/spicetify or ~/.config/spicetify on Unix)
//  3. The Spotilite user folder (os.UserConfigDir()/spotilite)
//
// Returns an error only if parsing fails. A missing file is NOT an error:
// the returned Config is zero-valued with path set to the preferred location.
func Load() (*Config, error) {
	candidates := candidatePaths()
	c := &Config{raw: sections{}}

	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if err := parseINI(string(data), c.raw); err != nil {
			return nil, fmt.Errorf("config %s: %w", p, err)
		}
		c.path = p
		c.applyRaw()
		c.resolveUserDirs()
		return c, nil
	}

	// No file found: set path to first candidate so Save() has somewhere to go.
	c.path = candidates[0]
	c.resolveUserDirs()
	return c, nil
}

// MustLoad is a convenience that returns an empty-but-valid config.
// Errors returned are logged but the returned Config is always usable.
func MustLoad() *Config {
	c, err := Load()
	if err != nil {
		c.path = candidatePaths()[0]
		c.raw = sections{}
		c.resolveUserDirs()
	}
	return c
}

// Save writes the config back to disk at Config.Path(), creating the parent
// directory if needed.
func (c *Config) Save() error {
	if c.path == "" {
		return fmt.Errorf("config: no path set")
	}
	c.applyToRaw()
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(c.path, []byte(c.serialize()), 0o644)
}

// --- path resolution ---------------------------------------------------------

func candidatePaths() []string {
	if env := os.Getenv("SPICETIFY_CONFIG"); env != "" {
		return []string{env}
	}
	out := make([]string, 0, 3)
	if p := spicetifyConfigDir(); p != "" {
		out = append(out, filepath.Join(p, configFilename))
	}
	if ucd, err := os.UserConfigDir(); err == nil {
		out = append(out, filepath.Join(ucd, "spotilite", configFilename))
	}
	if len(out) == 0 {
		// Fallback: ~/.spotilite (works even when UserConfigDir fails).
		if home, err := os.UserHomeDir(); err == nil {
			out = append(out, filepath.Join(home, ".spotilite", configFilename))
		}
	}
	return out
}

func spicetifyConfigDir() string {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return ""
		}
		return filepath.Join(appdata, "spicetify")
	}
	// posix: respect XDG_CONFIG_HOME, fall back to ~/.config
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "spicetify")
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "spicetify")
	}
	return ""
}

func spicetifyUserDataDir() string {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return ""
		}
		return appdata
	}
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	return ""
}

func (c *Config) resolveUserDirs() {
	base := spicetifyUserDataDir()
	if base != "" {
		c.UserExtensionsDir = filepath.Join(base, "spicetify", "Extensions")
		c.UserCustomAppsDir = filepath.Join(base, "spicetify", "CustomApps")
		c.UserThemesDir = filepath.Join(base, "spicetify", "Themes")
	} else if ucd, err := os.UserConfigDir(); err == nil {
		c.UserExtensionsDir = filepath.Join(ucd, "spotilite", "Extensions")
		c.UserCustomAppsDir = filepath.Join(ucd, "spotilite", "CustomApps")
		c.UserThemesDir = filepath.Join(ucd, "spotilite", "Themes")
	}
}

// --- apply raw -> typed -------------------------------------------------------

func (c *Config) applyRaw() {
	get := func(section, key string) string {
		if sec, ok := c.raw[section]; ok {
			return sec[key]
		}
		return ""
	}
	getBool := func(section, key string) bool {
		return strings.TrimSpace(get(section, key)) == "1"
	}

	c.SpotifyPath = get("Setting", "spotify_path")
	c.PrefsPath = get("Setting", "prefs_path")
	c.SpotifyLaunchFlags = get("Setting", "spotify_launch_flags")
	c.CheckSpicetifyUpg = getBool("Setting", "check_spicetify_upgrade")

	c.DisableSentry = getBool("Preprocesses", "disable_sentry")
	c.DisableUILogging = getBool("Preprocesses", "disable_ui_logging")
	c.RemoveRTLRule = getBool("Preprocesses", "remove_rtl_rule")
	c.ExposeAPIs = getBool("Preprocesses", "expose_apis")
	c.DisableWinCasts = getBool("Preprocesses", "disable_win_casts")

	c.CurrentTheme = get("AdditionalOptions", "current_theme")
	c.ColorScheme = get("AdditionalOptions", "color_scheme")
	c.Extensions = splitPipe(get("AdditionalOptions", "extensions"))
	c.CustomApps = splitPipe(get("AdditionalOptions", "custom_apps"))
	c.InjectCSS = getBool("AdditionalOptions", "inject_css")
	c.InjectThemeJS = getBool("AdditionalOptions", "inject_theme_js")
	c.ReplaceColors = getBool("AdditionalOptions", "replace_colors")
	c.OverwriteAssets = getBool("AdditionalOptions", "overwrite_assets")
	c.SidebarConfig = getBool("AdditionalOptions", "sidebar_config")
	c.HomeConfig = getBool("AdditionalOptions", "home_config")
	c.ExpFeatures = getBool("AdditionalOptions", "experimental_features")
}

// --- apply typed -> raw -------------------------------------------------------

func (c *Config) applyToRaw() {
	if c.raw == nil {
		c.raw = sections{}
	}
	setStr := func(section, key, value string) {
		if c.raw[section] == nil {
			c.raw[section] = map[string]string{}
		}
		c.raw[section][key] = value
	}
	setBool := func(section, key string, value bool) {
		setStr(section, key, boolToString(value))
	}

	setStr("Setting", "spotify_path", c.SpotifyPath)
	setStr("Setting", "prefs_path", c.PrefsPath)
	setStr("Setting", "spotify_launch_flags", c.SpotifyLaunchFlags)
	setBool("Setting", "check_spicetify_upgrade", c.CheckSpicetifyUpg)

	setBool("Preprocesses", "disable_sentry", c.DisableSentry)
	setBool("Preprocesses", "disable_ui_logging", c.DisableUILogging)
	setBool("Preprocesses", "remove_rtl_rule", c.RemoveRTLRule)
	setBool("Preprocesses", "expose_apis", c.ExposeAPIs)
	setBool("Preprocesses", "disable_win_casts", c.DisableWinCasts)

	setStr("AdditionalOptions", "current_theme", c.CurrentTheme)
	setStr("AdditionalOptions", "color_scheme", c.ColorScheme)
	setStr("AdditionalOptions", "extensions", joinPipe(c.Extensions))
	setStr("AdditionalOptions", "custom_apps", joinPipe(c.CustomApps))
	setBool("AdditionalOptions", "inject_css", c.InjectCSS)
	setBool("AdditionalOptions", "inject_theme_js", c.InjectThemeJS)
	setBool("AdditionalOptions", "replace_colors", c.ReplaceColors)
	setBool("AdditionalOptions", "overwrite_assets", c.OverwriteAssets)
	setBool("AdditionalOptions", "sidebar_config", c.SidebarConfig)
	setBool("AdditionalOptions", "home_config", c.HomeConfig)
	setBool("AdditionalOptions", "experimental_features", c.ExpFeatures)
}

// --- helpers ------------------------------------------------------------------

func splitPipe(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func joinPipe(parts []string) string {
	return strings.Join(parts, "|")
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
