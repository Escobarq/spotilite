package spicetify

import (
	"os/exec"
	"strings"
	"testing"

	"spotilite/internal/config"
	"spotilite/internal/customapps"
	"spotilite/internal/extensions"
)

// TestBundleSyntax pipes the concatenated Bundle() output through
// `node --check` to catch IIFE/ordering issues that per-file checks miss.
// Skips when node is not on PATH.
func TestBundleSyntax(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not available")
	}
	src := Bundle(BuildContext{})
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
	src := Bundle(BuildContext{})
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

// TestBundleIncludesConfig verifies that Bundle() with a populated BuildContext
// emits the Spicetify.Config assignment *after* the shims and the extensions.
func TestBundleIncludesConfig(t *testing.T) {
	cfg := &config.Config{CurrentTheme: "Dribbblish", ColorScheme: "dark"}
	src := Bundle(BuildContext{Config: cfg, LocalAPI: "http://localhost:8765"})
	idxShim := strings.Index(src, "initSpicetify()")
	idxConfig := strings.Index(src, "Spicetify.Config")
	idxLocalAPI := strings.Index(src, "Spicetify.localAPI")
	if idxShim < 0 || idxConfig < 0 || idxLocalAPI < 0 {
		t.Fatalf("missing elements (shim=%d config=%d api=%d)", idxShim, idxConfig, idxLocalAPI)
	}
	if !(idxShim < idxConfig) {
		t.Errorf("shims must precede Config: shim=%d config=%d", idxShim, idxConfig)
	}
	if !(idxConfig < idxLocalAPI) {
		t.Errorf("Config must precede localAPI assignment: config=%d api=%d", idxConfig, idxLocalAPI)
	}
	if !strings.Contains(src, "Dribbblish") {
		t.Error("current_theme value missing from injected config")
	}
}

// TestBundleIncludesHelpers verifies the three spicetify-cli helpers are
// embedded only when their respective config flag is set.
func TestBundleIncludesHelpers(t *testing.T) {
	// All off: bundle should not contain any helper body.
	srcOff := Bundle(BuildContext{Config: &config.Config{}})
	if strings.Contains(srcOff, "Spicetify.RemoteConfigResolver") {
		t.Error("expFeatures embedded without cfg.ExperimentalFeatures")
	}
	// All on: bundle should contain all three.
	cfg := &config.Config{SidebarConfig: true, HomeConfig: true, ExpFeatures: true}
	srcOn := Bundle(BuildContext{Config: cfg})
	if !strings.Contains(srcOn, "SidebarConfig") {
		t.Error("sidebarConfig missing despite cfg.SidebarConfig=true")
	}
	if !strings.Contains(srcOn, "SpicetifyHomeConfig") {
		t.Error("homeConfig missing despite cfg.HomeConfig=true")
	}
	if !strings.Contains(srcOn, "Spicetify.expFeatureOverride") {
		t.Error("expFeatures missing despite cfg.ExpFeatures=true")
	}
}

// TestBundleIncludesExtensionAndApp verifies user extensions and custom apps
// produce non-empty body fragments inside the bundle.
func TestBundleIncludesExtensionAndApp(t *testing.T) {
	cfg := &config.Config{}
	ext := extensions.Loaded{Name: "shuffle", Source: "console.log('shuffle loaded');", IsModule: false}
	app := &customapps.App{
		Name: "Reddit", FolderName: "reddit",
		Manifest: &customapps.Manifest{Name: "Reddit", Icon: "<svg/>", ActiveIcon: "<svg/>"},
		IndexJS:  "function render(c){}",
	}
	src := Bundle(BuildContext{Config: cfg, Extensions: []extensions.Loaded{ext}, CustomApps: []*customapps.App{app}})
	if !strings.Contains(src, "shuffle loaded") {
		t.Error("user extension source not in bundle")
	}
	if !strings.Contains(src, "function render(c){}") {
		t.Error("custom app index.js not in bundle")
	}
	if !strings.Contains(src, `"reddit"`) {
		t.Error("custom app folder name missing from webpack chunk id")
	}
}
