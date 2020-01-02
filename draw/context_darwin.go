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
	"github.com/richardwilkes/macos/cf"
	"github.com/richardwilkes/macos/cg"
	"github.com/richardwilkes/macos/ct"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
	"github.com/richardwilkes/ux/draw/linecap"
	"github.com/richardwilkes/ux/draw/linejoin"
	"github.com/richardwilkes/ux/draw/quality"
)

var _ Context = &context{}

// OSContext is the platform-specific drawing context on macOS.
type OSContext = cg.Context

type context struct {
	gc OSContext
}

func osNewContextForOSContext(gc OSContext) Context {
	return &context{gc: gc}
}

func (c *context) OSContext() OSContext {
	return c.gc
}

func (c *context) Save() {
	c.gc.SaveGState()
}

func (c *context) Restore() {
	c.gc.RestoreGState()
}

func (c *context) SetOpacity(opacity float64) {
	c.gc.SetAlpha(opacity)
}

func (c *context) SetPatternOffset(x, y float64) {
	c.gc.SetPatternPhase(x, y)
}

func (c *context) SetStrokeWidth(width float64) {
	c.gc.SetLineWidth(width)
}

func (c *context) SetLineCap(lineCap linecap.LineCap) {
	c.gc.SetLineCap(cg.LineCap(lineCap))
}

func (c *context) SetLineJoin(lineJoin linejoin.LineJoin) {
	c.gc.SetLineJoin(cg.LineJoin(lineJoin))
}

func (c *context) SetMiterLimit(limit float64) {
	c.gc.SetMiterLimit(limit)
}

func (c *context) SetLineDash(phase float64, segments ...float64) {
	c.gc.SetLineDash(phase, segments...)
}

func (c *context) SetInterpolationQualityHint(q quality.Quality) {
	c.gc.SetInterpolationQuality(cg.InterpolationQuality(q))
}

func (c *context) Translate(x, y float64) {
	c.gc.TranslateCTM(x, y)
}

func (c *context) Scale(x, y float64) {
	c.gc.ScaleCTM(x, y)
}

func (c *context) Rotate(angleInRadians float64) {
	c.gc.RotateCTM(angleInRadians)
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
	x, y, width, height := c.gc.GetClipBoundingBox()
	return geom.Rect{
		Point: geom.Point{
			X: x,
			Y: y,
		},
		Size: geom.Size{
			Width:  width,
			Height: height,
		},
	}
}

func (c *context) Clip() {
	c.gc.Clip()
}

func (c *context) ClipEvenOdd() {
	c.gc.EOClip()
}

func (c *context) DrawString(x, y float64, font *Font, ink Ink, str string) {
	if str != "" {
		c.gc.SaveGState()
		as := cf.AttributedStringCreateMutable(0)
		as.BeginEditing()
		s := cf.StringCreateWithString(str)
		length := s.GetLength()
		as.ReplaceString(0, 0, s)
		as.SetAttribute(0, length, ct.FontAttributeName, cf.Type(font.ref))
		as.SetAttribute(0, length, ct.ForegroundColorFromContextAttributeName, cf.Type(cf.BooleanTrue))
		as.EndEditing()
		line := ct.LineCreateWithAttributedString(cf.AttributedString(as))
		c.gc.SetTextMatrix(cg.AffineTransformIdentity)
		c.gc.TranslateCTM(0, font.ascent)
		c.gc.ScaleCTM(1, -1)
		c.gc.TranslateCTM(x, -y)
		ink.osPrepareForFill(c)
		line.Draw(c.gc)
		line.Release()
		as.Release()
		c.gc.RestoreGState()
	}
}

func (c *context) BeginPath() {
	c.gc.BeginPath()
}

func (c *context) MoveTo(x, y float64) {
	c.gc.MoveToPoint(x, y)
}

func (c *context) LineTo(x, y float64) {
	c.gc.AddLineToPoint(x, y)
}

func (c *context) QuadCurveTo(cpx, cpy, x, y float64) {
	c.gc.AddQuadCurveToPoint(cpx, cpy, x, y)
}

func (c *context) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	c.gc.AddCurveToPoint(cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c *context) Rect(rect geom.Rect) {
	c.gc.AddRect(rect.X, rect.Y, rect.Width, rect.Height)
}

func (c *context) RoundedRect(rect geom.Rect, cornerRadius float64) {
	c.gc.MoveToPoint(rect.X, rect.Y+cornerRadius)
	c.gc.AddQuadCurveToPoint(rect.X, rect.Y, rect.X+cornerRadius, rect.Y)
	c.gc.AddLineToPoint(rect.X+rect.Width-cornerRadius, rect.Y)
	c.gc.AddQuadCurveToPoint(rect.X+rect.Width, rect.Y, rect.X+rect.Width, rect.Y+cornerRadius)
	c.gc.AddLineToPoint(rect.X+rect.Width, rect.Y+rect.Height-cornerRadius)
	c.gc.AddQuadCurveToPoint(rect.X+rect.Width, rect.Y+rect.Height, rect.X+rect.Width-cornerRadius, rect.Y+rect.Height)
	c.gc.AddLineToPoint(rect.X+cornerRadius, rect.Y+rect.Height)
	c.gc.AddQuadCurveToPoint(rect.X, rect.Y+rect.Height, rect.X, rect.Y+rect.Height-cornerRadius)
	c.gc.ClosePath()
}

func (c *context) Ellipse(rect geom.Rect) {
	c.gc.AddEllipseInRect(rect.X, rect.Y, rect.Width, rect.Height)
}

func (c *context) Polygon(polygon poly.Polygon) {
	for _, cont := range polygon {
		for i, pt := range cont {
			if i == 0 {
				c.gc.MoveToPoint(pt.X, pt.Y)
			} else {
				c.gc.AddLineToPoint(pt.X, pt.Y)
			}
		}
	}
	c.gc.ClosePath()
}

func (c *context) ClosePath() {
	c.gc.ClosePath()
}

func (c *context) Dispose() {
	c.gc.Release()
}
