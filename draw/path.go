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

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/toolbox/xmath/geom/poly"
)

var (
	_ Pather     = &Path{}
	_ PathSender = &Path{}
)

// Pather defines the methods necessary for describing of a series of shapes
// or lines.
type Pather interface {
	// BeginPath discards any existing path data and starts a fresh path.
	BeginPath()
	// MoveTo begins a new sub-path at the specified point.
	MoveTo(x, y float64)
	// LineTo appends a straight line segment from the current point to the
	// specified point.
	LineTo(x, y float64)
	// QuadCurveTo appends a quadratic Bezier curve from the current point,
	// using a control point and an end point.
	QuadCurveTo(cpx, cpy, x, y float64)
	// CubicCurveTo appends a cubic Bezier curve from the current point, using
	// the provided controls points and end point.
	CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64)
	// Rect adds a rectangle to the path. The rectangle is a complete
	// sub-path, i.e. it starts with a MoveTo and ends with a ClosePath
	// operation.
	Rect(rect geom.Rect)
	// RoundedRect adds a rounded rectangle to the path. The rectangle is a
	// complete sub-path, i.e. it starts with a MoveTo and ends with a
	// ClosePath operation.
	RoundedRect(rect geom.Rect, cornerRadius float64)
	// Ellipse adds an ellipse to the path. The ellipse is a complete
	// sub-path, i.e. it starts with a MoveTo and ends with a ClosePath
	// operation.
	Ellipse(rect geom.Rect)
	// Polygon add a polygon to the path. The ellipse is a complete sub-path,
	// i.e. it starts with a MoveTo and ends with a ClosePath operation.
	Polygon(polygon poly.Polygon)
	// ClosePath closes and terminates the current sub-path.
	ClosePath()
}

// PathSender sends the contents of a path into the target Pather.
type PathSender interface {
	SendPath(pathBuilder Pather)
}

// Path is a description of a series of shapes or lines.
type Path struct {
	nodes []PathSender
}

// Append the other path to the end of this path.
func (p *Path) Append(other *Path) {
	p.nodes = append(p.nodes, other.nodes...)
}

// Clone creates an exact copy of this path.
func (p *Path) Clone() *Path {
	other := &Path{nodes: make([]PathSender, len(p.nodes))}
	copy(other.nodes, p.nodes)
	return other
}

// BeginPath discards any existing path data and starts a fresh path.
func (p *Path) BeginPath() {
	p.nodes = nil
}

// MoveTo begins a new sub-path at the specified point.
func (p *Path) MoveTo(x, y float64) {
	p.nodes = append(p.nodes, &moveToPathNode{x: x, y: y})
}

// LineTo appends a straight line segment from the current point to the
// specified point.
func (p *Path) LineTo(x, y float64) {
	p.nodes = append(p.nodes, &lineToPathNode{x: x, y: y})
}

// QuadCurveTo appends a quadratic Bezier curve from the current point, using
// a control point and an end point.
func (p *Path) QuadCurveTo(cpx, cpy, x, y float64) {
	p.nodes = append(p.nodes, &quadCurveToPathNode{cpx: cpx, cpy: cpy, x: x, y: y})
}

// CubicCurveTo appends a cubic Bezier curve from the current point, using the
// provided controls points and end point.
func (p *Path) CubicCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	p.nodes = append(p.nodes, &cubicCurveToPathNode{cp1x: cp1x, cp1y: cp1y, cp2x: cp2x, cp2y: cp2y, x: x, y: y})
}

// Rect adds a rectangle to the path. The rectangle is a complete sub-path,
// i.e. it starts with a MoveTo and ends with a ClosePath operation.
func (p *Path) Rect(rect geom.Rect) {
	p.nodes = append(p.nodes, &rectPathNode{rect: rect})
}

// RoundedRect adds a rounded rectangle to the path. The rectangle is a
// complete sub-path, i.e. it starts with a MoveTo and ends with a ClosePath
// operation.
func (p *Path) RoundedRect(rect geom.Rect, cornerRadius float64) {
	p.nodes = append(p.nodes, &roundedRectPathNode{rect: rect, cornerRadius: cornerRadius})
}

// Ellipse adds an ellipse to the path. The ellipse is a complete sub-path,
// i.e. it starts with a MoveTo and ends with a ClosePath operation.
func (p *Path) Ellipse(rect geom.Rect) {
	p.nodes = append(p.nodes, &ellipsePathNode{rect: rect})
}

// Polygon add a polygon to the path. The ellipse is a complete sub-path,
// i.e. it starts with a MoveTo and ends with a ClosePath operation.
func (p *Path) Polygon(polygon poly.Polygon) {
	p.nodes = append(p.nodes, &polygonPathNode{polygon: polygon})
}

// ClosePath closes and terminates the current sub-path.
func (p *Path) ClosePath() {
	p.nodes = append(p.nodes, &closePathNode{})
}

// SendPath sends the contents of the path into the target Pather.
func (p *Path) SendPath(pather Pather) {
	for _, n := range p.nodes {
		n.SendPath(pather)
	}
}

// Bounds returns the bounding rectangle for this path.
func (p *Path) Bounds() geom.Rect {
	x1 := -math.MaxFloat64 / 2
	y1 := x1
	x2 := math.MaxFloat64 / 2
	y2 := x2
	lastX := x1
	lastY := y1
	for _, n := range p.nodes {
		switch t := n.(type) {
		case *moveToPathNode:
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.x, t.y)
			lastX = t.x
			lastY = t.y
		case *lineToPathNode:
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.x, t.y)
			lastX = t.x
			lastY = t.y
		case *rectPathNode:
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.X, t.rect.Y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.Right(), t.rect.Bottom())
			lastX = t.rect.X
			lastY = t.rect.Y
		case *roundedRectPathNode:
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.X, t.rect.Y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.Right(), t.rect.Bottom())
			lastX = t.rect.X
			lastY = t.rect.Y
		case *ellipsePathNode:
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.X, t.rect.Y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, t.rect.Right(), t.rect.Bottom())
			lastX = t.rect.X
			lastY = t.rect.Y
		case *polygonPathNode:
			r := t.polygon.Bounds()
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, r.X, r.Y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, r.Right(), r.Bottom())
			if len(t.polygon) > 0 {
				c := t.polygon[len(t.polygon)-1]
				if len(c) > 0 {
					p := c[len(c)-1]
					lastX = p.X
					lastY = p.Y
				}
			}
		case *quadCurveToPathNode:
			minX, minY, maxX, maxY := bezierBounds(lastX, lastY, t.cpx, t.cpy, t.cpx, t.cpy, t.x, t.y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, minX, minY)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, maxX, maxY)
			lastX = t.x
			lastY = t.y
		case *cubicCurveToPathNode:
			minX, minY, maxX, maxY := bezierBounds(lastX, lastY, t.cp1x, t.cp1y, t.cp2x, t.cp2y, t.x, t.y)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, minX, minY)
			adjustBoundsForPoint(&x1, &y1, &x2, &y2, maxX, maxY)
			lastX = t.x
			lastY = t.y
		default:
		}
	}
	return geom.Rect{
		Point: geom.Point{
			X: x1,
			Y: y1,
		},
		Size: geom.Size{
			Width:  x2 - x1,
			Height: y2 - y1,
		},
	}
}

const zeroTolerance = 1e-12

func bezierBounds(sx, sy, cp1x, cp1y, cp2x, cp2y, ex, ey float64) (minX, minY, maxX, maxY float64) {
	tvalues := make([]float64, 0, 4)
	for i := 0; i < 2; i++ {
		var a, b, c float64
		if i == 0 {
			b = 6*sx - 12*cp1x + 6*cp2x
			a = -3*sx + 9*cp1x - 9*cp2x + 3*ex
			c = 3*cp1x - 3*sx
		} else {
			b = 6*sy - 12*cp1y + 6*cp2y
			a = -3*sy + 9*cp1y - 9*cp2y + 3*ey
			c = 3*cp1y - 3*sy
		}
		if math.Abs(a) < zeroTolerance {
			if math.Abs(b) < zeroTolerance {
				continue
			}
			t := -c / b
			if t > 0 && t < 1 {
				tvalues = append(tvalues, t)
			}
			continue
		}
		b2ac := b*b - 4*c*a
		if b2ac < 0 {
			if math.Abs(b2ac) < zeroTolerance {
				t := -b / (2 * a)
				if t > 0 && t < 1 {
					tvalues = append(tvalues, t)
				}
			}
			continue
		}
		sqrtb2ac := math.Sqrt(b2ac)
		t := (-b + sqrtb2ac) / (2 * a)
		if t > 0 && t < 1 {
			tvalues = append(tvalues, t)
		}
		t = (-b - sqrtb2ac) / (2 * a)
		if t > 0 && t < 1 {
			tvalues = append(tvalues, t)
		}
	}
	xvalues := make([]float64, len(tvalues))
	yvalues := make([]float64, len(tvalues))
	for i, t := range tvalues {
		mt := 1 - t
		mt2 := mt * mt
		mt3 := mt2 * mt
		t2 := t * t
		t3 := t2 * t
		xvalues[i] = (mt3 * sx) + (3 * mt2 * t * cp1x) + (3 * mt * t2 * cp2x) + (t3 * ex)
		yvalues[i] = (mt3 * sy) + (3 * mt2 * t * cp1y) + (3 * mt * t2 * cp2y) + (t3 * ey)
	}
	minX = math.Min(sx, ex)
	maxX = math.Max(sx, ex)
	for _, x := range xvalues {
		if minX > x {
			minX = x
		}
		if maxX < x {
			maxX = x
		}
	}
	minY = math.Min(sy, ey)
	maxY = math.Max(sy, ey)
	for _, y := range yvalues {
		if minY > y {
			minY = y
		}
		if maxY < y {
			maxY = y
		}
	}
	return
}

func adjustBoundsForPoint(x1, y1, x2, y2 *float64, x, y float64) {
	if *x1 > x {
		*x1 = x
	}
	if *x2 < x {
		*x2 = x
	}
	if *y1 > y {
		*y1 = y
	}
	if *y2 < y {
		*y2 = y
	}
}

type moveToPathNode struct {
	x, y float64
}

func (n *moveToPathNode) SendPath(p Pather) {
	p.MoveTo(n.x, n.y)
}

type lineToPathNode struct {
	x, y float64
}

func (n *lineToPathNode) SendPath(p Pather) {
	p.LineTo(n.x, n.y)
}

type rectPathNode struct {
	rect geom.Rect
}

func (n *rectPathNode) SendPath(p Pather) {
	p.Rect(n.rect)
}

type roundedRectPathNode struct {
	rect         geom.Rect
	cornerRadius float64
}

func (n *roundedRectPathNode) SendPath(p Pather) {
	p.RoundedRect(n.rect, n.cornerRadius)
}

type ellipsePathNode struct {
	rect geom.Rect
}

func (n *ellipsePathNode) SendPath(p Pather) {
	p.Ellipse(n.rect)
}

type polygonPathNode struct {
	polygon poly.Polygon
}

func (n *polygonPathNode) SendPath(p Pather) {
	p.Polygon(n.polygon)
}

type quadCurveToPathNode struct {
	cpx, cpy, x, y float64
}

func (n *quadCurveToPathNode) SendPath(p Pather) {
	p.QuadCurveTo(n.cpx, n.cpy, n.x, n.y)
}

type cubicCurveToPathNode struct {
	cp1x, cp1y, cp2x, cp2y, x, y float64
}

func (n *cubicCurveToPathNode) SendPath(p Pather) {
	p.CubicCurveTo(n.cp1x, n.cp1y, n.cp2x, n.cp2y, n.x, n.y)
}

type closePathNode struct {
}

func (n *closePathNode) SendPath(p Pather) {
	p.ClosePath()
}
