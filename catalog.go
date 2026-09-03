//go:build !wasm

package css

var (
	ColorPrimary     = Token{Name: "--color-primary", Dark: "#654FF0"}
	ColorOnPrimary   = Token{Name: "--color-on-primary", Dark: "#FFFFFF"}
	ColorSuccess     = Token{Name: "--color-success", Dark: "#1e7a30"}
	ColorOnSuccess   = Token{Name: "--color-on-success", Dark: "#FFFFFF"}
	ColorDanger      = Token{Name: "--color-danger", Dark: "#ba2c0d"}
	ColorOnDanger    = Token{Name: "--color-on-danger", Dark: "#FFFFFF"}
	ColorAccent      = Token{Name: "--color-accent", Dark: "#e8a33d"}
	ColorOnAccent    = Token{Name: "--color-on-accent", Dark: "#1C1C1E"}

	ColorBackground   = Token{Name: "--color-background", Light: "#FFFFFF", Dark: "#0D1117"}
	ColorOnBackground = Token{Name: "--color-on-background", Light: "#1C1C1E", Dark: "#E6EDF3"}
	ColorSurface      = Token{Name: "--color-surface", Light: "#F2F2F7", Dark: "#161B22"}
	ColorOnSurface    = Token{Name: "--color-on-surface", Light: "#1C1C1E", Dark: "#E6EDF3"}
	ColorOutline      = Token{Name: "--color-outline", Light: "#D1D1D6", Dark: "#30363D"}
	ColorMuted        = Token{Name: "--color-muted", Light: "#6E6E73", Dark: "#8B949E"}

	// Computed semantic tokens — derived from live var() references, no hex drift.

	// LightStatic: these three hold live var()/color-mix() expressions a
	// browser without color-mix() (Safari < 16.2) can't evaluate, so — unlike
	// every plain-hex token above, whose LightValue() derives automatically —
	// they need a precomputed static fallback. Computed once, from the
	// catalog's own default Light values: an app that overrides ColorSurface/
	// ColorPrimary via Theme(Set(...)) does not retroactively change this
	// literal — an accepted gap for the legacy-browser tier only.
	//
	// Dark is built from the constituent tokens' NestedEnhanced(), NOT
	// EnhancedVar() and NOT a raw "var(--color-surface)" string: this Dark
	// string becomes an ARGUMENT inside the outer color-mix() here, and
	// NestedEnhanced() is the one guaranteed to contain no var() anywhere —
	// see its doc comment in tokens.go for why even a var() to an
	// always-safe property poisons an outer color-mix()/light-dark() that a
	// legacy browser can't parse.
	ColorSurfaceSunken = Token{Name: "--color-surface-sunken", Dark: "color-mix(in oklab, " + ColorSurface.NestedEnhanced() + ", " + ColorOnSurface.NestedEnhanced() + " 8%)", LightStatic: staticMix(ColorSurface.Light, ColorOnSurface.Light, 0.08)}
	ColorSelection     = Token{Name: "--color-selection", Dark: "color-mix(in oklab, " + ColorPrimary.NestedEnhanced() + ", transparent 85%)", LightStatic: FadeStatic(ColorPrimary, 0.85)}
	ColorOnSelection   = Token{Name: "--color-on-selection", Dark: ColorOnSurface.NestedEnhanced(), LightStatic: ColorOnSurface.Light}

	// ColorAccentWash is Accent faded toward transparency, the same
	// construction as ColorSelection but off ColorAccent instead of
	// ColorPrimary: a light amber tint for a hover/preview state that must
	// read as "leans toward the accent color" without claiming the solid
	// Accent fill, which is reserved for an actual committed state (e.g.
	// selection).
	ColorAccentWash = Token{Name: "--color-accent-wash", Dark: "color-mix(in oklab, " + ColorAccent.NestedEnhanced() + ", transparent 85%)", LightStatic: FadeStatic(ColorAccent, 0.85)}

	// ColorAccentHover is Accent faded only 30% toward transparency — a
	// visibly softer amber than the fully committed Accent fill, but close
	// enough in strength that ColorOnPrimary (white) still reads on it. A
	// pairing like ColorAccentWash (85% faded, nearly the page background
	// already) cannot carry a white icon at any faded strength that low;
	// this token exists for callers that need "clearly not yet committed"
	// AND "white icon" at once, which the 85% wash cannot deliver.
	ColorAccentHover = Token{Name: "--color-accent-hover", Dark: "color-mix(in oklab, " + ColorAccent.NestedEnhanced() + ", transparent 30%)", LightStatic: FadeStatic(ColorAccent, 0.30)}

	// ColorDangerWash is Danger faded toward transparency, the same
	// construction as ColorAccentWash but off ColorDanger instead of
	// ColorAccent: a light red tint for a selection state that must read as
	// "leans toward danger" without claiming the solid Danger fill, which is
	// reserved for an actual destructive commit (e.g. the confirm button).
	ColorDangerWash = Token{Name: "--color-danger-wash", Dark: "color-mix(in oklab, " + ColorDanger.NestedEnhanced() + ", transparent 85%)", LightStatic: FadeStatic(ColorDanger, 0.85)}

	FontSans = Token{Name: "--font-sans", Dark: `"Roboto", system-ui, -apple-system, sans-serif`}
	TextXs   = Token{Name: "--text-xs", Dark: "0.75rem"}
	TextSm   = Token{Name: "--text-sm", Dark: "0.875rem"}
	TextBase = Token{Name: "--text-base", Dark: "1rem"}
	TextLg   = Token{Name: "--text-lg", Dark: "1.25rem"}
	TextXl   = Token{Name: "--text-xl", Dark: "1.5rem"}
	Text2xl  = Token{Name: "--text-2xl", Dark: "2rem"}

	LeadingNormal     = Token{Name: "--leading-normal", Dark: "1.5"}
	FontWeightRegular = Token{Name: "--font-weight-regular", Dark: "400"}
	FontWeightMedium  = Token{Name: "--font-weight-medium", Dark: "500"}
	FontWeightBold    = Token{Name: "--font-weight-bold", Dark: "700"}

	Space0  = Token{Name: "--space-0", Dark: "0"}
	Space1  = Token{Name: "--space-1", Dark: "0.25rem"}
	Space2  = Token{Name: "--space-2", Dark: "0.5rem"}
	Space3  = Token{Name: "--space-3", Dark: "0.75rem"}
	Space4  = Token{Name: "--space-4", Dark: "1rem"}
	Space6  = Token{Name: "--space-6", Dark: "1.5rem"}
	Space8  = Token{Name: "--space-8", Dark: "2rem"}
	Space12 = Token{Name: "--space-12", Dark: "3rem"}

	RadiusSm   = Token{Name: "--radius-sm", Dark: "4px"}
	RadiusMd   = Token{Name: "--radius-md", Dark: "8px"}
	RadiusLg   = Token{Name: "--radius-lg", Dark: "16px"}
	RadiusFull = Token{Name: "--radius-full", Dark: "9999px"}

	ShadowSm = Token{Name: "--shadow-sm", Dark: "0 1px 2px rgba(0, 0, 0, 0.05)"}
	ShadowMd = Token{Name: "--shadow-md", Dark: "0 4px 6px rgba(0, 0, 0, 0.10)"}
	ShadowLg = Token{Name: "--shadow-lg", Dark: "0 10px 15px rgba(0, 0, 0, 0.10)"}

	DurationFast = Token{Name: "--duration-fast", Dark: "150ms"}
	DurationBase = Token{Name: "--duration-base", Dark: "250ms"}
	DurationSlow = Token{Name: "--duration-slow", Dark: "400ms"}
	EaseInOut    = Token{Name: "--ease-in-out", Dark: "cubic-bezier(0.4, 0,   0.2, 1)"}

	ZBase     = Token{Name: "--z-base", Dark: "0"}
	ZDropdown = Token{Name: "--z-dropdown", Dark: "100"}
	ZSticky   = Token{Name: "--z-sticky", Dark: "200"}
	ZModal    = Token{Name: "--z-modal", Dark: "300"}
	ZToast    = Token{Name: "--z-toast", Dark: "400"}
	ZTooltip  = Token{Name: "--z-tooltip", Dark: "500"}

	BpSm = Token{Name: "--bp-sm", Dark: "640px"}
	BpMd = Token{Name: "--bp-md", Dark: "768px"}
	BpLg = Token{Name: "--bp-lg", Dark: "1024px"}
	BpXl = Token{Name: "--bp-xl", Dark: "1280px"}

	MaxWReadable = Token{Name: "--max-w-readable", Dark: "65ch"}

	ColumnNarrow = Token{Name: "--column-narrow", Dark: "12rem"}
	ColumnMedium = Token{Name: "--column-medium", Dark: "20rem"}
	ColumnWide   = Token{Name: "--column-wide", Dark: "30rem"}

	// How far a veil blurs what is behind it. A dimmed wash alone still lets
	// the page compete for attention; softening it is what makes a dialog read
	// as the only thing in focus.
	// 4px is the common step for this in modern UI — Tailwind's backdrop-blur-sm,
	// and roughly where design systems land when they blur a scrim at all.
	// Past ~8px the page stops reading as "behind" and starts reading as
	// "broken", and it compounds with the 60% wash the veil already applies.
	VeilBlur = Token{Name: "--veil-blur", Dark: "4px"}

	// The width every chip shares — a field's legend, a row's badge — so a
	// column of them lines up instead of each one hugging its own text.
	ChipWidth = Token{Name: "--chip-width", Dark: "7rem"}

	// The height every chip shares — a field's legend, a row's badge — so a
	// chip is a box of KNOWN size instead of an emergent one. Without this the
	// height comes from font-size × line-height and two chips only match by
	// accident; with it, OnEdge can mount a chip over a border line with real
	// margins instead of a transform, which is invisible to scroll-height
	// measurement and reserves no layout space.
	ChipHeight = Token{Name: "--chip-height", Dark: "1.25rem"}

	// The height every interactive row shares — a list row, a form field —
	// so the two read as the same rhythm instead of drifting apart.
	ControlHeight = Token{Name: "--control-height", Dark: "3.125rem"}

	// Rail widths — the fixed column of a Sidebar layout.
	RailNarrow = Token{Name: "--rail-narrow", Dark: "3.5rem"}
	RailWide   = Token{Name: "--rail-wide", Dark: "12rem"}

	// Interaction intensities — how far a state moves from its base colour.

	MixHover = Token{Name: "--mix-hover", Dark: "15%"}
	MixFocus = Token{Name: "--mix-focus", Dark: "30%"}
	MixPress = Token{Name: "--mix-press", Dark: "45%"}

	// Device geometry — insets reported by the device, and the viewport
	// height that shrinks with Safari iOS' collapsing URL bar.
	SafeTop    = Token{Name: "--safe-top", Dark: "env(safe-area-inset-top, 0px)"}
	SafeRight  = Token{Name: "--safe-right", Dark: "env(safe-area-inset-right, 0px)"}
	SafeBottom = Token{Name: "--safe-bottom", Dark: "env(safe-area-inset-bottom, 0px)"}
	SafeLeft   = Token{Name: "--safe-left", Dark: "env(safe-area-inset-left, 0px)"}
	ViewportH  = Token{Name: "--viewport-h", Dark: "100dvh"}
)

var (
	SurfacePrimary    = Pair{ColorPrimary, ColorOnPrimary}
	SurfacePanel      = Pair{ColorSurface, ColorOnSurface}
	SurfaceBackground = Pair{ColorBackground, ColorOnBackground}
	SurfaceSunken     = Pair{ColorSurfaceSunken, ColorOnSurface}
	SurfaceSelected   = Pair{ColorSelection, ColorOnSelection}
	SurfaceDanger     = Pair{ColorDanger, ColorOnDanger}
	SurfaceDangerWash = Pair{ColorDangerWash, ColorOnSurface}
	SurfaceAccent     = Pair{ColorAccent, ColorOnAccent}
	SurfaceSuccess    = Pair{ColorSuccess, ColorOnSuccess}
)
