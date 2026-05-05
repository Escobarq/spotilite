// Package spotify provides utilities to customize the Spotify web experience
// inside the Wails webview, such as injecting CSS to hide unwanted UI elements
// and adding a draggable region for frameless windows.
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
a[href="/download"],
a[href="/download/"],
[data-testid="download-button"],
[data-testid="download-app-button"],
button[aria-label*="Install app"],
button[aria-label*="Descargar app"] {
    display: none !important;
    visibility: hidden !important;
    opacity: 0 !important;
    pointer-events: none !important;
    height: 0 !important;
    width: 0 !important;
}

/* Hide "Connect to a device" button */
[data-testid="device-picker-icon-button"],
[data-testid="connect-device-button"],
button[aria-label*="Connect to a device"],
button[aria-label*="Conectar a un dispositivo"],
button[aria-label*="Dispositivos disponibles"] {
    display: none !important;
    visibility: hidden !important;
    opacity: 0 !important;
    pointer-events: none !important;
    height: 0 !important;
    width: 0 !important;
}

/* Make room for the draggable title bar */
#spotilite-drag-bar {
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
    --wails-draggable: drag;
}

/* Push Spotify content below the drag bar */
#spotilite-drag-bar ~ div[class*="main-view-container"],
#spotilite-drag-bar ~ .main-view-container__scroll-node,
#spotilite-drag-bar ~ [data-testid*="topbar"] {
    margin-top: 32px !important;
}

/* Ensure Spotify top bar doesn't overlap */
[class*="Root__top-bar"],
[data-testid*="topbar"] {
    margin-top: 32px !important;
}
`

// selectors is the list of DOM selectors we force-hide via inline styles as a
// fallback when the stylesheet injection hasn't caught up yet.
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

// dragBarScript creates a minimal draggable title bar inside the Spotify page.
// This bar has NO buttons (they can't call Go bindings from an external domain).
// Window controls are handled via the system tray menu and global shortcuts.
const dragBarScript = `
    var bar = document.getElementById('spotilite-drag-bar');
    if (!bar) {
        bar = document.createElement('div');
        bar.id = 'spotilite-drag-bar';

        var left = document.createElement('div');
        left.style.cssText = 'display:flex;align-items:center;padding-left:12px;gap:8px;height:100%;pointer-events:none;';
        left.innerHTML = '<svg width="16" height="16" viewBox="0 0 168 168" fill="none"><circle cx="84" cy="84" r="80" fill="#1DB954"/><path d="M122.4 92.3c-15.1-8.9-40-9.7-54.4-5.4-2.2.7-4.5-.6-5.1-2.8-.7-2.2.6-4.5 2.8-5.1 16.5-4.9 43.6-4 60.3 6 2 1.2 2.6 3.7 1.4 5.7-1.1 1.9-3.6 2.6-5 1.6zm-.5 14.2c-1.2-2-3.7-2.6-5.7-1.4-13.3 8.2-33.6 10.5-49.4 5.8-2.2-.7-4.5.6-5.1 2.8-.7 2.2.6 4.5 2.8 5.1 17.9 5.4 40.1 2.8 55.1-6.3 1.9-1.3 2.5-3.8 1.3-5.6v-.4zm-5.2-42.6c-17.5-10.4-46.4-11.4-63-6.3-2.2.7-4.5-.6-5.1-2.8-.7-2.2.6-4.5 2.8-5.1 19.1-5.7 50.2-4.6 69.8 7.1 2 1.2 2.6 3.7 1.4 5.7-1.2 1.8-3.7 2.5-5.9 1.4z" fill="white"/></svg><span style="color:#fff;font-size:13px;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif;font-weight:600;">Spotilite</span>';

        bar.appendChild(left);

        if (document.body) {
            document.body.insertBefore(bar, document.body.firstChild);
        }
    }
`

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

` + dragBarScript + `

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
