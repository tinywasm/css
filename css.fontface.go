//go:build !wasm

package css

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/font"
)

// FontFaces returns the @font-face block for the four faces of the declared
// family, served from urlPrefix (e.g. "/assets"). The prefix is caller data:
// this package neither knows nor invents server paths.
//
// Not part of RootCSS() or RenderCSS(): whoever serves the files decides when
// to inject it.
func FontFaces(d font.Declaration, urlPrefix string) *Stylesheet {
	family := d.Family()
	if family == "" {
		return NewStylesheet()
	}
	items := make([]item, 0, 4)
	for s := font.Regular; s <= font.BoldItalic; s++ {
		weight, style := faceWeightStyle(s)
		src := `url("` + joinURLPrefix(urlPrefix, family.Face(s)+".ttf") + `") format("truetype")`
		items = append(items, fontFaceItem{decls: []decl{
			{Prop: "font-family", Val: `"` + string(family) + `"`},
			{Prop: "font-style", Val: style},
			{Prop: "font-weight", Val: weight},
			{Prop: "font-display", Val: "swap"},
			{Prop: "src", Val: src},
		}})
	}
	return NewStylesheet(items...)
}

func faceWeightStyle(s font.Style) (weight, style string) {
	switch s {
	case font.Regular:
		return "400", "normal"
	case font.Bold:
		return "700", "normal"
	case font.Italic:
		return "400", "italic"
	case font.BoldItalic:
		return "700", "italic"
	}
	return "400", "normal"
}

func joinURLPrefix(prefix, file string) string {
	if prefix == "" {
		return file
	}
	if prefix[len(prefix)-1] == '/' {
		return prefix + file
	}
	return prefix + "/" + file
}

// fontFaceItem is the private @font-face at-rule, same pattern as layerItem,
// mediaItem and keyframesItem. Not exported: the DSL stays internal.
type fontFaceItem struct {
	decls []decl
}

func (f fontFaceItem) writeTo(b *fmt.Builder) {
	b.WriteString("@font-face {\n")
	for _, d := range f.decls {
		b.WriteString("  ")
		b.WriteString(d.Prop)
		b.WriteString(": ")
		b.WriteString(d.Val)
		b.WriteString(";\n")
	}
	b.WriteString("}\n\n")
}
