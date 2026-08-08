//go:build !wasm

package css

// resetRules are the base reset rules RenderCSS() emits in the `tokens`
// layer — the lowest of the four the widget style DSL declares. Unlayered,
// `svg { display: block }` would beat `@layer widgets { .part { display:
// none } }` regardless of specificity, so any component hiding an icon by
// state could not.
func resetRules() []item {
	return []item{
		rule(selector("*, *::before, *::after"),
			boxSizing(str("border-box")),
		),
		// text-size-adjust stops mobile browsers reflowing text at their own
		// scale. tap-highlight-color is the iOS-only grey wash painted over
		// anything tappable for ~300ms after a touch: Chrome Android uses a
		// different colour and duration, so the same button flashes differently
		// on each. Components own their press state; the UA's is noise.
		rule(selector("html"),
			rawRule("  -webkit-text-size-adjust: 100%;\n  text-size-adjust: 100%;\n  -webkit-tap-highlight-color: transparent;"),
		),
		rule(selector("body"),
			margin(zero),
			fontFamily(FontSans),
			fontSize(TextBase),
			lineHeight(LeadingNormal),
			color(ColorOnSurface),
			background(ColorBackground),
		),
		rule(selector(":focus-visible"),
			outline(str("2px solid "+ColorPrimary.Var())),
			outlineOffset(px(2)),
		),
		// User-agent margins and list chrome are geometry the style DSL cannot
		// express and cannot see. A <ul> used as a nav rail carries 40px of
		// padding-inline-start it never asked for, an <h1> carries 0.67em of
		// block margin: both silently widen or heighten the component around
		// them. Zero them here, in the lowest layer, so a part's box is exactly
		// what its rule declares.
		rule(selector("h1, h2, h3, h4, h5, h6, p, figure, blockquote, dl, dd"),
			margin(zero),
		),
		// Heading sizes are the DSL's business, not the user agent's. An <h1>
		// carries font-size: 2em by default, which compounds against whatever
		// the part around it declares and makes a title three times the size a
		// component asked for. FontSize()/FontWeight() say what a heading looks
		// like; its level stays a semantic choice.
		// An <a> arrives underlined and painted the user agent's link blue, which
		// fights every surface a component puts it on — a nav item ends up blue
		// text on a blue button. Components say what a link looks like; the
		// browser's default is chrome, like the heading sizes below.
		rule(selector("a"),
			rawRule("  color: inherit;\n  text-decoration: none;"),
		),
		rule(selector("h1, h2, h3, h4, h5, h6"),
			rawRule("  font-size: inherit;\n  font-weight: inherit;"),
		),
		// Form controls are the other UA-font holdout: an <input> arrives with
		// the user agent's family and size (Arial 13px on most platforms), which
		// shrinks the value inside a field the design measured for the app font.
		// Inheriting makes the control speak the same type as the text around it;
		// a part can still override size with FontSize().
		rule(selector("input, textarea, select"),
			rawRule("  font: inherit;"),
		),
		// A <button> is the widest cross-browser gap left. iOS Safari renders it
		// with appearance: push-button — its own corner radius, a vertical
		// gradient, a UA border and centred text — and that chrome outranks the
		// background and radius a part declares. Chrome Android paints a flat
		// grey box instead, so the same button is two different shapes. Zeroing
		// appearance, background, border and radius makes the element an empty
		// box on both, which is what a styled part expects to start from.
		// padding: 0 kills the UA's own button padding (1px block, 6px inline)
		// that a bare <button> inside a flex row still carries — a part that
		// declares Pad() rebuilds it anyway, but an unstyled one must not
		// arrive with a lopsided box the design never asked for.
		// The [type=…] selectors cover <input type="submit">, which is a button
		// everywhere except in the selector `button`.
		rule(selector("button, [type=\"button\"], [type=\"reset\"], [type=\"submit\"]"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  background-color: transparent;\n  background-image: none;\n  color: inherit;\n  font: inherit;\n  border: 0;\n  border-radius: 0;\n  padding: 0;"),
		),
		// iOS Safari paints an inset shadow and forces its own corner radius on
		// text fields; no border-radius a part declares removes it. :where()
		// keeps checkbox and radio out — appearance: none on those erases the
		// control entirely instead of flattening it, and their styling belongs
		// to tinywasm/form.
		rule(selector("input:where(:not([type=\"checkbox\"]):not([type=\"radio\"])), textarea"),
			rawRule("  -webkit-appearance: none;\n  appearance: none;\n  border-radius: 0;"),
		),
		// Firefox and Edge let a <select> inherit text-transform from an
		// ancestor; every other engine does not, so an uppercased container
		// silently uppercases the dropdown on two browsers only. The native
		// arrow stays: it is the only affordance the control has, and
		// tinywasm/form owns replacing it.
		rule(selector("select"),
			rawRule("  text-transform: none;"),
		),
		// Firefox ships placeholders at opacity 0.54, so the same muted colour
		// reads lighter there than on Chrome or Safari.
		rule(selector("::placeholder"),
			rawRule("  opacity: 1;"),
		),
		// Every engine renders a monospace default about 3px smaller than the
		// surrounding text — a historical quirk of the `monospace` keyword.
		// Re-stating the family resets that scaling, and 1em pins the size back
		// to the text around it.
		rule(selector("code, kbd, samp, pre"),
			rawRule("  font-family: monospace;\n  font-size: 1em;"),
		),
		rule(selector("ol, ul"),
			rawRule("  list-style: none;"),
			margin(zero),
			padding(zero),
		),
		// A <summary> draws a disclosure triangle the component did not ask
		// for and cannot style. Components supply their own affordance.
		rule(selector("summary"),
			rawRule("  list-style: none;"),
		),
		rule(selector("summary::-webkit-details-marker"),
			display(none),
		),
		rule(selector("img, svg, video"),
			display(block),
			maxWidth(pct(100)),
		),
		// Author styles outrank the user agent stylesheet whatever their layer,
		// so the display: block above silently defeats the UA's own
		// `[hidden] { display: none }` and an <img hidden> keeps rendering.
		// Restating the rule here is what makes the attribute work again.
		rule(selector("[hidden]"),
			display(none),
		),
		rule(selector(":root"),
			rawRule("color-scheme: light dark;"),
		),
		rule(selector("[data-theme=\"light\"]"),
			rawRule("color-scheme: light;"),
		),
		rule(selector("[data-theme=\"dark\"]"),
			rawRule("color-scheme: dark;"),
		),
	}
}
