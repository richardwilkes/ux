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
	"github.com/richardwilkes/macos/ns"
	"github.com/richardwilkes/toolbox/xmath/geom"
)

type osCursor struct {
	cursor *ns.Cursor
}

func osInitSystemCursors() {
	ArrowCursor = &Cursor{cursor: osCursor{cursor: ns.ArrowCursor()}}
	TextCursor = &Cursor{cursor: osCursor{cursor: ns.IBeamCursor()}}
	VerticalTextCursor = &Cursor{cursor: osCursor{cursor: ns.IBeamCursorForVerticalLayout()}}
	CrossHairCursor = &Cursor{cursor: osCursor{cursor: ns.CrosshairCursor()}}
	ClosedHandCursor = &Cursor{cursor: osCursor{cursor: ns.ClosedHandCursor()}}
	OpenHandCursor = &Cursor{cursor: osCursor{cursor: ns.OpenHandCursor()}}
	PointingHandCursor = &Cursor{cursor: osCursor{cursor: ns.PointingHandCursor()}}
	ResizeLeftCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeLeftCursor()}}
	ResizeRightCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeRightCursor()}}
	ResizeLeftRightCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeLeftRightCursor()}}
	ResizeUpCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeUpCursor()}}
	ResizeDownCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeDownCursor()}}
	ResizeUpDownCursor = &Cursor{cursor: osCursor{cursor: ns.ResizeUpDownCursor()}}
	DisappearingItemCursor = &Cursor{cursor: osCursor{cursor: ns.DisappearingItemCursor()}}
	NotAllowedCursor = &Cursor{cursor: osCursor{cursor: ns.OperationNotAllowedCursor()}}
	DragLinkCursor = &Cursor{cursor: osCursor{cursor: ns.DragLinkCursor()}}
	DragCopyCursor = &Cursor{cursor: osCursor{cursor: ns.DragCopyCursor()}}
	ContextMenuCursor = &Cursor{cursor: osCursor{cursor: ns.ContextualMenuCursor()}}
}

func osNewCursor(img *Image, hotSpot geom.Point) osCursor {
	return osCursor{cursor: ns.CursorInitWithImageHotSpotRetain(ns.ImageInitWithCGImageSizeRetain(img.osImage(), float64(img.LogicalWidth()), float64(img.LogicalHeight())), hotSpot.X, hotSpot.Y)}
}

func osHideCursorUntilMouseMoves() {
	ns.CursorSetHiddenUntilMouseMoves(true)
}

func (c *Cursor) osMakeCurrent() {
	if c.cursor.cursor != nil {
		c.cursor.cursor.Set()
	}
}
