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
	"math"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
	"github.com/richardwilkes/ux/draw/linecap"
	"github.com/richardwilkes/ux/draw/linejoin"
	"github.com/richardwilkes/ux/draw/quality"
	"github.com/richardwilkes/win32/d2d"
)

// RAW: Given the huge graphical issues I've seen in the early implementation
// using GDI, it's not going to be even remotely acceptable. Will have to
// look at Direct2D instead... but that will require writing a shim layer, as
// there is no C API for it, which is required for Go to make the calls.

// OSContext is the platform-specific drawing context on Windows.
type OSContext = *d2d.HWNDRenderTarget

type context struct {
	renderTarget *OSContext
	stack        []*contextState
	path         Path
}

func osNewContextForOSContext(gc *OSContext) Context {
	c := &context{
		renderTarget: gc,
		stack:        []*contextState{{strokeWidth: 1}},
	}
	rt := c.OSContext()
	rt.BeginDraw()
	c.osiSetMatrix(xmath.NewIdentityMatrix2D())
	rt.PushAxisAlignedClip(d2d.Rect{
		Left:   -math.MaxFloat32 / 2,
		Top:    -math.MaxFloat32 / 2,
		Right:  math.MaxFloat32,
		Bottom: math.MaxFloat32,
	}, false)
	return c
}

func (c *context) OSContext() OSContext {
	return *c.renderTarget
}

func (c *context) current() *contextState {
	return c.stack[len(c.stack)-1]
}

func (c *context) Save() {
	current := c.current()
	current.state = c.OSContext().SaveDrawingState()
	c.stack = append(c.stack, current.copy(c))
}

func (c *context) Restore() {
	if len(c.stack) > 0 {
		c.current().dispose()
		c.stack[len(c.stack)-1] = nil
		c.stack = c.stack[:len(c.stack)-1]
		current := c.current()
		c.OSContext().RestoreDrawingState(current.state)
		current.state = nil
		c.path = *current.clip.Clone()
		c.pushClip(current.clipWindingFillMode)
	}
}

func (c *context) SetOpacity(opacity float64) {
	// RAW: Implement
}

func (c *context) SetPatternOffset(x, y float64) {
	// RAW: Implement
}

func (c *context) SetStrokeWidth(width float64) {
	c.current().strokeWidth = float32(width)
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
	c.osiSetMatrix(xmath.NewTranslationMatrix2D(x, y))
}

func (c *context) Scale(x, y float64) {
	c.osiSetMatrix(xmath.NewScaleMatrix2D(x, y))
}

func (c *context) Rotate(angleInRadians float64) {
	c.osiSetMatrix(xmath.NewRotationMatrix2D(angleInRadians))
}

func (c *context) osiSetMatrix(matrix *xmath.Matrix2D) {
	c.OSContext().SetTransform(&d2d.Matrix3x2{
		A11: float32(matrix.XX),
		A12: float32(matrix.YX),
		A21: float32(matrix.XY),
		A22: float32(matrix.YY),
		A31: float32(matrix.X0),
		A32: float32(matrix.Y0),
	})
}

func (c *context) Fill(ink Ink) {
	ink.osFill(c)
	c.path.BeginPath()
}

func (c *context) FillEvenOdd(ink Ink) {
	ink.osFillEvenOdd(c)
	c.path.BeginPath()
}

func (c *context) Stroke(ink Ink) {
	ink.osStroke(c)
	c.path.BeginPath()
}

func (c *context) GetClipRect() geom.Rect {
	current := c.current()
	switch len(current.clip.nodes) {
	case 0:
		return geom.Rect{
			Point: geom.Point{
				X: -math.MaxFloat32,
				Y: -math.MaxFloat32,
			},
			Size: geom.Size{
				Width:  math.MaxFloat32,
				Height: math.MaxFloat32,
			},
		}
	case 1:
		if rpn, ok := current.clip.nodes[0].(*rectPathNode); ok {
			return rpn.rect
		}
		fallthrough
	default:
		return current.clip.Bounds()
	}
}

func (c *context) Clip() {
	c.pushClip(true)
}

func (c *context) ClipEvenOdd() {
	c.pushClip(false)
}

func (c *context) pushClip(windingFillMode bool) {
	c.popClip()
	current := c.current()
	current.clipWindingFillMode = windingFillMode
	current.clip.BeginPath()
	defer c.path.BeginPath()
	switch len(c.path.nodes) {
	case 0:
		current.clip.Rect(geom.Rect{})
		c.OSContext().PushAxisAlignedClip(d2d.Rect{}, false)
		return
	case 1:
		if rpn, ok := c.path.nodes[0].(*rectPathNode); ok {
			current.clip.Rect(rpn.rect)
			c.OSContext().PushAxisAlignedClip(d2d.Rect{
				Left:   float32(rpn.rect.X),
				Top:    float32(rpn.rect.Y),
				Right:  float32(rpn.rect.Right()),
				Bottom: float32(rpn.rect.Bottom()),
			}, false)
			return
		}
		fallthrough
	default:
		p, err := newWinPath(c, windingFillMode, false)
		if err != nil {
			jot.Error(err)
			return
		}
		c.path.SendPath(p)
		c.path.SendPath(&current.clip)
		p.gc.OSContext().PushLayer(&d2d.LayerParameters{GeometricMask: p.geometry(),}, nil)
		p.dispose()
	}
}

func (c *context) popClip() {
	current := c.current()
	switch len(current.clip.nodes) {
	case 0:
		c.OSContext().PopAxisAlignedClip()
	case 1:
		if _, ok := current.clip.nodes[0].(*rectPathNode); ok {
			c.OSContext().PopAxisAlignedClip()
			return
		}
		fallthrough
	default:
		c.OSContext().PopLayer()
	}
}

func (c *context) DrawImage(img *Image, where geom.Point) {
	// RAW: Implement
}

func (c *context) DrawImageInRect(img *Image, rect geom.Rect) {
	// RAW: Implement
}

func (c *context) DrawString(x, y float64, font *Font, ink Ink, str string) {
	// RAW: Use ink
	// win32.SelectObject(c.hdc, win32.HGDIOBJ(font.ref))
	// win32.TextOut(c.hdc, int(x), int(y), str)
}

func (c *context) BeginPath() {
	c.path.BeginPath()
}

func (c *context) MoveTo(x, y float64) {
	c.path.MoveTo(x, y)
}

func (c *context) LineTo(x, y float64) {
	c.path.LineTo(x, y)
}

func (c *context) QuadCurveTo(cpx, cpy, x, y float64) {
	c.path.QuadCurveTo(cpx, cpy, x, y)
}

func (c *context) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	c.path.CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c *context) Rect(rect geom.Rect) {
	c.path.Rect(rect)
}

func (c *context) RoundedRect(rect geom.Rect, cornerRadius float64) {
	c.path.RoundedRect(rect, cornerRadius)
}

func (c *context) Ellipse(rect geom.Rect) {
	c.path.Ellipse(rect)
}

func (c *context) Polygon(polygon poly.Polygon) {
	c.path.Polygon(polygon)
}

func (c *context) ClosePath() {
	c.path.ClosePath()
}

func (c *context) Dispose() {
	c.popClip()
	c.path.BeginPath()
	for _, one := range c.stack {
		one.dispose()
	}
	if t1, t2 := c.OSContext().EndDraw(); t1 != 0 || t2 != 0 {
		c.OSContext().Release()
		*c.renderTarget = nil
	}
}
