// Copyright ©2019-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package draw

func (c Color) osiPrepareForFill(gc Context) {
	gc.OSContext().SetRGBFillColor(c.RedIntensity(), c.GreenIntensity(), c.BlueIntensity(), c.AlphaIntensity())
}

func (c Color) osFill(gc Context) {
	c.osiPrepareForFill(gc)
	gc.OSContext().FillPath()
}

func (c Color) osFillEvenOdd(gc Context) {
	c.osiPrepareForFill(gc)
	gc.OSContext().EOFillPath()
}

func (c Color) osStroke(gc Context) {
	g := gc.OSContext()
	g.SetRGBStrokeColor(c.RedIntensity(), c.GreenIntensity(), c.BlueIntensity(), c.AlphaIntensity())
	g.StrokePath()
}
