package label

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/widget"
)

// Label represents non-interactive text and/or an image.
type Label struct {
	ux.Panel
	managed
}

// New creates a new, empty label.
func New() *Label {
	l := &Label{}
	l.managed.initialize()
	l.InitTypeAndID(l)
	l.SetSizer(l.DefaultSizes)
	l.DrawCallback = l.DefaultDraw
	return l
}

// DefaultSizes provides the default sizing.
func (l *Label) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = widget.LabelSize(l.text, l.font, l.image, l.side, l.gap)
	if b := l.Border(); b != nil {
		pref.AddInsets(b.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, pref
}

// DefaultDraw provides the default drawing.
func (l *Label) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	widget.DrawLabel(gc, l.ContentRect(false), l.hAlign, l.vAlign, l.text, l.font, l.ink, l.image, l.side, l.gap, l.Enabled())
}
