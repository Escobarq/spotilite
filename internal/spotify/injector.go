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
