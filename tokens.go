package css

// Token is a design token: a named visual decision with a fallback value.
// Industry-standard term (W3C Design Tokens CG, Material, Carbon, Primer, Spectrum).
type Token struct{ Name, Fallback string }

func (t Token) Var() string { return "var(" + t.Name + "," + t.Fallback + ")" }

// Source is a theme source token: a settable-only visual choice that does not support Var().
type Source struct{ Name, Fallback string }

type ValueGetter interface {
	GetName() string
	GetFallback() string
}

func (t Token) GetName() string     { return t.Name }
func (t Token) GetFallback() string { return t.Fallback }

func (s Source) GetName() string     { return s.Name }
func (s Source) GetFallback() string { return s.Fallback }

// Token catalog — every token from the legacy theme.css.
var (
	// Brand colors
	ColorPrimary     = Token{"--color-primary", "#1b5d8c"}
	ColorOnPrimary   = Token{"--color-on-primary", "#FFFFFF"}
	ColorSecondary   = Token{"--color-secondary", "#654FF0"}
	ColorOnSecondary = Token{"--color-on-secondary", "#FFFFFF"}
	ColorSuccess     = Token{"--color-success", "#1e7a30"}
	ColorOnSuccess   = Token{"--color-on-success", "#FFFFFF"}
	ColorError       = Token{"--color-error", "#ba2c0d"}
	ColorOnError     = Token{"--color-on-error", "#FFFFFF"}

	// Theme — active layer (consumed by components)
	ColorBackground    = Token{"--color-background", "#FFFFFF"}
	ColorSurface       = Token{"--color-surface", "#F2F2F7"}
	ColorSurfaceSunken = Token{"--color-surface-sunken", "#E5E5EA"}
	ColorOnSurface     = Token{"--color-on-surface", "#1C1C1E"}
	ColorOutline       = Token{"--color-outline", "#D1D1D6"}
	ColorMuted         = Token{"--color-muted", "#6E6E73"}
	ColorHover         = Token{"--color-hover", "#B8860B"}
	ColorSelection     = Token{"--color-selection", "#f5a623"}    // selected/active highlight
	ColorOnSelection   = Token{"--color-on-selection", "#1C1C1E"} // text on a selected row
	ColorDisabled      = Token{"--color-disabled", "#E5E5EA"}
	ColorOnDisabled    = Token{"--color-on-disabled", "#666666"}

	// Theme — source layer (apps redeclare these for rebrand)
	ColorBackgroundLight    = Source{"--color-background-light", "#FFFFFF"}
	ColorBackgroundDark     = Source{"--color-background-dark", "#0D1117"}
	ColorSurfaceLight       = Source{"--color-surface-light", "#F2F2F7"}
	ColorSurfaceDark        = Source{"--color-surface-dark", "#161B22"}
	ColorSurfaceSunkenLight = Source{"--color-surface-sunken-light", "#E5E5EA"}
	ColorSurfaceSunkenDark  = Source{"--color-surface-sunken-dark", "#21262D"}
	ColorOnSurfaceLight     = Source{"--color-on-surface-light", "#1C1C1E"}
	ColorOnSurfaceDark      = Source{"--color-on-surface-dark", "#E6EDF3"}
	ColorOutlineLight       = Source{"--color-outline-light", "#D1D1D6"}
	ColorOutlineDark        = Source{"--color-outline-dark", "#30363D"}
	ColorMutedLight         = Source{"--color-muted-light", "#6E6E73"}
	ColorMutedDark          = Source{"--color-muted-dark", "#8B949E"}
	ColorHoverLight         = Source{"--color-hover-light", "#B8860B"}
	ColorHoverDark          = Source{"--color-hover-dark", "#F7DF1E"}
	ColorSelectionLight     = Source{"--color-selection-light", "#f5a623"}
	ColorSelectionDark      = Source{"--color-selection-dark", "#9e6a2e"}
	ColorOnSelectionLight   = Source{"--color-on-selection-light", "#1C1C1E"}
	ColorOnSelectionDark    = Source{"--color-on-selection-dark", "#FFFFFF"}
	ColorDisabledLight      = Source{"--color-disabled-light", "#E5E5EA"}
	ColorDisabledDark       = Source{"--color-disabled-dark", "#21262D"}
	ColorOnDisabledLight    = Source{"--color-on-disabled-light", "#5e5e64"}
	ColorOnDisabledDark     = Source{"--color-on-disabled-dark", "#8b949e"}

	// Typography size scale
	TextXs   = Token{"--text-xs", "0.75rem"}
	TextSm   = Token{"--text-sm", "0.875rem"}
	TextBase = Token{"--text-base", "1rem"}
	TextLg   = Token{"--text-lg", "1.25rem"}
	TextXl   = Token{"--text-xl", "1.5rem"}
	Text2xl  = Token{"--text-2xl", "2rem"}

	// Line-height / weight
	LeadingNormal     = Token{"--leading-normal", "1.5"}
	FontWeightRegular = Token{"--font-weight-regular", "400"}
	FontWeightMedium  = Token{"--font-weight-medium", "500"}
	FontWeightBold    = Token{"--font-weight-bold", "700"}

	// Spacing (4px grid)
	Space0  = Token{"--space-0", "0"}
	Space1  = Token{"--space-1", "0.25rem"}
	Space2  = Token{"--space-2", "0.5rem"}
	Space3  = Token{"--space-3", "0.75rem"}
	Space4  = Token{"--space-4", "1rem"}
	Space6  = Token{"--space-6", "1.5rem"}
	Space8  = Token{"--space-8", "2rem"}
	Space12 = Token{"--space-12", "3rem"}

	// Border radius
	RadiusSm   = Token{"--radius-sm", "4px"}
	RadiusMd   = Token{"--radius-md", "8px"}
	RadiusLg   = Token{"--radius-lg", "16px"}
	RadiusFull = Token{"--radius-full", "9999px"}

	// Elevation
	ShadowSm = Token{"--shadow-sm", "0 1px 2px rgba(0, 0, 0, 0.05)"}
	ShadowMd = Token{"--shadow-md", "0 4px 6px rgba(0, 0, 0, 0.10)"}
	ShadowLg = Token{"--shadow-lg", "0 10px 15px rgba(0, 0, 0, 0.10)"}

	// Motion
	DurationFast = Token{"--duration-fast", "150ms"}
	DurationBase = Token{"--duration-base", "250ms"}
	DurationSlow = Token{"--duration-slow", "400ms"}
	EaseInOut    = Token{"--ease-in-out", "cubic-bezier(0.4, 0,   0.2, 1)"}

	// Z-index
	ZBase     = Token{"--z-base", "0"}
	ZDropdown = Token{"--z-dropdown", "100"}
	ZSticky   = Token{"--z-sticky", "200"}
	ZModal    = Token{"--z-modal", "300"}
	ZToast    = Token{"--z-toast", "400"}
	ZTooltip  = Token{"--z-tooltip", "500"}

	// Breakpoints
	BpSm = Token{"--bp-sm", "640px"}
	BpMd = Token{"--bp-md", "768px"}
	BpLg = Token{"--bp-lg", "1024px"}
	BpXl = Token{"--bp-xl", "1280px"}

	// Container widths
	MaxWReadable = Token{"--max-w-readable", "65ch"}

	// Grid columns
	ColumnNarrow = Token{"--column-narrow", "12rem"}
	ColumnMedium = Token{"--column-medium", "20rem"}
	ColumnWide   = Token{"--column-wide", "30rem"}

	// Focus ring
	ColorFocusRing = Token{"--color-focus-ring", "#1b5d8c"}

	// Interactive colors — hover, focus, press states
	ColorPrimaryHover = Token{"--color-primary-hover", "#164d74"}
	ColorPrimaryFocus = Token{"--color-primary-focus", "#123e5d"}
	ColorPrimaryPress = Token{"--color-primary-press", "#0e304a"}

	ColorSecondaryHover = Token{"--color-secondary-hover", "#523ee0"}
	ColorSecondaryFocus = Token{"--color-secondary-focus", "#422fd1"}
	ColorSecondaryPress = Token{"--color-secondary-press", "#3321c2"}

	ColorSuccessHover = Token{"--color-success-hover", "#186327"}
	ColorSuccessFocus = Token{"--color-success-focus", "#134d1e"}
	ColorSuccessPress = Token{"--color-success-press", "#0e3816"}

	ColorDangerHover = Token{"--color-danger-hover", "#9a240b"}
	ColorDangerFocus = Token{"--color-danger-focus", "#7a1c09"}
	ColorDangerPress = Token{"--color-danger-press", "#5b1507"}

	ColorErrorHover = Token{"--color-error-hover", "#9a240b"}
	ColorErrorFocus = Token{"--color-error-focus", "#7a1c09"}
	ColorErrorPress = Token{"--color-error-press", "#5b1507"}

	ColorWarningHover = Token{"--color-warning-hover", "#b45309"}
	ColorWarningFocus = Token{"--color-warning-focus", "#92400e"}
	ColorWarningPress = Token{"--color-warning-press", "#78350f"}

	ColorInfoHover = Token{"--color-info-hover", "#0369a1"}
	ColorInfoFocus = Token{"--color-info-focus", "#075985"}
	ColorInfoPress = Token{"--color-info-press", "#0c4a6e"}

	ColorNeutralHover = Token{"--color-neutral-hover", "#374151"}
	ColorNeutralFocus = Token{"--color-neutral-focus", "#1f2937"}
	ColorNeutralPress = Token{"--color-neutral-press", "#111827"}

	ColorMutedHover = Token{"--color-muted-hover", "#4b5563"}
	ColorMutedFocus = Token{"--color-muted-focus", "#374151"}
	ColorMutedPress = Token{"--color-muted-press", "#1f2937"}
)

// Pair represents a complete surface decision: background and foreground colors coupled.
// A background is never declared without its foreground, preventing the type of bug
// where a text color is accidentally used as a panel background.
type Pair struct{ Bg, Fg Token }

var (
	SurfacePrimary  = Pair{ColorPrimary, ColorOnPrimary}
	SurfacePanel    = Pair{ColorSurface, ColorOnSurface}
	SurfaceSunken   = Pair{ColorSurfaceSunken, ColorOnSurface}
	SurfaceSelected = Pair{ColorSelection, ColorOnSelection}
	SurfaceDanger   = Pair{ColorError, ColorOnError}
	SurfaceSuccess  = Pair{ColorSuccess, ColorOnSuccess}
	SurfaceDisabled = Pair{ColorDisabled, ColorOnDisabled}
)

type NamedPair struct {
	Name string
	Bg   ValueGetter
	Fg   ValueGetter
	Min  float64
}

func AllPairs() []NamedPair {
	return []NamedPair{
		{"SurfacePrimary", ColorPrimary, ColorOnPrimary, 4.5},
		{"SurfacePanel", ColorSurface, ColorOnSurface, 4.5},
		{"SurfaceSunken", ColorSurfaceSunken, ColorOnSurface, 4.5},
		{"SurfaceSelected", ColorSelection, ColorOnSelection, 4.5},
		{"SurfaceDanger", ColorError, ColorOnError, 4.5},
		{"SurfaceSuccess", ColorSuccess, ColorOnSuccess, 4.5},
		{"SurfaceDisabled", ColorDisabled, ColorOnDisabled, 4.5},

		// Light twins
		{"SurfacePanelLight", ColorSurfaceLight, ColorOnSurfaceLight, 4.5},
		{"SurfaceSunkenLight", ColorSurfaceSunkenLight, ColorOnSurfaceLight, 4.5},
		{"SurfaceSelectedLight", ColorSelectionLight, ColorOnSelectionLight, 4.5},
		{"SurfaceDisabledLight", ColorDisabledLight, ColorOnDisabledLight, 4.5},

		// Dark twins
		{"SurfacePanelDark", ColorSurfaceDark, ColorOnSurfaceDark, 4.5},
		{"SurfaceSunkenDark", ColorSurfaceSunkenDark, ColorOnSurfaceDark, 4.5},
		{"SurfaceSelectedDark", ColorSelectionDark, ColorOnSelectionDark, 4.5},
		{"SurfaceDisabledDark", ColorDisabledDark, ColorOnDisabledDark, 4.5},
	}
}
