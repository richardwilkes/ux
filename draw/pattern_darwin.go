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

type osPattern = cg.Pattern

type patternCallbacks struct {
	img *Image
}

func newPatternCallbacks(img *Image) *patternCallbacks {
	return &patternCallbacks{img: img}
}

func (c *patternCallbacks) PatternDraw(gc cg.Context) {
	gc.DrawImage(0, 0, float64(c.img.LogicalWidth()), float64(c.img.LogicalHeight()), c.img.osImage())
}

func (c *patternCallbacks) PatternRelease() {
	c.img = nil
}

func osNewPattern(img *Image) osPattern {
	w := float64(img.LogicalWidth())
	h := float64(img.LogicalHeight())
	return cg.PatternCreate(0, 0, w, h, cg.AffineTransformIdentity, w, h, cg.PatternTilingConstantSpacing, true, newPatternCallbacks(img))
}

func (p *Pattern) osiPrepareForFill(gc Context) {
	g := gc.OSContext()
	patternSpace := cg.ColorSpaceCreatePattern(0)
	g.SetFillColorSpace(patternSpace)
	patternSpace.Release()
	g.SetFillPattern(p.Resource.(*patternRef).osPattern, 1)
}

func (p *Pattern) osFill(gc Context) {
	p.osiPrepareForFill(gc)
	gc.OSContext().FillPath()
}

func (p *Pattern) osFillEvenOdd(gc Context) {
	p.osiPrepareForFill(gc)
	gc.OSContext().EOFillPath()
}

func (p *Pattern) osStroke(gc Context) {
	g := gc.OSContext()
	patternSpace := cg.ColorSpaceCreatePattern(0)
	g.SetFillColorSpace(patternSpace)
	patternSpace.Release()
	g.SetStrokePattern(p.Resource.(*patternRef).osPattern, 1)
	g.StrokePath()
}

func (r *patternRef) osIsValid() bool {
	return r.osPattern != 0
}

func (r *patternRef) osDispose() {
	if r.osIsValid() {
		r.osPattern.Release()
		r.osPattern = 0
	}
}
