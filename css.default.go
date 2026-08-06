//go:build !wasm

package css

// themeGroupDecls declares the adaptive theme colors plus, for every theme
// pair among them, the plain light/dark halves declareSplit needs — see
// Token.EnhancedVar for why widget/style reads those instead of the
// light-dark()-valued properties declared here.
func themeGroupDecls() []decl {
	decls := []decl{
		declare(ColorBackground),
		declare(ColorOnBackground),
		declare(ColorSurface),
		declare(ColorOnSurface),
		declare(ColorOutline),
		declare(ColorMuted),
		declare(ColorSurfaceSunken),
		declare(ColorSelection),
		declare(ColorOnSelection),
	}
	for _, t := range []Token{ColorBackground, ColorOnBackground, ColorSurface, ColorOnSurface, ColorOutline, ColorMuted} {
		decls = append(decls, declareSplit(t)...)
	}
	return decls
}

// defaultRoots declares every token group that is not brand identity: the
// adaptive theme surface colors plus the typography, spacing, radius,
// shadow, motion, z-index, breakpoint, and layout scales. These are the
// tokens an app inherits as-is; only brandRoot() is meant to be reskinned.
func defaultRoots() []item {
	return []item{
		root(themeGroupDecls()...),
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
			declare(FontSans),
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
		root(
			// Device geometry — insets reported by the device, and the viewport
			// height that shrinks with Safari iOS' collapsing URL bar.
			declare(SafeTop),
			declare(SafeRight),
			declare(SafeBottom),
			declare(SafeLeft),
			declare(ViewportH),
		),
	}
}
