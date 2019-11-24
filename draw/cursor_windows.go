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
