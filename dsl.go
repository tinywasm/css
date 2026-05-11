//go:build !wasm

package css

import (
	"github.com/tinywasm/fmt"
)

type Stylesheet struct{ items []item }
type item interface{ writeTo(b *fmt.Builder) }

func New(items ...item) *Stylesheet { return &Stylesheet{items} }

func (s *Stylesheet) String() string {
	b := fmt.GetConv()
	defer b.PutConv()
	for _, it := range s.items {
		it.writeTo(b)
	}
	return b.String()
}

// Selector is a raw CSS selector string used by the DSL.
type Selector string

func (s Selector) cssValue() string { return string(s) }

// Pseudo-class helpers on Class
func (c Class) Hover() Selector    { return Selector("." + string(c) + ":hover") }
func (c Class) Focus() Selector    { return Selector("." + string(c) + ":focus") }
func (c Class) Disabled() Selector { return Selector("." + string(c) + ":disabled") }

type RuleItem struct {
	sel   string
	decls []Decl
}

func (r RuleItem) writeTo(b *fmt.Builder) {
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

type RuleContent interface{ ruleContent() }

func (d Decl) ruleContent()    {}

type rawRule string

func (r rawRule) ruleContent() {}

func RawRule(s string) rawRule { return rawRule(s) }

func Rule(sel any, content ...RuleContent) item {
	var s string
	switch v := sel.(type) {
	case Class:
		s = "." + string(v)
	case Selector:
		s = string(v)
	case string:
		s = v
	}
	decls := make([]Decl, 0, len(content))
	for _, c := range content {
		switch v := c.(type) {
		case Decl:
			decls = append(decls, v)
		case rawRule:
			decls = append(decls, Decl{Prop: "raw", Val: string(v)})
		}
	}
	return RuleItem{sel: s, decls: decls}
}

func Root(decls ...Decl) item {
	return RuleItem{sel: ":root", decls: decls}
}

type MediaItem struct {
	query string
	items []item
}

func (m MediaItem) writeTo(b *fmt.Builder) {
	b.WriteString("@media ")
	b.WriteString(m.query)
	b.WriteString(" {\n")
	for _, it := range m.items {
		switch s := it.(type) {
		case RuleItem:
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
		case *RuleItem:
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

func Media(query string, items ...item) item {
	return MediaItem{query: query, items: items}
}

func MediaPrefersDark(items ...item) item {
	return Media("(prefers-color-scheme: dark)", items...)
}

// KeyframeStep is one step of a keyframes animation.
// At is the percentage or named position ("0%", "50%", "100%", "from", "to").
type KeyframeStep struct {
	At    string
	Decls []Decl
}

// At builds a KeyframeStep. Variadic Decls match the DSL's existing rule shape.
func At(at string, decls ...Decl) KeyframeStep {
	return KeyframeStep{At: at, Decls: decls}
}

type KeyframesItem struct {
	name  string
	steps []KeyframeStep
}

func (k KeyframesItem) writeTo(b *fmt.Builder) {
	b.WriteString("@keyframes ")
	b.WriteString(k.name)
	b.WriteString(" {\n")
	for _, s := range k.steps {
		b.WriteString("  ")
		b.WriteString(s.At)
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
func Keyframes(name string, steps ...KeyframeStep) item {
	return KeyframesItem{name: name, steps: steps}
}

type RawItem string

func (r RawItem) writeTo(b *fmt.Builder) {
	b.WriteString(string(r))
	b.WriteString("\n\n")
}

func Raw(css string) item { return RawItem(css) }

type Decl struct{ Prop, Val string }

type Value interface{ cssValue() string }

func (t Token) cssValue() string { return t.Var() }

type stringValue string

func (s stringValue) cssValue() string { return string(s) }

func Px(n int) Value      { return stringValue(fmt.Sprintf("%dpx", n)) }
func Rem(f float64) Value { return stringValue(fmt.Sprintf("%grem", f)) }
func Em(f float64) Value  { return stringValue(fmt.Sprintf("%gem", f)) }
func Pct(n int) Value     { return stringValue(fmt.Sprintf("%d%%", n)) }
func Hex(s string) Value  { return stringValue(s) }
func Str(s string) Value  { return stringValue(s) }

type kw string

func (k kw) cssValue() string { return string(k) }

var (
	Auto    Value = kw("auto")
	None    Value = kw("none")
	Block   Value = kw("block")
	Flex_   Value = kw("flex")
	Grid    Value = kw("grid")
	Inline  Value = kw("inline-block")
	Center  Value = kw("center")
	Zero    Value = kw("0")
	Pointer Value = kw("pointer")
)

func joinValues(vs []Value) string {
	b := fmt.GetConv()
	defer b.PutConv()
	for i, v := range vs {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(v.cssValue())
	}
	return b.String()
}

func Background(v Value) Decl      { return Decl{"background", v.cssValue()} }
func BackgroundColor(v Value) Decl { return Decl{"background-color", v.cssValue()} }
func BackgroundImage(v Value) Decl { return Decl{"background-image", v.cssValue()} }
func Color(v Value) Decl           { return Decl{"color", v.cssValue()} }
func Padding(v ...Value) Decl      { return Decl{"padding", joinValues(v)} }
func Margin(v ...Value) Decl       { return Decl{"margin", joinValues(v)} }
func Border(v ...Value) Decl       { return Decl{"border", joinValues(v)} }
func BorderColor(v Value) Decl     { return Decl{"border-color", v.cssValue()} }
func BorderRadius(v ...Value) Decl { return Decl{"border-radius", joinValues(v)} }
func BoxShadow(v Value) Decl       { return Decl{"box-shadow", v.cssValue()} }
func BoxSizing(v Value) Decl       { return Decl{"box-sizing", v.cssValue()} }
func Display(v Value) Decl         { return Decl{"display", v.cssValue()} }
func Flex(v ...Value) Decl         { return Decl{"flex", joinValues(v)} }
func FlexDirection(v Value) Decl   { return Decl{"flex-direction", v.cssValue()} }
func Gap(v Value) Decl             { return Decl{"gap", v.cssValue()} }
func JustifyContent(v Value) Decl  { return Decl{"justify-content", v.cssValue()} }
func AlignItems(v Value) Decl      { return Decl{"align-items", v.cssValue()} }
func Width(v Value) Decl           { return Decl{"width", v.cssValue()} }
func Height(v Value) Decl          { return Decl{"height", v.cssValue()} }
func MaxWidth(v Value) Decl        { return Decl{"max-width", v.cssValue()} }
func MinHeight(v Value) Decl       { return Decl{"min-height", v.cssValue()} }
func FontSize(v Value) Decl        { return Decl{"font-size", v.cssValue()} }
func FontWeight(v Value) Decl      { return Decl{"font-weight", v.cssValue()} }
func LineHeight(v Value) Decl      { return Decl{"line-height", v.cssValue()} }
func LetterSpacing(v Value) Decl   { return Decl{"letter-spacing", v.cssValue()} }
func Transition(v ...Value) Decl   { return Decl{"transition", joinValues(v)} }
func Animation(v ...Value) Decl    { return Decl{"animation", joinValues(v)} }
func Transform(v Value) Decl       { return Decl{"transform", v.cssValue()} }
func Cursor(v Value) Decl          { return Decl{"cursor", v.cssValue()} }
func Outline(v Value) Decl         { return Decl{"outline", v.cssValue()} }
func OutlineOffset(v Value) Decl   { return Decl{"outline-offset", v.cssValue()} }
func Opacity(v float64) Decl       { return Decl{"opacity", fmt.Sprintf("%g", v)} }
func PointerEvents(v Value) Decl   { return Decl{"pointer-events", v.cssValue()} }
func Position(v Value) Decl        { return Decl{"position", v.cssValue()} }
func Top(v Value) Decl             { return Decl{"top", v.cssValue()} }
func Right(v Value) Decl           { return Decl{"right", v.cssValue()} }
func Bottom(v Value) Decl          { return Decl{"bottom", v.cssValue()} }
func Left(v Value) Decl            { return Decl{"left", v.cssValue()} }
func ZIndex(v Value) Decl          { return Decl{"z-index", v.cssValue()} }
func FontFamily(v Value) Decl      { return Decl{"font-family", v.cssValue()} }

func Declare(t Token, value string) Decl {
	return Decl{t.Name, value}
}

func Bind(active, source Token) Decl {
	return Decl{active.Name, source.Var()}
}
