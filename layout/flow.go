package layout

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout/align"
	"math"
)

type flow struct {
	target   Layoutable
	hSpacing float64
	vSpacing float64
}

// NewFlow creates and assigns a layout on the specified target. This layout
// arranges the children of its target left-to-right, then top-to-bottom at
// their preferred sizes, if possible. hSpacing defines the spacing between
// columns. vSpacing defines the spacing between rows.
//
// Each child may have an Alignment set for its LayoutData to control vertical
// positioning within the row. If not present, Start is assumed.
func NewFlow(target Layoutable, hSpacing, vSpacing float64) {
	target.SetLayout(&flow{
		target:   target,
		hSpacing: hSpacing,
		vSpacing: vSpacing,
	})
}

func (l *flow) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	if l.hSpacing < 0 {
		l.hSpacing = 0
	}
	if l.vSpacing < 0 {
		l.vSpacing = 0
	}
	var insets geom.Insets
	b := l.target.Border()
	if b != nil {
		insets = b.Insets()
	}
	if hint.Width < 1 {
		hint.Width = math.MaxFloat32
	}
	if hint.Height < 1 {
		hint.Height = math.MaxFloat32
	}
	width := hint.Width - (insets.Left + insets.Right)
	pt := geom.Point{X: insets.Left, Y: insets.Top}
	result := geom.Size{Width: pt.Y, Height: pt.Y}
	availWidth := width
	availHeight := hint.Height - (insets.Top + insets.Bottom)
	var maxHeight float64
	var largestChildMin geom.Size
	for _, child := range l.target.ChildrenForLayout() {
		min, pref, _ := child.Sizes(geom.Size{})
		if largestChildMin.Width < min.Width {
			largestChildMin.Width = min.Width
		}
		if largestChildMin.Height < min.Height {
			largestChildMin.Height = min.Height
		}
		if pref.Width > availWidth {
			switch {
			case min.Width <= availWidth:
				pref.Width = availWidth
			case pt.X == insets.Left:
				pref.Width = min.Width
			default:
				pt.X = insets.Left
				pt.Y += maxHeight + l.vSpacing
				availWidth = width
				availHeight -= maxHeight + l.vSpacing
				maxHeight = 0
				if pref.Width > availWidth {
					if min.Width <= availWidth {
						pref.Width = availWidth
					} else {
						pref.Width = min.Width
					}
				}
			}
			savedWidth := pref.Width
			min, pref, _ = child.Sizes(geom.Size{Width: pref.Width})
			pref.Width = savedWidth
			if pref.Height > availHeight {
				if min.Height <= availHeight {
					pref.Height = availHeight
				} else {
					pref.Height = min.Height
				}
			}
		}
		extent := pt.X + pref.Width
		if result.Width < extent {
			result.Width = extent
		}
		extent = pt.Y + pref.Height
		if result.Height < extent {
			result.Height = extent
		}
		if maxHeight < pref.Height {
			maxHeight = pref.Height
		}
		availWidth -= pref.Width + l.hSpacing
		if availWidth <= 0 {
			pt.X = insets.Left
			pt.Y += maxHeight + l.vSpacing
			availWidth = width
			availHeight -= maxHeight + l.vSpacing
			maxHeight = 0
		} else {
			pt.X += pref.Width + l.hSpacing
		}
	}
	result.Width += insets.Right
	result.Height += insets.Bottom
	largestChildMin.Width += insets.Left + insets.Right
	largestChildMin.Height += insets.Top + insets.Bottom
	return largestChildMin, result, MaxSize(result)
}

func (l *flow) Layout() {
	var insets geom.Insets
	b := l.target.Border()
	if b != nil {
		insets = b.Insets()
	}
	size := l.target.FrameRect().Size
	width := size.Width - (insets.Left + insets.Right)
	pt := geom.Point{X: insets.Left, Y: insets.Top}
	availWidth := width
	availHeight := size.Height - (insets.Top + insets.Bottom)
	var maxHeight float64
	children := l.target.ChildrenForLayout()
	rects := make([]geom.Rect, len(children))
	start := 0
	for i, child := range children {
		min, pref, _ := child.Sizes(geom.Size{})
		if pref.Width > availWidth {
			switch {
			case min.Width <= availWidth:
				pref.Width = availWidth
			case pt.X == insets.Left:
				pref.Width = min.Width
			default:
				pt.X = insets.Left
				pt.Y += maxHeight + l.vSpacing
				availWidth = width
				availHeight -= maxHeight + l.vSpacing
				if i > start {
					l.applyRects(children[start:i], rects[start:i], maxHeight)
					start = i
				}
				maxHeight = 0
				if pref.Width > availWidth {
					if min.Width <= availWidth {
						pref.Width = availWidth
					} else {
						pref.Width = min.Width
					}
				}
			}
			savedWidth := pref.Width
			min, pref, _ = child.Sizes(geom.Size{Width: pref.Width})
			pref.Width = savedWidth
			if pref.Height > availHeight {
				if min.Height <= availHeight {
					pref.Height = availHeight
				} else {
					pref.Height = min.Height
				}
			}
		}
		rects[i] = geom.Rect{Point: pt, Size: pref}
		if maxHeight < pref.Height {
			maxHeight = pref.Height
		}
		availWidth -= pref.Width + l.hSpacing
		if availWidth <= 0 {
			pt.X = insets.Left
			pt.Y += maxHeight + l.vSpacing
			availWidth = width
			availHeight -= maxHeight + l.vSpacing
			l.applyRects(children[start:i+1], rects[start:i+1], maxHeight)
			start = i + 1
			maxHeight = 0
		} else {
			pt.X += pref.Width + l.hSpacing
		}
	}
	if start < len(children) {
		l.applyRects(children[start:], rects[start:], maxHeight)
	}
}

func (l *flow) applyRects(children []Layoutable, rects []geom.Rect, maxHeight float64) {
	for i, child := range children {
		vAlign, ok := child.LayoutData().(align.Alignment)
		if !ok {
			vAlign = align.Start
		}
		switch vAlign {
		case align.Middle:
			if rects[i].Height < maxHeight {
				rects[i].Y += (maxHeight - rects[i].Height) / 2
			}
		case align.End:
			rects[i].Y += maxHeight - rects[i].Height
		case align.Fill:
			rects[i].Height = maxHeight
		default: // same as Start
		}
		child.SetFrameRect(rects[i])
	}
}
