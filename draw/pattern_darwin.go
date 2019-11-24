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

func (p *Pattern) osPrepareForFill(gc Context) {
	g := gc.OSContext()
	patternSpace := cg.ColorSpaceCreatePattern(0)
	g.SetFillColorSpace(patternSpace)
	patternSpace.Release()
	g.SetFillPattern(p.Resource.(*patternRef).osPattern, 1)
}

func (p *Pattern) osFill(gc Context) {
	p.osPrepareForFill(gc)
	gc.OSContext().FillPath()
}

func (p *Pattern) osFillEvenOdd(gc Context) {
	p.osPrepareForFill(gc)
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
