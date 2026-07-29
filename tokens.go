package css

// Token is a design token: a named visual decision with a fallback value.
// Industry-standard term (W3C Design Tokens CG, Material, Carbon, Primer, Spectrum).
type Token struct{ Name, Fallback string }

// Var returns the CSS variable expression for the token, including its default fallback.
func (t Token) Var() string { return "var(" + t.Name + "," + t.Fallback + ")" }

// ValueGetter defines the interface for retrieving token metadata.
type ValueGetter interface {
	GetName() string
	GetFallback() string
}

// GetName returns the CSS property name of the token.
func (t Token) GetName() string     { return t.Name }

// GetFallback returns the default/fallback value of the token.
func (t Token) GetFallback() string { return t.Fallback }

// Token catalog — the single source of truth for all design decisions.
var (
	// Brand colors — Core identity colors

	// ColorPrimary is the primary brand color used for key call-to-actions, primary buttons, and highlighted states.
	ColorPrimary     = Token{"--color-primary", "#1b5d8c"}
	// ColorOnPrimary is the high-contrast foreground color designed to be placed on top of the primary color.
	ColorOnPrimary   = Token{"--color-on-primary", "#FFFFFF"}
	// ColorSuccess is the color used to convey success states, positive outcomes, completed steps, and confirmation messages.
	ColorSuccess     = Token{"--color-success", "#1e7a30"}
	// ColorOnSuccess is the high-contrast foreground color designed to be placed on top of the success color.
	ColorOnSuccess   = Token{"--color-on-success", "#FFFFFF"}
	// ColorDanger is the color used for destructive actions, validation errors, failure states, and critical alerts.
	ColorDanger       = Token{"--color-danger", "#ba2c0d"}
	// ColorOnDanger is the high-contrast foreground color designed to be placed on top of the danger color.
	ColorOnDanger     = Token{"--color-on-danger", "#FFFFFF"}

	// Theme colors — Adaptive layout and content layers (using light-dark responsive fallbacks)

	// ColorBackground is the default page canvas/background color, adapting seamlessly from light mode to dark mode.
	ColorBackground    = Token{"--color-background", "light-dark(#FFFFFF, #0D1117)"}
	// ColorOnBackground is the primary text and foreground color used on top of the default background.
	ColorOnBackground  = Token{"--color-on-background", "light-dark(#1C1C1E, #E6EDF3)"}
	// ColorSurface is the surface/container background color (e.g. panels, cards, inputs, buttons).
	ColorSurface       = Token{"--color-surface", "light-dark(#F2F2F7, #161B22)"}
	// ColorOnSurface is the default text and foreground color used on top of surface containers.
	ColorOnSurface     = Token{"--color-on-surface", "light-dark(#1C1C1E, #E6EDF3)"}
	// ColorOutline is the subtle boundary/border color used for separation, inputs, cards, and divider lines.
	ColorOutline       = Token{"--color-outline", "light-dark(#D1D1D6, #30363D)"}
	// ColorMuted is the secondary, lower-contrast text and icon color for secondary labels, hints, and helper texts.
	ColorMuted         = Token{"--color-muted", "light-dark(#6E6E73, #8B949E)"}

	// Typography size scale — Font-size choices following a consistent modular ratio

	// TextXs is the smallest font size (0.75rem / 12px), ideal for micro-copy, captions, and tooltips.
	TextXs   = Token{"--text-xs", "0.75rem"}
	// TextSm is the small font size (0.875rem / 14px), perfect for supporting text, labels, and helper copy.
	TextSm   = Token{"--text-sm", "0.875rem"}
	// TextBase is the standard body font size (1rem / 16px) used for main content, paragraphs, and list items.
	TextBase = Token{"--text-base", "1rem"}
	// TextLg is the large body/sub-header font size (1.25rem / 20px) for small headers or highlighted intro text.
	TextLg   = Token{"--text-lg", "1.25rem"}
	// TextXl is the secondary header font size (1.5rem / 24px) for section titles and medium headers.
	TextXl   = Token{"--text-xl", "1.5rem"}
	// Text2xl is the primary header font size (2rem / 32px) used for page titles and large display headers.
	Text2xl  = Token{"--text-2xl", "2rem"}

	// Typography weights and line-heights

	// LeadingNormal is the default, highly readable line-height ratio (1.5) for body copy and general layout.
	LeadingNormal     = Token{"--leading-normal", "1.5"}
	// FontWeightRegular represents the standard/regular weight (400) for body copy.
	FontWeightRegular = Token{"--font-weight-regular", "400"}
	// FontWeightMedium represents the medium/semi-bold weight (500) for emphasizeable text, buttons, and subheaders.
	FontWeightMedium  = Token{"--font-weight-medium", "500"}
	// FontWeightBold represents the bold weight (700) for strong hierarchy, primary headers, and page titles.
	FontWeightBold    = Token{"--font-weight-bold", "700"}

	// Spacing scale (4px grid) — Predictable margins, padding, and gap spacing

	// Space0 represents zero spacing (0px), useful for resetting margins/paddings explicitly.
	Space0  = Token{"--space-0", "0"}
	// Space1 is the smallest non-zero gap/padding (0.25rem / 4px), perfect for tight spacing or icon-text gaps.
	Space1  = Token{"--space-1", "0.25rem"}
	// Space2 is the secondary gap/padding (0.5rem / 8px), ideal for button padding or dense item gaps.
	Space2  = Token{"--space-2", "0.5rem"}
	// Space3 is the intermediate gap/padding (0.75rem / 12px), commonly used for standard list item gaps.
	Space3  = Token{"--space-3", "0.75rem"}
	// Space4 is the baseline layout spacing (1rem / 16px) for card paddings, page margins, and standard gaps.
	Space4  = Token{"--space-4", "1rem"}
	// Space6 is a larger layout spacing (1.5rem / 24px) for generous grouping of layout sections.
	Space6  = Token{"--space-6", "1.5rem"}
	// Space8 is a major section spacing (2rem / 32px) for separating major blocks of a page.
	Space8  = Token{"--space-8", "2rem"}
	// Space12 is the largest spacing primitive (3rem / 48px) used for large viewport boundaries.
	Space12 = Token{"--space-12", "3rem"}

	// Border radius scale — Standardized rounding of corners for consistent component feel

	// RadiusSm is a subtle corner rounding (4px) suited for small elements like checkboxes, badges, and tooltips.
	RadiusSm   = Token{"--radius-sm", "4px"}
	// RadiusMd is the standard corner rounding (8px) used for buttons, inputs, alerts, and card boundaries.
	RadiusMd   = Token{"--radius-md", "8px"}
	// RadiusLg is a generous corner rounding (16px) for larger dialogs, modals, panels, and dashboard widgets.
	RadiusLg   = Token{"--radius-lg", "16px"}
	// RadiusFull is a fully circular border-radius (9999px) for pill buttons, avatar images, and circular badges.
	RadiusFull = Token{"--radius-full", "9999px"}

	// Elevation shadows — Visual hierarchy and depth elevation scale

	// ShadowSm is a soft, minimal drop shadow for low-elevation cards, buttons, and inline dropdown containers.
	ShadowSm = Token{"--shadow-sm", "0 1px 2px rgba(0, 0, 0, 0.05)"}
	// ShadowMd is a medium-depth shadow for hovering components like dropdowns, popovers, and elevated buttons.
	ShadowMd = Token{"--shadow-md", "0 4px 6px rgba(0, 0, 0, 0.10)"}
	// ShadowLg is a deep, immersive drop shadow reserved for primary overlay elements like modals, toast notifications, and dialogs.
	ShadowLg = Token{"--shadow-lg", "0 10px 15px rgba(0, 0, 0, 0.10)"}

	// Motion transitions and easing curves

	// DurationFast is a rapid transition time (150ms) for high-frequency actions like hover, focus, and toggle states.
	DurationFast = Token{"--duration-fast", "150ms"}
	// DurationBase is the standard layout transition duration (250ms) for panels, page transitions, and drawer openings.
	DurationBase = Token{"--duration-base", "250ms"}
	// DurationSlow is an extended transition duration (400ms) for major visual transformations or dramatic reveals.
	DurationSlow = Token{"--duration-slow", "400ms"}
	// EaseInOut is a smooth, natural-feeling cubic-bezier easing curve for UI element animations.
	EaseInOut    = Token{"--ease-in-out", "cubic-bezier(0.4, 0,   0.2, 1)"}

	// Z-index stacking contract — Strictly layered z-axis positions preventing overlap collisions

	// ZBase is the default baseline layer (0).
	ZBase     = Token{"--z-base", "0"}
	// ZDropdown is the layer (100) allocated to floating selectors, dropdown menus, and popovers.
	ZDropdown = Token{"--z-dropdown", "100"}
	// ZSticky is the layer (200) for persistent navigation headers or sticky control bars.
	ZSticky   = Token{"--z-sticky", "200"}
	// ZModal is the layer (300) reserved for blocking overlay boxes, dialog sheets, and lightboxes.
	ZModal    = Token{"--z-modal", "300"}
	// ZToast is the layer (400) for floating message banners, feedback notifications, and alerts.
	ZToast    = Token{"--z-toast", "400"}
	// ZTooltip is the highest layer (500) dedicated strictly to fleeting helper tooltips.
	ZTooltip  = Token{"--z-tooltip", "500"}

	// Breakpoints — Viewport width thresholds

	// BpSm is the small screen/mobile viewport threshold width (640px).
	BpSm = Token{"--bp-sm", "640px"}
	// BpMd is the medium screen/tablet viewport threshold width (768px).
	BpMd = Token{"--bp-md", "768px"}
	// BpLg is the large screen/desktop viewport threshold width (1024px).
	BpLg = Token{"--bp-lg", "1024px"}
	// BpXl is the extra-large/widescreen desktop viewport threshold width (1280px).
	BpXl = Token{"--bp-xl", "1280px"}

	// Container widths

	// MaxWReadable defines the maximum comfortable horizontal length (65ch) for readable body paragraphs.
	MaxWReadable = Token{"--max-w-readable", "65ch"}

	// Grid columns

	// ColumnNarrow represents a narrow column layout unit width (12rem).
	ColumnNarrow = Token{"--column-narrow", "12rem"}
	// ColumnMedium represents a medium column layout unit width (20rem).
	ColumnMedium = Token{"--column-medium", "20rem"}
	// ColumnWide represents a wide column layout unit width (30rem).
	ColumnWide   = Token{"--column-wide", "30rem"}
)

// Pair represents a complete surface decision: background and foreground colors coupled.
// A background is never declared without its foreground, preventing the type of bug
// where a text color is accidentally used as a panel background.
type Pair struct{ Bg, Fg Token }

var (
	// SurfacePrimary couples the primary brand color with its appropriate on-primary high contrast text.
	SurfacePrimary  = Pair{ColorPrimary, ColorOnPrimary}
	// SurfacePanel couples the general container surface color with the standard on-surface body text.
	SurfacePanel    = Pair{ColorSurface, ColorOnSurface}
	// SurfaceBackground couples the default body background canvas with its standard on-background body text.
	SurfaceBackground = Pair{ColorBackground, ColorOnBackground}
	// SurfaceDanger couples the destructive/danger red color with its high contrast white text.
	SurfaceDanger   = Pair{ColorDanger, ColorOnDanger}
	// SurfaceSuccess couples the affirmative success green color with its high contrast white text.
	SurfaceSuccess  = Pair{ColorSuccess, ColorOnSuccess}
)

// NamedPair is a descriptor used for contrast compliance verification.
type NamedPair struct {
	Name string
	Bg   ValueGetter
	Fg   ValueGetter
	Min  float64
}

// AllPairs returns the 5 functional design decision pairs for automated contrast auditing.
func AllPairs() []NamedPair {
	return []NamedPair{
		{"SurfacePrimary", ColorPrimary, ColorOnPrimary, 4.5},
		{"SurfacePanel", ColorSurface, ColorOnSurface, 4.5},
		{"SurfaceBackground", ColorBackground, ColorOnBackground, 4.5},
		{"SurfaceDanger", ColorDanger, ColorOnDanger, 4.5},
		{"SurfaceSuccess", ColorSuccess, ColorOnSuccess, 4.5},
	}
}
