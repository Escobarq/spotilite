# Spotilite

Un cliente desktop ligero para Spotify construido con [Wails](https://wails.io/) y React. Proporciona una experiencia nativa integrada con el sistema operativo: ventana frameless con barra de información flotante, system tray minimalista, atajos de teclado globales y notificaciones nativas.

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

## Requisitos

- [Go](https://go.dev/) 1.23+
- [Node.js](https://nodejs.org/) 18+ con npm
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2.12+
- Windows 10/11 (soporte completo) / macOS / Linux

## Instalación Rápida

```bash
# Clonar el repositorio
git clone https://github.com/escobarq/spotilite.git
cd spotilite

# Instalar dependencias del frontend
cd frontend && npm install && cd ..

# Ejecutar en modo desarrollo
wails dev

# Compilar para producción
wails build
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

*(Diseño minimalista: solo lo esencial)*

## Arquitectura

Consulta [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) para una explicación detallada de la estructura del proyecto, decisiones de diseño y flujo de datos.

## Documentación

- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — Arquitectura y diseño
- [`docs/FEATURES.md`](docs/FEATURES.md) — Lista detallada de funcionalidades
- [`docs/SESSIONS.md`](docs/SESSIONS.md) — Log de sesiones de desarrollo
- [`docs/SETUP.md`](docs/SETUP.md) — Guía de instalación y desarrollo

## Tecnologías

- **Backend**: Go 1.23, Wails v2
- **Frontend**: React 18, Vite (solo para pantalla de carga inicial)
- **System Tray**: getlantern/systray (Windows)
- **Notificaciones**: gen2brain/beeep
- **Atajos Globales**: Windows API (RegisterHotKey, sin CGO)

## Licencia

MIT

## Autor

- **Escobarq** — [@escobarq](https://github.com/escobarq)
