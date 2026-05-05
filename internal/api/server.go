// Package api provides a local HTTP server that exposes window controls and
// settings endpoints. It allows the JavaScript injected into Spotify to
// communicate with the Go backend via fetch requests.
package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// DefaultPort is the local port for the API server.
const DefaultPort = "8765"

// Handler is the interface that App implements to handle API requests.
type Handler interface {
	Minimize()
	Maximize()
	UnMaximize()
	Close()
	SetLanguage(lang string)
	SetBackgroundMode(enabled bool)
	GetSettings() Settings
}

// Settings holds the current app settings.
type Settings struct {
	Language       string `json:"language"`
	BackgroundMode bool   `json:"backgroundMode"`
}

// Server is the local HTTP API server.
type Server struct {
	handler Handler
	server  *http.Server
}

// NewServer creates a new API server. The handler can be set later via SetHandler.
func NewServer(handler Handler) *Server {
	return &Server{handler: handler}
}

// SetHandler sets the handler after App has been created.
func (s *Server) SetHandler(handler Handler) {
	s.handler = handler
}

// Start launches the API server in a background goroutine.
func (s *Server) Start() {
	mux := http.NewServeMux()

	// CORS middleware to allow requests from any origin.
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

// Stop gracefully shuts down the API server.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	slog.Info("stopping spotilite api server")
	return s.server.Shutdown(ctx)
}

func (s *Server) handleMinimize(w http.ResponseWriter, r *http.Request) {
	if s.handler != nil {
		s.handler.Minimize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMaximize(w http.ResponseWriter, r *http.Request) {
	if s.handler != nil {
		s.handler.Maximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUnMaximize(w http.ResponseWriter, r *http.Request) {
	if s.handler != nil {
		s.handler.UnMaximize()
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleClose(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.handler != nil {
		settings := s.handler.GetSettings()
		json.NewEncoder(w).Encode(settings)
	} else {
		json.NewEncoder(w).Encode(Settings{Language: "es", BackgroundMode: true})
	}
}
