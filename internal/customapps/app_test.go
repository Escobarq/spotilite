package customapps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeApp(t *testing.T, dir, name string, manifestJSON, indexJS, styleCSS string) {
	t.Helper()
	d := filepath.Join(dir, name)
	if err := os.MkdirAll(d, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "manifest.json"), []byte(manifestJSON), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "index.js"), []byte(indexJS), 0o644); err != nil {
		t.Fatal(err)
	}
	if styleCSS != "" {
		if err := os.WriteFile(filepath.Join(d, "style.css"), []byte(styleCSS), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const sampleManifest = `{"name":"MyApp","icon":"<svg></svg>","active-icon":"<svg></svg>"}`
const sampleIndexJS = `function render(container) { container.innerHTML = "<h1>Hello</h1>"; }`

func TestLoadCustomApp(t *testing.T) {
	dir := t.TempDir()
	writeApp(t, dir, "my-app", sampleManifest, sampleIndexJS, ".foo{color:red}")
	m := NewManager([]string{dir})
	app, err := m.Load("my-app")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if app.Name != "MyApp" {
		t.Errorf("name = %q want MyApp", app.Name)
	}
	if !app.HasRender {
		t.Error("expected hasRender to be true")
	}
	if app.StylesCSS != ".foo{color:red}" {
		t.Errorf("stylesCSS = %q", app.StylesCSS)
	}
}

func TestJSInjectionContainsRoute(t *testing.T) {
	app := &App{
		Name: "MyApp",
		Dir:  t.TempDir(),
		Manifest: &Manifest{
			Name: "MyApp",
			Icon: "<svg></svg>",
		},
		JSBundle:  "function render(c){};",
		HasRender: true,
	}
	injection := app.JSInjection()
	if !strings.Contains(injection, "/spotilite/MyApp") {
		t.Errorf("missing route in injection:\n%s", injection)
	}
	if !strings.Contains(injection, "document.getElementById") {
		t.Errorf("missing container creation:\n%s", injection)
	}
	if !strings.Contains(injection, "handlePopState") {
		t.Errorf("missing popstate handler:\n%s", injection)
	}
}

func TestLoadMissingApp(t *testing.T) {
	m := NewManager([]string{t.TempDir()})
	if _, err := m.Load("nonexistent"); err == nil {
		t.Error("expected error for missing app")
	}
}

func TestListNames(t *testing.T) {
	dir := t.TempDir()
	writeApp(t, dir, "a", sampleManifest, sampleIndexJS, "")
	writeApp(t, dir, "b", sampleManifest, sampleIndexJS, "")
	os.MkdirAll(filepath.Join(dir, "noManifest"), 0o755)
	m := NewManager([]string{dir})
	names := m.ListNames()
	if len(names) != 2 {
		t.Fatalf("got %v want 2 names", names)
	}
}