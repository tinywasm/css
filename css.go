//go:build !wasm

package css

func RootCSS() *Stylesheet {
	return NewStylesheet(
		root(
			// Brand group — Pa100T reference palette: steel blue with white text.
			declare(ColorPrimary),
			declare(ColorOnPrimary),
			declare(ColorSuccess),
			declare(ColorOnSuccess),
			declare(ColorDanger),
			declare(ColorOnDanger),
			declare(ColorAccent),
			declare(ColorOnAccent),
		),
		root(
			// Theme group — Adaptive layout colors using CSS light-dark() fallbacks.
			declare(ColorBackground),
			declare(ColorOnBackground),
			declare(ColorSurface),
			declare(ColorOnSurface),
			declare(ColorOutline),
			declare(ColorMuted),
			declare(ColorSurfaceSunken),
			declare(ColorSelection),
			declare(ColorOnSelection),
		),
		root(
			// Interaction intensities
			declare(MixHover),
			declare(MixFocus),
			declare(MixPress),
		),
		root(
			// Typography scale
			declare(TextXs),
			declare(TextSm),
			declare(TextBase),
			declare(TextLg),
			declare(TextXl),
			declare(Text2xl),
		),
		root(
			// Spacing scale
			declare(Space0),
			declare(Space1),
			declare(Space2),
			declare(Space3),
			declare(Space4),
			declare(Space6),
			declare(Space8),
			declare(Space12),
		),
		root(
			// Border-radius scale
			declare(RadiusSm),
			declare(RadiusMd),
			declare(RadiusLg),
			declare(RadiusFull),
		),
		root(
			// Typography
			declare(LeadingNormal),
			declare(FontWeightRegular),
			declare(FontWeightMedium),
			declare(FontWeightBold),
		),
		root(
			// Elevation
			declare(ShadowSm),
			declare(ShadowMd),
			declare(ShadowLg),
		),
		root(
			// Motion
			declare(DurationFast),
			declare(DurationBase),
			declare(DurationSlow),
			declare(EaseInOut),
		),
		root(
			// Z-index
			declare(ZBase),
			declare(ZDropdown),
			declare(ZSticky),
			declare(ZModal),
			declare(ZToast),
			declare(ZTooltip),
		),
		root(
			// Breakpoints
			declare(BpSm),
			declare(BpMd),
			declare(BpLg),
			declare(BpXl),
		),
		root(
			// Container widths
			declare(MaxWReadable),
		),
		root(
			// Grid columns
			declare(ColumnNarrow),
			declare(ColumnMedium),
			declare(ColumnWide),
			// Rail widths
			declare(RailNarrow),
			declare(RailWide),
			// Shared control height and chip width
			declare(ControlHeight),
			declare(ChipWidth),
			declare(VeilBlur),
		),
	)
}

// RenderCSS is the base reset. Every rule sits in the `tokens` layer — the
// lowest of the four the widget style DSL declares. Unlayered, `svg { display:
// block }` would beat `@layer widgets { .part { display: none } }` regardless
// of specificity, so any component hiding an icon by state could not.
func RenderCSS() *Stylesheet {
	return NewStylesheet(layer("tokens",
		rule(selector("*, *::before, *::after"),
			boxSizing(str("border-box")),
		),
		rule(selector("html"),
			rawRule("  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;"),
		),
		rule(selector("body"),
			margin(zero),
			fontFamily(str("system-ui, -apple-system, \"Segoe UI\", Roboto, sans-serif")),
			fontSize(TextBase),
			lineHeight(LeadingNormal),
			color(ColorOnSurface),
			background(ColorBackground),
		),
		rule(selector(":focus-visible"),
			outline(str("2px solid "+ColorPrimary.Var())),
			outlineOffset(px(2)),
		),
		// User-agent margins and list chrome are geometry the style DSL cannot
		// express and cannot see. A <ul> used as a nav rail carries 40px of
		// padding-inline-start it never asked for, an <h1> carries 0.67em of
		// block margin: both silently widen or heighten the component around
		// them. Zero them here, in the lowest layer, so a part's box is exactly
		// what its rule declares.
		rule(selector("h1, h2, h3, h4, h5, h6, p, figure, blockquote, dl, dd"),
			margin(zero),
		),
		// Heading sizes are the DSL's business, not the user agent's. An <h1>
		// carries font-size: 2em by default, which compounds against whatever
		// the part around it declares and makes a title three times the size a
		// component asked for. FontSize()/FontWeight() say what a heading looks
		// like; its level stays a semantic choice.
		// An <a> arrives underlined and painted the user agent's link blue, which
		// fights every surface a component puts it on — a nav item ends up blue
		// text on a blue button. Components say what a link looks like; the
		// browser's default is chrome, like the heading sizes below.
		rule(selector("a"),
			rawRule("  color: inherit;\n  text-decoration: none;"),
		),
		rule(selector("h1, h2, h3, h4, h5, h6"),
			rawRule("  font-size: inherit;\n  font-weight: inherit;"),
		),
		rule(selector("ol, ul"),
			rawRule("  list-style: none;"),
			margin(zero),
			padding(zero),
		),
		// A <summary> draws a disclosure triangle the component did not ask
		// for and cannot style. Components supply their own affordance.
		rule(selector("summary"),
			rawRule("  list-style: none;"),
		),
		rule(selector("summary::-webkit-details-marker"),
			display(none),
		),
		rule(selector("img, svg, video"),
			display(block),
			maxWidth(pct(100)),
		),
		rule(selector(":root"),
			rawRule("color-scheme: light dark;"),
		),
		rule(selector("[data-theme=\"light\"]"),
			rawRule("color-scheme: light;"),
		),
		rule(selector("[data-theme=\"dark\"]"),
			rawRule("color-scheme: dark;"),
		),
	))
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
