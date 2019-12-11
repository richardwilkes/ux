// Code created from "widget.go.tmpl" - don't edit by hand

package label

import (
	"github.com/richardwilkes/ux/border"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
)

type managed struct {
	image  *draw.Image //nolint:structcheck
	text   string      //nolint:structcheck
	font   *draw.Font
	ink    draw.Ink
	gap    float64
	hAlign align.Alignment
	vAlign align.Alignment
	side   side.Side
}

func (m *managed) initialize() {
	m.font = draw.LabelFont
	m.ink = draw.LabelColor
	m.gap = 3
	m.hAlign = align.Start
	m.vAlign = align.Middle
	m.side = side.Left
}

// Image returns the image. May be nil.
func (l *Label) Image() *draw.Image {
	return l.image
}

// SetImage sets the image. May be nil.
func (l *Label) SetImage(value *draw.Image) *Label {
	if l.image != value {
		l.image = value
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Text returns the text content.
func (l *Label) Text() string {
	return l.text
}

// SetText sets the text content.
func (l *Label) SetText(value string) *Label {
	if l.text != value {
		l.text = value
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Font returns the font that will be used when drawing text content.
func (l *Label) Font() *draw.Font {
	return l.font
}

// SetFont sets the font that will be used when drawing text content. Pass in
// nil to use the default.
func (l *Label) SetFont(value *draw.Font) *Label {
	if value == nil {
		value = draw.LabelFont
	}
	if l.font != value {
		l.font = value
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// Ink returns the ink that will be used when drawing text content.
func (l *Label) Ink() draw.Ink {
	return l.ink
}

// SetInk sets the ink that will be used when drawing text content. Pass in
// nil to use the default.
func (l *Label) SetInk(value draw.Ink) *Label {
	if value == nil {
		value = draw.LabelColor
	}
	if l.ink != value {
		l.ink = value
		l.MarkForRedraw()
	}
	return l
}

// Gap returns the gap to put between the image and text.
func (l *Label) Gap() float64 {
	return l.gap
}

// SetGap sets the gap to put between the image and text.
func (l *Label) SetGap(value float64) *Label {
	if value < 0 {
		value = 0
	}
	if l.gap != value {
		l.gap = value
		l.MarkForLayoutAndRedraw()
	}
	return l
}

// HAlign returns the horizontal alignment.
func (l *Label) HAlign() align.Alignment {
	return l.hAlign
}

// SetHAlign sets the horizontal alignment.
func (l *Label) SetHAlign(value align.Alignment) *Label {
	if l.hAlign != value {
		l.hAlign = value
		l.MarkForRedraw()
	}
	return l
}

// VAlign returns the vertical alignment.
func (l *Label) VAlign() align.Alignment {
	return l.vAlign
}

// SetVAlign sets the vertical alignment.
func (l *Label) SetVAlign(value align.Alignment) *Label {
	if l.vAlign != value {
		l.vAlign = value
		l.MarkForRedraw()
	}
	return l
}

// Side returns the side of the text the image should be on.
func (l *Label) Side() side.Side {
	return l.side
}

// SetSide sets the side of the text the image should be on.
func (l *Label) SetSide(value side.Side) *Label {
	if l.side != value {
		l.side = value
		l.MarkForRedraw()
	}
	return l
}

// SetBorder sets the border. May be nil.
func (l *Label) SetBorder(value border.Border) *Label {
	l.Panel.SetBorder(value)
	return l
}

// SetEnabled sets enabled state.
func (l *Label) SetEnabled(enabled bool) *Label {
	l.Panel.SetEnabled(enabled)
	return l
}

// SetFocusable whether it can have the keyboard focus.
func (l *Label) SetFocusable(focusable bool) *Label {
	l.Panel.SetFocusable(focusable)
	return l
}
