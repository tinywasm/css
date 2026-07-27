# Justification for the Typed CSS DSL in Go

> **ARCHITECTURE NOTE:** This document justifies the design decisions of the typed CSS DSL in Go. Following the final redesign, this low-level DSL has been **unexported from the public surface** to prevent widget developers from escaping the design system by writing manual CSS properties. Instead, the DSL now acts as the internal emission engine, consumed exclusively by the semantic intent API of `github.com/tinywasm/widget/style`.

> Analysis document. It answers the question: **is this API the most intuitive, readable, and professional way to write CSS in Go for the tinywasm ecosystem?**

## 1. The Proposed Form (Internal Emission Engine)

```go
//go:build !wasm
package button

import . "github.com/tinywasm/css"

var (
    ClsBtn     Class = "btn"
    ClsPrimary Class = "btn-primary"
)

func (b *Button) RenderCSS() *Stylesheet {
    return NewStylesheet(
        rule(ClsBtn,
            padding(rem(0.5), rem(1)),
            borderRadius(radiusSm),
            cursor(pointer),
            fontSize(textBase),
        ),
        rule(ClsPrimary,
            background(ColorPrimary),
            color(ColorOnPrimary),
        ),
        rule(Hover(ClsPrimary),
            opacity(0.9),
        ),
    )
}
```

## 2. Evaluation Criteria

To answer honestly, we must define what "better" means. These are the three criteria defined:

| Criterion | Operational Definition |
|---|---|
| **Intuitive** | A developer who knows CSS recognizes the intent in seconds without reading documentation. |
| **Readable** | A reader who did not write the code can return weeks later and reconstruct the mental model effortlessly. |
| **Professional** | Coherent with established practices in consolidated frameworks; supports refactoring, testing, and IDE tools. |

To these, I add two non-negotiable technical criteria of the tinywasm ecosystem, because they disqualify several theoretically valid alternatives:

| Technical Criterion | Reason |
|---|---|
| **Zero CSS in the WASM binary** | The framework optimizes size with TinyGo. Any solution that drags CSS generation code into the frontend is disqualified. |
| **No generators** | `go generate` and external tools add debt; the ecosystem has already rejected this path. |

---

## 3. Evaluation Against Criteria

### 3.1 Intuitive — Is it understood without documentation?

**Yes, with two reasonable assumptions.**

| Line in the DSL | Obvious CSS Equivalent |
|---|---|
| `rule(ClsPrimary, background(ColorPrimary))` | `.btn-primary { background: var(--color-primary); }` |
| `rule(Hover(ClsPrimary), opacity(0.9))` | `.btn-primary:hover { opacity: 0.9; }` |
| `padding(rem(0.5), rem(1))` | `padding: 0.5rem 1rem;` |

The only non-trivial mapping is `Cls<X>` ↔ class selector. Once learned (one paragraph in README) the mapping is 1:1. Compared to TypeScript-CSS-in-JS (where one has to learn camelCase, units as strings vs numbers, theme objects, variants), the barrier is clearly lower.

### 3.2 Readable — Does it survive six months later?

**Better than the current raw CSS string.** Three reasons:

1. **Tokens are names, not duplicated values.** Today in `button.css` one reads `var(--color-primary, #00ADD8)` — the reader has to decide if `#00ADD8` is an intentional fallback or outdated junk. With the DSL, it reads `background(ColorPrimary)`: zero ambiguity.
2. **Selectors are references, not strings.** `Hover(ClsPrimary)` prevents the silent typo `.btn-primry:hover {}` which today is invisible until opening the browser.
3. **The IDE assists.** "Find references" on `ColorPrimary` lists all uses in the repository. On `--color-primary` in strings, the IDE gives partial and mixed false-positive results.

**Loses to raw CSS in one point:** visual density. `padding: 0.5rem 1rem;` takes fewer pixels than `padding(rem(0.5), rem(1))`. This is mitigated with the dot-import (without it, it would be worse) but not eliminated. This is the honest cost.

### 3.3 Professional — Does it hold up under industry rigor?

**Yes, with direct precedent.** The pattern (typed DSL + tokens as constants + classes as typed identifiers + static CSS extraction) is exactly what the following do:

| System | Language | Equivalent Pattern |
|---|---|---|
| **vanilla-extract** | TypeScript | `style({ background: vars.color.primary })` → CSS extracted in build |
| **Linaria** | TypeScript | Tagged templates with static extraction |
| **Stitches** | TypeScript | `styled('button', { variants: {...} })`, typed tokens |
| **JetBrains Compose HTML** | Kotlin | Rule builder DSL, typed tokens |
| **ScalaCSS** | Scala | Pure Scala DSL, generated class names |
| **W3C Design Tokens CG** | (specification) | Standardizes the concept of "token" as a typed entity |

This is not a local invention. It is the industrial convergence of the last decade applied to Go. The difference: TypeScript needs an additional compiler (vanilla-extract uses a Babel/esbuild plugin); in Go, `//go:build !wasm` is enough for the code to exist only on the server.

### 3.4 Zero CSS in the WASM binary

**Guaranteed by construction.** The layout of the `tinywasm/css` package:

| File | Build tag | Compiles to WASM |
|---|---|---|
| `tokens.go` (Class, Token, constants) | none | ✅ |
| `dsl.go` (Stylesheet, Rule, properties) | `!wasm` | ❌ |
| `css.go` (RootCSS, RenderCSS) | `!wasm` | ❌ |

The only things that cross to the WASM binary are the class name strings (`"btn-primary"`) that the HTML needs to emit. TinyGo also eliminates unreferenced tokens via dead-code elimination. The CSS generator does not exist in the frontend.

### 3.5 No generators

**Fulfilled.** There is no `go generate`, no `theme.css`, no additional build step. The Go compiler is the only tool.

---

## 4. Alternatives Evaluated and Discarded

| Alternative | Reason for Discarding |
|---|---|
| **Keep `.css` + `//go:embed`** (current state) | Stringly-typed; no errors are detected until opening the browser; renaming tokens is manual and prone to drift. |
| **`.css` + external linter** | Solves typo detection but adds external tool; still two languages. |
| **`theme.css → tokens.go` generator** | Maintains two representations; generator is permanent debt; contradicts "no generators". |
| **CSS-in-Go runtime styled-components style** | Drags the CSS engine into the WASM binary. Disqualified immediately. |
| **Text templates (`text/template`)** | Returns to untyped strings; loses compiler validation. |
| **Fluent DSL (builder with `.Padding(...).Color(...)`)** | Expressively equivalent to the variadic constructor, worse for long rules (uncomfortable vertical chaining); variadic flattens better. |
| **Sub-package `tinywasm/css/cssgo`** | Forces two imports, breaks dot-import. No real value. |

---

## 5. Honest Risks (Not Hidden)

A professional analysis must name what the pattern loses or complicates:

1. **Verbosity relative to raw CSS.** `padding(rem(0.5), rem(1))` is more characters than `padding: 0.5rem 1rem`. Mitigation: dot-import; authors learn to visually scan the `rule(sel, ...decls)` structure.

2. **`@media` and rare selectors go through `selector("...")` or `media("...")`.** For `@container`, complex attribute selectors, `:nth-child(...)`, the API becomes a string escape hatch. It is not elegant but it is honest: covering 100% of the CSS spec in typed Go is a disproportionate effort. The DSL prioritizes the common 90%.

3. **`css.go` can grow.** A component with 200 lines of CSS becomes 200 lines of Go in `css.go`. It is the same volume, no more. If it hurts in some specific component, a `css_styles.go` in the same package is allowed as an exception.

4. **Initial adoption curve.** A contributor who only knows classic CSS needs half an hour to internalize the mapping. Unique cost; the benefit is permanent.

5. **No CSS-aware formatting tooling.** `gofmt` does not know how to align declarations like `prettier` aligns CSS. Mitigation: the variadic constructor naturally produces one declaration per line; formatting is predictable.

### 5.1 On Verbosity: Variadic Units?

Natural question: if `padding(rem(0.5), rem(1))` is longer than `padding: 0.5rem 1rem`, why not make `rem` variadic → `padding(rem(0.5, 1))`?

**Discarded.** Reasons, in order of weight:

1. **Type muddle.** `rem` is *a unit*, not a list. Making it variadic turns `Value` into "atomic value or serialized sequence", contaminating the whole system. It enables compileable junk like `color(rem(0.5, 1))` → `color: 0.5rem 1rem;`.
2. **Negligible savings.** ~5 characters × ~15% shorthand declarations = ~75 chars in the whole project.
3. **The DSL already resolves shorthand in the right place:** the property (`padding`) is variadic, just like the CSS grammar (`padding: <length>{1,4}`). Moving variadicity to the unit shifts responsibility to the wrong place.
4. **Blocks real mixed units:** `boxShadow(em(0.1), em(0.1), em(0.2), ColorSurface)` mixes `em` with a token — impossible if `em` or `rem` become variadic.

If the team decides in the future that density matters more than type purity, the correct path would be *direct floats with implicit unit per property* (`padding(0.5, 1)` → rem), not variadic units. Separate analysis.

---

## 6. Is it *the Best* or Only *a Good One*?

Here we must distinguish two questions:

### 6.1 Is it the best way to write CSS in a typed language?
There is a legitimate debate. **vanilla-extract in TypeScript** is probably more mature today in absolute terms. But for a **Go-first + TinyGo + no generators** project, constraints eliminate TypeScript from the set of applicable solutions.

### 6.2 Is it the best way to write CSS in Go for tinywasm?
**Yes, within the set of solutions compatible with the project's constraints.** I do not know of an alternative that simultaneously fulfills:
- Zero CSS in WASM binary
- No generators
- Compile-time typo detection
- Tokens as first-class entities
- Shared class names between HTML and CSS
- A single way to do it

The first five exist isolated in other proposals; none gathers them all.

---

## 7. Verdict

**Yes, it is the most intuitive, readable, and professional way to write CSS in Go for tinywasm**, conditioned on accepting the costs named in section 5 (especially verbosity and the initial curve). It is professionally defensible because it replicates a pattern with ~10 years of industrial adoption in other typed languages, adapted to the specific constraints of Go + TinyGo + SSR architecture of the project.

If the verdict does not convince, the points to question first are:
1. Are the costs of section 5 acceptable to your team?
2. Is the 90% CSS coverage (with escape hatch for the rest) sufficient, or does the project need the full CSS spec?
3. Is dot-import acceptable as an ecosystem convention?

If all three are "yes", the pattern is the right one. If any is "no", the analysis must be revisited before executing the plans.

---

## References

- W3C Design Tokens Community Group: <https://design-tokens.github.io/community-group/>
- vanilla-extract: <https://vanilla-extract.style/>
- Linaria: <https://linaria.dev/>
- Stitches (archived but conceptual reference): <https://stitches.dev/>
- Lightning Design System (origin of the term "design token"): <https://www.lightningdesignsystem.com/design-tokens/>
