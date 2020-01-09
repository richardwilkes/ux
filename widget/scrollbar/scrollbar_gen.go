// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Code created from "widget.go.tmpl" - don't edit by hand

package scrollbar

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
	markInk              draw.Ink
	disabledMarkInk      draw.Ink
	barSize              float64
	initialRepeatDelay   time.Duration
	repeatDelay          time.Duration
}

func (m *managed) initialize() {
	m.backgroundInk = draw.ControlBackgroundInk
	m.focusedBackgroundInk = draw.ControlFocusedBackgroundInk
	m.pressedBackgroundInk = draw.ControlPressedBackgroundInk
	m.edgeInk = draw.ControlEdgeAdjColor
	m.markInk = draw.ControlTextColor
	m.disabledMarkInk = draw.DisabledControlTextColor
	m.barSize = 16
	m.initialRepeatDelay = time.Millisecond * 250
	m.repeatDelay = time.Millisecond * 75
}

// BackgroundInk returns the ink that will be used for the background when
// enabled but not pressed or focused.
func (s *ScrollBar) BackgroundInk() draw.Ink {
	return s.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background when
// enabled but not pressed or focused. Pass in nil to use the default.
func (s *ScrollBar) SetBackgroundInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.ControlBackgroundInk
	}
	if s.backgroundInk != value {
		s.backgroundInk = value
		s.MarkForRedraw()
	}
	return s
}

// FocusedBackgroundInk returns the ink that will be used for the background
// when enabled and focused.
func (s *ScrollBar) FocusedBackgroundInk() draw.Ink {
	return s.focusedBackgroundInk
}

// SetFocusedBackgroundInk sets the ink that will be used for the background
// when enabled and focused. Pass in nil to use the default.
func (s *ScrollBar) SetFocusedBackgroundInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.ControlFocusedBackgroundInk
	}
	if s.focusedBackgroundInk != value {
		s.focusedBackgroundInk = value
		s.MarkForRedraw()
	}
	return s
}

// PressedBackgroundInk returns the ink that will be used for the background
// when enabled and pressed.
func (s *ScrollBar) PressedBackgroundInk() draw.Ink {
	return s.pressedBackgroundInk
}

// SetPressedBackgroundInk sets the ink that will be used for the background
// when enabled and pressed. Pass in nil to use the default.
func (s *ScrollBar) SetPressedBackgroundInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.ControlPressedBackgroundInk
	}
	if s.pressedBackgroundInk != value {
		s.pressedBackgroundInk = value
		s.MarkForRedraw()
	}
	return s
}

// EdgeInk returns the ink that will be used for the edges.
func (s *ScrollBar) EdgeInk() draw.Ink {
	return s.edgeInk
}

// SetEdgeInk sets the ink that will be used for the edges. Pass in nil to
// use the default.
func (s *ScrollBar) SetEdgeInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.ControlEdgeAdjColor
	}
	if s.edgeInk != value {
		s.edgeInk = value
		s.MarkForRedraw()
	}
	return s
}

// MarkInk returns the ink that will be used for control marks when enabled.
func (s *ScrollBar) MarkInk() draw.Ink {
	return s.markInk
}

// SetMarkInk sets the ink that will be used for control marks when enabled.
// Pass in nil to use the default.
func (s *ScrollBar) SetMarkInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.ControlTextColor
	}
	if s.markInk != value {
		s.markInk = value
		s.MarkForRedraw()
	}
	return s
}

// DisabledMarkInk returns the ink that will be used for control marks when
// disabled.
func (s *ScrollBar) DisabledMarkInk() draw.Ink {
	return s.disabledMarkInk
}

// SetDisabledMarkInk sets the ink that will be used for control marks when
// disabled. Pass in nil to use the default.
func (s *ScrollBar) SetDisabledMarkInk(value draw.Ink) *ScrollBar {
	if value == nil {
		value = draw.DisabledControlTextColor
	}
	if s.disabledMarkInk != value {
		s.disabledMarkInk = value
		s.MarkForRedraw()
	}
	return s
}

// BarSize returns the height of a horizontal scrollbar or the width of a
// vertical scrollbar.
func (s *ScrollBar) BarSize() float64 {
	return s.barSize
}

// SetBarSize sets the height of a horizontal scrollbar or the width of a
// vertical scrollbar.
func (s *ScrollBar) SetBarSize(value float64) *ScrollBar {
	if value < 8 {
		value = 8
	}
	if s.barSize != value {
		s.barSize = value
		s.MarkForLayoutAndRedraw()
	}
	return s
}

// InitialRepeatDelay returns the amount of time to wait before triggering
// the first repeating event.
func (s *ScrollBar) InitialRepeatDelay() time.Duration {
	return s.initialRepeatDelay
}

// SetInitialRepeatDelay sets the amount of time to wait before triggering
// the first repeating event.
func (s *ScrollBar) SetInitialRepeatDelay(value time.Duration) *ScrollBar {
	if value < time.Millisecond*20 {
		value = time.Millisecond * 20
	}
	if s.initialRepeatDelay != value {
		s.initialRepeatDelay = value
	}
	return s
}

// RepeatDelay returns the amount of time to wait before triggering a
// repeating event.
func (s *ScrollBar) RepeatDelay() time.Duration {
	return s.repeatDelay
}

// SetRepeatDelay sets the amount of time to wait before triggering a
// repeating event.
func (s *ScrollBar) SetRepeatDelay(value time.Duration) *ScrollBar {
	if value < time.Millisecond*20 {
		value = time.Millisecond * 20
	}
	if s.repeatDelay != value {
		s.repeatDelay = value
	}
	return s
}

// SetBorder sets the border. May be nil.
func (s *ScrollBar) SetBorder(value border.Border) *ScrollBar {
	s.Panel.SetBorder(value)
	return s
}

// SetEnabled sets enabled state.
func (s *ScrollBar) SetEnabled(enabled bool) *ScrollBar {
	s.Panel.SetEnabled(enabled)
	return s
}

// SetFocusable whether it can have the keyboard focus.
func (s *ScrollBar) SetFocusable(focusable bool) *ScrollBar {
	s.Panel.SetFocusable(focusable)
	return s
}
