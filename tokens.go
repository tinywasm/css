//go:build !wasm

package css

// Token is a design token: a named visual decision with a fallback value.
// For static tokens, Dark holds the literal fallback and Light is empty.
// For theme tokens, Light and Dark hold the light/dark pair and Var/GetFallback
// generate light-dark(Light, Dark) automatically.
type Token struct {
	Name       string
	Light, Dark string // Dark solo = static fallback; ambos = light-dark pair
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

// Pair represents a complete surface decision: background and foreground colors coupled.
type Pair struct{ Bg, Fg Token }

// NamedPair is a descriptor used for contrast compliance verification.
type NamedPair struct {
	Name string
	Bg   ValueGetter
	Fg   ValueGetter
	Min  float64
}

// AllPairs returns the 7 functional design decision pairs for automated contrast auditing.
// SurfaceSunken and SurfaceSelected are excluded: their values are color-mix()
// expressions that resolveColor() cannot evaluate (known gap, documented in SPECS).
func AllPairs() []NamedPair {
	return []NamedPair{
		{"SurfacePrimary", ColorPrimary, ColorOnPrimary, 4.5},
		{"SurfacePanel", ColorSurface, ColorOnSurface, 4.5},
		{"SurfaceBackground", ColorBackground, ColorOnBackground, 4.5},
		{"SurfaceDanger", ColorDanger, ColorOnDanger, 4.5},
		{"SurfaceSuccess", ColorSuccess, ColorOnSuccess, 4.5},
	}
}
