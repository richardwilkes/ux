// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package scrollarea

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/widget/scrollarea/behavior"
	"github.com/richardwilkes/ux/widget/scrollbar"
)

// ScrollArea provides a panel that can hold another panel and show it through
// a scrollable viewport.
type ScrollArea struct {
	ux.Panel
	managed
	hBar     *scrollbar.ScrollBar
	vBar     *scrollbar.ScrollBar
	view     *ux.Panel
	content  *ux.Panel
	behavior behavior.Behavior
}

// New creates a new, empty ScrollArea.
func New() *ScrollArea {
	s := &ScrollArea{}
	s.managed.initialize()
	s.InitTypeAndID(s)
	s.SetBorder(s.unfocusedBorder)
	s.view = ux.NewPanel()
	s.hBar = scrollbar.New(true).SetTarget(s)
	s.vBar = scrollbar.New(false).SetTarget(s)
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
	return s
}

// Content returns the content panel. May be nil.
func (s *ScrollArea) Content() *ux.Panel {
	return s.content
}

// SetContent sets the content panel, replacing any existing one.
func (s *ScrollArea) SetContent(content *ux.Panel, behave behavior.Behavior) *ScrollArea {
	if s.content != nil {
		s.content.RemoveFromParent()
	}
	s.content = content
	s.behavior = behave
	if s.content != nil {
		s.view.AddChildAtIndex(s.content, 0)
	}
	s.MarkForLayoutAndRedraw()
	return s
}

// DefaultDraw provides the default drawing.
func (s *ScrollArea) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	gc.Rect(s.ContentRect(true))
	gc.Fill(s.backgroundInk)
}

// DefaultMouseWheel provides the default mouse wheel handling.
func (s *ScrollArea) DefaultMouseWheel(where, delta geom.Point, mod keys.Modifiers) bool {
	if delta.Y != 0 {
		s.vBar.SetScrolledPosition(s.ScrolledPosition(false) - delta.Y*s.ScrollAmount(false, delta.Y > 0, false))
	}
	if delta.X != 0 {
		s.hBar.SetScrolledPosition(s.ScrolledPosition(true) - delta.X*s.ScrollAmount(true, delta.X > 0, true))
	}
	return true
}

// DefaultFocusGained provides the default focus gained handling.
func (s *ScrollArea) DefaultFocusGained() {
	s.SetBorder(s.focusedBorder)
	s.MarkForRedraw()
}

// DefaultFocusLost provides the default focus lost handling.
func (s *ScrollArea) DefaultFocusLost() {
	s.SetBorder(s.unfocusedBorder)
	s.MarkForRedraw()
}

// DefaultKeyDown provides the default key down handling.
func (s *ScrollArea) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	switch keyCode {
	case keys.Up.Code, keys.NumpadUp.Code:
		s.vBar.SetScrolledPosition(s.ScrolledPosition(false) - s.ScrollAmount(false, true, false))
	case keys.Down.Code, keys.NumpadDown.Code:
		s.vBar.SetScrolledPosition(s.ScrolledPosition(false) + s.ScrollAmount(false, false, false))
	case keys.Left.Code, keys.NumpadLeft.Code:
		s.hBar.SetScrolledPosition(s.ScrolledPosition(true) - s.ScrollAmount(true, true, false))
	case keys.Right.Code, keys.NumpadRight.Code:
		s.hBar.SetScrolledPosition(s.ScrolledPosition(true) + s.ScrollAmount(true, false, false))
	case keys.Home.Code, keys.NumpadHome.Code:
		s.ScrollBar(mod.ShiftDown()).SetScrolledPosition(0)
	case keys.End.Code, keys.NumpadEnd.Code:
		s.ScrollBar(mod.ShiftDown()).SetScrolledPosition(s.ContentSize(mod.ShiftDown()))
	case keys.PageUp.Code, keys.NumpadPageUp.Code:
		s.ScrollBar(mod.ShiftDown()).SetScrolledPosition(s.ScrolledPosition(mod.ShiftDown()) - s.ScrollAmount(mod.ShiftDown(), true, true))
	case keys.PageDown.Code, keys.NumpadPageDown.Code:
		s.ScrollBar(mod.ShiftDown()).SetScrolledPosition(s.ScrolledPosition(mod.ShiftDown()) + s.ScrollAmount(mod.ShiftDown(), false, true))
	default:
		return false
	}
	return true
}

// ScrollBar returns the scrollbar for the given axis.
func (s *ScrollArea) ScrollBar(horizontal bool) *scrollbar.ScrollBar {
	if horizontal {
		return s.hBar
	}
	return s.vBar
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
	viewRect := s.view.ContentRect(true)
	hAdj := computeScrollAdj(rect.X, viewRect.X, rect.Y+rect.Width, viewRect.X+viewRect.Width)
	vAdj := computeScrollAdj(rect.Y, viewRect.Y, rect.Y+rect.Height, viewRect.Y+viewRect.Height)
	if hAdj != 0 || vAdj != 0 {
		if hAdj != 0 {
			s.hBar.SetScrolledPosition(s.ScrolledPosition(true) + hAdj)
		}
		if vAdj != 0 {
			s.vBar.SetScrolledPosition(s.ScrolledPosition(false) + vAdj)
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
