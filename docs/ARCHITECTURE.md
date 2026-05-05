# Arquitectura de Spotilite

Este documento describe la arquitectura técnica del proyecto, las decisiones de diseño y cómo interactúan los componentes.

## Visión General

Spotilite es una aplicación híbrida compuesta por:

1. **Backend Go** (Wails) — Controla la ventana nativa frameless, system tray, atajos de teclado, notificaciones e inyección de CSS.
2. **Spotify Web** (WebView directo) — La aplicación web de Spotify cargada directamente en el webview de Wails.
3. **Barra flotante inyectada** — Elemento visual inyectado en la página de Spotify mediante JavaScript.

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    Spotilite Desktop App                     │
│           (Wails WebView — ventana frameless)               │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Barra flotante Spotilite (inyectada en Spotify)      │  │
│  │  • Logo + "Spotilite" + badge "Desktop"               │  │
│  │  • Informativa / decorativa                           │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  open.spotify.com (cargado directamente en webview)   │  │
│  │  • CSS inyectado oculta botones no deseados           │  │
│  │  • Controles via: System Tray + Atajos globales       │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend Go (Wails)                        │
├─────────────────────────────────────────────────────────────┤
│  internal/app/app.go          │  Ciclo de vida + controles  │
│  internal/spotify/injector.go │  Navegación + CSS injection │
│  internal/systray/            │  Icono nativo en bandeja    │
│  internal/shortcut/           │  Atajos globales (Win API)  │
│  internal/i18n/i18n.go        │  Traducciones ES/EN         │
└─────────────────────────────────────────────────────────────┘
```

## Decisiones Clave de Diseño

### Frameless + Barra Flotante Inyectada

**Problema**: La barra de título nativa de Windows no se integra visualmente con la interfaz oscura de Spotify.

**Solución**:
- `Frameless: true` en `main.go` elimina completamente la barra nativa y los bordes.
- Una **barra flotante** se inyecta en la página de Spotify mediante `runtime.WindowExecJS`.
- La barra es puramente decorativa/informativa (logo + nombre + badge).
- No contiene botones porque desde el dominio externo de Spotify no se pueden llamar bindings de Go.

**Ventaja**: Apariencia limpia e integrada con Spotify. Sin elementos nativos que rompan la estética.

### Navegación Directa (Sin Iframe)

**Problema**: Spotify envía headers `X-Frame-Options` que bloquean la carga en iframes desde dominios externos.

**Solución**:
- El webview de Wails navega **directamente** a `https://open.spotify.com` mediante `window.location.replace`.
- Esto evita completamente las restricciones de iframe.

**Ventaja**: Spotify carga siempre, sin errores de seguridad.

### Controles de Ventana Alternativos

**Problema**: Sin barra nativa, no hay botones de minimizar/maximizar/cerrar visibles.

**Solución**:
- **System Tray**: Icono nativo con menú minimalista (Mostrar, Salir).
- **Atajo Global**: `Ctrl+Shift+S` alterna visibilidad desde cualquier aplicación.
- **Cerrar**: El botón de cerrar del system tray cierra la app definitivamente.
- **Minimizar/Ocultar**: Al intentar cerrar la ventana (`Alt+F4` o similares), se oculta a la bandeja.

**Ventaja**: Controles accesibles sin ocupar espacio visual en la interfaz.

### System Tray Nativo Minimalista

**Filosofía**: El system tray debe tener solo lo esencial.

**Implementación**:
- Solo dos opciones: **Mostrar Ventana** y **Salir**.
- No hay submenús, toggles ni configuraciones en el tray.
- El icono es siempre visible en la bandeja de Windows.

**Ventaja**: Simplicidad. El usuario no se pierde entre opciones.

### Idioma Automático

**Filosofía**: La app debe adaptarse al sistema, no requerir configuración manual.

**Implementación**:
- Al iniciar, detecta el idioma del sistema operativo.
- Si es español → usa español. Si no → inglés.
- No hay menú de cambio de idioma.

**Ventaja**: Cero configuración para el usuario.

### Atajos de Teclado sin CGO

**Problema**: Las librerías de atajos globales en Go suelen requerir CGO.

**Solución**:
- `internal/shortcut/shortcut_windows.go` usa directamente la API de Windows (`RegisterHotKey` en `user32.dll`) mediante `syscall`.
- No requiere CGO. Compilación pura en Go.

**Ventaja**: `go build` funciona sin configuración adicional.

## Flujo de Datos

### Arranque de la Aplicación

1. `main.go` crea instancias de `i18n`, `systray.Manager`, `app.App`.
2. `wails.Run` inicia el webview en modo frameless.
3. `app.Startup` se ejecuta:
   - Inicia el system tray (`systray.Start()`).
   - Registra el atajo global (`shortcut.Register`).
   - Inicia el inyector CSS (`injector.Start()`).
4. El inyector navega a Spotify y comienza a inyectar CSS periódicamente.

### Cierre de la Ventana

1. Usuario cierra la ventana (`Alt+F4` o similar).
2. `app.OnBeforeClose` intercepta el cierre.
3. La ventana se oculta (`runtime.Hide`).
4. Notificación toast: *"Spotilite se está ejecutando en segundo plano"*.
5. La app sigue viva. El icono permanece en la bandeja.

### Mostrar/Ocultar Ventana

1. Usuario presiona `Ctrl+Shift+S` o hace clic en "Mostrar Ventana" en el tray.
2. `app.ToggleWindowVisibility` alterna entre `runtime.Show` y `runtime.Hide`.

## Estructura de Directorios

```
spotilite/
├── main.go                        # Punto de entrada
├── internal/
│   ├── app/
│   │   └── app.go                 # Ciclo de vida Wails + controles
│   ├── i18n/
│   │   └── i18n.go                # Traducciones ES/EN
│   ├── shortcut/
│   │   ├── shortcut_windows.go    # Atajo global (RegisterHotKey, sin CGO)
│   │   └── shortcut_stub.go       # Stub para Linux/macOS
│   ├── spotify/
│   │   └── injector.go            # Navegación + CSS injection + barra flotante
│   └── systray/
│       ├── tray_windows.go        # Icono nativo + menú minimalista
│       └── tray_stub.go           # Stub para Linux/macOS
├── frontend/
│   ├── src/
│   │   ├── App.jsx                # Pantalla de carga inicial
│   │   ├── main.jsx               # Entry point React
│   │   └── style.css              # Estilos globales
│   ├── wailsjs/                   # Bindings auto-generados por Wails
│   └── dist/                      # Build output
├── build/
│   └── windows/
│       ├── icon.ico               # Icono del ejecutable y system tray
│       └── wails.exe.manifest     # Manifiesto de Windows
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

- **Navegación directa**: Al navegar directamente a Spotify, la app carga contenido de terceros. El usuario debe confiar en Spotify.
- **Inyección CSS**: El código inyectado solo oculta elementos visuales y añade la barra flotante. No accede a datos personales.
- **CGO desactivado**: La compilación usa `CGO_ENABLED=0` para mayor portabilidad.

## Extensiones Futuras

- Soporte para Linux y macOS en `shortcut` y `systray`.
- Preferencias persistentes (guardadas en `os.UserConfigDir`).
- Integración con la API de Spotify (requiere OAuth2) para controles multimedia nativos.
