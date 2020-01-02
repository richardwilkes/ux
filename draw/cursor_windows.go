// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
)

type osCursor struct {
}

func osInitSystemCursors() {
	// RAW: Implement
}

func osNewCursor(img *Image, hotSpot geom.Point) osCursor {
	// RAW: Implement
	return osCursor{}
}

func osHideCursorUntilMouseMoves() {
	// RAW: Implement
}

func (c *Cursor) osMakeCurrent() {
	// RAW: Implement
}
