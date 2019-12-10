package label

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
	"github.com/richardwilkes/ux/widget"
)

// Label represents non-interactive text and/or an image.
type Label struct {
	ux.Panel
	image  *draw.Image
	text   string
	font   *draw.Font
	ink    draw.Ink
	gap    float64
	hAlign align.Alignment
	vAlign align.Alignment
	side   side.Side
}

// New creates a new, empty label.
func New() *Label {
	l := &Label{
		font:   draw.LabelFont,
		ink:    draw.LabelColor,
		gap:    3,
		hAlign: align.Start,
		vAlign: align.Middle,
		side:   side.Left,
	}
	l.InitTypeAndID(l)
	l.SetSizer(l.DefaultSizes)
	l.DrawCallback = l.DefaultDraw
	return l
}

// Image returns the image. May be nil.
func (l *Label) Image() *draw.Image {
	return l.image
}

// SetImage sets the image. May be nil.
func (l *Label) SetImage(image *draw.Image) *Label {
	if l.image != image {
		l.image = image
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Text returns the text content.
func (l *Label) Text() string {
	return l.text
}

// SetText sets the text content.
func (l *Label) SetText(text string) *Label {
	if l.text != text {
		l.text = text
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Font returns the font that will be used when drawing text content.
func (l *Label) Font() *draw.Font {
	return l.font
}

// SetFont sets the font to use when drawing text content.
func (l *Label) SetFont(font *draw.Font) *Label {
	if font == nil {
		font = draw.LabelFont
	}
	if l.font != font {
		l.font = font
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Ink returns the ink that will be used when drawing text content.
func (l *Label) Ink() draw.Ink {
	return l.ink
}

// SetInk sets the ink to use when drawing text content.
func (l *Label) SetInk(ink draw.Ink) *Label {
	if ink == nil {
		ink = draw.LabelColor
	}
	if l.ink != ink {
		l.ink = ink
		l.MarkForRedraw()
	}
	return l
}

// Gap returns the gap to put between the image and text.
func (l *Label) Gap() float64 {
	return l.gap
}

// SetGap sets the gap to put between the image and text.
func (l *Label) SetGap(gap float64) *Label {
	if gap < 0 {
		gap = 0
	}
	if l.gap != gap {
		l.gap = gap
		if l.image != nil && l.text != "" {
			l.MarkForLayoutAndRedraw()
		}
	}
	return l
}

// HAlign returns the horizontal alignment.
func (l *Label) HAlign() align.Alignment {
	return l.hAlign
}

// SetHAlign sets the horizontal alignment.
func (l *Label) SetHAlign(hAlign align.Alignment) *Label {
	if l.hAlign != hAlign {
		l.hAlign = hAlign
		l.MarkForRedraw()
	}
	return l
}

// VAlign returns the vertical alignment.
func (l *Label) VAlign() align.Alignment {
	return l.hAlign
}

// SetVAlign sets the vertical alignment.
func (l *Label) SetVAlign(vAlign align.Alignment) *Label {
	if l.vAlign != vAlign {
		l.vAlign = vAlign
		l.MarkForRedraw()
	}
	return l
}

// Side returns the side of the text the image should be on.
func (l *Label) Side() side.Side {
	return l.side
}

// SetSide sets the side of the text the image should be on.
func (l *Label) SetSide(s side.Side) *Label {
	if l.side != s {
		l.side = s
		l.MarkForRedraw()
	}
	return l
}

// DefaultSizes provides the default sizing.
func (l *Label) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = widget.LabelSize(l.text, l.font, l.image, l.side, l.gap)
	if border := l.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, pref
}

// DefaultDraw provides the default drawing.
func (l *Label) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	widget.DrawLabel(gc, l.ContentRect(false), l.hAlign, l.vAlign, l.text, l.font, l.ink, l.image, l.side, l.gap, l.Enabled())
}
