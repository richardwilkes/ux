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
	"github.com/richardwilkes/macos/cg"
)

func (g *Gradient) osPrepareForFill(gc Context) {
	// Unsupported
}

func (g *Gradient) osFill(gc Context) {
	gc.Save()
	gc.Clip()
	g.draw(gc)
	gc.Restore()
}

func (g *Gradient) osFillEvenOdd(gc Context) {
	gc.Save()
	gc.ClipEvenOdd()
	g.draw(gc)
	gc.Restore()
}

func (g *Gradient) osStroke(gc Context) {
	gc.Save()
	gc.OSContext().ReplacePathWithStrokedPath()
	gc.Clip()
	g.draw(gc)
	gc.Restore()
}

func (g *Gradient) draw(gc Context) {
	rect := gc.GetClipRect()
	colorSpace := cg.ColorSpaceCreateDeviceRGB()
	count := len(g.Stops)
	components := make([]float64, count*4)
	locations := make([]float64, count)
	for i, one := range g.Stops {
		j := i * 4
		components[j] = one.Color.Color.RedIntensity()
		components[j+1] = one.Color.Color.GreenIntensity()
		components[j+2] = one.Color.Color.BlueIntensity()
		components[j+3] = one.Color.Color.AlphaIntensity()
		locations[i] = one.Location
	}
	gradient := cg.GradientCreateWithColorComponents(colorSpace, components, locations)
	gradient.Retain()
	colorSpace.Release()
	sx := rect.X + rect.Width*g.Start.X
	sy := rect.Y + rect.Height*g.Start.Y
	ex := rect.X + rect.Width*g.End.X
	ey := rect.Y + rect.Height*g.End.Y
	if g.StartRadius > 0 && g.EndRadius > 0 {
		gc.OSContext().DrawRadialGradient(gradient, sx, sy, g.StartRadius, ex, ey, g.EndRadius, cg.GradientDrawsBeforeStartLocation|cg.GradientDrawsAfterEndLocation)
	} else {
		gc.OSContext().DrawLinearGradient(gradient, sx, sy, ex, ey, cg.GradientDrawsBeforeStartLocation|cg.GradientDrawsAfterEndLocation)
	}
	gradient.Release()
}
