//go:build !wasm

package css

import (
	"github.com/tinywasm/fmt"
)

type Stylesheet struct{ items []item }
type item interface{ writeTo(b *fmt.Builder) }

func NewStylesheet(items ...item) *Stylesheet { return &Stylesheet{items} }

func (s *Stylesheet) String() string {
	// No PutConv here: Conv.String() already returns the object to the pool.
	// Releasing it twice puts the same pointer in the pool twice, so two later
	// Convert() calls hand back the same buffer and corrupt each other.
	b := fmt.GetConv()
	for _, it := range s.items {
		it.writeTo(b)
	}
	return b.String()
}

// Selector is a raw CSS selector string used by the DSL.
type selector string

func (s selector) cssValue() string { return string(s) }


type ruleItem struct {
	sel   string
	decls []decl
}

func (r ruleItem) writeTo(b *fmt.Builder) {
	b.WriteString(r.sel)
	b.WriteString(" {\n")
	for _, d := range r.decls {
		if d.Prop == "raw" {
			b.WriteString(d.Val)
			b.WriteString("\n")
			continue
		}
		b.WriteString("  ")
		b.WriteString(d.Prop)
		b.WriteString(": ")
		b.WriteString(d.Val)
		b.WriteString(";\n")
	}
	b.WriteString("}\n\n")
}

type ruleContent interface{ ruleContent() }

func (d decl) ruleContent() {}

type rawRuleType string

func (r rawRuleType) ruleContent() {}

// RawRule is a transitional escape hatch. New code SHOULD use the typed DSL.
// Each RawRule call site SHOULD carry a // TODO(css-dsl): add typed X comment
// naming the missing property, so reviewers can decide whether to extend the
// DSL or accept the raw use case (vendor-prefixed, exotic property, etc.).
func rawRule(s string) rawRuleType { return rawRuleType(s) }

func ensureSemicolon(s string) string {
	i := len(s) - 1
	for i >= 0 && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i--
	}
	if i >= 0 && s[i] != ';' && s[i] != '}' {
		return s + ";"
	}
	return s
}

func rule(sel selector, content ...ruleContent) item {
	s := string(sel)
	decls := make([]decl, 0, len(content))
	for _, c := range content {
		switch v := c.(type) {
		case decl:
			decls = append(decls, v)
		case rawRuleType:
			decls = append(decls, decl{Prop: "raw", Val: ensureSemicolon(string(v))})
		}
	}
	return ruleItem{sel: s, decls: decls}
}

func root(decls ...decl) item {
	return ruleItem{sel: ":root", decls: decls}
}

type mediaItem struct {
	query string
	items []item
}

func (m mediaItem) writeTo(b *fmt.Builder) {
	b.WriteString("@media ")
	b.WriteString(m.query)
	b.WriteString(" {\n")
	for _, it := range m.items {
		switch s := it.(type) {
		case ruleItem:
			b.WriteString("  ")
			b.WriteString(s.sel)
			b.WriteString(" {\n")
			for _, d := range s.decls {
				if d.Prop == "raw" {
					b.WriteString("    ")
					b.WriteString(d.Val)
					b.WriteString("\n")
					continue
				}
				b.WriteString("    ")
				b.WriteString(d.Prop)
				b.WriteString(": ")
				b.WriteString(d.Val)
				b.WriteString(";\n")
			}
			b.WriteString("  }\n\n")
		case *ruleItem:
			b.WriteString("  ")
			b.WriteString(s.sel)
			b.WriteString(" {\n")
			for _, d := range s.decls {
				if d.Prop == "raw" {
					b.WriteString("    ")
					b.WriteString(d.Val)
					b.WriteString("\n")
					continue
				}
				b.WriteString("    ")
				b.WriteString(d.Prop)
				b.WriteString(": ")
				b.WriteString(d.Val)
				b.WriteString(";\n")
			}
			b.WriteString("  }\n\n")
		default:
			it.writeTo(b)
		}
	}
	b.WriteString("}\n\n")
}

func media(query string, items ...item) item {
	return mediaItem{query: query, items: items}
}

func mediaPrefersDark(items ...item) item {
	return media("(prefers-color-scheme: dark)", items...)
}

// MediaDesktop wraps the canonical "landscape + hover" media query used by
// tinywasm layouts to distinguish desktop from mobile.
// Reference: appears 4 times verbatim in platformd Appendix A.
func mediaDesktop(items ...item) item {
	return media("(orientation: landscape) and (hover: hover)", items...)
}

// KeyframeStep is one step of a keyframes animation.
// At is the percentage or named position ("0%", "50%", "100%", "from", "to").
type keyframeStep struct {
	at    string
	Decls []decl
}

// At builds a KeyframeStep. Variadic Decls match the DSL's existing rule shape.
func at(at string, decls ...decl) keyframeStep {
	return keyframeStep{at: at, Decls: decls}
}

type keyframesItem struct {
	name  string
	steps []keyframeStep
}

func (k keyframesItem) writeTo(b *fmt.Builder) {
	b.WriteString("@keyframes ")
	b.WriteString(k.name)
	b.WriteString(" {\n")
	for _, s := range k.steps {
		b.WriteString("  ")
		b.WriteString(s.at)
		b.WriteString(" {\n")
		for _, d := range s.Decls {
			b.WriteString("    ")
			b.WriteString(d.Prop)
			b.WriteString(": ")
			b.WriteString(d.Val)
			b.WriteString(";\n")
		}
		b.WriteString("  }\n")
	}
	b.WriteString("}\n\n")
}

// Keyframes builds an @keyframes at-rule.
func keyframes(name string, steps ...keyframeStep) item {
	return keyframesItem{name: name, steps: steps}
}

type RawItem string

func (r RawItem) writeTo(b *fmt.Builder) {
	b.WriteString(string(r))
	b.WriteString("\n\n")
}

func Raw(css string) item { return RawItem(css) }

type decl struct{ Prop, Val string }

type value interface{ cssValue() string }

func (t Token) cssValue() string { return t.Var() }

type stringValue string

func (s stringValue) cssValue() string { return string(s) }

func px(n int) value         { return stringValue(fmt.Sprintf("%dpx", n)) }
func rem(f float64) value    { return stringValue(fmt.Sprintf("%grem", f)) }
func em(f float64) value     { return stringValue(fmt.Sprintf("%gem", f)) }
func pct(n int) value        { return stringValue(fmt.Sprintf("%d%%", n)) }
func vw(f float64) value     { return stringValue(fmt.Sprintf("%gvw", f)) }
func vh(f float64) value     { return stringValue(fmt.Sprintf("%gvh", f)) }
func calc(expr string) value { return stringValue(fmt.Sprintf("calc(%s)", expr)) }
func hex(s string) value     { return stringValue(s) }
func str(s string) value     { return stringValue(s) }

type kw string

func (k kw) cssValue() string { return string(k) }

var (
	auto    value = kw("auto")
	none    value = kw("none")
	block   value = kw("block")
	flex_   value = kw("flex")
	grid    value = kw("grid")
	inline  value = kw("inline-block")
	center  value = kw("center")
	zero    value = kw("0")
	pointer value = kw("pointer")

	fixed       value = kw("fixed")
	absolute    value = kw("absolute")
	unset       value = kw("unset")
	initial     value = kw("initial")
	flexEnd     value = kw("flex-end")
	spaceAround value = kw("space-around")
	row         value = kw("row")
	column      value = kw("column")
	hidden      value = kw("hidden")
	visible     value = kw("visible")
	uppercase   value = kw("uppercase")
	capitalize  value = kw("capitalize")
	rightText   value = kw("right")

	relative     value = kw("relative")
	spaceBetween value = kw("space-between")
	inlineFlex   value = kw("inline-flex")
	wrap         value = kw("wrap")
	flexStart    value = kw("flex-start")
	noRepeat     value = kw("no-repeat")
)

func joinValues(vs []value) string {
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = v.cssValue()
	}
	// See Stylesheet.String: Conv.String() already releases b to the pool.
	b := fmt.GetConv()
	for i, p := range parts {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(p)
	}
	return b.String()
}

func background(v value) decl          { return decl{"background", v.cssValue()} }
func backgroundColor(v value) decl     { return decl{"background-color", v.cssValue()} }
func backgroundImage(v value) decl     { return decl{"background-image", v.cssValue()} }
func color(v value) decl               { return decl{"color", v.cssValue()} }
func padding(v ...value) decl          { return decl{"padding", joinValues(v)} }
func margin(v ...value) decl           { return decl{"margin", joinValues(v)} }
func border(v ...value) decl           { return decl{"border", joinValues(v)} }
func borderColor(v value) decl         { return decl{"border-color", v.cssValue()} }
func borderRadius(v ...value) decl     { return decl{"border-radius", joinValues(v)} }
func boxShadow(v value) decl           { return decl{"box-shadow", v.cssValue()} }
func boxSizing(v value) decl           { return decl{"box-sizing", v.cssValue()} }
func display(v value) decl             { return decl{"display", v.cssValue()} }
func flex(v ...value) decl             { return decl{"flex", joinValues(v)} }
func flexDirection(v value) decl       { return decl{"flex-direction", v.cssValue()} }
func gap(v value) decl                 { return decl{"gap", v.cssValue()} }
func justifyContent(v value) decl      { return decl{"justify-content", v.cssValue()} }
func alignItems(v value) decl          { return decl{"align-items", v.cssValue()} }
func width(v value) decl               { return decl{"width", v.cssValue()} }
func height(v value) decl              { return decl{"height", v.cssValue()} }
func maxWidth(v value) decl            { return decl{"max-width", v.cssValue()} }
func minHeight(v value) decl           { return decl{"min-height", v.cssValue()} }
func fontSize(v value) decl            { return decl{"font-size", v.cssValue()} }
func fontWeight(v value) decl          { return decl{"font-weight", v.cssValue()} }
func lineHeight(v value) decl          { return decl{"line-height", v.cssValue()} }
func letterSpacing(v value) decl       { return decl{"letter-spacing", v.cssValue()} }
func transition(v ...value) decl       { return decl{"transition", joinValues(v)} }
func animation(v ...value) decl        { return decl{"animation", joinValues(v)} }
func transform(v value) decl           { return decl{"transform", v.cssValue()} }
func cursor(v value) decl              { return decl{"cursor", v.cssValue()} }
func outline(v value) decl             { return decl{"outline", v.cssValue()} }
func outlineOffset(v value) decl       { return decl{"outline-offset", v.cssValue()} }
func opacity(v float64) decl           { return decl{"opacity", fmt.Sprintf("%g", v)} }
func pointerEvents(v value) decl       { return decl{"pointer-events", v.cssValue()} }
func position(v value) decl            { return decl{"position", v.cssValue()} }
func top(v value) decl                 { return decl{"top", v.cssValue()} }
func right(v value) decl               { return decl{"right", v.cssValue()} }
func bottom(v value) decl              { return decl{"bottom", v.cssValue()} }
func left(v value) decl                { return decl{"left", v.cssValue()} }
func zIndex(v value) decl              { return decl{"z-index", v.cssValue()} }
func fontFamily(v value) decl          { return decl{"font-family", v.cssValue()} }
func minWidth(v value) decl            { return decl{"min-width", v.cssValue()} }
func maxHeight(v value) decl           { return decl{"max-height", v.cssValue()} }
func alignSelf(v value) decl           { return decl{"align-self", v.cssValue()} }
func overflow(v value) decl            { return decl{"overflow", v.cssValue()} }
func visibility(v value) decl          { return decl{"visibility", v.cssValue()} }
func textAlign(v value) decl           { return decl{"text-align", v.cssValue()} }
func textTransform(v value) decl       { return decl{"text-transform", v.cssValue()} }
func textDecoration(v value) decl      { return decl{"text-decoration", v.cssValue()} }
func textShadow(v ...value) decl       { return decl{"text-shadow", joinValues(v)} }
func userSelect(v value) decl          { return decl{"user-select", v.cssValue()} }
func touchAction(v value) decl         { return decl{"touch-action", v.cssValue()} }
func listStyleType(v value) decl       { return decl{"list-style-type", v.cssValue()} }
func gridArea(v value) decl            { return decl{"grid-area", v.cssValue()} }
func gridTemplate(v value) decl        { return decl{"grid-template", v.cssValue()} }
func marginLeft(v value) decl          { return decl{"margin-left", v.cssValue()} }
func marginRight(v value) decl         { return decl{"margin-right", v.cssValue()} }
func paddingBottom(v value) decl       { return decl{"padding-bottom", v.cssValue()} }
func listStyle(v value) decl           { return decl{"list-style", v.cssValue()} }
func all(v value) decl                 { return decl{"all", v.cssValue()} }
func overflowY(v value) decl           { return decl{"overflow-y", v.cssValue()} }
func gridTemplateRows(v value) decl    { return decl{"grid-template-rows", v.cssValue()} }
func gridTemplateColumns(v value) decl { return decl{"grid-template-columns", v.cssValue()} }
func borderRight(v ...value) decl      { return decl{"border-right", joinValues(v)} }
func paddingTop(v value) decl          { return decl{"padding-top", v.cssValue()} }
func paddingLeft(v value) decl         { return decl{"padding-left", v.cssValue()} }
func paddingRight(v value) decl        { return decl{"padding-right", v.cssValue()} }
func marginTop(v value) decl           { return decl{"margin-top", v.cssValue()} }
func marginBottom(v value) decl        { return decl{"margin-bottom", v.cssValue()} }
func flexWrap(v value) decl            { return decl{"flex-wrap", v.cssValue()} }
func flexGrow(v value) decl            { return decl{"flex-grow", v.cssValue()} }
func alignContent(v value) decl        { return decl{"align-content", v.cssValue()} }
func borderBottom(v ...value) decl     { return decl{"border-bottom", joinValues(v)} }
func borderLeft(v ...value) decl       { return decl{"border-left", joinValues(v)} }
func backgroundSize(v value) decl      { return decl{"background-size", v.cssValue()} }
func backgroundPosition(v value) decl  { return decl{"background-position", v.cssValue()} }
func backgroundRepeat(v value) decl    { return decl{"background-repeat", v.cssValue()} }

func declare(t Token) decl {
	return decl{t.Name, t.Fallback}
}
