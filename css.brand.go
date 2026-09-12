//go:build !wasm

package css

// brandRoot declares the identity palette an app reskins to become its own:
// primary, success, danger and accent, each paired with its on-color — plus
// ColorPrimary's default gradient partner. Kept apart from defaultRoots() so
// a white-label override touches one small group instead of hunting through
// the full token catalog.
func brandRoot() item {
	decls := []decl{
		// Brand group — Pa100T reference palette: steel blue with white text.
		declare(ColorPrimary),
		declare(ColorOnPrimary),
		declare(ColorSuccess),
		declare(ColorOnSuccess),
		declare(ColorDanger),
		declare(ColorOnDanger),
		declare(ColorAccent),
		declare(ColorOnAccent),
		declare(ColorPrimaryGradient),
	}
	decls = append(decls, defaultGradient(ColorPrimary, "135deg", ColorPrimary, ColorPrimaryGradient)...)
	return root(decls...)
}

// defaultGradient declares t's background-image companion (and its stops
// companion, ImageStopsVarName) as a permanent catalog default — the same
// shape SetGradient's Override produces for a one-off app override, but
// baked in so every consumer gets it without calling Theme(...). Theme()
// still appends any app override after this block, and CSS custom
// properties resolve last-declaration-wins, so Theme(SetGradient(...)) or
// Theme(ClearGradient(...)) still wins exactly like it does today.
func defaultGradient(t Token, angle string, from, to Token) []decl {
	stops := from.Var() + ", " + to.Var()
	return []decl{
		{t.ImageVarName(), "linear-gradient(" + angle + ", " + stops + ")"},
		{t.ImageStopsVarName(), stops},
	}
}
