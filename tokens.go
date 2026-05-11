package css

// Class is a CSS class name. Shared by HTML emission (WASM) and CSS emission (SSR).
// Only the string identity crosses the WASM boundary; pseudo-class helpers
// (Hover/Focus/Disabled) live in dsl.go (!wasm) because they only feed Rule().
type Class string

func (c Class) String() string { return string(c) }

// Token is a design token: a named visual decision with a fallback value.
// Industry-standard term (W3C Design Tokens CG, Material, Carbon, Primer, Spectrum).
type Token struct{ Name, Fallback string }

func (t Token) Var() string { return "var(" + t.Name + "," + t.Fallback + ")" }

// Token catalog — every token from the legacy theme.css.
var (
	// Brand colors
	ColorPrimary     = Token{"--color-primary", "#00ADD8"}
	ColorOnPrimary   = Token{"--color-on-primary", "#1C1C1E"}
	ColorSecondary   = Token{"--color-secondary", "#654FF0"}
	ColorOnSecondary = Token{"--color-on-secondary", "#FFFFFF"}
	ColorSuccess     = Token{"--color-success", "#3FB950"}
	ColorError       = Token{"--color-error", "#E34F26"}

	// Theme — active layer (consumed by components)
	ColorBackground = Token{"--color-background", "#FFFFFF"}
	ColorSurface    = Token{"--color-surface", "#F2F2F7"}
	ColorOnSurface  = Token{"--color-on-surface", "#1C1C1E"}
	ColorMuted      = Token{"--color-muted", "#6E6E73"}
	ColorHover      = Token{"--color-hover", "#B8860B"}

	// Theme — source layer (apps redeclare these for rebrand)
	ColorBackgroundLight = Token{"--color-background-light", "#FFFFFF"}
	ColorBackgroundDark  = Token{"--color-background-dark", "#0D1117"}
	ColorSurfaceLight    = Token{"--color-surface-light", "#F2F2F7"}
	ColorSurfaceDark     = Token{"--color-surface-dark", "#161B22"}
	ColorOnSurfaceLight  = Token{"--color-on-surface-light", "#1C1C1E"}
	ColorOnSurfaceDark   = Token{"--color-on-surface-dark", "#E6EDF3"}
	ColorMutedLight      = Token{"--color-muted-light", "#6E6E73"}
	ColorMutedDark       = Token{"--color-muted-dark", "#8B949E"}
	ColorHoverLight      = Token{"--color-hover-light", "#B8860B"}
	ColorHoverDark       = Token{"--color-hover-dark", "#F7DF1E"}

	// Typography size scale
	TextXs   = Token{"--text-xs", "0.75rem"}
	TextSm   = Token{"--text-sm", "0.875rem"}
	TextBase = Token{"--text-base", "1rem"}
	TextLg   = Token{"--text-lg", "1.25rem"}
	TextXl   = Token{"--text-xl", "1.5rem"}
	Text2xl  = Token{"--text-2xl", "2rem"}

	// Line-height / weight / tracking
	LeadingTight      = Token{"--leading-tight", "1.25"}
	LeadingNormal     = Token{"--leading-normal", "1.5"}
	LeadingRelaxed    = Token{"--leading-relaxed", "1.75"}
	FontWeightRegular = Token{"--font-weight-regular", "400"}
	FontWeightMedium  = Token{"--font-weight-medium", "500"}
	FontWeightBold    = Token{"--font-weight-bold", "700"}
	TrackingTight     = Token{"--tracking-tight", "-0.02em"}
	TrackingNormal    = Token{"--tracking-normal", "0"}
	TrackingWide      = Token{"--tracking-wide", "0.05em"}

	// Spacing (4px grid)
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
	ShadowSm = Token{"--shadow-sm", "0 1px 2px rgba(0,0,0,0.05)"}
	ShadowMd = Token{"--shadow-md", "0 4px 6px rgba(0,0,0,0.1)"}
	ShadowLg = Token{"--shadow-lg", "0 10px 15px rgba(0,0,0,0.1)"}
	ShadowXl = Token{"--shadow-xl", "0 20px 25px rgba(0,0,0,0.15)"}

	// Motion
	DurationFast = Token{"--duration-fast", "150ms"}
	DurationBase = Token{"--duration-base", "250ms"}
	DurationSlow = Token{"--duration-slow", "400ms"}
	EaseIn       = Token{"--ease-in", "cubic-bezier(0.4,0,1,1)"}
	EaseOut      = Token{"--ease-out", "cubic-bezier(0,0,0.2,1)"}
	EaseInOut    = Token{"--ease-in-out", "cubic-bezier(0.4,0,0.2,1)"}

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
	MaxWProse   = Token{"--max-w-prose", "65ch"}
	MaxWContent = Token{"--max-w-content", "1200px"}
	MaxWScreen  = Token{"--max-w-screen", "1440px"}
)
