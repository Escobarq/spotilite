# Log de Sesiones de Desarrollo

## Sesión 1 — Proyecto Inicial y Navegación Básica

**Fecha**: 2026-05-05

### Objetivo
Crear un cliente desktop ligero de Spotify usando Wails v2.

### Cambios Realizados
- Creado proyecto Wails base con `wails init`.
- Navegación directa a `https://open.spotify.com` mediante `runtime.WindowExecJS` con `window.location.replace`.
- Eliminado el uso de `<iframe>` del frontend.

### Lecciones Aprendidas
- Navegar a URLs externas directamente en el webview funciona, pero limita la interacción con los bindings de Wails.
- Spotify es una SPA (Single Page Application) que re-renderiza elementos periódicamente.

---

## Sesión 2 — Refactorización y Buenas Prácticas

**Fecha**: 2026-05-05

### Objetivo
Aplicar buenas prácticas de Go: separación de responsabilidades, logging estructurado y context cancellation.

### Cambios Realizados
- **Nueva estructura** de directorios bajo `internal/`:
  - `internal/spotify/injector.go` — Servicio dedicado a la inyección de CSS.
  - `internal/app/app.go` — Orquesta la aplicación Wails.
- Reemplazado `println` por `log/slog` con contexto estructurado.
- Goroutine del injector ahora respeta `ctx.Done()` para graceful shutdown.
- Eliminado `app.go` de la raíz; `main.go` ahora solo inicializa y corre.

### Lecciones Aprendidas
- `main.go` debe ser mínimo: solo crear instancias y llamar a la función de ejecución.
- `context.Context` es esencial para el control del ciclo de vida de goroutines.

---

## Sesión 3 — System Tray e Internacionalización

**Fecha**: 2026-05-05

### Objetivo
Añadir un menú contextual en la bandeja del sistema y soporte multi-idioma.

### Cambios Realizados
- **Nuevo paquete** `internal/i18n/i18n.go`:
  - Traducciones ES/EN con `sync.RWMutex` para thread-safety.
  - Fallback automático a inglés.
- **Nuevo paquete** `internal/tray/tray.go`:
  - Menú contextual usando `runtime.MenuSetApplicationMenu`.
  - Opciones: Mostrar/Ocultar, Idioma (submenú), Salir.
  - Menú se refresca al cambiar de idioma.
- **Nuevo hook** `OnBeforeClose`: al cerrar la ventana con `✕`, se oculta en lugar de destruirse si el modo segundo plano está activo.

### Lecciones Aprendidas
- Wails v2 no tiene `runtime.TraySetMenu`; se usa `runtime.MenuSetApplicationMenu`.
- Los menús se deben reconstruir completamente al cambiar de idioma.

---

## Sesión 4 — Atajos Globales y Notificaciones

**Fecha**: 2026-05-05

### Objetivo
Implementar atajos de teclado globales y notificaciones nativas al minimizar.

### Cambios Realizados
- **Nuevo paquete** `internal/shortcut/shortcut_windows.go`:
  - Usa `RegisterHotKey` de `user32.dll` mediante `syscall` (sin CGO).
  - `Ctrl+Shift+S` para mostrar/ocultar ventana.
  - Stub para Linux/macOS.
- **Notificaciones** con `github.com/gen2brain/beeep`:
  - Toast nativo al minimizar a la bandeja.
  - Texto localizado según el idioma activo.

### Problemas Encontrados
- `github.com/robotn/gohook` requiere CGO y no compilaba en Windows sin mingw.
- Solución: reemplazar por implementación directa con `syscall`.

### Lecciones Aprendidas
- Evitar dependencias con CGO para mayor portabilidad.
- `beeep` genera notificaciones nativas sin configuración adicional en Windows.

---

## Sesión 5 — Frameless y Barra de Título Personalizada

**Fecha**: 2026-05-05

### Objetivo
Eliminar la barra de título nativa de Windows (que seguía siendo blanca) y crear una propia integrada con Spotify.

### Cambios Realizados
- `Frameless: true` en `options.App`.
- Eliminado tema oscuro nativo de Windows (no funcionaba consistentemente).
- **Inyección de barra personalizada** en `injector.go`:
  - Script JavaScript que crea un `div` fijo en la parte superior de Spotify.
  - Botones: Minimizar, Maximizar, Cerrar.
  - Logo de Spotify + texto "Spotilite".
  - `--wails-draggable:drag` para mover la ventana.
- CSS que empuja el contenido de Spotify 32px hacia abajo para no taparlo.

### Problemas Encontrados
- Los botones inyectados en una página externa no pueden acceder a `window.go.app.App` porque los bindings de Wails no existen en ese contexto.
- Spotify tiene CSP (Content Security Policy) que bloquea scripts inyectados.

### Lecciones Aprendidas
- Los bindings de Wails solo existen en el contexto del frontend local (`localhost` / `wails://`).
- No se pueden inyectar controles nativos en páginas externas con CSP estricto.

---

## Sesión 6 — Correcciones Arquitectónicas y Estabilidad

**Fecha**: 2026-05-05

### Objetivo
Restaurar la funcionalidad de los botones de la barra de título.

### Cambios Realizados
- **Restaurada arquitectura iframe**:
  - El webview principal carga el frontend local de React.
  - React renderiza la barra de título personalizada + iframe de Spotify.
  - Los botones usan `window.go.app.App.*` directamente desde React (contexto local).
- **Correcciones**:
  - Ruta correcta de bindings: `window.go.app.App` (no `window.go.main.App`).
  - Botones con `--wails-draggable:no-drag` para evitar que el drag intercepte los clics.
  - `e.stopPropagation()` en `onmousedown` y `onclick`.
- **Eliminada** la barra inyectada desde Go (ya no es necesaria).
- **Documentación** completa:
  - `README.md` actualizado.
  - `docs/ARCHITECTURE.md` — explicación técnica detallada.
  - `docs/SETUP.md` — guía de instalación.
  - `docs/FEATURES.md` — funcionalidades implementadas y planificadas.
  - `docs/SESSIONS.md` — este documento.

### Decisiones Arquitectónicas Clave
- **Iframe + React** es la arquitectura correcta para esta app: mantiene los bindings funcionales mientras muestra Spotify.
- **Go maneja** system tray, atajos, notificaciones e inyección CSS en el iframe.
- **React maneja** la UI de la barra de título y el layout general.

### Estado Final
- ✅ Ventana frameless con barra de título personalizada funcional.
- ✅ System tray nativo con icono visible.
- ✅ Atajo global `Ctrl+Shift+S`.
- ✅ Notificaciones nativas.
- ✅ i18n ES/EN.
- ✅ Modo segundo plano configurable.
- ✅ Inyección CSS ocultando botones no deseados de Spotify.
- ✅ Documentación completa.
