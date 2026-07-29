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
		),
		root(
			// Theme group — Adaptive layout colors using CSS light-dark() fallbacks.
			declare(ColorBackground),
			declare(ColorOnBackground),
			declare(ColorSurface),
			declare(ColorOnSurface),
			declare(ColorOutline),
			declare(ColorMuted),
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
		),
	)
}

func RenderCSS() *Stylesheet {
	return NewStylesheet(
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
	)
}

// Override is the customized override of a single Token's value.
type Override struct {
	token Token
	value string
}

// Set builds an Override for a designated Token with the specified custom value.
func Set(t Token, value string) Override { return Override{t, value} }

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
