// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package border

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
)

// Border defines methods required of all border providers.
type Border interface {
	// Insets returns the insets describing the space the border occupies on
	// each side.
	Insets() geom.Insets
	// Draw the border into rect.
	Draw(gc draw.Context, rect geom.Rect, inLiveResize bool)
}
