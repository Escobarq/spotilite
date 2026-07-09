package themes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleColorINI = `
; spicetify-cli compatible color.ini
[base]
text     = #FFFFFF
subtext  = #B3B3B3
button   = #1DB954
button-active = #1ed760
main     = #121212
sidebar  = #000000
player   = #181818
card     = #282828

[light]
text     = #000000
main     = #FFFFFF
button   = #1DB954
`

func writeTheme(t *testing.T, dir, name string, colorINI string, userCSS string) {
	t.Helper()
	tDir := filepath.Join(dir, name)
	if err := os.MkdirAll(tDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tDir, "color.ini"), []byte(colorINI), 0o644); err != nil {
		t.Fatal(err)
	}
	if userCSS != "" {
		if err := os.WriteFile(filepath.Join(tDir, "user.css"), []byte(userCSS), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadTheme(t *testing.T) {
	dir := t.TempDir()
	writeTheme(t, dir, "TestTheme", sampleColorINI, "")

	m := NewManager([]string{dir})
	tm, err := m.Load("TestTheme")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tm.Schemes) != 2 {
		t.Fatalf("schemes = %v, want 2", tm.SchemeNames())
	}
	s := tm.Scheme("base")
	if s == nil {
		t.Fatal("missing base scheme")
	}
	if s.Colors["button"] != "#1DB954" {
		t.Errorf("button = %q want #1DB954", s.Colors["button"])
	}
	// Fallback: empty scheme name -> first scheme in map iteration order
	// (deterministic enough for the documented contract).
	if fallback := tm.Scheme(""); fallback == nil {
		t.Error("empty scheme name returned nil; expected fallback")
	}
}

func TestLoadMissingTheme(t *testing.T) {
	dir := t.TempDir()
	m := NewManager([]string{dir})
	if _, err := m.Load("Nope"); err == nil {
		t.Error("expected error for missing theme")
	}
}

func TestListNames(t *testing.T) {
	dir := t.TempDir()
	writeTheme(t, dir, "A", sampleColorINI, "")
	writeTheme(t, dir, "B", sampleColorINI, "")
	// Folder without color.ini should be skipped.
	if err := os.MkdirAll(filepath.Join(dir, "NoColor"), 0o755); err != nil {
		t.Fatal(err)
	}
	m := NewManager([]string{dir})
	names := m.ListNames()
	if len(names) != 2 {
		t.Fatalf("got %v want 2 names", names)
	}
}

func TestSpiceCSSIncludesMapAndTokens(t *testing.T) {
	dir := t.TempDir()
	writeTheme(t, dir, "T", sampleColorINI, "")
	m := NewManager([]string{dir})
	tm, err := m.Load("T")
	if err != nil {
		t.Fatal(err)
	}
	css := tm.SpiceCSS("base")
	if !strings.Contains(css, "--spice-button: #1DB954") {
		t.Errorf("missing spice-button declaration:\n%s", css)
	}
	if !strings.Contains(css, "--text-base: var(--spice-text)") {
		t.Errorf("missing token mapping for --text-base:\n%s", css)
	}
	if !strings.Contains(css, "--essential-bright-accent: var(--spice-button)") {
		t.Errorf("missing token mapping for button accent:\n%s", css)
	}
}

func TestFullCSSAppendsUserCSS(t *testing.T) {
	dir := t.TempDir()
	writeTheme(t, dir, "T", sampleColorINI, "/* my custom rules */\n.foo { color: red; }")
	m := NewManager([]string{dir})
	tm, err := m.Load("T")
	if err != nil {
		t.Fatal(err)
	}
	css := tm.FullCSS("base")
	if !strings.Contains(css, "user.css") {
		t.Errorf("missing user.css sentinel:\n%s", css)
	}
	if !strings.Contains(css, ".foo { color: red; }") {
		t.Errorf("user.css content lost:\n%s", css)
	}
}
