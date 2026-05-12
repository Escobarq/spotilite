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
	modulesScript := i.buildModulesScript()
	settingsBtnScript := i.buildSettingsBtnScript()

	script := `(function() {
		function doInject() {
			` + settingsBtnScript + `
			` + modulesScript + `
		}

		var btn = document.getElementById('spotilite-settings-float');
		if (!btn) {
			try {
				doInject();
				console.log('[Spotilite] Injected successfully');
			} catch(e) {
				console.error('[Spotilite] Injection error:', e);
			}
		}
	})();`

	runtime.WindowExecJS(ctx, script)
}

func (i *Injector) buildSettingsBtnScript() string {
	return `(function() {
	if (window.__spotilite_settings_injected) return;
	window.__spotilite_settings_injected = true;

	var btn = document.getElementById('spotilite-settings-float');
	if (btn) return;

	var style = document.createElement('style');
	style.id = 'spotilite-float-style';
	style.textContent = '#spotilite-settings-float{position:fixed;bottom:80px;right:20px;width:44px;height:44px;border-radius:50%;background:#1DB954;border:none;color:#000;cursor:pointer;z-index:2147483647;display:flex;align-items:center;justify-content:center;box-shadow:0 4px 16px rgba(0,0,0,0.5);transition:transform 0.2s}#spotilite-settings-float:hover{transform:scale(1.1);background:#1ed760}#spotilite-settings-panel{position:fixed;bottom:140px;right:20px;width:280px;background:#282828;border-radius:12px;box-shadow:0 8px 32px rgba(0,0,0,0.7);z-index:2147483647;flex-direction:column;padding:16px;gap:12px;display:none;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif}#spotilite-settings-panel h3{margin:0;color:#1DB954;font-size:14px;font-weight:700;text-transform:uppercase}.spotilite-setting-row{display:flex;align-items:center;justify-content:space-between;padding:8px 0}.spotilite-setting-label{color:#fff;font-size:13px}.spotilite-toggle{width:36px;height:20px;border-radius:10px;background:#555;position:relative;cursor:pointer;transition:background 0.2s;border:none;padding:0}.spotilite-toggle::after{content:"";position:absolute;top:2px;left:2px;width:16px;height:16px;border-radius:50%;background:#fff;transition:transform 0.2s}.spotilite-toggle.on{background:#1DB954}.spotilite-toggle.on::after{transform:translateX(16px)}.spotilite-divider{height:1px;background:rgba(255,255,255,0.1);margin:4px 0}.spotilite-lang-btns{display:flex;gap:8px}.spotilite-lang-btn{flex:1;padding:6px;background:#3e3e3e;border:none;color:#fff;border-radius:6px;cursor:pointer;font-size:12px;transition:background 0.2s}.spotilite-lang-btn.active{background:#1DB954;color:#000;font-weight:600}.spotilite-close-btn{position:absolute;top:8px;right:8px;width:24px;height:24px;border-radius:50%;background:transparent;border:none;color:#888;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:16px}.spotilite-close-btn:hover{background:rgba(255,255,255,0.1);color:#fff}';
	document.head.appendChild(style);

	btn = document.createElement('button');
	btn.id = 'spotilite-settings-float';
	btn.title = 'Spotilite Settings';
	btn.innerHTML = '<svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M19.14 12.94c.04-.31.06-.63.06-.94 0-.31-.02-.63-.06-.94l2.03-1.58a.49.49 0 00.12-.61l-1.92-3.32a.49.49 0 00-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.48.48 0 00-.48-.41h-3.84a.48.48 0 00-.48.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96a.49.49 0 00-.59.22L2.74 8.87a.49.49 0 00.12.61l2.03 1.58c-.05.31-.07.63-.07.94s.02.64.07.94l-2.03 1.58a.49.49 0 00-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.48-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32a.49.49 0 00-.12-.61l-2.01-1.58zM12 15.6A3.6 3.6 0 1115.6 12 3.6 3.6 0 0112 15.6z"/></svg>';
	btn.onclick = function() {
		var panel = document.getElementById('spotilite-settings-panel');
		if (panel) {
			panel.style.display = panel.style.display === 'none' ? 'flex' : 'none';
		}
	};
	document.body.appendChild(btn);

	var panel = document.createElement('div');
	panel.id = 'spotilite-settings-panel';
	panel.innerHTML = '<button class="spotilite-close-btn" onclick="this.parentElement.style.display=\'none\'">&times;</button><h3>Spotilite</h3><div class="spotilite-setting-row"><span class="spotilite-setting-label">Ad Blocker</span><button class="spotilite-toggle" id="spotilite-toggle-ad" onclick="window.__spotiliteToggle(\'adblock\', this)"></button></div><div class="spotilite-setting-row"><span class="spotilite-setting-label">Block Sections</span><button class="spotilite-toggle" id="spotilite-toggle-sec" onclick="window.__spotiliteToggle(\'sectionblock\', this)"></button></div><div class="spotilite-setting-row"><span class="spotilite-setting-label">Hide Premium</span><button class="spotilite-toggle" id="spotilite-toggle-prem" onclick="window.__spotiliteToggle(\'premium_spoof\', this)"></button></div><div class="spotilite-setting-row"><span class="spotilite-setting-label">Experiments</span><button class="spotilite-toggle" id="spotilite-toggle-exp" onclick="window.__spotiliteToggle(\'experiments\', this)"></button></div><div class="spotilite-setting-row"><span class="spotilite-setting-label">History</span><button class="spotilite-toggle" id="spotilite-toggle-hist" onclick="window.__spotiliteToggle(\'history\', this)"></button></div><div class="spotilite-divider"></div><div class="spotilite-setting-row"><span class="spotilite-setting-label">Language</span></div><div class="spotilite-lang-btns"><button class="spotilite-lang-btn" id="spotilite-lang-es" onclick="window.__spotiliteSetLang(\'es\')">ES</button><button class="spotilite-lang-btn" id="spotilite-lang-en" onclick="window.__spotiliteSetLang(\'en\')">EN</button></div>';
	document.body.appendChild(panel);

	var API = 'http://localhost:8765';
	var lang = localStorage.getItem('spotilite.lang') || 'es';

	window.__spotiliteToggle = function(module, el) {
		var currentState = el.classList.contains('on');
		var newState = !currentState;
		el.classList.toggle('on');
		localStorage.setItem('spotilite.' + module, newState);
		fetch(API + '/api/spotx/module', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({module: module, enabled: newState})
		}).catch(function(){});
	};

	window.__spotiliteSetLang = function(l) {
		lang = l;
		localStorage.setItem('spotilite.lang', l);
		document.getElementById('spotilite-lang-es').classList.toggle('active', l === 'es');
		document.getElementById('spotilite-lang-en').classList.toggle('active', l === 'en');
		fetch(API + '/api/settings/lang', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({lang: l})
		}).catch(function(){});
		setTimeout(function(){ location.reload(); }, 300);
	};

	function loadSettings() {
		var modules = ['adblock', 'sectionblock', 'premium_spoof', 'experiments', 'history'];
		var ids = ['spotilite-toggle-ad', 'spotilite-toggle-sec', 'spotilite-toggle-prem', 'spotilite-toggle-exp', 'spotilite-toggle-hist'];
		for (var i = 0; i < modules.length; i++) {
			var enabled = localStorage.getItem('spotilite.' + modules[i]) !== 'false';
			var el = document.getElementById(ids[i]);
			if (el) {
				el.classList.toggle('on', enabled);
			}
		}
		document.getElementById('spotilite-lang-es').classList.toggle('active', lang === 'es');
		document.getElementById('spotilite-lang-en').classList.toggle('active', lang === 'en');
	}

	loadSettings();
	console.log('[Spotilite] Settings button injected');
	})();`
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
	`a[href*="/premium"]`,
	`a[href*="/subscription"]`,
	`button[aria-label*="Upgrade to Premium"]`,
	`button[aria-label*="Get Premium"]`,
	`button[aria-label*="Actualizar a Premium"]`,
	`button[aria-label*="Obtén Premium"]`,
	`[data-encore-id="buttonTertiary"]`,
	`[data-encore-id="buttonSecondary"]`,
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
