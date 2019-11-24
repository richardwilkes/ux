package widget

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// DrawRectBase fills and strokes a rectangle.
func DrawRectBase(gc draw.Context, rect geom.Rect, fillInk, strokeInk draw.Ink) {
	gc.Rect(rect)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.Rect(rect)
	gc.Stroke(strokeInk)
}

// DrawRoundedRectBase fills and strokes a rounded rectangle.
func DrawRoundedRectBase(gc draw.Context, rect geom.Rect, cornerRadius float64, fillInk, strokeInk draw.Ink) {
	gc.RoundedRect(rect, cornerRadius)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.RoundedRect(rect, cornerRadius)
	gc.Stroke(strokeInk)
}

// DrawEllipseBase fills and strokes an ellipse.
func DrawEllipseBase(gc draw.Context, rect geom.Rect, fillInk, strokeInk draw.Ink) {
	gc.Ellipse(rect)
	gc.Fill(fillInk)
	rect.InsetUniform(0.5)
	gc.Ellipse(rect)
	gc.Stroke(strokeInk)
}
