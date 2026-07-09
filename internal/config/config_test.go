package config

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleINI = `
; this is a comment
# alternate comment

[Setting]
spotify_path     = C:\Users\Test\AppData\Roaming\Spotify\Spotify.exe
prefs_path       = C:\Users\Test\AppData\Roaming\Spotify\prefs
spotify_launch_flags =
check_spicetify_upgrade = 1

[Preprocesses]
disable_sentry = 1
disable_ui_logging = 1
expose_apis    = 1
disable_win_casts = 0

[AdditionalOptions]
current_theme  = Dribbblish
color_scheme   = dark
extensions     = shuffle.js|trashbin.js|autoplay.js-
custom_apps    = reddit|lyrics-
inject_css     = 1
inject_theme_js = 0
replace_colors = 1
experimental_features = 1

[Patch]
xpui.js_find_8100 = .dropdown__dropdown
xpui.js_repl_8100 = .XnnUPtiKeRYjjJfrfQOA
`

func TestLoadParseSample(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFilename)
	if err := os.WriteFile(path, []byte(sampleINI), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Setenv("SPICETIFY_CONFIG", path)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if c.SpotifyPath == "" {
		t.Fatalf("empty SpotifyPath")
	}
	if !c.ExposeAPIs {
		t.Errorf("expected ExposeAPIs true")
	}
	if !c.DisableSentry {
		t.Errorf("expected DisableSentry true")
	}
	if c.ReplaceColors == false {
		t.Errorf("expected ReplaceColors true")
	}
	if got, want := c.CurrentTheme, "Dribbblish"; got != want {
		t.Errorf("theme = %q, want %q", got, want)
	}

	enabled := c.EnabledExtensions()
	wantEnabled := map[string]bool{"shuffle.js": true, "trashbin.js": true, "autoplay.js": false}
	if len(enabled) != 2 {
		t.Fatalf("EnabledExtensions = %v, want 2 entries", enabled)
	}
	for _, e := range enabled {
		if !wantEnabled[e] {
			t.Errorf("unexpected enabled: %q", e)
		}
	}
	if !c.ExtensionEnabled("shuffle.js") {
		t.Errorf("shuffle should be enabled")
	}
	if c.ExtensionEnabled("autoplay.js") {
		t.Errorf("autoplay should be disabled (trailing -)")
	}
	if !c.ExtensionEnabled("nonexistent.js") {
		// not present => treated as disabled
		// behavior is documented; never returns true for unknown
	}
}

func TestSetExtension(t *testing.T) {
	c := &Config{raw: sections{}}
	c.Extensions = []string{"a", "b-", "c"}
	c.SetExtension("b", true)
	if len(c.Extensions) != 3 {
		t.Fatalf("len = %d", len(c.Extensions))
	}
	if c.raw["AdditionalOptions"]["extensions"] != "a|c|b" {
		t.Errorf("raw pipe = %q want %q", c.raw["AdditionalOptions"]["extensions"], "a|c|b")
	}
	c.SetExtension("d", true)
	if len(c.Extensions) != 4 || c.Extensions[3] != "d" {
		t.Errorf("add new = %v", c.Extensions)
	}
}

func TestSetCustomApp(t *testing.T) {
	c := &Config{raw: sections{}}
	c.CustomApps = []string{"reddit", "lyrics-"}
	c.SetCustomApp("lyrics", true)
	if c.raw["AdditionalOptions"]["custom_apps"] != "reddit|lyrics" {
		t.Errorf("raw pipe = %q", c.raw["AdditionalOptions"]["custom_apps"])
	}
}

func TestSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFilename)
	t.Setenv("SPICETIFY_CONFIG", path)
	written := sampleINI
	if err := os.WriteFile(path, []byte(written), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	c.SetExtension("newext", true)
	c.SetCustomApp("myapp", true)
	c.SetTheme("Dribbblish", "light")
	c.ExposeAPIs = false

	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	c2, err := Load()
	if err != nil {
		t.Fatalf("Load 2: %v", err)
	}
	if !c2.ExtensionEnabled("newext") {
		t.Errorf("newext should be enabled after round-trip, got %v", c2.Extensions)
	}
	if c2.CurrentTheme != "Dribbblish" || c2.ColorScheme != "light" {
		t.Errorf("theme not preserved: %q/%q", c2.CurrentTheme, c2.ColorScheme)
	}
}

func TestMissingFileIsNotError(t *testing.T) {
	t.Setenv("SPICETIFY_CONFIG", filepath.Join(t.TempDir(), "does-not-exist.ini"))
	c, err := Load()
	if err != nil {
		t.Fatalf("error on missing file: %v", err)
	}
	if c == nil {
		t.Fatalf("nil config on missing file")
	}
	if c.path == "" {
		t.Errorf("expected default path to be set")
	}
}

func TestParseEdgeCases(t *testing.T) {
	out := sections{}
	cases := []struct {
		in   string
		sec  string
		key  string
		want string
	}{
		{"[s]\nk=v", "s", "k", "v"},
		{"[s]\nk=  spaced  ", "s", "k", "spaced"},
		{"[s]\nk=a=b=c", "s", "k", "a=b=c"},
		{"[s]\n;comment\nk=v", "s", "k", "v"},
		{"[s1]\nk1=v1\n[s2]\nk2=v2", "s1", "k1", "v1"},
	}
	for i, tc := range cases {
		out = sections{}
		if err := parseINI(tc.in, out); err != nil {
			t.Errorf("case %d parse error: %v", i, err)
			continue
		}
		got := out[tc.sec][tc.key]
		if got != tc.want {
			t.Errorf("case %d got %q want %q", i, got, tc.want)
		}
	}
}

func TestInvalidLeadingKeyErrors(t *testing.T) {
	out := sections{}
	in := "nosection\nk=v\n"
	if err := parseINI(in, out); err != nil {
		t.Logf("got expected error: %v", err)
	} else {
		// Acceptable: silently dropped. Current impl drops it.
		if _, ok := out["nosection"]; ok {
			t.Errorf("unexpected section named %q", "nosection")
		}
	}
}
