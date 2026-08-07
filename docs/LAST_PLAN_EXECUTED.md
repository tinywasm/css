---
PLAN: "feat(css): token --chip-height para que un chip sea una caja de tamaño conocido"
TAG: v0.5.2
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> Es la **etapa 1 de 4** de una ola que cruza cuatro repos. Orden obligatorio:
> **css → widget → components → layout**. `widget` depende de este token; no
> empezar `widget` hasta que esto esté publicado.

# Plan — `--chip-height`: un chip con altura declarada, no emergente

## 1. Por qué

Hoy la altura de un chip (la leyenda de un `fieldset`, el badge de una fila de
`targetlist`) es **emergente**: sale de `font-size × line-height` más el padding
que le toque. Nadie la declara y nadie puede leerla.

Eso tiene dos consecuencias, y la segunda es la que rompe cosas:

1. **Dos chips coinciden por casualidad, no por construcción.** El comentario en
   `components/fieldset/css.go` ya lo dice con todas las letras: la leyenda tiene
   que medir lo mismo que el badge de al lado, y lo consigue evitando añadirle
   padding vertical. Es un acuerdo verbal entre dos archivos de repos distintos.
2. **`OnEdge` no puede reservar espacio.** Para montar un chip *sobre* una línea
   de borde, `widget/style` lo desplaza con `transform: translateY(±50%)` — exacto
   a cualquier tamaño de fuente, precisamente porque la altura es desconocida.
   Pero un `transform`:
   - es **invisible para `scrollHeight`/`clientHeight`** (por especificación), así
     que ningún contenedor puede reservar el hueco que el chip ocupa de verdad;
   - crea un *stacking context*, con efectos de apilamiento no declarados;
   - no reserva espacio en el flujo, así que el chip pisa a quien tenga debajo.

Ese tercer punto es el que produjo el defecto observado: el badge de la última
fila de una lista se solapaba con el botón de acción flotante, y **ninguna
cantidad de padding en el contenedor podía corregirlo**, porque el motor de
layout nunca "ve" dónde se pinta realmente el badge.

Con la altura en un token, el desplazamiento se puede calcular con márgenes
reales (`calc(-0.5 * var(--chip-height))`) en vez de con `transform`. La caja
pasa a existir para el layout, y todo lo de arriba se cae solo.

**Este plan solo añade el token.** Quien lo usa es `widget` (etapa 2).

## 2. Contexto del repo para un agente sin contexto previo

- Módulo: `github.com/tinywasm/css`. `docs/PLAN.md` va junto a `go.mod`.
- Los tokens viven en `catalog.go` como valores `Token{Name, Light, Dark, ...}`.
- Un token no sirve de nada si no se **declara** además en `css.default.go`
  (o `css.brand.go` si es identidad de marca): el catálogo define el valor, la
  declaración lo emite en `:root`.
- Hay un test que falla si se emite un token no registrado —
  `TestNoUndeclaredTokensInEmittedCSS` en `css_test.go` — con una lista explícita
  de tokens conocidos que **también** hay que actualizar.
- Nada de librería estándar en paquetes que compilan a WASM: usar `tinywasm/fmt`,
  nunca `errors`/`strconv`/`strings`.
- Prohibidas las cadenas repetidas en la lógica: todo literal repetido va a una
  constante con nombre.

## 3. Etapas

### Etapa 1 — declarar el token

En `catalog.go`, junto a `ChipWidth` (que ya existe y es su pareja natural):

```go
// La altura que comparte todo chip — la leyenda de un campo, el badge de una
// fila — para que un chip sea una caja de tamaño CONOCIDO y no emergente.
// Sin esto la altura sale de font-size × line-height y dos chips solo coinciden
// por casualidad; con esto, OnEdge puede montar el chip sobre una línea de borde
// con márgenes reales en vez de un transform, que es invisible para el cálculo
// de scroll y no reserva espacio.
ChipHeight = Token{Name: "--chip-height", Dark: "1.25rem"}
```

`1.25rem` (20px) es el valor que el chip mide **hoy** con `TextXs` (0.75rem) y el
`line-height` heredado — medido en el navegador sobre `.tw-field__label`. El token
no cambia nada visualmente: fija lo que ya ocurría.

Declararlo en `css.default.go`, en el mismo grupo donde ya está `ChipWidth`
(no es identidad de marca, así que **no** va en `css.brand.go`).

Añadirlo a la lista `allTokens` de `TestNoUndeclaredTokensInEmittedCSS` en
`css_test.go`, junto a `ChipWidth`.

**Aceptación:**
- `grep -n "chip-height" catalog.go css.default.go css_test.go` devuelve las tres.
- `go build ./... && go test ./... -count=1` en verde.
- El CSS emitido por `RootCSS()` contiene `--chip-height` en `:root`.

### Etapa 2 — no hacer nada más

No tocar `ChipWidth`, ni `ControlHeight`, ni la escala `Z*`. Este plan es un
token y su declaración. El consumo va en el plan de `widget`.

| Etapa | Archivos | Puerta |
|---|---|---|
| 1 | `catalog.go`, `css.default.go`, `css_test.go` | — |

## 4. Lo que este plan NO hace

- No cambia `OnEdge` (eso es `widget`).
- No revisa la escala de z-index. Los tokens `ZBase`…`ZTooltip` siguen siendo
  para **capas de overlay** (dropdown, modal, toast) y no para orden local entre
  hermanos; esa distinción se documenta en el plan de `widget`, donde vive el
  código que la aplica.
