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
	"github.com/richardwilkes/toolbox/xmath/geom"
)

// Display holds information about each available active display.
type Display struct {
	Frame         geom.Rect
	Usable        geom.Rect
	ScalingFactor float64
	Primary       bool
}

// Primary returns the primary display.
func Primary() *Display {
	for _, d := range All() {
		if d.Primary {
			return d
		}
	}
	return nil
}

// All returns all displays.
func All() []*Display {
	return osDisplays()
}
