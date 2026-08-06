//go:build !wasm

package css

import (
	clr "github.com/tinywasm/color"
	"github.com/tinywasm/font"

	. "github.com/tinywasm/fmt"
)

// FontStack returns a CSS font family stack starting with the given font family.
func FontStack(family font.Family) string {
	return `"` + string(family) + `", system-ui, -apple-system, sans-serif`
}

func RootCSS() *Stylesheet {
	items := append([]item{brandRoot()}, defaultRoots()...)
	return NewStylesheet(items...)
}

// RenderCSS is the base reset, emitted in the `tokens` layer — the lowest of
// the four the widget style DSL declares. Unlayered, `svg { display: block }`
// would beat `@layer widgets { .part { display: none } }` regardless of
// specificity, so any component hiding an icon by state could not.
func RenderCSS() *Stylesheet {
	return NewStylesheet(layer("tokens", resetRules()...))
}

// Override is the customized override of a single Token's value.
type Override struct {
	token Token
	value string

	// light and dark, when non-empty, additionally override the token's
	// LightVarName()/DarkVarName() split properties — see Token.EnhancedVar.
	// Only SetTheme populates these; a plain Set() on a theme token still
	// fills them (forcing one value for both halves, consistent with
	// forcing token.Name to a single value), but Set() on a static token
	// leaves them empty since there is no split to override.
	light, dark string
}

// Set builds an Override for a designated Token with the specified custom value.
func Set(t Token, value string) Override {
	o := Override{token: t, value: value}
	if t.Light != "" {
		o.light, o.dark = value, value
	}
	return o
}

// SetTheme builds an Override for a theme-aware Token with a custom light/dark pair.
func SetTheme(t Token, light, dark string) Override {
	return Override{token: t, value: "light-dark(" + light + ", " + dark + ")", light: light, dark: dark}
}

// Theme returns the entire RootCSS() catalog with custom overrides appended.
func Theme(overrides ...Override) *Stylesheet {
	catalog := RootCSS() // default catalog
	if len(overrides) == 0 {
		return catalog
	}
	var decls []decl
	for _, o := range overrides {
		decls = append(decls, decl{o.token.Name, o.value})
		if o.light != "" {
			decls = append(decls, decl{o.token.LightVarName(), o.light})
			decls = append(decls, decl{o.token.DarkVarName(), o.dark})
		}
	}
	return withRootTail(catalog, root(decls...))
}

func withRootTail(s *Stylesheet, it item) *Stylesheet {
	s.items = append(s.items, it)
	return s
}

// Hover, Focus and Press return the standard interaction-state derivation for
// any base token: the base mixed toward the theme's contrasting extreme.
// The mixer is light-dark(black, white) so a hover darkens on a light theme
// and lightens on a dark one. The intensity is a token, so an app can retune
// it with Theme(Set(MixHover, "22%")) without republishing this package.
func Hover(t Token) string  { return mixToward(t, MixHover) }
func Focus(t Token) string  { return mixToward(t, MixFocus) }
func Press(t Token) string  { return mixToward(t, MixPress) }

func mixToward(t, amount Token) string {
	return "color-mix(in oklab, " + t.NestedEnhanced() + ", light-dark(black, white) " + amount.NestedEnhanced() + ")"
}

// HoverStatic, FocusStatic and PressStatic are the browser-safe counterparts
// of Hover/Focus/Press: t's LightValue mixed toward black by the same
// intensity, computed once in Go instead of once per paint in the browser.
// Callers emit this as the first of a double declaration — see
// tinywasm/widget/style — so a browser without color-mix() support (Safari <
// 16.2) keeps it, permanently in the light theme, instead of an invalid
// declaration.
func HoverStatic(t Token) string { return staticMixToward(t, MixHover) }
func FocusStatic(t Token) string { return staticMixToward(t, MixFocus) }
func PressStatic(t Token) string { return staticMixToward(t, MixPress) }

func staticMixToward(t, amount Token) string {
	return staticMix(t.LightValue(), "#000000", parsePercent(amount.LightValue()))
}

// FadeStatic is the static counterpart of the color-mix(in <space>, TOKEN
// P%, transparent) pattern used for a token faded toward transparency (e.g.
// ColorSelection, or a veil/backdrop wash): t's LightValue faded toward
// transparent by transparentPct (0.0-1.0 — the weight transparent gets).
func FadeStatic(t Token, transparentPct float64) string {
	return staticMix(t.LightValue(), "transparent", transparentPct)
}

// staticMix is the shared computation every *Static derivation reduces to.
func staticMix(aHex, bHex string, pct float64) string {
	return string(clr.Mix(clr.Color(aHex), clr.Color(bHex), pct))
}

// parsePercent parses a token's plain percentage literal (e.g. "15%", the
// only shape MixHover/MixFocus/MixPress ever hold — static Dark-only tokens,
// never light-dark pairs) into a 0.0-1.0 weight. Malformed input — which
// cannot happen for the catalog's own tokens — yields 0 (no shift).
func parsePercent(s string) float64 {
	n := len(s)
	if n < 2 || s[n-1] != '%' {
		return 0
	}
	v, err := Convert(s[:n-1]).Float64()
	if err != nil {
		return 0
	}
	return v / 100.0
}
