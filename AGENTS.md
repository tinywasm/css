# AGENTS.md — tinywasm/css

Constraints for agents changing this library. Read before touching any file.

## Mission

Single source of truth for the ecosystem's design decisions: the token catalog and the
typed CSS emission engine. Stylesheets are produced **at build/SSR time**; the browser
receives plain CSS, never Go.

Consumers: `tinywasm/widget/style` (semantic styling API), `tinywasm/sitec` (SSR asset
registration), `tinywasm/components` and `tinywasm/layout` (their `//go:build !wasm`
`css.go` files).

---

## 1. Build-time only — nothing reaches the WASM binary

Every `.go` file in this package carries `//go:build !wasm`. This is not an optimization,
it is the boundary: with the tag, an import from frontend code fails to compile
(`build constraints exclude all Go files`); without it, the cost is a silent size
regression that depends on the linker's dead-code elimination.

- Values (`#1b5d8c`, `1rem`, `250ms`) describe CSS and are resolved when the stylesheet is
  emitted. There is nothing for the browser-side Go code to do with them.
- Identity strings that the frontend genuinely needs (widget and part names) belong to
  `tinywasm/widget`, which is identity-only and WASM-safe. Do not reintroduce them here.
- Interfaces are the worst thing to leave on the WASM side: satisfying one forces method
  sets and type descriptors to survive DCE. `ValueGetter`, `NamedPair` and `AllPairs()`
  exist solely for the contrast audit — audit code is `!wasm` code.
- If a frontend ever needs a variable name at runtime, pass the string down from
  build-time code. Do not untag the catalog: a `Token` without its values is useless, and
  with its values it is CSS inside the binary.

Guarded by `TestPackageIsBuildTimeOnly`. Under `GOOS=js` this repo's `./...` matches no
packages — a warning and exit 0. That is the expected result, not a broken suite.

## 2. A value is written once

The catalog is the only place a design value may appear.

- **Never repeat a token's value.** A composed token references `var(--other-token)`; it
  never re-states the hex that already lives in that token's literal.
- **Never bake a token's value into a Go string** (`"color-mix(..., " + ColorSurface.Light + ")"`).
  It compiles, and it breaks `SetTheme()`: an app's override only propagates into a
  derived value if that value stays a live `var()` the browser recomputes.
- Reading the value through `t.Var()` is not duplication — it takes it from the source of
  truth. Writing the same literal twice is.

## 3. No formula may freeze the theme direction

Both themes are served by one declaration. A derivation that only reads correctly in one
of them is a defect, not a tuning issue.

- Mix toward a token that already inverts by theme (`var(--color-on-surface)`) or toward
  `light-dark(black, white)`.
- Never mix toward a bare `black`, `white`, or `--color-background`: on a dark theme these
  invert the intent (a hover that darkens an already-dark surface disappears).

## 4. Magnitudes are design decisions too

An intensity, a ratio, a step is as much a design decision as a colour, and gets the same
treatment: a token in the catalog, referenced by `var()`. A percentage frozen inside a Go
function can only be changed by republishing this package and regenerating every
consumer's stylesheet — and can never be overridden under a local scope.

Accepted cost, and it must be documented: an override with a wrongly-typed value
invalidates the whole declaration at computed-value time. A `var()` fallback covers
"undefined", not "ill-defined".

## 5. The public surface is semantic, never raw CSS

`RootCSS`, `RenderCSS`, `Theme`/`Set`/`SetTheme`, the token catalog, the interaction
derivations, `Stylesheet`, `NewStylesheet`, `Raw`. That is the whole surface.

The property and selector helpers of `dsl.go` (`rule`, `root`, `declare`, `padding`,
`color`, …) are **unexported on purpose**: an exported way to write an arbitrary CSS
property is an escape hatch from the design system, which is the thing this library
exists to close. Do not export one, and do not add a second way to do something that
already has one. Rationale: `docs/JUSTIFICACION_DSL.md`.

## 6. Contrast is audited, and the audit has a known blind spot

`AllPairs()` lists the functional surface pairs at a 4.5 minimum. `resolveColor` parses
hex and `light-dark()` — it **cannot evaluate `color-mix()`**.

So: a pair whose background is a composed value cannot be added to `AllPairs()` as if it
were covered. Either justify the gap in writing (naming a proxy pair that *is* audited),
or close it by implementing the mix in Go inside `contrast_test.go`. Silence is not an
option — the accessibility guarantee is something this library advertises.

## 7. Theming is a browser-side mechanism

`light-dark()` in the token fallback plus `color-scheme` on `:root`. No source tokens, no
variable-to-variable `bind()`, no duplicated `prefers-color-scheme` blocks — all of that
was removed and must not come back. `Token.Light`/`Dark` both set = theme pair; `Dark`
alone = static value.

---

## Adding or changing a token

1. Literal in `catalog.go`, in its group. Naming follows the existing prefixes
   (`--color-*`, `--space-*`, `--text-*`, `--radius-*`, `--mix-*`); a foreground is
   `ColorOn<Role>` paired with `Color<Role>`.
2. `declare()` it in the matching `root(...)` group: brand identity colors go in
   `brandRoot()` (`css.brand.go`), every other scale in `defaultRoots()` (`css.default.go`).
   `RootCSS()` (`css.go`) only composes the two.
3. Add it to `allTokens` in `css_test.go` — `TestNoUndeclaredTokensInEmittedCSS` validates
   emitted CSS against that hardcoded list and fails otherwise.
4. If it is a surface pair with literal values, add it to `AllPairs()`; if its value is
   composed, apply §6.
5. `docs/SPECS.md` always; `docs/MIGRATION.md` when it replaces or removes something.
6. `gotest`.

Removing a token is a breaking change for `widget` and `components`. The replacement
mapping ships in `docs/MIGRATION.md` in the same release — never leave a consumer to
guess, and never tell a consumer to compute the value itself: a formula duplicated across
consumers is the drift this catalog eliminates.

## Testing

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest   # external agents have no global gotest
gotest
```

`gotest`, never `go test`. Stdlib assertions only (`testing`/`strings`, no testify).
Guards are part of the design, not incidental coverage: if a rule above can be checked by
a test, it has one, and weakening a guard requires the same justification as changing the
rule.

## Publishing

`gopush 'message'`, never a raw `git commit`/`push`. Documentation is updated **before**
publishing, in the same commit as the code it describes.

## Documentation

`docs/ARCHITECTURE.md` (what and why), `docs/SPECS.md` (exact values and API — the
authority; if code and SPECS disagree, one of them is wrong and both are fixed in the same
commit), `docs/MIGRATION.md` (consumer upgrade path), `docs/JUSTIFICACION_DSL.md` (why the
DSL exists in this shape). `README.md` indexes every file in `docs/`.

`docs/PLAN.md` is ephemeral: it is deleted in the commit that publishes the work it
describes, so no permanent document may link to it. Everything written in English.
