package spicetify

import "strings"

// Bundle returns the concatenated JS source for all shipped spicetify
// sub-modules, wrapped in an IIFE that installs them onto the window.
// Caller (typically injector) ships the result into the webview via
// wails runtime.WindowExecJS.
//
// Order matters: core.js must come first (it owns window.Spicetify),
// then uri (used by player), then player (consumes both), finally the
// optional React+ReactDOM bundle for components.
//
// The returned string is the literal JS to evaluate; callers should wrap
// it to swallow exceptions (typically with try/catch + console.error).
func Bundle() string {
	var b strings.Builder
	for _, src := range orderedSources() {
		b.WriteString(src)
		b.WriteString("\n\n")
	}
	return b.String()
}

func orderedSources() []string {
	return []string{
		coreJS,
		uriJS,
		cosmosJS,
		playerJS,
		keyboardJS,
		componentsJS,
	}
}
