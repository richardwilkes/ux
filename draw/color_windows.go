// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

import "github.com/richardwilkes/win32"

func (c Color) osPrepareForFill(gc Context) {
	osc := gc.(*context)
	osc.disposeBrush()
	osc.brush = win32.CreateSolidBrush(fromColorToWin32ColorRef(c))
}

func (c Color) osFill(gc Context) {
	c.osiFill(gc, win32.WINDING)
}

func (c Color) osFillEvenOdd(gc Context) {
	c.osiFill(gc, win32.ALTERNATE)
}

func (c Color) osiFill(gc Context, mode int) {
	c.osPrepareForFill(gc)
	osc := gc.(*context)
	win32.EndPath(osc.hdc)
	win32.SetPolyFillMode(osc.hdc, mode)
	win32.SelectObject(osc.hdc, win32.HGDIOBJ(osc.brush))
	win32.FillPath(osc.hdc)
}

func (c Color) osStroke(gc Context) {
	osc := gc.(*context)
	osc.disposePen()
	osc.pen = win32.CreatePen(win32.PS_SOLID, 1, fromColorToWin32ColorRef(c))
	win32.EndPath(osc.hdc)
	win32.SelectObject(osc.hdc, win32.HGDIOBJ(osc.pen))
	win32.StrokePath(osc.hdc)
}
