# Theming Architecture

## RootCSS and the Single-Winner Slot

In the `tinywasm` SSR pipeline, `assetmin` discovers the `RootCSS()` blocks of the dependency graph. The `:root` block (the token vocabulary) is a **single-winner slot by replacement**:

1. If the application (the root project) declares `func RootCSS() *css.Stylesheet`, that block **completely replaces** the default `RootCSS()` of the `tinywasm/css` library.
2. `RenderCSS()` (the rule and binding logic) is **additive**: the contributions of all modules are concatenated.

## Entrypoint: `Theme()`

To facilitate rebranding without losing the complete token catalog (spacing scales, typography, etc.), the library provides `css.Theme(...Override)`.

`Theme()` obtains the default catalog via `RootCSS()` and appends a final `:root` block with the provided overrides. This ensures that:

- The application does not need to redeclare tokens it does not want to change.
- The application's changes win the CSS cascade by appearing at the end of the block.
- The catalog remains complete so that components continue to function.

## Type-Safety with `Override`

The `Override` type is opaque and can only be constructed using `css.Set(Token, value)`. This prevents illegal states such as trying to inject arbitrary CSS properties into the `:root` block through the theme entrypoint, forcing only typed catalog tokens to be overridden.

```go
// Correct usage
css.Theme(css.Set(css.ColorPrimary, "#hex"))

// Does not compile (Set requires a Token)
// css.Theme(css.Set("padding", "20px"))
```
