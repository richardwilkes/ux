// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package radiobutton

import (
	"math"
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/selectable"
)

// RadioButton represents a clickable radio button with an optional label.
type RadioButton struct {
	selectable.Panel
	managed
	ClickCallback func()
	Pressed       bool
}

// New creates a new radio button.
func New() *RadioButton {
	r := &RadioButton{}
	r.managed.initialize()
	r.InitTypeAndID(r)
	r.SetFocusable(true)
	r.SetSizer(r.DefaultSizes)
	r.DrawCallback = r.DefaultDraw
	r.GainedFocusCallback = r.MarkForRedraw
	r.LostFocusCallback = r.MarkForRedraw
	r.MouseDownCallback = r.DefaultMouseDown
	r.MouseDragCallback = r.DefaultMouseDrag
	r.MouseUpCallback = r.DefaultMouseUp
	r.KeyDownCallback = r.DefaultKeyDown
	return r
}

// DefaultSizes provides the default sizing.
func (r *RadioButton) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = r.circleAndLabelSize()
	if border := r.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (r *RadioButton) circleAndLabelSize() geom.Size {
	circleSize := r.circleSize()
	if r.image == nil && r.text == "" {
		return geom.Size{Width: circleSize, Height: circleSize}
	}
	size := widget.LabelSize(r.text, r.font, r.image, r.side, r.gap)
	size.Width += r.gap + circleSize
	if size.Height < circleSize {
		size.Height = circleSize
	}
	return size
}

func (r *RadioButton) circleSize() float64 {
	return math.Ceil(r.font.Height())
}

// DefaultDraw provides the default drawing.
func (r *RadioButton) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !r.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := r.ContentRect(false)
	size := r.circleAndLabelSize()
	switch r.hAlign {
	case align.Middle, align.Fill:
		rect.X = math.Floor(rect.X + (rect.Width-size.Width)/2)
	case align.End:
		rect.X += rect.Width - size.Width
	default: // Start
	}
	switch r.vAlign {
	case align.Middle, align.Fill:
		rect.Y = math.Floor(rect.Y + (rect.Height-size.Height)/2)
	case align.End:
		rect.Y += rect.Height - size.Height
	default: // Start
	}
	rect.Size = size
	circleSize := r.circleSize()
	if r.image != nil || r.text != "" {
		rct := rect
		rct.X += circleSize + r.gap
		rct.Width -= circleSize + r.gap
		widget.DrawLabel(gc, rct, r.hAlign, r.vAlign, r.text, r.font, r.textInk, r.image, r.side, r.gap, r.Enabled())
	}
	if rect.Height > circleSize {
		rect.Y += math.Floor((rect.Height - circleSize) / 2)
	}
	rect.Width = circleSize
	rect.Height = circleSize
	widget.DrawEllipseBase(gc, rect, r.currentBackgroundInk(), r.edgeInk)
	if r.Selected() {
		rect.InsetUniform(0.5 + 0.2*circleSize)
		gc.Ellipse(rect)
		gc.Fill(r.currentMarkInk())
	}
}

func (r *RadioButton) currentBackgroundInk() draw.Ink {
	switch {
	case r.Pressed:
		return r.pressedBackgroundInk
	case r.Focused():
		return r.focusedBackgroundInk
	default:
		return r.backgroundInk
	}
}

func (r *RadioButton) currentMarkInk() draw.Ink {
	if r.Pressed || r.Focused() {
		return r.pressedTextInk
	}
	return r.textInk
}

// Click makes the radio button behave as if a user clicked on it.
func (r *RadioButton) Click() {
	r.SetSelected(true)
	pressed := r.Pressed
	r.Pressed = true
	r.MarkForRedraw()
	r.FlushDrawing()
	r.Pressed = pressed
	time.Sleep(r.clickAnimationTime)
	r.MarkForRedraw()
	if r.ClickCallback != nil {
		r.ClickCallback()
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (r *RadioButton) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	r.Pressed = true
	r.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (r *RadioButton) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := r.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if r.Pressed != pressed {
		r.Pressed = pressed
		r.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (r *RadioButton) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	r.Pressed = false
	r.MarkForRedraw()
	rect := r.ContentRect(false)
	if rect.ContainsPoint(where) {
		r.SetSelected(true)
		if r.ClickCallback != nil {
			r.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (r *RadioButton) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		r.Click()
		return true
	}
	return false
}
