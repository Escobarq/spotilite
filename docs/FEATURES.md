# Lista de Funcionalidades

## Funcionalidades Implementadas

### Interfaz de Usuario

- [x] **Ventana frameless** sin barra de título nativa del sistema operativo.
- [x] **Barra flotante decorativa** integrada en la parte superior de Spotify.
  - [x] Logo de Spotify (14px).
  - [x] Texto "Spotilite".
  - [x] Badge verde "DESKTOP".
  - [x] Gradiente oscuro Spotify (`#191414` → `#121212`).
  - [x] Empuja el contenido de Spotify 28px hacia abajo para no taparlo.
- [x] **Tema oscuro** con fondo negro puro (`#000000`).

### System Tray (Bandeja del Sistema)

- [x] **Icono nativo** visible en la bandeja de Windows.
- [x] **Menú contextual minimalista** con clic derecho:
  - [x] Mostrar Ventana.
  - [x] Salir.
- [x] **Diseño simplificado**: solo lo esencial, sin submenús ni toggles.

### Controles de Ventana

- [x] **Mostrar / Ocultar** ventana.
- [x] **Cerrar** aplicación (desde system tray).

### Atajos de Teclado

- [x] `Ctrl + Shift + S` — Mostrar / Ocultar ventana desde cualquier aplicación.

### Notificaciones

- [x] **Notificación nativa** al minimizar a la bandeja.
- [x] Texto localizado según el idioma del sistema.

### Internacionalización (i18n)

- [x] Soporte para **Español** (`es`).
- [x] Soporte para **Inglés** (`en`).
- [x] **Detección automática** del idioma del sistema operativo al inicio.
- [x] Fallback automático a inglés.

### Modo Segundo Plano

- [x] Al cerrar la ventana, se oculta en lugar de destruirse.
- [x] Icono permanece en la bandeja.
- [x] Notificación toast informa que sigue ejecutándose.
- [x] Siempre activo (no requiere configuración).

### Personalización de Spotify

- [x] **Ocultar botón "Instalar app de escritorio"**.
- [x] **Ocultar botón "Conectar a un dispositivo"**.
- [x] Inyección CSS periódica (cada 3 segundos) para persistir cambios en la SPA.

### Compatibilidad

- [x] Windows 10/11 (soporte completo: frameless, tray, shortcuts, notificaciones).
- [x] macOS (básico: frameless, ventana).
- [x] Linux (básico: frameless, ventana).

## Funcionalidades Planificadas

### Corto Plazo

- [ ] **Persistencia de preferencias**: Guardar estado entre sesiones.
- [ ] **Soporte completo macOS**: System tray nativo y atajos de teclado.
- [ ] **Soporte completo Linux**: System tray nativo y atajos de teclado.

### Medio Plazo

- [ ] **Integración con Spotify API**: Controles multimedia nativos.
- [ ] **Media Keys**: Soporte para teclas multimedia del teclado.
- [ ] **Mini player**: Ventana flotante pequeña con controles básicos.

### Largo Plazo

- [ ] **Temas personalizables**: Colores de la barra flotante configurable.
- [ ] **Plugins**: Sistema de extensiones.
- [ ] **Actualizaciones automáticas**: Integración con GoReleaser + auto-updater.

## Historial de Cambios

### Sesión 1 — Fundamentos
- Navegación directa a Spotify en webview.
- Inyección CSS para ocultar botones no deseados.

### Sesión 2 — Refactor y Buenas Prácticas
- Separación de responsabilidades (`internal/`).
- Logging con `slog`.
- Context cancellation en goroutines.

### Sesión 3 — System Tray e i18n
- Menú contextual nativo en bandeja.
- Traducciones ES/EN.
- Modo segundo plano configurable.

### Sesión 4 — Atajos y Notificaciones
- `Ctrl+Shift+S` global shortcut.
- Notificaciones nativas con `beeep`.
- Tema oscuro de Spotify.

### Sesión 5 — Frameless y Barra Personalizada
- Ventana sin marco nativo.
- Barra de título inyectada en Spotify.

### Sesión 6 — Correcciones Arquitectónicas
- Restauración de arquitectura iframe + React para controles funcionales.
- Corrección de bindings.
- Documentación completa.

### Sesión 7 — Arquitectura Limpia Final
- **Eliminación de barra nativa**: Frameless definitivo.
- **Simplificación del tray**: Solo Show y Quit.
- **Idioma automático**: Detección del sistema operativo.
- **Barra flotante decorativa**: Logo + nombre + badge, sin botones.
- **Controles alternativos**: Tray + atajo global.
