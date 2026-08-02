//go:build !wasm

package css

var (
	ColorPrimary     = Token{Name: "--color-primary", Dark: "#1b5d8c"}
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

	ColorSurfaceSunken = Token{Name: "--color-surface-sunken", Dark: "color-mix(in oklab, var(--color-surface), var(--color-on-surface) 8%)"}
	ColorSelection     = Token{Name: "--color-selection", Dark: "color-mix(in oklab, var(--color-primary), transparent 85%)"}
	ColorOnSelection   = Token{Name: "--color-on-selection", Dark: "var(--color-on-surface)"}

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
)

var (
	SurfacePrimary    = Pair{ColorPrimary, ColorOnPrimary}
	SurfacePanel      = Pair{ColorSurface, ColorOnSurface}
	SurfaceBackground = Pair{ColorBackground, ColorOnBackground}
	SurfaceSunken     = Pair{ColorSurfaceSunken, ColorOnSurface}
	SurfaceSelected   = Pair{ColorSelection, ColorOnSelection}
	SurfaceDanger     = Pair{ColorDanger, ColorOnDanger}
	SurfaceAccent     = Pair{ColorAccent, ColorOnAccent}
	SurfaceSuccess    = Pair{ColorSuccess, ColorOnSuccess}
)
