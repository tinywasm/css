# PLAN — Extend `tinywasm/css` DSL: Round 2 + `New` → `NewStylesheet` rename

> Two goals: (1) close the remaining `RawRule` gaps from `layout/platformd` and `layout/rightpanel`; (2) rename `New` → `NewStylesheet` to eliminate the dot-import name collision risk.

---

## 0. Breaking change — rename `New` → `NewStylesheet`

### Problem

`tinywasm/css` exports `New(items ...item) *Stylesheet`. When consumers use dot imports (`. "github.com/tinywasm/css"`), the name `New` is injected into the package namespace. This is a generic identifier that collides silently with:

- Constructors named `New` in the same package or other dot-imported packages (e.g. `tinywasm/dom`, `tinywasm/fmt`).
- Any local variable or function also named `New`.
- Future packages added to the dot-import chain.

The collision is currently latent but guaranteed to surface as the ecosystem grows. The fix must happen before the API solidifies further.

### Solution

Rename the constructor to `NewStylesheet` — unambiguous, self-describing, collision-free.

```go
// Before
func New(items ...item) *Stylesheet

// After
func NewStylesheet(items ...item) *Stylesheet
```

### Blast radius

Every `ssr.go` that calls `New(...)` must be updated. Known callers (confirmed by grep):

| File | Change |
| --- | --- |
| `tinywasm/css/ssr.go` (internal, 2 calls) | `New(` → `NewStylesheet(` |
| `tinywasm/css/dsl_test.go` (6 calls) | same |
| `layout/rightpanel/ssr.go` | same |
| `layout/platformd/ssr.go` | same |
| `tinywasm/components/contentcard/ssr.go` | same |
| `tinywasm/components/themetoggle/ssr.go` | same |
| `tinywasm/components/datatable/ssr.go` | same |
| `tinywasm/components/actionbutton/ssr.go` | same |
| `tinywasm/components/selectsearch/ssr.go` | same |
| `tinywasm/components/dialog/ssr.go` | same |
| `tinywasm/components/navbar/ssr.go` | same |

> The agent executing this plan MUST grep the entire tinywasm monorepo for `\bNew(` in `ssr.go` files before proceeding, since new components may have been added since this plan was written.
>
> ```bash
> grep -rn "\bNew(" --include="ssr.go" /path/to/tinywasm/
> ```

### No deprecation shim

Do **not** keep `New` as an alias (`var New = NewStylesheet`). An alias defeats the entire purpose — it keeps the collision in place. This is a clean rename; callers update in the same PR or immediately following PRs per module. Since all callers are in the same monorepo, the migration is atomic.

### Version bump

This rename is a **major breaking change** in semver terms. The CSS package version MUST be bumped to `v0.1.0` (or higher) to signal incompatibility. Do not publish as a patch.

---

Previous round (v0.0.5) added: `MinWidth`, `MaxHeight`, `AlignSelf`, `Overflow`, `Visibility`, `TextAlign`, `TextTransform`, `TextDecoration`, `TextShadow`, `UserSelect`, `TouchAction`, `ListStyleType`, `GridArea`, `GridTemplate`, plus keywords `Fixed/Absolute/Unset/Initial/FlexEnd/SpaceAround/Row/Column/Hidden/Visible/Uppercase/Capitalize/RightText` and constructors `Vw`, `Vh`, `Calc`.

This round closes what that sweep revealed as still missing.

---

## 1. Gap inventory

Sourced from `grep -n "RawRule" layout/platformd/ssr.go layout/rightpanel/ssr.go` after v0.0.5, excluding vendor-prefixed properties (those remain as `RawRule` permanently).

### From `layout/platformd/ssr.go`

| Line(s) | Raw CSS | Typed addition needed |
| --- | --- | --- |
| 44 | `list-style: none` | `ListStyle(v Value)` |
| 69, 240, 253, 304, 309, 321, 326 | `margin-left: <v>` | `MarginLeft(v Value)` |
| 327 | `margin-right: .4rem` | `MarginRight(v Value)` |
| 245 | `grid-template-columns: 1fr 3fr 1fr` | `GridTemplateColumns(v Value)` |
| 336 | `all: initial` | `All(v Value)` |

### From `layout/rightpanel/ssr.go`

| Line(s) | Raw CSS | Typed addition needed |
| --- | --- | --- |
| 46, 54, 109 | `overflow: hidden` | `Overflow` already in v0.0.5 — rightpanel just needs to call it (no new DSL needed) |
| 88, 122, 151, 154 | `overflow-y: auto / visible` | `OverflowY(v Value)` |
| 51, 106 | `grid-template-rows: auto 1fr` | `GridTemplateRows(v Value)` |
| 55 | `border-right: 0.1vw solid var(...)` | `BorderRight(v ...Value)` |
| 84 | `padding-bottom: var(--space-1)` | `PaddingBottom(v Value)` |

---

## 2. Additions

All land in `dsl.go`. No new files. One-liner each, following the established pattern.

### 2.1 Property functions (9)

```go
func MarginLeft(v Value) Decl          { return Decl{"margin-left", v.cssValue()} }
func MarginRight(v Value) Decl         { return Decl{"margin-right", v.cssValue()} }
func PaddingBottom(v Value) Decl       { return Decl{"padding-bottom", v.cssValue()} }
func ListStyle(v Value) Decl           { return Decl{"list-style", v.cssValue()} }
func All(v Value) Decl                 { return Decl{"all", v.cssValue()} }
func OverflowY(v Value) Decl           { return Decl{"overflow-y", v.cssValue()} }
func GridTemplateRows(v Value) Decl    { return Decl{"grid-template-rows", v.cssValue()} }
func GridTemplateColumns(v Value) Decl { return Decl{"grid-template-columns", v.cssValue()} }
func BorderRight(v ...Value) Decl      { return Decl{"border-right", joinValues(v)} }
```

### 2.2 Keywords (0 new)

All keywords needed by the above properties are already in the DSL:
- `MarginLeft(Auto)`, `MarginLeft(Zero)` → `Auto`, `Zero` ✅
- `OverflowY(Auto)`, `OverflowY(Visible)` → `Auto`, `Visible` ✅
- `ListStyle(None)` → `None` ✅
- `All(Initial)` → `Initial` ✅

### 2.3 Value constructors (0 new)

- `MarginLeft(Px(5))`, `MarginRight(Rem(0.4))` → `Px`, `Rem` ✅
- `BorderRight(Vw(0.1), Str("solid"), token.Var())` → `Vw`, `Str` ✅
- `GridTemplateRows(Str("auto 1fr"))` → `Str` ✅
- `GridTemplateColumns(Str("1fr 3fr 1fr"))` → `Str` ✅

---

## 3. Consumer migrations (after this plan lands)

Once the new version is published, two packages sweep their `RawRule` TODO comments:

### `layout/platformd/ssr.go`

| Before | After |
| --- | --- |
| `RawRule("list-style: none;")` | `ListStyle(None)` |
| `RawRule("margin-left: 5px;")` | `MarginLeft(Px(5))` |
| `RawRule("margin-left: 0;")` | `MarginLeft(Zero)` |
| `RawRule("margin-left: auto;")` | `MarginLeft(Auto)` |
| `RawRule("margin-left: -100vw;")` | `MarginLeft(Vw(-100))` |
| `RawRule("margin-left: .4rem;")` | `MarginLeft(Rem(0.4))` |
| `RawRule("margin-right: .4rem;")` | `MarginRight(Rem(0.4))` |
| `RawRule("grid-template-columns: 1fr 3fr 1fr;")` | `GridTemplateColumns(Str("1fr 3fr 1fr"))` |
| `RawRule("all: initial;")` | `All(Initial)` |

Remaining permanent `RawRule` (vendor-prefixed — leave as-is):
- `-webkit-box-sizing`, `-moz-box-sizing`
- `-webkit-user-select`, `-khtml-user-select`, `-moz-user-select`, `-ms-user-select`

### `layout/rightpanel/ssr.go`

| Before | After |
| --- | --- |
| `RawRule("overflow: hidden;")` | `Overflow(Hidden)` |
| `RawRule("overflow-y: auto;")` | `OverflowY(Auto)` |
| `RawRule("overflow-y: visible;")` | `OverflowY(Visible)` |
| `RawRule("grid-template-rows: auto 1fr;")` | `GridTemplateRows(Str("auto 1fr"))` |
| `RawRule("border-right: 0.1vw solid "+token.Var()+";")` | `BorderRight(Vw(0.1), Str("solid"), token)` |
| `RawRule("padding-bottom: "+Space1.Var()+";")` | `PaddingBottom(Space1)` |

After the sweep, `grep -n "RawRule" layout/platformd/ssr.go layout/rightpanel/ssr.go` must return only vendor-prefixed lines.

---

## 4. Stages

| # | Stage | Output | Verify |
| --- | --- | --- | --- |
| 0 | **Rename `New` → `NewStylesheet`** | In `tinywasm/css`: rename in `dsl.go`, `ssr.go`, `dsl_test.go`. Grep all `ssr.go` files in the monorepo and update every caller. | `go build ./...` and `go test ./...` across all affected modules green. |
| 1 | **Property functions** | Add 9 functions to `dsl.go`. | `go test ./css/...` green. One assertion per function in `dsl_test.go`. |
| 2 | **Bump version to `v0.1.0`** | Tag and publish. Version must be `v0.1.0` (breaking change). | `go get github.com/tinywasm/css@v0.1.0` resolves in all consumer modules. |
| 3 | **platformd sweep** | Replace all non-vendor `RawRule` in `layout/platformd/ssr.go` per §3 table. Update `go.mod` to `v0.1.0`. | `go build ./platformd/...` and `go test ./platformd/...` green. |
| 4 | **rightpanel sweep** | Replace all `RawRule` in `layout/rightpanel/ssr.go` per §3 table. Update `go.mod` to `v0.1.0`. | `go build ./rightpanel/...` and `go test ./rightpanel/...` green. |
| 5 | **components sweep** | Update all `tinywasm/components/**/ssr.go` callers (`New(` → `NewStylesheet(`). Update each component's `go.mod` to `v0.1.0`. | `go build ./...` green across components. |
| 6 | **Final grep checks** | `grep -rn "\bNew(" --include="ssr.go"` returns zero hits. `grep -n "RawRule" layout/platformd/ssr.go layout/rightpanel/ssr.go` returns only vendor-prefixed lines. | CI green. |

---

## 5. Acceptance criteria

- `New` does not exist in `tinywasm/css` public API. `grep -rn "func New\b" css/` returns zero hits.
- `grep -rn "\bNew(" --include="ssr.go"` across the monorepo returns zero hits.
- 9 new typed property functions added to `dsl.go`, each with a unit test.
- `layout/platformd/ssr.go` — zero non-vendor `RawRule`.
- `layout/rightpanel/ssr.go` — zero `RawRule`.
- `go test ./...` green across `tinywasm/css`, `layout/platformd`, `layout/rightpanel`, and all components.
- `tinywasm/css` published at `v0.1.0`.
- `css/README.md` "Supported properties" table updated; constructor entry reads `NewStylesheet` not `New`.
