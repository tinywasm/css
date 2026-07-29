# Migration Guide: Upgrading to tinywasm/css v0.3.3 (Interaction Derivation API)

This guide provides instructions and mapping specifications to upgrade applications and dependent libraries from the legacy 69-color-token system to the modernized, streamlined 12-token system. This document is optimized for both human reading and direct consumption by Large Language Models (LLMs) to perform automated code migrations.

---

## 1. Architectural Changes

### The 12-Token Paradigm
Previously, `tinywasm/css` maintained separate *active* tokens and *source* tokens for light/dark mode (e.g. `ColorBackground`, `ColorBackgroundLight`, `ColorBackgroundDark`). Responsive theme switching was accomplished via manual variable binding inside `RenderCSS()`.

In `v0.2.0`, **all source tokens are eliminated**. Responsive light/dark theme adaptation is achieved directly in browser space using the CSS `light-dark()` function within the active token fallback value. Manual bindings and prefers-color-scheme rule copying are removed entirely.

---

## 2. API Signature Changes

| Legacy Type / Function | v0.2.0 Type / Function | Notes |
|---|---|---|
| `css.Source` | *Removed* | Completely deleted. All references now use `css.Token`. |
| `css.Set(css.Source, string)` | `css.Set(css.Token, string)` | `css.Set` now directly accepts `css.Token` (e.g. `css.ColorPrimary`). |
| `css.bind(...)` | *Removed* (Unexported) | Variable-to-variable bindings in `RenderCSS()` are obsolete. |
| `css.declareSource(...)` | *Removed* (Unexported) | Source token declaration helper is obsolete. |

---

## 3. Color Token Mapping Reference

Use the following lookup table to replace legacy variables with the new 12 active tokens:

### Brand Group

| Legacy Go Variable | New Go Variable | CSS Custom Property | Default Fallback Value |
|---|---|---|---|
| `ColorPrimary` | `ColorPrimary` | `--color-primary` | `#1b5d8c` |
| `ColorOnPrimary` | `ColorOnPrimary` | `--color-on-primary` | `#FFFFFF` |
| `ColorSecondary` | *Removed* (Use primary or canvas) | — | — |
| `ColorOnSecondary` | *Removed* (Use primary or canvas) | — | — |
| `ColorSuccess` | `ColorSuccess` | `--color-success` | `#1e7a30` |
| `ColorOnSuccess` | `ColorOnSuccess` | `--color-on-success` | `#FFFFFF` |
| `ColorError` | `ColorDanger` | `--color-danger` | `#ba2c0d` |
| `ColorOnError` | `ColorOnDanger` | `--color-on-danger` | `#FFFFFF` |

### Theme/Canvas Group

| Legacy Go Variable | New Go Variable | CSS Custom Property | Light | Dark |
|---|---|---|---|---|
| `ColorBackground` | `ColorBackground` | `--color-background` | `#FFFFFF` | `#0D1117` |
| — | `ColorOnBackground` | `--color-on-background` | `#1C1C1E` | `#E6EDF3` |
| `ColorSurface` | `ColorSurface` | `--color-surface` | `#F2F2F7` | `#161B22` |
| `ColorOnSurface` | `ColorOnSurface` | `--color-on-surface` | `#1C1C1E` | `#E6EDF3` |
| `ColorOutline` | `ColorOutline` | `--color-outline` | `#D1D1D6` | `#30363D` |
| `ColorMuted` | `ColorMuted` | `--color-muted` | `#6E6E73` | `#8B949E` |

### Restored as computed tokens (registered in `:root`)
The following tokens were removed in v0.3.0 and restored in v0.3.3 as computed
values that reference live `var()` expressions instead of hardcoded hex:

| Token | CSS Property | Fallback (v0.3.3) |
|---|---|---|
| `ColorSurfaceSunken` | `--color-surface-sunken` | `color-mix(in oklab, var(--color-surface), var(--color-on-surface) 8%)` |
| `ColorSelection` | `--color-selection` | `color-mix(in oklab, var(--color-primary), transparent 85%)` |
| `ColorOnSelection` | `--color-on-selection` | `var(--color-on-surface)` |

These are declared in `:root` and can be overridden with `css.Set()` like any
other token. Consumers must NOT re-derive them inline — use the token.

### Replaced by CSS API (no token needed)
The following concepts no longer have dedicated tokens but are covered by
`css.Hover/Focus/Press()` functions:

- **27 Interactive Hover/Focus/Press Twins** (`ColorPrimaryHover`…): use
  `css.Hover(css.ColorPrimary)`, `css.Focus(css.ColorSurface)`, etc. The
  derivation (color-mix toward `light-dark(black, white)`) is centralised in
  `css`; component-level `color-mix()` is prohibited.
- **`ColorFocusRing`**: Replace references with `ColorPrimary`.
- **`ColorHover`**: Use `css.Hover(css.ColorSurface)`.
- **`ColorDisabled` / `ColorOnDisabled`**: Use `ColorSurface` / `ColorMuted`.

### Deleted entirely
- **22 Light/Dark Source Twins** (`ColorBackgroundLight`, `ColorBackgroundDark`,
  etc.): Deleted. All mode logic resides inside the main token's `light-dark()`.
- **Source type**: Deleted. Use `Token`.
- **`ColorSecondary` / `ColorOnSecondary`**: Deleted. Use `ColorSurface` / `ColorOnSurface`.
- **`ColorError` / `ColorOnError`**: Renamed to `ColorDanger` / `ColorOnDanger`.

---

## 4. UI Pair Mapping

If you reference design system `css.Pair` variables, use this mapping:

| Legacy Pair | New Pair | New Background | New Foreground |
|---|---|---|---|
| `SurfacePrimary` | `SurfacePrimary` | `ColorPrimary` | `ColorOnPrimary` |
| `SurfacePanel` | `SurfacePanel` | `ColorSurface` | `ColorOnSurface` |
| `SurfaceSunken` | `SurfaceSunken` | `ColorSurfaceSunken` | `ColorOnSurface` |
| `SurfaceSelected` | `SurfaceSelected` | `ColorSelection` | `ColorOnSelection` |
| `SurfaceDanger` | `SurfaceDanger` | `ColorDanger` | `ColorOnDanger` |
| `SurfaceSuccess` | `SurfaceSuccess` | `ColorSuccess` | `ColorOnSuccess` |
| `SurfaceDisabled` | `SurfaceDisabled` | `Token{"--color-disabled", "var(--color-surface)"}` | `Token{"--color-on-disabled", "var(--color-muted)"}` |

---

## 5. Migration Recipes

### Recipe A: Overriding/Theming in the Application

#### Legacy Code:
```go
func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.ColorPrimary, "#FF6B35"),
        css.Set(css.ColorBackgroundLight, "#FAFAFA"),
        css.Set(css.ColorBackgroundDark, "#121212"),
    )
}
```

#### New Code (v0.3.3):
```go
func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.ColorPrimary, "#FF6B35"),
        css.SetTheme(css.ColorBackground, "#FAFAFA", "#121212"),
        css.Set(css.MixHover, "22%"),
    )
}
```

---

### Recipe B: Component Theme Switcher & HTML Integration

Previously, theme-switching required re-binding CSS variables using manual attributes or `@media` queries. Now, the browser does this natively when the `color-scheme` property is manipulated.

#### HTML/JS Implementation:
Ensure your theme switcher toggle continues to apply `data-theme="light"` or `data-theme="dark"` on the `<html>` or `<body>` element.
`RenderCSS()` automatically applies:
```css
[data-theme="light"] {
  color-scheme: light;
}
[data-theme="dark"] {
  color-scheme: dark;
}
```
The browser's CSS rendering engine will immediately resolve all `light-dark()` tokens accordingly. No JavaScript variable injection is required.

---

## 6. Upgrading to v0.4.0 (Device viewport classes + rail widths)

This release is purely additive. No consumer changes are required.

### What was added

- **`Device`** — closed enum (`Mobile`, `Tablet`, `Desktop`) with typed
  `Query()` strings. This is the sanctioned replacement for any hand-written
  `@media` condition in `widget/style` and in components.
- **`Query(devices ...Device) string`** — joins multiple conditions into a
  deterministic OR-list.
- **`RailNarrow` / `RailWide`** — width tokens for sidebar/fixed-column
  layouts, declared in `:root`.

### Migration guide

| Before (hand-written, prohibited going forward) | After |
|---|---|
| `"@media (min-width: 1024px)"` | `css.Desktop.Query()` |
| `"@media (max-width: 639.98px), (min-width: 1024px)"` | `css.Query(css.Mobile, css.Desktop)` |
| hardcoded `3.5rem` rail width | `css.RailNarrow.Var()` |
| hardcoded `12rem` rail width | `css.RailWide.Var()` |

No code changes are forced. New responsive layout code should use `Device`;
existing code can migrate at any time.
