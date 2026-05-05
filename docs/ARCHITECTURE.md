# Arquitectura de Spotilite

Este documento describe la arquitectura técnica del proyecto, las decisiones de diseño y cómo interactúan los componentes.

## Visión General

Spotilite es una aplicación híbrida compuesta por:

1. **Backend Go** (Wails) — Controla la ventana nativa, system tray, atajos de teclado, notificaciones e inyección de CSS.
2. **Frontend React** (local) — Renderiza la barra de título personalizada y el iframe de Spotify.
3. **Spotify Web** (iframe) — La aplicación web de Spotify cargada dentro del iframe.

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    Spotilite Desktop App                     │
│  (Wails WebView — ventana frameless sin barra de título)    │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Barra de Título (React)                              │  │
│  │  • Logo + "Spotilite"                                 │  │
│  │  • Botones: ─ □ ✕  →  window.go.app.App.*            │  │
│  │  • Arrastrable (--wails-draggable:drag)               │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Iframe: open.spotify.com                             │  │
│  │  • CSS inyectado oculta botones no deseados           │  │
│  │  • Se mantiene dentro de la ventana frameless         │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend Go (Wails)                        │
├─────────────────────────────────────────────────────────────┤
│  internal/app/app.go          │  Ciclo de vida de la app    │
│  internal/spotify/injector.go │  Inyección CSS en iframe    │
│  internal/tray/tray.go        │  Menú contextual de ventana │
│  internal/systray/            │  Icono nativo en bandeja    │
│  internal/shortcut/           │  Atajos globales (Win API)  │
│  internal/i18n/i18n.go        │  Traducciones ES/EN         │
└─────────────────────────────────────────────────────────────┘
```

## Decisiones Clave de Diseño

### Frameless + Barra de Título Propia

**Problema**: La barra de título nativa de Windows no se integra visualmente con la interfaz oscura de Spotify.

**Solución**:
- `Frameless: true` en `main.go` elimina completamente la barra nativa.
- React renderiza una barra personalizada con los colores de Spotify (`#191414`).
- Los botones llaman directamente a los bindings de Wails (`window.go.app.App.Minimize()`).

**Ventaja**: Control total sobre la apariencia y los controles.

### Iframe en lugar de Navegación Directa

**Problema**: Cuando Wails navega directamente a `open.spotify.com`, los bindings de Go no están disponibles en el contexto externo debido a las políticas de seguridad (CSP) de Spotify.

**Solución**:
- El webview principal carga el frontend local de React (donde los bindings existen).
- React incluye un `<iframe src="https://open.spotify.com">` que carga Spotify.
- La barra de título personalizada (React) puede comunicarse con Go sin restricciones.

**Ventaja**: Los botones de la barra funcionan perfectamente y se mantiene la experiencia nativa.

### System Tray Nativo (getlantern/systray)

**Problema**: Wails no proporciona un system tray con menú contextual nativo y icono visible por defecto.

**Solución**:
- `internal/systray/tray_windows.go` usa `getlantern/systray` para crear un icono real en la bandeja de Windows.
- El menú contextual incluye: Mostrar/Ocultar, toggle de segundo plano, idioma, salir.
- Se refresca dinámicamente cuando cambia el idioma.

**Ventaja**: Icono persistente visible junto al reloj de Windows.

### Atajos de Teclado sin CGO

**Problema**: Las librerías de atajos globales en Go suelen requerir CGO, lo que complica la compilación cruzada.

**Solución**:
- `internal/shortcut/shortcut_windows.go` usa directamente la API de Windows (`RegisterHotKey` en `user32.dll`) mediante `syscall`.
- No requiere CGO. Compilación pura en Go.

**Ventaja**: `go build` funciona sin configuración adicional.

## Flujo de Datos

### Arranque de la Aplicación

1. `main.go` crea instancias de `i18n`, `systray.Manager`, `app.App`.
2. `wails.Run` inicia el webview en modo frameless.
3. El webview carga `frontend/dist/index.html` (React local).
4. React renderiza la barra de título + iframe de Spotify.
5. `app.Startup` se ejecuta:
   - Inicia el system tray (`systray.Start()`).
   - Registra el atajo global (`shortcut.Register`).
   - Inicia el inyector CSS (`injector.Start()`).

### Cierre de la Ventana (Modo Segundo Plano)

1. Usuario pulsa `✕` → React llama `window.go.app.App.Close()`.
2. `app.Close()` ejecuta `runtime.Hide(ctx)`.
3. Se muestra una notificación nativa: *"Spotilite se está ejecutando en segundo plano"*.
4. La app sigue viva. El icono permanece en la bandeja.

### Cambio de Idioma

1. Usuario selecciona un idioma en el menú de la bandeja.
2. `systray.Manager` actualiza `i18n.SetLanguage()`.
3. `systray.Manager.Refresh()` reconstruye el menú con las nuevas traducciones.

## Estructura de Directorios

```
spotilite/
├── main.go                        # Punto de entrada
├── internal/
│   ├── app/
│   │   └── app.go                 # Ciclo de vida Wails + controles de ventana
│   ├── i18n/
│   │   └── i18n.go                # Motor de traducciones (ES/EN)
│   ├── shortcut/
│   │   ├── shortcut_windows.go    # Atajo global (RegisterHotKey, sin CGO)
│   │   └── shortcut_stub.go       # Stub para Linux/macOS
│   ├── spotify/
│   │   └── injector.go            # Inyección CSS en iframe de Spotify
│   ├── systray/
│   │   ├── tray_windows.go        # Icono nativo + menú contextual
│   │   └── tray_stub.go           # Stub para Linux/macOS
│   └── tray/
│       └── tray.go                # Menú de la ventana (Wails MenuBar)
├── frontend/
│   ├── src/
│   │   ├── App.jsx                # Barra de título + iframe
│   │   ├── main.jsx               # Entry point React
│   │   └── style.css              # Estilos globales
│   ├── wailsjs/                   # Bindings auto-generados por Wails
│   └── dist/                      # Build output
├── build/
│   └── windows/
│       ├── icon.ico               # Icono del ejecutable y system tray
│       └── wails.exe.manifest     # Manifiesto de Windows (compatibilidad Win10/11)
├── docs/
│   ├── ARCHITECTURE.md            # Este documento
│   ├── FEATURES.md                # Lista de funcionalidades
│   ├── SESSIONS.md                # Log de sesiones de desarrollo
│   └── SETUP.md                   # Guía de instalación
├── README.md
├── wails.json
├── go.mod
└── go.sum
```

## Consideraciones de Seguridad

- **CSP de Spotify**: El iframe de Spotify no puede ejecutar scripts del contexto padre por restricciones de seguridad del navegador. Por eso la barra de título debe estar fuera del iframe.
- **CGO desactivado**: La compilación usa `CGO_ENABLED=0` para mayor portabilidad. Las funciones nativas de Windows se acceden mediante `syscall` en lugar de CGO.

## Extensiones Futuras

- Soporte para Linux y macOS en `shortcut` y `systray`.
- Preferencias persistentes (idioma, modo segundo plano) usando `os.UserConfigDir`.
- Integración con la API de Spotify (requiere OAuth2) para controles multimedia nativos.
- Tema personalizable (colores de la barra de título).
