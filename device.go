//go:build !wasm

package css

import "github.com/tinywasm/fmt"

// Device is the closed set of viewport classes. It exists because a media query
// condition cannot read a custom property: the pixel thresholds must be baked into
// the query string, and this is the only file in the ecosystem allowed to hold them.
//
// The three classes are mutually exclusive and jointly exhaustive: every viewport
// width matches exactly one. That property is asserted by TestDeviceClassesPartition.
type Device uint8

const (
	Mobile  Device = iota
	Tablet
	Desktop
)

func (d Device) String() string {
	switch d {
	case Mobile:
		return "Mobile"
	case Tablet:
		return "Tablet"
	case Desktop:
		return "Desktop"
	default:
		return "Unknown"
	}
}

// Query returns the media-query condition for exactly this class, without the
// leading "@media ". Thresholds mirror BpSm (640px) and BpLg (1024px).
func (d Device) Query() string {
	switch d {
	case Mobile:
		return "(max-width: 639.98px)"
	case Tablet:
		return "(min-width: 640px) and (max-width: 1023.98px)"
	case Desktop:
		return "(min-width: 1024px)"
	default:
		return ""
	}
}

// Query joins several device classes into one condition list. Callers that mean
// "tablet and desktop" pass both rather than emitting two blocks.
// Duplicate and unknown values are dropped; the result is ordered Mobile,
// Tablet, Desktop regardless of argument order, so emission is deterministic.
func Query(devices ...Device) string {
	if len(devices) == 0 {
		return ""
	}

	var hasMobile, hasTablet, hasDesktop bool
	for _, d := range devices {
		switch d {
		case Mobile:
			hasMobile = true
		case Tablet:
			hasTablet = true
		case Desktop:
			hasDesktop = true
		}
	}

	b := fmt.GetConv()
	count := 0
	appendIf := func(present bool, query string) {
		if !present {
			return
		}
		if count > 0 {
			b.WriteString(", ")
		}
		b.WriteString(query)
		count++
	}

	appendIf(hasMobile, Mobile.Query())
	appendIf(hasTablet, Tablet.Query())
	appendIf(hasDesktop, Desktop.Query())

	return b.String()
}
