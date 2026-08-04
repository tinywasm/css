//go:build !wasm

package css

import "github.com/tinywasm/font"

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
}

// Set builds an Override for a designated Token with the specified custom value.
func Set(t Token, value string) Override { return Override{t, value} }

// SetTheme builds an Override for a theme-aware Token with a custom light/dark pair.
func SetTheme(t Token, light, dark string) Override {
	return Override{t, "light-dark(" + light + ", " + dark + ")"}
}

// Theme returns the entire RootCSS() catalog with custom overrides appended.
func Theme(overrides ...Override) *Stylesheet {
	catalog := RootCSS() // default catalog
	if len(overrides) == 0 {
		return catalog
	}
	decls := make([]decl, len(overrides))
	for i, o := range overrides {
		decls[i] = decl{o.token.Name, o.value}
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
	return "color-mix(in oklab, " + t.Var() + ", light-dark(black, white) " + amount.Var() + ")"
}
