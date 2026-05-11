//go:build !wasm

package css

import _ "embed"

//go:embed theme.css
var rootCSS string

//go:embed render.css
var renderCSS string

// RootCSS returns the canonical design tokens (vocabulary) as a CSS string.
// assetmin uses it as the framework default for the open slot.
//
// Resolution is single-winner: if the app root declares its own RootCSS(),
// that block fully replaces this one. Apps that want to inherit framework
// defaults compose explicitly with `css.RootCSS() + override`.
//
// Contains ONLY :root {} value declarations: brand colors, source-layer
// theme tokens (-light/-dark), typography scale, spacing scale, radius scale.
func RootCSS() string { return rootCSS }

// RenderCSS returns the framework switching rules that bind active tokens
// (`--color-X`) to source tokens (`--color-X-light` / `--color-X-dark`),
// including `@media (prefers-color-scheme: dark)` for automatic OS-driven
// dark mode.
//
// assetmin injects RenderCSS additively — it is never replaced by an app.
// This guarantees dark-mode wiring survives any app-level RootCSS override,
// provided the app supplies the source-layer tokens.
func RenderCSS() string { return renderCSS }
