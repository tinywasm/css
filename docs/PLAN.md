# PLAN — Typed CSS for tinywasm/css

## Goal

Replace the embed-based `theme.css` + `render.css` model with a Go-typed DSL. `tinywasm/css` becomes the single source of truth for both the token catalog and the SSR rendering API. No `.css` files, no embeds, no code generation.

## Why

The current API is stringly-typed:
- `theme.css` and `render.css` are opaque to the compiler.
- Components reference tokens like `var(--color-primary, #00ADD8)` — renaming a token does not break the build, and the hex fallback duplicates the source value and silently drifts.
- `assetmin` parses `ssr.go` with `go/ast` only because the CSS payload is a string; this extractor is fragile (3 supported forms, every new authoring pattern requires an AST patch).

The DSL eliminates all three problems at the source: every selector, declaration, and token reference is a Go expression.

## Final API (target state)

### Package layout

```
tinywasm/css/
├── tokens.go          // no build tag — Class, Token, all token vars
├── dsl.go             //go:build !wasm — Stylesheet, Rule, Root, properties
├── ssr.go             //go:build !wasm — RootCSS(), RenderCSS()
├── dsl_test.go
└── ssr_test.go
```

### `tokens.go` — shared across SSR and WASM

```go
package css

// Class is a CSS class name. Shared by HTML emission (WASM) and CSS emission (SSR).
// Only the string identity crosses the WASM boundary; pseudo-class helpers
// (Hover/Focus/Disabled) live in dsl.go (!wasm) because they only feed Rule().
type Class string
func (c Class) String() string { return string(c) }

// Token is a design token: a named visual decision with a fallback value.
// Industry-standard term (W3C Design Tokens CG, Material, Carbon, Primer, Spectrum).
type Token struct{ Name, Fallback string }
func (t Token) Var() string { return "var(" + t.Name + "," + t.Fallback + ")" }

// Token catalog — every token from the legacy theme.css.
var (
    // Brand colors
    ColorPrimary    = Token{"--color-primary",     "#00ADD8"}
    ColorOnPrimary  = Token{"--color-on-primary",  "#1C1C1E"}
    ColorSecondary  = Token{"--color-secondary",   "#654FF0"}
    ColorOnSecondary= Token{"--color-on-secondary","#FFFFFF"}
    ColorSuccess    = Token{"--color-success",     "#3FB950"}
    ColorError      = Token{"--color-error",       "#E34F26"}

    // Theme — active layer (consumed by components)
    ColorBackground = Token{"--color-background", "#FFFFFF"}
    ColorSurface    = Token{"--color-surface",    "#F2F2F7"}
    ColorOnSurface  = Token{"--color-on-surface", "#1C1C1E"}
    ColorMuted      = Token{"--color-muted",      "#6E6E73"}
    ColorHover      = Token{"--color-hover",      "#B8860B"}

    // Theme — source layer (apps redeclare these for rebrand)
    ColorBackgroundLight = Token{"--color-background-light", "#FFFFFF"}
    ColorBackgroundDark  = Token{"--color-background-dark",  "#0D1117"}
    ColorSurfaceLight    = Token{"--color-surface-light",    "#F2F2F7"}
    ColorSurfaceDark     = Token{"--color-surface-dark",     "#161B22"}
    ColorOnSurfaceLight  = Token{"--color-on-surface-light", "#1C1C1E"}
    ColorOnSurfaceDark   = Token{"--color-on-surface-dark",  "#E6EDF3"}
    ColorMutedLight      = Token{"--color-muted-light",      "#6E6E73"}
    ColorMutedDark       = Token{"--color-muted-dark",       "#8B949E"}
    ColorHoverLight      = Token{"--color-hover-light",      "#B8860B"}
    ColorHoverDark       = Token{"--color-hover-dark",       "#F7DF1E"}

    // Typography size scale
    TextXs   = Token{"--text-xs",   "0.75rem"}
    TextSm   = Token{"--text-sm",   "0.875rem"}
    TextBase = Token{"--text-base", "1rem"}
    TextLg   = Token{"--text-lg",   "1.25rem"}
    TextXl   = Token{"--text-xl",   "1.5rem"}
    Text2xl  = Token{"--text-2xl",  "2rem"}

    // Line-height / weight / tracking
    LeadingTight       = Token{"--leading-tight",        "1.25"}
    LeadingNormal      = Token{"--leading-normal",       "1.5"}
    LeadingRelaxed     = Token{"--leading-relaxed",      "1.75"}
    FontWeightRegular  = Token{"--font-weight-regular",  "400"}
    FontWeightMedium   = Token{"--font-weight-medium",   "500"}
    FontWeightBold     = Token{"--font-weight-bold",     "700"}
    TrackingTight      = Token{"--tracking-tight",       "-0.02em"}
    TrackingNormal     = Token{"--tracking-normal",      "0"}
    TrackingWide       = Token{"--tracking-wide",        "0.05em"}

    // Spacing (4px grid)
    Space1  = Token{"--space-1",  "0.25rem"}
    Space2  = Token{"--space-2",  "0.5rem"}
    Space3  = Token{"--space-3",  "0.75rem"}
    Space4  = Token{"--space-4",  "1rem"}
    Space6  = Token{"--space-6",  "1.5rem"}
    Space8  = Token{"--space-8",  "2rem"}
    Space12 = Token{"--space-12", "3rem"}

    // Border radius
    RadiusSm   = Token{"--radius-sm",   "4px"}
    RadiusMd   = Token{"--radius-md",   "8px"}
    RadiusLg   = Token{"--radius-lg",   "16px"}
    RadiusFull = Token{"--radius-full", "9999px"}

    // Elevation
    ShadowSm = Token{"--shadow-sm", "0 1px 2px rgba(0,0,0,0.05)"}
    ShadowMd = Token{"--shadow-md", "0 4px 6px rgba(0,0,0,0.1)"}
    ShadowLg = Token{"--shadow-lg", "0 10px 15px rgba(0,0,0,0.1)"}
    ShadowXl = Token{"--shadow-xl", "0 20px 25px rgba(0,0,0,0.15)"}

    // Motion
    DurationFast = Token{"--duration-fast", "150ms"}
    DurationBase = Token{"--duration-base", "250ms"}
    DurationSlow = Token{"--duration-slow", "400ms"}
    EaseIn       = Token{"--ease-in",       "cubic-bezier(0.4,0,1,1)"}
    EaseOut      = Token{"--ease-out",      "cubic-bezier(0,0,0.2,1)"}
    EaseInOut    = Token{"--ease-in-out",   "cubic-bezier(0.4,0,0.2,1)"}

    // Z-index
    ZBase     = Token{"--z-base",     "0"}
    ZDropdown = Token{"--z-dropdown", "100"}
    ZSticky   = Token{"--z-sticky",   "200"}
    ZModal    = Token{"--z-modal",    "300"}
    ZToast    = Token{"--z-toast",    "400"}
    ZTooltip  = Token{"--z-tooltip",  "500"}

    // Breakpoints
    BpSm = Token{"--bp-sm", "640px"}
    BpMd = Token{"--bp-md", "768px"}
    BpLg = Token{"--bp-lg", "1024px"}
    BpXl = Token{"--bp-xl", "1280px"}

    // Container widths
    MaxWProse   = Token{"--max-w-prose",   "65ch"}
    MaxWContent = Token{"--max-w-content", "1200px"}
    MaxWScreen  = Token{"--max-w-screen",  "1440px"}
)
```

### `dsl.go` — SSR-only DSL

```go
//go:build !wasm
package css

import . "github.com/tinywasm/fmt"

type Stylesheet struct{ items []item }
type item interface{ writeTo(*Builder) }

func New(items ...item) *Stylesheet { return &Stylesheet{items} }

func (s *Stylesheet) String() string {
    b := &Builder{}
    for _, it := range s.items { it.writeTo(b) }
    return b.String()
}

// Note: tinywasm/fmt is used exclusively for string building.
// Do NOT import the stdlib "strings" package anywhere in tinywasm/*,
// even under //go:build !wasm — see tinywasm-wide convention.

// Selector is a raw CSS selector string used by the DSL.
// Lives here (not in tokens.go) because no WASM code consumes it.
type Selector string

// Pseudo-class helpers on Class — SSR-only so they don't bloat the WASM binary.
func (c Class) Hover() Selector    { return Selector(string(c) + ":hover") }
func (c Class) Focus() Selector    { return Selector(string(c) + ":focus") }
func (c Class) Disabled() Selector { return Selector(string(c) + ":disabled") }

// Rule constructors — accept any selector source
func Rule(sel any, decls ...Decl) item    { /* normalize Class | Selector | string */ }
func Root(decls ...Decl) item              { /* :root {} */ }
func MediaPrefersDark(items ...item) item  { /* @media (prefers-color-scheme: dark) */ }
func Media(query string, items ...item) item
func Keyframes(name string, frames ...KeyframeStep) item  // KeyframeStep = struct{ At string; Decls []Decl }
func Raw(css string) item                  // escape hatch (last resort)

// Declarations — one constructor per CSS property in active use.
// Audit of the existing .css files yields the minimum surface:
type Decl struct{ Prop, Val string }

func Background(v Value) Decl
func BackgroundColor(v Value) Decl
func BackgroundImage(v Value) Decl
func Color(v Value) Decl
func Padding(v ...Value) Decl
func Margin(v ...Value) Decl
func Border(v ...Value) Decl
func BorderColor(v Value) Decl
func BorderRadius(v ...Value) Decl
func BoxShadow(v Value) Decl
func BoxSizing(v Value) Decl
func Display(v Value) Decl
func Flex(v ...Value) Decl
func FlexDirection(v Value) Decl
func Gap(v Value) Decl
func JustifyContent(v Value) Decl
func AlignItems(v Value) Decl
func Width(v Value) Decl
func Height(v Value) Decl
func MaxWidth(v Value) Decl
func MinHeight(v Value) Decl
func FontSize(v Value) Decl
func FontWeight(v Value) Decl
func LineHeight(v Value) Decl
func LetterSpacing(v Value) Decl
func Transition(v ...Value) Decl
func Animation(v ...Value) Decl
func Transform(v Value) Decl
func Cursor(v Value) Decl
func Outline(v Value) Decl
func Opacity(v float64) Decl
func PointerEvents(v Value) Decl
func Position(v Value) Decl
func Top(v Value) Decl ; func Right(v Value) Decl ; func Bottom(v Value) Decl ; func Left(v Value) Decl
func ZIndex(v Value) Decl

// Value is anything renderable to a CSS value: Token, raw string, number, keyword.
type Value interface{ cssValue() string }
// Token already implements cssValue() via Var()
// Adapters for literals:
func Px(n int) Value
func Rem(f float64) Value
func Em(f float64) Value
func Pct(n int) Value
func Hex(s string) Value
func Str(s string) Value      // arbitrary string when no constructor fits

// Keywords (the small fixed set used today)
var (
    Auto    Value = kw("auto")
    None    Value = kw("none")
    Block   Value = kw("block")
    Flex_   Value = kw("flex")
    Grid    Value = kw("grid")
    Inline  Value = kw("inline-block")
    Center  Value = kw("center")
    Zero    Value = kw("0")
    Pointer Value = kw("pointer")
)

// Root-scope helpers for declaring/binding tokens (only valid inside Root())
func Declare(t Token, value string) Decl   // emits "--color-primary: #00ADD8;"
func Bind(active, source Token) Decl       // emits "--color-X: var(--color-X-light);"
```

### `ssr.go` — replaces theme.css + render.css

```go
//go:build !wasm
package css

func RootCSS() *Stylesheet {
    return New(
        Root(
            // Brand
            Declare(ColorPrimary,      "#00ADD8"),
            Declare(ColorOnPrimary,    "#1C1C1E"),
            Declare(ColorSecondary,    "#654FF0"),
            Declare(ColorOnSecondary,  "#FFFFFF"),
            Declare(ColorSuccess,      "#3FB950"),
            Declare(ColorError,        "#E34F26"),
            // Source layer (light)
            Declare(ColorBackgroundLight, "#FFFFFF"),
            Declare(ColorSurfaceLight,    "#F2F2F7"),
            Declare(ColorOnSurfaceLight,  "#1C1C1E"),
            Declare(ColorMutedLight,      "#6E6E73"),
            Declare(ColorHoverLight,      "#B8860B"),
            // Source layer (dark)
            Declare(ColorBackgroundDark,  "#0D1117"),
            Declare(ColorSurfaceDark,     "#161B22"),
            Declare(ColorOnSurfaceDark,   "#E6EDF3"),
            Declare(ColorMutedDark,       "#8B949E"),
            Declare(ColorHoverDark,       "#F7DF1E"),
            // Typography, spacing, radius, shadows, motion, z-index, breakpoints, widths
            // ... full vocabulary
        ),
    )
}

func RenderCSS() *Stylesheet {
    return New(
        // Active = light defaults
        Root(
            Bind(ColorBackground, ColorBackgroundLight),
            Bind(ColorSurface,    ColorSurfaceLight),
            Bind(ColorOnSurface,  ColorOnSurfaceLight),
            Bind(ColorMuted,      ColorMutedLight),
            Bind(ColorHover,      ColorHoverLight),
        ),
        // OS dark mode
        MediaPrefersDark(
            Root(
                Bind(ColorBackground, ColorBackgroundDark),
                Bind(ColorSurface,    ColorSurfaceDark),
                Bind(ColorOnSurface,  ColorOnSurfaceDark),
                Bind(ColorMuted,      ColorMutedDark),
                Bind(ColorHover,      ColorHoverDark),
            ),
        ),
        // Minimal reset
        Rule(Selector("*,*::before,*::after"), BoxSizing(Str("border-box"))),
        Rule(Selector("body"),
            Margin(Zero),
            FontSize(TextBase),
            LineHeight(LeadingNormal),
            BackgroundColor(ColorBackground),
            Color(ColorOnSurface),
        ),
        Rule(Selector(":focus-visible"), Outline(Str("2px solid "+ColorPrimary.Var()))),
        Rule(Selector("img,svg,video"),
            Display(Block),
            MaxWidth(Pct(100)),
        ),
    )
}
```

## Files removed

- `tinywasm/css/theme.css`
- `tinywasm/css/render.css`
- Existing `ssr.go` embed wiring (rewritten as DSL).

## Files added

- `tokens.go` (no build tag)
- `dsl.go` (`//go:build !wasm`)
- New `ssr.go` returning `*Stylesheet`.
- `dsl_test.go` — verifies String() output of every constructor.
- `ssr_test.go` — golden-file test that `RootCSS().String()` and `RenderCSS().String()` produce byte-equivalent output to the legacy `theme.css` / `render.css` (modulo whitespace).

## Steps

1. Snapshot current `theme.css` + `render.css` rendered output (used as golden in step 6).
2. Implement `tokens.go` — straight transcription of `theme.css` values.
3. Implement `dsl.go` — start with `Stylesheet`, `Rule`, `Root`, `Decl`, `Value`, then add property constructors one-by-one.
4. Implement `Declare` and `Bind` helpers for Root scope.
5. Rewrite `ssr.go` using DSL only.
6. Add golden test comparing DSL output to the snapshot from step 1.
7. Delete `theme.css`, `render.css`.
8. Update `README.md`: remove all references to embed files, theme.css, render.css; add DSL examples; keep the SSR contract section (single-winner RootCSS, additive RenderCSS) intact — it is API-level, not implementation-level.

## Acceptance

- `go test ./...` passes including golden equivalence with previous CSS output.
- `theme.css` and `render.css` are deleted.
- `RootCSS()` and `RenderCSS()` return `*Stylesheet`.
- `Class`, `Token`, `Selector` are usable from WASM builds (no build tag).
- `Stylesheet`, `Rule`, `Decl`, property constructors are SSR-only (`!wasm`).
- A consumer can write `import . "github.com/tinywasm/css"` and use everything unprefixed.

## Out of scope

- Container queries (`@container`) — not used today; add when first consumer needs it.
- CSS layers (`@layer`) — not used today.
- Animations beyond `@keyframes pulse-url` (move to its owning component).
