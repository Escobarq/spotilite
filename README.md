# Spotilite

Un cliente desktop ligero para Spotify construido con [Wails](https://wails.io/) y React. Proporciona una experiencia nativa integrada con el sistema operativo: ventana frameless con controles personalizados, system tray, atajos de teclado globales, notificaciones nativas e internacionalización.

![Spotilite](build/appicon.png)

## Características

- 🎵 **Spotify Web** embebido en una ventana nativa sin barras de título del sistema operativo.
- 🖥️ **Ventana frameless** con barra de título personalizada estilo Spotify.
- 🔔 **System Tray** con icono nativo y menú contextual.
- ⌨️ **Atajo global** `Ctrl+Shift+S` para mostrar/ocultar la ventana desde cualquier lugar.
- 🌍 **Internacionalización** (i18n) con soporte para Español e Inglés.
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

La barra de título personalizada incluye:

| Botón | Acción |
|---|---|
| `─` | Minimizar ventana |
| `□` / `❐` | Maximizar / Restaurar |
| `✕` | Cerrar (respeta modo segundo plano) |

### Atajos de Teclado

| Atajo | Acción |
|---|---|
| `Ctrl + Shift + S` | Mostrar / Ocultar ventana |

### Menú de Bandeja (System Tray)

Haz clic derecho en el icono de Spotilite en la bandeja del sistema:

- **Mostrar Ventana** / **Hide Window**
- **[x] Ejecutar en segundo plano** — toggle
- **Idioma** → Español / English
- **Salir**

## Arquitectura

Consulta [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) para una explicación detallada de la estructura del proyecto, decisiones de diseño y flujo de datos.

## Documentación

- [`docs/SETUP.md`](docs/SETUP.md) — Guía completa de instalación y desarrollo
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — Arquitectura y diseño
- [`docs/FEATURES.md`](docs/FEATURES.md) — Lista detallada de funcionalidades
- [`docs/SESSIONS.md`](docs/SESSIONS.md) — Log de sesiones de desarrollo

## Tecnologías

- **Backend**: Go 1.23, Wails v2
- **Frontend**: React 18, Vite
- **System Tray**: getlantern/systray (Windows)
- **Notificaciones**: gen2brain/beeep
- **Atajos Globales**: Windows API (RegisterHotKey, sin CGO)

## Licencia

MIT

## Autor

- **Escobarq** — [@escobarq](https://github.com/escobarq)
