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
	Image     *draw.Image
	Text      string
	Font      *draw.Font      // The font to use
	TextInk   draw.Ink        // The text ink to use
	Gap       float64         // The gap to put between the image and text
	HAlign    align.Alignment // The horizontal alignment to use
	VAlign    align.Alignment // The vertical alignment to use
	ImageSide side.Side       // The side of the text the image should be on
}

// NewWithText creates a new label with the specified text.
func NewWithText(text string) *Label {
	return New(nil, text)
}

// NewWithImage creates a new label with the specified image.
func NewWithImage(image *draw.Image) *Label {
	return New(image, "")
}

// New creates a new label. Both image and text are optional.
func New(image *draw.Image, text string) *Label {
	l := &Label{
		Image:     image,
		Text:      text,
		Font:      draw.LabelFont,
		TextInk:   draw.LabelColor,
		Gap:       3,
		HAlign:    align.Start,
		VAlign:    align.Middle,
		ImageSide: side.Left,
	}
	l.InitTypeAndID(l)
	l.SetSizer(l.DefaultSizes)
	l.DrawCallback = l.DefaultDraw
	return l
}

// DefaultSizes provides the default sizing.
func (l *Label) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = Size(l.Text, l.Font, l.Image, l.ImageSide, l.Gap)
	if border := l.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, pref
}

// DefaultDraw provides the default drawing.
func (l *Label) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	Draw(gc, l.ContentRect(false), l.HAlign, l.VAlign, l.Text, l.Font, l.TextInk, l.Image, l.ImageSide, l.Gap, l.Enabled())
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
