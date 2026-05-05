# Log de Sesiones de Desarrollo

## Sesión 1 — Proyecto Inicial y Navegación Básica

**Fecha**: 2026-05-05

### Objetivo
Crear un cliente desktop ligero de Spotify usando Wails v2.

### Cambios Realizados
- Creado proyecto Wails base.
- Navegación directa a `https://open.spotify.com` mediante `runtime.WindowExecJS`.

### Lecciones Aprendidas
- Spotify es una SPA que re-renderiza elementos periódicamente.

---

## Sesión 2 — Refactorización y Buenas Prácticas

**Fecha**: 2026-05-05

### Objetivo
Aplicar buenas prácticas de Go.

### Cambios Realizados
- Estructura `internal/`: `spotify/injector.go`, `app/app.go`.
- `log/slog` en lugar de `println`.
- Context cancellation en goroutines.

---

## Sesión 3 — System Tray e Internacionalización

**Fecha**: 2026-05-05

### Cambios Realizados
- Paquete `internal/i18n/i18n.go` con ES/EN.
- Paquete `internal/tray/tray.go` con menú contextual.
- Hook `OnBeforeClose` para modo segundo plano.

---

## Sesión 4 — Atajos Globales y Notificaciones

**Fecha**: 2026-05-05

### Cambios Realizados
- `internal/shortcut/shortcut_windows.go` con `RegisterHotKey` (sin CGO).
- Notificaciones con `beeep`.

### Problemas
- `gohook` requiere CGO. Solución: implementación directa con `syscall`.

---

## Sesión 5 — Frameless y Barra de Título Personalizada

**Fecha**: 2026-05-05

### Cambios Realizados
- `Frameless: true`.
- Barra inyectada en Spotify con botones.

### Problemas
- Los botones inyectados no pueden acceder a bindings de Go desde dominio externo.

---

## Sesión 6 — Correcciones Arquitectónicas

**Fecha**: 2026-05-05

### Cambios Realizados
- Restauración de arquitectura iframe + React.
- Corrección de bindings (`window.go.app.App`).
- Documentación completa en `docs/`.

### Problemas
- Spotify bloquea iframes con `X-Frame-Options`.

---

## Sesión 7 — Arquitectura Limpia Final

**Fecha**: 2026-05-05

### Objetivo
Crear una arquitectura estable, limpia y mantenible.

### Decisiones Clave
1. **Frameless definitivo**: Sin barra nativa, sin iframe.
2. **Navegación directa a Spotify**: El webview carga Spotify directamente.
3. **Barra flotante decorativa**: Solo visual (logo + nombre + badge), sin botones.
4. **System tray minimalista**: Solo "Mostrar" y "Salir".
5. **Idioma automático**: Detectado del sistema operativo, sin menú de configuración.
6. **Controles alternativos**: `Ctrl+Shift+S` + system tray.

### Cambios Realizados
- Simplificación de `internal/app/app.go`: eliminados métodos de ventana innecesarios.
- Simplificación de `internal/systray/tray_windows.go`: solo Show y Quit.
- Eliminación de `internal/tray/`: ya no se necesita menú de ventana.
- Actualización de documentación.

### Lecciones Aprendidas
- **No se pueden inyectar controles funcionales en páginas externas** con CSP estricto.
- **La simplicidad es clave**: menos opciones = menos errores.
- **Usar los mecanismos nativos del sistema** (tray, shortcuts) en lugar de reinventar controles de ventana.

### Estado Final
- ✅ Ventana frameless limpia.
- ✅ Spotify carga directamente sin bloqueos.
- ✅ Barra flotante decorativa integrada.
- ✅ System tray minimalista funcional.
- ✅ Atajo global `Ctrl+Shift+S`.
- ✅ Notificaciones nativas.
- ✅ i18n automático.
- ✅ Documentación completa.
