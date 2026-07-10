// Package customapps loads Spicetify custom apps (React-based side panels)
// from disk and returns the JS needed to mount them into the web player.
//
// On the desktop client, custom apps are compiled into webpack chunks that
// Spotify's router navigates to. On the web player we cannot mutate the
// webpack manifest, so this package instead:
//
//   - reads the app's manifest.json + index.js + optional subfiles
//   - injects a DOM container + sidebar entry
//   - uses history.pushState + popstate to toggle the app when the user
//     clicks the sidebar icon
//   - calls the app's render(container) function (the required lifecycle hook)
//
// This is a simplified implementation that supports the most common custom
// apps (lyrics-plus, reddit, etc.) as long as they don't rely on webpack
// chunk IDs or deep Spotify internals.
package customapps

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manifest mirrors the custom app manifest.json shape used by spicetify-cli.
type Manifest struct {
	Name                  interface{} `json:"name"` // string or map[string]string (translations)
	Icon                  string      `json:"icon"`
	ActiveIcon            string      `json:"active-icon"`
	Subfiles              []string    `json:"subfiles"`
	SubfilesExtension     []string    `json:"subfiles_extension"`
	AssetsDir             string      `json:"assets"`
	EnableExperimentalAPI bool        `json:"enable_experimental_api"`
}

func (m *Manifest) GetName() string {
	if s, ok := m.Name.(string); ok {
		return s
	}
	if mmap, ok := m.Name.(map[string]interface{}); ok {
		if en, ok := mmap["en"].(string); ok {
			return en
		}
		for _, v := range mmap {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// App represents a loaded custom app ready for injection.
type App struct {
	Name       string      // display name from manifest (translation-aware)
	FolderName string      // directory name on disk (used as route path and webpack chunk id)
	Dir        string      // absolute path
	Manifest   *Manifest   // parsed manifest.json
	IndexJS    string      // contents of index.js (must define render(container))
	Subfiles   []string    // concatenated JS of subfiles (optional)
	StylesCSS  string      // contents of style.css if present
	JSBundle   string      // full script: subfiles + index.js
	HasRender  bool        // true if index.js defines a global render function
}

// Manager resolves custom apps from one or more search directories.
type Manager struct {
	searchDirs []string
	cache      map[string]*App
}

// NewManager builds a Manager that searches the given directories in order.
func NewManager(searchDirs []string) *Manager {
	return &Manager{searchDirs: searchDirs, cache: make(map[string]*App)}
}

// Load returns a custom app by folder name. Errors if the folder or its
// manifest.json/index.js are missing or unparseable. Successful loads are cached.
func (m *Manager) Load(name string) (*App, error) {
	if name == "" {
		return nil, fmt.Errorf("customapps: empty name")
	}
	if app, ok := m.cache[name]; ok {
		return app, nil
	}
	for _, dir := range m.searchDirs {
		app, err := loadFromDir(filepath.Join(dir, name))
		if err == nil {
			m.cache[name] = app
			return app, nil
		}
	}
	return nil, fmt.Errorf("customapps: %q not found in %v", name, m.searchDirs)
}

// ListNames returns the set of custom app folder names available across all
// search directories (deduplicated, first-seen order).
func (m *Manager) ListNames() []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, dir := range m.searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			n := e.Name()
			if seen[n] {
				continue
			}
			// Quick presence ping: must have manifest.json + index.js
			if _, err := os.Stat(filepath.Join(dir, n, "manifest.json")); err != nil {
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, n, "index.js")); err != nil {
				continue
			}
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

func loadFromDir(dir string) (*App, error) {
	manifestPath := filepath.Join(dir, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var mf Manifest
	if err := json.Unmarshal(manifestData, &mf); err != nil {
		return nil, fmt.Errorf("customapps: parse manifest.json %s: %w", dir, err)
	}

	indexPath := filepath.Join(dir, "index.js")
	indexJS, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	var subfiles []string
	for _, sf := range mf.Subfiles {
		p := filepath.Join(dir, sf)
		data, e := os.ReadFile(p)
		if e != nil {
			continue
		}
		subfiles = append(subfiles, string(data))
	}
	for _, ext := range mf.SubfilesExtension {
		p := filepath.Join(dir, ext)
		data, e := os.ReadFile(p)
		if e != nil {
			continue
		}
		subfiles = append(subfiles, string(data))
	}

	var stylesCSS string
	if css, e := os.ReadFile(filepath.Join(dir, "style.css")); e == nil {
		stylesCSS = string(css)
	}

	hasRender := strings.Contains(string(indexJS), "function render(") ||
		strings.Contains(string(indexJS), "const render =") ||
		strings.Contains(string(indexJS), "var render =") ||
		strings.Contains(string(indexJS), "export function render")

	app := &App{
		Name:       mf.GetName(),
		FolderName: filepath.Base(dir),
		Dir:        dir,
		Manifest:   &mf,
		IndexJS:     string(indexJS),
		Subfiles:   subfiles,
		StylesCSS:  stylesCSS,
		HasRender:  hasRender,
	}
	app.JSBundle = strings.Join(subfiles, "\n") + "\n" + app.IndexJS
	return app, nil
}

// JSInjection returns the JS script that, when injected into the web player,
// mounts the custom app into a DOM container and registers a sidebar entry.
// The script:
//   - creates a <div id="spotilite-app-<name>"> container
//   - injects a <li> into the sidebar with the app's icon (from manifest)
//   - listens for popstate events and toggles the container when the route
//     matches '/spotilite/<name>'
//   - calls render(container) if defined (the required lifecycle hook)
func (a *App) JSInjection() string {
	name := a.Name
	icon := a.Manifest.Icon
	activeIcon := a.Manifest.ActiveIcon
	if activeIcon == "" {
		activeIcon = icon
	}
	js := fmt.Sprintf("(function(){var id='spotilite-app-%s';var path='/spotilite/%s';var active=false;var icon=%s;var activeIcon=%s;function log(m){console.log('[Spotilite CustomApp %s] '+m);}function inject(){log('Attempting...');var s=['li.main-navBar-navBarItem','.main-navBar-navBarItem','[class*=main-navBar-navBarItem]','.main-navBar-navBar','.main-appShell-sideBar li','.main-appShell-sideBar ul','.main-appShell-sideBar'];var n=null;for(var i=0;i<s.length;i++){var e=document.querySelector(s[i]);if(e){n=e;log('Found:'+s[i]);break;}}if(!n){log('Not found');setTimeout(inject,2000);return;}var p=n.parentElement||n;var li=document.createElement('li');li.className='main-navBar-navBarItem spotilite-app-link';li.dataset.name='%s';li.innerHTML=icon;li.title='%s';li.style.cssText='cursor:pointer;display:flex;align-items:center;justify-content:center;';li.onclick=function(e){e.preventDefault();e.stopPropagation();active=!active;if(active){history.pushState({app:'%s'},'',path);li.innerHTML=activeIcon;li.classList.add('main-navBar-navBarLinkActive');}else{history.back();li.innerHTML=icon;li.classList.remove('main-navBar-navBarLinkActive');}};p.appendChild(li);log('Injected');}function handlePopState(e){var c=document.getElementById(id);if(!c)return;var st=(e&&e.state)||{};if(st.app==='%s'){c.style.display='block';active=true;}else{c.style.display='none';active=false;}}window.addEventListener('popstate',handlePopState);setTimeout(inject,2000);})();",
		name, name, jsEscape(icon), jsEscape(activeIcon), name,
		name, name, name, name)
	return js
}

func jsEscape(s string) string {
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

// CSSInjection returns the theme CSS for the app (style.css) if present.
func (a *App) CSSInjection() string {
	return a.StylesCSS
}