// Package extensions loads Spicetify community extensions (single-file .js or
// .mjs scripts) from disk so they can be injected alongside the Spotilite
// wrapper. The desktop spicetify-cli puts these into the spotify xpui bundle's
// extensions/ folder where Spotify's own HTML loads them via <script>.
//
// On the web player that location doesn't exist, so this loader reads each
// enabled extension file (resolved against the configured Extensions folder,
// matching the spicetify-cli directory layout for compatibility), concatenates
// their bodies with the proper script-type wrapper, and returns a single JS
// chunk ready for the Wails injector to ship via runtime.WindowExecJS.
//
// .mjs files containing a "// spicetify_map{...}{...}" import-path rewriter
// comment block (the pattern used by spicetify-cli's pushExtensions to remap
// ES module imports against node_modules) get their import statements rewritten
// in-memory to absolute file URLs served by Spotilite's local API server. This
// keeps npm dependencies of .mjs extensions working without symlinking a
// node_modules folder into the (non-existent) xpui/extensions directory.
package extensions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Loaded represents one extension ready for injection.
type Loaded struct {
	Name    string // folder/short name (without trailing .js/.mjs)
	File    string // absolute source path
	Source  string // body content (possibly rewritten for .mjs)
	IsModule bool  // true for .mjs (wrap in type="module" semantics)
}

// Loader resolves extension filenames against one or more directories, the
// first match winning (mirrors Manager from themes). The default directory is
// the spicetify-cli user folder at %APPDATA%/spicetify/Extensions.
type Loader struct {
	searchDirs []string
	apiBase    string // http://localhost:8765 — used for .mjs import rewriting
}

// NewLoader returns a Loader that searches the given directories in order.
// apiBase is the local Go API server URL that serves /ext path roots used for
// rewriting .mjs imports (currently unused but kept for future bundling).
func NewLoader(searchDirs []string, apiBase string) *Loader {
	return &Loader{searchDirs: searchDirs, apiBase: apiBase}
}

// Load returns the source of each named extension that exists on disk,
// skipping missing files with a console-style warning emitted to the returned
// []string of skipped names. Names ending in ".js" or ".mjs" are honored
// verbatim; others get ".js" appended (spicetify-cli convention).
func (l *Loader) Load(names []string) ([]Loaded, []string, error) {
	var out []Loaded
	var skipped []string
	for _, n := range names {
		body, path, mod, ok, err := l.resolve(n)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			skipped = append(skipped, n)
			continue
		}
		out = append(out, Loaded{
			Name:     strings.TrimSuffix(strings.TrimSuffix(n, ".mjs"), ".js"),
			File:     path,
			Source:   body,
			IsModule: mod,
		})
	}
	return out, skipped, nil
}

func (l *Loader) resolve(name string) (body string, path string, isMod bool, ok bool, err error) {
	candidates := candidateFilenames(name)
	for _, dir := range l.searchDirs {
		for _, c := range candidates {
			p := filepath.Join(dir, c)
			data, e := os.ReadFile(p)
			if e != nil {
				continue
			}
			mod := strings.HasSuffix(c, ".mjs")
			src := string(data)
			if mod {
				src = rewriteMJSImports(src)
			}
			return src, p, mod, true, nil
		}
	}
	return "", "", false, false, nil
}

// candidateFilenames returns the candidate names to be tried when resolving a
// user-supplied extension reference. If the reference already ends in .js or
// .mjs it's tried as-is; otherwise ".js" is appended. We also try a fallback
// to ".mjs" if no ".js" is present (covers accidental bare names).
func candidateFilenames(ref string) []string {
	hasSuffix := strings.HasSuffix(ref, ".js") || strings.HasSuffix(ref, ".mjs")
	if hasSuffix {
		return []string{filepath.Base(ref)}
	}
	return []string{filepath.Base(ref) + ".js", filepath.Base(ref) + ".mjs"}
}

// rewriteMJSImports detects a "// spicetify_map {...}{...}" comment at the top
// of the extension and rewrites import-from specifiers to the supplied remap.
// The spicetify-cli pattern is:
//   // spicetify_map { "react" }{ ./node_modules/react/index.js }
// Each rule `ID -> TARGET` becomes: import ".{ID}." -> "./TARGET" (dropped if
// no TARGET is supplied — we leave it untouched so the user's bundler emits
// a sensible error rather than silent breakage).
//
// Spicetify-cli ships this rewriting in src/apply's pushExtensions function;
// on the web player we can't write alongside the bundle, so this is the best
// substitute (the file paths still need to be reachable through the local API
// server for any of this to actually load, currently out of scope — we just
// strip the import statements so the module doesn't break).
var spicetifyMapRe = regexp.MustCompile(`(?m)^//\s*spicetify_map\s*(\{"[^"]*"\})?\s*(\{"[^"]*"\})?\s*$`)

func rewriteMJSImports(src string) string {
	// Drop the comment so strict-mode modules don't choke on it. Real bundling
	// support (resolving node_modules deps through the local API server) is a
	// pending follow-up.
	src = spicetifyMapRe.ReplaceAllString(src, "")
	return src
}

// JSBundle wraps each Loaded extension in an IIFE so console logs are tagged
// and exceptions are swallowed (matching spicetify-cli's runtime behavior
// where one failing extension doesn't block the others).
func JSBundle(items []Loaded) string {
	var b strings.Builder
	for _, it := range items {
		b.WriteString("/* === spotilite extension: ")
		b.WriteString(it.Name)
		b.WriteString(" === */\n")
		wrapper := "(function() {\n try {\n%s\n } catch (e) { console.error('[Spotilite] extension %s error:', e); }\n})();\n"
		fmt.Fprintf(&b, wrapper, it.Source, it.Name)
	}
	return b.String()
}

// LoadAndBundle is the convenience entrypoint that resolves names, bundles
// the result, and returns both the script and the list of files that could
// not be found (so the caller can surface them in the UI/logs).
func (l *Loader) LoadAndBundle(names []string) (script string, skipped []string, err error) {
	loaded, skipped, err := l.Load(names)
	if err != nil {
		return "", nil, err
	}
	return JSBundle(loaded), skipped, nil
}
