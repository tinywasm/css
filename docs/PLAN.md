# PLAN — one value, one place; one way, not two

Execution document. Steps, reference code, test strategy. **Ephemeral**: not
indexed by `README.md`, and no permanent document links here.

Breaking release. `tinywasm/css` has no published tag beyond `v0.3.0` and its only
consumers are inside this suite, so the cost of breaking is at its minimum.

---

## Development Rules

- **Documentation first.** Update `docs/ARCHITECTURE.md` alongside the code; it
  describes a theming model this plan changes.
- **Dependency direction is one-way.** `widget` may import `css`. **`css` must
  never import `widget`, `widget/style`, or `ssr`** — not for a type, not for a
  test fixture. A test asserts the module graph.
- **One value, one place.** A value is declared exactly once. If the same number
  or colour appears in two files, one of them is wrong and nothing will tell you
  which.
- **One way, not two.** If two APIs can accomplish the same thing, one of them is
  deleted. A choice with no criterion is a defect, not flexibility.
- **css owns values. It does not own component identity.** Anything that names a
  *part of a component* belongs to `widget`.
- **WASM-safe by default.** `tokens.go` carries no build tag and must stay
  importable from a WASM binary; emission logic stays behind `//go:build !wasm`.

---

## 1. Goal

Make the guarantee this library exists to provide actually true. Everything below
follows from one question the review asked and answered: **when two files hold
the same value, which one reaches the screen?**

---

## 2. Findings

All measured against `4c01e6b`, by executing the code and computing the values —
not read off the source.

### F-1 — The contrast guarantee is verified against colours that never render

This is the headline, and it invalidates the reason the rest of the suite depends
on this library.

`contrast_test.go` computes ratios from `Token.Fallback`. `RootCSS()` declares a
*different* value for the same tokens. So the test measures one palette and the
browser paints another:

| pair | tested (`Fallback`) | renders (`:root`) | AA ≥ 4.5 |
|---|---|---|---|
| `SurfacePrimary` | 5.75 : 1 | **3.83 : 1** | fails as shipped |
| `SurfaceSuccess` | 5.40 : 1 | **2.54 : 1** | fails as shipped |
| `SurfaceDanger` | 5.33 : 1 | **3.88 : 1** | fails as shipped |
| `SurfaceDisabledLight` | 4.57 : 1 | **2.60 : 1** | fails as shipped |
| `SurfaceDisabledDark` | 4.95 : 1 | **3.31 : 1** | fails as shipped |

The test is green. Five surfaces fail WCAG AA at the colours users see. White
text on the success surface sits at 2.54 : 1.

`tinywasm/widget` builds its entire "you do not need to know design" premise on
this test. That premise is currently false, and the cause is exactly the
duplication this plan removes.

### F-2 — Twelve tokens hold two different values

The direct cause of F-1. `Token.Fallback` and the `declare(...)` in `RootCSS()`
are two independent sources for one value, and twelve have already drifted:

```
ColorPrimary          fallback #1E6B9E   :root #3f88bf
ColorSuccess          fallback #1F7A31   :root #3FB950
ColorError            fallback #D12200   :root #E34F26
ColorOnDisabledLight  fallback #666666   :root #8E8E93
ColorOnDisabledDark   fallback #8B949E   :root #6E7681
ShadowSm/Md/Lg/Xl     whitespace differences
EaseIn/Out/InOut      whitespace differences
```

The fallback is invisible while `:root` is present, so the drift is silent until
a component renders without the root sheet — at which point the brand colour
changes.

### F-3 — Eleven active tokens are never declared in `:root`

`ColorBackground`, `ColorSurface`, `ColorSurfaceSunken`, `ColorOnSurface`,
`ColorOutline`, `ColorMuted`, `ColorHover`, `ColorSelection`, `ColorOnSelection`,
`ColorDisabled`, `ColorOnDisabled`.

They are `bind()`-ed in `RenderCSS()` and absent from `RootCSS()`. An application
that ships `RootCSS()` alone — which `docs/ARCHITECTURE.md` presents as the
theming entry point — gets components with no theme layer and no dark mode,
falling back to hardcoded hex.

### F-4 — `Class` exists in two packages

`css.Class` and `widget.Class` are the same type with the same two methods:

```go
type Class string
func (c Class) String() string
func (c Class) AsAttr() fmt.KeyValue
```

`css.Class` is referenced by **no consumer** — zero hits across `widget` and
`ssr`. Inside `css` it is reachable only from `rule()`'s type switch and from
`Hover`/`Focus`/`Disabled`, which nothing calls.

It is also in the wrong library. A class name is *component identity*, derived
from a widget's `Name` and `Part`. `css` owns values; a CSS class identifier is
not a value.

**And neither copy is safe.** `widget.go` claims writing `Class("anything")`
outside the package does not compile. It does — `type Class string` permits
conversion from any untyped string constant. Verified:

```
forged class compiles: i-was-never-derived-from-a-Name
and produces an attribute: {class i-was-never-derived-from-a-Name}
```

So the "markup and stylesheet agree by construction" guarantee is currently
decorative. Fixing it is `widget`'s work; deleting the duplicate is this
repository's.

### F-5 — The pair catalogue exists three times

1. `Pair` values in `tokens.go`: `SurfacePrimary`, `SurfacePanel`, …
2. A hand-written list of fifteen pairs inside `contrast_test.go`.
3. `Surface.Resolve()` in `widget/style`, which rebuilds the same couples from
   individual tokens and ignores `Pair` entirely.

Adding a surface means editing three places, and the test list can silently stop
covering the values in use — which is half of how F-1 happened.

### F-6 — `rule()` accepts three kinds of selector

```go
func rule(sel any, content ...ruleContent) item {
    switch v := sel.(type) {
    case Class:    …
    case selector: …
    case string:   …
    }
}
```

Three ways to say the same thing, and `any` defeats the type system at the one
place a stylesheet DSL exists to constrain.

### F-7 — The focus ring is about to be declared twice

`RenderCSS()` already emits a global rule:

```go
rule(selector(":focus-visible"), outline(str("2px solid "+ColorPrimary.Var())), …)
```

`tinywasm/widget` is planning a per-surface focus ring inside `Interactive()`.
Two mechanisms for one visual decision, and they would disagree: a
primary-coloured ring is invisible on the primary surface itself.

### F-8 — Tokens that are unreachable by design

Forty-five tokens are referenced by no consumer. Most are legitimate — the
twenty-two Light/Dark source tokens are `bind()`-ed rather than referenced, and
the breakpoints are for the application shell. But some encode a choice the
suite has already closed:

| Token | Why unreachable |
|---|---|
| `EaseIn`, `EaseOut` | `widget` fixes easing at `--ease-in-out` — "the easing is not chosen" |
| `ShadowXl` | `Elevation` has four steps mapping to `none`/`sm`/`md`/`lg` |
| `TrackingTight/Normal/Wide` | no letter-spacing option exists, in any consumer |
| `LeadingTight`, `LeadingRelaxed` | no line-height option exists; `LeadingNormal` is used by the base reset |
| `MaxWContent`, `MaxWScreen` | no consumer; `MaxWContent` also reads as a twin of `Size.Content`, which means something else entirely |

Keeping a token no typed API can reach is the same defect as a scale step that
resolves to its neighbour: it advertises a choice that does not exist.

---

## 3. Ownership rule

The rule every disposition below follows.

| Concern | Owner | Reason |
|---|---|---|
| What a colour, space, duration or z-level **is** | `css` | It is a value |
| Which foreground pairs with which background | `css` | It is a contrast guarantee about values |
| Light/dark switching, theming, rebrand | `css` | It is value substitution |
| What a **part of a component** is called | `widget` | It is identity, not a value |
| Which token applies to which part in which state | `widget/style` | It is a decision about values |
| Border, radius and padding of a surface | `widget/style` | Component shape, not palette |

`css` never imports `widget`. Today it does not — but it holds two of `widget`'s
concepts (`Class`, and a surface vocabulary), which is the same violation
expressed as duplication instead of an import.

---

## 4. Disposition

### 4.1 Delete

| Symbol | Why |
|---|---|
| `Class`, and its `String()` / `AsAttr()` | Duplicate of `widget.Class`; used by no consumer; component identity, not a value (F-4) |
| `Hover(Class)`, `Focus(Class)`, `Disabled(Class)` | Only callers of `Class`; nothing calls them; `widget/style` builds its own selectors |
| The `case Class` and `case string` arms of `rule()` | One way to name a selector (F-6) |
| `EaseIn`, `EaseOut` | Easing is fixed at `--ease-in-out` by the only consumer (F-8) |
| `ShadowXl` | Unreachable from a four-step `Elevation` (F-8) |
| `TrackingTight`, `TrackingNormal`, `TrackingWide` | No consumer can express letter-spacing (F-8) |
| `LeadingTight`, `LeadingRelaxed` | No consumer can express line-height (F-8) |
| `MaxWContent`, `MaxWScreen` | No consumer; `MaxWContent` reads as a twin of `Size.Content` (F-8) |
| The hand-written pair list in `contrast_test.go` | Third copy of the pair catalogue (F-5) |

### 4.2 Keep, unchanged

Colour actives (11), colour sources (22), brand colours (8), `Space0`–`Space12`,
`Radius*`, `Text*`, `FontWeight*`, `LeadingNormal`, `Shadow` sm/md/lg,
`Duration*`, `EaseInOut`, `Z*` (6), `Bp*` (4), `Token`, `Override`, `Set`,
`Theme`, `Stylesheet`, `NewStylesheet`, `Raw`.

`Z*` and `Bp*` are unreferenced today and stay: `widget` is about to map
`Kind.Layer()` onto `--z-*`, and the breakpoints are the sanctioned shell-level
mechanism now that no component emits a query.

### 4.3 Add — required by `tinywasm/widget`'s release

| Token | For |
|---|---|
| `--color-<family>-hover`, `-focus`, `-press` for the eight interactive families | `Interactive()`; these live in `widget/style` today with hardcoded hex, outside the contrast test |
| `--color-focus-ring` | So the ring is visible on every surface, including `Primary` (F-7) |
| `--column-narrow`, `--column-medium`, `--column-wide` | `Grid` column minimums, hardcoded `rem` today |
| `--max-w-readable`, replacing `--max-w-prose` | A scale step is named after the token it emits; `Size.Readable` must mirror it |

### 4.4 Restructure

**One value, one place (F-2).** `declare` stops taking a value:

```go
// before — two sources of truth
declare(ColorPrimary, "#3f88bf")

// after — the token's own value, declared once in tokens.go
declare(ColorPrimary)
```

`Token.Fallback` becomes *the* value. `RootCSS()` becomes a list of tokens to
emit, not a second place to write numbers. Drift is then unrepresentable rather
than merely tested for.

**Active tokens are declared, not only bound (F-3).** `RootCSS()` emits the
eleven actives with their fallbacks, so a sheet shipped without `RenderCSS()`
still has a complete vocabulary.

**Referenceable and settable become different types.** Today all 91 tokens are
both, so a component can reference `ColorSurfaceDark` directly and silently break
light mode:

```go
type Token struct{ Name, Value string }   // referenceable: has Var()
type Source struct{ Name, Value string }  // settable only: no Var()

func Set(s Source, value string) Override
```

The twenty-two Light/Dark tokens become `Source`. `Set` accepts only `Source`;
`Var()` exists only on `Token`. Referencing a theme source in a rule stops
compiling.

**The pair catalogue becomes the single source (F-5).**

```go
func AllPairs() []NamedPair   // the catalogue, iterated by the contrast test
```

The test iterates it instead of restating it, so a pair that is added is
contrast-checked by construction and one that is removed cannot leave a stale
assertion behind.

**The focus ring keeps one owner (F-7).** The global `:focus-visible` rule in
`RenderCSS()` switches from `ColorPrimary` to `ColorFocusRing`, and
`tinywasm/widget` does **not** emit an outline in `Interactive()` — only the
background change. One mechanism.

---

## 5. Implementation order

Dependency order.

**Step 1 — collapse the two value sources.** `tokens.go` + `css.go`. Change
`declare(Token)` to read `Token.Fallback`; delete every literal from `RootCSS()`.
Where the two disagreed (F-2), **the `:root` value wins** — it is what has been
rendering, so adopting the fallback would change five surfaces' appearance
without anyone asking. Copy those twelve values into `tokens.go` and delete them
from `css.go`.

**Step 2 — fix the contrast failures the merge exposes.** Step 1 makes the test
see the real palette, and five pairs will fail (F-1). Adjust the *foreground* or
the *background* until each clears 4.5 : 1, and record the chosen values. This is
a visual change and needs a human decision — it is the one step in this plan that
is not mechanical.

**Step 3 — declare the actives.** Add the eleven bound-only tokens to
`RootCSS()` (F-3).

**Step 4 — split `Token` and `Source`.** Retype the twenty-two Light/Dark
tokens; `Set` takes `Source`; `Var()` stays on `Token` only.

**Step 5 — delete.** Everything in §4.1, in one commit.

**Step 6 — `rule()` takes a `selector`.** Drop `any` and the three-arm switch
(F-6).

**Step 7 — add the new tokens.** §4.3, contrast-tested with the rest.

**Step 8 — `AllPairs()` and the iterating contrast test.** Extend the catalogue
to the ten families `widget/style` will resolve through.

**Step 9 — the focus ring.** `RenderCSS()` uses `ColorFocusRing`.

**Step 10 — documentation.** `docs/ARCHITECTURE.md` describes a theming model
this plan changes in three ways: `declare` no longer carries values, `Set` no
longer accepts any token, and `RootCSS()` is now complete on its own. Update it,
and index every `docs/` file from `README.md` except this one.

---

## 6. Test strategy

| Test | Asserts | Closes |
|---|---|---|
| `TestContrastFromRenderedValues` | the contrast test iterates `AllPairs()` and computes from the value that reaches `:root` | F-1, F-5 |
| `TestNoValueDeclaredTwice` | no literal value appears in `css.go`; every declaration reads its token | F-2 |
| `TestRootIsComplete` | every referenceable `Token` is declared by `RootCSS()` | F-3 |
| `TestSourceIsNotReferenceable` | compile-time: `Source` has no `Var()`; a `go vet`-style check or a `_ = Source{}.Var()` negative build fixture | F-3 |
| `TestNoWidgetImport` | `go list -deps` for this module contains no `tinywasm/widget` path | ownership rule |
| `TestWasmSafe` | `GOOS=js GOARCH=wasm go build` succeeds for `tokens.go`'s package surface | dev rule |
| existing `TestNoUndeclaredTokensInEmittedCSS` | keep; it becomes meaningful once §4.4 lands | — |

`TestNoValueDeclaredTwice` is the highest-value item: it is the mechanical
enforcement of "one value, one place", and F-1 is what its absence cost.

---

## 7. Coordination

This release **blocks** `tinywasm/widget`'s closed-API release, which needs §4.3.
Ship in this order:

1. `css` steps 1–9 → tag.
2. `widget` picks up the new tokens, deletes its eighteen local ones, drops the
   focus ring from `Interactive()`, and makes `Class` unforgeable — `type Class
   struct{ s string }`, since F-4 shows the current form is not.
3. `ssr` needs nothing from this release.

`widget` also stops emitting the literal `"0"` for its no-space step and
references `Space0` instead — the last invented value on that side.
