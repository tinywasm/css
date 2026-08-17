# Cómo un tema llega al color de un componente

Un sitio declara su paleta con una sola llamada:

```go
func RootCSS() *css.Stylesheet {
	return css.Theme(
		css.Set(css.ColorPrimary, "#16a34a"),
		css.Set(css.ColorSurface, "#ffffff"),
		css.Set(css.ColorOnSurface, "#111827"),
	)
}
```

Para que eso funcione, dos cosas tienen que sostenerse a la vez, y son opuestas.

## Las dos restricciones

**1. Lo que emite un componente tiene que ser alcanzable desde `:root`.**
`widget/style` pinta cada superficie con `Token.EnhancedVar()`. Si eso devuelve
un literal horneado en tiempo de compilación de Go, ningún `Theme(Set(...))`
puede cambiarlo. Y no es un rincón de la API: es el fondo, el texto y el borde
de **toda** superficie del ecosistema.

**2. La propiedad que alcanza tiene que ser computable en cualquier motor.**
Aquí está la asimetría que hace esto sutil:

| Situación | Qué hace el navegador |
|---|---|
| Declaración que **no puede parsear** (`background: light-dark(#a,#b)` sin soporte) | La descarta. Una declaración hermana anterior para la misma propiedad **queda en pie**. |
| Declaración con `var()` cuyo valor sustituido **no parsea** | Es *invalid at computed-value time*. **No** cae a la hermana anterior: cae al valor inicial — `transparent`, para un color. |

Es decir: una custom property que **contiene** `light-dark()` envenena toda
referencia `var()` a ella en un motor sin esa función. No "color equivocado":
transparente.

## La forma que cumple las dos

El veneno se quita en el origen. `:root` declara el valor **estático** sin
guardia, y vuelve a declarar el **adaptativo** dentro de una consulta de
capacidad:

```css
:root { --color-surface: #F2F2F7; }

@supports (color: light-dark(#000, #fff)) and (color: color-mix(in oklab, #000, #fff)) {
  :root { --color-surface: light-dark(#F2F2F7, #161B22); }
}

/* lo que el sitio declara, al final de la hoja */
:root { --color-surface: #ffffff; }
```

Con eso `--color-surface` **siempre** contiene algo computable, así que
`EnhancedVar()` puede ser el `var()` pelado y un componente emite:

```css
background-color: #F2F2F7;              /* estático, por si la propiedad faltara */
background-color: var(--color-surface); /* el que un tema alcanza */
```

- Motor sin `light-dark()`: el bloque `@supports` no aplica, la propiedad vale
  `#F2F2F7`, el `var()` computa. Tema claro permanente — el comportamiento que
  esta protección siempre buscó.
- Motor moderno: la propiedad vale `light-dark(...)` y se resuelve según
  `color-scheme`.
- Sitio con tema: su override va al final, gana sobre ambas mitades.

El orden es la parte que hay que respetar: estático → `@supports` → overrides
de la app. `withRootTail` pone los overrides al final justamente por esto.

## Dónde vive cada pieza

| Pieza | Archivo |
|---|---|
| `declare(t)` — valor estático en `:root` | `dsl.go` |
| `enhancedDecls(ts...)` — valor adaptativo | `dsl.go` |
| `ModernColorSupport` — la condición | `dsl.go` |
| Los tokens que se parten en dos | `themeTokens()` en `css.default.go` |
| `Token.EnhancedVar()` | `tokens.go` |

## Lo que sigue siendo literal, a propósito

`Token.NestedEnhanced()` **no** usa `var()` y no debe empezar a usarlo. Es para
cuando el token es un **argumento dentro de otra** llamada `color-mix()`/
`light-dark()` (los tokens compuestos de `catalog.go`, `mixToward` en `css.go`).
Ahí la declaración externa tiene que fallar en tiempo de *parseo* para que su
hermana estática sobreviva; meter un `var()` dentro la haría parsear y fallar
después, que es el caso malo de la tabla de arriba.

## Tema único

Un sitio sin modo oscuro — la mayoría de los sitios de marca — fija qué mitad
gana con `color-scheme`:

```css
:root { color-scheme: light }
```

`css.reset.go` ya emite `color-scheme: light` para `[data-theme="light"]`, así
que ponerlo en el `<html>` consigue lo mismo por la vía tipada.

## Qué se rompió antes de esto

`veltylabs/mjosefa-website`, migración v2 → v3. El sitio declaraba su paleta
completa con la API sancionada y se pintaba con la de la librería. En un
visitante con el sistema en oscuro, la web de una clínica cuya v2 no tiene tema
oscuro salía casi negra. La regresión está fijada en
`TestThemeOverrideReachesComponentDeclarations` (`css_test.go`).
