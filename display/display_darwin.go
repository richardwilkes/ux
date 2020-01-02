// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package display

import (
	"github.com/richardwilkes/macos/cg"
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

func osDisplays() []*Display {
	screens := ns.Screens()
	displays := make([]*Display, len(screens))
	for i := range screens {
		id := screens[i].DisplayID()
		x, y, width, height := cg.DisplayBounds(id)
		vx, vy, vwidth, vheight := screens[i].VisibleFrame()
		_, fy, _, fheight := screens[i].Frame()
		displays[i] = &Display{
			Frame: geom.Rect{
				Point: geom.Point{
					X: x,
					Y: y,
				},
				Size: geom.Size{
					Width:  width,
					Height: height,
				},
			},
			Usable: geom.Rect{
				Point: geom.Point{
					X: vx,
					Y: y + (fy + fheight - (vy + vheight)),
				},
				Size: geom.Size{
					Width:  vwidth,
					Height: vheight,
				},
			},
			ScalingFactor: screens[i].BackingScaleFactor(),
			Primary:       cg.DisplayIsMain(id),
		}
	}
	return displays
}
