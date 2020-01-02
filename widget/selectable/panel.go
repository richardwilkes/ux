// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package selectable

import "github.com/richardwilkes/ux"

// Panel wraps a ux.Panel for use in a Group.
type Panel struct {
	ux.Panel
	group    *Group
	selected bool
}

// AsSelectable returns the object as a selectable panel.
func (p *Panel) AsSelectable() *Panel {
	return p
}

// Selected returns true if the panel is currently selected.
func (p *Panel) Selected() bool {
	return p.selected
}

// SetSelected sets the panel's selected state.
func (p *Panel) SetSelected(selected bool) {
	if p.group != nil {
		p.group.Select(p)
	} else {
		p.setSelected(selected)
	}
}

func (p *Panel) setSelected(selected bool) {
	if p.selected != selected {
		p.selected = selected
		p.MarkForRedraw()
	}
}
