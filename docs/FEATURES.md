# Lista de Funcionalidades

## Funcionalidades Implementadas

### Interfaz de Usuario

- [x] **Ventana frameless** sin barra de título nativa del sistema operativo.
- [x] **Barra de título personalizada** con colores de Spotify (`#191414`).
  - [x] Logo de Spotify + nombre "Spotilite".
  - [x] Botón Minimizar (`─`).
  - [x] Botón Maximizar / Restaurar (`□` / `❐`).
  - [x] Botón Cerrar (`✕`) con hover rojo.
- [x] **Arrastre de ventana** desde la barra de título personalizada.
- [x] **Iframe de Spotify** sin bordes ni barras de desplazamiento visibles.

### System Tray (Bandeja del Sistema)

- [x] **Icono nativo** visible en la bandeja de Windows.
- [x] **Menú contextual** con clic derecho:
  - [x] Mostrar Ventana.
  - [x] Ocultar Ventana.
  - [x] Toggle "Ejecutar en segundo plano" con checkmark visual `[x]`.
  - [x] Submenú de idioma (Español / English).
  - [x] Salir.
- [x] **Menú refrescable** al cambiar de idioma.

### Controles de Ventana

- [x] **Minimizar** ventana.
- [x] **Maximizar** ventana.
- [x] **Restaurar** ventana desde maximizado.
- [x] **Cerrar** ventana (respeta modo segundo plano).

### Atajos de Teclado

- [x] `Ctrl + Shift + S` — Mostrar / Ocultar ventana desde cualquier aplicación.

### Notificaciones

- [x] **Notificación nativa** al minimizar a la bandeja (modo segundo plano activo).
- [x] Texto localizado según el idioma seleccionado.

### Internacionalización (i18n)

- [x] Soporte para **Español** (`es`).
- [x] Soporte para **Inglés** (`en`).
- [x] Cambio de idioma en tiempo real desde el menú de la bandeja.
- [x] Fallback automático a inglés si falta una traducción.

### Modo Segundo Plano

- [x] Al cerrar la ventana con `✕`, si está activado, se oculta en lugar de cerrar.
- [x] Icono permanece en la bandeja.
- [x] Notificación toast informa que sigue ejecutándose.
- [x] Toggle desde el menú contextual de la bandeja.

### Personalización de Spotify

- [x] **Ocultar botón "Instalar app de escritorio"** (`data-testid="download-button"`).
- [x] **Ocultar botón "Conectar a un dispositivo"** (`data-testid="device-picker-icon-button"`).
- [x] Inyección CSS periódica (cada 3 segundos) para persistir cambios en la SPA.

### Compatibilidad

- [x] Windows 10/11 (soporte completo: frameless, tray, shortcuts, notificaciones).
- [x] macOS (básico: frameless, ventana, controles).
- [x] Linux (básico: frameless, ventana, controles).

## Funcionalidades Planificadas

### Corto Plazo

- [ ] **Persistencia de preferencias**: Guardar idioma y modo segundo plano entre sesiones.
- [ ] **Soporte completo macOS**: System tray nativo y atajos de teclado.
- [ ] **Soporte completo Linux**: System tray nativo y atajos de teclado.

### Medio Plazo

- [ ] **Integración con Spotify API**: Controles multimedia nativos (play/pause/next/prev).
- [ ] **Media Keys**: Soporte para teclas multimedia del teclado.
- [ ] **Mini player**: Ventana flotante pequeña con controles básicos.

### Largo Plazo

- [ ] **Temas personalizables**: Colores de la barra de título configurable.
- [ ] **Plugins**: Sistema de extensiones para personalizar la experiencia.
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

### Sesión 6 — Correcciones y Estabilidad
- Restauración de arquitectura iframe + React para controles funcionales.
- Corrección de bindings (`window.go.app.App`).
- Corrección de draggable en botones (`--wails-draggable:no-drag`).
- Documentación completa (`docs/`).
