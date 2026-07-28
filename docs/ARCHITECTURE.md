# Theming Architecture

## RootCSS and the Single-Winner Slot

In the `tinywasm` SSR pipeline, `assetmin` discovers the `RootCSS()` blocks of the dependency graph. The `:root` block (the token vocabulary) is a **single-winner slot by replacement**:

1. If the application (the root project) declares `func RootCSS() *css.Stylesheet`, that block **completely replaces** the default `RootCSS()` of the `tinywasm/css` library.
2. `RenderCSS()` (the rule and binding logic) is **additive**: the contributions of all modules are concatenated.

Because `RootCSS()` is completely self-contained and declares all active/bound tokens too, it is complete on its own.

## Entrypoint: `Theme()`

To facilitate rebranding without losing the complete token catalog (spacing scales, typography, etc.), the library provides `css.Theme(...Override)`.

`Theme()` obtains the default catalog via `RootCSS()` and appends a final `:root` block with the provided overrides. This ensures that:

- The application does not need to redeclare tokens it does not want to change.
- The application's changes win the CSS cascade by appearing at the end of the block.
- The catalog remains complete so that components continue to function.

## Type-Safety with `Override` and `Source`

To prevent design-system evasion, we enforce a strict separation between referenceable active tokens (`Token`) and settable-only source tokens (`Source`):

- **`Token`** is referenceable in component rules (has `.Var()`), but cannot be passed to `Set()`.
- **`Source`** represents the theme source variables (like `ColorSurfaceLight`, `ColorBackgroundDark`) which are settable using `css.Set(Source, value)` but cannot be directly referenced in rules (has no `.Var()`). This enforces that all references map to the dynamic light/dark active layers.

The `Override` type is opaque and can only be constructed using `css.Set(Source, value)`.

```go
// Correct usage (overriding a theme source)
css.Theme(css.Set(css.ColorSurfaceLight, "#FAFAFA"))

// Does not compile (Set requires a Source token, ColorPrimary is a Token)
// css.Theme(css.Set(css.ColorPrimary, "#FF0000"))
```

## Declaration with Single-Source Values

`declare` helper in the DSL reads the `Fallback` of a token directly, guaranteeing that all token default values reside exactly once in `tokens.go` rather than being duplicated in `css.go`:

```go
// Declares the default fallback value defined exactly once in tokens.go
declare(ColorPrimary)
```
