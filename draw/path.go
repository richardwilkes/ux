package draw

import (
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
