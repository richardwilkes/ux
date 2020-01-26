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
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
	"github.com/richardwilkes/win32/d2d"
)

var _ Pather = &winPath{}

type winPath struct {
	gc              *context
	pathGeometry    *d2d.PathGeometry
	sink            *d2d.GeometrySink
	result          []*d2d.Geometry
	windingFillMode bool
	isLine          bool
	figureOpen      bool
	hasContent      bool
}

func newWinPath(gc *context, windingFillMode, isLine bool) (*winPath, error) {
	p := &winPath{
		gc:              gc,
		windingFillMode: windingFillMode,
		isLine:          isLine,
	}
	if err := p.setup(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *winPath) setup() error {
	if p.pathGeometry = p.gc.renderTarget.Factory().CreatePathGeometry(); p.pathGeometry == nil {
		return errs.New("unable to create path geometry")
	}
	if p.sink = p.pathGeometry.Open(); p.sink == nil {
		p.pathGeometry.Release()
		p.pathGeometry = nil
		return errs.New("unable to create sink")
	}
	p.sink.SetFillMode(p.windingFillMode)
	return nil
}

func (p *winPath) BeginPath() {
	if p.figureOpen {
		p.sink.EndFigure(false)
	}
	p.sink.BeginFigure(d2d.Point{}, p.isLine)
	p.figureOpen = true
}

func (p *winPath) MoveTo(x, y float64) {
	if p.figureOpen {
		p.sink.EndFigure(false)
	}
	p.sink.BeginFigure(d2d.Point{
		X: float32(x),
		Y: float32(y),
	}, p.isLine)
	p.figureOpen = true
}

func (p *winPath) LineTo(x, y float64) {
	p.sink.AddLine(d2d.Point{
		X: float32(x),
		Y: float32(y),
	})
	p.hasContent = true
}

func (p *winPath) QuadCurveTo(cpx, cpy, x, y float64) {
	p.sink.AddQuadraticBezier(d2d.QuadraticBezierSegment{
		Point1: d2d.Point{
			X: float32(cpx),
			Y: float32(cpy),
		},
		Point2: d2d.Point{
			X: float32(x),
			Y: float32(y),
		},
	})
	p.hasContent = true
}

func (p *winPath) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	p.sink.AddBezier(d2d.BezierSegment{
		Point1: d2d.Point{
			X: float32(cp1x),
			Y: float32(cp1y),
		},
		Point2: d2d.Point{
			X: float32(cp2x),
			Y: float32(cp2y),
		},
		Point3: d2d.Point{
			X: float32(x),
			Y: float32(y),
		},
	})
	p.hasContent = true
}

func (p *winPath) Rect(rect geom.Rect) {
	rg := p.gc.renderTarget.Factory().CreateRectangleGeometry(d2d.Rect{
		Left:   float32(rect.X),
		Top:    float32(rect.Y),
		Right:  float32(rect.Right()),
		Bottom: float32(rect.Bottom()),
	})
	if rg == nil {
		jot.Error(errs.New("unable to create rectangle geometry"))
		return
	}
	p.addOtherGeometry(&rg.Geometry)
}

func (p *winPath) RoundedRect(rect geom.Rect, cornerRadius float64) {
	rg := p.gc.renderTarget.Factory().CreateRoundedRectangleGeometry(d2d.RoundedRect{
		Rect: d2d.Rect{
			Left:   float32(rect.X),
			Top:    float32(rect.Y),
			Right:  float32(rect.Right()),
			Bottom: float32(rect.Bottom()),
		},
		RadiusX: float32(cornerRadius),
		RadiusY: float32(cornerRadius),
	})
	if rg == nil {
		jot.Error(errs.New("unable to create rounded rect geometry"))
		return
	}
	p.addOtherGeometry(&rg.Geometry)
}

func (p *winPath) Ellipse(rect geom.Rect) {
	eg := p.gc.renderTarget.Factory().CreateEllipseGeometry(d2d.Ellipse{
		Point: d2d.Point{
			X: float32(rect.CenterX()),
			Y: float32(rect.CenterY()),
		},
		RadiusX: float32(rect.Width / 2),
		RadiusY: float32(rect.Height / 2),
	})
	if eg == nil {
		jot.Error(errs.New("unable to create ellipse geometry"))
		return
	}
	p.addOtherGeometry(&eg.Geometry)
}

func (p *winPath) Polygon(polygon poly.Polygon) {
	if p.figureOpen {
		p.sink.EndFigure(false)
		p.figureOpen = false
	}
	for _, cont := range polygon {
		for i, pt := range cont {
			if i == 0 {
				p.MoveTo(pt.X, pt.Y)
			} else {
				p.LineTo(pt.X, pt.Y)
			}
		}
		p.ClosePath()
	}
}

func (p *winPath) ClosePath() {
	if p.figureOpen {
		p.sink.EndFigure(true)
		p.figureOpen = false
	}
}

func (p *winPath) addOtherGeometry(geom *d2d.Geometry) {
	if p.figureOpen {
		p.sink.EndFigure(true)
		p.figureOpen = false
	}
	if p.hasContent {
		p.closeSink()
		p.result = append(p.result, &p.pathGeometry.Geometry)
		p.hasContent = false
		if err := p.setup(); err != nil {
			jot.Error(err)
		}
	}
	p.result = append(p.result, geom)
}

func (p *winPath) geometry() *d2d.Geometry {
	if p.figureOpen {
		p.sink.EndFigure(false)
		p.figureOpen = false
	}
	p.closeSink()
	if p.hasContent {
		p.result = append(p.result, &p.pathGeometry.Geometry)
	}
	if len(p.result) == 1 {
		return p.result[0]
	}
	group := p.gc.renderTarget.Factory().CreateGeometryGroup(p.windingFillMode, p.result)
	if group == nil {
		jot.Error(errs.New("unable to create geometry group"))
		return nil
	}
	// release the result here
	return &group.Geometry
}

func (p *winPath) closeSink() {
	if p.sink != nil {
		p.sink.Close()
		p.sink.Release()
		p.sink = nil
	}
}

func (p *winPath) dispose() {
	p.closeSink()
	if p.pathGeometry != nil {
		p.pathGeometry.Release()
		p.pathGeometry = nil
	}
	p.gc = nil
}
