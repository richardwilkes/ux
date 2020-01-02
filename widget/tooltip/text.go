// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package tooltip

import (
	"strings"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/flex"
	"github.com/richardwilkes/ux/widget/label"
)

// NewBase returns the base for a tooltip.
func NewBase() *ux.Panel {
	tip := ux.NewPanel()
	tip.SetBorder(border.NewCompound(border.NewLine(draw.UnemphasizedSelectedContentBackgroundColor, 0, geom.NewUniformInsets(1), false), border.NewEmpty(geom.Insets{Top: 2, Left: 4, Bottom: 2, Right: 4})))
	tip.DrawCallback = func(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
		gc.Rect(dirty)
		gc.Fill(draw.ControlBackgroundColor)
	}
	return tip
}

// NewWithText creates a standard text tooltip panel.
func NewWithText(text string) *ux.Panel {
	tip := NewBase()
	flex.New().Apply(tip)
	for _, str := range strings.Split(text, "\n") {
		l := label.New().SetText(str).SetFont(draw.SystemFont)
		tip.AddChild(l.AsPanel())
	}
	return tip
}

// NewWithSecondaryText creates a text tooltip panel containing a primary
// piece of text along with a secondary piece of text in a slightly smaller
// font.
func NewWithSecondaryText(primary, secondary string) *ux.Panel {
	tip := NewBase()
	flex.New().Apply(tip)
	for _, str := range strings.Split(primary, "\n") {
		l := label.New().SetText(str).SetFont(draw.SystemFont)
		tip.AddChild(l.AsPanel())
	}
	if secondary != "" {
		for _, str := range strings.Split(secondary, "\n") {
			l := label.New().SetText(str).SetFont(draw.SmallSystemFont)
			tip.AddChild(l.AsPanel())
		}
	}
	return tip
}
