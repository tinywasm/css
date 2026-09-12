---
PLAN: "feat(css): degradado primary→gradient como default de marca, no overhead por app"
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> Es la **etapa 1 de 2** de una ola de dos repos. Orden obligatorio:
> **css → app-demo**. No tocar `app-demo` hasta que este repo esté publicado
> (el override que se borra en `app-demo` depende del nuevo default de aquí).

# Plan — el degradado violeta→cian pasa a ser el default de `css.Theme()`

## 1. Por qué

Hoy `css.ColorPrimary` ya tiene como valor de catálogo `#654FF0` (el violeta
de WebAssembly) — eso ya es el default real de todo proyecto webtyp, sin que
nadie lo declare. Pero el **degradado** que se ve en `webtyp/app-demo` (una
diagonal de ese violeta hacia `#00ADD8`, el cian del gopher de Go) NO es un
default del framework: es un override que `app-demo/config/css.go` declara
por su cuenta (`css.Set` + `css.SetGradient`, con una constante local
`GoCyan`).

Consecuencia real observada: un segundo proyecto (`mjosefa-cms`) que
simplemente quiere "verse bien por defecto" solo puede obtener el color
plano (`#654FF0` sólido), nunca el degradado, a menos que copie el mismo
bloque de `Set`/`SetGradient` que `app-demo` ya tiene. Eso es exactamente lo
que la regla de "piezas de lego" de este ecosistema prohíbe: una decisión
visual que se repite en dos proyectos debió vivir en la librería desde la
primera vez.

Este plan mueve el degradado al catálogo de `webtyp.com/css`, para que
**cualquier proyecto que llame `css.Theme()` sin overrides lo reciba
automáticamente** — cero configuración, un solo lugar que mantener.

## 2. Design gate (API pública nueva)

Este plan agrega un token exportado (`ColorPrimaryGradient`) y una función
exportada (`ClearGradient`). Por regla de la skill `api-design`, van las
cinco respuestas:

1. **Prior art**: Next.js/Vercel arrancan todo proyecto nuevo con su propio
   acabado visual de marca (el "glow" negro) sin que el usuario configure
   nada; Material Design 3 arranca con un color semilla por defecto
   (púrpura) que cualquier app hereda hasta reskinearlo; Stripe expone sus
   propios starters con el degradado morado-azul de marca ya aplicado. El
   patrón común: un framework con identidad visual propia arranca opinado,
   no neutro — y siempre deja un mecanismo de override de una sola llamada.
   Este cambio sigue exactamente ese patrón.
2. **Novice-name test**: "`ColorPrimaryGradient` es el segundo color hacia el
   que se degrada `ColorPrimary`" — se lee bien. "`ClearGradient` apaga el
   degradado de un token" — se lee bien.
3. **Complexity ledger**: +1 token, +1 función, +2 líneas en `brandRoot()`.
   A cambio se borra el bloque completo de override en `app-demo`
   (`goToken` + el cuerpo de `Theme.RootCSS()`, ver etapa 2). Con ≥2
   proyectos consumiendo el default, el balance ya es negativo, y mejora con
   cada proyecto nuevo que no tiene que volver a declararlo.
4. **Dónde vive**: `webtyp.com/css` ya es, por diseño documentado en
   `widget/docs/ARCHITECTURE.md` §2, el dueño de los *valores* globales
   (tokens de color/espacio/duración). Un acabado visual por defecto del
   token de marca es un valor, no una decisión de widget ni una
   composición de app — pertenece aquí, no en `webtyp/widget` ni en cada
   `config/css.go`.
5. **Qué borra**: el override de `app-demo/config/css.go` completo (etapa 3)
   queda redundante y se elimina en la misma ola.

## 3. Contexto técnico que ya existe (no inventar nada nuevo)

- `Token.ImageVarName()` / `Token.ImageStopsVarName()` (`tokens.go:81,90`) ya
  son las propiedades CSS que `widget/style` lee para pintar un degradado de
  fondo (`var(--color-primary-image, none)`), con `none` como fallback hoy.
- `SetGradient(t, angle, from, to)` (`css.go`) ya sabe construir el string
  `linear-gradient(...)` y sus dos declaraciones — pero solo como `Override`
  aplicado vía `Theme(overrides...)`, nunca como parte del catálogo base
  (`RootCSS()`/`brandRoot()`). Este plan replica esa misma construcción
  directamente en `brandRoot()`, como default permanente en vez de override.
- `Theme(overrides...)` (`css.go`) hace `append` de cualquier override
  **después** del catálogo base (`withRootTail`), así que un override
  posterior sigue ganando dentro del mismo `:root` — un app que llame
  `Theme(Set(ColorPrimary, "#16a34a"))` sigue funcionando exactamente igual
  que hoy; solo cambia el degradado por defecto cuando el app no dice nada.

## 4. Etapas

### Etapa 1 — nuevo token en `catalog.go`

Archivo: `catalog.go`. Agregar, junto al bloque de `ColorPrimary`/
`ColorOnPrimary` (línea 6-7):

```go
// ColorPrimaryGradient is the second stop ColorPrimary fades into by
// default — see brandRoot(). An app that wants a flat solid primary calls
// Theme(ClearGradient(ColorPrimary)); an app that wants a different second
// stop calls Theme(Set(ColorPrimaryGradient, "#yourColor")).
ColorPrimaryGradient = Token{Name: "--color-primary-gradient", Dark: "#00ADD8"}
```

Va dentro del mismo bloque `var (...)` existente, no en uno nuevo.

### Etapa 2 — degradado por defecto en `css.brand.go`

Archivo: `css.brand.go`. Reemplazar el cuerpo de `brandRoot()` completo por:

```go
// brandRoot declares the identity palette an app reskins to become its own:
// primary, success, danger and accent, each paired with its on-color — plus
// ColorPrimary's default gradient partner. Kept apart from defaultRoots() so
// a white-label override touches one small group instead of hunting through
// the full token catalog.
func brandRoot() item {
	decls := []decl{
		declare(ColorPrimary),
		declare(ColorOnPrimary),
		declare(ColorSuccess),
		declare(ColorOnSuccess),
		declare(ColorDanger),
		declare(ColorOnDanger),
		declare(ColorAccent),
		declare(ColorOnAccent),
		declare(ColorPrimaryGradient),
	}
	decls = append(decls, defaultGradient(ColorPrimary, "135deg", ColorPrimary, ColorPrimaryGradient)...)
	return root(decls...)
}

// defaultGradient declares t's background-image companion (and its stops
// companion, ImageStopsVarName) as a permanent default — the same shape
// SetGradient's Override produces for a one-off app override, but baked
// into the catalog so every consumer gets it without calling Theme(...).
// Theme() still appends any app override after this block, and CSS custom
// properties resolve last-declaration-wins, so Theme(SetGradient(...)) or
// Theme(ClearGradient(...)) still wins exactly like it does today.
func defaultGradient(t Token, angle string, from, to Token) []decl {
	stops := from.Var() + ", " + to.Var()
	return []decl{
		{t.ImageVarName(), "linear-gradient(" + angle + ", " + stops + ")"},
		{t.ImageStopsVarName(), stops},
	}
}
```

No tocar `defaultRoots()` ni ningún otro archivo de `catalog.go` fuera del
token agregado en la etapa 1.

### Etapa 3 — vía de escape: `ClearGradient` en `css.go`

Archivo: `css.go`, inmediatamente debajo de la función `SetGradient`
existente. Agregar:

```go
// ClearGradient turns off token t's default gradient, restoring a flat
// solid fill — the opposite of SetGradient. Use it when an app overrides
// t's own color (Theme(Set(t, ...))) and wants that override to render
// flat instead of inheriting t's catalog default gradient (see
// ColorPrimary / ColorPrimaryGradient in brandRoot()).
//
// It clears only ImageVarName() (what widget/style actually paints).
// ImageStopsVarName() is left as-is: nothing reads the stops companion
// without also reading the image var first, so there is nothing to
// desynchronize — do not add a second decl here for it.
func ClearGradient(t Token) Override {
	return Override{token: t, gradient: "none"}
}
```

No modificar la struct `Override` ni el switch dentro de `Theme()` en
`css.go` — `Theme()` ya sabe emitir `o.gradient` como
`decl{o.token.ImageVarName(), o.gradient}` sin cambios.

### Etapa 4 — tests

Archivo: `css_test.go`. Agregar, junto a los tests existentes
`TestRootCSS_*`:

```go
func TestRootCSS_DefaultsPrimaryToAGradient(t *testing.T) {
	css := RootCSS().String()
	if !strings.Contains(css, "--color-primary-image: linear-gradient(135deg,") {
		t.Fatalf("expected a default --color-primary-image gradient, got:\n%s", css)
	}
	if !strings.Contains(css, "--color-primary-gradient: #00ADD8") {
		t.Fatalf("expected --color-primary-gradient default declared, got:\n%s", css)
	}
}

func TestTheme_ClearGradientOverridesDefault(t *testing.T) {
	out := Theme(ClearGradient(ColorPrimary)).String()
	if !strings.Contains(out, "--color-primary-image: none") {
		t.Fatalf("expected ClearGradient to emit --color-primary-image: none, got:\n%s", out)
	}
}
```

Ajustar imports (`strings`) si el archivo no lo importa ya — revisar el
encabezado de `css_test.go` antes de agregarlo duplicado.

### Etapa 5 — criterios de aceptación

- `go test ./...` pasa en `webtyp.com/css`.
- `grep -n "func brandRoot" css.brand.go` muestra la nueva firma con
  `defaultGradient(...)`.
- `grep -n "ColorPrimaryGradient" catalog.go` y `grep -n "func ClearGradient" css.go`
  devuelven exactamente una coincidencia cada uno (sin duplicados).
- Publicar (el ejecutor NO corre `gopush`/`codejob` — eso es un paso externo
  al agente, ver skill agents-workflow).

## 5. Etapa 2 de la ola (repo `webtyp/app-demo`, plan separado)

Una vez publicado este cambio, `app-demo/config/css.go` queda con un
override que produce exactamente el mismo resultado que el nuevo default —
se vuelve redundante. Ese repo recibe su propio `docs/PLAN.md` que:

1. Borra `goToken`, el import de `"webtyp.com/css"` si queda sin uso, y el
   cuerpo actual de `func (Theme) RootCSS() *css.Stylesheet`.
2. Lo reemplaza por `return css.Theme()` (igual que
   `veltylabs/mjosefa-cms/config/css.go` hoy).
3. Verifica visualmente (captura de pantalla o `browser_screenshot` del MCP
   de webtyp) que el degradado violeta→cian se sigue viendo igual.

No escribir ese plan todavía — depende de que este publique primero.

| Etapa | Archivo | Acción |
|---|---|---|
| 1 | `catalog.go` | Agregar token `ColorPrimaryGradient` |
| 2 | `css.brand.go` | Reescribir `brandRoot()` + agregar `defaultGradient()` |
| 3 | `css.go` | Agregar `ClearGradient(t Token) Override` |
| 4 | `css_test.go` | Agregar los dos tests nuevos |
| 5 | — | Verificar criterios de aceptación |
