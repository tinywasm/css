//go:build !wasm

package css

import (
	"strings"
	"testing"

	"github.com/tinywasm/font"
)

func TestFontFaces_FourFaces(t *testing.T) {
	got := FontFaces(font.Declare("Roboto", "x"), "/assets").String()
	want := `@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: url("/assets/Roboto-Regular.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: normal;
  font-weight: 700;
  font-display: swap;
  src: url("/assets/Roboto-Bold.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: italic;
  font-weight: 400;
  font-display: swap;
  src: url("/assets/Roboto-Italic.ttf") format("truetype");
}

@font-face {
  font-family: "Roboto";
  font-style: italic;
  font-weight: 700;
  font-display: swap;
  src: url("/assets/Roboto-BoldItalic.ttf") format("truetype");
}

`
	if got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestFontFaces_WeightAndStylePerFace(t *testing.T) {
	got := FontFaces(font.Declare("Roboto", ""), "/assets").String()
	cases := []struct {
		face   string
		weight string
		style  string
	}{
		{"Roboto-Regular", "400", "normal"},
		{"Roboto-Bold", "700", "normal"},
		{"Roboto-Italic", "400", "italic"},
		{"Roboto-BoldItalic", "700", "italic"},
	}
	for _, c := range cases {
		idx := strings.Index(got, c.face)
		if idx < 0 {
			t.Errorf("missing face file %s", c.face)
			continue
		}
		// Look back to the start of this @font-face block.
		blockStart := strings.LastIndex(got[:idx], "@font-face")
		blockEnd := strings.Index(got[idx:], "}")
		block := got[blockStart : idx+blockEnd]
		if !strings.Contains(block, "font-weight: "+c.weight) {
			t.Errorf("%s: want font-weight %s in:\n%s", c.face, c.weight, block)
		}
		if !strings.Contains(block, "font-style: "+c.style) {
			t.Errorf("%s: want font-style %s in:\n%s", c.face, c.style, block)
		}
		if !strings.Contains(block, `format("truetype")`) {
			t.Errorf("%s: missing format(\"truetype\")", c.face)
		}
		if !strings.Contains(block, "font-display: swap") {
			t.Errorf("%s: missing font-display: swap", c.face)
		}
	}
}

func TestFontFaces_URLPrefixJoining(t *testing.T) {
	for _, c := range []struct {
		prefix string
		want   string
	}{
		{"/assets", `url("/assets/Roboto-Regular.ttf")`},
		{"/assets/", `url("/assets/Roboto-Regular.ttf")`},
		{"", `url("Roboto-Regular.ttf")`},
	} {
		got := FontFaces(font.Declare("Roboto", ""), c.prefix).String()
		if !strings.Contains(got, c.want) {
			t.Errorf("prefix %q: want %s in:\n%s", c.prefix, c.want, got)
		}
		if strings.Contains(got, "//Roboto") {
			t.Errorf("prefix %q produced double slash: %s", c.prefix, got)
		}
	}
}

func TestFontFaces_EmptyFamilyEmitsNothing(t *testing.T) {
	got := FontFaces(font.Declare("", "x"), "/assets").String()
	if got != "" {
		t.Errorf("empty family must emit no rules, got:\n%q", got)
	}
	if strings.Contains(got, "@font-face") {
		t.Errorf("empty family must not emit @font-face")
	}
}

func TestFontFaces_NotInRootOrRender(t *testing.T) {
	if strings.Contains(RootCSS().String(), "@font-face") {
		t.Error("RootCSS must not contain @font-face")
	}
	if strings.Contains(RenderCSS().String(), "@font-face") {
		t.Error("RenderCSS must not contain @font-face")
	}
}
