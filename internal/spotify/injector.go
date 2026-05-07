package spotify

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/spotify/modules"
)

const (
	SpotifyURL    = "https://open.spotify.com"
	APIBaseURL    = "http://localhost:8765"
	initialDelay  = 2 * time.Second
	retryInterval = 3 * time.Second
	styleID       = "spotilite-css-inject"
)

type Injector struct {
	modules []modules.Module
}

func NewInjector(modList ...modules.Module) *Injector {
	if len(modList) == 0 {
		modList = DefaultModules()
	}
	return &Injector{modules: modList}
}

func DefaultModules() []modules.Module {
	return []modules.Module{
		modules.NewAdBlockModule(true),
		modules.NewSectionBlockModule(true),
		modules.NewPremiumSpoofModule(true),
		modules.NewExperimentModule(true),
		modules.NewHistoryModule(true),
		modules.NewCustomCSSModule(true),
	}
}

func (i *Injector) GetModules() []modules.Module {
	return i.modules
}

func (i *Injector) GetModule(name string) modules.Module {
	for _, m := range i.modules {
		if m.Name() == name {
			return m
		}
	}
	return nil
}

func (i *Injector) Start(ctx context.Context) {
	slog.Info("navigating to spotify", "url", SpotifyURL)
	runtime.WindowExecJS(ctx, `window.location.replace('`+SpotifyURL+`')`)

	go i.run(ctx)
}

func (i *Injector) run(ctx context.Context) {
	select {
	case <-time.After(initialDelay):
	case <-ctx.Done():
		slog.Info("spotify injector stopped before initial delay")
		return
	}

	i.inject(ctx)

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			i.inject(ctx)
		case <-ctx.Done():
			slog.Info("stopping spotify injector")
			return
		}
	}
}

func (i *Injector) inject(ctx context.Context) {
	titleScript := i.buildTitleBarScript()
	modulesScript := i.buildModulesScript()

	bootstrap := `(function() {
		function doInject() {
			` + titleScript + `
			` + modulesScript + `
		}
		if (document.readyState === 'loading') {
			document.addEventListener('DOMContentLoaded', doInject);
		} else {
			doInject();
		}
	})();`

	runtime.WindowExecJS(ctx, bootstrap)
}

func (i *Injector) buildTitleBarScript() string {
	return `(function() {` +
		`try {` +
		`var styleId = '` + styleID + `-titlebar';` +
		`var style = document.getElementById(styleId);` +
		`if (!style) {` +
		`style = document.createElement('style');` +
		`style.id = styleId;` +
		`style.textContent = ` + jsString(titleBarCSS) + `;` +
		`document.head.appendChild(style);` +
		`}` +
		`} catch(e) { console.error('[Spotilite] CSS injection error:', e); }` +
		`try {` +
		titleBarScript +
		`} catch(e) { console.error('[Spotilite] Titlebar error:', e); }` +
		`})();`
}

func (i *Injector) buildModulesScript() string {
	var cssBuilder strings.Builder
	var jsBuilder strings.Builder
	var selectors []string

	for _, m := range i.modules {
		if !m.Enabled() {
			continue
		}
		if css := m.CSS(); css != "" {
			cssBuilder.WriteString(css)
			cssBuilder.WriteString("\n")
		}
		if js := m.JS(); js != "" {
			jsBuilder.WriteString("(function() { try { ")
			jsBuilder.WriteString(js)
			jsBuilder.WriteString(" } catch(e) { console.error('[Spotilite] Module error:', e); } })();\n")
		}
		if sels := m.Selectors(); len(sels) > 0 {
			selectors = append(selectors, sels...)
		}
	}

	js := jsBuilder.String()
	selArray := buildSelectorArray(selectors)

	return `(function() {` +
		`try {` +
		`var styleId = '` + styleID + `-modules';` +
		`var style = document.getElementById(styleId);` +
		`if (!style) {` +
		`style = document.createElement('style');` +
		`style.id = styleId;` +
		`style.textContent = ` + jsString(cssBuilder.String()) + `;` +
		`document.head.appendChild(style);` +
		`}` +
		`} catch(e) { console.error('[Spotilite] Module CSS error:', e); }` +
		`try {` +
		js +
		`var selectors = [` + selArray + `];` +
		`selectors.forEach(function(sel) {` +
		`try {` +
		`var nodes = document.querySelectorAll(sel);` +
		`nodes.forEach(function(node) {` +
		`node.style.display = 'none';` +
		`node.style.visibility = 'hidden';` +
		`node.style.opacity = '0';` +
		`node.style.pointerEvents = 'none';` +
		`});` +
		`} catch (e) {}` +
		`});` +
		`} catch(e) { console.error('[Spotilite] Module JS error:', e); }` +
		`})();`
}

func jsString(s string) string {
	var b strings.Builder
	b.WriteRune('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteRune('"')
	return b.String()
}

func (i *Injector) UpdateCustomCSS(ctx context.Context) {
	js := `(function() {
	var css = localStorage.getItem('spotilite.custom_css');
	if (!css) { document.getElementById('spotilite-custom-css')?.remove(); return; }
	var id = 'spotilite-custom-css';
	var style = document.getElementById(id);
	if (!style) {
		style = document.createElement('style');
		style.id = id;
		document.head.appendChild(style);
	}
	style.textContent = css;
})();`
	runtime.WindowExecJS(ctx, js)
}

const titleBarCSS = `
html, body { margin: 0 !important; padding: 0 !important; overflow: hidden !important; }
#spotilite-title-bar {
  position: fixed; top: 0; left: 0; width: 100%; height: 28px;
  z-index: 2147483647; background: #191414;
  display: flex; align-items: center; justify-content: space-between;
  --wails-draggable: drag;
  user-select: none;
}
#spotilite-title-bar .spotilite-left {
  display: flex; align-items: center; gap: 8px;
  padding-left: 12px; height: 100%;
  --wails-draggable: drag;
}
#spotilite-title-bar .spotilite-logo {
  width: 18px; height: 18px; background: #1DB954;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: bold; color: #000;
}
#spotilite-title-bar .spotilite-title {
  color: #fff; font-size: 12px; font-weight: 600;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}
#spotilite-title-bar .spotilite-right { display: flex; height: 100%; }
#spotilite-title-bar .spotilite-btn {
  width: 46px; height: 100%; background: transparent; border: none;
  color: #b3b3b3; font-size: 14px; cursor: pointer;
  transition: background 0.15s, color 0.15s;
}
#spotilite-title-bar .spotilite-btn:hover { background: #282828; color: #fff; }
#spotilite-title-bar .spotilite-close:hover { background: #e81123; color: #fff; }
`

var baseSelectors = []string{
	`a[href="/download"]`,
	`a[href="/download/"]`,
	`[data-testid="download-button"]`,
	`[data-testid="download-app-button"]`,
	`button[aria-label*="Install app"]`,
	`button[aria-label*="Descargar app"]`,
	`[data-testid="device-picker-icon-button"]`,
	`[data-testid="connect-device-button"]`,
	`button[aria-label*="Connect to a device"]`,
	`button[aria-label*="Conectar a un dispositivo"]`,
	`button[aria-label*="Dispositivos disponibles"]`,
}

func buildSelectorArray(selectors []string) string {
	all := append(baseSelectors, selectors...)
	out := ""
	for idx, sel := range all {
		if idx > 0 {
			out += ",\n"
		}
		out += ` '` + sel + `'`
	}
	return out
}
