package border

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Compound provides stacking of borders together.
type Compound struct {
	borders []Border
}

// NewCompound creates a border that contains other borders. The first one
// will be drawn in the outermost position, with each successive one moving
// further into the interior.
func NewCompound(borders ...Border) *Compound {
	return &Compound{borders: borders}
}

// Insets returns the insets describing the space the border occupies on each
// side.
func (b *Compound) Insets() geom.Insets {
	insets := geom.Insets{}
	for _, one := range b.borders {
		insets.Add(one.Insets())
	}
	return insets
}

// Draw the border into rect.
func (b *Compound) Draw(gc draw.Context, rect geom.Rect, inLiveResize bool) {
	for _, one := range b.borders {
		gc.Save()
		one.Draw(gc, rect, inLiveResize)
		gc.Restore()
		rect.Inset(one.Insets())
	}
}
