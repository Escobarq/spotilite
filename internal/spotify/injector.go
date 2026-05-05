// Package spotify provides utilities to customize the Spotify web experience
// inside the Wails webview, such as injecting CSS to hide unwanted UI elements.
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

// Injector handles navigation to Spotify and periodic CSS injection.
type Injector struct{}

// NewInjector creates a new Injector instance.
func NewInjector() *Injector {
	return &Injector{}
}

// Start navigates to Spotify and begins the background CSS injection loop.
// The loop respects context cancellation for graceful shutdown.
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
	// Build the JS payload dynamically so selectors and cssRules are inlined.
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
