package proxy

import (
	"bufio"
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// AdBlockProxy is a local HTTP proxy that blocks ad requests
type AdBlockProxy struct {
	server   *http.Server
	port     string
	enabled  bool
	mu       sync.RWMutex
	blockedCount int64
}

// Ad URL patterns to block
var adPatterns = []string{
	"doubleclick.net",
	"googlesyndication.com",
	"googleadservices.com",
	"moatads.com",
	"adservice.google.com",
	"audio-ads.spotify.com",
	"ads-audio.spotify.com",
	"ad-logger.spotify.com",
	"ad-handler.spotify.com",
	"spotify.com/ads/",
	"spotify.com/audioad",
	"audio-sp.spotify.com/ad",
	"audio-fa.scdn.co/ad",
	"heads-fa.scdn.co/ad",
	"partnerakamai.spotify.com/ad",
	"megaphone.fm/ad",
	"adclick",
	"adserver",
	"adtech",
	"/ad/",
	"/ads/",
}

// NewAdBlockProxy creates a new ad-blocking proxy
func NewAdBlockProxy(port string) *AdBlockProxy {
	return &AdBlockProxy{
		port:    port,
		enabled: true,
	}
}

// Start starts the proxy server
func (p *AdBlockProxy) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.server != nil {
		return nil // Already running
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handleRequest)

	p.server = &http.Server{
		Addr:         "127.0.0.1:" + p.port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("ad-block proxy starting", "addr", p.server.Addr)
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("proxy server error", "error", err)
		}
	}()

	go func() {
		<-ctx.Done()
		p.Stop()
	}()

	return nil
}

// Stop stops the proxy server
func (p *AdBlockProxy) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.server.Shutdown(ctx)
	p.server = nil
	slog.Info("ad-block proxy stopped", "blocked_count", p.blockedCount)
	return err
}

// handleRequest handles incoming proxy requests
func (p *AdBlockProxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	if !p.enabled {
		p.forwardRequest(w, r)
		return
	}

	// Check if this is an ad request
	if p.isAdRequest(r) {
		p.mu.Lock()
		p.blockedCount++
		count := p.blockedCount
		p.mu.Unlock()

		slog.Debug("blocked ad request", "url", r.URL.String(), "total_blocked", count)
		
		// Return empty response for ad requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}

	// Forward non-ad requests
	p.forwardRequest(w, r)
}

// isAdRequest checks if a request is for an advertisement
func (p *AdBlockProxy) isAdRequest(r *http.Request) bool {
	url := r.URL.String()
	host := r.Host
	
	// Check URL and host against ad patterns
	for _, pattern := range adPatterns {
		if strings.Contains(url, pattern) || strings.Contains(host, pattern) {
			return true
		}
	}
	
	return false
}

// forwardRequest forwards the request to the actual destination
func (p *AdBlockProxy) forwardRequest(w http.ResponseWriter, r *http.Request) {
	// Handle CONNECT method for HTTPS
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}

	// Create a new request
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Copy the original request
	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Send the request
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	bufio.NewReader(resp.Body).WriteTo(w)
}

// handleConnect handles HTTPS CONNECT requests
func (p *AdBlockProxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	// Check if this is an ad domain
	if p.isAdRequest(r) {
		p.mu.Lock()
		p.blockedCount++
		p.mu.Unlock()
		
		slog.Debug("blocked ad CONNECT", "host", r.Host)
		http.Error(w, "Blocked", http.StatusForbidden)
		return
	}

	// Establish connection to the destination
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	// Hijack the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	// Send 200 Connection Established
	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Bidirectional copy
	go func() {
		bufio.NewReader(destConn).WriteTo(clientConn)
	}()
	bufio.NewReader(clientConn).WriteTo(destConn)
}

// SetEnabled enables or disables the proxy
func (p *AdBlockProxy) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
	slog.Info("proxy ad-blocking", "enabled", enabled)
}

func (p *AdBlockProxy) configureSystemProxy(enable bool) {
	if runtime.GOOS == "windows" {
		var args []string
		if enable {
			args = []string{"winhttp", "set", "proxy", "127.0.0.1:" + p.port}
		} else {
			args = []string{"winhttp", "reset", "proxy"}
		}
		cmd := exec.Command("netsh", args...)
		if err := cmd.Run(); err != nil {
			slog.Warn("failed to configure winhttp proxy", "error", err)
		}
	}
}

// GetBlockedCount returns the number of blocked requests
func (p *AdBlockProxy) GetBlockedCount() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.blockedCount
}
