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

package list

import (
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/label"
)

type managed struct {
	factory                widget.CellFactory
	backgroundInk          draw.Ink
	alternateBackgroundInk draw.Ink
	selectedBackgroundInk  draw.Ink
}

func (m *managed) initialize() {
	m.factory = &label.CellFactory{}
	m.backgroundInk = draw.TextBackgroundColor
	m.alternateBackgroundInk = draw.TextAlternateBackgroundColor
	m.selectedBackgroundInk = draw.ControlAccentColor
}

// Factory returns the cell factory.
func (l *List) Factory() widget.CellFactory {
	return l.factory
}

// SetFactory sets the cell factory. Pass in nil to use the default.
func (l *List) SetFactory(value widget.CellFactory) *List {
	if value == nil {
		value = &label.CellFactory{}
	}
	if l.factory != value {
		l.factory = value
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// BackgroundInk returns the ink that will be used for the background on even
// rows when not selected.
func (l *List) BackgroundInk() draw.Ink {
	return l.backgroundInk
}

// SetBackgroundInk sets the ink that will be used for the background on even
// rows when not selected. Pass in nil to use the default.
func (l *List) SetBackgroundInk(value draw.Ink) *List {
	if value == nil {
		value = draw.TextBackgroundColor
	}
	if l.backgroundInk != value {
		l.backgroundInk = value
		l.MarkForRedraw()
	}
	return l
}

// AlternateBackgroundInk returns the ink that will be used for the
// background on odd rows when not selected.
func (l *List) AlternateBackgroundInk() draw.Ink {
	return l.alternateBackgroundInk
}

// SetAlternateBackgroundInk sets the ink that will be used for the
// background on odd rows when not selected. Pass in nil to use the default.
func (l *List) SetAlternateBackgroundInk(value draw.Ink) *List {
	if value == nil {
		value = draw.TextAlternateBackgroundColor
	}
	if l.alternateBackgroundInk != value {
		l.alternateBackgroundInk = value
		l.MarkForRedraw()
	}
	return l
}

// SelectedBackgroundInk returns the ink that will be used for the background
// when selected.
func (l *List) SelectedBackgroundInk() draw.Ink {
	return l.selectedBackgroundInk
}

// SetSelectedBackgroundInk sets the ink that will be used for the background
// when selected. Pass in nil to use the default.
func (l *List) SetSelectedBackgroundInk(value draw.Ink) *List {
	if value == nil {
		value = draw.ControlAccentColor
	}
	if l.selectedBackgroundInk != value {
		l.selectedBackgroundInk = value
		l.MarkForRedraw()
	}
	return l
}

// SetBorder sets the border. May be nil.
func (l *List) SetBorder(value border.Border) *List {
	l.Panel.SetBorder(value)
	return l
}

// SetEnabled sets enabled state.
func (l *List) SetEnabled(enabled bool) *List {
	l.Panel.SetEnabled(enabled)
	return l
}

// SetFocusable whether it can have the keyboard focus.
func (l *List) SetFocusable(focusable bool) *List {
	l.Panel.SetFocusable(focusable)
	return l
}
