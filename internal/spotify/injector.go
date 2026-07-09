package spotify

import (
	"context"
	"encoding/base64"
	"log/slog"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"spotilite/internal/spotify/modules"
)

const (
	SpotifyURL    = "https://open.spotify.com"
	APIBaseURL    = "http://localhost:8765"
	initialDelay  = 5 * time.Second
	retryInterval = 5 * time.Second
	styleID       = "spotilite-css-inject"
)

type Injector struct {
	modules  []modules.Module
	extraJS  string
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

func (i *Injector) SetExtraJS(js string) {
	i.extraJS = js
}

func (i *Injector) Start(ctx context.Context) {
	slog.Info("navigating to spotify", "url", SpotifyURL)
	runtime.WindowExecJS(ctx, "console.log('[Spotilite] Starting navigation');window.location.href='"+SpotifyURL+"';console.log('[Spotilite] Navigation triggered')")
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
	modulesScript := i.buildModulesScript()
	extra := i.extraJS

	slog.Info("[Spotilite] Injection called", "extraLen", len(extra), "modulesLen", len(modulesScript))

	// Build full code: modules + extra
	fullCode := modulesScript + extra

	// Use base64 + decodeURIComponent(escape(atob())) to handle UTF-8 correctly
	b64 := b64Encode(fullCode)
	evalScript := "try{var s=atob('" + b64 + "');var decoded=decodeURIComponent(escape(s));eval(decoded);console.log('[Spotilite] Injected successfully');}catch(e){console.error('[Spotilite] Injection error:',e.message||String(e));};"
	runtime.WindowExecJS(ctx, evalScript)
	slog.Info("[Spotilite] Injection sent")
}

func b64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
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
			jsBuilder.WriteString("try { ")
			jsBuilder.WriteString(js)
			jsBuilder.WriteString(" } catch(e) { console.error('[Spotilite] Module error:', e); } ")
		}
		if sels := m.Selectors(); len(sels) > 0 {
			selectors = append(selectors, sels...)
		}
	}

	css := cssBuilder.String()
	js := jsBuilder.String()
	selArray := buildSelectorArray(selectors)

	// Build hiding selectors script
	hideScript := "var selectors=[" + selArray + "];selectors.forEach(function(sel){try{document.querySelectorAll(sel).forEach(function(node){node.style.display='none';node.style.visibility='hidden';node.style.opacity='0';node.style.pointerEvents='none';});}catch(e){}});"

	// Build final script that injects CSS and runs JS - escape the css string for JS
	return "try{var styleId='" + styleID + "-modules';var style=document.getElementById(styleId);if(!style){style=document.createElement('style');style.id=styleId;style.textContent=" + jsString(css) + ";document.head.appendChild(style);}}catch(e){console.error('[Spotilite] Module CSS error:',e);}" + js + hideScript
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
	`[data-testid="install-app-button"]`,
	`button[aria-label*="Install app"]`,
	`button[aria-label*="Descargar app"]`,
	`button[aria-label*="Install application"]`,
	`[data-testid="device-picker-icon-button"]`,
	`[data-testid="connect-device-button"]`,
	`button[aria-label*="Connect to a device"]`,
	`button[aria-label*="Conectar a un dispositivo"]`,
	`button[aria-label*="Dispositivos disponibles"]`,
	`[data-testid="upgrade-button"]`,
	`[data-testid="premium-link"]`,
	`[data-testid="upgrade-to-premium"]`,
	`[data-testid="upsell-banner"]`,
	`button[aria-label*="Upgrade to Premium"]`,
	`button[aria-label*="Get Premium"]`,
	`button[aria-label*="Actualizar a Premium"]`,
	`button[aria-label*="Obtén Premium"]`,
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
