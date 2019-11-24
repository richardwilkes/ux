package radiobutton

import (
	"math"
	"time"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/label"
	"github.com/richardwilkes/ux/widget/selectable"
)

// RadioButton represents a clickable radio button with an optional label.
type RadioButton struct {
	selectable.Panel
	ClickCallback        func()
	Text                 string
	Image                *draw.Image
	Font                 *draw.Font      // The font to use
	BackgroundInk        draw.Ink        // The background ink when enabled but not pressed or focused
	FocusedBackgroundInk draw.Ink        // The background ink when enabled and focused
	PressedBackgroundInk draw.Ink        // The background ink when enabled and pressed
	EdgeInk              draw.Ink        // The ink to use on the edges
	TextInk              draw.Ink        // The text ink to use
	PressedTextInk       draw.Ink        // The text ink when enabled and pressed
	Gap                  float64         // The gap to put between the radio button, image and text
	CornerRadius         float64         // The amount of rounding to use on the corners
	ClickAnimationTime   time.Duration   // The amount of time to spend animating the click action
	ImageSide            side.Side       // The side of the text the image should be on
	HAlign               align.Alignment // The horizontal alignment to use
	VAlign               align.Alignment // The vertical alignment to use
	Pressed              bool
}

// NewWithText creates a new radio button with the specified text.
func NewWithText(text string) *RadioButton {
	return New(text, nil)
}

// NewWithImage creates a new radio button with the specified image.
func NewWithImage(img *draw.Image) *RadioButton {
	return New("", img)
}

// New creates a new radio button with the specified
// text and image.
func New(text string, img *draw.Image) *RadioButton {
	c := &RadioButton{
		Text:                 text,
		Image:                img,
		Font:                 draw.SystemFont,
		BackgroundInk:        draw.ControlBackgroundInk,
		FocusedBackgroundInk: draw.ControlFocusedBackgroundInk,
		PressedBackgroundInk: draw.ControlPressedBackgroundInk,
		EdgeInk:              draw.ControlEdgeAdjColor,
		TextInk:              draw.ControlTextColor,
		PressedTextInk:       draw.AlternateSelectedControlTextColor,
		Gap:                  3,
		CornerRadius:         4,
		ClickAnimationTime:   time.Millisecond * 100,
		ImageSide:            side.Left,
		HAlign:               align.Start,
		VAlign:               align.Middle,
	}
	c.InitTypeAndID(c)
	c.SetFocusable(true)
	c.SetSizer(c.DefaultSizes)
	c.DrawCallback = c.DefaultDraw
	c.GainedFocusCallback = c.MarkForRedraw
	c.LostFocusCallback = c.MarkForRedraw
	c.MouseDownCallback = c.DefaultMouseDown
	c.MouseDragCallback = c.DefaultMouseDrag
	c.MouseUpCallback = c.DefaultMouseUp
	c.KeyDownCallback = c.DefaultKeyDown
	return c
}

// DefaultSizes provides the default sizing.
func (c *RadioButton) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = c.circleAndLabelSize()
	if border := c.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (c *RadioButton) circleAndLabelSize() geom.Size {
	circleSize := c.circleSize()
	if c.Image == nil && c.Text == "" {
		return geom.Size{Width: circleSize, Height: circleSize}
	}
	size := label.Size(c.Text, c.Font, c.Image, c.ImageSide, c.Gap)
	size.Width += c.Gap + circleSize
	if size.Height < circleSize {
		size.Height = circleSize
	}
	return size
}

func (c *RadioButton) circleSize() float64 {
	return math.Ceil(c.Font.Height())
}

// DefaultDraw provides the default drawing.
func (c *RadioButton) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !c.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := c.ContentRect(false)
	size := c.circleAndLabelSize()
	switch c.HAlign {
	case align.Middle, align.Fill:
		rect.X = math.Floor(rect.X + (rect.Width-size.Width)/2)
	case align.End:
		rect.X += rect.Width - size.Width
	default: // Start
	}
	switch c.VAlign {
	case align.Middle, align.Fill:
		rect.Y = math.Floor(rect.Y + (rect.Height-size.Height)/2)
	case align.End:
		rect.Y += rect.Height - size.Height
	default: // Start
	}
	rect.Size = size
	circleSize := c.circleSize()
	if c.Image != nil || c.Text != "" {
		r := rect
		r.X += circleSize + c.Gap
		r.Width -= circleSize + c.Gap
		label.Draw(gc, r, c.HAlign, c.VAlign, c.Text, c.Font, c.TextInk, c.Image, c.ImageSide, c.Gap, c.Enabled())
	}
	if rect.Height > circleSize {
		rect.Y += math.Floor((rect.Height - circleSize) / 2)
	}
	rect.Width = circleSize
	rect.Height = circleSize
	widget.DrawEllipseBase(gc, rect, c.currentBackgroundInk(), c.EdgeInk)
	if c.Selected() {
		rect.InsetUniform(0.5 + 0.2*circleSize)
		gc.Ellipse(rect)
		gc.Fill(c.currentMarkInk())
	}
}

func (c *RadioButton) currentBackgroundInk() draw.Ink {
	switch {
	case c.Pressed:
		return c.PressedBackgroundInk
	case c.Focused():
		return c.FocusedBackgroundInk
	default:
		return c.BackgroundInk
	}
}

func (c *RadioButton) currentMarkInk() draw.Ink {
	if c.Pressed || c.Focused() {
		return c.PressedTextInk
	}
	return c.TextInk
}

// Click makes the radio button behave as if a user clicked on it.
func (c *RadioButton) Click() {
	c.SetSelected(true)
	pressed := c.Pressed
	c.Pressed = true
	c.MarkForRedraw()
	c.FlushDrawing()
	c.Pressed = pressed
	time.Sleep(c.ClickAnimationTime)
	c.MarkForRedraw()
	if c.ClickCallback != nil {
		c.ClickCallback()
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (c *RadioButton) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	c.Pressed = true
	c.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (c *RadioButton) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := c.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if c.Pressed != pressed {
		c.Pressed = pressed
		c.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (c *RadioButton) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	c.Pressed = false
	c.MarkForRedraw()
	rect := c.ContentRect(false)
	if rect.ContainsPoint(where) {
		c.SetSelected(true)
		if c.ClickCallback != nil {
			c.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (c *RadioButton) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		c.Click()
		return true
	}
	return false
}
