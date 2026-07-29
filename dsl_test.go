//go:build !wasm

package css

import (
	"strings"
	"testing"

	"github.com/tinywasm/fmt"
)

func TestDSL_Rule(t *testing.T) {
	sheet := NewStylesheet(
		rule(selector(".btn"),
			backgroundColor(hex("#fff")),
			color(ColorPrimary),
		),
	)
	got := sheet.String()
	want := ".btn {\n  background-color: #fff;\n  color: var(--color-primary,#1b5d8c);\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Padding_JoinValues(t *testing.T) {
	sheet := NewStylesheet(
		rule(selector(".test"),
			padding(Space1, Space2, Space3, Space4),
		),
	)
	got := sheet.String()
	want := "  padding: var(--space-1,0.25rem) var(--space-2,0.5rem) var(--space-3,0.75rem) var(--space-4,1rem);\n"
	if !strings.Contains(got, want) {
		t.Errorf("got:\n%s\nexpected it to contain:\n%s", got, want)
	}
}

func TestDSL_RawRule_Semicolon(t *testing.T) {
	sheet := NewStylesheet(
		rule(selector(".test"),
			rawRule("grid-template: 1fr"),
			rawRule("gap: 1rem"),
			rawRule("  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;"),
		),
	)
	got := sheet.String()
	want1 := "grid-template: 1fr;\n"
	want2 := "gap: 1rem;\n"
	want3 := "  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;\n"
	if !strings.Contains(got, want1) || !strings.Contains(got, want2) || !strings.Contains(got, want3) {
		t.Errorf("got:\n%s\nmissing expected raw rules with semicolon", got)
	}
}

func TestDSL_NewAdditions(t *testing.T) {
	sheet := NewStylesheet(
		rule(selector(".test"),
			minWidth(px(100)),
			maxHeight(vh(50)),
			alignSelf(flexEnd),
			overflow(hidden),
			visibility(visible),
			textAlign(rightText),
			textTransform(uppercase),
			textDecoration(none),
			textShadow(px(1), px(1), hex("#000")),
			userSelect(none),
			touchAction(auto),
			listStyleType(none),
			gridArea(str("content")),
			gridTemplate(calc("100% - 20px")),
			width(vw(80)),
			position(fixed),
			top(unset),
			bottom(initial),
			flexDirection(row),
			justifyContent(spaceAround),
			marginLeft(px(5)),
			marginRight(rem(0.4)),
			paddingBottom(Space1),
			listStyle(none),
			all(initial),
			overflowY(auto),
			gridTemplateRows(str("auto 1fr")),
			gridTemplateColumns(str("1fr 3fr 1fr")),
			borderRight(vw(0.1), str("solid"), hex("#ccc")),
			paddingTop(px(10)),
			paddingLeft(px(15)),
			paddingRight(px(20)),
			marginTop(rem(0.5)),
			marginBottom(rem(0.8)),
			flexWrap(wrap),
			flexGrow(px(1)),
			alignContent(spaceBetween),
			borderBottom(px(2), str("solid"), hex("#000")),
			borderLeft(px(1), str("dashed"), hex("#999")),
			backgroundSize(str("cover")),
			backgroundPosition(center),
			backgroundRepeat(noRepeat),
			position(relative),
		),
		mediaDesktop(
			rule(selector(".desktop"), display(grid), flexDirection(column)),
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
		"position: relative;",
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
		"padding-top: 10px;",
		"padding-left: 15px;",
		"padding-right: 20px;",
		"margin-top: 0.5rem;",
		"margin-bottom: 0.8rem;",
		"flex-wrap: wrap;",
		"flex-grow: 1px;",
		"align-content: space-between;",
		"border-bottom: 2px solid #000;",
		"border-left: 1px dashed #999;",
		"background-size: cover;",
		"background-position: center;",
		"background-repeat: no-repeat;",
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
		keyframes("pulse",
			at("0%",
				transform(str("scale(1)")),
				opacity(1),
			),
			at("100%",
				transform(str("scale(1.1)")),
				opacity(0),
				color(ColorPrimary),
			),
		),
	)
	got := sheet.String()
	want := "@keyframes pulse {\n  0% {\n    transform: scale(1);\n    opacity: 1;\n  }\n  100% {\n    transform: scale(1.1);\n    opacity: 0;\n    color: var(--color-primary,#1b5d8c);\n  }\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Root(t *testing.T) {
	sheet := NewStylesheet(
		root(
			declare(ColorPrimary),
			declare(ColorBackground),
		),
	)
	got := sheet.String()
	want := ":root {\n  --color-primary: #1b5d8c;\n  --color-background: light-dark(#FFFFFF, #0D1117);\n}\n\n"
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDSL_Media(t *testing.T) {
	sheet := NewStylesheet(
		mediaPrefersDark(
			rule(selector("body"), background(ColorBackground)),
		),
	)
	got := sheet.String()
	if !strings.Contains(got, "@media (prefers-color-scheme: dark)") {
		t.Errorf("missing media query: %s", got)
	}
	if !strings.Contains(got, "background: var(--color-background,light-dark(#FFFFFF, #0D1117))") {
		t.Errorf("missing background: %s", got)
	}
}

// Stylesheet.String() must not release its Conv twice: Conv.String() already
// returns the object to the pool, so an extra PutConv leaves the same pointer
// in the pool twice and two later Convert() calls share one buffer.
func TestStylesheetStringDoesNotPoisonConvPool(t *testing.T) {
	_ = NewStylesheet(Raw(".x{color:red}")).String()

	outer := fmt.Convert()
	outer.WriteString("outer-")
	inner := fmt.Convert()
	inner.WriteString("inner")

	if got := inner.String(); got != "inner" {
		t.Fatalf("inner Conv shares its buffer with the outer one: got %q, want %q", got, "inner")
	}
	if got := outer.String(); got != "outer-" {
		t.Fatalf("outer Conv was clobbered: got %q, want %q", got, "outer-")
	}
}

// joinValues has the same contract, and the whole sheet must render identically
// however many times it is stringified.
func TestStylesheetStringIsRepeatable(t *testing.T) {
	sheet := NewStylesheet(
		rule(selector(".test"),
			padding(Space1, Space2),
			border(px(1), ColorOutline),
		),
	)

	first := sheet.String()
	for i := 0; i < 3; i++ {
		// Interleave an unrelated Conv so a poisoned pool would surface here.
		scratch := fmt.Convert()
		scratch.WriteString("scratch")
		if got := scratch.String(); got != "scratch" {
			t.Fatalf("scratch Conv contaminated on pass %d: got %q", i, got)
		}
		if got := sheet.String(); got != first {
			t.Fatalf("pass %d differs:\ngot:\n%q\nfirst:\n%q", i, got, first)
		}
	}
}
