package border

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Empty provides and empty border with the specified insets.
type Empty struct {
	insets geom.Insets
}

// NewEmpty creates a new empty border with the specified insets.
func NewEmpty(insets geom.Insets) *Empty {
	return &Empty{insets: insets}
}

// Insets returns the insets describing the space the border occupies on each
// side.
func (b *Empty) Insets() geom.Insets {
	return b.insets
}

// Draw the border into rect.
func (b *Empty) Draw(gc draw.Context, rect geom.Rect, inLiveResize bool) {
}
