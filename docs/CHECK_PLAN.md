# PLAN — `tinywasm/css`: entrypoint tipado de tema (`Theme`)

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Gate 1 del [PLAN_TEMA_Y_CAPACIDADES](../../docs/PLAN_TEMA_Y_CAPACIDADES.md). Bloquea a `layout` y
> `mjosefa-cms` (`components` es verificación, no depende de este gate). Alinéate con
> `tinywasm/docs/ARNES_DE_CONSTRUCCION.md`.

## Contexto (estado actual — ya implementado, NO reescribir)

`tinywasm/css` ya es dueño del catálogo tipado de tokens y emite el `:root`. Lo relevante:

- `tokens.go` — `type Token struct{ Name, Fallback string }` con `Var() string` →
  `"var(--x,fallback)"`. Catálogo completo: `ColorPrimary/Secondary/Success/Error`,
  capa activa `ColorBackground/Surface/OnSurface/Muted/Hover`, capa fuente
  `Color*Light/*Dark`, escalas `Text*`, `Space*`, `Radius*`, `Shadow*`, motion, z-index, etc.
- `css.go`:
  - `RootCSS() *Stylesheet` — declara marca + capa fuente (Light/Dark) + escalas en `:root`.
  - `RenderCSS() *Stylesheet` — reset + `body` + **enlaza la capa activa** con
    `Bind(ColorSurface, ColorSurfaceLight)` (etc.) en `:root`, y un `MediaPrefersDark` que
    reenlaza a las variantes Dark.
- `dsl.go` — `Root(decls ...Decl)`, `Declare(Token, value) Decl` (asumido; ver más abajo),
  `Bind(active, source Token) Decl`, `NewStylesheet(items ...item)`, `Token.cssValue()`.

**Mecanismo de rebrand que ya existe pero NO está formalizado:** una app puede re-`Declare`
los tokens fuente (`ColorSurfaceLight`, `ColorSecondary`, …) en un `:root` posterior. El
comentario en `tokens.go` lo insinúa ("apps redeclare these for rebrand"), pero **no hay un
entrypoint tipado único** que lo haga, así que cada consumidor improvisa (p.ej. `platformd`
inventó su propia paleta `--pd-*`). Eso es el hueco del arnés que cierra este plan.

## Objetivo

Exponer **una sola forma tipada** de que una app aplique su tema (override de tokens), de modo
que:

- La app declare su paleta **una vez**, referenciando símbolos `Token` (no strings de nombres
  de variables), y el override gane sobre `RootCSS`/`RenderCSS`.
- Ningún consumidor (layout/components) necesite declarar tokens propios ni hex.
- Un override mal escrito **no compile** (se exige `Token`, no un nombre libre).

## Diseño

### 0. Cómo se entrega el tema — el pipeline SSR (verificado, condiciona el diseño)

`tinywasm/ssr` descubre por regex los `RootCSS()`/`RenderCSS()` del grafo y `tinywasm/assetmin`
los enruta a `web/public/style.css`. **Hecho crítico** (`assetmin/ssr_loader.go`): el bloque
`:root` (lo que produce `RootCSS()`) es un **slot de un solo ganador por REEMPLAZO**:

- Solo el **proyecto raíz** (la app) o el módulo `tinywasm/css` pueden declarar `RootCSS()`;
  cualquier otro módulo que lo declare se **ignora con warning**.
- Si la app declara su propio `RootCSS()`, **reemplaza por completo** el de `tinywasm/css`
  (`if fromRoot != nil { … } else if fromCss != nil { … }`).
- `RenderCSS()` es aparte: se concatena de todos los módulos (reset + `Bind` de la capa activa
  de css + platformd + rightpanel + components). No es de un solo ganador.

**Consecuencia de diseño:** una app rebranda declarando `func RootCSS() *css.Stylesheet` en su
raíz; como **reemplaza**, ese `RootCSS()` debe emitir el **catálogo completo** (no solo los
overrides), o se perderían Space*/Radius*/tipografía/etc. Por eso `Theme()` **no** puede devolver
solo el bloque de overrides: debe devolver el catálogo completo con los overrides aplicados al
final. Este plan **no** cambia `ssr` ni `assetmin` (el mecanismo ya existe); solo provee el
helper tipado con la forma correcta.

### 1. `func Theme(overrides ...Override) *Stylesheet`  (nuevo, en `css.go`)

Entrypoint canónico de rebrand. Devuelve el **catálogo `:root` completo** (equivalente a
`RootCSS()`) con los `overrides` aplicados al final del mismo stylesheet (el último `:root` gana
en la cascada dentro del bloque). Así es un **reemplazo válido** que la app retorna desde su
propio `RootCSS()` sin perder tokens.

**El parámetro es un tipo dedicado `Override`, NO `...Decl`.** `Declare` devuelve `Decl`, el mismo
tipo que `Color`, `Padding`, `PaddingBottom`, etc. Si `Theme` recibiera `...Decl`,
`Theme(Padding(Px(4)))` compilaría y produciría `:root{padding:4px}` — estado ilegal
representable (#3/#6). Un `Override` solo se construye con `Set(Token, value)`, así que en `Theme`
solo entran overrides de token del catálogo:

```go
// Override es el cambio de valor de UN token. Campos no exportados: solo Set lo construye.
type Override struct {
    token Token
    value string
}

// Set declara el override de un token del catálogo. Token tipado (no un nombre libre);
// value es el borde de I/O.
func Set(t Token, value string) Override { return Override{t, value} }

// Theme devuelve el catálogo :root COMPLETO (como RootCSS) con los overrides al final.
// Pensado como el RootCSS() del proyecto raíz — assetmin REEMPLAZA el :root de css por el de
// la app, por eso trae el catálogo entero, no solo los overrides.
func Theme(overrides ...Override) *Stylesheet {
    root := RootCSS() // catálogo por defecto
    decls := make([]Decl, len(overrides))
    for i, o := range overrides {
        decls[i] = Declare(o.token, o.value)
    }
    return withRootTail(root, Root(decls...)) // helper interno: añade un :root final; ver Etapa 1
}
```

Uso — la app lo expone como su `RootCSS()` de raíz (documéntalo en README):

```go
// css.go del PROYECTO RAÍZ de la app  //go:build !wasm
func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.ColorSecondary, "#3f88bf"),
        css.Set(css.ColorOnSecondary, "#ffffff"),   // blanco = texto sobre el acento
        css.Set(css.ColorSurfaceLight, "#e9e9e9"),
        css.Set(css.ColorBackgroundLight, "#e9e9e9"),
    )
}
```

**Por qué no un struct con un campo por token** (`ThemeSpec{Secondary, Surface, …}`):
- Duplicaría el catálogo (`tokens.go` es la fuente única de qué tokens existen); cada token nuevo
  obligaría a añadir un campo → dos listas que derivan (rompe #4 y #1 "reutiliza tipos ya declarados").
- En la práctica solo cubriría colores; cerraría el override de `Space*`/`Radius*`/`Text*`.
- Override parcial exigiría distinguir "no seteado" de "vacío" (punteros/sentinelas) → footgun.
- El valor seguiría siendo `string`; no gana type-safety y pierde el reuso del `Token` en la clave.

`...Override` da "exactamente lo requerido" (solo overrides de token válidos) **reutilizando** el
catálogo, sin inventar una segunda lista.

### 2. Garantizar override en un único punto

Verifica (y ajusta si hace falta) que sobre-escribir un token **fuente** (p.ej.
`ColorSurfaceLight`) se propague a la capa activa (`--color-surface`) por el `Bind` existente,
y que sobre-escribir un token **activo** directamente (`ColorSurface`) también funcione. El
objetivo: la app elige rebrand "por fuente" (respeta light/dark) o "por activo" (fija un color),
ambos vía `Declare`. No dupliques el `Bind`; solo confirma con un test que el override gana.

### 3. Sin `Theme()` mágico de nombres — ni `Decl` arbitrario

`Theme` NO acepta `map[string]string`, `...any`, ni nombres de variable como string. Tampoco
`...Decl` (aceptaría `Padding`/`Color`/cualquier propiedad dentro de `:root`). Solo `...Override`,
construidos con `Set(Token, value)`: el **token** es tipado; el `value` es el borde de I/O. No
cambies la firma de `Declare` (se sigue usando internamente para emitir el `:root`).

## Restricciones del arnés (obligatorias)

- Sin stdlib en el paquete: usa `tinywasm/fmt` (nunca `fmt`, `strings`, `strconv`).
- Cero `any`/`interface{}`/`...any` en la API nueva. `Theme` recibe `...Override` (solo
  construible con `Set(Token, value)`), no `...Decl` (que aceptaría cualquier propiedad).
- Superficie mínima: exporta `Theme`, `Override` y `Set`. Los campos de `Override` son no
  exportados (solo `Set` lo construye); helpers de composición (`withRootTail`) sin exportar.
- No hardcodear nombres de variables CSS como string literal en lógica: usa los símbolos `Token`.
- Una sola forma: `Theme` es EL camino de rebrand. No agregues variantes.

## Tests (`gotest`, no `go test`)

Añade a `css_test.go` (assertions de stdlib, dual WASM/stdlib):

1. `Theme(Set(ColorSecondary,"#3f88bf"))` contiene el **catálogo completo** (aparecen p.ej.
   `--space-2`, `--radius-md`, `--text-xl`, la capa fuente `--color-surface-light`) Y, más abajo,
   `--color-secondary:#3f88bf` como último valor de ese token. Es un reemplazo válido de
   `RootCSS()` — NO solo el bloque de overrides.
2. El valor efectivo de `--color-secondary` en `Theme(...)` es el override (última aparición), no
   el default `#654FF0`. Verifica por orden de aparición del string.
3. `Theme()` sin overrides es equivalente en tokens a `RootCSS()` (no pierde ninguno) — así una
   app que solo declara `RootCSS(){ return Theme() }` no rompe nada.

## Documentación (antes de cerrar)

- `README.md`: sección "Theming una app" con la tabla "quiero X → uso Y" y el ejemplo: la app
  expone `func RootCSS() *css.Stylesheet { return css.Theme(Declare(...)) }` en su raíz. Explica
  que assetmin **reemplaza** el `:root` de css por el de la app (por eso `Theme` trae el catálogo
  completo), y que `RenderCSS()` de css sigue aportando el `Bind` de la capa activa.
- `docs/ARCHITECTURE.md` (o `docs/API.md` si existe): documentar `Theme` como el entrypoint
  único de rebrand y su relación con el slot `RootCSS` de `assetmin` (un solo ganador, reemplazo).
- Re-indexar `README.md` para que enlace todo `docs/`.

## Etapas

| Etapa | Entregable | Criterio de hecho |
|-------|-----------|-------------------|
| 1 | `type Override` + `func Set(Token,string) Override` + `func Theme(...Override) *Stylesheet` (catálogo completo + overrides) + helper interno `withRootTail`, en `css.go` | Compila; solo acepta `Override`; `Theme()` ≡ `RootCSS()` en tokens |
| 2 | Confirmar que el override gana (última aparición) y que no se pierde ningún token | Test 1–3 verdes |
| 3 | Tests en `css_test.go` | `gotest` verde en stdlib y WASM |
| 4 | Docs (README + ARCHITECTURE/API) — incl. relación con slot RootCSS de assetmin | `Theme` documentado como forma única |

## Referencia (código actual, para reciclar — NO reescribir lo que ya existe)

`css.go` capa activa (ya existe, base sobre la que Theme debe ganar):

```go
Root(
    Bind(ColorBackground, ColorBackgroundLight),
    Bind(ColorSurface, ColorSurfaceLight),
    Bind(ColorOnSurface, ColorOnSurfaceLight),
    Bind(ColorMuted, ColorMutedLight),
    Bind(ColorHover, ColorHoverLight),
),
MediaPrefersDark(Root(
    Bind(ColorSurface, ColorSurfaceDark), /* ... */
)),
```
