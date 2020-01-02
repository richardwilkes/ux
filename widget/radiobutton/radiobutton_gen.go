// Code created from "widget.go.tmpl" - don't edit by hand
//
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
	"time"

	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
)

type managed struct {
	image                *draw.Image //nolint:structcheck
	text                 string      //nolint:structcheck
	font                 *draw.Font
	backgroundInk        draw.Ink
	focusedBackgroundInk draw.Ink
	pressedBackgroundInk draw.Ink
	edgeInk              draw.Ink
	textInk              draw.Ink
	pressedTextInk       draw.Ink
	gap                  float64
	cornerRadius         float64
	clickAnimationTime   time.Duration
	hAlign               align.Alignment
	vAlign               align.Alignment
	side                 side.Side
}

func (m *managed) initialize() {
	m.font = draw.SystemFont
	m.backgroundInk = draw.ControlBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.textInk = draw.ControlTextColor
	m.pressedTextInk = draw.AlternateSelectedControlTextColor
	m.gap = 3
	m.cornerRadius = 4
	m.clickAnimationTime = time.Millisecond * 100
	m.hAlign = align.Middle
	m.vAlign = align.Middle
	m.side = side.Left
}

// Image returns the image. May be nil.
func (r *RadioButton) Image() *draw.Image {
	return r.image
}

// SetImage sets the image. May be nil.
func (r *RadioButton) SetImage(value *draw.Image) *RadioButton {
	if r.image != value {
		r.image = value
		r.MarkForLayoutAndRedraw()
	}
	return r
}

// Text returns the text content.
func (r *RadioButton) Text() string {
	return r.text
}

// SetText sets the text content.
func (r *RadioButton) SetText(value string) *RadioButton {
	if r.text != value {
		r.text = value
		r.MarkForLayoutAndRedraw()
	}
	return r
}

// Font returns the font that will be used when drawing text content.
func (r *RadioButton) Font() *draw.Font {
	return r.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (r *RadioButton) SetFont(value *draw.Font) *RadioButton {
	if value == nil {
		value = draw.SystemFont
	}
	if r.font != value {
		r.font = value
		r.MarkForLayoutAndRedraw()
	}
	return r
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not pressed or focused.
func (r *RadioButton) BackgroundInk() draw.Ink {
	return r.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not pressed or focused. Pass in nil to use the default.
func (r *RadioButton) SetBackgroundInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if r.backgroundInk != value {
		r.backgroundInk = value
		r.MarkForRedraw()
	}
	return r
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (r *RadioButton) FocusedBackgroundInk() draw.Ink {
	return r.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (r *RadioButton) SetFocusedBackgroundInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if r.focusedBackgroundInk != value {
		r.focusedBackgroundInk = value
		r.MarkForRedraw()
	}
	return r
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (r *RadioButton) PressedBackgroundInk() draw.Ink {
	return r.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (r *RadioButton) SetPressedBackgroundInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if r.pressedBackgroundInk != value {
		r.pressedBackgroundInk = value
		r.MarkForRedraw()
	}
	return r
}

// EdgeInk returns the ink that will be used for the edges.
func (r *RadioButton) EdgeInk() draw.Ink {
	return r.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (r *RadioButton) SetEdgeInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if r.edgeInk != value {
		r.edgeInk = value
		r.MarkForRedraw()
	}
	return r
}

// TextInk returns the ink that will be used for the text when disabled or
// not pressed.
func (r *RadioButton) TextInk() draw.Ink {
	return r.textInk
}

// SetTextInk sets the ink that will be used for the text when disabled or
// not pressed. Pass in nil to use the default.
func (r *RadioButton) SetTextInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.ControlTextColor
	}
	if r.textInk != value {
		r.textInk = value
		r.MarkForRedraw()
	}
	return r
}

// PressedTextInk returns the ink that will be used for the text when enabled
// and pressed.
func (r *RadioButton) PressedTextInk() draw.Ink {
	return r.pressedTextInk
}

// SetPressedTextInk sets the ink that will be used for the text when enabled
// and pressed. Pass in nil to use the default.
func (r *RadioButton) SetPressedTextInk(value draw.Ink) *RadioButton {
	if value == nil {
		value = draw.AlternateSelectedControlTextColor
	}
	if r.pressedTextInk != value {
		r.pressedTextInk = value
		r.MarkForRedraw()
	}
	return r
}

// Gap returns the gap to put between the image and text.
func (r *RadioButton) Gap() float64 {
	return r.gap
}

// SetGap sets the gap to put between the image and text.
func (r *RadioButton) SetGap(value float64) *RadioButton {
	if value < 0 {
		value = 0
	}
	if r.gap != value {
		r.gap = value
		r.MarkForLayoutAndRedraw()
	}
	return r
}

// CornerRadius returns the amount of rounding to use on the corners.
func (r *RadioButton) CornerRadius() float64 {
	return r.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (r *RadioButton) SetCornerRadius(value float64) *RadioButton {
	if value < 0 {
		value = 0
	}
	if r.cornerRadius != value {
		r.cornerRadius = value
		r.MarkForRedraw()
	}
	return r
}

// ClickAnimationTime returns the amount of time to spend animating the click
// action.
func (r *RadioButton) ClickAnimationTime() time.Duration {
	return r.clickAnimationTime
}

// SetClickAnimationTime sets the amount of time to spend animating the click
// action.
func (r *RadioButton) SetClickAnimationTime(value time.Duration) *RadioButton {
	if value < time.Millisecond*50 {
		value = time.Millisecond * 50
	}
	if r.clickAnimationTime != value {
		r.clickAnimationTime = value
	}
	return r
}

// HAlign returns the horizontal alignment.
func (r *RadioButton) HAlign() align.Alignment {
	return r.hAlign
}

// SetHAlign sets the horizontal alignment.
func (r *RadioButton) SetHAlign(value align.Alignment) *RadioButton {
	if r.hAlign != value {
		r.hAlign = value
		r.MarkForRedraw()
	}
	return r
}

// VAlign returns the vertical alignment.
func (r *RadioButton) VAlign() align.Alignment {
	return r.vAlign
}

// SetVAlign sets the vertical alignment.
func (r *RadioButton) SetVAlign(value align.Alignment) *RadioButton {
	if r.vAlign != value {
		r.vAlign = value
		r.MarkForRedraw()
	}
	return r
}

// Side returns the side of the text the image should be on.
func (r *RadioButton) Side() side.Side {
	return r.side
}

// SetSide sets the side of the text the image should be on.
func (r *RadioButton) SetSide(value side.Side) *RadioButton {
	if r.side != value {
		r.side = value
		r.MarkForRedraw()
	}
	return r
}

// SetBorder sets the border. May be nil.
func (r *RadioButton) SetBorder(value border.Border) *RadioButton {
	r.Panel.SetBorder(value)
	return r
}

// SetEnabled sets enabled state.
func (r *RadioButton) SetEnabled(enabled bool) *RadioButton {
	r.Panel.SetEnabled(enabled)
	return r
}

// SetFocusable whether it can have the keyboard focus.
func (r *RadioButton) SetFocusable(focusable bool) *RadioButton {
	r.Panel.SetFocusable(focusable)
	return r
}

// SetSelected sets the panel's selected state.
func (r *RadioButton) SetSelected(selected bool) *RadioButton {
	r.Panel.SetSelected(selected)
	return r
}
