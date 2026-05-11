//go:build !wasm

package css

import (
	"strings"
	"testing"
)

func TestDSL_Rule(t *testing.T) {
	sheet := New(
		Rule(".btn",
			BackgroundColor(Hex("#fff")),
			Color(ColorPrimary),
		),
	)
	got := sheet.String()
	want := ".btn {\n  background-color: #fff;\n  color: var(--color-primary,#00ADD8);\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Keyframes(t *testing.T) {
	sheet := New(
		Keyframes("pulse",
			At("0%",
				Transform(Str("scale(1)")),
				Opacity(1),
			),
			At("100%",
				Transform(Str("scale(1.1)")),
				Opacity(0),
				Color(ColorPrimary),
			),
		),
	)
	got := sheet.String()
	want := "@keyframes pulse {\n  0% {\n    transform: scale(1);\n    opacity: 1;\n  }\n  100% {\n    transform: scale(1.1);\n    opacity: 0;\n    color: var(--color-primary,#00ADD8);\n  }\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Root(t *testing.T) {
	sheet := New(
		Root(
			Declare(ColorPrimary, "#00ADD8"),
			Bind(ColorBackground, ColorBackgroundLight),
		),
	)
	got := sheet.String()
	want := ":root {\n  --color-primary: #00ADD8;\n  --color-background: var(--color-background-light,#FFFFFF);\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Media(t *testing.T) {
	sheet := New(
		MediaPrefersDark(
			Root(Bind(ColorBackground, ColorBackgroundDark)),
		),
	)
	got := sheet.String()
	if !strings.Contains(got, "@media (prefers-color-scheme: dark)") {
		t.Errorf("missing media query: %s", got)
	}
	if !strings.Contains(got, "--color-background: var(--color-background-dark") {
		t.Errorf("missing binding: %s", got)
	}
}

func TestDSL_Pseudo(t *testing.T) {
	btn := Class("btn")
	sheet := New(
		Rule(btn.Hover(),
			Opacity(0.8),
		),
	)
	got := sheet.String()
	want := ".btn:hover {\n  opacity: 0.8;\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
