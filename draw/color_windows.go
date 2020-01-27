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
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/win32/d2d"
)

func (c Color) osPrepareForFill(gc Context) {
	// Not supported
}

func (c Color) osFill(gc Context) {
	c.osiFill(gc, true)
}

func (c Color) osFillEvenOdd(gc Context) {
	c.osiFill(gc, false)
}

func (c Color) osiFill(gc Context, windingFillMode bool) {
	osc := gc.(*context)
	if len(osc.path.nodes) > 0 {
		brush := c.osiNewBrush(osc)
		if brush == nil {
			jot.Error(errs.New("unable to create brush"))
			return
		}
		defer brush.Release()
		p, err := newWinPath(osc, windingFillMode, false)
		if err != nil {
			jot.Error(err)
			return
		}
		osc.path.SendPath(p)
		p.gc.OSContext().FillGeometry(p.geometry(), &brush.Brush, nil)
		p.dispose()
	}
}

func (c Color) osiNewBrush(gc *context) *d2d.SolidColorBrush {
	return gc.OSContext().CreateSolidColorBrush(d2d.Color{
		R: float32(c.RedIntensity()),
		G: float32(c.GreenIntensity()),
		B: float32(c.BlueIntensity()),
		A: float32(c.AlphaIntensity()),
	}, nil)
}

func (c Color) osStroke(gc Context) {
	osc := gc.(*context)
	if len(osc.path.nodes) > 0 {
		brush := c.osiNewBrush(osc)
		if brush == nil {
			jot.Error(errs.New("unable to create brush"))
			return
		}
		defer brush.Release()
		p, err := newWinPath(osc, true, true)
		if err != nil {
			jot.Error(err)
			return
		}
		osc.path.SendPath(p)
		current := p.gc.current()
		p.gc.OSContext().DrawGeometry(p.geometry(), &brush.Brush, current.strokeWidth, current.strokeStyle)
		p.dispose()
	}
}
