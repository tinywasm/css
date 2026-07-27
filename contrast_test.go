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

func TestContrastRatios(t *testing.T) {
	pairs := []struct {
		name string
		bg   Token
		fg   Token
		min  float64
	}{
		{"SurfacePrimary", ColorPrimary, ColorOnPrimary, 4.5},
		{"SurfacePanel", ColorSurface, ColorOnSurface, 4.5},
		{"SurfaceSunken", ColorSurfaceSunken, ColorOnSurface, 4.5},
		{"SurfaceSelected", ColorSelection, ColorOnSelection, 4.5},
		{"SurfaceDanger", ColorError, ColorOnError, 4.5},
		{"SurfaceSuccess", ColorSuccess, ColorOnSuccess, 4.5},
		{"SurfaceDisabled", ColorDisabled, ColorOnDisabled, 4.5},

		// Light twins
		{"SurfacePanelLight", ColorSurfaceLight, ColorOnSurfaceLight, 4.5},
		{"SurfaceSunkenLight", ColorSurfaceSunkenLight, ColorOnSurfaceLight, 4.5},
		{"SurfaceSelectedLight", ColorSelectionLight, ColorOnSelectionLight, 4.5},
		{"SurfaceDisabledLight", ColorDisabledLight, ColorOnDisabledLight, 4.5},

		// Dark twins
		{"SurfacePanelDark", ColorSurfaceDark, ColorOnSurfaceDark, 4.5},
		{"SurfaceSunkenDark", ColorSurfaceSunkenDark, ColorOnSurfaceDark, 4.5},
		{"SurfaceSelectedDark", ColorSelectionDark, ColorOnSelectionDark, 4.5},
		{"SurfaceDisabledDark", ColorDisabledDark, ColorOnDisabledDark, 4.5},
	}

	for _, tc := range pairs {
		ratio := contrastRatio(tc.bg.Fallback, tc.fg.Fallback)
		if ratio < tc.min {
			t.Errorf("Pair %s (Bg: %s %s, Fg: %s %s) has contrast ratio %.2f:1, expected >= %.2f:1",
				tc.name, tc.bg.Name, tc.bg.Fallback, tc.fg.Name, tc.fg.Fallback, ratio, tc.min)
		} else {
			t.Logf("Pair %s (Bg: %s, Fg: %s) contrast ratio is %.2f:1 (PASS)", tc.name, tc.bg.Fallback, tc.fg.Fallback, ratio)
		}
	}
}
