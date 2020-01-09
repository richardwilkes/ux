// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Code created from "widget.go.tmpl" - don't edit by hand

package scrollarea

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	backgroundInk   draw.Ink
	focusedBorder   border.Border
	unfocusedBorder border.Border
}

func (m *managed) initialize() {
	m.backgroundInk = draw.TextBackgroundColor
	m.focusedBorder = border.NewCompound(border.NewLine(draw.ControlAccentColor, 0, geom.NewUniformInsets(1), false), border.NewLine(draw.ControlAccentColor, 0, geom.NewUniformInsets(1), true))
	m.unfocusedBorder = border.NewCompound(border.NewLine(draw.ControlEdgeAdjColor, 0, geom.NewUniformInsets(1), false), border.NewLine(draw.ARGB(0, 0, 0, 0), 0, geom.NewUniformInsets(1), true))
}

// BackgroundInk returns the ink that will be used for the background.
func (s *ScrollArea) BackgroundInk() draw.Ink {
	return s.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background. Pass
// in nil to use the default.
func (s *ScrollArea) SetBackgroundInk(value draw.Ink) *ScrollArea {
	if value == nil {
		value = draw.TextBackgroundColor
	}
	if s.backgroundInk != value {
		s.backgroundInk = value
		s.MarkForRedraw()
	}
	return s
}

// FocusedBorder returns the border to use when focused. Note that the border
// should present the same insets as the unfocused border or the display will
// not appear correct.
func (s *ScrollArea) FocusedBorder() border.Border {
	return s.focusedBorder
}

// SetFocusedBorder sets the border to use when focused. Note that the border
// should present the same insets as the unfocused border or the display will
// not appear correct. Pass in nil to use the default.
func (s *ScrollArea) SetFocusedBorder(value border.Border) *ScrollArea {
	if value == nil {
		value = border.NewCompound(border.NewLine(draw.ControlAccentColor, 0, geom.NewUniformInsets(1), false), border.NewLine(draw.ControlAccentColor, 0, geom.NewUniformInsets(1), true))
	}
	if s.focusedBorder != value {
		s.focusedBorder = value
		s.MarkForLayoutAndRedraw()
	}
	return s
}

// UnfocusedBorder returns the border to use when not focused. Note that the
// border should present the same insets as the focused border or the display
// will not appear correct.
func (s *ScrollArea) UnfocusedBorder() border.Border {
	return s.unfocusedBorder
}

// SetUnfocusedBorder sets the border to use when not focused. Note that the
// border should present the same insets as the focused border or the display
// will not appear correct. Pass in nil to use the default.
func (s *ScrollArea) SetUnfocusedBorder(value border.Border) *ScrollArea {
	if value == nil {
		value = border.NewCompound(border.NewLine(draw.ControlEdgeAdjColor, 0, geom.NewUniformInsets(1), false), border.NewLine(draw.ARGB(0, 0, 0, 0), 0, geom.NewUniformInsets(1), true))
	}
	if s.unfocusedBorder != value {
		s.unfocusedBorder = value
		s.MarkForLayoutAndRedraw()
	}
	return s
}

// SetBorder sets the border. May be nil.
func (s *ScrollArea) SetBorder(value border.Border) *ScrollArea {
	s.Panel.SetBorder(value)
	return s
}

// SetEnabled sets enabled state.
func (s *ScrollArea) SetEnabled(enabled bool) *ScrollArea {
	s.Panel.SetEnabled(enabled)
	return s
}

// SetFocusable whether it can have the keyboard focus.
func (s *ScrollArea) SetFocusable(focusable bool) *ScrollArea {
	s.Panel.SetFocusable(focusable)
	return s
}
