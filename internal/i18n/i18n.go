// Package i18n provides simple internationalization support for spotilite.
// It supports runtime language switching and fallback to English.
package i18n

import (
	"fmt"
	"sync"
)

// Supported language codes.
const (
	LangEnglish = "en"
	LangSpanish = "es"
)

// DefaultLang is used when no language is explicitly set.
const DefaultLang = LangEnglish

// Translator holds translations and the current language.
type Translator struct {
	mu   sync.RWMutex
	lang string
	data map[string]map[string]string
}

// New creates a Translator preloaded with default translations.
func New() *Translator {
	t := &Translator{
		lang: DefaultLang,
		data: make(map[string]map[string]string),
	}
	t.loadDefaults()
	return t
}

// SetLanguage changes the active language.
func (t *Translator) SetLanguage(lang string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lang = lang
}

// Language returns the currently active language code.
func (t *Translator) Language() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.lang
}

// T retrieves the translation for a given key in the active language.
// If the key or language is missing, it falls back to English, then returns the key itself.
func (t *Translator) T(key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if msgs, ok := t.data[t.lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to English.
	if msgs, ok := t.data[LangEnglish]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}

// Tf is a convenience wrapper around T that formats the string with args.
func (t *Translator) Tf(key string, args ...any) string {
	return fmt.Sprintf(t.T(key), args...)
}

// loadDefaults registers the built-in translations.
func (t *Translator) loadDefaults() {
	t.data[LangEnglish] = map[string]string{
		"app.title": "Spotilite",
		"tray.show": "Show Window",
		"tray.hide": "Hide Window",
		"tray.quit": "Quit",
		"tray.runInBackground": "Run in Background",
		"tray.language": "Language",
		"tray.lang.en": "English",
		"tray.lang.es": "Spanish",
		"notif.minimizedToTray": "Spotilite is running in the background",
		"notif.minimizedBody": "Click the tray icon to show the window again.",
		"dialog.confirmQuit": "Are you sure you want to quit Spotilite?",
		"spotx.adblock": "Ad Blocker",
		"spotx.sections": "Block Sections",
		"spotx.premium": "Hide Premium",
		"spotx.experiments": "Disable Experiments",
		"spotx.lyrics_theme": "Lyrics Theme",
		"spotx.history": "Track History",
	}

	t.data[LangSpanish] = map[string]string{
		"app.title": "Spotilite",
		"tray.show": "Mostrar Ventana",
		"tray.hide": "Ocultar Ventana",
		"tray.quit": "Salir",
		"tray.runInBackground": "Ejecutar en segundo plano",
		"tray.language": "Idioma",
		"tray.lang.en": "Inglés",
		"tray.lang.es": "Español",
		"notif.minimizedToTray": "Spotilite se está ejecutando en segundo plano",
		"notif.minimizedBody": "Haz clic en el icono de la bandeja para volver a mostrar la ventana.",
		"dialog.confirmQuit": "¿Estás seguro de que quieres salir de Spotilite?",
		"spotx.adblock": "Bloquear anuncios",
		"spotx.sections": "Bloquear secciones",
		"spotx.premium": "Ocultar Premium",
		"spotx.experiments": "Desactivar experiments",
		"spotx.lyrics_theme": "Tema de letras",
		"spotx.history": "Historial de tracks",
	}

	t.data[LangSpanish] = map[string]string{
		"app.title":                "Spotilite",
		"tray.show":                "Mostrar Ventana",
		"tray.hide":                "Ocultar Ventana",
		"tray.quit":                "Salir",
		"tray.runInBackground":     "Ejecutar en segundo plano",
		"tray.language":            "Idioma",
		"tray.lang.en":             "Inglés",
		"tray.lang.es":             "Español",
		"notif.minimizedToTray":    "Spotilite se está ejecutando en segundo plano",
		"notif.minimizedBody":      "Haz clic en el icono de la bandeja para volver a mostrar la ventana.",
		"dialog.confirmQuit":       "¿Estás seguro de que quieres salir de Spotilite?",
	}
}
