//go:build !wasm

package css

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func parseHex(hex string) (r, g, b float64) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return 0, 0, 0
	}
	ri, _ := strconv.ParseUint(hex[0:2], 16, 32)
	gi, _ := strconv.ParseUint(hex[2:4], 16, 32)
	bi, _ := strconv.ParseUint(hex[4:6], 16, 32)
	return float64(ri) / 255.0, float64(gi) / 255.0, float64(bi) / 255.0
}

func relativeLuminance(hex string) float64 {
	r, g, b := parseHex(hex)
	c := func(v float64) float64 {
		if v <= 0.03928 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*c(r) + 0.7152*c(g) + 0.0722*c(b)
}

func contrastRatio(hex1, hex2 string) float64 {
	l1 := relativeLuminance(hex1)
	l2 := relativeLuminance(hex2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

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
		ratioLight := contrastRatio(bgLight, fgLight)
		if ratioLight < tc.Min {
			t.Errorf("Pair %s (Light Mode) (Bg: %s %s, Fg: %s %s) has contrast ratio %.2f:1, expected >= %.2f:1",
				tc.Name, tc.Bg.GetName(), bgLight, tc.Fg.GetName(), fgLight, ratioLight, tc.Min)
		} else {
			t.Logf("Pair %s (Light Mode) (Bg: %s, Fg: %s) contrast ratio is %.2f:1 (PASS)", tc.Name, bgLight, fgLight, ratioLight)
		}

		// Test Dark mode contrast
		bgDark := resolveColor(tc.Bg.GetFallback(), true)
		fgDark := resolveColor(tc.Fg.GetFallback(), true)
		ratioDark := contrastRatio(bgDark, fgDark)
		if ratioDark < tc.Min {
			t.Errorf("Pair %s (Dark Mode) (Bg: %s %s, Fg: %s %s) has contrast ratio %.2f:1, expected >= %.2f:1",
				tc.Name, tc.Bg.GetName(), bgDark, tc.Fg.GetName(), fgDark, ratioDark, tc.Min)
		} else {
			t.Logf("Pair %s (Dark Mode) (Bg: %s, Fg: %s) contrast ratio is %.2f:1 (PASS)", tc.Name, bgDark, fgDark, ratioDark)
		}
	}
}
