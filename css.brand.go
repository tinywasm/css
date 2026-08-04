//go:build !wasm

package css

// brandRoot declares the identity palette an app reskins to become its own:
// primary, success, danger and accent, each paired with its on-color. Kept
// apart from defaultRoots() so a white-label override touches one small
// group instead of hunting through the full token catalog.
func brandRoot() item {
	return root(
		// Brand group — Pa100T reference palette: steel blue with white text.
		declare(ColorPrimary),
		declare(ColorOnPrimary),
		declare(ColorSuccess),
		declare(ColorOnSuccess),
		declare(ColorDanger),
		declare(ColorOnDanger),
		declare(ColorAccent),
		declare(ColorOnAccent),
	)
}
