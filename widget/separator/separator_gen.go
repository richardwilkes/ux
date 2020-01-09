// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Code created from "widget.go.tmpl" - don't edit by hand

package separator

import (
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
)

type managed struct {
	fillInk draw.Ink
}

func (m *managed) initialize() {
	m.fillInk = draw.SeparatorColor
}

// FillInk returns the ink that will be used for the separator line.
func (s *Separator) FillInk() draw.Ink {
	return s.fillInk
}

// SetFillInk sets the ink that will be used for the separator line. Pass in
// nil to use the default.
func (s *Separator) SetFillInk(value draw.Ink) *Separator {
	if value == nil {
		value = draw.SeparatorColor
	}
	if s.fillInk != value {
		s.fillInk = value
		s.MarkForRedraw()
	}
	return s
}

// SetBorder sets the border. May be nil.
func (s *Separator) SetBorder(value border.Border) *Separator {
	s.Panel.SetBorder(value)
	return s
}

// SetEnabled sets enabled state.
func (s *Separator) SetEnabled(enabled bool) *Separator {
	s.Panel.SetEnabled(enabled)
	return s
}

// SetFocusable whether it can have the keyboard focus.
func (s *Separator) SetFocusable(focusable bool) *Separator {
	s.Panel.SetFocusable(focusable)
	return s
}
