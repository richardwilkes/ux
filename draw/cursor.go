package draw

import "github.com/richardwilkes/toolbox/xmath/geom"

// Available system cursors
var (
	ArrowCursor            *Cursor
	TextCursor             *Cursor
	VerticalTextCursor     *Cursor
	CrossHairCursor        *Cursor
	ClosedHandCursor       *Cursor
	OpenHandCursor         *Cursor
	PointingHandCursor     *Cursor
	ResizeLeftCursor       *Cursor
	ResizeRightCursor      *Cursor
	ResizeLeftRightCursor  *Cursor
	ResizeUpCursor         *Cursor
	ResizeDownCursor       *Cursor
	ResizeUpDownCursor     *Cursor
	DisappearingItemCursor *Cursor
	NotAllowedCursor       *Cursor
	DragLinkCursor         *Cursor
	DragCopyCursor         *Cursor
	ContextMenuCursor      *Cursor
)

// Cursor provides a graphical cursor for the mouse location.
type Cursor struct {
	cursor osCursor
}

// NewCursor creates a new custom cursor from an image.
func NewCursor(img *Image, hotSpot geom.Point) *Cursor {
	return &Cursor{cursor: osNewCursor(img, hotSpot)}
}

// HideCursorUntilMouseMoves hides the cursor until the mouse is moved.
func HideCursorUntilMouseMoves() {
	osHideCursorUntilMouseMoves()
}

// MakeCurrent makes this cursor the current cursor.
func (c *Cursor) MakeCurrent() {
	c.osMakeCurrent()
}
