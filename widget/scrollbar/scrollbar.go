package scrollbar

import (
	"math"
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/widget"
)

const (
	none scrollbarPart = iota
	thumb
	lineUp
	lineDown
	pageUp
	pageDown
)

type scrollbarPart int

// ScrollPager objects can provide line and page information for scrolling.
type ScrollPager interface {
	// ScrollAmount is called to determine how far to scroll in the given
	// direction. 'page' is true if the desire is to reveal a full page of
	// content, or false to reveal a full line of content. A positive value
	// should be returned regardless of the direction, although negative
	// values will behave as if they were positive.
	ScrollAmount(horizontal, towardsStart, page bool) float64
}

// Scrollable objects can respond to ScrollBars.
type Scrollable interface {
	ScrollPager
	// ScrolledPosition is called to determine the current position of the
	// Scrollable.
	ScrolledPosition(horizontal bool) float64
	// SetScrolledPosition is called to set the current position of the
	// Scrollable.
	SetScrolledPosition(horizontal bool, position float64)
	// VisibleSize is called to determine the size of the visible portion of
	// the Scrollable.
	VisibleSize(horizontal bool) float64
	// ContentSize is called to determine the total size of the Scrollable.
	ContentSize(horizontal bool) float64
}

// ScrollBar is a widget that controls scrolling.
type ScrollBar struct {
	ux.Panel
	managed
	target     Scrollable
	thumbDown  float64
	sequence   int
	pressed    scrollbarPart
	horizontal bool
}

// NewHorizontal creates a new horizontal scrollbar.
func NewHorizontal() *ScrollBar {
	return New(true)
}

// NewVertical creates a new vertical scrollbar.
func NewVertical() *ScrollBar {
	return New(false)
}

// New creates a new scrollbar.
func New(horizontal bool) *ScrollBar {
	s := &ScrollBar{horizontal: horizontal}
	s.managed.initialize()
	s.InitTypeAndID(s)
	s.SetSizer(s.DefaultSizes)
	s.DrawCallback = s.DefaultDraw
	s.MouseDownCallback = s.DefaultMouseDown
	s.MouseDragCallback = s.DefaultMouseDrag
	s.MouseUpCallback = s.DefaultMouseUp
	return s
}

// Target returns the scrollable target. May be nil.
func (s *ScrollBar) Target() Scrollable {
	return s.target
}

// SetTarget sets the scrollable target. May be nil.
func (s *ScrollBar) SetTarget(target Scrollable) *ScrollBar {
	if s.target != target {
		s.target = target
		s.MarkForRedraw()
	}
	return s
}

// DefaultSizes provides the default sizing.
func (s *ScrollBar) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	if s.horizontal {
		pref.Width = s.barSize * 2
		pref.Height = s.barSize
		max.Width = layout.DefaultMaxSize
		max.Height = s.barSize
	} else {
		pref.Width = s.barSize
		pref.Height = s.barSize * 2
		max.Width = s.barSize
		max.Height = layout.DefaultMaxSize
	}
	if border := s.Border(); border != nil {
		insets := border.Insets()
		pref.AddInsets(insets)
		max.AddInsets(insets)
	}
	return pref, pref, max
}

// DefaultDraw provides the default drawing.
func (s *ScrollBar) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	rect := s.ContentRect(false)
	if s.horizontal {
		rect.Height = s.barSize
	} else {
		rect.Width = s.barSize
	}
	widget.DrawRectBase(gc, rect, s.currentBackgroundInk(none), s.edgeInk)
	s.drawLineButton(gc, lineDown)
	if s.pressed == pageUp || s.pressed == pageDown {
		rect = s.partRect(s.pressed)
		if !rect.IsEmpty() {
			if s.horizontal {
				rect.Y++
				rect.Height -= 2
			} else {
				rect.X++
				rect.Width -= 2
			}
			gc.Rect(rect)
			gc.Fill(s.currentBackgroundInk(s.pressed))
		}
	}
	s.drawThumb(gc)
	s.drawLineButton(gc, lineUp)
}

func (s *ScrollBar) drawLineButton(gc draw.Context, linePart scrollbarPart) {
	rect := s.partRect(linePart)
	widget.DrawRectBase(gc, rect, s.currentBackgroundInk(linePart), s.edgeInk)
	rect.InsetUniform(1.5)
	if s.horizontal {
		triHeight := rect.Width * 0.75
		triWidth := triHeight / 2
		left := rect.X + (rect.Width-triWidth)/2
		right := left + triWidth
		top := rect.Y + (rect.Height-triHeight)/2
		bottom := top + triHeight
		if linePart == lineUp {
			left, right = right, left
		}
		gc.MoveTo(left, top)
		gc.LineTo(left, bottom)
		gc.LineTo(right, top+(bottom-top)/2)
	} else {
		triWidth := rect.Height * 0.75
		triHeight := triWidth / 2
		left := rect.X + (rect.Width-triWidth)/2
		right := left + triWidth
		top := rect.Y + (rect.Height-triHeight)/2
		bottom := top + triHeight
		if linePart == lineUp {
			top, bottom = bottom, top
		}
		gc.MoveTo(left, top)
		gc.LineTo(right, top)
		gc.LineTo(left+(right-left)/2, bottom)
	}
	gc.ClosePath()
	gc.Fill(s.currentMarkColor(linePart))
}

func (s *ScrollBar) drawThumb(gc draw.Context) {
	if rect := s.partRect(thumb); !rect.IsEmpty() {
		widget.DrawRectBase(gc, rect, s.currentBackgroundInk(thumb), s.edgeInk)
		var v0, v1, v2 float64
		if s.horizontal {
			v0 = math.Floor(rect.X + rect.Width/2)
			d := math.Ceil(rect.Height * 0.2)
			v1 = rect.Y + d
			v2 = rect.Y + rect.Height - (d + 1)
		} else {
			v0 = math.Floor(rect.Y + rect.Height/2)
			d := math.Ceil(rect.Width * 0.2)
			v1 = rect.X + d
			v2 = rect.X + rect.Width - (d + 1)
		}
		for i := -1; i < 2; i++ {
			if s.horizontal {
				x := v0 + float64(i*2)
				gc.MoveTo(x, v1)
				gc.LineTo(x, v2)
			} else {
				y := v0 + float64(i*2)
				gc.MoveTo(v1, y)
				gc.LineTo(v2, y)
			}
		}
		gc.Stroke(s.currentMarkColor(thumb))
	}
}

func (s *ScrollBar) currentBackgroundInk(which scrollbarPart) draw.Ink {
	switch {
	case which != none && s.pressed == which:
		return s.pressedBackgroundInk
	case s.Focused():
		return s.focusedBackgroundInk
	default:
		return s.backgroundInk
	}
}

func (s *ScrollBar) currentMarkColor(which scrollbarPart) draw.Ink {
	if s.partEnabled(which) {
		return s.markInk
	}
	return s.disabledMarkInk
}

func (s *ScrollBar) partEnabled(which scrollbarPart) bool {
	if s.Enabled() && s.target != nil {
		switch which {
		case lineUp, pageUp:
			return s.target.ScrolledPosition(s.horizontal) > 0
		case lineDown, pageDown:
			return s.target.ScrolledPosition(s.horizontal) < s.target.ContentSize(s.horizontal)-s.target.VisibleSize(s.horizontal)
		case thumb:
			pos := s.target.ScrolledPosition(s.horizontal)
			return pos > 0 || pos < s.target.ContentSize(s.horizontal)-s.target.VisibleSize(s.horizontal)
		default:
		}
	}
	return false
}

func (s *ScrollBar) partRect(which scrollbarPart) geom.Rect {
	var rect geom.Rect
	switch which {
	case thumb:
		if s.target != nil {
			content := s.target.ContentSize(s.horizontal)
			visible := s.target.VisibleSize(s.horizontal)
			if content-visible > 0 {
				pos := s.target.ScrolledPosition(s.horizontal)
				full := s.ContentRect(false)
				if s.horizontal {
					full.X += s.barSize - 1
					full.Width -= s.barSize*2 - 2
					full.Height = s.barSize
					if full.Width > 0 {
						scale := full.Width / content
						visible *= scale
						min := s.barSize * 0.75
						if visible < min {
							scale = (full.Width + visible - min) / content
							visible = min
						}
						pos *= scale
						full.X += pos
						full.Width = visible + 1
						rect = full
					}
				} else {
					full.Y += s.barSize - 1
					full.Height -= s.barSize*2 - 2
					full.Width = s.barSize
					if full.Height > 0 {
						scale := full.Height / content
						visible *= scale
						min := s.barSize * 0.75
						if visible < min {
							scale = (full.Height + visible - min) / content
							visible = min
						}
						pos *= scale
						full.Y += pos
						full.Height = visible + 1
						rect = full
					}
				}
			}
		}
	case lineUp:
		rect = s.ContentRect(false)
		rect.Width = s.barSize
		rect.Height = s.barSize
	case lineDown:
		rect = s.ContentRect(false)
		if s.horizontal {
			rect.X += rect.Width - s.barSize
		} else {
			rect.Y += rect.Height - s.barSize
		}
		rect.Width = s.barSize
		rect.Height = s.barSize
	case pageUp:
		rect = s.partRect(lineUp)
		thumb := s.partRect(thumb)
		if s.horizontal {
			rect.X += rect.Width
			rect.Width = thumb.X - rect.X
		} else {
			rect.Y += rect.Height
			rect.Height = thumb.Y - rect.Y
		}
	case pageDown:
		rect = s.partRect(lineDown)
		thumb := s.partRect(thumb)
		if s.horizontal {
			x := thumb.X + thumb.Width
			rect.Width = rect.X - x
			rect.X = x
		} else {
			y := thumb.Y + thumb.Height
			rect.Height = rect.Y - y
			rect.Y = y
		}
	}
	rect.Align()
	return rect
}

// SetScrolledPosition attempts to set the current scrolled position of this
// ScrollBar to the specified value. The value will be clipped to the
// available range. If no target has been set, then nothing will happen.
func (s *ScrollBar) SetScrolledPosition(position float64) {
	if s.target != nil {
		position = math.Max(math.Min(position, s.target.ContentSize(s.horizontal)-s.target.VisibleSize(s.horizontal)), 0)
		if s.target.ScrolledPosition(s.horizontal) != position {
			s.target.SetScrolledPosition(s.horizontal, position)
			s.MarkForRedraw()
		}
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (s *ScrollBar) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	s.sequence++
	what := s.over(where)
	if s.partEnabled(what) {
		s.pressed = what
		switch what {
		case thumb:
			if s.horizontal {
				s.thumbDown = where.X - s.partRect(what).X
			} else {
				s.thumbDown = where.Y - s.partRect(what).Y
			}
		case lineUp, lineDown, pageUp, pageDown:
			s.scheduleRepeat(what, s.initialRepeatDelay)
		}
		s.MarkForRedraw()
	}
	return true
}

func (s *ScrollBar) over(where geom.Point) scrollbarPart {
	for i := thumb; i <= pageDown; i++ {
		rect := s.partRect(i)
		if rect.ContainsPoint(where) {
			return i
		}
	}
	return none
}

func (s *ScrollBar) scheduleRepeat(which scrollbarPart, delay time.Duration) {
	if window := s.Window(); window != nil {
		current := s.sequence
		switch which {
		case lineUp:
			s.SetScrolledPosition(s.target.ScrolledPosition(s.horizontal) - math.Abs(s.target.ScrollAmount(s.horizontal, true, false)))
		case lineDown:
			s.SetScrolledPosition(s.target.ScrolledPosition(s.horizontal) + math.Abs(s.target.ScrollAmount(s.horizontal, false, false)))
		case pageUp:
			s.SetScrolledPosition(s.target.ScrolledPosition(s.horizontal) - math.Abs(s.target.ScrollAmount(s.horizontal, true, true)))
		case pageDown:
			s.SetScrolledPosition(s.target.ScrolledPosition(s.horizontal) + math.Abs(s.target.ScrollAmount(s.horizontal, false, true)))
		default:
			return
		}
		ux.InvokeAfter(func() {
			if current == s.sequence && s.pressed == which {
				s.scheduleRepeat(which, s.repeatDelay)
			}
		}, delay)
	}
}

// DefaultMouseDrag provides the default mouse drag handling.
func (s *ScrollBar) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	if s.pressed == thumb {
		var pos float64
		rect := s.partRect(lineUp)
		if s.horizontal {
			pos = where.X - (s.thumbDown + rect.X + rect.Width - 1)
		} else {
			pos = where.Y - (s.thumbDown + rect.Y + rect.Height - 1)
		}
		s.SetScrolledPosition(pos / s.thumbScale())
	}
}

func (s *ScrollBar) thumbScale() float64 {
	var scale float64 = 1
	content := s.target.ContentSize(s.horizontal)
	visible := s.target.VisibleSize(s.horizontal)
	if content-visible > 0 {
		var size float64
		min := s.barSize * 0.75
		rect := s.ContentRect(false)
		if s.horizontal {
			size = rect.Width
		} else {
			size = rect.Height
		}
		size -= s.barSize*2 + 2
		if size > 0 {
			scale = size / content
			visible *= scale
			if visible < min {
				scale = (size + visible - min) / content
			}
		}
	}
	return scale
}

// DefaultMouseUp provides the default mouse up handling.
func (s *ScrollBar) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	s.pressed = none
	s.MarkForRedraw()
}
