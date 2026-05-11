# Justificación del DSL tipado de CSS en Go

> Documento de análisis. Responde a la pregunta: **¿es esta API la forma más intuitiva, legible y profesional de escribir CSS en Go para el ecosistema tinywasm?**

## 1. La forma propuesta

```go
//go:build !wasm
package button

import . "github.com/tinywasm/css"

var (
    ClsBtn     Class = "btn"
    ClsPrimary Class = "btn-primary"
)

func (b *Button) RenderCSS() *Stylesheet {
    return New(
        Rule(ClsBtn,
            Padding(Rem(0.5), Rem(1)),
            BorderRadius(RadiusSm),
            Cursor(Pointer),
            FontSize(TextBase),
        ),
        Rule(ClsPrimary,
            Background(ColorPrimary),
            Color(ColorOnPrimary),
        ),
        Rule(ClsPrimary.Hover(),
            Opacity(0.9),
        ),
    )
}
```

## 2. Criterios de evaluación

Para responder con honestidad hay que fijar qué significa "mejor". Estos son los tres criterios que el usuario enunció:

| Criterio | Definición operativa |
|---|---|
| **Intuitivo** | Un desarrollador que sabe CSS reconoce la intención en segundos sin leer documentación. |
| **Legible** | Un lector que no escribió el código puede volver semanas después y reconstruir el modelo mental sin esfuerzo. |
| **Profesional** | Coherente con prácticas establecidas en frameworks consolidados; soporta refactor, tests, herramientas de IDE. |

A esos tres añado dos criterios técnicos no negociables del ecosistema tinywasm, porque son los que descalifican a varias alternativas teóricamente válidas:

| Criterio técnico | Razón |
|---|---|
| **Cero CSS en el binario WASM** | El framework optimiza tamaño con TinyGo. Cualquier solución que arrastre código de generación de CSS al frontend está descalificada. |
| **Sin generadores** | `go generate` y herramientas externas añaden deuda; el ecosistema ya rechazó esa vía. |

---

## 3. Evaluación frente a los criterios

### 3.1 Intuitivo — ¿se entiende sin documentación?

**Sí, con dos asunciones razonables.**

| Línea en el DSL | Equivalente CSS evidente |
|---|---|
| `Rule(ClsPrimary, Background(ColorPrimary))` | `.btn-primary { background: var(--color-primary); }` |
| `Rule(ClsPrimary.Hover(), Opacity(0.9))` | `.btn-primary:hover { opacity: 0.9; }` |
| `Padding(Rem(0.5), Rem(1))` | `padding: 0.5rem 1rem;` |

La única traducción no trivial es `Cls<X>` ↔ selector de clase. Una vez aprendida (un párrafo de README) el mapeo es 1:1. Comparado con TypeScript-CSS-in-JS (donde hay que aprender camelCase, units como strings vs numbers, theme objects, variants), la barrera es claramente menor.

### 3.2 Legible — ¿sobrevive seis meses después?

**Mejor que el CSS string actual.** Tres razones:

1. **Los tokens son nombres, no valores duplicados.** Hoy en `button.css` se lee `var(--color-primary, #00ADD8)` — el lector tiene que decidir si el `#00ADD8` es un fallback intencional o basura desactualizada. Con el DSL se lee `Background(ColorPrimary)`: cero ambigüedad.
2. **Los selectores son referencias, no strings.** `ClsPrimary.Hover()` impide el typo silencioso `.btn-primry:hover {}` que hoy es invisible hasta abrir el navegador.
3. **El IDE colabora.** "Find references" sobre `ColorPrimary` lista todos los usos en el repo. Sobre `--color-primary` en strings, el IDE da resultados parciales y mezclados con falsos positivos.

**Pierde frente al CSS crudo en un punto:** densidad visual. `padding: 0.5rem 1rem;` ocupa menos pixeles que `Padding(Rem(0.5), Rem(1))`. Esto se mitiga con el dot-import (sin él sería peor) pero no se elimina. Es el coste honesto.

### 3.3 Profesional — ¿se sostiene con el rigor de la industria?

**Sí, con precedente directo.** El patrón (DSL tipado + tokens como constantes + clases como identificadores tipados + extracción estática a CSS) es exactamente lo que hacen:

| Sistema | Lenguaje | Patrón equivalente |
|---|---|---|
| **vanilla-extract** | TypeScript | `style({ background: vars.color.primary })` → CSS extraído en build |
| **Linaria** | TypeScript | Tagged templates con extracción estática |
| **Stitches** | TypeScript | `styled('button', { variants: {...} })`, tokens tipados |
| **JetBrains Compose HTML** | Kotlin | DSL builder de reglas, tokens tipados |
| **ScalaCSS** | Scala | DSL puro Scala, class names generadas |
| **W3C Design Tokens CG** | (especificación) | Estandariza el concepto de "token" como entidad tipada |

No es invención local. Es la convergencia industrial de la última década aplicada a Go. La diferencia: TypeScript necesita un compilador adicional (vanilla-extract usa Babel/esbuild plugin); en Go basta `//go:build !wasm` para que el código exista solo en el servidor.

### 3.4 Cero CSS en el binario WASM

**Garantizado por construcción.** El layout del paquete `tinywasm/css`:

| Archivo | Build tag | Compila a WASM |
|---|---|---|
| `tokens.go` (Class, Token, constantes) | ninguno | ✅ |
| `dsl.go` (Stylesheet, Rule, propiedades) | `!wasm` | ❌ |
| `ssr.go` (RootCSS, RenderCSS) | `!wasm` | ❌ |

Lo único que cruza al binario WASM son las strings con nombres de clase (`"btn-primary"`) que el HTML necesita emitir. TinyGo además elimina por dead-code los tokens no referenciados. El generador de CSS no existe en el frontend.

### 3.5 Sin generadores

**Cumplido.** No hay `go generate`, no hay `theme.css`, no hay paso de build adicional. El compilador Go es la única herramienta.

---

## 4. Alternativas evaluadas y descartadas

| Alternativa | Razón de descarte |
|---|---|
| **Mantener `.css` + `//go:embed`** (estado actual) | Stringly-typed; ningún error se detecta hasta abrir el navegador; renombrar tokens es manual y propenso a drift. |
| **`.css` + linter externo** | Resuelve detección de typos pero añade herramienta externa; sigue siendo dos lenguajes. |
| **Generador `theme.css → tokens.go`** | Mantiene dos representaciones; el generador es deuda permanente; contradice "sin generadores". |
| **CSS-in-Go runtime estilo styled-components** | Arrastra el motor CSS al binario WASM. Descalifica de inmediato. |
| **Templates de texto (`text/template`)** | Devuelve a strings sin tipo; pierde validación del compilador. |
| **DSL fluido (builder con `.Padding(...).Color(...)`)** | Equivalente expresivo al constructor variádico, peor para reglas largas (encadenamiento vertical incómodo); el variádico aplana mejor. |
| **Sub-paquete `tinywasm/css/cssgo`** | Obliga dos imports, rompe el dot-import. Sin valor real. |

---

## 5. Riesgos honestos (no se ocultan)

Un análisis profesional debe nombrar lo que el patrón pierde o complica:

1. **Verbosidad relativa al CSS crudo.** `Padding(Rem(0.5), Rem(1))` es más caracteres que `padding: 0.5rem 1rem`. Mitigación: dot-import; los autores aprenden a escanear visualmente la estructura `Rule(sel, ...decls)`.

2. **`@media` y selectores raros pasan por `Selector("...")` o `Media("...")`.** Para `@container`, attribute selectors complejos, `:nth-child(...)`, la API se vuelve un escape hatch a string. No es elegante pero es honesto: cubrir 100% de CSS spec en Go tipado es esfuerzo desproporcionado. El DSL prioriza el 90% común.

3. **`ssr.go` puede crecer.** Un componente con 200 líneas de CSS se convierte en 200 líneas de Go en `ssr.go`. Es el mismo volumen, no más. Si duele en algún componente puntual, se permite un `ssr_styles.go` en el mismo package como excepción.

4. **Curva de adopción inicial.** Un colaborador que solo conoce CSS clásico necesita media hora para internalizar el mapeo. Coste único; el beneficio es permanente.

5. **No hay tooling de formato CSS-aware.** `gofmt` no sabe alinear declaraciones como `prettier` alinea CSS. Mitigación: el constructor variádico produce naturalmente una declaración por línea; el formateo es predecible.

### 5.1 Sobre la verbosidad: ¿unidades variádicas?

Pregunta natural: si `Padding(Rem(0.5), Rem(1))` es más largo que `padding: 0.5rem 1rem`, ¿por qué no hacer `Rem` variádico → `Padding(Rem(0.5, 1))`?

**Descartado.** Razones, en orden de peso:

1. **Muddle de tipos.** `Rem` es *una unidad*, no una lista. Hacerlo variádico convierte `Value` en "valor atómico o secuencia serializada", contaminando todo el sistema. Habilita basura compilable como `Color(Rem(0.5, 1))` → `color: 0.5rem 1rem`.
2. **Ahorro despreciable.** ~5 caracteres × ~15% de declaraciones shorthand = ~75 chars en todo el proyecto.
3. **El DSL ya resuelve shorthand en el lugar correcto**: la propiedad (`Padding`) es variádica, igual que la gramática CSS (`padding: <length>{1,4}`). Mover variadicidad a la unidad traslada la responsabilidad al lugar equivocado.
4. **Bloquea unidades mixtas reales**: `BoxShadow(Em(0.1), Em(0.1), Em(0.2), ColorSurface)` mezcla `em` con un token — imposible si `Em` o `Rem` se vuelven variádicos.

Si el equipo decide en el futuro que la densidad importa más que la pureza de tipos, el camino correcto sería *floats directos con unidad implícita por propiedad* (`Padding(0.5, 1)` → rem), no unidades variádicas. Análisis aparte.

---

## 6. ¿Es *la mejor* o solo *una buena*?

Aquí hay que distinguir dos preguntas:

### 6.1 ¿Es la mejor forma de escribir CSS en un lenguaje tipado?
Hay debate legítimo. **vanilla-extract en TypeScript** es probablemente más maduro hoy en absoluto. Pero para un proyecto **Go-first + TinyGo + sin generadores**, las restricciones eliminan a TypeScript del set de soluciones aplicables.

### 6.2 ¿Es la mejor forma de escribir CSS en Go para tinywasm?
**Sí, dentro del set de soluciones compatibles con las restricciones del proyecto.** No conozco una alternativa que cumpla simultáneamente:
- Cero CSS en binario WASM
- Sin generadores
- Detección de typos en compile-time
- Tokens como entidades de primera clase
- Compartir nombres de clase entre HTML y CSS
- Una sola forma de hacerlo

Las cinco primeras existen aisladas en otras propuestas; ninguna las junta todas.

---

## 7. Veredicto

**Sí, es la forma más intuitiva, legible y profesional de escribir CSS en Go para tinywasm**, condicionado a que se acepten los costes nombrados en la sección 5 (sobre todo: verbosidad y curva inicial). Es defendible profesionalmente porque replica un patrón con ~10 años de adopción industrial en otros lenguajes tipados, adaptado a las restricciones específicas de Go + TinyGo + arquitectura SSR del proyecto.

Si la respuesta no convence, los puntos a cuestionar primero son:
1. ¿Los costes de la sección 5 son aceptables para tu equipo?
2. ¿La cobertura del 90% de CSS (con escape hatch para el resto) es suficiente, o el proyecto necesita CSS-spec completo?
3. ¿El dot-import es aceptable como convención del ecosistema?

Si los tres son "sí", el patrón es el indicado. Si alguno es "no", hay que revisitar el análisis antes de ejecutar los PLAN.md.

---

## Referencias

- W3C Design Tokens Community Group: <https://design-tokens.github.io/community-group/>
- vanilla-extract: <https://vanilla-extract.style/>
- Linaria: <https://linaria.dev/>
- Stitches (archivado pero referencia conceptual): <https://stitches.dev/>
- Lightning Design System (origen del término "design token"): <https://www.lightningdesignsystem.com/design-tokens/>
- Plan técnico de implementación: [`PLAN_typed_css.md`](./PLAN_typed_css.md)
