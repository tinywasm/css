//go:build !wasm

package css

import (
	"strings"
	"testing"

	twcolor "github.com/tinywasm/color"
)

// resolveColor resolves a color fallback value which might be a CSS light-dark() function.
// For example, resolveColor("light-dark(#FFFFFF, #0D1117)", true) returns "#0D1117".
func resolveColor(val string, dark bool) string {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "light-dark(") && strings.HasSuffix(val, ")") {
		inner := val[len("light-dark(") : len(val)-1]
		parts := strings.Split(inner, ",")
		if len(parts) == 2 {
			if dark {
				return strings.TrimSpace(parts[1])
			}
			return strings.TrimSpace(parts[0])
		}
	}
	return val
}

func TestContrastRatios(t *testing.T) {
	for _, tc := range AllPairs() {
		// Test Light mode contrast
		bgLight := resolveColor(tc.Bg.GetFallback(), false)
		fgLight := resolveColor(tc.Fg.GetFallback(), false)
		ratioLight := twcolor.Contrast(twcolor.Color(bgLight), twcolor.Color(fgLight))
		if ratioLight < tc.Min {
			t.Errorf("Pair %s (Light Mode) (Bg: %s %s, Fg: %s %s) has contrast ratio %.2f:1, expected >= %.2f:1",
				tc.Name, tc.Bg.GetName(), bgLight, tc.Fg.GetName(), fgLight, ratioLight, tc.Min)
		} else {
			t.Logf("Pair %s (Light Mode) (Bg: %s, Fg: %s) contrast ratio is %.2f:1 (PASS)", tc.Name, bgLight, fgLight, ratioLight)
		}

		// Test Dark mode contrast
		bgDark := resolveColor(tc.Bg.GetFallback(), true)
		fgDark := resolveColor(tc.Fg.GetFallback(), true)
		ratioDark := twcolor.Contrast(twcolor.Color(bgDark), twcolor.Color(fgDark))
		if ratioDark < tc.Min {
			t.Errorf("Pair %s (Dark Mode) (Bg: %s %s, Fg: %s %s) has contrast ratio %.2f:1, expected >= %.2f:1",
				tc.Name, tc.Bg.GetName(), bgDark, tc.Fg.GetName(), fgDark, ratioDark, tc.Min)
		} else {
			t.Logf("Pair %s (Dark Mode) (Bg: %s, Fg: %s) contrast ratio is %.2f:1 (PASS)", tc.Name, bgDark, fgDark, ratioDark)
		}
	}
}
