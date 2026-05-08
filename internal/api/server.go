package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

const DefaultPort = "8765"

type Handler interface {
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

type Settings struct {
	Language      string `json:"language"`
	BackgroundMode bool  `json:"backgroundMode"`
}

type SpotXSettings struct {
	AdBlock      bool   `json:"adBlock"`
	SectionBlock bool   `json:"sectionBlock"`
	PremiumSpoof bool   `json:"premiumSpoof"`
	Experiments  bool   `json:"experiments"`
	LyricsTheme  string `json:"lyricsTheme"`
	TrackHistory bool   `json:"trackHistory"`
}

type Server struct {
	handler Handler
	server  *http.Server
}

func NewServer(handler Handler) *Server {
	return &Server{handler: handler}
}

func (s *Server) SetHandler(handler Handler) {
	s.handler = handler
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	cors := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

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

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	slog.Info("stopping spotilite api server")
	return s.server.Shutdown(ctx)
}

func (s *Server) handleMinimize(w http.ResponseWriter, _ *http.Request) {
	if s.handler != nil {
		s.handler.Minimize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMaximize(w http.ResponseWriter, _ *http.Request) {
	if s.handler != nil {
		s.handler.Maximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUnMaximize(w http.ResponseWriter, _ *http.Request) {
	if s.handler != nil {
		s.handler.UnMaximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleClose(w http.ResponseWriter, _ *http.Request) {
	if s.handler != nil {
		s.handler.Close()
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
	if s.handler != nil {
		s.handler.SetLanguage(req.Lang)
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
	if s.handler != nil {
		s.handler.SetBackgroundMode(req.Enabled)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetSettings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.handler != nil {
		settings := s.handler.GetSettings()
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
	if s.handler != nil {
		s.handler.SetModuleEnabled(req.Module, req.Enabled)
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
	if s.handler != nil {
		s.handler.SetLyricsTheme(req.Theme)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetSpotXSettings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.handler != nil {
		settings := s.handler.GetSpotXSettings()
		json.NewEncoder(w).Encode(settings)
	} else {
		json.NewEncoder(w).Encode(SpotXSettings{
			AdBlock: true, SectionBlock: true, PremiumSpoof: true,
			Experiments: true, LyricsTheme: "neon", TrackHistory: true,
		})
	}
}
