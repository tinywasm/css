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
- `FontSans` (`--font-sans`): `"Roboto", system-ui, -apple-system, sans-serif`
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

### Grid Columns
- `ColumnNarrow` (`--column-narrow`): `12rem`
- `ColumnMedium` (`--column-medium`): `20rem`
- `ColumnWide` (`--column-wide`): `30rem`

### Rail Widths
- `RailNarrow` (`--rail-narrow`): `3.5rem` — icon-only sidebar column
- `RailWide` (`--rail-wide`): `12rem` — icon + label sidebar column

### Device Geometry

Static tokens (`Dark` alone). Components write `var(--safe-top)` / `var(--viewport-h)`;
they never write a bare `env(...)` or `100vh`.

| Token | CSS Property | Value | Purpose |
|---|---|---|---|
| `SafeTop` | `--safe-top` | `env(safe-area-inset-top, 0px)` | Notch / Dynamic Island inset |
| `SafeRight` | `--safe-right` | `env(safe-area-inset-right, 0px)` | Right safe-area inset |
| `SafeBottom` | `--safe-bottom` | `env(safe-area-inset-bottom, 0px)` | Home-indicator / gesture bar inset |
| `SafeLeft` | `--safe-left` | `env(safe-area-inset-left, 0px)` | Left safe-area inset |
| `ViewportH` | `--viewport-h` | `100dvh` | Visible viewport height (shrinks with Safari iOS URL bar) |

`env()` returns `0px` on devices without insets. Visible effect also requires
`<meta name="viewport" content="…, viewport-fit=cover">` from `tinywasm/html` /
`tinywasm/assetmin` — not emitted here.

`dvh` is older than this library's baseline (`light-dark()`, `color-mix()`), so no
`@supports` guard is emitted.

**Where to apply:** `RootCSS()` is vocabulary only. Padding a header with
`var(--safe-top)` or sizing a shell with `var(--viewport-h)` is a layout decision
owned by `tinywasm/layout` / `tinywasm/widget/style`.

### Form control text size (usage constraint, not a reset rule)

A form control must not carry a text size below `TextBase` (`1rem` / 16px) on
touch devices: iOS Safari zooms on focus when the computed `font-size` is under
16px. The reset cannot enforce this — `RenderCSS()` sits in `@layer tokens`, so
any `@layer widgets` rule wins by layer order. The constraint lives in
`tinywasm/widget/style`; runtime detection is `browser_audit_mobile` in
`tinywasm/devbrowser`.

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

### RootCSS, RenderCSS and FontFaces

The library provides three CSS emitters. `RootCSS` and `RenderCSS` match the
`assetmin` SSR pipeline contract; `FontFaces` is injected separately by whoever
serves the font files.

### 1. `RootCSS() *Stylesheet`
Emits the core CSS Custom Property `:root` block declarations. It holds the vocabulary of design decisions.
```go
// Returns a Stylesheet declaring all variables inside :root {}
root := css.RootCSS()
```

### 2. `FontFaces(d font.Declaration, urlPrefix string) *Stylesheet`

Emits one `@font-face` rule per face of the declared family. **Not** part of
`RootCSS()` or `RenderCSS()`: whoever serves the font files decides when to
inject the stylesheet (typically `assetmin` with its own URL prefix).

Weight and style are derived from `font.Style` — never received as strings:

| `font.Style` | `font-weight` | `font-style` |
|---|---|---|
| `Regular` | 400 | normal |
| `Bold` | 700 | normal |
| `Italic` | 400 | italic |
| `BoldItalic` | 700 | italic |

File names come from `d.Family().Face(s) + ".ttf"` (derivation lives in
`tinywasm/font`). Format is always `format("truetype")` — one TTF per face for
web and PDF. Every rule sets `font-display: swap`.

An empty family (`Declare("", …)`) emits an empty stylesheet (no broken rules).

Exact output for `FontFaces(font.Declare("Roboto", "x"), "/assets")`:

```css
@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: url("/assets/Roboto-Regular.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: url("/assets/Roboto-Bold.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: italic;
  font-weight: 400;
  font-display: swap;
  src: url("/assets/Roboto-Italic.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: italic;
  font-weight: 700;
  font-display: swap;
  src: url("/assets/Roboto-BoldItalic.ttf") format("truetype");
}
```

URL prefix joining: `"/assets"`, `"/assets/"` and `""` all produce valid URLs
(no double slash, no missing slash).

### 3. `RenderCSS() *Stylesheet`
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

#### Cross-browser normalization guarantees

Beyond the box-model and user-agent-margin resets, `RenderCSS()` closes the points where
two shipping engines disagree, so a part's own rules land on the same box everywhere:

| Rule | Divergence it closes |
|---|---|
| `-webkit-tap-highlight-color: transparent` on `html` | iOS washes tapped elements in grey; Chrome Android uses another colour and duration |
| `appearance/background/border/radius` zeroed on `button, [type=button|reset|submit]` | iOS renders buttons as `push-button` chrome — its own radius, gradient and border — that outranks a part's declarations |
| `appearance: none; border-radius: 0` on text fields | iOS forces an inset shadow and its own corner radius on inputs |
| `text-transform: none` on `select` | Firefox and Edge let `<select>` inherit `text-transform`; other engines do not |
| `opacity: 1` on `::placeholder` | Firefox ships placeholders at `0.54` |
| `font-family: monospace; font-size: 1em` on `code, kbd, samp, pre` | every engine renders the `monospace` default ~3px smaller |
| `[hidden] { display: none }` | author styles outrank the UA stylesheet whatever their layer, so the `img, svg, video { display: block }` rule would otherwise defeat the `hidden` attribute |

Checkbox and radio are deliberately excluded from the text-field rule via `:where()`:
`appearance: none` erases those controls rather than flattening them, and their styling
belongs to `tinywasm/form`.

Two divergences are **out of scope for any reset**: `system-ui` resolves to different
typefaces per platform (SF Pro vs Roboto) with different metrics, and native control
chrome (`<select>` dropdowns, date and file pickers) is owned by the OS.

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

---

## 6. Device Viewport Classes

`Device` is a closed enum of three viewport-width classes used in media queries.
A custom property (`var()`) cannot appear inside an `@media` condition — the
pixel thresholds must be baked into the Go string at build time. The `BpSm` /
`BpMd` / `BpLg` / `BpXl` tokens remain useful for container queries and JS.

### Definition

```go
type Device uint8

const (
    Mobile  Device = iota // (max-width: 639.98px)
    Tablet                // (min-width: 640px) and (max-width: 1023.98px)
    Desktop               // (min-width: 1024px)
)
```

### Query strings

| Class | Media condition | Covers |
|---|---|---|
| `Mobile` | `(max-width: 639.98px)` | 0 – 639.98px |
| `Tablet` | `(min-width: 640px) and (max-width: 1023.98px)` | 640px – 1023.98px |
| `Desktop` | `(min-width: 1024px)` | 1024px+ |

The `.98px` fractions avoid a half-pixel gap where no class matches due to
fractional viewport widths (zoom, scrollbar, device pixel ratio).

The three classes are mutually exclusive and jointly exhaustive: every viewport
width matches exactly one.

### Join function

```go
Query(devices ...Device) string
```

Joins multiple device conditions into an OR-list (`, ` separator). Duplicates
are dropped; order is always Mobile, Tablet, Desktop for deterministic emission.

### Rule: no `var()` in `@media`

A token's value (`--bp-sm: 640px`) cannot be used inside a media query because
custom properties are not resolved at `@media` evaluation time. The literal
pixel values live exclusively in `device.go`. Consumers that need the same
thresholds in a container query or JS should reference the `BpSm`/`BpMd`/`BpLg`/`BpXl`
tokens. Those tokens are NOT a substitute for a `Device` media condition.

### Rule: the two copies of each threshold are guarded, not trusted

`640` and `1024` necessarily appear twice — once as a literal in `device.go` and
once as a `--bp-*` token value — because the media condition cannot read the
token. That duplication is the only sanctioned exception to "never state a value
twice", and it is safe solely because `TestDeviceThresholdsMatchBreakpointTokens`
ties the two together: it parses `BpSm`/`BpLg` and asserts that `Tablet` and
`Desktop` open exactly at those widths, and that `Mobile` and `Tablet` close
`0.02px` below them. Changing a breakpoint token without changing `device.go`
fails the build. Do not delete or weaken that test.
