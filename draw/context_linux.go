package draw

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
	"github.com/richardwilkes/ux/draw/linecap"
	"github.com/richardwilkes/ux/draw/linejoin"
	"github.com/richardwilkes/ux/draw/quality"
)

// OSContext is the platform-specific drawing context on Windows.
type OSContext = int // RAW: Implement

type context struct {
}

func osNewContextForOSContext(gc OSContext) Context {
	return &context{} // RAW: Implement
}

func (c *context) OSContext() OSContext {
	return 0 // RAW: Implement
}

func (c *context) Save() {
	// RAW: Implement
}

func (c *context) Restore() {
	// RAW: Implement
}

func (c *context) SetOpacity(opacity float64) {
	// RAW: Implement
}

func (c *context) SetPatternOffset(x, y float64) {
	// RAW: Implement
}

func (c *context) SetStrokeWidth(width float64) {
	// RAW: Implement
}

func (c *context) SetLineCap(lineCap linecap.LineCap) {
	// RAW: Implement
}

func (c *context) SetLineJoin(lineJoin linejoin.LineJoin) {
	// RAW: Implement
}

func (c *context) SetMiterLimit(limit float64) {
	// RAW: Implement
}

func (c *context) SetLineDash(phase float64, lengths ...float64) {
	// RAW: Implement
}

func (c *context) SetInterpolationQualityHint(q quality.Quality) {
	// RAW: Implement
}

func (c *context) Translate(x, y float64) {
	// RAW: Implement
}

func (c *context) Scale(x, y float64) {
	// RAW: Implement
}

func (c *context) Rotate(angleInRadians float64) {
	// RAW: Implement
}

func (c *context) Fill(ink Ink) {
	ink.osFill(c)
}

func (c *context) FillEvenOdd(ink Ink) {
	ink.osFillEvenOdd(c)
}

func (c *context) Stroke(ink Ink) {
	ink.osStroke(c)
}

func (c *context) GetClipRect() geom.Rect {
	return geom.Rect{} // RAW: Implement
}

func (c *context) Clip() {
	// RAW: Implement
}

func (c *context) ClipEvenOdd() {
	// RAW: Implement
}

func (c *context) DrawImage(img *Image, where geom.Point) {
	// RAW: Implement
}

func (c *context) DrawImageInRect(img *Image, rect geom.Rect) {
	// RAW: Implement
}

func (c *context) DrawString(x, y float64, font *Font, ink Ink, str string) {
	// RAW: Implement
}

func (c *context) BeginPath() {
	// RAW: Implement
}

func (c *context) MoveTo(x, y float64) {
	// RAW: Implement
}

func (c *context) LineTo(x, y float64) {
	// RAW: Implement
}

func (c *context) QuadCurveTo(cpx, cpy, x, y float64) {
	// RAW: Implement
}

func (c *context) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	// RAW: Implement
}

func (c *context) Rect(rect geom.Rect) {
	// RAW: Implement
}

func (c *context) RoundedRect(rect geom.Rect, cornerRadius float64) {
	// RAW: Implement
}

func (c *context) Ellipse(rect geom.Rect) {
	// RAW: Implement
}

func (c *context) Polygon(polygon poly.Polygon) {
	for _, cont := range polygon {
		for i, pt := range cont {
			if i == 0 {
				c.MoveTo(pt.X, pt.Y)
			} else {
				c.LineTo(pt.X, pt.Y)
			}
		}
	}
	c.ClosePath()
}

func (c *context) ClosePath() {
	// RAW: Implement
}

func (c *context) Dispose() {
	// RAW: Implement
}
