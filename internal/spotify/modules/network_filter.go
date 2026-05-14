package modules

import (
	"sync"
)

// NetworkFilterModule implements network-level ad blocking by intercepting
// HTTP/HTTPS requests at the WebView2 layer before they reach the Spotify web application.
type NetworkFilterModule struct {
	BaseModule
	blockList []string
	mu        sync.RWMutex
}

// NewNetworkFilterModule creates a new NetworkFilterModule with the specified enabled state
// and initializes it with a default block list of known ad domains.
func NewNetworkFilterModule(enabled bool) *NetworkFilterModule {
	return &NetworkFilterModule{
		BaseModule: BaseModule{
			name:    "network_filter",
			enabled: enabled,
		},
		blockList: defaultBlockList(),
	}
}

// defaultBlockList returns the default list of ad domain patterns to block.
// These patterns are matched against request URLs to identify and block ad requests.
func defaultBlockList() []string {
	return []string{
		"ads.spotify.com",
		"spclient.wg.spotify.com/ads",
		"audio-ads.spotify.com",
		"ads-audio.spotify.com",
		"ad-logger.spotify.com",
		"pubads.g.doubleclick.net",
		"doubleclick.net",
		"googlesyndication.com",
		"moatads.com",
		"/ad/",
		"/ads/",
	}
}

// CSS returns an empty string as NetworkFilter does not use CSS injection.
// This satisfies the Module interface while maintaining semantic correctness.
func (m *NetworkFilterModule) CSS() string {
	return ""
}

// JS returns the JavaScript code for network request interception.
// This code wraps the native fetch and XMLHttpRequest APIs to block ad requests.
func (m *NetworkFilterModule) JS() string {
	// TODO: Implement JavaScript interception code in subsequent tasks
	return ""
}

// Selectors returns an empty slice as NetworkFilter does not use DOM selectors.
// This satisfies the Module interface while maintaining semantic correctness.
func (m *NetworkFilterModule) Selectors() []string {
	return []string{}
}

// AddPattern adds a new URL pattern to the block list.
// The pattern will be used to match and block ad requests.
// This method is thread-safe.
func (m *NetworkFilterModule) AddPattern(pattern string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blockList = append(m.blockList, pattern)
}

// RemovePattern removes a URL pattern from the block list.
// This method is thread-safe.
func (m *NetworkFilterModule) RemovePattern(pattern string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for i, p := range m.blockList {
		if p == pattern {
			m.blockList = append(m.blockList[:i], m.blockList[i+1:]...)
			return
		}
	}
}

// GetBlockList returns a copy of the current block list.
// This method is thread-safe.
func (m *NetworkFilterModule) GetBlockList() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to prevent external modification
	result := make([]string, len(m.blockList))
	copy(result, m.blockList)
	return result
}

// SetBlockList replaces the entire block list with a new list of patterns.
// This method is thread-safe.
func (m *NetworkFilterModule) SetBlockList(patterns []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blockList = patterns
}
