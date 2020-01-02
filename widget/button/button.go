// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package button

import (
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/selectable"
)

// Button represents a clickable button.
type Button struct {
	selectable.Panel
	managed
	ClickCallback func()
	Pressed       bool
}

// New creates a new button.
func New() *Button {
	b := &Button{}
	b.managed.initialize()
	b.InitTypeAndID(b)
	b.SetFocusable(true)
	b.SetSizer(b.DefaultSizes)
	b.DrawCallback = b.DefaultDraw
	b.GainedFocusCallback = b.MarkForRedraw
	b.LostFocusCallback = b.MarkForRedraw
	b.MouseDownCallback = b.DefaultMouseDown
	b.MouseDragCallback = b.DefaultMouseDrag
	b.MouseUpCallback = b.DefaultMouseUp
	b.KeyDownCallback = b.DefaultKeyDown
	return b
}

// DefaultSizes provides the default sizing.
func (b *Button) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	text := b.text
	if b.image == nil && text == "" {
		text = "M"
	}
	pref = widget.LabelSize(text, b.font, b.image, b.side, b.gap)
	if theBorder := b.Border(); theBorder != nil {
		pref.AddInsets(theBorder.Insets())
	}
	pref.Width += b.horizontalMargin()*2 + 2
	pref.Height += b.verticalMargin()*2 + 2
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (b *Button) horizontalMargin() float64 {
	if b.text == "" && b.image != nil {
		return b.imageOnlyHMargin
	}
	return b.hMargin
}

func (b *Button) verticalMargin() float64 {
	if b.text == "" && b.image != nil {
		return b.imageOnlyVMargin
	}
	return b.vMargin
}

// DefaultDraw provides the default drawing.
func (b *Button) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !b.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := b.ContentRect(false)
	widget.DrawRoundedRectBase(gc, rect, b.cornerRadius, b.currentBackgroundInk(), b.edgeInk)
	rect.InsetUniform(1.5)
	rect.X += b.horizontalMargin()
	rect.Y += b.verticalMargin()
	rect.Width -= b.horizontalMargin() * 2
	rect.Height -= b.verticalMargin() * 2
	widget.DrawLabel(gc, rect, b.hAlign, b.vAlign, b.text, b.font, b.currentTextInk(), b.image, b.side, b.gap, b.Enabled())
}

func (b *Button) currentBackgroundInk() draw.Ink {
	switch {
	case b.Pressed:
		return b.pressedBackgroundInk
	case b.Focused():
		return b.focusedBackgroundInk
	case b.sticky && b.Selected():
		return b.selectedBackgroundInk
	default:
		return b.backgroundInk
	}
}

func (b *Button) currentTextInk() draw.Ink {
	if b.Pressed || b.Focused() {
		return b.pressedTextInk
	}
	return b.textInk
}

// Click makes the button behave as if a user clicked on it.
func (b *Button) Click() {
	b.SetSelected(true)
	pressed := b.Pressed
	b.Pressed = true
	b.MarkForRedraw()
	b.FlushDrawing()
	b.Pressed = pressed
	time.Sleep(b.clickAnimationTime)
	b.MarkForRedraw()
	if b.ClickCallback != nil {
		b.ClickCallback()
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (b *Button) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	b.Pressed = true
	b.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (b *Button) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := b.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if b.Pressed != pressed {
		b.Pressed = pressed
		b.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (b *Button) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	b.Pressed = false
	b.MarkForRedraw()
	rect := b.ContentRect(false)
	if rect.ContainsPoint(where) {
		b.SetSelected(true)
		if b.ClickCallback != nil {
			b.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (b *Button) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		b.Click()
		return true
	}
	return false
}
