# PLAN — `ssr.go` → split por extensión (`css.go`)

## Objetivo

Renombrar `css/ssr.go` a `css/css.go` para alinearse con la nueva convención
del motor de `assetmin`: los assets SSR se descubren por archivos con nombre de
extensión (`css.go`, `js.go`, `html.go`, `svg.go`), todos `//go:build !wasm`.
El nombre reservado `ssr.go` se elimina del ecosistema.

## Justificación

`ssr.go` es un nombre mágico que no comunica su contenido. `css.go` es
autoexplicativo y alinea con SRP (`core-principles`). `tinywasm/css` es el
proveedor canónico del tema `:root` (vía `RootCSS()`) y del reset base (vía
`RenderCSS()`), así que es la pieza de referencia: su archivo debe modelar la
convención correcta para todo adoptante.

Breaking change coordinado a nivel monorepo — ver el stage homónimo en
`assetmin/docs/PLAN.md`.

## Estado actual

`css/ssr.go` contiene **funciones package-level** (sin receiver):

- `RootCSS() *Stylesheet` — todos los design tokens (`:root`).
- `RenderCSS() *Stylesheet` — reset base + bindings de tema claro/oscuro.

Ambos métodos producen CSS, así que ambos van al **mismo `css.go`** según la
regla de mapeo (RootCSS y RenderCSS → `css.go`).

## Cambios

- Renombrar `css/ssr.go` → `css/css.go`. Contenido **literal**: mismo build
  tag, mismo package, mismas dos funciones. No cambia lógica ni firmas.
- Verificar que ningún doc de `css/` (README, `JUSTIFICACION_DSL.md`) siga
  refiriéndose a `ssr.go` como el archivo donde declarar el tema.

## Precondición técnica

`assetmin` debe estar publicado con la whitelist `ssrSourceFiles`
(`css.go/js.go/svg.go/html.go`) y sin reconocer ya `ssr.go`. Aplicar este
renombrado en el **mismo PR coordinado** que el cambio de motor, para no dejar
el extractor incapaz de descubrir el tema base.

```bash
go list -m github.com/tinywasm/assetmin
```

## Tests y validación

- `go test ./...` verde en `tinywasm/css` (los tests no dependen del nombre del
  archivo).
- Verificar vía el flujo de extracción de `assetmin` que `RootCSS()` sigue
  ganando el slot `open` (regla single-override) tras el renombrado.

## Stages

| # | Tarea | Done |
|---|---|---|
| 1 | Confirmar precondición: `assetmin` con whitelist `ssrSourceFiles` publicado | [ ] |
| 2 | Renombrar `css/ssr.go` → `css/css.go` (contenido literal) | [ ] |
| 3 | Barrer docs de `css/` para no mencionar `ssr.go` | [ ] |
| 4 | `go test ./...` verde | [ ] |
| 5 | Verificar extracción del tema base vía `assetmin` sin regresiones | [ ] |
