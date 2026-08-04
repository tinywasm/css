# PLAN — Una sola tipografía para web y PDF, autoalojada

## Antes de escribir código: lee [CONSTRUCTION_HARNESS.md](CONSTRUCTION_HARNESS.md)

**Es vinculante, no orientativo.** Los principios que gobiernan este trabajo:

| # | Principio | Cómo se aplica aquí |
|---|---|---|
| 1 | Typed over `any` | La familia tipográfica es hoy un literal `string` dentro del reset. Pasa a ser un `Token` del catálogo. |
| 4 | One way to do each thing | Una sola forma de declarar la tipografía del producto, y un solo lugar donde vive el valor. |
| 9 | Lego pieces, never forks | **El más importante de este plan.** La entrega del binario de la fuente es de `assetmin`, no de `css`. Ver §3.2. |

Y la regla que descarta la solución que proponía la versión anterior:

> **Never wrap a library to fix its behaviour.** A wrapper that patches a defect is a fork
> with a friendlier name. Fix it where it lives and publish.

---

## 1. Problema

`RenderCSS()` fija la familia como literal en la regla `body` (`css.reset.go`):

```
font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
```

`system-ui` resuelve a **SF Pro en iOS/macOS** y a **Roboto en Android**: distinta altura-x,
distinto avance horizontal, el texto corta líneas en sitios distintos. Ningún reset lo
normaliza.

Además el producto emite **PDF** con `tinywasm/pdf`, y hoy usa otra tipografía distinta
(DroidSans). Web y PDF del mismo documento no se parecen.

Restricción de despliegue: los entornos objetivo **no dependen de internet**. Google Fonts y
cualquier CDN quedan descartados; la fuente viaja con la aplicación.

---

## 2. Hallazgos medidos

### 2.1 La tipografía debe satisfacer al medio más restrictivo, que es el PDF

`tinywasm/pdf` embebe vía `AddUTF8FontFromBytes`, que lee las tablas `glyf`/`loca` de
TrueType. Eso impone dos exclusiones **que deciden la elección para ambos medios**:

- **Nada de fuentes variables.** El motor ignora los ejes de variación: cargaría siempre la
  instancia por defecto, así que la negrita saldría idéntica a la regular.
- **Nada de OTF/CFF.** Verificado sobre `Inter-Regular.otf`: cabecera `OTTO`, contornos CFF,
  sin tabla `glyf` — no lo puede cargar ninguna ruta del paquete.

Por tanto: **caras estáticas TrueType, una por estilo.** Esto invalida la recomendación de
«Inter Variable» de la versión anterior de este plan, que razonaba sólo desde la web.

### 2.2 Pesos reales de los subsets

Subseteado al rango latino (ASCII + Latin-1 + comillas tipográficas + `€`) con
`pyftsubset`, conservando `kern,liga,clig`. **Un solo archivo por cara, TTF, servido a
los dos medios:**

| Roboto Regular | Bytes |
|---|---|
| TTF crudo | 27.828 B |
| TTF servido con gzip | 18.240 B |
| **TTF servido con brotli** | **16.132 B** |

### Por qué un formato y no dos

La versión anterior de este plan proponía `.woff2` para web y `.ttf` para PDF. **Es
peor en las dos dimensiones que importan**, medido:

| | Bytes que viajan | Peticiones |
|---|---|---|
| **Un TTF compartido** (brotli) | **16.132 B** | **1** |
| Un WOFF2 + un TTF | 27.890 B | 2 |

El WOFF2 pesa 13.536 B, poco menos que el TTF comprimido. Pero el motor de PDF no lo
lee, así que habría que bajar además el TTF entero: casi el doble de bytes.

Y la razón de fondo es la **caché del navegador**: el PDF se genera en el frontend, así
que pide el mismo archivo que la página ya descargó para su `@font-face`. Es un acierto
de caché, no una petición nueva. Con dos formatos se paga todo dos veces.

El navegador lee TTF sin problema (`format("truetype")`). El motor de PDF **no** lee
WOFF2. TTF es el único formato que sirve a ambos.

**DroidSans queda descartada por contenido, no por peso:** no tiene el glifo `€` —relevante
en un producto de cotizaciones— y no tiene cursiva real; los consumidores actuales apuntan
los estilos `"I"`/`"BI"` a los archivos rectos.

**Recomendación: Roboto**, y con las tablas de kerning conservadas la diferencia es
mayor de lo que parecía. Las cuatro caras, servidas con brotli:

| | Crudo | Servido |
|---|---|---|
| **Roboto** | 115.096 B | **69.978 B** |
| Inter | 209.288 B | 95.725 B |

Inter era la más liviana mientras el TTF se recortaba *sin* tablas de layout —cuando
sólo servía al PDF, que las ignora—. Al conservarlas para que el navegador aplique
kerning, sus tablas OpenType, mucho más extensas, la ponen un 37% por encima.

Además Roboto es lo que Android ya renderiza de forma nativa: esos usuarios no ven
cambio e iOS converge hacia ellos. En PDF el peso del archivo fuente es casi
irrelevante, porque el motor vuelve a subsetear al embeber (medido: un PDF de 18 KB con
85 KB de fuentes registradas).

### 2.3 `assetmin` no sabe entregar fuentes

```go
// assetmin/assetmin.go:143
func (c *AssetMin) SupportedExtensions() []string {
    return []string{".js", ".css", ".svg", ".html"}
}
```

No hay tipo de asset binario. Un `.ttf` de fuente hoy no tiene por dónde salir. **Este es el
bloqueo real del plan**, y §3.2 explica por qué no se resuelve aquí.

---

## 3. Propuesta

### 3.1 Fase 1 — La familia pasa a ser un token *(en este repo)*

Hoy es un literal, lo que contradice el principio 1 y la propia AGENTS.md §2 de este paquete
(«un valor se escribe una vez»). Una tipografía es una decisión de diseño y le toca el mismo
trato que a un color o a un espaciado:

```go
FontSans = Token{Name: "--font-sans", Dark: `system-ui, -apple-system, "Segoe UI", Roboto, sans-serif`}
```

`RenderCSS()` pasa a `fontFamily(FontSans)`. Una app cambia su tipografía con
`Theme(Set(FontSans, ...))` sin republicar este paquete, que es el mecanismo de override que
ya existe y funciona — una sola forma de hacerlo (principio 4).

La cadena de respaldo se conserva siempre, así el texto renderiza aunque el asset falle:

```
--font-sans: "Roboto", system-ui, -apple-system, sans-serif;
```

**Esta fase es autónoma:** arregla el literal y no depende de nada de lo que sigue.

### 3.2 Fase 2 — La entrega del `.ttf` es de `assetmin`, no de `css`

La versión anterior de este plan proponía que `css` hiciera `go:embed` de la fuente y
emitiera su base64 dentro del `@font-face`. **Es incorrecto bajo el harness.** `assetmin`
es la pieza que posee la entrega de assets; que no soporte fuentes es un hueco *suyo*.
Que `css` lo rodee embebiendo binarios es literalmente lo que la regla prohíbe:

> Never wrap a library to fix its behaviour. A wrapper that patches a defect is a fork with
> a friendlier name. **Fix it where it lives and publish.**

Y el corolario del documento explica por qué esto importa más de lo que parece: el hueco
aflora en la hoja —aquí, `css`—, donde el agente no tiene autoridad para publicar aguas
arriba, así que parchea localmente. La deuda no es accidental: el flujo la garantiza.

**Acción:** abrir `docs/PLAN.md` en `tinywasm/assetmin` para añadir el tipo de asset
binario (`.ttf`) al pipeline. Este plan queda **bloqueado** en su fase 2 hasta entonces.
No se implementa aquí ningún sustituto provisional.

### 3.3 Fase 3 — El PDF consume las mismas caras *(en `tinywasm/pdf`)*

Cubierto por `pdf/docs/PLAN.md`, cuya API resultante es:

```go
func LoadTypeface(regular, bold, italic, boldItalic string) (Typeface, error)
```

Cuatro parámetros distintos, así que apuntar dos estilos al mismo archivo —lo que hacen hoy
los consumidores— se ve al escribirlo.

---

## 4. La pieza que falta: `tinywasm/font`

> Corrección de una versión anterior de este plan, que respondía «no hace falta paquete».
> Esa respuesta evaluaba un contrato de datos **directo entre `css` y `pdf`** — y para eso
> sigue siendo correcta: no se importan, no se encuentran, nada cruza entre ellos. Pero la
> tipografía no viaja por ahí. Viaja por el pipeline de assets, donde el molde ya existe.

### La restricción que fija el diseño

Dos exigencias del producto, aparentemente opuestas:

1. **La tipografía nunca va embebida.** Cada proyecto declara la suya; ninguna librería
   trae una por defecto.
2. **El binario WASM sólo carga nombres, nunca bytes.** El frontend necesita *saber* qué
   cara pedir; no necesita llevarla dentro.

Se resuelven con la misma partición que el ecosistema ya usa, y que la AGENTS.md de este
paquete describe:

> *Identity strings that the frontend genuinely needs belong to `tinywasm/widget`, which is
> identity-only and WASM-safe.*

`tinywasm/widget` no tiene un solo build tag: es identidad pura y por eso cruza a WASM.
`tinywasm/css` es `//go:build !wasm` entero: valores, nunca identidad. La tipografía se
parte igual.

### El molde ya existe: `image.go`

No hay que inventar convención. `assetmin` ya declara assets binarios así:

```go
// assetmin/image_processor.go
// ImageProcessor procesa imágenes declaradas en los image.go de los módulos.
// Implementado por github.com/tinywasm/image/min; inyectado por el composition root (app).
type ImageProcessor interface {
    LoadImages() error
    ReloadModule(moduleDir string) error
    UnobservedFiles() []string
}
```

Contrato tipado en la pieza que posee el concern, implementación en librería aparte,
inyección desde la raíz de composición. Es el patrón lego del harness, ya funcionando para
un asset binario declarado en un archivo `!wasm`. **Las fuentes son el mismo caso.**

### La partición

| Pieza | Build tag | Qué contiene | Qué NO contiene |
|---|---|---|---|
| `tinywasm/font` | ninguno → **WASM-safe** | `Family`, `Face`: identidad y derivación de nombres | rutas, bytes, CSS |
| `config/fonts.go` del proyecto | **ninguno** | qué familia y en qué subcarpeta | valores CSS |
| `config/css.go` del proyecto | `//go:build !wasm` | el `Theme` con sus valores | la declaración |
| `tinywasm/assetmin` | — | contrato `FontProcessor` + el `@font-face` | la implementación |
| `tinywasm/css` | `!wasm` | `--font-sans` a partir de `font.Family` | los archivos |
| `tinywasm/pdf` | dual | pide las caras por nombre en runtime | los archivos |

El proyecto declara una vez, en `config/fonts.go`, **sin build tag**:

```go
package config

func Fonts() font.Declaration {
    return font.Declare("Roboto", "fonts/")
}
```

La omisión del tag es el diseño, no un olvido: un archivo `!wasm` no se compila para el
navegador, y el PDF se genera justo ahí — ese código necesita saber que la familia es
`"Roboto"` para pedir `fonts/Roboto-Bold.ttf`. Marcarlo lo dejaría fuera del alcance de
su propio consumidor. Una declaración es identidad, e identidad cruza.

Su vecino `config/css.go` **sí** lleva el tag, porque devuelve valores. Y como ambos son
el mismo paquete Go, `RootCSS()` llama a `Fonts()` directamente:

```go
//go:build !wasm

package config

func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.FontSans, css.FontStack(Fonts().Family())),
    )
}
```

De ahí salen las dos vistas sin duplicar la decisión: `css` toma la familia para
`--font-sans`, `pdf` deriva los nombres de cara, y `assetmin` entrega los archivos. **El
binario WASM recibe `"Roboto"` y la regla de derivación — nunca un byte de fuente.**

**`tinywasm/ssr` no participa.** Se evaluó añadirle un productor `RenderFonts()` y se
descartó al ver lo anterior: el valor viaja dentro de `RootCSS()`, que `ssr` ya extrae.

### Por qué esto sí es un lego y no ceremonia

La versión anterior descartaba el paquete porque sería «datos sin comportamiento importado
por dos librerías que nunca interactúan». Aquí no es así: `font` posee una
responsabilidad —la identidad de la tipografía y la derivación de sus caras— y expone un
contrato tipado que **cuatro** piezas consumen por caminos distintos. Y sustituye al
`string` suelto que hoy nombra familias en ambos extremos, que es el mismo agujero
diagnosticado en `pdf/docs/PLAN.md` (principio 1).

El drift deja de ser evitable-por-disciplina y pasa a ser irrepresentable: no hay dos
lugares donde escribir el nombre.

### Planes que esto abre

| Repo | Alcance | Estado |
|---|---|---|
| `tinywasm/font` | `Family`, `Style`, `Face`, `Declaration`. WASM-safe, sin mapas. | ✅ **v0.0.3 publicado** |
| `tinywasm/pdf` | Harness de tipografía tipado | ✅ **v0.1.0 publicado** |
| `tinywasm/pdf` | `LoadDeclared(font.Declaration)` — el puente que aún falta | `pdf/docs/PLAN.md` |
| `tinywasm/css` | Fase 1 de este plan, con `FontSans` alimentado por `font.Family`. | este documento |
| `tinywasm/assetmin` | Contrato `FontProcessor` calcado de `ImageProcessor` + el `@font-face` con sus URLs. | `assetmin/docs/PLAN.md` |

Orden obligado: `font` ya está —los demás lo importan—, luego `css` y `pdf` en paralelo,
y `assetmin` al final, cuando ambos consumidores estén definidos.

El estado del conjunto vive en `app-releases/docs/TYPOGRAPHY_MASTER_PLAN.md`.

---

## 5. La otra pieza: `tinywasm/color`

Este plan trata de tipografía, pero el mismo patrón —una decisión sin pieza que la
posea— aparece aquí con el color, y el trabajo cae en **este** repositorio.

### 5.1 La ciencia del color vive dentro de un test

`contrast_test.go` declara tres funciones que no son de test:

| Función | Qué es |
|---|---|
| `parseHex` | parser de hex con `strconv.ParseUint`, tres veces |
| `relativeLuminance` | luminancia relativa WCAG, con la curva sRGB |
| `contrastRatio` | ratio de contraste WCAG |

Ninguna es específica de CSS. Es ciencia del color de propósito general, atrapada donde
nadie más puede usarla — y `pdf/color.go` tiene su propia copia del parser. Dos
implementaciones del mismo cálculo (principio 4).

**Corrección:** `github.com/tinywasm/color` ya existe y posee el concern. Este repo
importa `Luminance` y `Contrast` de allí y borra sus copias.

**El sentido de la dependencia es correcto:** este paquete es `//go:build !wasm` y
`color` no tiene tag. `css → color` compila; al revés sería imposible.

### 5.2 La carencia de la auditoría, y algo que no estaba documentado

`tokens.go:49` dice:

```go
// AllPairs returns the 7 functional design decision pairs for automated contrast auditing.
// SurfaceSunken and SurfaceSelected are excluded: their values are color-mix()
// expressions that resolveColor() cannot evaluate (known gap, documented in SPECS).
```

Dos cosas:

1. **El comentario dice 7 y la lista tiene 5.** No cuadra ni descontando las dos
   exclusiones que nombra.
2. **`SurfaceAccent` también está fuera, y el comentario no lo menciona.** Sus dos
   colores son hex literal (`#e8a33d` / `#1C1C1E`), así que `resolveColor` los resuelve
   sin problema: la exclusión no tiene la justificación que sí tienen las otras dos.
   Calculado, su ratio es **7,89** — pasa el mínimo de 4,5 con holgura. O sea que no se
   excluyó por fallar, sino por descuido, y hoy **nadie vigila que siga pasando** si
   alguien retoca el acento.

**Corrección inmediata:** añadir `SurfaceAccent` a `AllPairs()` y corregir el comentario
para que el número coincida con la lista.

**Corrección de fondo:** `SurfaceSunken` y `SurfaceSelected` sólo entran cuando exista
`Mix` en `tinywasm/color` — evaluar `color-mix(in oklab, …)` es trabajo de una librería
de color, no de una de CSS. Está en `color/docs/PLAN.md` §3.3. Mientras tanto la
exclusión sigue, pero **documentada como pendiente con dueño**, no como límite
permanente.

---

## 6. Lo que no resuelve

- **Controles nativos.** El desplegable de `<select>`, los date pickers y el input de
  archivo son chrome del sistema operativo.
- **Rasterizado.** El antialiasing subpíxel difiere entre motores.
- **Física de scroll.** El rubber-banding de iOS no tiene equivalente en Android.

---

## 7. Verificación

1. `RootCSS().String()` contiene `--font-sans`.
2. `RenderCSS().String()` usa `var(--font-sans` en `body`, sin literal `system-ui` suelto.
3. `Theme(Set(FontSans, "X"))` emite el override en el `:root` de cola.
4. `--font-sans` añadido a `allTokens` en `css_test.go` — lo exige
   `TestNoUndeclaredTokensInEmittedCSS`.
5. `docs/SPECS.md` actualizado con el token nuevo.

Del color (§5):

6. `grep -n "ParseUint\|math.Pow" contrast_test.go` no encuentra nada: la luminancia y el
   ratio vienen de `tinywasm/color`.
7. `AllPairs()` incluye `SurfaceAccent`, y el comentario que la precede declara un número
   que coincide con la lista.
8. La exclusión de `SurfaceSunken` y `SurfaceSelected` sigue documentada, pero como
   **pendiente con dueño** (`color/docs/PLAN.md` §3.3), no como límite del ecosistema.
9. `gotest`.

La fase 2 no se da por verificada hasta que `assetmin` sirva el `.ttf` y una página real
lo cargue sin petición de red externa.

---

## 8. Decisiones pendientes

1. ¿Se confirma Roboto, o se prefiere Inter por su dibujo pese a los 13 KB extra por visita?
2. ¿Se abre ya el `PLAN.md` de `assetmin` para el asset binario, o la fase 1 va sola primero?
3. ¿El test anti-drift vive en la app, o `assetmin` debería exponerlo al registrar la fuente?
