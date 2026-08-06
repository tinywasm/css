//go:build !wasm

package css

import (
	"os/exec"
	"strings"
	"testing"
)

var composedTokens = []Token{ColorSurfaceSunken, ColorSelection, ColorOnSelection}

// A composed token's fallback is read as an ARGUMENT inside another
// color-mix()/light-dark() call (see catalog.go's own construction, and
// css.go's mixToward) or, via declare(), as a declaration's own top-level
// :root value — either way it must contain NO var() anywhere.
//
// This inverts what this test checked before EnhancedVar/NestedEnhanced
// existed (composed tokens used to read other tokens live, via
// "var(--color-surface)"). The reason is a CSS Custom Properties rule that
// is easy to miss: ANY var() anywhere in a declaration's value — even one
// that would itself always resolve to something valid — defers validity
// checking for the WHOLE declaration to computed-value time. A browser that
// can't parse the outer color-mix()/light-dark() then falls to the
// property's initial value instead of an earlier sibling declaration,
// exactly the bug the double-declaration Safari-legacy fallback exists to
// avoid (see docs/PLAN or the git history around the iPhone-7-all-blue
// investigation). Confirmed empirically, not just by spec reading: a
// minimal repro with the unsupported function's arguments as var()
// references reproduces the failure; the same repro with plain literals
// does not.
//
// The real cost: an app's Theme(Set(...))/SetTheme() does not reach these
// three tokens' Dark — only republishing this package's catalog defaults
// does. Accepted, since dark/light switching for THESE composed derivations
// only matters on browsers new enough to have color-mix() in the first
// place.
func TestComposedTokensContainNoVar(t *testing.T) {
	for _, tok := range composedTokens {
		if strings.Contains(tok.GetFallback(), "var(") {
			t.Errorf("%s must contain no var() anywhere — it is read as an argument inside an outer color-mix()/light-dark() call, and any var() there defers the WHOLE declaration to computed-value time, breaking the legacy-Safari fallback: %s", tok.Name, tok.GetFallback())
		}
	}
}

// Composed tokens reference variables that must exist in the same :root.
func TestComposedTokensAreDeclared(t *testing.T) {
	emitted := RootCSS().String()
	for _, tok := range composedTokens {
		if !strings.Contains(emitted, tok.Name+":") {
			t.Errorf("%s is not declared in RootCSS(); its inner var() has no guaranteed referent", tok.Name)
		}
	}
}

// An interaction derivation must not freeze the theme direction: the mixer is
// light-dark(black, white), never a bare black or white.
func TestInteractionDerivationIsThemeAware(t *testing.T) {
	for name, got := range map[string]string{
		"Hover": Hover(ColorPrimary),
		"Focus": Focus(ColorPrimary),
		"Press": Press(ColorPrimary),
	} {
		if !strings.Contains(got, "light-dark(black, white)") {
			t.Errorf("%s freezes the theme direction: %s", name, got)
		}
	}
}

// Hover/Focus/Press must bake their intensity as a literal, not a live
// var() reference to --mix-hover/-focus/-press — see
// TestComposedTokensContainNoVar for why: mixToward's whole point is a
// color-mix() call a legacy browser can't parse, and ANY var() anywhere in
// that declaration (even to MixHover, an always-safe static token) would
// defer the WHOLE thing to computed-value time and break the fallback. Cost
// accepted: Theme(Set(MixHover, ...)) does not retune these three — only
// republishing this package's catalog default does.
func TestInteractionIntensityIsLiteral(t *testing.T) {
	for _, c := range []struct {
		name string
		got  string
		tok  Token
	}{
		{"Hover", Hover(ColorPrimary), MixHover},
		{"Focus", Focus(ColorPrimary), MixFocus},
		{"Press", Press(ColorPrimary), MixPress},
	} {
		if strings.Contains(c.got, "var(") {
			t.Errorf("%s must contain no var() anywhere: %s", c.name, c.got)
		}
		if !strings.Contains(c.got, c.tok.Dark) {
			t.Errorf("%s must bake %s's literal intensity (%s): %s", c.name, c.tok.Name, c.tok.Dark, c.got)
		}
	}
}

// The package is build-time only: nothing here may reach a wasm binary.
func TestPackageIsBuildTimeOnly(t *testing.T) {
	cmd := exec.Command("go", "list", "-f", "{{len .GoFiles}}", ".")
	cmd.Env = append(cmd.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("css compiles for wasm (%s files); every file must be //go:build !wasm", strings.TrimSpace(string(out)))
	}
	if !strings.Contains(string(out), "build constraints exclude all Go files") {
		t.Fatalf("unexpected go list failure: %s", out)
	}
}
