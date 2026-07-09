package spicetify

import (
	"os/exec"
	"strings"
	"testing"
)

// TestBundleSyntax pipes the concatenated Bundle() output through
// `node --check` to catch IIFE/ordering issues that per-file checks miss.
// Skips when node is not on PATH.
func TestBundleSyntax(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not available")
	}
	src := Bundle()
	if strings.TrimSpace(src) == "" {
		t.Fatal("Bundle returned empty string; embeds missing?")
	}
	cmd := exec.Command("node", "--check")
	cmd.Stdin = strings.NewReader(src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bundle fails node --check: %v\n--- output ---\n%s", err, out)
	}
}

// TestBundleOrder verifies core.js precedes uri.js in the bundle (core owns the
// window.Spicetify object that later files attach to).
func TestBundleOrder(t *testing.T) {
	src := Bundle()
	coreIdx := strings.Index(src, "initSpicetify()")
	uriIdx := strings.Index(src, "initSpicetifyURI")
	playerIdx := strings.Index(src, "initSpicetifyPlayer")
	cosmosIdx := strings.Index(src, "initSpicetifyCosmos")
	if coreIdx < 0 || uriIdx < 0 || playerIdx < 0 || cosmosIdx < 0 {
		t.Fatalf("missing sub-init markers in bundle (core=%d uri=%d player=%d cosmos=%d)",
			coreIdx, uriIdx, playerIdx, cosmosIdx)
	}
	if !(coreIdx < uriIdx) {
		t.Errorf("core.js must precede uri.js: core=%d uri=%d", coreIdx, uriIdx)
	}
	if !(uriIdx < playerIdx) {
		t.Errorf("uri.js must precede player.js: uri=%d player=%d", uriIdx, playerIdx)
	}
	if !(cosmosIdx < playerIdx) {
		t.Errorf("cosmos.js must precede player.js: cosmos=%d player=%d", cosmosIdx, playerIdx)
	}
}
