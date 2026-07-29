---
PLAN: "feat: typed Device viewport classes and rail width tokens"
TAG: v0.4.0
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# Plan — `css`: typed `Device` viewport classes + rail width tokens

## 0. Why this exists (read before touching code)

`tinywasm/layout`'s application chassis (`platformd`) renders a blank page today. The
root cause chain ends here: the chassis needs to say *"the nav rail is a fixed column on
desktop and a slide-in drawer on mobile"*, and no library in the ecosystem currently
offers a **typed, closed, tested** way to name a viewport class. The only thing that ever
existed was an unexported helper in this repo (`mediaDesktop`, `dsl.go:162`) hardcoding
`(orientation: landscape) and (hover: hover)`, plus four `--bp-*` custom properties that
**cannot be used in a media query at all** (see §1.2).

This plan closes that gap in `css` — the layer that owns every value in the ecosystem.
`widget/style` and components consume what this plan exports; they must never write a
media query string themselves.

This repo is a **build-time-only** package (`TestPackageIsBuildTimeOnly` in
`tokens_test.go` enforces it). Everything added here is `//go:build !wasm` or lives in an
already-untagged file that never reaches a WASM binary. Do not add a `wasm` build tag to
anything.

---

## 1. Background facts the executor must know

### 1.1 The DSL is closed on purpose

Since v0.3.x every DSL constructor in `dsl.go` is **unexported** (`rule`, `media`, `root`,
`keyframes`, `at`, …). The only exported construction path is `NewStylesheet(items ...item)`
and `Raw(css string) item`. **Do not export any existing lowercase helper.** This plan adds
exactly one new exported type plus one exported function and two tokens — nothing else.

### 1.2 `var()` does not work inside `@media` — this is the whole reason for the type

```css
/* INVALID. Silently never matches. This is not a browser bug. */
@media (min-width: var(--bp-md)) { … }
```

Media-query *conditions* are evaluated before custom properties resolve. That is why
`BpSm`/`BpMd`/`BpLg`/`BpXl` (`catalog.go:68-71`) are declared in `:root` but are useless
for the job. The literal pixel value must be baked into the Go string that produces the
query. Keeping those literals **in this file and nowhere else** is the point of `Device`.

Do **not** delete the `Bp*` tokens: they remain useful to JS and to container queries.
Do **not** try to make `Device.Query()` reference them via `var()`.

### 1.3 Existing token shape

`Token` has `Name` and `Dark` fields; a single-value (non-theme) token puts its value in
`Dark`. See `catalog.go:60-77` for the pattern used by `--z-*` and `--column-*`.

`TestNoUndeclaredTokensInEmittedCSS` (`css_test.go:140`) fails if a token is referenced but
never declared in `RootCSS()`. Every token added below MUST also be added to `RootCSS()`.

---

## 2. Stage 1 — `Device`: the closed set of viewport classes

Create a new file **`device.go`** at the repo root (untagged — it declares no DSL item and
must be readable from anywhere).

```go
package css

// Device is the closed set of viewport classes. It exists because a media query
// condition cannot read a custom property: the pixel thresholds must be baked into
// the query string, and this is the only file in the ecosystem allowed to hold them.
//
// The three classes are mutually exclusive and jointly exhaustive: every viewport
// width matches exactly one. That property is asserted by TestDeviceClassesPartition.
type Device uint8

const (
	Mobile  Device = iota // narrow, touch-first
	Tablet                // medium
	Desktop               // wide, pointer-first
)

func (d Device) String() string {
	switch d {
	case Mobile:
		return "Mobile"
	case Tablet:
		return "Tablet"
	case Desktop:
		return "Desktop"
	default:
		return "Unknown"
	}
}

// Query returns the media-query condition for exactly this class, without the
// leading "@media ". Thresholds mirror BpSm (640px) and BpLg (1024px).
func (d Device) Query() string {
	switch d {
	case Mobile:
		return "(max-width: 639.98px)"
	case Tablet:
		return "(min-width: 640px) and (max-width: 1023.98px)"
	case Desktop:
		return "(min-width: 1024px)"
	default:
		return ""
	}
}

// Query joins several device classes into one condition list. Callers that mean
// "tablet and desktop" pass both rather than emitting two blocks.
// Duplicate and unknown values are dropped; the result is ordered Mobile,
// Tablet, Desktop regardless of argument order, so emission is deterministic.
func Query(devices ...Device) string { … }
```

### 2.1 Exact rules for `Query(devices ...Device)`

- Deduplicate, then sort ascending by the `Device` value, then join each `Query()` with
  `", "` (comma = logical OR in a media query list).
- `Query()` with no arguments returns `""`.
- `Query(Mobile, Tablet, Desktop)` returns the three conditions joined — it does **not**
  collapse to `"all"`. Determinism beats cleverness here and the minifier handles it.
- Use `github.com/tinywasm/fmt` for string work. **No `strings`, no `sort`, no `fmt`
  from the stdlib** anywhere in this repo's non-test files.

### 2.2 The `.98px` detail — do not "clean it up"

`639.98px` / `1023.98px` are deliberate. A browser reporting a fractional viewport width
(zoom, scrollbar, device pixel ratio) at exactly `639.5px` must still match `Mobile`; using
`639px` leaves a half-pixel gap where no class matches and the element falls back to its
base rule — precisely the class of silent layout bug this plan exists to eliminate. Leave
the fractions in place.

---

## 3. Stage 2 — rail width tokens

`widget/style` will gain a `Sidebar` primitive whose fixed column needs a width from a
closed scale. Per this repo's rule *"never invent a value"*, the values live here.

In `catalog.go`, immediately after the `ColumnNarrow/Medium/Wide` block:

```go
	// Rail widths — the fixed column of a Sidebar layout.
	RailNarrow = Token{Name: "--rail-narrow", Dark: "3.5rem"}  // icon only
	RailWide   = Token{Name: "--rail-wide", Dark: "12rem"}     // icon + label
```

In `css.go`, inside `RootCSS()`, extend the existing "Grid columns" `root(...)` group (do
not add a fourth `root(...)` block):

```go
		root(
			// Grid columns
			declare(ColumnNarrow),
			declare(ColumnMedium),
			declare(ColumnWide),
			// Rail widths
			declare(RailNarrow),
			declare(RailWide),
		),
```

`3.5rem` fits a `2.5em` icon plus its padding — the value the previous hand-written
chassis converged on after a documented bug where a `4vw` rail collapsed below the icon
width. Do not change it to a viewport unit.

---

## 4. Stage 3 — tests (`device_test.go`)

`gotest`, never `go test`. Stdlib `testing` + `strings` only; no testify.

| Test | Asserts |
|---|---|
| `TestDeviceClassesPartition` | For widths `320, 639, 640, 767, 1023, 1024, 1440`, exactly **one** of the three `Query()` strings matches. Parse the min/max out of each query and compare numerically — do not string-match. |
| `TestDeviceQueryHasNoVar` | No `Query()` result contains `var(` — §1.2. |
| `TestQueryJoinIsDeterministic` | `Query(Desktop, Mobile)` == `Query(Mobile, Desktop)`, and both equal `Mobile.Query() + ", " + Desktop.Query()`. |
| `TestQueryDeduplicates` | `Query(Mobile, Mobile)` == `Mobile.Query()`. |
| `TestQueryEmpty` | `Query()` == `""`. |
| `TestRailTokensDeclared` | `RootCSS().String()` contains `--rail-narrow` and `--rail-wide`. |
| `TestDeviceStringRoundTrip` | Each of the three constants has a non-`"Unknown"` `String()`. |

Existing `TestNoUndeclaredTokensInEmittedCSS` and `TestComposedTokensAreDeclared` must stay
green — that is the check that catches a forgotten `declare()`.

---

## 5. Documentation (same commit as the code, before publishing)

- **`docs/SPECS.md`** — new section `Device`: the three classes, their exact query strings
  as a table, the partition guarantee, and the `var()`-in-`@media` prohibition from §1.2
  stated as a rule. Add `--rail-narrow` / `--rail-wide` to the token table.
- **`docs/ARCHITECTURE.md`** — one paragraph: why the breakpoint literal lives in Go and
  not in a custom property, cross-referencing the `Bp*` tokens as the JS/container-query
  counterpart.
- **`docs/MIGRATION.md`** — new section "Upgrading to v0.4.0": purely additive, no consumer
  changes required; `Device` is the sanctioned replacement for any hand-written `@media`
  string in `widget/style` and in components.
- **`README.md`** — must index every file in `docs/`. Verify nothing is missing.

Everything in English. Diagrams `flowchart TD`, no `subgraph`, `<br/>` for line breaks.

---

## 6. Explicit non-goals — do NOT do these

- Do **not** export `media`, `rule`, `root`, `keyframes` or any other lowercase DSL helper.
  The chassis is being rebuilt on `widget/style`, not on a reopened `css` DSL.
- Do **not** delete `BpSm`/`BpMd`/`BpLg`/`BpXl`.
- Do **not** delete or export the unexported `mediaDesktop` helper. It stays unexported and
  unused; `widget/style` will use `Device` instead. (A later plan removes it once no
  internal caller remains — not this one.)
- Do **not** add an orientation or `hover:` condition to `Device.Query()`. Width-based
  classes are testable for exhaustivity; capability-based ones are not, and the old
  `(orientation: landscape) and (hover: hover)` query matched most tablets as "desktop",
  which is one of the reasons the chassis broke.
- Do **not** touch the colour tokens, the `light-dark()` machinery, or `Hover/Focus/Press`.

---

## 7. Acceptance criteria

1. `gotest` green in this repo.
2. `grep -rn "@media" --include=*.go .` shows media-query strings **only** in `device.go`
   and the pre-existing unexported `media*` helpers in `dsl.go`.
3. `grep -rn "var(--bp-" --include=*.go .` → empty.
4. `RootCSS().String()` contains `--rail-narrow:3.5rem` and `--rail-wide:12rem`.
5. The exported surface added is exactly: `Device`, `Mobile`, `Tablet`, `Desktop`,
   `(Device).String`, `(Device).Query`, `Query`, `RailNarrow`, `RailWide`. Nothing else.
6. `GOOS=js GOARCH=wasm go list -deps ./...` — this package must remain absent from any
   WASM dependency graph (unchanged behaviour; just confirm the guard test still passes).

---

## 8. Stages

| # | Stage | Files | Gate |
|---|---|---|---|
| 1 | `Device` type, `Query` join | `device.go` (new) | — |
| 2 | Rail width tokens | `catalog.go`, `css.go` | — |
| 3 | Tests | `device_test.go` (new) | blocks 4 |
| 4 | Docs | `docs/SPECS.md`, `docs/ARCHITECTURE.md`, `docs/MIGRATION.md`, `README.md` | — |

Stages 1 and 2 are independent and may be done in either order. Stage 3 must pass before
stage 4 is written.

**Downstream:** `tinywasm/widget` consumes this release
(https://github.com/tinywasm/widget/blob/main/docs/PLAN.md), and `tinywasm/layout` consumes
that one (https://github.com/tinywasm/layout/blob/main/docs/PLAN.md). Neither can start
until this release is published.
