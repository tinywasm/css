---
PLAN: "feat: tokens de safe area y altura de viewport dinámica"
TAG: v0.5.1
---

## Antes de escribir código: lee [CONSTRUCTION_HARNESS.md](CONSTRUCTION_HARNESS.md)

**Es vinculante, no orientativo.**

| # | Principio | Cómo se aplica aquí |
|---|---|---|
| 4 | Magnitudes son decisiones de diseño (AGENTS §4) | El inset del notch y la altura del viewport son magnitudes: token en el catálogo, referenciado por `var()`. Nunca un `env(...)` suelto dentro de una regla de componente |
| 9 | Lego pieces, never forks | Esta librería posee el vocabulario CSS. `assetmin` y `devbrowser` no declaran variables; consumen éstas |
| 5 | Minimal surface | Cinco tokens y ninguna función nueva. La `@media` para safe areas no existe: `env()` ya devuelve `0px` donde no aplica |
| 2 | Explicit over implicit | `--safe-top` dice qué es. `env(safe-area-inset-top, 0px)` repetido en cada componente, no |

---

## 1. El hueco

Este catálogo ya arrastra una historia con Safari iOS. El plan anterior
(`LAST_PLAN_EXECUTED.md`, v0.5.0) nombró el síntoma con estas palabras:

> *la app no se ve igual en Safari iOS que en Chrome Android*

y cerró la parte tipográfica con `FontFaces`. Este plan cierra la siguiente capa
del mismo síntoma: **la geometría del dispositivo**.

Hoy no existe forma de expresar dos magnitudes que sólo un iPhone conoce:

| Magnitud | Estado | Qué se rompe sin ella |
|---|---|---|
| Los cuatro insets de la pantalla (notch, Dynamic Island, barra de gestos) | **no existe** ningún token. `grep 'env(' *.go` → 0 resultados | Un header pegado arriba queda bajo la Dynamic Island; un botón abajo, bajo la barra de gestos. En el escritorio no se nota |
| Altura del viewport **visible** | **no existe**. Un componente sólo puede escribir `100vh` | En Safari iOS la barra de URL se contrae al hacer scroll: `100vh` mide el viewport *largo*, así que el contenido se corta por abajo de forma permanente |

Lo que el reset **sí** cubre ya, y no hay que volver a tocar: `css.reset.go` trae
`-webkit-text-size-adjust`, `-webkit-tap-highlight-color`, el `appearance` de
botones e inputs y el radio inset de los campos de texto — todos con su
justificación escrita, todos apuntando a Safari iOS. La familia del arreglo ya
vive aquí; le faltan estas dos piezas.

### 1.1 Precondición externa

`env(safe-area-inset-*)` devuelve **`0px` en todos los dispositivos** salvo que el
documento se sirva con `<meta name="viewport" content="…, viewport-fit=cover">`.
Esa etiqueta no es CSS y no se emite aquí: es de `tinywasm/html` (shell de
`Document()`) y de `tinywasm/assetmin` (`index.html` del bundle), cada uno con su
propio plan.

Consecuencia operativa: los tokens de §2.1 son **correctos y verificables por
test** desde el primer día, pero **no producen efecto visible en un dispositivo**
hasta que aquellos dos planes estén ejecutados. No es motivo para retrasar éste
—el orden natural es publicar el vocabulario antes que su consumo— pero sí para
no dar por fallido el resultado si se prueba antes de tiempo.

---

## 2. Cambio

### 2.1 Cuatro tokens de safe area

En `catalog.go`, grupo nuevo, con la nomenclatura de prefijos existente:

```go
SafeTop    = Token{Name: "--safe-top", Dark: "env(safe-area-inset-top, 0px)"}
SafeRight  = Token{Name: "--safe-right", Dark: "env(safe-area-inset-right, 0px)"}
SafeBottom = Token{Name: "--safe-bottom", Dark: "env(safe-area-inset-bottom, 0px)"}
SafeLeft   = Token{Name: "--safe-left", Dark: "env(safe-area-inset-left, 0px)"}
```

Notas de conformidad con AGENTS.md:

- **`Dark` a solas = valor estático** (§7). Un inset no depende del tema; no
  lleva par `Light`/`Dark`.
- El fallback `0px` del `env()` **no** es duplicación de valor (§2): es el valor
  que la especificación define para un dispositivo sin insets, y el `env()` sin
  segundo argumento en un navegador que no conoce la variable invalida la
  declaración entera. Es la diferencia entre "indefinido" y "mal definido" que la
  propia AGENTS §4 señala.
- Un componente **nunca** escribe `env(safe-area-inset-top)`. Escribe
  `var(--safe-top)`. Ésa es la razón de que existan: hoy la única forma de usar
  un inset sería un valor crudo dentro de una regla, que es exactamente lo que
  este catálogo existe para cerrar.

### 2.2 Un token de altura de viewport

```go
ViewportH = Token{Name: "--viewport-h", Dark: "100dvh"}
```

`dvh` (*dynamic viewport height*) mide el viewport **visible ahora**: encoge y
crece con la barra de URL de Safari iOS, que es justo lo que `vh` no hace.

**Sobre el soporte de navegador, y por qué no lleva `@supports`:** este catálogo
ya emite `light-dark()` en las parejas de tema (AGENTS §7) y `color-mix()` en los
tokens compuestos (`ColorSurfaceSunken`, `catalog.go:24`). `light-dark()` exige
Chrome 123 y Safari 17.5 (2024); `dvh` está disponible desde Safari 15.4 y
Chrome 108 (2022). **`dvh` es estrictamente más antiguo que la línea base que
esta librería ya requiere**, así que un bloque `@supports` no protegería a ningún
navegador que pudiera renderizar el resto del stylesheet: sería código muerto con
apariencia de rigor. No se añade, y ésta es la justificación por escrito.

*(Corolario: no se añade ningún emisor `supports()` al DSL. Si algún día hace
falta, será por otra razón y con su propio análisis.)*

### 2.3 Declaración y registro

Los cinco tokens entran en `css.default.go`, en `defaultRoots()`, como un grupo
nuevo al final —no son identidad de marca, así que no van en `brandRoot()`:

```go
root(
    // Device geometry — insets reported by the device, and the viewport
    // height that shrinks with Safari iOS' collapsing URL bar.
    declare(SafeTop),
    declare(SafeRight),
    declare(SafeBottom),
    declare(SafeLeft),
    declare(ViewportH),
),
```

Y en `allTokens` de `css_test.go:210`, que es lo que hace pasar
`TestNoUndeclaredTokensInEmittedCSS`.

### 2.4 Lo que este plan **no** añade al reset

Ninguna regla nueva en `resetRules()`. En particular, **no** se añade la guarda
anti-zoom de iOS, y la razón es una limitación real de la cascada que conviene
dejar documentada para que nadie la reintente:

iOS Safari hace zoom automático al enfocar un `input` con `font-size` computado
menor que 16px. La tentación es una regla como
`@media (pointer: coarse) { input { font-size: max(1rem, 1em) } }` en el reset.
**No funcionaría.** `RenderCSS()` emite todo el reset dentro de `@layer tokens`
(`css.go:22`), la más baja de las cuatro capas; cualquier `@layer widgets { .field { font-size: var(--text-sm) } }`
la gana por orden de capa, sin importar la especificidad. Es decir: la guarda
sería inofensiva justo en el caso que pretende cubrir, y redundante en el otro
—por defecto un control hereda `font: inherit` (`css.reset.go`), que resuelve a
`TextBase` = `1rem` = 16px, que ya es seguro.

El invariante correcto, entonces, no es una regla CSS sino una **restricción de
uso**, y se escribe en `docs/SPECS.md`:

> Un control de formulario no debe llevar un tamaño de texto por debajo de
> `TextBase` en dispositivos táctiles: iOS Safari hace zoom al enfocarlo.

Su punto de aplicación está en quien escribe esas reglas (`tinywasm/widget/style`),
no aquí, y su detección en tiempo de ejecución la aporta `browser_audit_mobile`
de `tinywasm/devbrowser`. Anotado para que la restricción tenga dueño y no se
quede en una frase suelta.

---

## 3. Verificación

Aserciones de stdlib, `gotest`, sin testify:

1. `TestNoUndeclaredTokensInEmittedCSS` pasa con los cinco tokens en `allTokens`
   (es la guarda que ya existe; el test nuevo es no romperla).
2. `RootCSS().String()` contiene las cinco declaraciones, cada una con su
   `env(...)` y su fallback `0px`.
3. Guarda de fallback: ningún `env(` emitido carece de segundo argumento. Un
   `env()` sin fallback invalida la declaración en cualquier navegador que no
   conozca la variable, y es un fallo silencioso — exactamente lo que el
   principio 6 prohíbe.
4. `ViewportH.Var()` devuelve `var(--viewport-h)`, consumible desde el DSL de
   componentes igual que cualquier otro token.

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
gotest
```

---

## 4. Alcance

### Dentro

- `catalog.go`: los cinco literales.
- `css.default.go`: el grupo `root(...)` nuevo.
- `css_test.go`: `allTokens` + los tests de §3.
- `docs/SPECS.md`: valores y API de los cinco tokens, más el invariante de §2.4.
- `docs/MIGRATION.md`: no aplica — el cambio es puramente aditivo, no sustituye
  ni elimina ningún token.

### Fuera

- Aplicar los tokens a algún selector (`body { padding-top: var(--safe-top) }`).
  `RootCSS()` es **vocabulario**; dónde se aplica el inset es una decisión de
  layout, y su dueño es `tinywasm/layout` / `tinywasm/widget/style`.
- La etiqueta `<meta name="viewport">`: no es CSS (§1.1).
- Cualquier regla nueva en `resetRules()` (§2.4).
- Un emisor `supports()` en el DSL (§2.2).
