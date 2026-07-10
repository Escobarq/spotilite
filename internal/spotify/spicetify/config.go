package spicetify

import (
	"encoding/json"
	"strings"

	"spotilite/internal/config"
	"spotilite/internal/customapps"
	"spotilite/internal/extensions"
)

// ConfigInjection returns the JS snippet that defines Spicetify.Config on the
// window. This mirrors what spicetify-cli writes into the patched xpui
// index.html (see cli-main/src/apply/apply.go htmlMod): a literal
// `Spicetify.Config = { ... }` assignment carrying the version, active
// theme/scheme, the enabled extensions and custom apps, and the Additional
// option flags the runtime JS consults.
//
// We emit the object via JSON marshalling (with proper escaping) and assign
// it. Extensions and custom_apps entries preserve the trailing "-" disabled
// marker convention used by spicetify-cli so that any code inspecting the raw
// list sees the same shape it would on the desktop client.
func ConfigInjection(cfg *config.Config, exts []extensions.Loaded, apps []*customapps.App) string {
	if cfg == nil {
		return "window.Spicetify = window.Spicetify || {};\nSpicetify.Config = {};\n"
	}

	type cfgObj struct {
		Version            string   `json:"version"`
		CurrentTheme       string   `json:"current_theme"`
		ColorScheme        string   `json:"color_scheme"`
		Extensions         []string `json:"extensions"`
		CustomApps         []string `json:"custom_apps"`
		CheckSpicetifyUpdate bool    `json:"check_spicetify_update"`
		InjectThemeJS      bool     `json:"inject_theme_js"`
		InjectCSS          bool     `json:"inject_css"`
		ReplaceColors      bool     `json:"replace_colors"`
		SidebarConfig      bool     `json:"sidebar_config"`
		HomeConfig         bool     `json:"home_config"`
		ExperimentalFeatures bool  `json:"experimental_features"`
	}

	extNames := make([]string, 0, len(exts))
	for _, e := range exts {
		n := e.Name + ".js"
		if e.IsModule {
			n = e.Name + ".mjs"
		}
		extNames = append(extNames, n)
	}
	appNames := make([]string, 0, len(apps))
	for _, a := range apps {
		appNames = append(appNames, a.Name)
	}

	obj := cfgObj{
		Version:              "spotilite-1.0",
		CurrentTheme:        cfg.CurrentTheme,
		ColorScheme:         cfg.ColorScheme,
		Extensions:          extNames,
		CustomApps:          appNames,
		CheckSpicetifyUpdate: cfg.CheckSpicetifyUpg,
		InjectThemeJS:       cfg.InjectThemeJS,
		InjectCSS:           cfg.InjectCSS,
		ReplaceColors:       cfg.ReplaceColors,
		SidebarConfig:       cfg.SidebarConfig,
		HomeConfig:          cfg.HomeConfig,
		ExperimentalFeatures: cfg.ExpFeatures,
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return "window.Spicetify = window.Spicetify || {};\nSpicetify.Config = {};\n"
	}

	var b strings.Builder
	b.WriteString("window.Spicetify = window.Spicetify || {};\n")
	b.WriteString("if(!Spicetify.Config){Spicetify.Config={};}\n")
	b.WriteString("try{Object.assign(Spicetify.Config, ")
	b.Write(data)
	b.WriteString(");}catch(e){console.error('[Spotilite] Config injection error:',e);}\n")
	return b.String()
}
