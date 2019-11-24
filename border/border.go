package border

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Border defines methods required of all border providers.
type Border interface {
	// Insets returns the insets describing the space the border occupies on
	// each side.
	Insets() geom.Insets
	// Draw the border into rect.
	Draw(gc draw.Context, rect geom.Rect, inLiveResize bool)
}
