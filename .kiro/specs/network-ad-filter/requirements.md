# Requirements Document

## Introduction

Este documento define los requisitos para implementar un sistema de filtrado de red en Spotilite que bloquee anuncios a nivel de WebView2 sin afectar la interfaz de usuario de Spotify. El sistema reemplazará el módulo AdBlock actual basado en inyección CSS/JS que causa problemas visuales y de rendimiento.

## Glossary

- **Network_Filter**: El componente que intercepta y bloquea requests HTTP/HTTPS de anuncios usando las capacidades de WebView2
- **WebView2**: El componente de Microsoft Edge WebView2 usado por Wails para renderizar contenido web en Windows
- **Ad_Domain**: Un dominio o patrón de URL conocido por servir anuncios de Spotify
- **Block_List**: La lista configurable de dominios y patrones de URL a bloquear
- **AdBlock_Module**: El módulo actual de Spotilite que usa inyección CSS/JS para bloquear anuncios
- **Injector**: El componente de Spotilite que gestiona e inyecta módulos en el WebView
- **API_Server**: El servidor HTTP local (puerto 8765) que expone endpoints para controlar Spotilite
- **WebResourceRequested**: El evento de WebView2 que se dispara antes de cargar un recurso web

## Requirements

### Requirement 1: Network Request Interception

**User Story:** Como desarrollador del sistema, quiero interceptar todas las requests de red del WebView2, para poder analizar y bloquear requests de anuncios antes de que se carguen.

#### Acceptance Criteria

1. WHEN WebView2 intenta cargar un recurso web, THE Network_Filter SHALL intercept the request before it is sent
2. THE Network_Filter SHALL capture the request URL, method, and headers
3. THE Network_Filter SHALL allow the request to proceed if it is not an ad request
4. THE Network_Filter SHALL block the request if it matches an Ad_Domain pattern
5. WHEN a request is blocked, THE Network_Filter SHALL return an empty response with HTTP 200 status
6. THE Network_Filter SHALL not introduce latency greater than 10ms for non-blocked requests

### Requirement 2: Ad Domain Blocking

**User Story:** Como usuario de Spotilite, quiero que los dominios conocidos de anuncios de Spotify sean bloqueados automáticamente, para que no se carguen anuncios durante la reproducción.

#### Acceptance Criteria

1. THE Block_List SHALL include the following Ad_Domain patterns:
   - `ads.spotify.com`
   - `spclient.wg.spotify.com/ads`
   - `audio-ads.spotify.com`
   - `ads-audio.spotify.com`
   - `ad-logger.spotify.com`
   - `pubads.g.doubleclick.net`
   - `doubleclick.net`
   - `googlesyndication.com`
   - `moatads.com`
2. WHEN a request URL contains any Ad_Domain pattern, THE Network_Filter SHALL block the request
3. WHEN a request URL contains `/ad/` or `/ads/` path segments from spotify.com domains, THE Network_Filter SHALL block the request
4. THE Network_Filter SHALL use case-insensitive pattern matching
5. THE Network_Filter SHALL support both exact domain matches and substring pattern matches

### Requirement 3: Block List Configuration

**User Story:** Como desarrollador del sistema, quiero que la lista de dominios bloqueados sea configurable, para poder agregar o remover patrones sin recompilar la aplicación.

#### Acceptance Criteria

1. THE Block_List SHALL be stored as a Go slice in the Network_Filter component
2. THE Network_Filter SHALL provide a method to add new Ad_Domain patterns at runtime
3. THE Network_Filter SHALL provide a method to remove Ad_Domain patterns at runtime
4. THE Network_Filter SHALL provide a method to retrieve the current Block_List
5. WHEN the Block_List is modified, THE Network_Filter SHALL apply changes immediately to subsequent requests

### Requirement 4: Module Integration

**User Story:** Como desarrollador del sistema, quiero integrar el Network_Filter con el sistema de módulos existente, para mantener consistencia con la arquitectura actual de Spotilite.

#### Acceptance Criteria

1. THE Network_Filter SHALL implement the Module interface defined in `internal/spotify/modules/module.go`
2. THE Network_Filter SHALL be registered in the Injector's DefaultModules list
3. WHEN the Network_Filter module is disabled, THE Network_Filter SHALL not intercept or block any requests
4. WHEN the Network_Filter module is enabled, THE Network_Filter SHALL intercept and block Ad_Domain requests
5. THE Network_Filter SHALL respond to enable/disable commands from the API_Server
6. THE Network_Filter SHALL maintain its enabled/disabled state across the application lifecycle

### Requirement 5: AdBlock Module Replacement

**User Story:** Como desarrollador del sistema, quiero reemplazar el AdBlock_Module actual con el Network_Filter, para eliminar los problemas de UI causados por la inyección CSS/JS.

#### Acceptance Criteria

1. THE Network_Filter SHALL not inject any CSS into the Spotify web page
2. THE Network_Filter SHALL not inject any JavaScript into the Spotify web page
3. THE Network_Filter SHALL not modify the DOM of the Spotify web page
4. THE Network_Filter SHALL not use CSS selectors to hide elements
5. WHEN the Network_Filter is enabled, THE AdBlock_Module SHALL be disabled or removed from the Injector
6. THE Network_Filter SHALL provide equivalent or better ad blocking than the current AdBlock_Module

### Requirement 6: WebView2 Event Handling

**User Story:** Como desarrollador del sistema, quiero usar el evento WebResourceRequested de WebView2, para interceptar requests de manera nativa y eficiente.

#### Acceptance Criteria

1. THE Network_Filter SHALL register a handler for the WebResourceRequested event
2. THE Network_Filter SHALL set the appropriate resource context filter to intercept all resource types
3. WHEN a WebResourceRequested event is triggered, THE Network_Filter SHALL evaluate the request URL against the Block_List
4. WHEN a request matches the Block_List, THE Network_Filter SHALL call `PutResponse` with an empty response
5. WHEN a request does not match the Block_List, THE Network_Filter SHALL allow the default behavior
6. THE Network_Filter SHALL handle WebResourceRequested events synchronously to avoid race conditions

### Requirement 7: API Control Endpoints

**User Story:** Como usuario de Spotilite, quiero poder habilitar o deshabilitar el filtrado de red desde la API local, para tener control sobre el comportamiento del bloqueador de anuncios.

#### Acceptance Criteria

1. THE API_Server SHALL expose an endpoint to enable the Network_Filter module
2. THE API_Server SHALL expose an endpoint to disable the Network_Filter module
3. THE API_Server SHALL expose an endpoint to retrieve the Network_Filter enabled status
4. WHEN the enable endpoint is called, THE Network_Filter SHALL start intercepting requests
5. WHEN the disable endpoint is called, THE Network_Filter SHALL stop intercepting requests
6. THE API_Server SHALL return HTTP 200 status for successful enable/disable operations

### Requirement 8: Performance and Compatibility

**User Story:** Como usuario de Spotilite, quiero que el filtrado de red no afecte el rendimiento de carga de Spotify, para mantener una experiencia fluida.

#### Acceptance Criteria

1. THE Network_Filter SHALL not increase page load time by more than 100ms
2. THE Network_Filter SHALL not cause visual flickering or UI glitches
3. THE Network_Filter SHALL be compatible with Windows 10 version 1809 or later
4. THE Network_Filter SHALL be compatible with Windows 11
5. THE Network_Filter SHALL not interfere with legitimate Spotify API requests
6. THE Network_Filter SHALL not block requests to `open.spotify.com` main domain
7. THE Network_Filter SHALL not block requests to `scdn.co` (Spotify CDN) except for ad-specific paths

### Requirement 9: Logging and Debugging

**User Story:** Como desarrollador del sistema, quiero que el Network_Filter registre las requests bloqueadas, para poder depurar y verificar el comportamiento del filtrado.

#### Acceptance Criteria

1. WHEN a request is blocked, THE Network_Filter SHALL log the blocked URL using slog
2. THE Network_Filter SHALL log at INFO level for blocked requests
3. THE Network_Filter SHALL include the matched Ad_Domain pattern in the log message
4. THE Network_Filter SHALL log initialization and shutdown events
5. THE Network_Filter SHALL log enable/disable state changes
6. THE Network_Filter SHALL not log non-blocked requests to avoid log spam

### Requirement 10: Error Handling

**User Story:** Como desarrollador del sistema, quiero que el Network_Filter maneje errores de manera robusta, para que fallos en el filtrado no afecten la funcionalidad general de Spotilite.

#### Acceptance Criteria

1. WHEN WebView2 event registration fails, THE Network_Filter SHALL log the error and continue without blocking
2. WHEN pattern matching fails for a specific URL, THE Network_Filter SHALL allow the request and log the error
3. WHEN the Block_List is empty, THE Network_Filter SHALL allow all requests
4. THE Network_Filter SHALL not crash the application if WebResourceRequested handler encounters an error
5. WHEN an error occurs during request interception, THE Network_Filter SHALL default to allowing the request
6. THE Network_Filter SHALL log all errors at ERROR level using slog
