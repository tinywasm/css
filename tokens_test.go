//go:build !wasm

package css

import (
	"os/exec"
	"strings"
	"testing"
)

var composedTokens = []Token{ColorSurfaceSunken, ColorSelection, ColorOnSelection}

// A fallback that composes other variables must not re-state their values.
func TestComposedTokensContainNoHex(t *testing.T) {
	for _, tok := range composedTokens {
		if strings.Contains(tok.GetFallback(), "#") {
			t.Errorf("%s hardcodes a hex inside a composed fallback: %s", tok.Name, tok.GetFallback())
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

// The intensity must stay a live var() so an app can retune it with
// Theme(Set(MixHover, ...)) — a baked percentage is only changeable by
// republishing this package and regenerating every consumer's stylesheet.
func TestInteractionIntensityIsOverridable(t *testing.T) {
	for _, c := range []struct {
		name string
		got  string
		tok  Token
	}{
		{"Hover", Hover(ColorPrimary), MixHover},
		{"Focus", Focus(ColorPrimary), MixFocus},
		{"Press", Press(ColorPrimary), MixPress},
	} {
		if !strings.Contains(c.got, c.tok.Var()) {
			t.Errorf("%s bakes its intensity instead of referencing %s: %s", c.name, c.tok.Name, c.got)
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
