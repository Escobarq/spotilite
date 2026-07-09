# Guía de Instalación y Desarrollo

## Requisitos Previos

### Go

Instala Go 1.22 o superior desde [go.dev](https://go.dev/dl/).

Verifica la instalación:
```bash
go version
```

### Wails CLI

Instala la herramienta de línea de comandos de Wails (necesaria sólo si quieres empaquetar iconos/manifiestos o ejecutar en modo dev):

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verifica la instalación:
```bash
wails version
```

> **No se requiere Node.js, npm ni Vite.** El frontend React fue eliminado; el splash inicial es un `index.html` embebido en el binario vía `//go:embed`.

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

```bash
go mod tidy
```

Eso es todo. Sólo dependencias Go.

## Ejecutar en Modo Desarrollo

```bash
wails dev
```

Esto:
1. Compila el backend Go en modo hot-reload.
2. Abre la ventana de Spotilite con el splash nativo.
3. Recarga automáticamente cuando cambias archivos Go (runtime.WindowReload no es necesario; los cambios se inyectan cada vez que guardas).

Sin Wails CLI también puedes:
```bash
go build -o spotilite.exe . && ./spotilite.exe
```

## Compilar para Producción

### Windows

```bash
wails build
```

El ejecutable se genera en `build/bin/Spotilite.exe`.

Alternativa sin Wails CLI:
```bash
go build -ldflags="-H windowsgui -s -w" -o spotilite.exe .
```

### Con Icono Personalizado

Asegúrate de que `build/windows/icon.ico` existe. Wails lo incrusta automáticamente en el ejecutable; `go build` directo sólo respeta este icono si la build se hace vía Wails.

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

### Error leyendo `config-xpui.ini`

Spotilite intenta leer y mantener el archivo que spicetify-cli usa. Si no existe (primera vez), arranca con un Config vacío y crea el archivo la primera vez que guardes. Para apuntar a un config diferente, exporta `SPICETIFY_CONFIG=/ruta/al/archivo` antes de lanzar Spotilite.

### Las extensiones no se cargan

1. Verifica que el archivo existe en `%APPDATA%\spicetify\Extensions\` (o donde indique `UserExtensionsDir`).
2. Revisa la consola del webview (`Ctrl+Shift+I` si `Wails dev` está en modo dev): las extensiones que fallen aparecen como `[Spotilite] extension <nombre> error:`.
3. Algunas extensiones Spicetify requieren `Spicetify.GraphQL.Definitions` o acceso a `sp://` que **no están disponibles** en la web player. Esas extensiones no pueden funcionar en Spotilite (ver tabla en `README.md`).

### El icono no aparece en la bandeja

El ejecutable busca el icono en `build/windows/icon.ico` relativo a su ubicación. Si mueves el `.exe`, asegúrate de que la carpeta `build/windows/` esté junto a él.

## Comandos Útiles

| Comando | Descripción |
|---|---|
| `wails dev` | Modo desarrollo con hot-reload |
| `wails build` | Compilar para producción |
| `go mod tidy` | Limpiar dependencias Go |
| `go build ./...` | Compilar Go sin ejecutar (todos los paquetes) |
| `go test ./...` | Ejecutar la suite de tests |
| `go vet ./...` | Validación estática |
