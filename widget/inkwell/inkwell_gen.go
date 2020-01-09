// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Code created from "widget.go.tmpl" - don't edit by hand

package inkwell

import (
	"time"

	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	backgroundInk        draw.Ink
	focusedBackgroundInk draw.Ink
	pressedBackgroundInk draw.Ink
	edgeInk              draw.Ink
	edgeHighlightInk     draw.Ink
	imageScale           float64
	contentSize          float64
	cornerRadius         float64
	clickAnimationTime   time.Duration
}

func (m *managed) initialize() {
	m.backgroundInk = draw.ControlBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.edgeHighlightInk = draw.ControlEdgeHighlightAdjColor
	m.imageScale = 0.5
	m.contentSize = 20
	m.cornerRadius = 4
	m.clickAnimationTime = time.Millisecond * 100
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not pressed or focused.
func (well *InkWell) BackgroundInk() draw.Ink {
	return well.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not pressed or focused. Pass in nil to use the default.
func (well *InkWell) SetBackgroundInk(value draw.Ink) *InkWell {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if well.backgroundInk != value {
		well.backgroundInk = value
		well.MarkForRedraw()
	}
	return well
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (well *InkWell) FocusedBackgroundInk() draw.Ink {
	return well.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (well *InkWell) SetFocusedBackgroundInk(value draw.Ink) *InkWell {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if well.focusedBackgroundInk != value {
		well.focusedBackgroundInk = value
		well.MarkForRedraw()
	}
	return well
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (well *InkWell) PressedBackgroundInk() draw.Ink {
	return well.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (well *InkWell) SetPressedBackgroundInk(value draw.Ink) *InkWell {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if well.pressedBackgroundInk != value {
		well.pressedBackgroundInk = value
		well.MarkForRedraw()
	}
	return well
}

// EdgeInk returns the ink that will be used for the edges.
func (well *InkWell) EdgeInk() draw.Ink {
	return well.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (well *InkWell) SetEdgeInk(value draw.Ink) *InkWell {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if well.edgeInk != value {
		well.edgeInk = value
		well.MarkForRedraw()
	}
	return well
}

// EdgeHighlightInk returns the ink that will be used just inside the edges.
func (well *InkWell) EdgeHighlightInk() draw.Ink {
	return well.edgeHighlightInk
}

// SetEdgeHighlightInk sets the ink that will be used just inside the edges.
// Pass in nil to use the default.
func (well *InkWell) SetEdgeHighlightInk(value draw.Ink) *InkWell {
	if value == nil {
		value = draw.ControlEdgeHighlightAdjColor
	}
	if well.edgeHighlightInk != value {
		well.edgeHighlightInk = value
		well.MarkForRedraw()
	}
	return well
}

// ImageScale returns the image scale to use for images dropped onto the
// well. Defaults to 0.5 to support retina displays.
func (well *InkWell) ImageScale() float64 {
	return well.imageScale
}

// SetImageScale sets the image scale to use for images dropped onto the
// well. Defaults to 0.5 to support retina displays.
func (well *InkWell) SetImageScale(value float64) *InkWell {
	if value < 0.05 {
		value = 0.05
	}
	if well.imageScale != value {
		well.imageScale = value
		well.MarkForRedraw()
	}
	return well
}

// ContentSize returns the content width and height.
func (well *InkWell) ContentSize() float64 {
	return well.contentSize
}

// SetContentSize sets the content width and height.
func (well *InkWell) SetContentSize(value float64) *InkWell {
	if value < 8 {
		value = 8
	}
	if well.contentSize != value {
		well.contentSize = value
		well.MarkForLayoutAndRedraw()
	}
	return well
}

// CornerRadius returns the amount of rounding to use on the corners.
func (well *InkWell) CornerRadius() float64 {
	return well.cornerRadius
}

// SetCornerRadius sets the amount of rounding to use on the corners.
func (well *InkWell) SetCornerRadius(value float64) *InkWell {
	if value < 0 {
		value = 0
	}
	if well.cornerRadius != value {
		well.cornerRadius = value
		well.MarkForRedraw()
	}
	return well
}

// ClickAnimationTime returns the amount of time to spend animating the click
// action.
func (well *InkWell) ClickAnimationTime() time.Duration {
	return well.clickAnimationTime
}

// SetClickAnimationTime sets the amount of time to spend animating the click
// action.
func (well *InkWell) SetClickAnimationTime(value time.Duration) *InkWell {
	if value < time.Millisecond*50 {
		value = time.Millisecond * 50
	}
	if well.clickAnimationTime != value {
		well.clickAnimationTime = value
	}
	return well
}

// SetBorder sets the border. May be nil.
func (well *InkWell) SetBorder(value border.Border) *InkWell {
	well.Panel.SetBorder(value)
	return well
}

// SetEnabled sets enabled state.
func (well *InkWell) SetEnabled(enabled bool) *InkWell {
	well.Panel.SetEnabled(enabled)
	return well
}

// SetFocusable whether it can have the keyboard focus.
func (well *InkWell) SetFocusable(focusable bool) *InkWell {
	well.Panel.SetFocusable(focusable)
	return well
}
