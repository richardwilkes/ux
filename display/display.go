package display

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
)

// Display holds information about each available active display.
type Display struct {
	Frame         geom.Rect
	Usable        geom.Rect
	ScalingFactor float64
	Primary       bool
}

// Primary returns the primary display.
func Primary() *Display {
	for _, d := range All() {
		if d.Primary {
			return d
		}
	}
	return nil
}

// All returns all displays.
func All() []*Display {
	return osDisplays()
}
