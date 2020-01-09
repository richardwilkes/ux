// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package checkbox

import (
	"math"
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/checkbox/state"
)

// CheckBox represents a clickable checkbox with an optional label.
type CheckBox struct {
	ux.Panel
	managed
	ClickCallback func()
	Pressed       bool
}

// New creates a new checkbox.
func New() *CheckBox {
	c := &CheckBox{}
	c.managed.initialize()
	c.InitTypeAndID(c)
	c.SetFocusable(true)
	c.SetSizer(c.DefaultSizes)
	c.DrawCallback = c.DefaultDraw
	c.GainedFocusCallback = c.MarkForRedraw
	c.LostFocusCallback = c.MarkForRedraw
	c.MouseDownCallback = c.DefaultMouseDown
	c.MouseDragCallback = c.DefaultMouseDrag
	c.MouseUpCallback = c.DefaultMouseUp
	c.KeyDownCallback = c.DefaultKeyDown
	return c
}

// DefaultSizes provides the default sizing.
func (c *CheckBox) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = c.boxAndLabelSize()
	if border := c.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (c *CheckBox) boxAndLabelSize() geom.Size {
	boxSize := c.boxSize()
	if c.image == nil && c.text == "" {
		return geom.Size{Width: boxSize, Height: boxSize}
	}
	size := widget.LabelSize(c.text, c.font, c.image, c.side, c.gap)
	size.Width += c.gap + boxSize
	if size.Height < boxSize {
		size.Height = boxSize
	}
	return size
}

func (c *CheckBox) boxSize() float64 {
	return math.Ceil(c.font.Height())
}

// DefaultDraw provides the default drawing.
func (c *CheckBox) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	rect := c.ContentRect(false)
	size := c.boxAndLabelSize()
	switch c.hAlign {
	case align.Middle, align.Fill:
		rect.X = math.Floor(rect.X + (rect.Width-size.Width)/2)
	case align.End:
		rect.X += rect.Width - size.Width
	default: // Start
	}
	switch c.vAlign {
	case align.Middle, align.Fill:
		rect.Y = math.Floor(rect.Y + (rect.Height-size.Height)/2)
	case align.End:
		rect.Y += rect.Height - size.Height
	default: // Start
	}
	rect.Size = size
	boxSize := c.boxSize()
	if c.image != nil || c.text != "" {
		r := rect
		r.X += boxSize + c.gap
		r.Width -= boxSize + c.gap
		widget.DrawLabel(gc, r, c.hAlign, c.vAlign, c.text, c.font, c.textInk, c.image, c.side, c.gap, c.Enabled())
	}
	if rect.Height > boxSize {
		rect.Y += math.Floor((rect.Height - boxSize) / 2)
	}
	rect.Width = boxSize
	rect.Height = boxSize
	if !c.Enabled() {
		gc.SetOpacity(0.33)
	}
	widget.DrawRoundedRectBase(gc, rect, c.cornerRadius, c.currentBackgroundInk(), c.edgeInk)
	rect.InsetUniform(0.5)
	switch c.State() {
	case state.Mixed:
		gc.SetStrokeWidth(2)
		gc.MoveTo(rect.X+rect.Width*0.25, rect.Y+rect.Height*0.5)
		gc.LineTo(rect.X+rect.Width*0.7, rect.Y+rect.Height*0.5)
		gc.Stroke(c.currentMarkInk())
	case state.On:
		gc.SetStrokeWidth(2)
		gc.MoveTo(rect.X+rect.Width*0.25, rect.Y+rect.Height*0.55)
		gc.LineTo(rect.X+rect.Width*0.45, rect.Y+rect.Height*0.7)
		gc.LineTo(rect.X+rect.Width*0.75, rect.Y+rect.Height*0.3)
		gc.Stroke(c.currentMarkInk())
	}
}

func (c *CheckBox) currentBackgroundInk() draw.Ink {
	switch {
	case c.Pressed:
		return c.pressedBackgroundInk
	case c.Focused():
		return c.focusedBackgroundInk
	default:
		return c.backgroundInk
	}
}

func (c *CheckBox) currentMarkInk() draw.Ink {
	if c.Pressed || c.Focused() {
		return c.pressedTextInk
	}
	return c.textInk
}

// Click makes the checkbox behave as if a user clicked on it.
func (c *CheckBox) Click() {
	c.updateState()
	pressed := c.Pressed
	c.Pressed = true
	c.MarkForRedraw()
	c.FlushDrawing()
	c.Pressed = pressed
	time.Sleep(c.clickAnimationTime)
	c.MarkForRedraw()
	if c.ClickCallback != nil {
		c.ClickCallback()
	}
}

func (c *CheckBox) updateState() {
	if c.State() == state.On {
		c.SetState(state.Off)
	} else {
		c.SetState(state.On)
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (c *CheckBox) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	c.Pressed = true
	c.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (c *CheckBox) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := c.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if c.Pressed != pressed {
		c.Pressed = pressed
		c.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (c *CheckBox) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	c.Pressed = false
	c.MarkForRedraw()
	rect := c.ContentRect(false)
	if rect.ContainsPoint(where) {
		c.updateState()
		if c.ClickCallback != nil {
			c.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (c *CheckBox) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		c.Click()
		return true
	}
	return false
}
