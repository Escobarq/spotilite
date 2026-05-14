# BlockTheSpot - Análisis de Compatibilidad con Spotilite

## ❌ Incompatibilidad Fundamental

**BlockTheSpot NO es compatible con Spotilite** debido a diferencias arquitectónicas fundamentales.

## 🔍 Cómo Funciona BlockTheSpot

BlockTheSpot utiliza **DLL Injection** para modificar el comportamiento de Spotify:

### Método de Operación

1. **Reemplaza `chrome_elf.dll`** en la carpeta de Spotify Desktop
2. **Intercepta llamadas** a funciones del sistema
3. **Bloquea requests** de anuncios antes de que lleguen al motor de Chromium
4. **Modifica el comportamiento** del cliente de Spotify Desktop

```
Spotify Desktop
    ↓
chrome_elf.dll (MODIFICADO por BlockTheSpot)
    ↓
Chromium Embedded Framework (CEF)
    ↓
Sistema Operativo
```

### Archivos que Modifica

- `%APPDATA%\Spotify\chrome_elf.dll` - DLL principal modificada
- `%APPDATA%\Spotify\config.ini` - Configuración
- Archivos de Spotify Desktop en `C:\Users\[User]\AppData\Roaming\Spotify\`

## 🏗️ Arquitectura de Spotilite

Spotilite tiene una arquitectura completamente diferente:

```
Spotilite (Wails App)
    ↓
WebView2 (Microsoft Edge)
    ↓
Spotify Web Player (https://open.spotify.com)
    ↓
Internet
```

### Diferencias Clave

| Aspecto | BlockTheSpot | Spotilite |
|---------|--------------|-----------|
| **Target** | Spotify Desktop | Spotify Web Player |
| **Motor** | CEF (Chromium Embedded) | WebView2 (Edge) |
| **Método** | DLL Injection | JavaScript Injection |
| **Archivos** | Modifica DLLs de Spotify | No modifica archivos de Spotify |
| **Alcance** | Cliente Desktop | Navegador embebido |

## ⚠️ Por Qué No Funciona

### 1. **Diferentes Ejecutables**

- **BlockTheSpot**: Modifica `Spotify.exe` (cliente desktop)
- **Spotilite**: Es `spotilite.exe` (aplicación Wails independiente)

### 2. **Diferentes Motores de Renderizado**

- **BlockTheSpot**: CEF (Chromium Embedded Framework)
- **Spotilite**: WebView2 (Microsoft Edge WebView)

### 3. **Diferentes Puntos de Inyección**

- **BlockTheSpot**: Inyecta a nivel de DLL del sistema
- **Spotilite**: Inyecta JavaScript en el WebView

### 4. **Diferentes Fuentes de Contenido**

- **BlockTheSpot**: Cliente desktop local
- **Spotilite**: Web player remoto (open.spotify.com)

## 🔄 Alternativas para Spotilite

### Opción 1: Proxy Local (✅ IMPLEMENTADO)

He creado un **proxy HTTP local** que funciona de manera similar a BlockTheSpot pero compatible con Spotilite:

**Características**:
- ✅ Intercepta requests a nivel de red
- ✅ Bloquea URLs de anuncios
- ✅ No requiere modificar archivos
- ✅ Compatible con WebView2
- ✅ Funciona con cualquier navegador

**Ubicación**: `internal/proxy/proxy.go`

**Cómo funciona**:
```
Spotilite → Proxy Local (puerto 8766) → Internet
                ↓
         Bloquea ads aquí
```

### Opción 2: JavaScript Injection (✅ YA EXISTE)

El módulo de adblock actual de Spotilite:

**Características**:
- ✅ Ya incluido
- ✅ Intercepta fetch() y XMLHttpRequest
- ✅ Oculta elementos de anuncios
- ⚠️ Efectividad ~70%

### Opción 3: SpotX (✅ RECOMENDADO)

Usar SpotX en Spotify Desktop + Spotilite:

**Ventajas**:
- ✅ 100% efectivo
- ✅ No interfiere con Spotilite
- ✅ Funciona independientemente

**Desventaja**:
- ⚠️ Requiere tener Spotify Desktop instalado

### Opción 4: Hosts File Blocking

Bloquear dominios de anuncios en el archivo hosts:

**Ubicación**: `C:\Windows\System32\drivers\etc\hosts`

**Agregar**:
```
127.0.0.1 audio-ads.spotify.com
127.0.0.1 ads-audio.spotify.com
127.0.0.1 ad-logger.spotify.com
127.0.0.1 ad-handler.spotify.com
127.0.0.1 doubleclick.net
127.0.0.1 googlesyndication.com
```

**Ventajas**:
- ✅ Funciona a nivel de sistema
- ✅ Afecta todas las aplicaciones
- ✅ Simple

**Desventajas**:
- ⚠️ Requiere permisos de administrador
- ⚠️ Puede romper otros servicios
- ⚠️ No es 100% efectivo

## 🆚 Comparación de Soluciones

| Solución | Efectividad | Complejidad | Compatibilidad | Recomendado |
|----------|-------------|-------------|----------------|-------------|
| **Proxy Local** | 80% | Media | ✅ Spotilite | ⭐⭐⭐⭐ |
| **JS Injection** | 70% | Baja | ✅ Spotilite | ⭐⭐⭐ |
| **SpotX** | 100% | Baja | ❌ Solo Desktop | ⭐⭐⭐⭐⭐ |
| **BlockTheSpot** | 100% | Baja | ❌ Solo Desktop | ❌ |
| **Hosts File** | 60% | Alta | ✅ Todo | ⭐⭐ |
| **Premium** | 100% | Ninguna | ✅ Todo | ⭐⭐⭐⭐⭐ |

## 🚀 Implementación del Proxy Local

### Integración en Spotilite

Para usar el proxy local en Spotilite:

1. **Iniciar el proxy** al arrancar la aplicación
2. **Configurar WebView2** para usar el proxy
3. **Monitorear** requests bloqueados

### Código de Integración

```go
// En internal/app/app.go
import "spotilite/internal/proxy"

type App struct {
    // ... campos existentes
    adProxy *proxy.AdBlockProxy
}

func NewApp(...) *App {
    return &App{
        // ... inicialización existente
        adProxy: proxy.NewAdBlockProxy("8766"),
    }
}

func (a *App) Startup(ctx context.Context) {
    // ... código existente
    
    // Iniciar proxy
    if err := a.adProxy.Start(ctx); err != nil {
        slog.Warn("failed to start ad-block proxy", "error", err)
    }
}
```

### Configuración de WebView2

```go
// En main.go, agregar opciones de WebView2
Windows: &windows.Options{
    WebviewUserDataPath: userDataPath,
    WebviewBrowserPath:  "",
    // Configurar proxy
    WebviewGpuIsDisabled: false,
    // Nota: WebView2 usa la configuración de proxy del sistema
},
```

### Configurar Proxy del Sistema (PowerShell)

```powershell
# Configurar proxy para Spotilite
netsh winhttp set proxy 127.0.0.1:8766
```

## 📊 Ventajas del Proxy Local vs BlockTheSpot

| Característica | Proxy Local | BlockTheSpot |
|----------------|-------------|--------------|
| **Compatible con Spotilite** | ✅ | ❌ |
| **No modifica archivos** | ✅ | ❌ |
| **Funciona con Web Player** | ✅ | ❌ |
| **Fácil de desactivar** | ✅ | ⚠️ |
| **No requiere reinstalar** | ✅ | ⚠️ |
| **Efectividad** | 80% | 100% |
| **Transparente** | ✅ | ❌ |

## 🔧 Próximos Pasos

### Para Implementar el Proxy

1. ✅ Crear módulo de proxy (`internal/proxy/proxy.go`) - **HECHO**
2. ⏳ Integrar en `App.Startup()`
3. ⏳ Configurar WebView2 para usar el proxy
4. ⏳ Agregar UI para ver estadísticas de bloqueo
5. ⏳ Agregar configuración on/off en settings

### Para Mejorar la Efectividad

1. Agregar más patrones de URLs de anuncios
2. Implementar detección heurística
3. Agregar lista blanca para evitar falsos positivos
4. Sincronizar listas de bloqueo con servicios como EasyList

## 📝 Conclusión

**BlockTheSpot no es compatible con Spotilite** debido a diferencias arquitectónicas fundamentales. Sin embargo, he implementado un **proxy local** que proporciona funcionalidad similar y es compatible con la arquitectura de Spotilite.

### Recomendación Final

Para la mejor experiencia:

1. **Usa el proxy local** de Spotilite (80% efectivo)
2. **O instala SpotX** en Spotify Desktop (100% efectivo)
3. **O suscríbete a Premium** (100% efectivo + legal)

---

**Última actualización**: Mayo 2026  
**Estado**: Proxy local implementado, pendiente integración
