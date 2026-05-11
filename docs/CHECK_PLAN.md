# PLAN — `@keyframes` support for tinywasm/css

> The previous PLAN (typed CSS DSL + token catalog + `RootCSS`/`RenderCSS`) has been executed and published. This file replaces it with the remaining typed-CSS work blocking the component migration.

## Goal

Add typed `@keyframes` support to the DSL so animation frames can reference tokens through the same compile-time-safe API as regular rules. Eliminate the need for `Raw("@keyframes ...")` as the canonical way to declare animations.

## Why

The audit of `tinywasm/components` revealed exactly one `@keyframes` in active use (`button.pulse-url`) and it references `--color-secondary` twice. Migrating that animation via `Raw()` would embed those token references as hardcoded strings inside an opaque payload, defeating the DSL's main guarantee: that renaming a token breaks the build everywhere it is used.

Adding `Keyframes()` is the only way to keep the DSL's invariant intact: **no token reference escapes the type system.**

Secondary motivation: `@keyframes` is structurally identical to `@media` (at-rule with block body). The DSL already tipa `Media()` and `MediaPrefersDark()`; not having `Keyframes()` is a visible asymmetry.

## API addition

```go
//go:build !wasm
package css

// KeyframeStep is one step of a keyframes animation.
// At is the percentage or named position ("0%", "50%", "100%", "from", "to").
type KeyframeStep struct {
    At    string
    Decls []Decl
}

// At builds a KeyframeStep. Variadic Decls match the DSL's existing rule shape.
func At(at string, decls ...Decl) KeyframeStep {
    return KeyframeStep{At: at, Decls: decls}
}

// Keyframes builds an @keyframes at-rule.
func Keyframes(name string, steps ...KeyframeStep) item { ... }
```

### Usage

```go
Keyframes("pulse-url",
    At("0%",
        Transform(Str("scale(0.90)")),
        BoxShadow(Str("0 0 0 0 "+ColorSecondary.Var())),
    ),
    At("70%",
        Transform(Str("scale(1)")),
        BoxShadow(Str("0 0 0 1vw transparent")),
    ),
    At("100%",
        Background(ColorSecondary),
        Transform(Str("scale(0.90)")),
    ),
)
```

Renaming `ColorSecondary` in `tokens.go` now breaks this code at compile time — the property the DSL exists to provide.

## Emitted CSS

```css
@keyframes pulse-url {
  0% { transform: scale(0.90); box-shadow: 0 0 0 0 var(--color-secondary,#654FF0); }
  70% { transform: scale(1); box-shadow: 0 0 0 1vw transparent; }
  100% { background: var(--color-secondary,#654FF0); transform: scale(0.90); }
}
```

Formatting (whitespace, newlines) matches the existing DSL output style for consistency with `Rule()` and `Media()`.

## Files modified

- `dsl.go` — add `KeyframeStep` type, `At()` constructor, `Keyframes()` constructor + internal `writeTo()` implementation.
- `dsl_test.go` — add tests covering:
  - Single step (`At("from", ...)`)
  - Multi-step percentage form
  - Token references inside steps render via `var(...)`
  - Empty keyframes (no steps) renders an empty `@keyframes name {}` (defensive — or panic, decide in implementation).

## Files added

None.

## Steps

1. Implement `KeyframeStep`, `At()`, `Keyframes()` in `dsl.go`. Keep `KeyframeStep` exported (consumers need to build the slice); `at` and `Decls` are exported fields so the type is plain-old-data.
2. Implement `writeTo()` for keyframes — mirrors the structure of `Media`/`Rule` emission, but each step emits as `<at> { <decls> }`.
3. Add tests in `dsl_test.go`. Include one golden-style test with a token reference to lock the rendering.
4. Update `README.md` DSL Reference section: add `Keyframes(name, At(...)...)` alongside `Media`, `Rule`, etc.

## Acceptance

- `Keyframes()`, `At()`, `KeyframeStep` are exported from `tinywasm/css`.
- A keyframe step referencing a `Token` emits `var(--name,fallback)` in the output.
- `go test ./...` passes including new keyframe tests.
- `README.md` mentions `Keyframes` in the DSL Reference section.

## Non-goals

- Animation shorthand helpers (`Animation()` accepting a `Keyframes` reference by name) — separate enhancement; `Animation(Str("pulse-url 2s infinite"))` is already expressible.
- Cubic-bezier or motion-token integration — orthogonal.
- `@property` at-rule, `@scope`, `@layer` — not required by any audited consumer.

## Dependency

This plan must land **before** `tinywasm/components/docs/PLAN.md` migrates `button`, which is the only component currently blocked on this API.
