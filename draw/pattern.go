package draw

import (
	"github.com/richardwilkes/toolbox/softref"
	"github.com/richardwilkes/ux/globals"
)

var _ Ink = &Pattern{}

// Pattern holds a pattern for drawing.
type Pattern softref.SoftRef

type patternRef struct {
	key       string
	img       *Image
	osPattern osPattern
}

func (r *patternRef) Key() string {
	return r.key
}

func (r *patternRef) Release() {
	r.osDispose()
	r.img = nil
}

// NewPattern creates a new pattern from an image.
func NewPattern(img *Image) *Pattern {
	r := &patternRef{
		key:       "p:" + img.Key,
		img:       img,
		osPattern: osNewPattern(img),
	}
	ref, existedPreviously := globals.Pool.NewSoftRef(r)
	if existedPreviously {
		r.Release()
	}
	return (*Pattern)(ref)
}

// IsValid returns true if the pattern is still valid (i.e. hasn't been
// disposed).
func (p *Pattern) IsValid() bool {
	return p.Resource.(*patternRef).osIsValid()
}

// Image returns the underlying image.
func (p *Pattern) Image() *Image {
	if p.IsValid() {
		return p.Resource.(*patternRef).img
	}
	return nil
}
