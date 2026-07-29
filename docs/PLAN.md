---
PLAN: "refactor: variables de color minimas"
TAG: v0.2.0
---

# problema:
por que tantas variables de color? cual es el minimo que deberia tener un sitio web para que se vea bien? para que añador tanta variabilidad?..cual es el estandar en la industria? lo que se busa es simplifciar la paleta de colores de cualquier sitio/aplicaicon..ya que la camtidad de tojken actual es abrimadoramente  excesiva.

# investigacion:
Cuántas hay realmente
En tokens.go hay 69 tokens de color:

Bloque	Cantidad
Brand (primary/secondary/success/error + on-*)	8
Theme activo (background, surface, outline, muted, hover, selection, disabled…)	11
Theme source (los -light / -dark)	22
Focus ring	1
Interactivos (hover/focus/press × 9 familias)	27
Los otros grupos (espaciado, radios, z-index, tipografía) están bien dimensionados. El problema es exclusivamente el color.

Peso muerto concreto que encontré
ColorDanger* y ColorError* son idénticos — mismos hex (#9a240b, #7a1c09, #5b1507). Uno de los dos sobra literalmente.
Warning, Info, Neutral tienen hover/focus/press pero no tienen token base. No puedes pintar un botón warning, solo su hover. Son estados huérfanos.
11 tokens se declaran pero nunca se leen. emit.go:346-363 (familyTokens) es el único consumidor real, y no toca ColorError{Hover,Focus,Press}, ColorWarning*, ColorInfo*, ColorMutedHover/Press. Se emiten en el CSS de cada página sin que nadie los use.
Por qué se acumularon
Las dos capas grandes (source light/dark = 22, e interactivos = 27) no son decisiones de diseño, son mecanismos. Alguien necesitaba dark mode sin JS y estados de interacción, y en 2020 la única forma era declarar cada variante a mano. Hoy el navegador las calcula.

El mínimo real
El consenso de la industria en paletas semánticas:

shadcn/ui (la referencia actual más usada): ~11 decisiones de hue reales.
GOV.UK Design System: 5 colores funcionales. Es el sistema accesible más auditado que existe.
Material 3: ~26 roles — es el outlier, y está diseñado para generar temas dinámicos desde una foto de wallpaper. No es tu caso.
Tailwind: cero tokens semánticos, solo una escala.
El mínimo funcional son 6 decisiones: fondo, texto, borde, texto secundario, acción (marca), destructivo. Con eso se construye un sitio completo y accesible.


# Propuesta para este proyecto
12 tokens en lugar de 69, sin perder ninguna capacidad:


ColorBackground  / ColorOnBackground
ColorSurface     / ColorOnSurface
ColorOutline
ColorMuted
ColorPrimary     / ColorOnPrimary
ColorSuccess     / ColorOnSuccess
ColorDanger      / ColorOnDanger
Lo que se elimina y cómo se reemplaza:

Se borra	Reemplazo
22 tokens -light/-dark	light-dark(#F2F2F7, #161B22) dentro del token activo — un solo token declara ambos temas
27 tokens hover/focus/press	color-mix(in oklab, var(--color-primary), black 12%) calculado en familyTokens
ColorFocusRing	es ColorPrimary
ColorSelection / ColorOnSelection	color-mix de primary al 15% sobre el fondo
ColorDisabled / ColorOnDisabled	ColorSurface + ColorMuted
ColorSecondary	rara vez justificado; un sistema con primary + neutral cubre el 95%
ColorSurfaceSunken	color-mix de surface con el fondo
light-dark() y color-mix() son Baseline desde 2024 y 2023 respectivamente — soporte universal hoy.

Ganancia colateral: AllPairs() pasa de 15 pares a 5. Rebrandear una app pasa de redefinir 22 variables source a redefinir 3 (--color-primary, --color-background, --color-surface); el resto se deriva solo y mantiene el contraste por construcción.

# entregables:
- [ ] refactor de tokens.go para que emita solo los 12 tokens propuestos.
- [ ] corregir tests que fallen por la eliminación de tokens.
- [ ] crear docs/SPECS.md con la especificacion actal de toda la libreria tinywasm/css y enlazarlo al readme.
- [ ] actualizar todo la documentacion relevante
- [ ] comentar todas la variables/tokens publicos con su uso para tener una claro idea de su propósito