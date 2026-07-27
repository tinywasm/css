//go:build !wasm

package css

func RootCSS() *Stylesheet {
	return NewStylesheet(
		root(
			// Brand group — Pa100T reference palette: steel blue with white text.
			declare(ColorPrimary, "#3f88bf"),
			declare(ColorOnPrimary, "#FFFFFF"),
			declare(ColorSecondary, "#654FF0"),
			declare(ColorOnSecondary, "#FFFFFF"),
			declare(ColorSuccess, "#3FB950"),
			declare(ColorOnSuccess, "#FFFFFF"),
			declare(ColorError, "#E34F26"),
			declare(ColorOnError, "#FFFFFF"),
		),
		root(
			// Theme group
			declare(ColorBackgroundLight, "#FFFFFF"),
			declare(ColorBackgroundDark, "#0D1117"),
			declare(ColorSurfaceLight, "#F2F2F7"),
			declare(ColorSurfaceDark, "#161B22"),
			declare(ColorSurfaceSunkenLight, "#E5E5EA"),
			declare(ColorSurfaceSunkenDark, "#21262D"),
			declare(ColorOnSurfaceLight, "#1C1C1E"),
			declare(ColorOnSurfaceDark, "#E6EDF3"),
			declare(ColorOutlineLight, "#D1D1D6"),
			declare(ColorOutlineDark, "#30363D"),
			declare(ColorMutedLight, "#6E6E73"),
			declare(ColorMutedDark, "#8B949E"),
			declare(ColorHoverLight, "#B8860B"),
			declare(ColorHoverDark, "#F7DF1E"),
			declare(ColorSelectionLight, "#f5a623"),
			declare(ColorSelectionDark, "#9e6a2e"),
			declare(ColorOnSelectionLight, "#1C1C1E"),
			declare(ColorOnSelectionDark, "#FFFFFF"),
			declare(ColorDisabledLight, "#E5E5EA"),
			declare(ColorDisabledDark, "#21262D"),
			declare(ColorOnDisabledLight, "#8E8E93"),
			declare(ColorOnDisabledDark, "#6E7681"),
		),
		root(
			// Typography scale
			declare(TextXs, "0.75rem"),
			declare(TextSm, "0.875rem"),
			declare(TextBase, "1rem"),
			declare(TextLg, "1.25rem"),
			declare(TextXl, "1.5rem"),
			declare(Text2xl, "2rem"),
		),
		root(
			// Spacing scale
			declare(Space0, "0"),
			declare(Space1, "0.25rem"),
			declare(Space2, "0.5rem"),
			declare(Space3, "0.75rem"),
			declare(Space4, "1rem"),
			declare(Space6, "1.5rem"),
			declare(Space8, "2rem"),
			declare(Space12, "3rem"),
		),
		root(
			// Border-radius scale
			declare(RadiusSm, "4px"),
			declare(RadiusMd, "8px"),
			declare(RadiusLg, "16px"),
			declare(RadiusFull, "9999px"),
		),
		root(
			// Typography
			declare(LeadingTight, "1.25"),
			declare(LeadingNormal, "1.5"),
			declare(LeadingRelaxed, "1.75"),
			declare(FontWeightRegular, "400"),
			declare(FontWeightMedium, "500"),
			declare(FontWeightBold, "700"),
			declare(TrackingTight, "-0.02em"),
			declare(TrackingNormal, "0"),
			declare(TrackingWide, "0.05em"),
		),
		root(
			// Elevation
			declare(ShadowSm, "0 1px 2px rgba(0, 0, 0, 0.05)"),
			declare(ShadowMd, "0 4px 6px rgba(0, 0, 0, 0.10)"),
			declare(ShadowLg, "0 10px 15px rgba(0, 0, 0, 0.10)"),
			declare(ShadowXl, "0 20px 25px rgba(0, 0, 0, 0.15)"),
		),
		root(
			// Motion
			declare(DurationFast, "150ms"),
			declare(DurationBase, "250ms"),
			declare(DurationSlow, "400ms"),
			declare(EaseIn, "cubic-bezier(0.4, 0,   1,   1)"),
			declare(EaseOut, "cubic-bezier(0,   0,   0.2, 1)"),
			declare(EaseInOut, "cubic-bezier(0.4, 0,   0.2, 1)"),
		),
		root(
			// Z-index
			declare(ZBase, "0"),
			declare(ZDropdown, "100"),
			declare(ZSticky, "200"),
			declare(ZModal, "300"),
			declare(ZToast, "400"),
			declare(ZTooltip, "500"),
		),
		root(
			// Breakpoints
			declare(BpSm, "640px"),
			declare(BpMd, "768px"),
			declare(BpLg, "1024px"),
			declare(BpXl, "1280px"),
		),
		root(
			// Container widths
			declare(MaxWProse, "65ch"),
			declare(MaxWContent, "1200px"),
			declare(MaxWScreen, "1440px"),
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
		root(
			bind(ColorBackground, ColorBackgroundLight),
			bind(ColorSurface, ColorSurfaceLight),
			bind(ColorSurfaceSunken, ColorSurfaceSunkenLight),
			bind(ColorOnSurface, ColorOnSurfaceLight),
			bind(ColorOutline, ColorOutlineLight),
			bind(ColorMuted, ColorMutedLight),
			bind(ColorHover, ColorHoverLight),
			bind(ColorSelection, ColorSelectionLight),
			bind(ColorOnSelection, ColorOnSelectionLight),
			bind(ColorDisabled, ColorDisabledLight),
			bind(ColorOnDisabled, ColorOnDisabledLight),
		),
		mediaPrefersDark(
			root(
				bind(ColorBackground, ColorBackgroundDark),
				bind(ColorSurface, ColorSurfaceDark),
				bind(ColorSurfaceSunken, ColorSurfaceSunkenDark),
				bind(ColorOnSurface, ColorOnSurfaceDark),
				bind(ColorOutline, ColorOutlineDark),
				bind(ColorMuted, ColorMutedDark),
				bind(ColorHover, ColorHoverDark),
				bind(ColorSelection, ColorSelectionDark),
				bind(ColorOnSelection, ColorOnSelectionDark),
				bind(ColorDisabled, ColorDisabledDark),
				bind(ColorOnDisabled, ColorOnDisabledDark),
			),
		),
		// Explicit theme override via the [data-theme] attribute (set by a theme
		// toggle on <html>). These come after the :root/@media defaults and have
		// higher specificity, so a manual choice wins over the OS preference.
		// Removing the attribute (ThemeAuto) falls back to the @media query above.
		rule(selector("[data-theme=\"light\"]"),
			bind(ColorBackground, ColorBackgroundLight),
			bind(ColorSurface, ColorSurfaceLight),
			bind(ColorSurfaceSunken, ColorSurfaceSunkenLight),
			bind(ColorOnSurface, ColorOnSurfaceLight),
			bind(ColorOutline, ColorOutlineLight),
			bind(ColorMuted, ColorMutedLight),
			bind(ColorHover, ColorHoverLight),
			bind(ColorSelection, ColorSelectionLight),
			bind(ColorOnSelection, ColorOnSelectionLight),
			bind(ColorDisabled, ColorDisabledLight),
			bind(ColorOnDisabled, ColorOnDisabledLight),
		),
		rule(selector("[data-theme=\"dark\"]"),
			bind(ColorBackground, ColorBackgroundDark),
			bind(ColorSurface, ColorSurfaceDark),
			bind(ColorSurfaceSunken, ColorSurfaceSunkenDark),
			bind(ColorOnSurface, ColorOnSurfaceDark),
			bind(ColorOutline, ColorOutlineDark),
			bind(ColorMuted, ColorMutedDark),
			bind(ColorHover, ColorHoverDark),
			bind(ColorSelection, ColorSelectionDark),
			bind(ColorOnSelection, ColorOnSelectionDark),
			bind(ColorDisabled, ColorDisabledDark),
			bind(ColorOnDisabled, ColorOnDisabledDark),
		),
	)
}

// Override es el cambio de valor de UN token. Campos no exportados: solo Set lo construye.
type Override struct {
	token Token
	value string
}

// Set declara el override de un token del catálogo. Token tipado (no un nombre libre);
// value es el borde de I/O.
func Set(t Token, value string) Override { return Override{t, value} }

// Theme devuelve el catálogo :root COMPLETO (como RootCSS) con los overrides al final.
// Pensado como el RootCSS() del proyecto raíz — assetmin REEMPLAZA el :root de css por el de
// la app, por eso trae el catálogo entero, no solo los overrides.
func Theme(overrides ...Override) *Stylesheet {
	catalog := RootCSS() // catálogo por defecto
	if len(overrides) == 0 {
		return catalog
	}
	decls := make([]decl, len(overrides))
	for i, o := range overrides {
		decls[i] = declare(o.token, o.value)
	}
	return withRootTail(catalog, root(decls...))
}

func withRootTail(s *Stylesheet, it item) *Stylesheet {
	s.items = append(s.items, it)
	return s
}
