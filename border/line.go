package border

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Line private a lined border.
type Line struct {
	insets       geom.Insets
	ink          draw.Ink
	cornerRadius float64
	noInset      bool
}

// NewLine creates a new line border. The cornerRadius specifies the amount of
// rounding to use on the corners. The insets represent how thick the border
// will be drawn on that edge. If noInset is true, the Insets() method will
// return zeroes.
func NewLine(ink draw.Ink, cornerRadius float64, insets geom.Insets, noInset bool) *Line {
	return &Line{
		insets:       insets,
		ink:          ink,
		cornerRadius: cornerRadius,
		noInset:      noInset,
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
	if b.cornerRadius > 0 {
		b.drawRoundedRect(gc, rect, b.cornerRadius)
		b.drawRoundedRect(gc, clip, math.Max(b.cornerRadius-((b.insets.Top+b.insets.Left+b.insets.Bottom+b.insets.Right)/4), 1))
	} else {
		b.drawRect(gc, rect)
		b.drawRect(gc, clip)
	}
	gc.FillEvenOdd(b.ink)
}

func (b *Line) drawRect(gc draw.Context, rect geom.Rect) {
	gc.MoveTo(rect.X, rect.Y)
	gc.LineTo(rect.X+rect.Width, rect.Y)
	gc.LineTo(rect.X+rect.Width, rect.Y+rect.Height)
	gc.LineTo(rect.X, rect.Y+rect.Height)
	gc.LineTo(rect.X, rect.Y)
}

func (b *Line) drawRoundedRect(gc draw.Context, rect geom.Rect, cornerRadius float64) {
	gc.MoveTo(rect.X, rect.Y+cornerRadius)
	gc.QuadCurveTo(rect.X, rect.Y, rect.X+cornerRadius, rect.Y)
	gc.LineTo(rect.X+rect.Width-cornerRadius, rect.Y)
	gc.QuadCurveTo(rect.X+rect.Width, rect.Y, rect.X+rect.Width, rect.Y+cornerRadius)
	gc.LineTo(rect.X+rect.Width, rect.Y+rect.Height-cornerRadius)
	gc.QuadCurveTo(rect.X+rect.Width, rect.Y+rect.Height, rect.X+rect.Width-cornerRadius, rect.Y+rect.Height)
	gc.LineTo(rect.X+cornerRadius, rect.Y+rect.Height)
	gc.QuadCurveTo(rect.X, rect.Y+rect.Height, rect.X, rect.Y+rect.Height-cornerRadius)
	gc.LineTo(rect.X, rect.Y+cornerRadius)
}
