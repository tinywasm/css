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

func TestRenderCSS_SitsInTheLowestLayer(t *testing.T) {
	// An unlayered rule outranks every layered one, so an unlayered
	// `svg { display: block }` would beat a component's
	// `@layer widgets { .part { display: none } }` and make any icon hidden
	// by state impossible to hide.
	got := RenderCSS().String()
	if !strings.HasPrefix(strings.TrimSpace(got), "@layer tokens {") {
		t.Errorf("the base reset must be wrapped in @layer tokens\nGot:\n%s", got)
	}
	if !strings.Contains(got, "img, svg, video") {
		t.Errorf("expected the replaced-element reset to survive the layer wrap\nGot:\n%s", got)
	}
}

// The reset's job is that a part's own rules land on the same box in every
// engine. Each case below is a place where two shipping browsers disagree
// unless the reset says otherwise, so each is a guard, not coverage.
func TestRenderCSS_NormalizesCrossBrowserDefaults(t *testing.T) {
	got := RenderCSS().String()
	for _, c := range []struct{ rule, why string }{
		{"-webkit-tap-highlight-color: transparent",
			"iOS paints a grey wash over anything tapped; Chrome Android uses another colour"},
		{"appearance: none",
			"iOS renders <button> as push-button chrome that outranks a part's background and radius"},
		{"background-image: none",
			"iOS gives <button> a vertical gradient no part asked for"},
		{"text-transform: none",
			"Firefox and Edge let <select> inherit text-transform; other engines do not"},
		{"opacity: 1",
			"Firefox ships ::placeholder at opacity 0.54"},
		{"font-size: 1em",
			"every engine renders the monospace default about 3px smaller"},
	} {
		if !strings.Contains(got, c.rule) {
			t.Errorf("reset missing %q — %s", c.rule, c.why)
		}
	}
}

// appearance: none flattens a text field's chrome but leaves the UA's own
// `border: 2px inset` standing, so every input arrived carrying a heavy dark
// box. A skin cannot undo that by painting a flat fill — there is no border
// declaration to override, only chrome to remove — so the reset has to. This
// asserts it on the text-field rule specifically: `border: 0` also appears in
// the button rule above, and a whole-sheet Contains would pass on that alone.
func TestRenderCSS_StripsTheUABorderFromTextFields(t *testing.T) {
	got := RenderCSS().String()
	i := strings.Index(got, `input:where(:not([type="checkbox"]):not([type="radio"]))`)
	if i == -1 {
		t.Fatal("expected a text-field reset rule")
	}
	end := strings.Index(got[i:], "}")
	if end == -1 {
		t.Fatal("malformed text-field reset rule")
	}
	if block := got[i : i+end]; !strings.Contains(block, "border: 0") {
		t.Errorf("the text-field reset must drop the UA border, block:\n%s", block)
	}
}

// appearance: none erases a checkbox and a radio instead of flattening them:
// the control disappears rather than losing its chrome. The text-field rule
// must keep both out.
func TestRenderCSS_LeavesCheckboxAndRadioNative(t *testing.T) {
	got := RenderCSS().String()
	if !strings.Contains(got, `input:where(:not([type="checkbox"]):not([type="radio"]))`) {
		t.Errorf("the text-field appearance reset must exclude checkbox and radio\nGot:\n%s", got)
	}
}

// Author styles outrank the UA stylesheet whatever their layer, so
// `img, svg, video { display: block }` defeats the UA's own [hidden] rule
// unless the reset restates it.
func TestRenderCSS_KeepsHiddenAttributeWorking(t *testing.T) {
	got := RenderCSS().String()
	img := strings.Index(got, "img, svg, video")
	hidden := strings.Index(got, "[hidden]")
	if hidden == -1 {
		t.Fatalf("reset must restate [hidden] { display: none }\nGot:\n%s", got)
	}
	if img == -1 {
		t.Fatalf("expected the replaced-element rule to be present\nGot:\n%s", got)
	}
	if hidden < img {
		t.Errorf("[hidden] must come after the img/svg/video rule it defends against")
	}
}

func TestGoldenEquivalence(t *testing.T) {
	// RootCSS golden test (partial, checking key values are present)
	root := RootCSS().String()
	tokens := []string{
		"--color-primary: #654FF0",
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

func TestSetGradient_OnlyTouchesImageVar(t *testing.T) {
	got := Theme(SetGradient(ColorPrimary, "135deg", ColorPrimary, ColorAccent)).String()

	if strings.Contains(got, "--color-primary: ;") || strings.Contains(got, "--color-primary:;") {
		t.Errorf("SetGradient alone must not emit an empty value for the token itself, got: %s", got)
	}
	if !strings.Contains(got, "--color-primary-image: linear-gradient(135deg,") {
		t.Errorf("SetGradient must emit the token's ImageVarName() declaration, got: %s", got)
	}
}

// TestSetGradient_ReferencesLiveOverride guards against a real bug caught
// while building this: from/to must resolve through the CASCADE (var()), not
// bake in the Go struct's catalog default — SetGradient(ColorPrimary, ...,
// ColorPrimary, ColorAccent) combined with Set(ColorPrimary, "#16a34a") in
// the SAME Theme() call must produce a gradient that follows the app's own
// override, not ColorPrimary's catalog default (a stale color the moment an
// app overrides it would be a silent, hard-to-notice bug — the gradient
// direction/stops would look right in the source but wrong on screen).
func TestSetGradient_ReferencesLiveOverride(t *testing.T) {
	got := Theme(
		Set(ColorPrimary, "#16a34a"),
		SetGradient(ColorPrimary, "135deg", ColorPrimary, ColorAccent),
	).String()

	if !strings.Contains(got, "linear-gradient(135deg, var(--color-primary") {
		t.Errorf("gradient must reference --color-primary through var(), not a baked-in static value, got: %s", got)
	}
	if strings.Contains(got, "linear-gradient(135deg, #654FF0") {
		t.Errorf("gradient baked in ColorPrimary's catalog default instead of referencing the live custom property, got: %s", got)
	}
}

func TestNoUndeclaredTokensInEmittedCSS(t *testing.T) {
	// Gather all known token names from the catalog
	knownTokens := map[string]bool{}
	allTokens := []ValueGetter{
		ColorPrimary, ColorOnPrimary, ColorSuccess, ColorOnSuccess, ColorDanger, ColorOnDanger,
		ColorAccent, ColorOnAccent,
		ColorBackground, ColorOnBackground, ColorSurface, ColorOnSurface, ColorOutline, ColorMuted,
		ColorSurfaceSunken, ColorSelection, ColorOnSelection,
		ColorAccentWash, ColorAccentHover,
		MixHover, MixFocus, MixPress,
		FontSans, TextXs, TextSm, TextBase, TextLg, TextXl, Text2xl,
		LeadingNormal, FontWeightRegular, FontWeightMedium, FontWeightBold,
		Space0, Space1, Space2, Space3, Space4, Space6, Space8, Space12,
		RadiusSm, RadiusMd, RadiusLg, RadiusFull,
		ShadowSm, ShadowMd, ShadowLg,
		DurationFast, DurationBase, DurationSlow, EaseInOut,
		ZBase, ZDropdown, ZSticky, ZModal, ZToast, ZTooltip,
		BpSm, BpMd, BpLg, BpXl,
		MaxWReadable,
		ColumnNarrow, ColumnMedium, ColumnWide,
		RailNarrow, RailWide,
		ControlHeight, ChipWidth, ChipHeight, VeilBlur,
		SafeTop, SafeRight, SafeBottom, SafeLeft, ViewportH,
	}
	for _, tok := range allTokens {
		knownTokens[tok.GetName()] = true
	}
	// Every theme pair also declares plain light/dark half properties (see
	// declareSplit, Token.EnhancedVar) — legitimate, not drift.
	for _, t := range []Token{ColorBackground, ColorOnBackground, ColorSurface, ColorOnSurface, ColorOutline, ColorMuted} {
		knownTokens[t.LightVarName()] = true
		knownTokens[t.DarkVarName()] = true
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

func TestRootCSS_ContainsDeviceGeometryTokens(t *testing.T) {
	got := RootCSS().String()
	for _, decl := range []string{
		"--safe-top: env(safe-area-inset-top, 0px)",
		"--safe-right: env(safe-area-inset-right, 0px)",
		"--safe-bottom: env(safe-area-inset-bottom, 0px)",
		"--safe-left: env(safe-area-inset-left, 0px)",
		"--viewport-h: 100dvh",
	} {
		if !strings.Contains(got, decl) {
			t.Errorf("RootCSS missing device geometry declaration %q\nGot:\n%s", decl, got)
		}
	}
}

// An env() without a fallback invalidates the whole declaration in any browser
// that does not know the variable — a silent failure. Every emitted env() must
// carry its second argument.
func TestRootCSS_EveryEnvHasFallback(t *testing.T) {
	got := RootCSS().String()
	idx := 0
	for {
		start := strings.Index(got[idx:], "env(")
		if start == -1 {
			break
		}
		startIdx := idx + start
		// Find matching close paren for this env( call.
		depth := 0
		endIdx := -1
		for i := startIdx + 3; i < len(got); i++ {
			switch got[i] {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					endIdx = i
				}
			}
			if endIdx != -1 {
				break
			}
		}
		if endIdx == -1 {
			t.Fatalf("unclosed env( starting at %d in RootCSS output", startIdx)
		}
		call := got[startIdx : endIdx+1]
		// env(name, fallback) must contain a comma separating name from fallback.
		inner := call[len("env(") : len(call)-1]
		if !strings.Contains(inner, ",") {
			t.Errorf("env() without fallback (silent invalidation risk): %s", call)
		}
		idx = endIdx + 1
	}
}

func TestViewportH_Var(t *testing.T) {
	got := ViewportH.Var()
	if got != "var(--viewport-h,100dvh)" {
		t.Errorf("ViewportH.Var() = %q, want %q", got, "var(--viewport-h,100dvh)")
	}
}
