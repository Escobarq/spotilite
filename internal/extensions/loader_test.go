package extensions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadJS(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ext.js"), []byte("console.log('hi');"), 0o644); err != nil {
		t.Fatal(err)
	}
	l := NewLoader([]string{dir}, "http://localhost:8765")
	items, skipped, err := l.Load([]string{"ext"})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(items) != 1 || len(skipped) != 0 {
		t.Fatalf("got items=%v skipped=%v", items, skipped)
	}
	if items[0].IsModule {
		t.Error("expected classic script")
	}
	if items[0].Source != "console.log('hi');" {
		t.Errorf("source = %q", items[0].Source)
	}
}

func TestLoadMJSRevision(t *testing.T) {
	dir := t.TempDir()
	src := `// spicetify_map {"react"}{"./node_modules/react/index.js"}
import React from "react";
console.log("hi");
`
	if err := os.WriteFile(filepath.Join(dir, "ext.mjs"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	l := NewLoader([]string{dir}, "http://localhost:8765")
	items, _, err := l.Load([]string{"ext"})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(items) != 1 {
		t.Fatal("missing item")
	}
	if !items[0].IsModule {
		t.Error("expected module")
	}
	if strings.Contains(items[0].Source, "spicetify_map") {
		t.Errorf("spicetify_map comment should have been stripped:\n%s", items[0].Source)
	}
}

func TestLoadMissingSkipped(t *testing.T) {
	dir := t.TempDir()
	l := NewLoader([]string{dir}, "http://localhost:8765")
	_, skipped, err := l.Load([]string{"nope", "alsoNope"})
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 2 {
		t.Fatalf("expected 2 skipped, got %v", skipped)
	}
}

func TestCandidateFilenames(t *testing.T) {
	cases := map[string][]string{
		"foo":         {"foo.js", "foo.mjs"},
		"bar.js":      {"bar.js"},
		"bar.mjs":     {"bar.mjs"},
		"baz/qux.js":  {"qux.js"},
	}
	for in, want := range cases {
		got := candidateFilenames(in)
		if len(got) != len(want) {
			t.Errorf("candidateFilenames(%q) = %v want %v", in, got, want)
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("candidateFilenames(%q)[%d] = %q want %q", in, i, got[i], want[i])
			}
		}
	}
}

func TestJSBundleWrapsEach(t *testing.T) {
	items := []Loaded{
		{Name: "a", Source: "var x=1;"},
		{Name: "b", Source: "var y=2;"},
	}
	out := JSBundle(items)
	if !strings.Contains(out, "extension: a") {
		t.Errorf("missing tag for extension a:\n%s", out)
	}
	if !strings.Contains(out, "extension: b") {
		t.Errorf("missing tag for extension b:\n%s", out)
	}
	if strings.Count(out, "catch (e)") != 2 {
		t.Errorf("expected 2 catch blocks, got %d", strings.Count(out, "catch"))
	}
}

func TestLoadAndBundleEndToEnd(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "first.js"), []byte("var x=1;"), 0o644)
	os.WriteFile(filepath.Join(dir, "second.js"), []byte("var y=2;"), 0o644)
	l := NewLoader([]string{dir}, "")
	script, skipped, err := l.LoadAndBundle([]string{"first", "second", "missing"})
	if err != nil {
		t.Fatalf("LoadAndBundle: %v", err)
	}
	if !strings.Contains(script, "extension: first") || !strings.Contains(script, "extension: second") {
		t.Errorf("expected both bundled:\n%s", script)
	}
	if len(skipped) != 1 || skipped[0] != "missing" {
		t.Errorf("expected ['missing'] got %v", skipped)
	}
}
