# tinywasm/css

CSS design tokens for the tinywasm framework, served via Go embed.

> **En análisis:** propuesta de DSL tipado para reemplazar el CSS embebido. Ver justificación en español: [`docs/JUSTIFICACION_DSL.md`](./docs/JUSTIFICACION_DSL.md) — y plan técnico en inglés: [`docs/PLAN_typed_css.md`](./docs/PLAN_typed_css.md).

This module exposes **both** `RootCSS()` and `RenderCSS()` with strictly separate responsibilities:

- `RootCSS()` → **vocabulary** (`theme.css`): static value declarations — brand, source tokens, scales.
- `RenderCSS()` → **switching logic** (`render.css`): active-token bindings + `@media (prefers-color-scheme)`.

The split exists so apps can replace the vocabulary without breaking the automatic dark-mode wiring. See [SSR contract](#ssr-contract-rootcss-vs-rendercss) below.

## SSR contract: `RootCSS` vs `RenderCSS`

`assetmin` recognizes two CSS functions with strictly separate roles:

| Function | Slot | Replacement | Content |
|---|---|---|---|
| `RootCSS() string` | `open` | **Single-winner** — app replaces framework | `:root {}` value declarations (vocabulary) |
| `RenderCSS() string` | `middle` | **Additive** — every module's contribution is preserved | CSS rules that consume tokens via `var()` (logic) |

The split is the key to safe theming: vocabulary is replaceable so apps can rebrand; logic is additive so dark-mode switching cannot be deleted by accident.

### What goes in `RootCSS()`

Pure value declarations only:

- Brand colors (`--color-primary`, `--color-success`, …)
- Source-layer theme tokens (`--color-background-light`, `--color-background-dark`, …)
- Typography, spacing, radius scales

### What goes in `RenderCSS()`

`RenderCSS` is a **slot identifier**, not a content-type description. Across the ecosystem it covers three distinct families of CSS — all unified by the same contract: *additive, injected after `RootCSS`, consume tokens via `var()`*.

| Family | Example | Aported by |
|---|---|---|
| **Token bindings** | `:root { --color-background: var(--color-background-light); }` | `tinywasm/css` |
| **Mode switching** | `@media (prefers-color-scheme: dark) { :root { … } }`, `[data-theme="dark"] { … }` | `tinywasm/css` (OS), `themeswitch` (manual) |
| **Component rules** | `.card { padding: var(--space-4); }` | Any component (`button`, `card`, `modal`, …) |

The name leans toward the third family (where `RenderCSS` most literally "renders elements"), but the slot's defining property is **how it is injected**, not what it visually does. Read `RenderCSS()` as "*everything that is not the single canonical `RootCSS` block*".

`tinywasm/css` owns the OS-preference branch here. Manual `[data-theme="…"]` switching lives in [`tinywasm/components/themeswitch`](../components/themeswitch/) — that's a separate `RenderCSS()` contribution from that component.

### App override pattern

```go
// ssr.go at the app root
import "github.com/tinywasm/css"

func RootCSS() string {
    return css.RootCSS() + `
    :root {
      --color-primary:          #FF6B35;
      --color-background-light: #FAFAFA;
      --color-background-dark:  #121212;
    }`
}
```

The app does **not** need to redeclare the active-layer bindings or the `@media` rule — those live in `RenderCSS()` and are always present.

#### Why replacement (and not merge) for RootCSS

A merge approach (injecting both `RootCSS` blocks and letting the cascade resolve) silently breaks `@media (prefers-color-scheme)` whenever the app redeclares an *active* token outside of any media query — the unconditional override lands after the framework's `@media` and wins in dark mode too. Single-winner replacement keeps the cascade order entirely the app's responsibility, while keeping the switching wiring safe in `RenderCSS`.

> `RootCSS` may only be declared in **`tinywasm/css`** or the **app root**. Other modules that declare it are ignored with a warning.

`theme.css` is the **single source of truth** for all design tokens. Eleven groups:

| Group | Tokens | Purpose |
|---|---|---|
| Color — Brand | `--color-primary/secondary/…` | Fixed identity colors |
| Color — Theme | `--color-background/surface/…` | Adaptive light/dark colors |
| Typography — Size | `--text-xs` … `--text-2xl` | Font-size scale (Major Third ratio) |
| Typography — Extras | `--leading-*`, `--font-weight-*`, `--tracking-*` | Line-height, weight, letter-spacing |
| Spacing | `--space-1` … `--space-12` | Margin/padding/gap scale (4px grid) |
| Border-radius | `--radius-sm` … `--radius-full` | Consistent corner rounding |
| Elevation | `--shadow-sm` … `--shadow-xl` | Box-shadow scale |
| Motion | `--duration-*`, `--ease-*` | Animation timing + easing curves |
| Z-index | `--z-dropdown/sticky/modal/toast/tooltip` | Stacking contract |
| Breakpoints | `--bp-sm/md/lg/xl` | Viewport widths (container queries / JS) |
| Container widths | `--max-w-prose/content/screen` | Max-width primitives |

---

## Design Philosophy

- **Semantic names over values** — `--color-on-surface` not `--color-white`. Names describe *intent*; values can change.
- **Scales over magic numbers** — typography and spacing follow mathematical ratios so all values are proportional and limited.
- **Two-layer color pattern** — separates *source* values (per mode) from *active* tokens (used by components). `@media (prefers-color-scheme)` switches modes without JS. Manual switching via `[data-theme]` is handled by `tinywasm/components/themeswitch` — not this module.
- **Single override point** — apps only need to change source-layer or scale variables; the rest cascades automatically.

---

## Color Tokens

### Brand Group (Fixed)

Identity colors that never flip between light and dark modes.

| Token | Value | Semantics |
|---|---|---|
| `--color-primary` | `#00ADD8` | Go cyan — primary brand color |
| `--color-on-primary` | `#1C1C1E` | Content (text/icon) displayed ON primary |
| `--color-secondary` | `#654FF0` | WASM purple — interactive accent |
| `--color-on-secondary` | `#FFFFFF` | Content displayed ON secondary |
| `--color-success` | `#3FB950` | Go gopher green |
| `--color-error` | `#E34F26` | HTML5 orange-red |

### Theme Group (Adaptive)

Tokens that change value based on active mode (light/dark).

| Token | Semantics | Light | Dark |
|---|---|---|---|
| `--color-background` | Page background | `#FFFFFF` | `#0D1117` |
| `--color-surface` | Panel / Card background | `#F2F2F7` | `#161B22` |
| `--color-on-surface` | Main text color | `#1C1C1E` | `#E6EDF3` |
| `--color-muted` | Secondary text / subtle borders | `#6E6E73` | `#8B949E` |
| `--color-hover` | Interactive hover state | `#B8860B` | `#F7DF1E` |

### Two-Layer Pattern

Color customization uses two layers to keep mode-switching logic in the framework:

1. **Source layer** — `--color-X-light` / `--color-X-dark` → apps override *these*.
2. **Active layer** — `--color-X` → what components reference; assigned automatically.

```css
/* in tinywasm/css — automatic, no JS required */
:root                              { --color-background: var(--color-background-light); }
@media (prefers-color-scheme: dark){ --color-background: var(--color-background-dark);  }

/* in tinywasm/components/themeswitch — only when that component is used */
[data-theme="dark"]                { --color-background: var(--color-background-dark);  }
[data-theme="light"]               { --color-background: var(--color-background-light); }
```

**App override rule** — to keep `@media`-based mode switching alive, redeclare the **source layer** tokens (`--color-X-light` / `--color-X-dark`), never the active ones. Combine with the `css.RootCSS() + override` pattern shown in the [SSR contract](#ssr-contract-rootcss-vs-rendercss) section.

---

## Typography Scale

Follows a **Major Third ratio (×1.25)** starting at `1rem` (browser default = 16px).
Use these for all `font-size` declarations — never hardcode pixel values.

| Token | Value | px | Typical use |
|---|---|---|---|
| `--text-xs` | `0.75rem` | 12px | Captions, labels, badges |
| `--text-sm` | `0.875rem` | 14px | Secondary / helper text |
| `--text-base` | `1rem` | 16px | Body copy |
| `--text-lg` | `1.25rem` | 20px | Lead text, card titles |
| `--text-xl` | `1.5rem` | 24px | Section headings |
| `--text-2xl` | `2rem` | 32px | Page / hero headings |

```css
.card-title { font-size: var(--text-lg); }
.caption    { font-size: var(--text-xs); color: var(--color-muted); }
```

---

## Spacing Scale

Built on a **4px grid** (multiples of `0.25rem`). Use for `margin`, `padding`, and `gap`.

| Token | Value | px | Common use |
|---|---|---|---|
| `--space-1` | `0.25rem` | 4px | Icon nudge, tight gaps |
| `--space-2` | `0.5rem` | 8px | Inline padding, small gaps |
| `--space-3` | `0.75rem` | 12px | Input padding, list items |
| `--space-4` | `1rem` | 16px | Card padding, section gaps |
| `--space-6` | `1.5rem` | 24px | Between sections |
| `--space-8` | `2rem` | 32px | Large section padding |
| `--space-12` | `3rem` | 48px | Page-level vertical rhythm |

```css
.card   { padding: var(--space-4); gap: var(--space-3); }
.button { padding: var(--space-2) var(--space-4); }
```

---

## Border-Radius Scale

| Token | Value | Typical use |
|---|---|---|
| `--radius-sm` | `4px` | Inputs, tags |
| `--radius-md` | `8px` | Cards, buttons |
| `--radius-lg` | `16px` | Modals, panels |
| `--radius-full` | `9999px` | Pills, avatars |

---

## Typography — Line-Height / Weight / Tracking

Complement the size scale. Use for `line-height`, `font-weight`, `letter-spacing`.

| Token | Value | Typical use |
|---|---|---|
| `--leading-tight` | `1.25` | Headings, dense UI |
| `--leading-normal` | `1.5` | Body copy (default in reset) |
| `--leading-relaxed` | `1.75` | Long-form reading |
| `--font-weight-regular` | `400` | Body |
| `--font-weight-medium` | `500` | Emphasis without bold |
| `--font-weight-bold` | `700` | Headings, strong emphasis |
| `--tracking-tight` | `-0.02em` | Large display text |
| `--tracking-normal` | `0` | Default |
| `--tracking-wide` | `0.05em` | Small-caps, all-caps labels |

---

## Elevation (Box-Shadow)

| Token | Typical use |
|---|---|
| `--shadow-sm` | Subtle lift — inputs, hover hints |
| `--shadow-md` | Cards, dropdowns |
| `--shadow-lg` | Modals, popovers |
| `--shadow-xl` | Hero panels, side sheets |

```css
.card { box-shadow: var(--shadow-md); }
```

---

## Motion (Duration + Easing)

| Token | Value | Typical use |
|---|---|---|
| `--duration-fast` | `150ms` | Hover, focus rings |
| `--duration-base` | `250ms` | Generic UI transitions |
| `--duration-slow` | `400ms` | Page-level animations |
| `--ease-in` | `cubic-bezier(0.4, 0, 1, 1)` | Element leaving the screen |
| `--ease-out` | `cubic-bezier(0, 0, 0.2, 1)` | Element entering the screen |
| `--ease-in-out` | `cubic-bezier(0.4, 0, 0.2, 1)` | State-to-state changes |

```css
.button { transition: background var(--duration-fast) var(--ease-out); }
```

---

## Z-Index — Stacking Contract

Avoid arbitrary `z-index` numbers. All layered components must reference these tokens so stacking conflicts become impossible.

| Token | Value | Owner |
|---|---|---|
| `--z-base` | `0` | Default flow |
| `--z-dropdown` | `100` | Dropdowns, select menus |
| `--z-sticky` | `200` | Sticky headers, sidebars |
| `--z-modal` | `300` | Modals, dialogs |
| `--z-toast` | `400` | Notifications, toasts |
| `--z-tooltip` | `500` | Tooltips (always on top) |

---

## Breakpoints

| Token | Value | Class of device |
|---|---|---|
| `--bp-sm` | `640px` | Phones (landscape) |
| `--bp-md` | `768px` | Tablets |
| `--bp-lg` | `1024px` | Laptops |
| `--bp-xl` | `1280px` | Desktops |

> **CSS limitation:** `@media` rules cannot read CSS custom properties. For viewport media queries, still write the literal pixel value. The tokens remain useful for:
> - `@container (min-width: var(--bp-md)) { … }` — container queries
> - JS (`getComputedStyle(document.documentElement).getPropertyValue('--bp-md')`)
> - Documentation of the canonical breakpoints

---

## Container Widths

| Token | Value | Typical use |
|---|---|---|
| `--max-w-prose` | `65ch` | Readable line length for body text |
| `--max-w-content` | `1200px` | Main content column |
| `--max-w-screen` | `1440px` | Full-page max width |

```css
.article { max-width: var(--max-w-prose); margin-inline: auto; }
```

---

## Framework Reset (in `render.css`)

`RenderCSS()` ships a **minimal reset** so apps start from a sane baseline without re-inventing it:

- `*, *::before, *::after { box-sizing: border-box; }`
- `body` — zero margin, system font stack, body font-size + line-height bound to tokens, surface colors applied.
- `:focus-visible` — visible focus ring on `--color-primary`.
- `img, svg, video` — `display: block; max-width: 100%`.

This is intentionally small. Apps that want a full reset (Eric Meyer, Normalize, modern-normalize) layer it on top from their own `RenderCSS()`.

---

## Convention: `@font-face`

There is no dedicated method for font loading. Declare `@font-face` in the **app root's `RenderCSS()`** (or in a dedicated component if a font is component-scoped). The framework intentionally does not bundle fonts — choosing a font family is a product decision.

```go
// ssr.go in your app root
func RenderCSS() string {
    return `
    @font-face {
      font-family: "Inter";
      src: url("/fonts/inter-var.woff2") format("woff2-variations");
      font-weight: 100 900;
      font-display: swap;
    }
    :root { --font-sans: "Inter", system-ui, sans-serif; }
    body  { font-family: var(--font-sans); }
    `
}
```

Note the optional `--font-sans` token pattern — once exposed at `:root`, every component can switch from `system-ui` to the loaded font automatically.

