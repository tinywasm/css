//go:build !wasm

package css

func RootCSS() *Stylesheet {
	return NewStylesheet(
		Root(
			// Brand group
			Declare(ColorPrimary, "#00ADD8"),
			Declare(ColorOnPrimary, "#1C1C1E"),
			Declare(ColorSecondary, "#654FF0"),
			Declare(ColorOnSecondary, "#FFFFFF"),
			Declare(ColorSuccess, "#3FB950"),
			Declare(ColorError, "#E34F26"),
		),
		Root(
			// Theme group
			Declare(ColorBackgroundLight, "#FFFFFF"),
			Declare(ColorBackgroundDark, "#0D1117"),
			Declare(ColorSurfaceLight, "#F2F2F7"),
			Declare(ColorSurfaceDark, "#161B22"),
			Declare(ColorOnSurfaceLight, "#1C1C1E"),
			Declare(ColorOnSurfaceDark, "#E6EDF3"),
			Declare(ColorMutedLight, "#6E6E73"),
			Declare(ColorMutedDark, "#8B949E"),
			Declare(ColorHoverLight, "#B8860B"),
			Declare(ColorHoverDark, "#F7DF1E"),
		),
		Root(
			// Typography scale
			Declare(TextXs, "0.75rem"),
			Declare(TextSm, "0.875rem"),
			Declare(TextBase, "1rem"),
			Declare(TextLg, "1.25rem"),
			Declare(TextXl, "1.5rem"),
			Declare(Text2xl, "2rem"),
		),
		Root(
			// Spacing scale
			Declare(Space1, "0.25rem"),
			Declare(Space2, "0.5rem"),
			Declare(Space3, "0.75rem"),
			Declare(Space4, "1rem"),
			Declare(Space6, "1.5rem"),
			Declare(Space8, "2rem"),
			Declare(Space12, "3rem"),
		),
		Root(
			// Border-radius scale
			Declare(RadiusSm, "4px"),
			Declare(RadiusMd, "8px"),
			Declare(RadiusLg, "16px"),
			Declare(RadiusFull, "9999px"),
		),
		Root(
			// Typography
			Declare(LeadingTight, "1.25"),
			Declare(LeadingNormal, "1.5"),
			Declare(LeadingRelaxed, "1.75"),
			Declare(FontWeightRegular, "400"),
			Declare(FontWeightMedium, "500"),
			Declare(FontWeightBold, "700"),
			Declare(TrackingTight, "-0.02em"),
			Declare(TrackingNormal, "0"),
			Declare(TrackingWide, "0.05em"),
		),
		Root(
			// Elevation
			Declare(ShadowSm, "0 1px 2px rgba(0, 0, 0, 0.05)"),
			Declare(ShadowMd, "0 4px 6px rgba(0, 0, 0, 0.10)"),
			Declare(ShadowLg, "0 10px 15px rgba(0, 0, 0, 0.10)"),
			Declare(ShadowXl, "0 20px 25px rgba(0, 0, 0, 0.15)"),
		),
		Root(
			// Motion
			Declare(DurationFast, "150ms"),
			Declare(DurationBase, "250ms"),
			Declare(DurationSlow, "400ms"),
			Declare(EaseIn, "cubic-bezier(0.4, 0,   1,   1)"),
			Declare(EaseOut, "cubic-bezier(0,   0,   0.2, 1)"),
			Declare(EaseInOut, "cubic-bezier(0.4, 0,   0.2, 1)"),
		),
		Root(
			// Z-index
			Declare(ZBase, "0"),
			Declare(ZDropdown, "100"),
			Declare(ZSticky, "200"),
			Declare(ZModal, "300"),
			Declare(ZToast, "400"),
			Declare(ZTooltip, "500"),
		),
		Root(
			// Breakpoints
			Declare(BpSm, "640px"),
			Declare(BpMd, "768px"),
			Declare(BpLg, "1024px"),
			Declare(BpXl, "1280px"),
		),
		Root(
			// Container widths
			Declare(MaxWProse, "65ch"),
			Declare(MaxWContent, "1200px"),
			Declare(MaxWScreen, "1440px"),
		),
	)
}

func RenderCSS() *Stylesheet {
	return NewStylesheet(
		Rule(Selector("*, *::before, *::after"),
			BoxSizing(Str("border-box")),
		),
		Rule(Selector("html"),
			RawRule("  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;"),
		),
		Rule(Selector("body"),
			Margin(Zero),
			FontFamily(Str("system-ui, -apple-system, \"Segoe UI\", Roboto, sans-serif")),
			FontSize(TextBase),
			LineHeight(LeadingNormal),
			Color(ColorOnSurface),
			Background(ColorBackground),
		),
		Rule(Selector(":focus-visible"),
			Outline(Str("2px solid " + ColorPrimary.Var())),
			OutlineOffset(Px(2)),
		),
		Rule(Selector("img, svg, video"),
			Display(Block),
			MaxWidth(Pct(100)),
		),
		Root(
			Bind(ColorBackground, ColorBackgroundLight),
			Bind(ColorSurface, ColorSurfaceLight),
			Bind(ColorOnSurface, ColorOnSurfaceLight),
			Bind(ColorMuted, ColorMutedLight),
			Bind(ColorHover, ColorHoverLight),
		),
		MediaPrefersDark(
			Root(
				Bind(ColorBackground, ColorBackgroundDark),
				Bind(ColorSurface, ColorSurfaceDark),
				Bind(ColorOnSurface, ColorOnSurfaceDark),
				Bind(ColorMuted, ColorMutedDark),
				Bind(ColorHover, ColorHoverDark),
			),
		),
	)
}
