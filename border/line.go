package border

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Line private a lined border.
type Line struct {
	insets  geom.Insets
	ink     draw.Ink
	noInset bool
}

// NewLine creates a new line border. The insets represent how thick the
// border will be drawn on that edge. If noInset is true, the Insets() method
// will return zeroes.
func NewLine(ink draw.Ink, insets geom.Insets, noInset bool) *Line {
	return &Line{
		insets:  insets,
		ink:     ink,
		noInset: noInset,
	}
}

// Insets returns the insets describing the space the border occupies on each
// side.
func (b *Line) Insets() geom.Insets {
	if b.noInset {
		return geom.Insets{}
	}
	return b.insets
}

// Draw the border into rect.
func (b *Line) Draw(gc draw.Context, rect geom.Rect, inLiveResize bool) {
	clip := rect
	clip.Inset(b.insets)
	gc.BeginPath()
	gc.MoveTo(rect.X, rect.Y)
	gc.LineTo(rect.X+rect.Width, rect.Y)
	gc.LineTo(rect.X+rect.Width, rect.Y+rect.Height)
	gc.LineTo(rect.X, rect.Y+rect.Height)
	gc.LineTo(rect.X, rect.Y)
	gc.MoveTo(clip.X, clip.Y)
	gc.LineTo(clip.X+clip.Width, clip.Y)
	gc.LineTo(clip.X+clip.Width, clip.Y+clip.Height)
	gc.LineTo(clip.X, clip.Y+clip.Height)
	gc.LineTo(clip.X, clip.Y)
	gc.FillEvenOdd(b.ink)
}
