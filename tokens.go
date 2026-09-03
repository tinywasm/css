//go:build !wasm

package css

// Token is a design token: a named visual decision with a fallback value.
// For static tokens, Dark holds the literal fallback and Light is empty.
// For theme tokens, Light and Dark hold the light/dark pair and Var/GetFallback
// generate light-dark(Light, Dark) automatically.
type Token struct {
	Name        string
	Light, Dark string // Dark solo = static fallback; ambos = light-dark pair

	// LightStatic overrides LightValue() for a token whose Dark (or Light)
	// holds a live expression — color-mix(...) or var(--other-token) — that a
	// browser without color-mix()/light-dark() support cannot evaluate.
	// Leave empty for every plain-hex token: LightValue() derives it
	// automatically from Light (or Dark). Only the handful of composed color
	// tokens in catalog.go need to set this, precomputed via color.Mix at
	// package-init time.
	LightStatic string
}

// Var returns the CSS variable expression for the token, including its default fallback.
func (t Token) Var() string {
	return "var(" + t.Name + "," + t.fallback() + ")"
}

// ValueGetter defines the interface for retrieving token metadata.
type ValueGetter interface {
	GetName() string
	GetFallback() string
}

// GetName returns the CSS property name of the token.
func (t Token) GetName() string { return t.Name }

// GetFallback returns the default/fallback value of the token.
func (t Token) GetFallback() string { return t.fallback() }

func (t Token) fallback() string {
	if t.Light == "" {
		return t.Dark
	}
	return "light-dark(" + t.Light + ", " + t.Dark + ")"
}

// LightValue returns the token's browser-safe static value: no light-dark()
// or color-mix() involved, ever. A browser that can't evaluate those (Safari
// < 17.5 / < 16.2 — iOS 15-era Safari on older hardware, among others) gets
// this value permanently, regardless of system or app theme — the light
// theme, never dark. Callers emit it as the first of a double declaration,
// immediately followed by EnhancedVar(), never Var() — see EnhancedVar for
// why that distinction is load-bearing.
func (t Token) LightValue() string {
	if t.LightStatic != "" {
		return t.LightStatic
	}
	if t.Light != "" {
		return t.Light
	}
	return t.Dark
}

// LightVarName and DarkVarName are the custom-property names of a theme
// token's plain, browser-safe halves — declared separately from Name itself
// (see declareSplit in dsl.go). Not used by EnhancedVar/NestedEnhanced below
// (see their comments for why a var() reference, even to one of these
// always-valid halves, is not safe inside a light-dark()/color-mix() call);
// they exist as public custom properties for an app's OWN CSS to build the
// same double-declaration pattern against. Only meaningful for a theme token
// (Light != "") — a static token has no light/dark split to name.
func (t Token) LightVarName() string { return t.Name + "-light" }
func (t Token) DarkVarName() string  { return t.Name + "-dark" }

// ImageVarName names the token's optional background-image companion
// property — see SetGradient. Declared unconditionally by widget/style as
// `background-image: var(--x-image, none);` alongside every surface's
// background-color, so it stays inert (the "none" fallback) until an app
// sets it, and paints over the solid color without replacing it when it
// does — no separate "is this surface a gradient" flag to keep in sync.
func (t Token) ImageVarName() string { return t.Name + "-image" }

// ImageStopsVarName names the companion that holds only the gradient's colour
// stops ("var(--from), var(--to)"), without the angle — set by SetGradient
// alongside ImageVarName(). It lets one surface repaint the SAME theme gradient
// at a different angle (widget/style's GradientAngle) instead of every surface
// re-origining the one baked-in direction. Unset until an app calls
// SetGradient, so a `linear-gradient(<angle>, var(--x-image-stops))` that
// references it is invalid — and therefore ignored — until then.
func (t Token) ImageStopsVarName() string { return t.Name + "-image-stops" }

// EnhancedVar returns the value used as the SECOND half of a double
// declaration when the token is the declaration's ENTIRE top-level value
// (e.g. `background-color: <this>;`) — see NestedEnhanced for the
// alternative required when it is instead one ARGUMENT inside another
// light-dark()/color-mix() call.
//
// It is Var(), and that is the whole point. It used to bake the catalog's
// light-dark(<light>,<dark>) in as a literal, because t.Name's registered
// :root value WAS that expression and reaching it through var() would make
// the declaration invalid at computed-value time on a browser that cannot
// evaluate it — which, per the Custom Properties cascade, falls to the
// property's initial value rather than to an earlier sibling declaration.
// The literal dodged that by failing at PARSE time instead, which does leave
// the sibling standing.
//
// It also meant no app could change any of these colours: a literal baked at
// Go build time is not reachable from Theme(Set(...)), so a site could
// declare its entire palette and still be painted in the library's default
// one. That is not a corner of the API — it is the background, text and
// border of every surface in the ecosystem.
//
// The poison is now removed at the source instead: :root declares the static
// value unguarded and re-declares the adaptive one inside ModernColorSupport
// (see declare/enhancedDecls in dsl.go), so t.Name always holds something
// every engine can compute. var() is therefore safe, and an app's override —
// which lands after both — finally reaches the declaration.
// It is the BARE var() — no fallback. Var()'s fallback is the adaptive
// expression, and a fallback still counts as a colour function sitting inside
// a var() as far as any reader (or drift guard) can tell, even though it is
// unreachable while the property is declared. It is also genuinely redundant
// here: :root declares this property unconditionally, so the fallback can
// never be substituted. Leaving it out keeps the emitted declaration exactly
// as safe and says so plainly.
func (t Token) EnhancedVar() string {
	return "var(" + t.Name + ")"
}

// NestedEnhanced is EnhancedVar's counterpart for a token used as one
// ARGUMENT inside another light-dark()/color-mix() call (see catalog.go's
// composed tokens, and css.go's mixToward) rather than as a declaration's
// own top-level value.
//
// The distinction matters because the "any var() anywhere defers validity
// to computed-value time" rule (see EnhancedVar) applies to the OUTER
// declaration as a whole, not just to the specific sub-expression that
// contains it: a static token's Var() is completely safe used on its own
// (EnhancedVar returns it, correctly, for exactly that reason) but the
// moment that SAME var() call sits inside an outer color-mix()/light-dark()
// that a legacy browser cannot parse, its mere presence defers the ENTIRE
// containing declaration to computed-value time too — even though the
// nested var() would itself have resolved just fine. That silently turns
// the outer function's parse-time failure into the same
// falls-to-initial-value bug EnhancedVar exists to avoid. NestedEnhanced
// closes that gap by never using var(), for any token kind.
func (t Token) NestedEnhanced() string {
	if t.Light != "" {
		return "light-dark(" + t.Light + "," + t.Dark + ")"
	}
	return t.Dark
}

// Pair represents a complete surface decision: background and foreground colors coupled.
type Pair struct{ Bg, Fg Token }

// NamedPair is a descriptor used for contrast compliance verification.
type NamedPair struct {
	Name string
	Bg   ValueGetter
	Fg   ValueGetter
	Min  float64
}

// AllPairs returns the 6 active functional design decision pairs for automated contrast auditing.
// SurfaceSunken, SurfaceSelected and SurfaceDangerWash are excluded (remaining 3 of the 9 total pairs):
// their values are color-mix() expressions that resolveColor() cannot evaluate.
// Evaluating color-mix(in oklab, ...) is a known gap, pending resolution in
// github.com/tinywasm/color (see color/docs/PLAN.md §3.3).
func AllPairs() []NamedPair {
	return []NamedPair{
		{"SurfacePrimary", ColorPrimary, ColorOnPrimary, 4.5},
		{"SurfacePanel", ColorSurface, ColorOnSurface, 4.5},
		{"SurfaceBackground", ColorBackground, ColorOnBackground, 4.5},
		{"SurfaceDanger", ColorDanger, ColorOnDanger, 4.5},
		{"SurfaceSuccess", ColorSuccess, ColorOnSuccess, 4.5},
		{"SurfaceAccent", ColorAccent, ColorOnAccent, 4.5},
	}
}
