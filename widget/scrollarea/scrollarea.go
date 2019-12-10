package scrollarea

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/widget/scrollarea/behavior"
	"github.com/richardwilkes/ux/widget/scrollbar"
)

// ScrollArea provides a panel that can hold another panel and show it through
// a scrollable viewport.
type ScrollArea struct {
	ux.Panel
	HBar          *scrollbar.ScrollBar
	VBar          *scrollbar.ScrollBar
	view          *ux.Panel
	content       *ux.Panel
	FocusedBorder border.Border // The border to use when focused.
	BackgroundInk draw.Ink
	behavior      behavior.Behavior
}

// New creates a new ScrollArea with the specified panel as its
// content. The content may be nil.
func New(content *ux.Panel, behave behavior.Behavior) *ScrollArea {
	s := &ScrollArea{
		FocusedBorder: border.NewLine(draw.ControlAccentColor, geom.NewUniformInsets(2), true),
		BackgroundInk: draw.TextBackgroundColor,
	}
	s.InitTypeAndID(s)
	s.SetBorder(border.NewLine(draw.ControlEdgeAdjColor, geom.NewUniformInsets(1), false))
	s.view = ux.NewPanel()
	s.HBar = scrollbar.New(s, true)
	s.VBar = scrollbar.New(s, false)
	s.AddChild(s.view)
	s.SetFocusable(true)
	s.SetLayout(&scrollAreaLayout{scrollArea: s})
	s.DrawCallback = s.DefaultDraw
	s.MouseWheelCallback = s.DefaultMouseWheel
	s.GainedFocusCallback = s.DefaultFocusGained
	s.LostFocusCallback = s.DefaultFocusLost
	s.KeyDownCallback = s.DefaultKeyDown
	s.FrameChangeInChildHierarchyCallback = s.DefaultFrameChangeInChildHierarchy
	s.ScrollRectIntoViewCallback = s.DefaultScrollRectIntoView
	if content != nil {
		s.SetContent(content, behave)
	}
	return s
}

// SetContent sets the content panel, replacing any existing one.
func (s *ScrollArea) SetContent(content *ux.Panel, behave behavior.Behavior) {
	if s.content != nil {
		s.content.RemoveFromParent()
	}
	s.content = content
	s.behavior = behave
	if s.content != nil {
		s.view.AddChildAtIndex(s.content, 0)
	}
	s.MarkForLayoutAndRedraw()
}

// DefaultDraw provides the default drawing.
func (s *ScrollArea) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	gc.Rect(dirty)
	gc.Fill(s.BackgroundInk)
}

// DefaultMouseWheel provides the default mouse wheel handling.
func (s *ScrollArea) DefaultMouseWheel(where, delta geom.Point, mod keys.Modifiers) bool {
	if delta.Y != 0 {
		s.VBar.SetScrolledPosition(s.ScrolledPosition(false) - delta.Y*s.ScrollAmount(false, delta.Y > 0, false))
	}
	if delta.X != 0 {
		s.HBar.SetScrolledPosition(s.ScrolledPosition(true) - delta.X*s.ScrollAmount(true, delta.X > 0, true))
	}
	return true
}

// DefaultFocusGained provides the default focus gained handling.
func (s *ScrollArea) DefaultFocusGained() {
	s.view.SetBorder(s.FocusedBorder)
	s.view.MarkForRedraw()
}

// DefaultFocusLost provides the default focus lost handling.
func (s *ScrollArea) DefaultFocusLost() {
	s.view.SetBorder(nil)
	s.view.MarkForRedraw()
}

// DefaultKeyDown provides the default key down handling.
func (s *ScrollArea) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	switch keyCode {
	case keys.Up.Code, keys.NumpadUp.Code:
		s.VBar.SetScrolledPosition(s.ScrolledPosition(false) - s.ScrollAmount(false, true, false))
	case keys.Down.Code, keys.NumpadDown.Code:
		s.VBar.SetScrolledPosition(s.ScrolledPosition(false) + s.ScrollAmount(false, false, false))
	case keys.Left.Code, keys.NumpadLeft.Code:
		s.HBar.SetScrolledPosition(s.ScrolledPosition(true) - s.ScrollAmount(true, true, false))
	case keys.Right.Code, keys.NumpadRight.Code:
		s.HBar.SetScrolledPosition(s.ScrolledPosition(true) + s.ScrollAmount(true, false, false))
	case keys.Home.Code, keys.NumpadHome.Code:
		s.barForMod(mod).SetScrolledPosition(0)
	case keys.End.Code, keys.NumpadEnd.Code:
		s.barForMod(mod).SetScrolledPosition(s.ContentSize(mod.ShiftDown()))
	case keys.PageUp.Code, keys.NumpadPageUp.Code:
		s.barForMod(mod).SetScrolledPosition(s.ScrolledPosition(mod.ShiftDown()) - s.ScrollAmount(mod.ShiftDown(), true, true))
	case keys.PageDown.Code, keys.NumpadPageDown.Code:
		s.barForMod(mod).SetScrolledPosition(s.ScrolledPosition(mod.ShiftDown()) + s.ScrollAmount(mod.ShiftDown(), false, true))
	default:
		return false
	}
	return true
}

func (s *ScrollArea) barForMod(mod keys.Modifiers) *scrollbar.ScrollBar {
	if mod.ShiftDown() {
		return s.HBar
	}
	return s.VBar
}

// DefaultFrameChangeInChildHierarchy provides the default frame change in
// child hierarchy handling.
func (s *ScrollArea) DefaultFrameChangeInChildHierarchy(panel *ux.Panel) {
	if s.content != nil {
		vs := s.view.ContentRect(false).Size
		rect := s.content.FrameRect()
		nl := rect.Point
		if rect.Y != 0 && vs.Height > rect.Y+rect.Height {
			nl.Y = math.Min(vs.Height-rect.Height, 0)
		}
		if rect.X != 0 && vs.Width > rect.X+rect.Width {
			nl.X = math.Min(vs.Width-rect.Width, 0)
		}
		if nl != rect.Point {
			rect.Point = nl
			s.content.SetFrameRect(rect)
		}
	}
}

// ScrollAmount implements the ScrollPager interface.
func (s *ScrollArea) ScrollAmount(horizontal, towardsStart, page bool) float64 {
	if s.content != nil {
		if pager, ok := s.content.Self().(scrollbar.ScrollPager); ok {
			return pager.ScrollAmount(horizontal, towardsStart, page)
		}
	}
	if !page {
		return 16
	}
	return s.VisibleSize(horizontal)
}

// ScrolledPosition implements the Scrollable interface.
func (s *ScrollArea) ScrolledPosition(horizontal bool) float64 {
	if s.content == nil {
		return 0
	}
	rect := s.content.FrameRect()
	if horizontal {
		return -rect.X
	}
	return -rect.Y
}

// SetScrolledPosition implements the Scrollable interface.
func (s *ScrollArea) SetScrolledPosition(horizontal bool, position float64) {
	if s.content != nil {
		rect := s.content.FrameRect()
		if horizontal {
			rect.X = -position
		} else {
			rect.Y = -position
		}
		s.content.SetFrameRect(rect)
	}
}

// VisibleSize implements the Scrollable interface.
func (s *ScrollArea) VisibleSize(horizontal bool) float64 {
	rect := s.view.FrameRect()
	if horizontal {
		return rect.Width
	}
	return rect.Height
}

// ContentSize implements the Scrollable interface.
func (s *ScrollArea) ContentSize(horizontal bool) float64 {
	if s.content == nil {
		return 0
	}
	rect := s.content.FrameRect()
	if horizontal {
		return rect.Width
	}
	return rect.Height
}

// DefaultScrollRectIntoView provides the default scroll rect into view
// handling.
func (s *ScrollArea) DefaultScrollRectIntoView(rect geom.Rect) bool {
	viewRect := s.view.FrameRect()
	hAdj := computeScrollAdj(rect.X, viewRect.X, rect.Y+rect.Width, viewRect.X+viewRect.Width)
	vAdj := computeScrollAdj(rect.Y, viewRect.Y, rect.Y+rect.Height, viewRect.Y+viewRect.Height)
	if hAdj != 0 || vAdj != 0 {
		if hAdj != 0 {
			s.HBar.SetScrolledPosition(s.ScrolledPosition(true) + hAdj)
		}
		if vAdj != 0 {
			s.VBar.SetScrolledPosition(s.ScrolledPosition(false) + vAdj)
		}
		return true
	}
	return false
}

func computeScrollAdj(upper1, upper2, lower1, lower2 float64) float64 {
	if upper1 < upper2 {
		return upper1 - upper2
	}
	if lower1 > lower2 {
		if lower1-upper1 <= lower2-upper2 {
			return lower1 - lower2
		}
		return upper1 - upper2
	}
	return 0
}
