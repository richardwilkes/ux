package label

import (
	"math"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
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
	pref = Size(l.text, l.font, l.image, l.side, l.gap)
	if border := l.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, pref
}

// DefaultDraw provides the default drawing.
func (l *Label) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	Draw(gc, l.ContentRect(false), l.hAlign, l.vAlign, l.text, l.font, l.ink, l.image, l.side, l.gap, l.Enabled())
}

// Size returns the preferred size of a label. Provided as a standalone
// function so that other types of panels can make use of it.
func Size(text string, font *draw.Font, image *draw.Image, imgSide side.Side, imgGap float64) geom.Size {
	var size geom.Size
	if text != "" {
		size = font.Extents(text)
		size.GrowToInteger()
	}
	adjustLabelSizeForImage(text, image, imgSide, imgGap, &size)
	size.GrowToInteger()
	return size
}

func adjustLabelSizeForImage(text string, image *draw.Image, imgSide side.Side, imgGap float64, size *geom.Size) {
	if image != nil {
		logicalSize := image.LogicalGeomSize()
		switch {
		case text == "":
			*size = logicalSize
		case imgSide.Horizontal():
			size.Width += logicalSize.Width + imgGap
			if size.Height < logicalSize.Height {
				size.Height = logicalSize.Height
			}
		default:
			size.Height += logicalSize.Height + imgGap
			if size.Width < logicalSize.Width {
				size.Width = logicalSize.Width
			}
		}
	}
}

// Draw draws the label. Provided as a standalone function so that other types
// of panels can make use of it.
func Draw(gc draw.Context, rect geom.Rect, hAlign, vAlign align.Alignment, text string, font *draw.Font, textInk draw.Ink, image *draw.Image, imgSide side.Side, imgGap float64, enabled bool) {
	// Determine overall size of content
	var size, txtSize geom.Size
	if text != "" {
		txtSize = font.Extents(text)
		txtSize.GrowToInteger()
		size = txtSize
	}
	adjustLabelSizeForImage(text, image, imgSide, imgGap, &size)

	// Adjust the working area for the content size
	switch hAlign {
	case align.Middle, align.Fill:
		rect.X = math.Floor(rect.X + (rect.Width-size.Width)/2)
	case align.End:
		rect.X += rect.Width - size.Width
	default: // Start
	}
	switch vAlign {
	case align.Middle, align.Fill:
		rect.Y = math.Floor(rect.Y + (rect.Height-size.Height)/2)
	case align.End:
		rect.Y += rect.Height - size.Height
	default: // Start
	}
	rect.Size = size

	// Determine image and text areas
	imgX := rect.X
	imgY := rect.Y
	txtX := rect.X
	txtY := rect.Y
	if text != "" && image != nil {
		logicalSize := image.LogicalGeomSize()
		switch imgSide {
		case side.Top:
			txtY += logicalSize.Height + imgGap
			if logicalSize.Width > txtSize.Width {
				txtX = math.Floor(txtX + (logicalSize.Width-txtSize.Width)/2)
			} else {
				imgX = math.Floor(imgX + (txtSize.Width-logicalSize.Width)/2)
			}
		case side.Left:
			txtX += logicalSize.Width + imgGap
			if logicalSize.Height > txtSize.Height {
				txtY = math.Floor(txtY + (logicalSize.Height-txtSize.Height)/2)
			} else {
				imgY = math.Floor(imgY + (txtSize.Height-logicalSize.Height)/2)
			}
		case side.Bottom:
			imgY += rect.Height - logicalSize.Height
			txtY = imgY - (imgGap + txtSize.Height)
			if logicalSize.Width > txtSize.Width {
				txtX = math.Floor(txtX + (logicalSize.Width-txtSize.Width)/2)
			} else {
				imgX = math.Floor(imgX + (txtSize.Width-logicalSize.Width)/2)
			}
		case side.Right:
			imgX += rect.Width - logicalSize.Width
			txtX = imgX - (imgGap + txtSize.Width)
			if logicalSize.Height > txtSize.Height {
				txtY = math.Floor(txtY + (logicalSize.Height-txtSize.Height)/2)
			} else {
				imgY = math.Floor(imgY + (txtSize.Height-logicalSize.Height)/2)
			}
		}
	}

	if !enabled {
		gc.SetOpacity(0.33)
	}

	// Draw the image
	if image != nil {
		rect.X = imgX
		rect.Y = imgY
		rect.Size = image.LogicalGeomSize()
		image.DrawInRect(gc, rect)
	}

	// Draw the text
	if text != "" {
		gc.DrawString(txtX, txtY, font, textInk, text)
	}
}
