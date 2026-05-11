# tinywasm/css

Typed CSS DSL and design tokens for the tinywasm framework.

This module replaces string-based `.css` files with a Go-typed DSL. It exposes **both** `RootCSS()` and `RenderCSS()` with strictly separate responsibilities:

- `RootCSS()` → **vocabulary**: design token declarations — brand, source tokens, scales.
- `RenderCSS()` → **logic**: minimal reset + active-token bindings + `@media (prefers-color-scheme)`.

The DSL ensures that every selector, declaration, and token reference is a Go expression, providing compile-time safety and eliminating hex-fallback drift.

## Usage

```go
import . "github.com/tinywasm/css"

func MyComponent() {
    // WASM: Use Class or Token.Var()
    btnClass := Class("btn")
    color := ColorPrimary.Var()
}

// SSR: Use the DSL to generate CSS
func Styles() *Stylesheet {
    return New(
        Rule(".btn",
            BackgroundColor(ColorPrimary),
            Color(Hex("#fff")),
            Padding(Space2, Space4),
            BorderRadius(RadiusMd),
        ),
        Rule(Class("btn").Hover(),
            Opacity(0.8),
        ),
    )
}
```

## SSR contract: `RootCSS` vs `RenderCSS`

`assetmin` recognizes two CSS functions with strictly separate roles:

| Function | Slot | Replacement | Content |
|---|---|---|---|
| `RootCSS() *Stylesheet` | `open` | **Single-winner** — app replaces framework | `:root {}` value declarations (vocabulary) |
| `RenderCSS() *Stylesheet` | `middle` | **Additive** — every module's contribution is preserved | CSS rules that consume tokens via `var()` (logic) |

The split is the key to safe theming: vocabulary is replaceable so apps can rebrand; logic is additive so dark-mode switching cannot be deleted by accident.

### App override pattern

```go
// ssr.go at the app root
import "github.com/tinywasm/css"

func RootCSS() *css.Stylesheet {
    return css.New(
        css.RootCSS(), // inherit framework defaults
        css.Root(
            css.Declare(css.ColorPrimary, "#FF6B35"),
            css.Declare(css.ColorBackgroundLight, "#FAFAFA"),
            css.Declare(css.ColorBackgroundDark, "#121212"),
        ),
    )
}
```

The app does **not** need to redeclare the active-layer bindings or the `@media` rule — those live in `RenderCSS()` and are always present.

---

## Design Tokens

Tokens are the single source of truth for all design decisions.

| Group | Purpose |
|---|---|
| Color — Brand | Fixed identity colors |
| Color — Theme | Adaptive light/dark colors |
| Typography — Size | Font-size scale (Major Third ratio) |
| Typography — Extras | Line-height, weight, letter-spacing |
| Spacing | Margin/padding/gap scale (4px grid) |
| Border-radius | Consistent corner rounding |
| Elevation | Box-shadow scale |
| Motion | Animation timing + easing curves |
| Z-index | Stacking contract |
| Breakpoints | Viewport widths (container queries / JS) |
| Container widths | Max-width primitives |

---

## Design Philosophy

- **Semantic names over values** — `ColorOnSurface` not `#ffffff`. Names describe *intent*; values can change.
- **Scales over magic numbers** — typography and spacing follow mathematical ratios so all values are proportional and limited.
- **Two-layer color pattern** — separates *source* values (per mode) from *active* tokens (used by components). `@media (prefers-color-scheme)` switches modes without JS.
- **Single override point** — apps only need to change source-layer or scale variables; the rest cascades automatically.

---

## DSL Reference

The DSL provides type-safe constructors for CSS properties:

- `BackgroundColor(Value)`, `Color(Value)`, `FontSize(Value)`, etc.
- `Padding(Value...)`, `Margin(Value...)`
- `Px(int)`, `Rem(float64)`, `Pct(int)`, `Hex(string)`, `Str(string)`
- `Rule(selector, declarations...)`
- `Root(declarations...)`
- `Media(query, items...)`

Keywords like `Auto`, `None`, `Block`, `Flex_`, `Center`, `Zero` are also provided.
