package spicetify

import (
	"strings"

	"spotilite/internal/config"
	"spotilite/internal/customapps"
	"spotilite/internal/extensions"
	"spotilite/internal/themes"
)

// BuildContext bundles all the inputs the Spicetify injection needs at runtime
// so that Bundle(ctx) can produce the exact same payload spicetify-cli writes
// into Apps/xpui/index.html (see cli-main/src/apply/apply.go htmlMod):
//
//  1. window.Spicetify = {...}                  core.js
//  2. Spicetify.URI                            uri.js
//  3. Spicetify.CosmosAsync + localAPI         cosmos.js
//  4. Spicetify.Player                         player.js
//  5. Spicetify.Keyboard (Mousetrap alias)      keyboard.js
//  6. Spicetify components shim                 components.js
//  7. Spicetify.Config = {...}                  config.go (ConfigInjection)
//  8. Spicetify.localAPI = 'http://localhost:8765'
//  9. theme.js (if inject_theme_js)             wrapped in IIFE try/catch
// 10. helper: sidebarConfig (if cfg.SidebarConfig)
// 11. helper: homeConfig    (if cfg.HomeConfig)
// 12. helper: expFeatures    (if cfg.ExpFeatures)
// 13. user extensions (in order, each IIFE-wrapped)
// 14. custom apps (webpackChunk push + RouteHook per app)
//
// The accompanying CSS (themes.Theme.FullCSS) is returned separately via
// ctx.ThemeCSS() so the injector can ship it as a <style> tag in addition
// to the JS blob.
type BuildContext struct {
	Config     *config.Config
	Theme      *themes.Theme
	Extensions []extensions.Loaded
	CustomApps []*customapps.App
	LocalAPI   string // e.g. "http://localhost:8765"; "" disables assignment
}

// Bundle returns the concatenated JS source for all shipped spicetify
// sub-modules plus the runtime context (config, theme.js, extensions, custom
// apps). The returned string is the literal JS to evaluate; callers should
// wrap it to swallow exceptions (typically with try/catch + console.error).
//
// Order matters and mirrors spicetify-cli apply.go htmlMod. core.js must come
// first (it owns window.Spicetify).
func Bundle(ctx BuildContext) string {
	var b strings.Builder

	// 1-6: shims that build window.Spicetify.Player/URI/CosmosAsync/Keyboard.
	for _, src := range orderedShimSources() {
		b.WriteString(src)
		b.WriteString("\n\n")
	}

	// 7: Spicetify.Config = {...}
	b.WriteString(ConfigInjection(ctx.Config, ctx.Extensions, ctx.CustomApps))
	b.WriteString("\n")

	// 8: Spicetify.localAPI (bridge so extensions can call our Go API server)
	if ctx.LocalAPI != "" {
		b.WriteString("try{window.Spicetify=window.Spicetify||{};Spicetify.localAPI=")
		b.WriteString(jsStringSp(ctx.LocalAPI))
		b.WriteString(";}catch(e){console.error('[Spotilite] localAPI assign error:',e);}\n")
	}

	// 9: theme.js (optional extension-like script bundled with the theme)
	if ctx.Config != nil && ctx.Theme != nil && ctx.Theme.ThemeJS != "" && ctx.Config.InjectThemeJS {
		b.WriteString("/* === spotilite theme.js === */\n")
		b.WriteString(extensionIIFE("theme", ctx.Theme.ThemeJS))
	}

	// 10-12: helper scripts. These expect window.Spicetify to already exist,
	// which it does at this point.
	if ctx.Config != nil {
		if ctx.Config.SidebarConfig {
			b.WriteString("/* === spotilite helper: sidebarConfig === */\n")
			b.WriteString(extensionIIFE("sidebarConfig", sidebarConfigJS))
		}
		if ctx.Config.HomeConfig {
			b.WriteString("/* === spotilite helper: homeConfig === */\n")
			b.WriteString(extensionIIFE("homeConfig", homeConfigJS))
		}
		if ctx.Config.ExpFeatures {
			b.WriteString("/* === spotilite helper: expFeatures === */\n")
			b.WriteString(extensionIIFE("expFeatures", expFeaturesJS))
		}
	}

	// 13: user extensions.
	for _, ext := range ctx.Extensions {
		b.WriteString("/* === spotilite extension: ")
		b.WriteString(ext.Name)
		b.WriteString(" === */\n")
		b.WriteString(extensionIIFE(ext.Name, ext.Source))
	}

	// 14: custom apps (webpack chunk push + router hook).
	for _, app := range ctx.CustomApps {
		b.WriteString(WebpackChunkTemplate(app, ""))
		b.WriteString("\n")
		b.WriteString(RouteHook(app))
	}

	return b.String()
}

// ThemeCSS returns the active theme's full CSS (spice variables + user.css)
// or "" when there is no theme. The injector ships this as a <style> tag.
func (c BuildContext) ThemeCSS() string {
	if c.Theme == nil || c.Config == nil {
		return ""
	}
	if !c.Config.ReplaceColors && !c.Config.InjectCSS {
		return ""
	}
	return c.Theme.FullCSS(c.Config.ColorScheme)
}

// orderedShimSources returns the JS shim bodies in load order. Split out from
// the legacy orderedSources() in this file's prior avatar so that callers can
// still reach the unmodified shims from this single source of truth.
func orderedShimSources() []string {
	return []string{
		coreJS,
		uriJS,
		cosmosJS,
		playerJS,
		keyboardJS,
		componentsJS,
	}
}

// extensionIIFE wraps a JS body in an immediately invoked function expression
// with try/catch so one failing extension does not block the others. Matches
// the runtime behavior spicetify-cli achieves by loading each extension as a
// separate <script defer> tag.
func extensionIIFE(name, body string) string {
	var b strings.Builder
	b.WriteString("(function(){try{\n")
	b.WriteString(body)
	b.WriteString("\n}catch(e){console.error('[Spotilite] extension ")
	b.WriteString(name)
	b.WriteString(" error:',e);}})();\n")
	return b.String()
}
