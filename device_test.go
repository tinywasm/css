//go:build !wasm

package css

import (
	"strconv"
	"strings"
	"testing"
)

func parseMaxWidth(query string) (float64, bool) {
	_, after, ok := strings.Cut(query, "max-width: ")
	if !ok {
		return 0, false
	}
	end := strings.Index(after, "px")
	if end == -1 {
		return 0, false
	}
	v, err := strconv.ParseFloat(after[:end], 64)
	return v, err == nil
}

func parseMinWidth(query string) (float64, bool) {
	_, after, ok := strings.Cut(query, "min-width: ")
	if !ok {
		return 0, false
	}
	end := strings.Index(after, "px")
	if end == -1 {
		return 0, false
	}
	v, err := strconv.ParseFloat(after[:end], 64)
	return v, err == nil
}

func parsePx(s string) (float64, bool) {
	trimmed, ok := strings.CutSuffix(s, "px")
	if !ok {
		return 0, false
	}
	v, err := strconv.ParseFloat(trimmed, 64)
	return v, err == nil
}

// The thresholds exist in two places: as literals inside Device.Query(), where a
// var() is invalid (SPECS §6), and as the --bp-* token values. Nothing in the
// language ties them together, so a change to one would drift silently past every
// consumer. This test is that tie: it is the reason the duplication is safe.
func TestDeviceThresholdsMatchBreakpointTokens(t *testing.T) {
	sm, ok := parsePx(BpSm.Dark)
	if !ok {
		t.Fatalf("BpSm value %q is not a px length", BpSm.Dark)
	}
	lg, ok := parsePx(BpLg.Dark)
	if !ok {
		t.Fatalf("BpLg value %q is not a px length", BpLg.Dark)
	}

	// The gap a fractional viewport width could fall into. Kept in sync with the
	// ".98" fractions in device.go.
	const epsilon = 0.02

	cases := []struct {
		name  string
		got   func() (float64, bool)
		want  float64
		token string
	}{
		{"Tablet min-width", func() (float64, bool) { return parseMinWidth(Tablet.Query()) }, sm, "BpSm"},
		{"Mobile max-width", func() (float64, bool) { return parseMaxWidth(Mobile.Query()) }, sm - epsilon, "BpSm"},
		{"Desktop min-width", func() (float64, bool) { return parseMinWidth(Desktop.Query()) }, lg, "BpLg"},
		{"Tablet max-width", func() (float64, bool) { return parseMaxWidth(Tablet.Query()) }, lg - epsilon, "BpLg"},
	}

	for _, c := range cases {
		got, ok := c.got()
		if !ok {
			t.Errorf("%s: could not be parsed out of the query", c.name)
			continue
		}
		if got != c.want {
			t.Errorf("%s = %g, want %g (derived from %s = %q); device.go and catalog.go have drifted apart",
				c.name, got, c.want, c.token, c.token)
		}
	}
}

func TestDeviceClassesPartition(t *testing.T) {
	tests := []struct {
		width    float64
		expected Device
	}{
		{320, Mobile},
		{639, Mobile},
		// Fractional widths are the whole reason for the ".98" fractions: at a
		// plain 639px/1023px boundary these fall into a gap matching no class.
		{639.5, Mobile},
		{639.98, Mobile},
		{640, Tablet},
		{767, Tablet},
		{1023, Tablet},
		{1023.5, Tablet},
		{1023.98, Tablet},
		{1024, Desktop},
		{1440, Desktop},
	}

	for _, tt := range tests {
		var matching Device
		var matched bool

		for _, d := range []Device{Mobile, Tablet, Desktop} {
			q := d.Query()
			matches := true

			if strings.Contains(q, "max-width: ") {
				mw, ok := parseMaxWidth(q)
				if ok && tt.width > mw {
					matches = false
				}
			}
			if strings.Contains(q, "min-width: ") {
				mw, ok := parseMinWidth(q)
				if ok && tt.width < mw {
					matches = false
				}
			}

			if matches {
				if matched {
					t.Errorf("width %.0f matches both %s and %s", tt.width, matching, d)
				}
				matching = d
				matched = true
			}
		}

		if !matched {
			t.Errorf("width %.0f matches no device class", tt.width)
		}
		if matched && matching != tt.expected {
			t.Errorf("width %.0f matched %s, expected %s", tt.width, matching, tt.expected)
		}
	}
}

func TestDeviceQueryHasNoVar(t *testing.T) {
	for _, d := range []Device{Mobile, Tablet, Desktop} {
		if strings.Contains(d.Query(), "var(") {
			t.Errorf("%s.Query() contains var(): %s", d, d.Query())
		}
	}
}

func TestQueryJoinIsDeterministic(t *testing.T) {
	a := Query(Desktop, Mobile)
	b := Query(Mobile, Desktop)
	if a != b {
		t.Errorf("Query(Desktop, Mobile) != Query(Mobile, Desktop): %q vs %q", a, b)
	}
	want := Mobile.Query() + ", " + Desktop.Query()
	if a != want {
		t.Errorf("Query(Mobile, Desktop) = %q, want %q", a, want)
	}
}

func TestQueryDeduplicates(t *testing.T) {
	single := Query(Mobile)
	dup := Query(Mobile, Mobile)
	if single != dup {
		t.Errorf("Query(Mobile, Mobile) = %q, want %q", dup, single)
	}
}

func TestQueryEmpty(t *testing.T) {
	if got := Query(); got != "" {
		t.Errorf("Query() = %q, want empty", got)
	}
}

func TestRailTokensDeclared(t *testing.T) {
	css := RootCSS().String()
	for _, name := range []string{"--rail-narrow", "--rail-wide"} {
		if !strings.Contains(css, name+":") {
			t.Errorf("RootCSS() missing declaration for %s", name)
		}
	}
}

func TestDeviceStringRoundTrip(t *testing.T) {
	for _, d := range []Device{Mobile, Tablet, Desktop} {
		if d.String() == "Unknown" {
			t.Errorf("%d.String() = Unknown", d)
		}
	}
}

func TestQueryIgnoresUnknown(t *testing.T) {
	got := Query(Mobile, 99, Desktop)
	want := Query(Mobile, Desktop)
	if got != want {
		t.Errorf("Query with unknown = %q, want %q", got, want)
	}
}

func TestQueryAllReturnsThree(t *testing.T) {
	got := Query(Mobile, Tablet, Desktop)
	parts := strings.Split(got, ", ")
	if len(parts) != 3 {
		t.Errorf("Query(all) has %d parts, want 3: %q", len(parts), got)
	}
}
