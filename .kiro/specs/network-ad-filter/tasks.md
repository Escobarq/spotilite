# Implementation Plan: Network Ad Filter

## Overview

This implementation plan converts the Network Ad Filter design into actionable coding tasks. The feature replaces the current AdBlock module's CSS/JS DOM manipulation approach with network-level request interception using JavaScript fetch/XHR wrapping. This eliminates visual glitches and performance issues while providing more reliable ad blocking.

The implementation follows a 6-phase approach: Core Module → JavaScript Interception → API Integration → AdBlock Deprecation → Testing → Documentation.

## Tasks

- [x] 1. Set up core NetworkFilter module structure
  - [x] 1.1 Create network_filter.go file with NetworkFilterModule struct
    - Create `internal/spotify/modules/network_filter.go`
    - Define `NetworkFilterModule` struct with `BaseModule` embedding
    - Add `blockList []string` field for URL patterns
    - Add `mu sync.RWMutex` for thread-safe block list access
    - Implement constructor `NewNetworkFilterModule(enabled bool)` with default block list
    - _Requirements: 4.1, 4.2, 3.1_
  
  - [x] 1.2 Implement Module interface methods
    - Implement `Name()` returning "network_filter"
    - Implement `Enabled()` and `SetEnabled()` using BaseModule
    - Implement `CSS()` returning empty string (no DOM styling)
    - Implement `Selectors()` returning empty slice (no element hiding)
    - Implement `JS()` method stub (will add JavaScript code in Phase 2)
    - _Requirements: 4.1, 5.1, 5.2, 5.3, 5.4_
  
  - [x] 1.3 Implement block list management methods
    - Implement `AddPattern(pattern string)` with validation and thread safety
    - Implement `RemovePattern(pattern string)` with thread safety
    - Implement `GetBlockList() []string` with read lock
    - Implement `SetBlockList(patterns []string)` with write lock
    - Add `defaultBlockList()` function returning default ad domain patterns
    - Add `isValidPattern(pattern string)` validation function
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 10.3_

- [ ] 2. Checkpoint - Verify module structure compiles
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 3. Implement JavaScript network interception code
  - [ ] 3.1 Create URL pattern matching logic in JavaScript
    - Write `isAdUrl(url)` function with case-insensitive substring matching
    - Implement pattern matching against block list array
    - Add whitelist protection for legitimate Spotify domains (open.spotify.com, api.spotify.com, scdn.co)
    - Add special handling for `/ad/` and `/ads/` path patterns
    - Include error handling with fail-open behavior
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 8.5, 8.6, 8.7, 10.2, 10.5_
  
  - [ ] 3.2 Implement fetch API interception
    - Wrap native `window.fetch` with custom implementation
    - Store original fetch as `window.__origFetch`
    - Check request URL against block list before sending
    - Return empty Response with HTTP 200 for blocked requests
    - Call original fetch for allowed requests
    - Add console logging for blocked requests
    - _Requirements: 1.1, 1.3, 1.4, 1.5, 9.1, 9.2, 9.3_
  
  - [ ] 3.3 Implement XMLHttpRequest interception
    - Wrap `XMLHttpRequest.prototype.open` to capture request URL
    - Wrap `XMLHttpRequest.prototype.send` to block flagged requests
    - Set `__adBlocked` flag on XHR instance for blocked requests
    - Return synthetic HTTP 200 response for blocked XHR requests
    - Add console logging for blocked XHR requests
    - _Requirements: 1.1, 1.3, 1.4, 1.5, 9.1, 9.2, 9.3_
  
  - [ ] 3.4 Create response generation functions
    - Implement `createEmptyResponse()` returning JSON response with HTTP 200
    - Implement `createEmptyAudioResponse()` returning empty audio buffer with HTTP 200
    - Set appropriate Content-Type headers for each response type
    - _Requirements: 1.5_
  
  - [ ] 3.5 Integrate JavaScript code into JS() method
    - Embed complete JavaScript interception code in `JS()` method
    - Use template string or raw string literal for readability
    - Include default block list patterns in JavaScript code
    - Add initialization guard to prevent double-injection
    - Add console log for successful initialization
    - _Requirements: 1.1, 1.2, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [ ] 4. Checkpoint - Test JavaScript interception manually
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Integrate NetworkFilter with module system
  - [ ] 5.1 Register NetworkFilter in DefaultModules
    - Add `modules.NewNetworkFilterModule(true)` to `DefaultModules()` in injector.go
    - Position NetworkFilter before AdBlockModule in module list
    - Ensure NetworkFilter is enabled by default
    - _Requirements: 4.2, 4.3_
  
  - [ ] 5.2 Add module retrieval helper in Injector
    - Verify `GetModule(name string)` method exists in Injector
    - Ensure it can retrieve NetworkFilter by name "network_filter"
    - _Requirements: 4.5_

- [ ] 6. Implement API control endpoints
  - [ ] 6.1 Add NetworkFilter enable/disable to API server
    - Update `internal/api/server.go` to handle "network_filter" module name
    - Ensure POST `/api/spotx/module` endpoint supports NetworkFilter
    - Add enable/disable logic using `SetEnabled()` method
    - Return HTTP 200 for successful operations
    - _Requirements: 7.1, 7.2, 7.4, 7.5, 7.6_
  
  - [ ] 6.2 Add NetworkFilter status endpoint
    - Update GET `/api/spotx/settings` or equivalent to include NetworkFilter state
    - Return NetworkFilter enabled status in JSON response
    - _Requirements: 7.3_
  
  - [ ] 6.3 Update App startup to expose NetworkFilter control
    - Verify `internal/app/app.go` exposes module control methods
    - Ensure NetworkFilter can be controlled via API server
    - _Requirements: 4.5, 4.6_

- [ ] 7. Implement AdBlock module deprecation logic
  - [ ] 7.1 Add automatic AdBlock disabling when NetworkFilter is enabled
    - Add logic in `App.Startup()` to check NetworkFilter state
    - If NetworkFilter is enabled, automatically disable AdBlockModule
    - Log the migration action using slog at INFO level
    - _Requirements: 5.5_
  
  - [ ] 7.2 Ensure mutual exclusivity of NetworkFilter and AdBlock
    - When NetworkFilter is enabled via API, disable AdBlock
    - When AdBlock is enabled via API, optionally disable NetworkFilter (user choice)
    - _Requirements: 5.5_

- [ ] 8. Checkpoint - Integration testing
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. Add logging and error handling
  - [ ] 9.1 Implement Go-side logging
    - Add slog.Info log when NetworkFilter is initialized
    - Add slog.Info log when NetworkFilter is enabled/disabled
    - Add slog.Error log if module initialization fails
    - Add slog.Warn log for invalid patterns
    - _Requirements: 9.4, 9.5, 10.1_
  
  - [ ] 9.2 Implement JavaScript-side logging
    - Add console.log for each blocked request with URL and matched pattern
    - Add console.log for NetworkFilter initialization
    - Add console.error for pattern matching errors
    - Ensure logs use "[NetworkFilter]" prefix for easy identification
    - _Requirements: 9.1, 9.2, 9.3_
  
  - [ ] 9.3 Add error handling for edge cases
    - Handle empty block list (allow all requests)
    - Handle pattern matching exceptions (fail open)
    - Handle JavaScript injection failures (log error, continue)
    - Ensure errors never crash the application
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6_

- [ ] 10. Write unit tests for NetworkFilter module
  - [ ]* 10.1 Write block list management tests
    - Test `AddPattern()` with valid patterns
    - Test `AddPattern()` with empty pattern (should reject)
    - Test `AddPattern()` with invalid patterns (should reject)
    - Test `RemovePattern()` removes existing pattern
    - Test `GetBlockList()` returns current patterns
    - Test `SetBlockList()` replaces entire list
    - Test thread safety with concurrent access
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_
  
  - [ ]* 10.2 Write pattern validation tests
    - Test `isValidPattern()` rejects empty strings
    - Test `isValidPattern()` rejects control characters
    - Test `isValidPattern()` rejects wildcard patterns ("*", ".")
    - Test `isValidPattern()` rejects patterns blocking Spotify domains
    - Test `isValidPattern()` accepts legitimate ad domain patterns
    - _Requirements: 10.1, 10.2_
  
  - [ ]* 10.3 Write Module interface tests
    - Test `Name()` returns "network_filter"
    - Test `Enabled()` and `SetEnabled()` work correctly
    - Test `CSS()` returns empty string
    - Test `Selectors()` returns empty slice
    - Test `JS()` returns non-empty JavaScript code
    - _Requirements: 4.1, 5.1, 5.2, 5.3, 5.4_

- [ ] 11. Write integration tests
  - [ ]* 11.1 Write API integration tests
    - Test POST `/api/spotx/module` enables NetworkFilter
    - Test POST `/api/spotx/module` disables NetworkFilter
    - Test GET endpoint returns NetworkFilter status
    - Test API returns HTTP 200 for successful operations
    - _Requirements: 7.1, 7.2, 7.3, 7.6_
  
  - [ ]* 11.2 Write module system integration tests
    - Test NetworkFilter is registered in DefaultModules
    - Test Injector can retrieve NetworkFilter by name
    - Test NetworkFilter JavaScript is injected when enabled
    - Test NetworkFilter JavaScript is not injected when disabled
    - Test AdBlock is disabled when NetworkFilter is enabled
    - _Requirements: 4.2, 4.3, 4.4, 5.5_

- [ ] 12. Update documentation
  - [ ] 12.1 Update ARCHITECTURE.md
    - Add NetworkFilter module to architecture diagram
    - Document network interception approach
    - Explain difference from AdBlock module
    - _Requirements: All_
  
  - [ ] 12.2 Update FEATURES.md
    - Add NetworkFilter feature description
    - Document API endpoints for control
    - Add usage examples
    - Document migration from AdBlock
    - _Requirements: All_
  
  - [ ] 12.3 Update README.md if needed
    - Add NetworkFilter to feature list
    - Update ad blocking description
    - _Requirements: All_

- [ ] 13. Final checkpoint - Manual testing and validation
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional test tasks and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation throughout implementation
- The NetworkFilter uses JavaScript-based interception (not native WebView2 events) due to Wails API limitations
- The design is architected to support future native WebView2 integration if Wails exposes the API
- AdBlock module is deprecated but not removed, allowing users to fall back if needed
- All network filtering errors fail open (allow requests) to prevent breaking Spotify functionality
- Block list patterns use simple substring matching for performance (<10ms per request)
- Protected domains (open.spotify.com, api.spotify.com, scdn.co) are never blocked

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["1.2", "1.3"] },
    { "id": 2, "tasks": ["3.1", "3.2", "3.3", "3.4"] },
    { "id": 3, "tasks": ["3.5", "5.1", "5.2"] },
    { "id": 4, "tasks": ["6.1", "6.2", "6.3", "7.1", "7.2"] },
    { "id": 5, "tasks": ["9.1", "9.2", "9.3"] },
    { "id": 6, "tasks": ["10.1", "10.2", "10.3", "11.1", "11.2"] },
    { "id": 7, "tasks": ["12.1", "12.2", "12.3"] }
  ]
}
```
