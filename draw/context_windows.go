// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
	"github.com/richardwilkes/ux/draw/linecap"
	"github.com/richardwilkes/ux/draw/linejoin"
	"github.com/richardwilkes/ux/draw/quality"
	"github.com/richardwilkes/win32"
)

// OSContext is the platform-specific drawing context on Windows.
type OSContext = win32.HDC

type context struct {
	hdc   OSContext
	brush win32.HBRUSH
	pen   win32.HPEN
}

func osNewContextForOSContext(gc OSContext) Context {
	c := &context{hdc: gc}
	win32.SetTextAlign(c.hdc, win32.TA_TOP|win32.TA_LEFT)
	return c
}

func (c *context) OSContext() OSContext {
	return c.hdc
}

func (c *context) Save() {
	win32.SaveDC(c.hdc)
}

func (c *context) Restore() {
	win32.RestoreDC(c.hdc, -1)
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
	matrix := xmath.NewTranslationMatrix2D(x, y)
	win32.SetWorldTransform(c.hdc, &win32.XFORM{
		EM11: float32(matrix.XX),
		EM12: float32(matrix.YX),
		EM21: float32(matrix.XY),
		EM22: float32(matrix.YY),
		EDx:  float32(matrix.X0),
		EDy:  float32(matrix.Y0),
	})
}

func (c *context) Scale(x, y float64) {
	matrix := xmath.NewScaleMatrix2D(x, y)
	win32.SetWorldTransform(c.hdc, &win32.XFORM{
		EM11: float32(matrix.XX),
		EM12: float32(matrix.YX),
		EM21: float32(matrix.XY),
		EM22: float32(matrix.YY),
		EDx:  float32(matrix.X0),
		EDy:  float32(matrix.Y0),
	})
}

func (c *context) Rotate(angleInRadians float64) {
	matrix := xmath.NewRotationMatrix2D(angleInRadians)
	win32.SetWorldTransform(c.hdc, &win32.XFORM{
		EM11: float32(matrix.XX),
		EM12: float32(matrix.YX),
		EM21: float32(matrix.XY),
		EM22: float32(matrix.YY),
		EDx:  float32(matrix.X0),
		EDy:  float32(matrix.Y0),
	})
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
	win32.EndPath(c.hdc)
	win32.SetPolyFillMode(c.hdc, win32.WINDING)
	win32.SelectClipPath(c.hdc, win32.RGN_AND)
}

func (c *context) ClipEvenOdd() {
	win32.EndPath(c.hdc)
	win32.SetPolyFillMode(c.hdc, win32.ALTERNATE)
	win32.SelectClipPath(c.hdc, win32.RGN_AND)
}

func (c *context) DrawImage(img *Image, where geom.Point) {
	// RAW: Implement
}

func (c *context) DrawImageInRect(img *Image, rect geom.Rect) {
	// RAW: Implement
}

func (c *context) DrawString(x, y float64, font *Font, ink Ink, str string) {
	// RAW: Use ink
	win32.SelectObject(c.hdc, win32.HGDIOBJ(font.ref))
	win32.TextOut(c.hdc, int(x), int(y), str)
}

func (c *context) BeginPath() {
	win32.BeginPath(c.hdc)
}

func (c *context) MoveTo(x, y float64) {
	win32.MoveToEx(c.hdc, int(x), int(y), nil)
}

func (c *context) LineTo(x, y float64) {
	win32.LineTo(c.hdc, int(x), int(y))
}

func (c *context) QuadCurveTo(cpx, cpy, x, y float64) {
	var pos win32.POINT
	win32.GetCurrentPositionEx(c.hdc, &pos)
	var pts [3]win32.POINT
	pts[0].X = int32(float64(pos.X) + (cpx-float64(pos.X))*2/3)
	pts[0].Y = int32(float64(pos.Y) + (cpy-float64(pos.Y))*2/3)
	pts[1].X = int32(float64(x) + (cpx-x)*2/3)
	pts[1].Y = int32(float64(y) + (cpy-y)*2/3)
	pts[2].X = int32(x)
	pts[2].Y = int32(y)
	win32.PolyBezierTo(c.hdc, pts[:])
}

func (c *context) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	var pts [3]win32.POINT
	pts[0].X = int32(cp1x)
	pts[0].Y = int32(cp1y)
	pts[1].X = int32(cp2x)
	pts[1].Y = int32(cp2y)
	pts[2].X = int32(x)
	pts[2].Y = int32(y)
	win32.PolyBezierTo(c.hdc, pts[:])
}

func (c *context) Rect(rect geom.Rect) {
	c.BeginPath()
	c.MoveTo(rect.X, rect.Y)
	c.LineTo(rect.Right(), rect.Y)
	c.LineTo(rect.Right(), rect.Bottom())
	c.LineTo(rect.X, rect.Bottom())
	c.ClosePath()
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
	win32.CloseFigure(c.hdc)
}

func (c *context) disposeBrush() {
	if c.brush != 0 {
		win32.DeleteObject(win32.HGDIOBJ(c.brush))
		c.brush = 0
	}
}

func (c *context) disposePen() {
	if c.pen != 0 {
		win32.DeleteObject(win32.HGDIOBJ(c.pen))
		c.pen = 0
	}
}

func (c *context) Dispose() {
	c.disposeBrush()
	c.disposePen()
}

func fromColorToWin32ColorRef(c Color) win32.COLORREF {
	// RAW: Nuke alpha channel... which isn't acceptable
	return win32.COLORREF(c >> 8)
}
