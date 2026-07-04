//go:build !wasm

package css

import (
	"strings"
	"testing"
)

func TestDSL_Rule(t *testing.T) {
	sheet := NewStylesheet(
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

func TestDSL_NewAdditions(t *testing.T) {
	sheet := NewStylesheet(
		Rule(".test",
			MinWidth(Px(100)),
			MaxHeight(Vh(50)),
			AlignSelf(FlexEnd),
			Overflow(Hidden),
			Visibility(Visible),
			TextAlign(RightText),
			TextTransform(Uppercase),
			TextDecoration(None),
			TextShadow(Px(1), Px(1), Hex("#000")),
			UserSelect(None),
			TouchAction(Auto),
			ListStyleType(None),
			GridArea(Str("content")),
			GridTemplate(Calc("100% - 20px")),
			Width(Vw(80)),
			Position(Fixed),
			Top(Unset),
			Bottom(Initial),
			FlexDirection(Row),
			JustifyContent(SpaceAround),
			MarginLeft(Px(5)),
			MarginRight(Rem(0.4)),
			PaddingBottom(Space1),
			ListStyle(None),
			All(Initial),
			OverflowY(Auto),
			GridTemplateRows(Str("auto 1fr")),
			GridTemplateColumns(Str("1fr 3fr 1fr")),
			BorderRight(Vw(0.1), Str("solid"), Hex("#ccc")),
		),
		MediaDesktop(
			Rule(".desktop", Display(Grid), FlexDirection(Column)),
		),
	)
	got := sheet.String()

	// Check for a few key properties to ensure they are rendered correctly
	expectations := []string{
		"min-width: 100px;",
		"max-height: 50vh;",
		"align-self: flex-end;",
		"overflow: hidden;",
		"visibility: visible;",
		"text-align: right;",
		"text-transform: uppercase;",
		"text-decoration: none;",
		"text-shadow: 1px 1px #000;",
		"user-select: none;",
		"touch-action: auto;",
		"list-style-type: none;",
		"grid-area: content;",
		"grid-template: calc(100% - 20px);",
		"width: 80vw;",
		"position: fixed;",
		"top: unset;",
		"bottom: initial;",
		"flex-direction: row;",
		"justify-content: space-around;",
		"margin-left: 5px;",
		"margin-right: 0.4rem;",
		"padding-bottom: var(--space-1,0.25rem);",
		"list-style: none;",
		"all: initial;",
		"overflow-y: auto;",
		"grid-template-rows: auto 1fr;",
		"grid-template-columns: 1fr 3fr 1fr;",
		"border-right: 0.1vw solid #ccc;",
		"@media (orientation: landscape) and (hover: hover)",
		"flex-direction: column;",
	}

	for _, want := range expectations {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected output %q in:\n%s", want, got)
		}
	}
}

func TestDSL_Keyframes(t *testing.T) {
	sheet := NewStylesheet(
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
	sheet := NewStylesheet(
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
	sheet := NewStylesheet(
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
	sheet := NewStylesheet(
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
