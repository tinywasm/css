# PLAN — Extend `tinywasm/css` DSL Coverage (scoped to `platformd`)

> Add the **minimum** set of typed properties, keywords, value constructors, and media helpers so that `layout/platformd/ssr.go` can be implemented with zero `RawRule` calls for properties that have a standard CSS equivalent.

Driver: while implementing `layout/platformd`, the missing DSL coverage forced `RawRule(...)` fallbacks. Each `RawRule` is a leak in the "no CSS-as-strings in Go" guarantee. This plan closes the specific gaps that block `platformd`.

**Scope discipline**: Every addition listed below is justified by one concrete usage in `layout/platformd/docs/PLAN.md` Appendix A. Anything that is *nice to have* but not used by the platform reference CSS is deferred.

---

## 1. Goals & non-goals

### Goals
- Add typed Go functions for every CSS property that `platformd` actually uses (and currently cannot express without `RawRule`).
- Add the keyword values and value constructors required by those properties.
- Add **one** media helper (`MediaDesktop`) for the query repeated 4× in the platform reference CSS.
- Preserve full backward compatibility (additions only).
- Provide one unit test per new function asserting the emitted CSS text.

### Non-goals
- Adding properties / keywords / selectors that no current consumer uses. They can be added incrementally when a real consumer appears. Adding speculative API now would dilute the DSL and slow review.
- Refactoring `dsl.go` internals.
- Splitting `dsl.go` into multiple files (premature for ~15 new functions).
- A CI gate for `RawRule` usage. Code review enforces this; a CI grep can come later if needed.
- Vendor-prefixed properties (`-webkit-*`). Still go through `RawRule` until justified by a non-prefixed-equivalent property is also adopted.

---

## 2. Current coverage (already typed — do not duplicate)

Property functions in `dsl.go`: `Background`, `BackgroundColor`, `BackgroundImage`, `Color`, `Padding`, `Margin`, `Border`, `BorderColor`, `BorderRadius`, `BoxShadow`, `BoxSizing`, `Display`, `Flex`, `FlexDirection`, `Gap`, `JustifyContent`, `AlignItems`, `Width`, `Height`, `MaxWidth`, `MinHeight`, `FontSize`, `FontWeight`, `LineHeight`, `LetterSpacing`, `Transition`, `Animation`, `Transform`, `Cursor`, `Outline`, `OutlineOffset`, `Opacity`, `PointerEvents`, `Position`, `Top`, `Right`, `Bottom`, `Left`, `ZIndex`, `FontFamily`, `Declare`, `Bind`.

Keywords: `Auto`, `None`, `Block`, `Flex_`, `Grid`, `Inline`, `Center`, `Zero`, `Pointer`.

Value constructors: `Px`, `Rem`, `Em`, `Pct`, `Hex`, `Str`.

Selectors: `Selector`, `Class.Hover`, `Class.Focus`, `Class.Disabled`.

Structure helpers: `New`, `Root`, `Rule`, `Media`, `Keyframes`, `RawRule`.

---

## 3. Additions

Each row cites the **exact line** in `layout/platformd/docs/PLAN.md` Appendix A that justifies inclusion. Anything not cited stays out.

### 3.1 Property functions (14)

| Function | CSS property | Justifying usage |
| --- | --- | --- |
| `MinWidth(v Value)` | `min-width` | `menu.css` `min-width: 2.5em` (desktop nav icon) |
| `MaxHeight(v Value)` | `max-height` | `body.css` `max-height: var(--header-height)` (desktop header) |
| `AlignSelf(v Value)` | `align-self` | `body.css` `align-self: flex-end` (mobile header), `align-self: unset` (desktop) |
| `Overflow(v Value)` | `overflow` | `body.css` `overflow: hidden` (desktop root) |
| `Visibility(v Value)` | `visibility` | `user-message.css` `visibility: hidden` (auto-dismiss fade) |
| `TextAlign(v Value)` | `text-align` | `user-message.css` `text-align: right` (USER_AREA), `text-align: center` (toasts) |
| `TextTransform(v Value)` | `text-transform` | `user-message.css` `text-transform: uppercase` (USER_AREA), `capitalize` (USER_NAME) |
| `TextDecoration(v Value)` | `text-decoration` | `body.css` `text-decoration: none` (links) |
| `TextShadow(v ...Value)` | `text-shadow` | `user-message.css` `text-shadow: 0.1em 0.1em 0.1em #ffffff` |
| `UserSelect(v Value)` | `user-select` | `add-default.css` `user-select: none` (global reset) |
| `TouchAction(v Value)` | `touch-action` | `user-message.css` `touch-action: none` (orientation warn overlay) |
| `ListStyleType(v Value)` | `list-style-type` | `add-default.css` `list-style-type: none` (global reset) |
| `GridArea(v Value)` | `grid-area` | `slider-panel.css` `grid-area: module-content` |
| `GridTemplate(v Value)` | `grid-template` | `body.css` multi-line `grid-template: "header header" ... / ...` (used with `Str("…")` for the multi-line literal) |

> **No long-hand `MarginLeft/Right/...` or `PaddingLeft/...`**: the existing `Margin(...)` and `Padding(...)` shorthands cover every Appendix-A usage (e.g. `margin-left: auto` → `Margin(Zero, Zero, Zero, Auto)`, `margin-right: .4rem` → `Margin(Zero, Rem(0.4), Zero, Zero)`). Adding longhand later is trivial if real ergonomic pain emerges.
> **No `Border` longhand split** (`BorderTop`, `BorderStyle`, `BorderWidth`): the platform reference does not need them — the only border-style/border-width rules were on `.foka`/`.ferr`, which are deliberately dropped from the typed port (§A.5 port note).
> **No `BackgroundRepeat/Size/Position`, `Filter`, `URL`**: the only consumers in Appendix A are `body::before` (logo background), which the port note explicitly drops.

### 3.2 Keyword values (12)

```go
var (
    // position
    Fixed    Value = kw("fixed")    // menu.css, user-message.css
    Absolute Value = kw("absolute") // slider-panel.css, user-message.css
    Unset    Value = kw("unset")    // body.css header reset (desktop)
    Initial  Value = kw("initial")  // user-message.css `all: initial`

    // flex / grid alignment
    FlexEnd     Value = kw("flex-end")     // body.css align-self
    SpaceAround Value = kw("space-around") // menu.css navbar

    // flex-direction
    Row    Value = kw("row")    // menu.css mobile
    Column Value = kw("column") // menu.css desktop

    // overflow / visibility
    Hidden  Value = kw("hidden")  // body.css overflow, user-message.css visibility
    Visible Value = kw("visible") // user-message.css keyframes

    // text-transform
    Uppercase  Value = kw("uppercase")  // user-message.css USER_AREA
    Capitalize Value = kw("capitalize") // user-message.css USER_NAME
)
```

> Deferred (no Appendix-A usage): `Sticky`, `Relative`, `Inherit`, `FlexStart`, `SpaceBetween`, `SpaceEvenly`, `Stretch`, `Baseline`, `Start`, `End`, `RowReverse`, `ColumnReverse`, `Scroll`, `Lowercase`, `LeftAlign`, `RightAlign` (text-align uses `Center` keyword and the new `RightAlignKw` only if needed — but Appendix A uses `text-align: right` so add **one** keyword for it: see below), `Justify`, `PreWrap`, `NoWrap`, `Underline`, `Solid`, `Dashed`, `Dotted`, `BorderBoxKw`, `ContentBoxKw`, `InlineBlock`, `InlineFlex`.

Single exception added because `text-align: right` is used and `Right` is already taken by the `Right()` property function:

```go
RightText Value = kw("right") // user-message.css USER_AREA
```

Center is already covered by the existing `Center` keyword.

### 3.3 Value constructors (3)

| Constructor | Output | Justifying usage |
| --- | --- | --- |
| `Vw(n int) Value` | `Nvw` | `100vw`, `5vw`, `95vw` (5+ occurrences across body.css, menu.css, slider-panel.css) |
| `Vh(n int) Value` | `Nvh` | `8vh`, `5vh`, `94vh`, `95vh`, `-100vh` (10+ occurrences) |
| `Calc(expr string) Value` | `calc(expr)` | `user-message.css` `calc(.5em + .5vh)`, `calc(.4em + .4vh)` |

> Deferred (no Appendix-A usage): `Vmin`, `Vmax`, `Ch`, `Rgba`, `RgbaHex`, `URL`, timing-function keyword constants (`Ease`, `EaseInOut`, ...). Transitions in Appendix A read `transition: 0.6s all ease` / `all .6s ease-in-out`; these expand to `Transition(Em(0.6), Str("all"), Str("ease"))` etc. — one `Str("ease")` per call is acceptable until a real frequency emerges.

`Calc` takes a `string` for the inner expression on purpose: building a typed arithmetic AST for CSS `calc()` would dwarf the gain. The escape hatch is **scoped to the arithmetic body only**, not arbitrary CSS.

### 3.4 Media helper (1)

```go
// MediaDesktop wraps the canonical "landscape + hover" media query used by
// tinywasm layouts to distinguish desktop from mobile.
// Reference: appears 4 times verbatim in platformd Appendix A.
func MediaDesktop(items ...item) item {
    return Media("(orientation: landscape) and (hover: hover)", items...)
}
```

> Deferred: `MediaPortrait` (one usage in Appendix A, `Media(raw)` is fine), `MediaMaxWidth` (zero usages; the orientation-warn rule uses a unique compound query handled via raw `Media(...)`).

### 3.5 Selector helpers (0)

**Nothing added.** Every selector pattern in Appendix A is one of:
- a plain class (`Rule(clsFoo, ...)`) — already typed,
- `:hover` — already available via `Class.Hover()`,
- a descendant combinator like `.menu-container:hover .link-text` — expressible as `Selector("."+string(clsMenu)+":hover ."+string(clsLinkText))` in one line.

No new combinator helpers are justified by current consumers. They can be added when the friction is concrete.

### 3.6 `RawRule` policy

Doc comment on `RawRule` updated to:

> `RawRule` is a transitional escape hatch. New code SHOULD use the typed DSL. Each `RawRule` call site SHOULD carry a `// TODO(css-dsl): add typed X` comment naming the missing property, so reviewers can decide whether to extend the DSL or accept the raw use case (vendor-prefixed, exotic property, etc.).

No CI gate. Code review enforces.

---

## 4. File layout

All additions land directly in **`dsl.go`** (and `dsl_test.go` for tests). No new files.

Rationale: adding ~15 functions, 13 keywords, 3 constructors, and 1 helper is well within the scope of `dsl.go`. Splitting into 4 new files is premature organization for a single coherent feature addition. The split is reversible later if the file grows past a natural size threshold.

---

## 5. Stages

| # | Stage | Output | Verify |
| --- | --- | --- | --- |
| 1 | **Keywords & constructors** | Add `Fixed`, `Absolute`, `Unset`, `Initial`, `FlexEnd`, `SpaceAround`, `Row`, `Column`, `Hidden`, `Visible`, `Uppercase`, `Capitalize`, `RightText` to the keyword block. Add `Vw`, `Vh`, `Calc` to value constructors. | Per-keyword test: emit a Rule using each new keyword, assert string contains the keyword text. |
| 2 | **Property functions** | Add the 14 functions from §3.1, each as a one-liner mirroring the existing pattern. | One assertion per function: emit a Rule, assert `property: value;` in the output. |
| 3 | **`MediaDesktop`** | Add the helper from §3.4. | Test that `MediaDesktop(Rule(clsX, Color(ColorPrimary)))` emits `@media (orientation: landscape) and (hover: hover) { .x { color: ... } }`. |
| 4 | **`RawRule` doc** | Update the `RawRule` doc comment per §3.6. | None — doc only. |
| 5 | **README update** | Update `css/README.md` "Supported properties" table to include the additions. | Read-through. |

After Stage 5, `platformd`'s implementation can do a sweep replacing every `RawRule` whose property is now typed. That sweep is **not part of this plan** — it lives in the `platformd` plan to avoid cross-repo coupling in this PR.

---

## 6. Backward compatibility

Pure additions. No existing signatures, types, or keyword names change. The only naming collision (`right` keyword vs `Right()` property function) is resolved by naming the new keyword `RightText` — explicit, single-purpose, no shadow risk.

---

## 7. Acceptance criteria

- All 14 properties, 13 keywords (12 + `RightText`), 3 constructors, and 1 media helper added with passing unit tests.
- `go test ./css/...` green.
- `css/README.md` lists the new surface.
- `platformd`'s upcoming sweep PR can replace every `RawRule(...)` for the properties in §3.1 with typed equivalents. (Verified by `platformd` maintainers; out of this PR's scope.)
- Binary-size impact: negligible (each addition is one `Decl{name, value}` literal or one `kw(...)` declaration).
