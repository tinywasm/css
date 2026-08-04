# tinywasm/css
<img src="docs/img/badges.svg">

Typed CSS design tokens and emission engine for the tinywasm framework.

This module acts as the single source of truth for design decisions and theme construction. It replaces string-based `.css` files with Go-typed design tokens, exposing **both** `RootCSS()` and `RenderCSS()` with strictly separate responsibilities:

- `RootCSS()` → **vocabulary**: design token declarations — brand, source tokens, scales.
- `RenderCSS()` → **logic**: minimal reset + active-token bindings + `@media (prefers-color-scheme)`.

The internal DSL ensures that every selector, declaration, and token reference is a Go expression, providing compile-time safety and eliminating hex-fallback drift.

## Public API & Architecture

To prevent design-system evasion, the low-level CSS properties and free-value constructors have been unexported from the public surface.

The public API consists solely of:
- `Token` and `Pair` design tokens.
- The design token catalog (e.g., `ColorPrimary`, `ColorSurface`, `Space2`, `MixHover`, etc.).
- `Hover(t)`, `Focus(t)`, `Press(t)` — interaction derivations via `color-mix()`.
- `Device`, `Mobile`, `Tablet`, `Desktop`, `Query(...)` — typed viewport classes for media queries.
- `Stylesheet` and `NewStylesheet` for compilation.
- `Theme`, `Set`, `SetTheme` for rebranded app themes.
- `RootCSS`, `RenderCSS`, `FontFaces`, `Raw` — CSS emitters (`FontFaces` is separate from the SSR root/render slots).

All component styling is expressed using the semantic visual intention API in `github.com/tinywasm/widget/style`, which compiles down to CSS rules.

### Component Styling Example

Components do not write raw CSS or lower-level property functions. Instead, they express intent using high-level layouts (`Stack`, `Row`, `Split`), semantic surfaces (`On(Page)`, `On(Panel)`), and scale exceptions (`Fill()`, `Round()`, `Pad()`):

```go
package targetlist

import (
	"github.com/tinywasm/widget"
	"github.com/tinywasm/widget/style"
)

const (
	nameTargetList = widget.Name("targetlist")
	partRow        = widget.Part("row")
)

func (l *TargetList) WidgetName() widget.Name { return nameTargetList }
func (l *TargetList) WidgetKind() widget.Kind { return widget.Listbox }

func (l *TargetList) Style() *style.Sheet {
	return style.Of(nameTargetList).
		Root(style.Stack(style.Space1), style.On(style.Sunken), style.Scrolls()).
		Part(partRow, style.Row(style.Space2), style.On(style.Panel), style.Pad(style.Space2), style.Round(style.RadiusSm)).
		When(widget.Selected, partRow, style.On(style.Selected))
}
```

---

## SSR contract: `RootCSS` vs `RenderCSS`

`assetmin` recognizes two CSS functions with strictly separate roles:

| Function | Slot | Replacement | Content |
|---|---|---|---|
| `RootCSS() *Stylesheet` | `open` | **Single-winner** — app replaces framework | `:root {}` value declarations (vocabulary) |
| `RenderCSS() *Stylesheet` | `middle` | **Additive** — every module's contribution is preserved | CSS rules that consume tokens via `var()` (logic) |

The split is the key to safe theming: vocabulary is replaceable so apps can rebrand; logic is additive so dark-mode switching cannot be deleted by accident.

### Theming an App (Rebrand)

To apply a theme or rebrand to an application, the root project exposes its own `RootCSS()`. Because `assetmin` treats the `:root` block as a **single-winner slot**, the app's `RootCSS()` completely replaces the library defaults.

```go
// config/css.go in the application (!wasm)
import "github.com/tinywasm/css"

func RootCSS() *css.Stylesheet {
    return css.Theme(
        css.Set(css.ColorPrimary, "#FF6B35"),
        css.SetTheme(css.ColorBackground, "#FAFAFA", "#121212"),
        css.SetTheme(css.ColorSurface, "#F2F2F7", "#161B22"),
        css.Set(css.MixHover, "22%"),
    )
}
```

`Theme()` returns a stylesheet with the token declarations that the app needs to overwrite, appended at the end of the default catalog block. This ensures that the app's branding overrides cascade correctly.

The app **does not** need to redeclare active layer bindings (`--color-surface`, etc.) or `@media (prefers-color-scheme)` logic; those reside in `RenderCSS()` and remain always active.

| Action | API |
|---|---|---|
| Rebrand brand color | `css.Set(css.ColorPrimary, "#hex")` |
| Change background (theme pair) | `css.SetTheme(css.ColorBackground, "#light", "#dark")` |
| Adjust interaction intensity | `css.Set(css.MixHover, "22%")` |
| Adjust global border-radius | `css.Set(css.RadiusMd, "12px")` |
| Adjust typographic scale | `css.Set(css.TextBase, "1.1rem")` |

---

## Design Tokens

Tokens are the single source of truth for all design decisions.

| Group | Purpose |
|---|---|
| Color — Brand | Fixed identity colors (e.g. `ColorPrimary`, `ColorOnPrimary`) |
| Color — Theme | Adaptive light/dark active layers (e.g. `ColorBackground`, `ColorSurface`) |
| Color — Pairs | Coupled background/foreground surface decisions (e.g. `SurfacePanel`, `SurfaceSunken`) |
| Color — Computed | Derived values via `color-mix()` (e.g. `ColorSurfaceSunken`, `ColorSelection`) |
| Interaction — Intensity | Mix-percentage scale for Hover/Focus/Press derivations (`MixHover`, `MixFocus`, `MixPress`) |
| Typography — Size | Font-size scale (Major Third ratio) |
| Typography — Extras | Line-height, weight, letter-spacing |
| Spacing | Margin/padding/gap scale (4px grid) |
| Border-radius | Consistent corner rounding |
| Elevation | Box-shadow scale |
| Motion | Animation timing + easing curves |
| Z-index | Stacking contract |
| Viewport classes | Typed media-query enum (`Mobile`, `Tablet`, `Desktop`) for responsive layout |
| Breakpoints | Viewport widths (`BpSm`/`BpMd`/`BpLg`/`BpXl` — container queries / JS only, not usable in `@media`) |
| Container widths | Max-width primitives |
| Rail widths | Sidebar/fixed-column scale (`RailNarrow`, `RailWide`) |

---

## Design Philosophy

- **Semantic names over values** — `ColorOnSurface` not `#ffffff`. Names describe *intent*; values can change.
- **Contrast safety by type** — Background/foreground design decisions are coupled in `Pair` structures. Contrast ratio tests guarantee compliance with WCAG >= 4.5:1.
- **Scales over magic numbers** — Typography and spacing follow mathematical ratios so all values are proportional and limited.
- **Two-layer color pattern** — Separates *source* values (per mode) from *active* tokens (used by components). `@media (prefers-color-scheme)` switches modes without JS.
- **Single override point** — Apps only need to change source-layer or scale variables; the rest cascades automatically.

## Documentation

- [AGENTS.md](AGENTS.md) — Constraints for anyone (human or agent) changing this library: the WASM boundary, the no-duplicated-value rule, and the checklist for adding a token.
- [Technical Specifications](docs/SPECS.md) — Detailed specifications of all active design tokens, browser-level light-dark handling, and the complete library contract.
- [Migration Guide](docs/MIGRATION.md) — Upgrade paths between releases.
- [Theming Architecture](docs/ARCHITECTURE.md) — Detailed guide on RootCSS, single-winner slot, theme overrides, and type safety constraints.
- [Go-Typed CSS DSL Justification](docs/JUSTIFICACION_DSL.md) — Deep-dive on design rationale, technical constraints, and comparison with alternatives.
