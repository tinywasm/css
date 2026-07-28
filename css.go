//go:build !wasm

package css

func RootCSS() *Stylesheet {
	return NewStylesheet(
		root(
			// Brand group — Pa100T reference palette: steel blue with white text.
			declare(ColorPrimary),
			declare(ColorOnPrimary),
			declare(ColorSecondary),
			declare(ColorOnSecondary),
			declare(ColorSuccess),
			declare(ColorOnSuccess),
			declare(ColorError),
			declare(ColorOnError),
		),
		root(
			// Theme group
			declareSource(ColorBackgroundLight),
			declareSource(ColorBackgroundDark),
			declareSource(ColorSurfaceLight),
			declareSource(ColorSurfaceDark),
			declareSource(ColorSurfaceSunkenLight),
			declareSource(ColorSurfaceSunkenDark),
			declareSource(ColorOnSurfaceLight),
			declareSource(ColorOnSurfaceDark),
			declareSource(ColorOutlineLight),
			declareSource(ColorOutlineDark),
			declareSource(ColorMutedLight),
			declareSource(ColorMutedDark),
			declareSource(ColorHoverLight),
			declareSource(ColorHoverDark),
			declareSource(ColorSelectionLight),
			declareSource(ColorSelectionDark),
			declareSource(ColorOnSelectionLight),
			declareSource(ColorOnSelectionDark),
			declareSource(ColorDisabledLight),
			declareSource(ColorDisabledDark),
			declareSource(ColorOnDisabledLight),
			declareSource(ColorOnDisabledDark),
		),
		root(
			// Active group
			declare(ColorBackground),
			declare(ColorSurface),
			declare(ColorSurfaceSunken),
			declare(ColorOnSurface),
			declare(ColorOutline),
			declare(ColorMuted),
			declare(ColorHover),
			declare(ColorSelection),
			declare(ColorOnSelection),
			declare(ColorDisabled),
			declare(ColorOnDisabled),
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
		root(
			// Focus ring
			declare(ColorFocusRing),
		),
		root(
			// Interactive colors — hover, focus, press states
			declare(ColorPrimaryHover),
			declare(ColorPrimaryFocus),
			declare(ColorPrimaryPress),

			declare(ColorSecondaryHover),
			declare(ColorSecondaryFocus),
			declare(ColorSecondaryPress),

			declare(ColorSuccessHover),
			declare(ColorSuccessFocus),
			declare(ColorSuccessPress),

			declare(ColorDangerHover),
			declare(ColorDangerFocus),
			declare(ColorDangerPress),

			declare(ColorErrorHover),
			declare(ColorErrorFocus),
			declare(ColorErrorPress),

			declare(ColorWarningHover),
			declare(ColorWarningFocus),
			declare(ColorWarningPress),

			declare(ColorInfoHover),
			declare(ColorInfoFocus),
			declare(ColorInfoPress),

			declare(ColorNeutralHover),
			declare(ColorNeutralFocus),
			declare(ColorNeutralPress),

			declare(ColorMutedHover),
			declare(ColorMutedFocus),
			declare(ColorMutedPress),
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
			outline(str("2px solid "+ColorFocusRing.Var())),
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
	token Source
	value string
}

// Set declara el override de un token del catálogo. Token tipado (no un nombre libre);
// value es el borde de I/O.
func Set(t Source, value string) Override { return Override{t, value} }

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
		decls[i] = decl{o.token.Name, o.value}
	}
	return withRootTail(catalog, root(decls...))
}

func withRootTail(s *Stylesheet, it item) *Stylesheet {
	s.items = append(s.items, it)
	return s
}
