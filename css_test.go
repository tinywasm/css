//go:build !wasm

package css

import (
	"strings"
	"testing"
)

func TestRootCSS_NotEmpty(t *testing.T) {
	got := RootCSS().String()
	if got == "" {
		t.Error("RootCSS() returned an empty string")
	}
}

func TestRootCSS_ContainsRootSelector(t *testing.T) {
	got := RootCSS().String()
	if !strings.Contains(got, ":root") {
		t.Errorf("RootCSS() output does not contain ':root'\nGot:\n%s", got)
	}
}

func TestRootCSS_ContainsCoreToken(t *testing.T) {
	got := RootCSS().String()
	if !strings.Contains(got, "--space-2") {
		t.Errorf("RootCSS() output does not contain core token '--space-2'\nGot:\n%s", got)
	}
}

func TestRootCSS_DoesNotContainSwitchingLogic(t *testing.T) {
	got := RootCSS().String()
	if strings.Contains(got, "@media (") {
		t.Errorf("RootCSS() must not contain @media rules (belongs in RenderCSS)\nGot:\n%s", got)
	}
}

func TestRenderCSS_ContainsColorScheme(t *testing.T) {
	got := RenderCSS().String()
	if !strings.Contains(got, "color-scheme: light dark;") {
		t.Errorf("RenderCSS() output does not contain standard light-dark color-scheme on :root\nGot:\n%s", got)
	}
}

func TestGoldenEquivalence(t *testing.T) {
	// RootCSS golden test (partial, checking key values are present)
	root := RootCSS().String()
	tokens := []string{
		"--color-primary: #1b5d8c",
		"--color-background: light-dark(#FFFFFF, #0D1117)",
		"--text-base: 1rem",
		"--space-4: 1rem",
		"--radius-md: 8px",
		"--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05)",
		"--duration-base: 250ms",
		"--z-modal: 300",
		"--bp-md: 768px",
	}
	for _, tok := range tokens {
		if !strings.Contains(root, tok) {
			t.Errorf("RootCSS missing expected token: %s", tok)
		}
	}

	// RenderCSS golden test
	render := RenderCSS().String()
	rules := []string{
		"box-sizing: border-box",
		"margin: 0",
		"font-size: var(--text-base",
		"outline: 2px solid var(--color-primary",
		"display: block",
		"color-scheme: light dark",
		"color-scheme: light",
		"color-scheme: dark",
	}
	for _, rule := range rules {
		if !strings.Contains(render, rule) {
			t.Errorf("RenderCSS missing expected rule: %s", rule)
		}
	}
}

func TestTheme_NoOverrides(t *testing.T) {
	got := Theme().String()
	want := RootCSS().String()
	if got != want {
		t.Error("Theme() without overrides must be identical to RootCSS()")
	}
}

func TestTheme_ContainsFullCatalog(t *testing.T) {
	got := Theme(Set(ColorSurface, "#FF0000")).String()
	tokens := []string{
		"--space-2",
		"--radius-md",
		"--text-xl",
		"--color-surface",
	}
	for _, tok := range tokens {
		if !strings.Contains(got, tok) {
			t.Errorf("Theme() output does not contain catalog token '%s'", tok)
		}
	}
}

func TestTheme_OverrideTakesPrecedence(t *testing.T) {
	overrideVal := "#3f88bf"
	got := Theme(Set(ColorSurface, overrideVal)).String()

	// Find all occurrences of --color-surface
	tokenName := "--color-surface"
	indices := []int{}
	lastIdx := 0
	for {
		idx := strings.Index(got[lastIdx:], tokenName)
		if idx == -1 {
			break
		}
		indices = append(indices, lastIdx+idx)
		lastIdx += idx + len(tokenName)
	}

	if len(indices) < 2 {
		t.Fatalf("Expected at least 2 occurrences of %s (default + override), got %d", tokenName, len(indices))
	}

	// The last occurrence should be the override
	lastOccurrence := got[indices[len(indices)-1]:]
	if !strings.Contains(lastOccurrence, overrideVal) {
		t.Errorf("Last occurrence of %s does not contain override value %s\nContext: %s", tokenName, overrideVal, lastOccurrence)
	}

	// Verify the default value is also present (as Theme appends overrides)
	if !strings.Contains(got, "light-dark(#F2F2F7, #161B22)") {
		t.Errorf("Default value for %s (light-dark(#F2F2F7, #161B22)) missing from output", tokenName)
	}
}

func TestNoUndeclaredTokensInEmittedCSS(t *testing.T) {
	// Gather all known token names from the catalog
	knownTokens := map[string]bool{}
	allTokens := []ValueGetter{
		ColorPrimary, ColorOnPrimary, ColorSuccess, ColorOnSuccess, ColorDanger, ColorOnDanger,
		ColorBackground, ColorOnBackground, ColorSurface, ColorOnSurface, ColorOutline, ColorMuted,
		TextXs, TextSm, TextBase, TextLg, TextXl, Text2xl,
		LeadingNormal, FontWeightRegular, FontWeightMedium, FontWeightBold,
		Space0, Space1, Space2, Space3, Space4, Space6, Space8, Space12,
		RadiusSm, RadiusMd, RadiusLg, RadiusFull,
		ShadowSm, ShadowMd, ShadowLg,
		DurationFast, DurationBase, DurationSlow, EaseInOut,
		ZBase, ZDropdown, ZSticky, ZModal, ZToast, ZTooltip,
		BpSm, BpMd, BpLg, BpXl,
		MaxWReadable,
		ColumnNarrow, ColumnMedium, ColumnWide,
	}
	for _, tok := range allTokens {
		knownTokens[tok.GetName()] = true
	}

	cssStr := RootCSS().String() + "\n" + RenderCSS().String()

	idx := 0
	for {
		start := strings.Index(cssStr[idx:], "--")
		if start == -1 {
			break
		}
		startIdx := idx + start
		endIdx := startIdx
		for endIdx < len(cssStr) {
			c := cssStr[endIdx]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
				endIdx++
			} else {
				break
			}
		}
		propName := cssStr[startIdx:endIdx]
		if !knownTokens[propName] {
			t.Errorf("Undeclared token found in emitted CSS: %s", propName)
		}
		idx = endIdx
	}
}
