package flow

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
)

// Flow holds the flow layout information.
type Flow struct {
	target   layout.Layoutable
	hSpacing float64
	vSpacing float64
}

// New creates a new Flow layout. This layout arranges the children of its
// target left-to-right, then top-to-bottom at their preferred sizes, if
// possible. Each child of the target may have an Alignment set for its
// LayoutData to control vertical positioning within the row. If not present,
// Start is assumed.
func New() *Flow {
	return &Flow{
		hSpacing: layout.DefaultHSpacing,
		vSpacing: layout.DefaultVSpacing,
	}
}

// HSpacing sets the spacing between columns. Defaults to DefaultHSpacing.
func (f *Flow) HSpacing(hSpacing float64) *Flow {
	f.hSpacing = hSpacing
	return f
}

// VSpacing sets the spacing between rows. Defaults to DefaultVSpacing.
func (f *Flow) VSpacing(vSpacing float64) *Flow {
	f.vSpacing = vSpacing
	return f
}

// Apply the layout to the target. A copy is made of this layout and that is
// applied to the target, so this layout may be applied to other targets.
func (f *Flow) Apply(target layout.Layoutable) {
	flow := *f
	flow.target = target
	target.SetLayout(&flow)
}

// Sizes implements Layout.
func (f *Flow) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	if f.hSpacing < 0 {
		f.hSpacing = 0
	}
	if f.vSpacing < 0 {
		f.vSpacing = 0
	}
	var insets geom.Insets
	b := f.target.Border()
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
	for _, child := range f.target.ChildrenForLayout() {
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
				pt.Y += maxHeight + f.vSpacing
				availWidth = width
				availHeight -= maxHeight + f.vSpacing
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
		availWidth -= pref.Width + f.hSpacing
		if availWidth <= 0 {
			pt.X = insets.Left
			pt.Y += maxHeight + f.vSpacing
			availWidth = width
			availHeight -= maxHeight + f.vSpacing
			maxHeight = 0
		} else {
			pt.X += pref.Width + f.hSpacing
		}
	}
	result.Width += insets.Right
	result.Height += insets.Bottom
	largestChildMin.Width += insets.Left + insets.Right
	largestChildMin.Height += insets.Top + insets.Bottom
	return largestChildMin, result, layout.MaxSize(result)
}

// Layout implements Layout.
func (f *Flow) Layout() {
	var insets geom.Insets
	b := f.target.Border()
	if b != nil {
		insets = b.Insets()
	}
	size := f.target.FrameRect().Size
	width := size.Width - (insets.Left + insets.Right)
	pt := geom.Point{X: insets.Left, Y: insets.Top}
	availWidth := width
	availHeight := size.Height - (insets.Top + insets.Bottom)
	var maxHeight float64
	children := f.target.ChildrenForLayout()
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
				pt.Y += maxHeight + f.vSpacing
				availWidth = width
				availHeight -= maxHeight + f.vSpacing
				if i > start {
					f.applyRects(children[start:i], rects[start:i], maxHeight)
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
		availWidth -= pref.Width + f.hSpacing
		if availWidth <= 0 {
			pt.X = insets.Left
			pt.Y += maxHeight + f.vSpacing
			availWidth = width
			availHeight -= maxHeight + f.vSpacing
			f.applyRects(children[start:i+1], rects[start:i+1], maxHeight)
			start = i + 1
			maxHeight = 0
		} else {
			pt.X += pref.Width + f.hSpacing
		}
	}
	if start < len(children) {
		f.applyRects(children[start:], rects[start:], maxHeight)
	}
}

func (f *Flow) applyRects(children []layout.Layoutable, rects []geom.Rect, maxHeight float64) {
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
