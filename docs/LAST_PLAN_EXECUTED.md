---
PLAN: "feat: emitir el @font-face de la familia declarada"
TAG: v0.5.0
---
## Antes de escribir código: lee [CONSTRUCTION_HARNESS.md](CONSTRUCTION_HARNESS.md)

**Es vinculante, no orientativo.**

| # | Principio | Cómo se aplica aquí |
|---|---|---|
| 4 | One way to do each thing | El texto CSS lo escribe esta pieza. Si `assetmin` formatea un `@font-face`, hay dos formas de escribir CSS. |
| 5 | Minimal surface | Una función. El prefijo de URL entra como parámetro; `css` no aprende rutas de servidor. |
| 1 | Typed over `any` | El peso y el estilo se derivan de `font.Style`, no se reciben como strings. |

---

## 1. El hueco

`css` ya sabe **nombrar** la tipografía del producto: `FontStack(font.Family)` (`css.go:8`)
alimenta el token `FontSans` (`catalog.go:28`). Eso resuelve *qué* fuente pedir.

No sabe **declararla**. No hay forma de emitir el `@font-face` que le dice al navegador
dónde está el archivo, y sin él `--font-sans: "Roboto", …` cae al `system-ui` de después
— que es exactamente el síntoma que arrancó todo esto: la app no se ve igual en Safari
iOS que en Chrome Android.

El hueco tiene consecuencia documentada: una versión anterior de este mismo plan proponía
embeber la fuente en base64 dentro del CSS para esquivarlo, y la versión actual de
`assetmin/docs/PLAN.md` estuvo a punto de armar la regla con un `fmt.Sprintf` en ese otro
repo. Las dos son la misma falta —*a fork with a friendlier name*— y las dos desaparecen
cuando la pieza que posee la sintaxis CSS ofrece la función.

### Por qué la URL entra como parámetro

La regla necesita la URL final del archivo, y esa la decide `assetmin`
(`AssetsURLPrefix` + `OutputDir`). `css` no la conoce y adivinarla lo acoplaría a una
convención que no controla. Recibirla como un `string` no es acoplamiento: es el llamador
diciendo dónde sirve sus assets.

---

## 2. Cambio

Archivo nuevo `css.fontface.go`, con el mismo `//go:build !wasm` que el resto del paquete:

```go
// FontFaces devuelve el bloque @font-face de las cuatro caras de la familia
// declarada, servidas desde urlPrefix (p. ej. "/assets"). El prefijo es un dato
// del llamador: este paquete no conoce ni inventa rutas de servidor.
//
// No entra en RootCSS() ni en RenderCSS(): quien sirve los archivos decide
// cuándo inyectarlo.
func FontFaces(d font.Declaration, urlPrefix string) *Stylesheet
```

Una regla por cara, en el orden de los cuatro `font.Style`. `font` ya está en `go.mod`
(`css.go:5` lo importa).

### 2.1 Peso y estilo se derivan, no se inventan

| `font.Style` | `font-weight` | `font-style` |
|---|---|---|
| `Regular` | 400 | normal |
| `Bold` | 700 | normal |
| `Italic` | 400 | italic |
| `BoldItalic` | 700 | italic |

Un `switch` sobre los cuatro constantes. Nada de mapas (regla del ecosistema) y nada de
recibir el peso por parámetro: `font.Style` ya lo dice.

### 2.2 El nombre del archivo lo deriva `font`

`d.Family().Face(s) + ".ttf"`. **Nunca** concatenar `"-Bold"` a mano: la derivación vive
en `font` y duplicarla es el defecto que ese repo está cerrando ahora mismo
(`font/docs/PLAN.md`, `Face(Regular)` → `-Regular`).

**Prerrequisito:** publicar `font` v0.1.0 antes de ejecutar esto, o las URLs emitidas
apuntarán al nombre viejo de la cara regular.

### 2.3 `format("truetype")`, no `woff2`

El ecosistema sirve **un solo TTF por cara para web y PDF**: el documento se genera en el
frontend, así que pide el mismo archivo que la página ya bajó — acierto de caché en vez
de una segunda descarga (16.132 B y 1 petición, frente a 27.890 B y 2). Escribir
`format("woff2")` haría que el navegador **rechace** el archivo, y el fallo se ve como
«la fuente no carga», no como un error. Números en
`app-releases/docs/TYPOGRAPHY_MASTER_PLAN.md`.

`font-display: swap` en las cuatro: el texto se pinta con la fuente de respaldo y se
recompone al llegar la real, en vez de quedar invisible.

### 2.4 Se construye con el DSL interno, no con `Raw`

`Raw` es la válvula de escape **para consumidores**, no para este paquete. `@font-face`
es una at-rule y aquí ya hay tres implementadas con el mismo patrón —`layerItem`
(`dsl.go:107`), media (`:124`), keyframes (`:205`)—: un tipo privado con su `writeTo`.
Seguir ese patrón. Si `FontFaces` se implementara con `Raw`, esta pieza estaría usando
su propia puerta de emergencia teniendo la puerta.

El DSL sigue **sin exportarse**: se exporta la función, no la at-rule.

---

## 3. Documentación

- `docs/SPECS.md`: la tabla de §2.1 como aserción de test, la salida exacta para
  `Declare("Roboto","x")` con prefijo `/assets`, y la nota de que `FontFaces` no forma
  parte de `RootCSS()`/`RenderCSS()`.
- `README.md`: una línea en la superficie pública. Es el quinto emisor de CSS junto a
  `RootCSS`, `RenderCSS`, `Stylesheet` y `Raw`.

---

## 4. Verificación

1. `FontFaces(font.Declare("Roboto", ""), "/assets")` emite **cuatro** bloques, uno por
   cara, con las URLs `/assets/Roboto-Regular.ttf`, `-Bold`, `-Italic`, `-BoldItalic`.
2. Los cuatro declaran `format("truetype")` y `font-display:swap`.
3. Peso y estilo coinciden con la tabla de §2.1, cara por cara.
4. El prefijo se une sin barra doble ni barra ausente: `"/assets"`, `"/assets/"` y `""`
   producen URLs válidas.
5. Una familia vacía (`Declare("", …)`, que `font` permite construir) no emite reglas
   rotas: o no emite nada, o la ausencia se ve. Decidirlo y afirmarlo en el test.
6. El paquete sigue siendo build-time-only: `TestPackageIsBuildTimeOnly`
   (`tokens_test.go`) sigue verde y el archivo nuevo lleva `//go:build !wasm`.
7. Ningún `map[` en el código nuevo.
8. `gotest`.

---

## 5. Lo que este plan NO hace

- **No copia archivos ni conoce `OutputDir`.** Eso es de `assetmin`
  (`assetmin/docs/PLAN.md`), que llamará a esta función con su prefijo.
- **No mete el `@font-face` en `RootCSS()`.** `RootCSS()` sale por la extracción de `ssr`
  y no sabe nada de URLs de assets.
- **No toca `FontStack` ni el token `FontSans`.** Ya funcionan y resuelven otra mitad del
  problema: ésta declara la fuente, aquélla la pide.
