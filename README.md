# Spotilite

Un cliente desktop ligero para Spotify construido con [Wails](https://wails.io/) y Go puro. Embebe el reproductor web de Spotify en una ventana frameless e implementa su propia capa de compatibilidad con [Spicetify](https://spicetify.app/) (extensiones, temas, custom apps) sobre la web pública, sin tocar el bundle del cliente desktop.

![Spotilite](build/appicon.png)

## Características

- 🎵 **Spotify Web** embebido directamente en una ventana frameless sin barras de título del sistema operativo.
- 🖥️ **Ventana frameless** con barra de información flotante estilo Spotify integrada en la interfaz.
- 🔔 **System Tray** minimalista con icono nativo.
- ⌨️ **Atajo global** `Ctrl+Shift+S` para mostrar/ocultar la ventana desde cualquier lugar.
- 🌍 **Idioma automático** detectado del sistema operativo (Español / Inglés).
- 🛡️ **Modo segundo plano**: minimiza a la bandeja en lugar de cerrar.
- 🎨 **Tema oscuro** integrado con los colores de Spotify.
- 🔕 **Oculta elementos de Spotify**: botón de descarga de app y conexión a dispositivos.
- 🧩 **Compatibilidad Spicetify**: lee `config-xpui.ini` y carga extensiones, temas (`color.ini`) y custom apps desde la carpeta estándar de Spicetify.

## Requisitos

- [Go](https://go.dev/) 1.22+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2.12+ (solo para empaquetar)
- Windows 10/11 (soporte completo) / macOS / Linux

## Instalación Rápida

```bash
# Clonar el repositorio
git clone https://github.com/escobarq/spotilite.git
cd spotilite

# Compilar para producción
wails build

# O, sin Wails CLI (solo binario Go):
go build -o spotilite.exe .
```

## Uso

### Controles de Ventana

Como la ventana es **frameless** (sin barra nativa), los controles son:

| Acción | Método |
|---|---|
| **Mover ventana** | Arrastra desde la barra de información flotante de Spotilite |
| **Mostrar / Ocultar** | `Ctrl + Shift + S` desde cualquier aplicación |
| **Cerrar app** | Menú de bandeja → **Salir** |

### Atajos de Teclado

| Atajo | Acción |
|---|---|
| `Ctrl + Shift + S` | Mostrar / Ocultar ventana |

### Menú de Bandeja (System Tray)

Haz clic derecho en el icono de Spotilite en la bandeja del sistema:

- **Mostrar Ventana**
- **Salir**

## Compatibilidad con Spicetify

Spotilite reimplementa dentro del webview la superficie de la API `window.Spicetify` que el CLI oficial expone mutando el bundle del cliente desktop. Lo hace así porque en `https://open.spotify.com` no se puede reescribir el webpack bundle: el código se sirve por red y todos los parches regex de `preprocess.go` (que dependen de strings minificados de `xpui.js`) son inaccesibles.

Lo que **funciona** en Spotilite:

| Característica Spicetify | Estado en Spotilite |
|---|---|
| `Spicetify.Player` (play/pause/next/seek/eventos) | Eventos vía `MediaSession` + `DOM` polling, comandos vía click simulado en los botones del reproductor |
| `Spicetify.URI` (parser/serializer `spotify:track:...`) | Porte limpio, sin dependencias |
| `Spicetify.CosmosAsync` (get/post/request) | `fetch` a la API web pública; fallback a través del proxy Go local para endpoints CORS-locked |
| `Spicetify.LocalStorage` / `Mousetrap`-alias `Keyboard` | Implementaciones nativas |
| `Spicetify.Menu` / `Spicetify.ContextMenu` | Menús flotantes con estilos forzados |
| `Spicetify.Topbar.Button` / `Spicetify.Playbar.Button` | Mount DOM en contenedores reales de Spotify |
| `Spicetify.PopupModal` / `Spicetify.Panel` | Overlays flotantes |
| `Spicetify.React` / `Spicetify.ReactDOM` | Shim que mapea a nodos DOM (no es React real) |
| Temas (`color.ini` + `user.css`) | Soportados; `--spice-*` se traduce a `--background-base`/`--text-base`/etc. |
| Extensiones `.js` / `.mjs` | Cargadas desde la carpeta `Extensions/` (mismo path que spicetify-cli) |
| Custom apps (manifest.json + index.js) | Soportado vía `history.pushState` + sidebar inyectado |

Lo que **NO funciona** (y por qué):

- `Spicetify.GraphQL.Definitions` — se popula por regex sobre `xpui.js` en el cliente desktop. En web player no es alcanzable.
- `Spicetify.Platform.LibraryAPI` / `PlaylistAPI` / `PlaybackAPI` — vienen de `Symbol.for("PlaybackAPI")` en el Registry inyectado por los patches de webpack. En web hay que navegar el `webpackChunkclient_web` cache en runtime (frágil, no soportado).
- Esquemas `sp://` / `wg://` / `hm://` — registrados por el CEF bridge del cliente desktop; reescritos a `https://api.spotify.com/` / `https://spclient.wg.spotify.com/` cuando aplica.

Una extensión que use `Spicetify.GraphQL.Definitions["someQuery"]` registrará un warning en consola y no funcionará. Una extensión que use `Spicetify.Player`, `Spicetify.LocalStorage`, `Spicetify.CosmosAsync.get(...)` para endpoints de la API pública funcionará normalmente.

### Configuración

Spotilite lee y escribe el mismo archivo que spicetify-cli:

- **Windows**: `%APPDATA%\spicetify\config-xpui.ini`
- **macOS/Linux**: `~/.config/spicetify/config-xpui.ini`

Secciones relevantes:

```ini
[AdditionalOptions]
current_theme = Dribbblish
color_scheme  = dark
extensions    = shuffle.js|trashbin.js|autoplay.js-
custom_apps   = reddit|lyrics-
```

Las extensiones habilitadas se buscan en `Extensions/`, los temas en `Themes/<name>/color.ini`, y los custom apps en `CustomApps/<name>/{manifest.json,index.js}`. Si esa carpeta no existe, Spotilite también acepta `os.UserConfigDir()/spotilite/` como fallback.

## Arquitectura

Consulta [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) para una explicación detallada de la estructura del proyecto, decisiones de diseño y flujo de datos.

## Documentación

- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — Arquitectura y diseño
- [`docs/FEATURES.md`](docs/FEATURES.md) — Lista detallada de funcionalidades
- [`docs/SESSIONS.md`](docs/SESSIONS.md) — Log de sesiones de desarrollo
- [`docs/SETUP.md`](docs/SETUP.md) — Guía de instalación y desarrollo

## Tecnologías

- **Backend**: Go 1.22+, Wails v2.12
- **Frontend**: ninguno (HTML/CSS incrustado en el binario vía `//go:embed` para el splash inicial)
- **System Tray**: getlantern/systray (Windows)
- **Notificaciones**: gen2brain/beeep
- **Atajos Globales**: Windows API nativa vía `syscall` (sin CGO)

## Licencia

MIT

## Autor

- **Escobarq** — [@escobarq](https://github.com/escobarq)
