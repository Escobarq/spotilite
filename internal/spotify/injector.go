// Package spotify provides utilities to customize the Spotify web experience
// inside the Wails webview, such as injecting CSS to hide unwanted UI elements
// and adding a functional title bar that communicates with the local API server.
package spotify

import (
	"context"
	"log/slog"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	// SpotifyURL is the target web application URL.
	SpotifyURL = "https://open.spotify.com"

	// APIBaseURL is the local API server URL for window controls.
	APIBaseURL = "http://localhost:8765"

	// initialDelay waits for the SPA to bootstrap before the first injection.
	initialDelay = 4 * time.Second

	// retryInterval re-applies styles in case React re-renders the DOM.
	retryInterval = 3 * time.Second

	// styleID is the DOM id for the injected <style> tag.
	styleID = "spotilite-css-inject"
)

// cssRules contains the stylesheet that hides specific Spotify UI elements.
const cssRules = `
/* Hide "Install desktop app" button */
a[href="/download"], a[href="/download/"],
[data-testid="download-button"], [data-testid="download-app-button"],
button[aria-label*="Install app"], button[aria-label*="Descargar app"] {
    display: none !important;
    visibility: hidden !important;
    opacity: 0 !important;
    pointer-events: none !important;
    height: 0 !important;
    width: 0 !important;
}

/* Hide "Connect to a device" button */
[data-testid="device-picker-icon-button"], [data-testid="connect-device-button"],
button[aria-label*="Connect to a device"], button[aria-label*="Conectar a un dispositivo"],
button[aria-label*="Dispositivos disponibles"] {
    display: none !important;
    visibility: hidden !important;
    opacity: 0 !important;
    pointer-events: none !important;
    height: 0 !important;
    width: 0 !important;
}

/* Spotilite title bar */
#spotilite-title-bar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 32px;
    z-index: 99999;
    background: #191414;
    display: flex;
    align-items: center;
    justify-content: space-between;
    user-select: none;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    box-shadow: 0 2px 8px rgba(0,0,0,0.4);
    --wails-draggable: drag;
}
#spotilite-title-bar .spotilite-left {
    display: flex;
    align-items: center;
    padding-left: 12px;
    gap: 8px;
    height: 100%;
    pointer-events: none;
}
#spotilite-title-bar .spotilite-logo {
    width: 16px;
    height: 16px;
}
#spotilite-title-bar .spotilite-text {
    color: #fff;
    font-size: 13px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    font-weight: 600;
}
#spotilite-title-bar .spotilite-menu-btn {
    height: 100%;
    padding: 0 12px;
    background: transparent;
    border: none;
    color: #b3b3b3;
    font-size: 13px;
    cursor: pointer;
    transition: color 0.15s;
}
#spotilite-title-bar .spotilite-menu-btn:hover {
    color: #fff;
}
#spotilite-title-bar .spotilite-dropdown {
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    background: #282828;
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
    z-index: 99999;
    min-width: 200px;
    padding: 4px 0;
    display: none;
}
#spotilite-title-bar .spotilite-dropdown.active {
    display: block;
}
#spotilite-title-bar .spotilite-dropdown-item {
    width: 100%;
    text-align: left;
    padding: 8px 12px;
    background: transparent;
    border: none;
    color: #fff;
    font-size: 13px;
    cursor: pointer;
    font-family: inherit;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
#spotilite-title-bar .spotilite-dropdown-item:hover {
    background: #3e3e3e;
}
#spotilite-title-bar .spotilite-dropdown-header {
    padding: 6px 12px;
    color: #b3b3b3;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    border-bottom: 1px solid rgba(255,255,255,0.1);
    margin-bottom: 4px;
}
#spotilite-title-bar .spotilite-dropdown-sep {
    height: 1px;
    background: rgba(255,255,255,0.1);
    margin: 4px 0;
}
#spotilite-title-bar .spotilite-win-btn {
    width: 46px;
    height: 100%;
    background: transparent;
    border: none;
    color: #b3b3b3;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}
#spotilite-title-bar .spotilite-win-btn:hover {
    background: #282828;
    color: #fff;
}
#spotilite-title-bar .spotilite-win-btn.close:hover {
    background: #e81123;
    color: #fff;
}

/* Push ONLY the root Spotify container below the title bar.
   The bar is injected as the first child of <body>, so the 
   very next sibling is Spotify's root app container. */
#spotilite-title-bar + div,
#spotilite-title-bar + #main,
#spotilite-title-bar + .Root {
    padding-top: 32px !important;
}
`

// titleBarScript creates the complete title bar with working buttons via fetch.
const titleBarScript = `
(function() {
    var bar = document.getElementById('spotilite-title-bar');
    if (bar) return;

    var API = '` + APIBaseURL + `';
    var lang = localStorage.getItem('spotilite.lang') || 'es';
    var bgMode = localStorage.getItem('spotilite.bg') !== 'false';

    var t = {
        es: { menu:'Menú', langHeader:'Idioma / Language', es:'Español', en:'English', bgMode:'Ejecutar en segundo plano', minimize:'Minimizar', maximize:'Maximizar', restore:'Restaurar', close:'Cerrar' },
        en: { menu:'Menu', langHeader:'Language / Idioma', es:'Spanish', en:'English', bgMode:'Run in background', minimize:'Minimize', maximize:'Maximize', restore:'Restore', close:'Close' }
    };
    var txt = t[lang] || t.es;

    function apiPost(path, body) {
        fetch(API + path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body || {}) }).catch(function(e){});
    }

    bar = document.createElement('div');
    bar.id = 'spotilite-title-bar';

    // Left: Logo + Title
    var left = document.createElement('div');
    left.className = 'spotilite-left';
    left.innerHTML = '<svg class="spotilite-logo" viewBox="0 0 168 168" fill="none"><circle cx="84" cy="84" r="80" fill="#1DB954"/><path d="M122.4 92.3c-15.1-8.9-40-9.7-54.4-5.4-2.2.7-4.5-.6-5.1-2.8-.7-2.2.6-4.5 2.8-5.1 16.5-4.9 43.6-4 60.3 6 2 1.2 2.6 3.7 1.4 5.7-1.1 1.9-3.6 2.6-5 1.6zm-.5 14.2c-1.2-2-3.7-2.6-5.7-1.4-13.3 8.2-33.6 10.5-49.4 5.8-2.2-.7-4.5.6-5.1 2.8-.7 2.2.6 4.5 2.8 5.1 17.9 5.4 40.1 2.8 55.1-6.3 1.9-1.3 2.5-3.8 1.3-5.6v-.4zm-5.2-42.6c-17.5-10.4-46.4-11.4-63-6.3-2.2.7-4.5-.6-5.1-2.8-.7-2.2.6-4.5 2.8-5.1 19.1-5.7 50.2-4.6 69.8 7.1 2 1.2 2.6 3.7 1.4 5.7-1.2 1.8-3.7 2.5-5.9 1.4z" fill="white"/></svg><span class="spotilite-text">Spotilite</span>';
    bar.appendChild(left);

    // Center: Menu
    var menuContainer = document.createElement('div');
    menuContainer.style.cssText = 'position:relative;height:100%;';

    var menuBtn = document.createElement('button');
    menuBtn.className = 'spotilite-menu-btn';
    menuBtn.innerHTML = '&#9776; ' + txt.menu;
    menuContainer.appendChild(menuBtn);

    var dropdown = document.createElement('div');
    dropdown.className = 'spotilite-dropdown';
    dropdown.innerHTML = '<div class="spotilite-dropdown-header">' + txt.langHeader + '</div>' +
        '<button class="spotilite-dropdown-item" data-lang="es">🇪🇸 ' + txt.es + '<span id="spotilite-check-es"></span></button>' +
        '<button class="spotilite-dropdown-item" data-lang="en">🇺🇸 ' + txt.en + '<span id="spotilite-check-en"></span></button>' +
        '<div class="spotilite-dropdown-sep"></div>' +
        '<button class="spotilite-dropdown-item" id="spotilite-bg-toggle"><span>' + txt.bgMode + '</span><span id="spotilite-bg-check" style="color:#1db954;font-weight:700;"></span></button>';
    menuContainer.appendChild(dropdown);

    menuBtn.onclick = function(e) {
        e.stopPropagation();
        dropdown.classList.toggle('active');
    };

    document.addEventListener('click', function(e) {
        if (!menuContainer.contains(e.target)) {
            dropdown.classList.remove('active');
        }
    });

    dropdown.querySelector('[data-lang="es"]').onclick = function() {
        lang = 'es'; localStorage.setItem('spotilite.lang', 'es'); apiPost('/api/settings/lang', {lang:'es'}); location.reload();
    };
    dropdown.querySelector('[data-lang="en"]').onclick = function() {
        lang = 'en'; localStorage.setItem('spotilite.lang', 'en'); apiPost('/api/settings/lang', {lang:'en'}); location.reload();
    };
    dropdown.querySelector('#spotilite-bg-toggle').onclick = function() {
        bgMode = !bgMode; localStorage.setItem('spotilite.bg', bgMode); apiPost('/api/settings/background', {enabled:bgMode}); updateChecks();
    };

    function updateChecks() {
        var es = document.getElementById('spotilite-check-es');
        var en = document.getElementById('spotilite-check-en');
        var bg = document.getElementById('spotilite-bg-check');
        if (es) es.textContent = lang === 'es' ? '✓' : '';
        if (en) en.textContent = lang === 'en' ? '✓' : '';
        if (bg) bg.textContent = bgMode ? '[x]' : '[ ]';
    }
    setTimeout(updateChecks, 100);

    bar.appendChild(menuContainer);

    // Right: Window controls
    var right = document.createElement('div');
    right.style.cssText = 'display:flex;align-items:center;height:100%;';

    function winBtn(label, action, isClose) {
        var btn = document.createElement('button');
        btn.className = 'spotilite-win-btn' + (isClose ? ' close' : '');
        btn.innerHTML = label;
        btn.onclick = function(e) { e.stopPropagation(); action(); };
        return btn;
    }

    var maximized = false;
    right.appendChild(winBtn('&#9472;', function() { apiPost('/api/window/minimize'); }));
    var maxBtn = winBtn('&#9633;', function() {
        maximized = !maximized;
        if (maximized) apiPost('/api/window/maximize');
        else apiPost('/api/window/unmaximize');
        maxBtn.innerHTML = maximized ? '&#9744;' : '&#9633;';
    });
    right.appendChild(maxBtn);
    right.appendChild(winBtn('&#10005;', function() { apiPost('/api/window/close'); }, true));

    bar.appendChild(right);

    if (document.body) {
        document.body.insertBefore(bar, document.body.firstChild);
    }
})();
`

// selectors is the list of DOM selectors we force-hide via inline styles.
var selectors = []string{
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

// Injector handles navigation to Spotify and periodic CSS injection.
type Injector struct{}

// NewInjector creates a new Injector instance.
func NewInjector() *Injector {
	return &Injector{}
}

// Start navigates to Spotify and begins the background CSS injection loop.
func (i *Injector) Start(ctx context.Context) {
	slog.Info("navigating to spotify", "url", SpotifyURL)
	runtime.WindowExecJS(ctx, `window.location.replace('`+SpotifyURL+`')`)

	go i.run(ctx)
}

// run waits for the initial page load and then ticks periodically to re-inject
// styles. It stops when the provided context is cancelled.
func (i *Injector) run(ctx context.Context) {
	select {
	case <-time.After(initialDelay):
		// proceed
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
			slog.Info("stopping spotify css injector")
			return
		}
	}
}

// inject builds and executes the JavaScript payload inside the webview.
func (i *Injector) inject(ctx context.Context) {
	script := `(function() {
    var styleId = '` + styleID + `';
    var style = document.getElementById(styleId);
    if (!style) {
        style = document.createElement('style');
        style.id = styleId;
        style.textContent = ` + "`" + cssRules + "`" + `;
        if (document.head) {
            document.head.appendChild(style);
        }
    }

` + titleBarScript + `

    var selectors = [
` + buildSelectorArray() + `
    ];

    selectors.forEach(function(sel) {
        try {
            var nodes = document.querySelectorAll(sel);
            nodes.forEach(function(node) {
                node.style.display = 'none';
                node.style.visibility = 'hidden';
                node.style.opacity = '0';
                node.style.pointerEvents = 'none';
            });
        } catch (e) {}
    });
})();`

	runtime.WindowExecJS(ctx, script)
}

// buildSelectorArray converts the Go selectors slice into a JS array literal.
func buildSelectorArray() string {
	out := ""
	for idx, sel := range selectors {
		if idx > 0 {
			out += ",\n"
		}
		out += `        '` + sel + `'`
	}
	return out
}
