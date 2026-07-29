# tinywasm/css Technical Specification

This document provides a comprehensive technical specification for the `tinywasm/css` library. It outlines the core architecture, design token definitions, API contract, and theme customizability model.

---

## 1. Design Token Architecture

The library acts as a single, typed source of truth for all layout, spacing, typography, and color decisions. It replaces error-prone string-based CSS with strongly-typed Go structures.

### Core Type: `Token`

Each decision is modeled as a `Token`:

```go
type Token struct {
    Name       string // CSS Custom Property identifier (e.g., "--color-primary")
    Light, Dark string // Dark solo = static fallback; ambos = light-dark pair
}

// Var returns the standard var() syntax for consuming this token in CSS.
func (t Token) Var() string {
    if t.Light == "" {
        return "var(" + t.Name + "," + t.Dark + ")"
    }
    return "var(" + t.Name + ",light-dark(" + t.Light + ", " + t.Dark + "))"
}
```

---

## 2. Color Palette & Theming (12-Token System)

Rather than maintaining separate and verbose light/dark mode source variables (which ballooned to over 69 variables previously), `tinywasm/css` leverages a modern, streamlined **12-token system**.

Theme adaptation (Light/Dark mode) is computed natively in the browser using the standard CSS `light-dark()` function, reducing complexity and cascade leaks.

### Active Theme Tokens (15-Color Catalog)

| Token Variable | CSS Property | Default Fallback Value | Description / Usage |
|---|---|---|---|
| `ColorPrimary` | `--color-primary` | `#1b5d8c` | Primary brand accent color for links, primary buttons, and key UI elements. |
| `ColorOnPrimary` | `--color-on-primary` | `#FFFFFF` | Text/foreground color used directly on top of `ColorPrimary`. |
| `ColorSuccess` | `--color-success` | `#1e7a30` | Visual color indicating success, positive, or completed states. |
| `ColorOnSuccess` | `--color-on-success` | `#FFFFFF` | Text/foreground color used directly on top of `ColorSuccess`. |
| `ColorDanger` | `--color-danger` | `#ba2c0d` | Destructive/Error accent color for alerts, warnings, and dangerous actions. |
| `ColorOnDanger` | `--color-on-danger` | `#FFFFFF` | Text/foreground color used directly on top of `ColorDanger`. |
| `ColorBackground` | `--color-background` | `light-dark(#FFFFFF, #0D1117)` | Page canvas background, responding to user/system preferences. |
| `ColorOnBackground` | `--color-on-background` | `light-dark(#1C1C1E, #E6EDF3)` | High contrast body text on the canvas background. |
| `ColorSurface` | `--color-surface` | `light-dark(#F2F2F7, #161B22)` | Surface/Card/Panel background color. |
| `ColorOnSurface` | `--color-on-surface` | `light-dark(#1C1C1E, #E6EDF3)` | Primary text on cards/panels. |
| `ColorOutline` | `--color-outline` | `light-dark(#D1D1D6, #30363D)` | Subtle borders, lines, and structural dividers. |
| `ColorMuted` | `--color-muted` | `light-dark(#6E6E73, #8B949E)` | Lower-contrast typography for captions, metadata, or supporting descriptions. |
| `ColorSurfaceSunken` | `--color-surface-sunken` | `color-mix(in oklab, var(--color-surface), var(--color-on-surface) 8%)` | Inset surface, e.g. for sunken panels; derives from `--color-surface` via `var()`. |
| `ColorSelection` | `--color-selection` | `color-mix(in oklab, var(--color-primary), transparent 85%)` | Highlighted/selected item background; tinte translúcido sobre `--color-surface`. |
| `ColorOnSelection` | `--color-on-selection` | `var(--color-on-surface)` | Text on `ColorSelection`; válido solo cuando el fondo es `--color-surface`. |

---

## 3. Scale Definitions

### Spacing Scale (4px Base Grid)
- `Space0` (`--space-0`): `0`
- `Space1` (`--space-1`): `0.25rem` (4px)
- `Space2` (`--space-2`): `0.5rem` (8px)
- `Space3` (`--space-3`): `0.75rem` (12px)
- `Space4` (`--space-4`): `1rem` (16px)
- `Space6` (`--space-6`): `1.5rem` (24px)
- `Space8` (`--space-8`): `2rem` (32px)
- `Space12` (`--space-12`): `3rem` (48px)

### Border Radius Scale
- `RadiusSm` (`--radius-sm`): `4px`
- `RadiusMd` (`--radius-md`): `8px`
- `RadiusLg` (`--radius-lg`): `16px`
- `RadiusFull` (`--radius-full`): `9999px`

### Typography Scale
- `TextXs` (`--text-xs`): `0.75rem`
- `TextSm` (`--text-sm`): `0.875rem`
- `TextBase` (`--text-base`): `1rem`
- `TextLg` (`--text-lg`): `1.25rem`
- `TextXl` (`--text-xl`): `1.5rem`
- `Text2xl` (`--text-2xl`): `2rem`

### Interaction Intensity Scale
- `MixHover` (`--mix-hover`): `15%`
- `MixFocus` (`--mix-focus`): `30%`
- `MixPress` (`--mix-press`): `45%`

Values reproduce the legacy hex scale (`#164d74`/`#123e5d`/`#0e304a` from `#1b5d8c`).
These are declared in `:root` for discoverability and can be overridden with
`Set(MixHover, "22%")`.

---

## 4. API Surface & Asset Generation

### Interaction Derivation
Interaction states are derived from any base token via CSS `color-mix()`:

```go
func Hover(t Token) string  // "color-mix(in oklab, <t>, light-dark(black, white) var(--mix-hover,15%))"
func Focus(t Token) string  // same with --mix-focus
func Press(t Token) string  // same with --mix-press
```

Devuelven un string (expresión CSS inline), no un `Token`: no inflan el catálogo
con 27 variables de estado. La intensidad es un token (`MixHover`/`MixFocus`/
`MixPress`), por lo que una app puede reajustarla con `Set(MixHover, "22%")`.

El mezclador es `light-dark(black, white)`, no un blanco o negro fijo, para que
la dirección del tema sea correcta en ambos modos (oscurece en light, aclara en
dark).

### RootCSS and RenderCSS

The library provides exactly two functions matching the `assetmin` SSR build pipeline contract:

### 1. `RootCSS() *Stylesheet`
Emits the core CSS Custom Property `:root` block declarations. It holds the vocabulary of design decisions.
```go
// Returns a Stylesheet declaring all variables inside :root {}
root := css.RootCSS()
```

### 2. `RenderCSS() *Stylesheet`
Emits the actual functional CSS rules (resets, layout, element resets).
It configures standard `color-scheme` rules on `:root`, `[data-theme="light"]`, and `[data-theme="dark"]`, enabling seamless support for modern CSS `light-dark()` functions:

```css
:root {
  color-scheme: light dark;
}
[data-theme="light"] {
  color-scheme: light;
}
[data-theme="dark"] {
  color-scheme: dark;
}
```

---

## 5. Overriding & Theme Customization

Applications can customize brand themes without losing the complete token scales by using `css.Theme(...Override)`:

```go
func Theme(overrides ...Override) *Stylesheet
```

Overrides are generated via `Set(Token, string)` for static values or
`SetTheme(Token, light, dark string)` for theme-aware pairs:

```go
import "github.com/tinywasm/css"

func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.ColorPrimary, "#FF6B35"),
        css.SetTheme(css.ColorBackground, "#FAFAFA", "#121212"),
        css.Set(css.MixHover, "22%"),
    )
}
```

Two levels of override:
1. **Change the base** (`SetTheme(ColorSurface, ...)`) — all derived tokens
   (`--color-surface-sunken`) recompute via live `var()` in the browser.
2. **Change the derived** (`Set(ColorSurfaceSunken, "...")`,
   `Set(MixHover, "22%")`) — substitutes the formula or intensity directly.

**Warning:** a `Set(MixHover, "not-a-percentage")` makes any declaration using
it invalid at computed-value time (the `background-color` becomes `unset`). The
`var()` fallback covers "undefined", not "ill-defined". Override with valid CSS.

This generates the baseline design catalog variables first, then appends custom overridden definitions at the end of the `:root` block to let the CSS cascade natively pick up overridden decisions.
