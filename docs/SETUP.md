# Guía de Instalación y Desarrollo

## Requisitos Previos

### Go

Instala Go 1.23 o superior desde [go.dev](https://go.dev/dl/).

Verifica la instalación:
```bash
go version
```

### Node.js y npm

Instala Node.js 18+ desde [nodejs.org](https://nodejs.org/).

Verifica la instalación:
```bash
node --version
npm --version
```

### Wails CLI

Instala la herramienta de línea de comandos de Wails:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verifica la instalación:
```bash
wails version
```

### Dependencias de Windows

- **Windows 10/11**
- **WebView2 Runtime** (generalmente ya instalado en Windows 11)
- Si no lo tienes, descárgalo desde [Microsoft Edge WebView2](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)

## Clonar el Proyecto

```bash
git clone https://github.com/escobarq/spotilite.git
cd spotilite
```

## Instalar Dependencias

### Backend (Go)

```bash
go mod tidy
```

### Frontend (React)

```bash
cd frontend
npm install
cd ..
```

## Ejecutar en Modo Desarrollo

```bash
wails dev
```

Esto:
1. Inicia el frontend en modo hot-reload (Vite).
2. Compila el backend Go.
3. Abre la ventana de Spotilite.
4. Recarga automáticamente cuando cambias archivos del frontend.

## Compilar para Producción

### Windows

```bash
wails build
```

El ejecutable se genera en `build/bin/Spotilite.exe`.

### Con Icono Personalizado

Asegúrate de que `build/windows/icon.ico` existe. Wails lo incrusta automáticamente en el ejecutable.

### Para Otras Plataformas

```bash
# macOS
wails build -platform darwin/amd64

# Linux
wails build -platform linux/amd64
```

**Nota**: `shortcut` y `systray` son stubs en macOS/Linux. La funcionalidad completa requiere implementaciones específicas para cada plataforma.

## Solución de Problemas

### Error: "WebView2 not found"

Instala el runtime de WebView2:
https://developer.microsoft.com/en-us/microsoft-edge/webview2/

### Error: "wails.exe.manifest" no encontrado

Asegúrate de tener el archivo `build/windows/wails.exe.manifest`. Si lo perdiste, puedes regenerar la estructura con:
```bash
wails init
```

### La barra de título sigue siendo blanca

Esto puede ocurrir si el manifiesto de Windows no declara compatibilidad con Windows 10/11. Verifica que `build/windows/wails.exe.manifest` contenga el nodo `<compatibility>` con el GUID `{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}`.

### Los botones de la barra no funcionan

1. Verifica que los bindings de Wails estén generados:
   ```bash
   wails generate module
   ```
2. Revisa la consola del desarrollador en la ventana de Wails (`Ctrl+Shift+I`).
3. Asegúrate de que `window.go.app.App` existe en el contexto del webview.

### El icono no aparece en la bandeja

El ejecutable busca el icono en `build/windows/icon.ico` relativo a su ubicación. Si mueves el `.exe`, asegúrate de que la carpeta `build/windows/` esté junto a él.

En modo desarrollo (`wails dev`), el icono se busca en `./build/windows/icon.ico` desde la raíz del proyecto.

## Comandos Útiles

| Comando | Descripción |
|---|---|
| `wails dev` | Modo desarrollo con hot-reload |
| `wails build` | Compilar para producción |
| `wails generate module` | Regenerar bindings Go→JS |
| `go mod tidy` | Limpiar dependencias Go |
| `go build ./...` | Compilar Go sin ejecutar |
| `cd frontend && npm run build` | Compilar frontend manualmente |
