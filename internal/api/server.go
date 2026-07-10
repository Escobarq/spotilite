package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const DefaultPort = "8765"

// WindowHandler mirrors the desktop integration surface of the Spotilite app
// (window controls, language, etc.). The Server is constructed with one of
// these; pass nil during very early boot and SetHandler later.
type WindowHandler interface {
	Minimize()
	Maximize()
	UnMaximize()
	Close()
	SetLanguage(lang string)
	SetBackgroundMode(enabled bool)
	GetSettings() Settings
	SetModuleEnabled(name string, enabled bool)
	SetLyricsTheme(theme string)
	GetSpotXSettings() SpotXSettings
}

// SpicetifyHandler is implemented by *app.App and gives the API server
// read/write access to the Spicetify config, extension list, theme chooser and
// reload trigger so webview-side extensions can interact with the Go side.
//
// All errors in this surface are best-effort. Spicetify-style extensions are
// "fail-open" — they should keep running if the bridge is unavailable.
type SpicetifyHandler interface {
	GetSpicetifyConfig() SpicetifyConfigDTO
	GetSpicetifyExtensions() []ExtensionDTO
	SetSpicetifyExtension(name string, enabled bool) error
	GetSpicetifyThemes() []string
	SetSpicetifyTheme(name, colorScheme string) error
	GetSpicetifyCustomApps() []string
	SetSpicetifyCustomApp(name string, enabled bool) error
	ReloadInjection(ctx context.Context) error
}

// Settings is the public settings DTO returned by /api/settings.
type Settings struct {
	Language       string `json:"language"`
	BackgroundMode bool   `json:"backgroundMode"`
}

// SpotXSettings mirrors the SpotX-style toggle module flags referenced in the
// existing /api/spotx handlers.
type SpotXSettings struct {
	AdBlock      bool   `json:"adBlock"`
	SectionBlock bool   `json:"sectionBlock"`
	PremiumSpoof bool   `json:"premiumSpoof"`
	Experiments  bool   `json:"experiments"`
	LyricsTheme  string `json:"lyricsTheme"`
	TrackHistory bool   `json:"trackHistory"`
}

// SpicetifyConfigDTO is the snapshot exposed to webview-side extensions via
// /api/spicetify/config. The fields are tuned to the same names used inside
// the injected JS payload (Spicetify.Config) so an extension can consult both
// without translation.
type SpicetifyConfigDTO struct {
	Version              string   `json:"version"`
	CurrentTheme         string   `json:"current_theme"`
	ColorScheme          string   `json:"color_scheme"`
	Extensions           []string `json:"extensions"`
	CustomApps           []string `json:"custom_apps"`
	InjectCSS            bool     `json:"inject_css"`
	InjectThemeJS        bool     `json:"inject_theme_js"`
	ReplaceColors        bool     `json:"replace_colors"`
	SidebarConfig        bool     `json:"sidebar_config"`
	HomeConfig           bool     `json:"home_config"`
	ExperimentalFeatures bool     `json:"experimental_features"`
	LocalAPI             string   `json:"local_api"`
}

// ExtensionDTO represents one Spicetify extension file known to the
// webview (name + whether enabled). The webview side receives this so it can
// surface a UI without needing Go bindings.
type ExtensionDTO struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Enabled bool   `json:"enabled"`
	IsMJS   bool   `json:"is_mjs"`
}

// Server is the local HTTP API front-end used by extensions running inside the
// Spotify webview. It exposes two surfaces:
//   - /api/window/*    window controls (legacy)
//   - /api/settings/*  language + background mode toggles (legacy)
//   - /api/spotx/*     SpotX toggle features (legacy)
//   - /api/spicetify/* config + extensions + themes bridge (new)
type Server struct {
	window     WindowHandler
	spice      SpicetifyHandler
	server     *http.Server
	appVersion string
}

// NewServer builds a Server with the given window handler. The Spicetify
// surface is wired later via SetSpicetifyHandler once main has constructed
// the App.
func NewServer(handler WindowHandler) *Server {
	return &Server{window: handler, appVersion: "spotilite-1.0"}
}

// SetHandler updates the window-control handler (replaces the constructor arg
// in main.go when finalization happens after App.NewApp).
func (s *Server) SetHandler(handler WindowHandler) {
	s.window = handler
}

// SetSpicetifyHandler installs the Spicetify config bridge. Call once during
// boot after App.NewApp.
func (s *Server) SetSpicetifyHandler(h SpicetifyHandler) {
	s.spice = h
}

// Version returns the version string reported in API DTOs.
func (s *Server) Version() string { return s.appVersion }

// Start listens on :DefaultPort and serves until Stop() is called. Safe to
// call from a goroutine.
func (s *Server) Start() {
	mux := http.NewServeMux()

	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Legacy endpoints preserved from the previous version.
	mux.HandleFunc("/api/window/minimize", cors(s.handleMinimize))
	mux.HandleFunc("/api/window/maximize", cors(s.handleMaximize))
	mux.HandleFunc("/api/window/unmaximize", cors(s.handleUnMaximize))
	mux.HandleFunc("/api/window/close", cors(s.handleClose))
	mux.HandleFunc("/api/settings/lang", cors(s.handleSetLang))
	mux.HandleFunc("/api/settings/background", cors(s.handleSetBackground))
	mux.HandleFunc("/api/settings", cors(s.handleGetSettings))
	mux.HandleFunc("/api/spotx/module", cors(s.handleSetModule))
	mux.HandleFunc("/api/spotx/lyrics_theme", cors(s.handleSetLyricsTheme))
	mux.HandleFunc("/api/spotx/settings", cors(s.handleGetSpotXSettings))

	// Spicetify endpoints (new).
	mux.HandleFunc("/api/spicetify/config", cors(s.handleSpiceGetConfig))
	mux.HandleFunc("/api/spicetify/extensions", cors(s.handleSpiceExtensions))
	mux.HandleFunc("/api/spicetify/extension/toggle", cors(s.handleSpiceExtensionToggle))
	mux.HandleFunc("/api/spicetify/themes", cors(s.handleSpiceGetThemes))
	mux.HandleFunc("/api/spicetify/theme", cors(s.handleSpiceSetTheme))
	mux.HandleFunc("/api/spicetify/customapps", cors(s.handleSpiceGetCustomApps))
	mux.HandleFunc("/api/spicetify/customapp/toggle", cors(s.handleSpiceCustomAppToggle))
	mux.HandleFunc("/api/spicetify/reload", cors(s.handleSpiceReload))

	s.server = &http.Server{
		Addr:         ":" + DefaultPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("starting spotilite api server", "address", "http://localhost:"+DefaultPort)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("api server error", "error", err)
		}
	}()
}

// Stop shuts down the API server.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	slog.Info("stopping spotilite api server")
	return s.server.Shutdown(ctx)
}

// --- legacy handlers ----------------------------------------------------------

func (s *Server) handleMinimize(w http.ResponseWriter, _ *http.Request) {
	if s.window != nil {
		s.window.Minimize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMaximize(w http.ResponseWriter, _ *http.Request) {
	if s.window != nil {
		s.window.Maximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUnMaximize(w http.ResponseWriter, _ *http.Request) {
	if s.window != nil {
		s.window.UnMaximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleClose(w http.ResponseWriter, _ *http.Request) {
	if s.window != nil {
		s.window.Close()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSetLang(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Lang string `json:"lang"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.window != nil {
		s.window.SetLanguage(req.Lang)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSetBackground(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.window != nil {
		s.window.SetBackgroundMode(req.Enabled)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetSettings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.window != nil {
		settings := s.window.GetSettings()
		json.NewEncoder(w).Encode(settings)
	} else {
		json.NewEncoder(w).Encode(Settings{Language: "es", BackgroundMode: true})
	}
}

func (s *Server) handleSetModule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Module  string `json:"module"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.window != nil {
		s.window.SetModuleEnabled(req.Module, req.Enabled)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSetLyricsTheme(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Theme string `json:"theme"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if s.window != nil {
		s.window.SetLyricsTheme(req.Theme)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetSpotXSettings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.window != nil {
		settings := s.window.GetSpotXSettings()
		json.NewEncoder(w).Encode(settings)
	} else {
		json.NewEncoder(w).Encode(SpotXSettings{
			AdBlock: true, SectionBlock: true, PremiumSpoof: true,
			Experiments: true, LyricsTheme: "neon", TrackHistory: true,
		})
	}
}

// --- spicetify handlers -------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) handleSpiceGetConfig(w http.ResponseWriter, r *http.Request) {
	if s.spice == nil {
		writeJSON(w, http.StatusServiceUnavailable, SpicetifyConfigDTO{LocalAPI: snapLocalURL(r)})
		return
	}
	dto := s.spice.GetSpicetifyConfig()
	dto.LocalAPI = snapLocalURL(r)
	writeJSON(w, http.StatusOK, dto)
}

func (s *Server) handleSpiceExtensions(w http.ResponseWriter, _ *http.Request) {
	if s.spice == nil {
		writeJSON(w, http.StatusOK, []ExtensionDTO{})
		return
	}
	writeJSON(w, http.StatusOK, s.spice.GetSpicetifyExtensions())
}

func (s *Server) handleSpiceExtensionToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	if s.spice == nil {
		writeError(w, http.StatusServiceUnavailable, "spice handler offline")
		return
	}
	var req struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.spice.SetSpicetifyExtension(strings.TrimSpace(req.Name), req.Enabled); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSpiceGetThemes(w http.ResponseWriter, _ *http.Request) {
	if s.spice == nil {
		writeJSON(w, http.StatusOK, []string{})
		return
	}
	writeJSON(w, http.StatusOK, s.spice.GetSpicetifyThemes())
}

func (s *Server) handleSpiceSetTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	if s.spice == nil {
		writeError(w, http.StatusServiceUnavailable, "spice handler offline")
		return
	}
	var req struct {
		Name   string `json:"name"`
		Scheme string `json:"scheme"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.spice.SetSpicetifyTheme(strings.TrimSpace(req.Name), strings.TrimSpace(req.Scheme)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSpiceGetCustomApps(w http.ResponseWriter, _ *http.Request) {
	if s.spice == nil {
		writeJSON(w, http.StatusOK, []string{})
		return
	}
	writeJSON(w, http.StatusOK, s.spice.GetSpicetifyCustomApps())
}

func (s *Server) handleSpiceCustomAppToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	if s.spice == nil {
		writeError(w, http.StatusServiceUnavailable, "spice handler offline")
		return
	}
	var req struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.spice.SetSpicetifyCustomApp(strings.TrimSpace(req.Name), req.Enabled); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSpiceReload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	if s.spice == nil {
		writeError(w, http.StatusServiceUnavailable, "spice handler offline")
		return
	}
	if err := s.spice.ReloadInjection(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// snapLocalURL returns the base URL of the API server as the webview would
// see it (always http://localhost:<port>). Used to populate local_api in the
// Spicetify config DTO so extensions and the injected JS agree on the same
// address.
func snapLocalURL(r *http.Request) string {
	return "http://localhost:" + DefaultPort
}
