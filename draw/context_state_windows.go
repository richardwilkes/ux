// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import "github.com/richardwilkes/win32/d2d"

type contextState struct {
	state               *d2d.DrawingStateBlock
	clip                Path
	strokeStyle         *d2d.StrokeStyle
	strokeWidth         float32
	clipWindingFillMode bool
}

func (cs *contextState) copy(c *context) *contextState {
	state := &contextState{
		strokeStyle:         cs.strokeStyle,
		strokeWidth:         cs.strokeWidth,
		clipWindingFillMode: cs.clipWindingFillMode,
	}
	state.clip = *state.clip.Clone()
	if state.strokeStyle != nil {
		state.strokeStyle.AddRef()
	}
	return state
}

func (cs *contextState) dispose() {
	if cs.state != nil {
		cs.state.Release()
		cs.state = nil
	}
	if cs.strokeStyle != nil {
		cs.strokeStyle.Release()
		cs.strokeStyle = nil
	}
}
