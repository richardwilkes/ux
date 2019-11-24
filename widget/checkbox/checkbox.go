package checkbox

import (
	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ux"
	"github.com/richardwilkes/ux/draw"
	"github.com/richardwilkes/ux/keys"
	"github.com/richardwilkes/ux/layout"
	"github.com/richardwilkes/ux/layout/align"
	"github.com/richardwilkes/ux/layout/side"
	"github.com/richardwilkes/ux/widget"
	"github.com/richardwilkes/ux/widget/checkbox/state"
	"github.com/richardwilkes/ux/widget/label"
	"math"
	"time"
)

// CheckBox represents a clickable checkbox with an optional label.
type CheckBox struct {
	ux.Panel
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
	Gap                  float64         // The gap to put between the checkbox, image and text
	CornerRadius         float64         // The amount of rounding to use on the corners
	ClickAnimationTime   time.Duration   // The amount of time to spend animating the click action
	ImageSide            side.Side       // The side of the text the image should be on
	HAlign               align.Alignment // The horizontal alignment to use
	VAlign               align.Alignment // The vertical alignment to use
	State                state.State
	Pressed              bool
}

// NewWithText creates a new checkbox with the specified text.
func NewWithText(text string) *CheckBox {
	return New(text, nil)
}

// NewWithImage creates a new checkbox with the specified image.
func NewWithImage(img *draw.Image) *CheckBox {
	return New("", img)
}

// New creates a new checkbox with the specified text and
// image.
func New(text string, img *draw.Image) *CheckBox {
	c := &CheckBox{
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
func (c *CheckBox) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	pref = c.boxAndLabelSize()
	if border := c.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (c *CheckBox) boxAndLabelSize() geom.Size {
	boxSize := c.boxSize()
	if c.Image == nil && c.Text == "" {
		return geom.Size{Width: boxSize, Height: boxSize}
	}
	size := label.Size(c.Text, c.Font, c.Image, c.ImageSide, c.Gap)
	size.Width += c.Gap + boxSize
	if size.Height < boxSize {
		size.Height = boxSize
	}
	return size
}

func (c *CheckBox) boxSize() float64 {
	return math.Ceil(c.Font.Height())
}

// DefaultDraw provides the default drawing.
func (c *CheckBox) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	rect := c.ContentRect(false)
	size := c.boxAndLabelSize()
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
	boxSize := c.boxSize()
	if c.Image != nil || c.Text != "" {
		r := rect
		r.X += boxSize + c.Gap
		r.Width -= boxSize + c.Gap
		label.Draw(gc, r, c.HAlign, c.VAlign, c.Text, c.Font, c.TextInk, c.Image, c.ImageSide, c.Gap, c.Enabled())
	}
	if rect.Height > boxSize {
		rect.Y += math.Floor((rect.Height - boxSize) / 2)
	}
	rect.Width = boxSize
	rect.Height = boxSize
	if !c.Enabled() {
		gc.SetOpacity(0.33)
	}
	widget.DrawRoundedRectBase(gc, rect, c.CornerRadius, c.currentBackgroundInk(), c.EdgeInk)
	rect.InsetUniform(0.5)
	switch c.State {
	case state.Mixed:
		gc.SetStrokeWidth(2)
		gc.MoveTo(rect.X+rect.Width*0.25, rect.Y+rect.Height*0.5)
		gc.LineTo(rect.X+rect.Width*0.7, rect.Y+rect.Height*0.5)
		gc.Stroke(c.currentMarkInk())
	case state.Checked:
		gc.SetStrokeWidth(2)
		gc.MoveTo(rect.X+rect.Width*0.25, rect.Y+rect.Height*0.55)
		gc.LineTo(rect.X+rect.Width*0.45, rect.Y+rect.Height*0.7)
		gc.LineTo(rect.X+rect.Width*0.75, rect.Y+rect.Height*0.3)
		gc.Stroke(c.currentMarkInk())
	}
}

func (c *CheckBox) currentBackgroundInk() draw.Ink {
	switch {
	case c.Pressed:
		return c.PressedBackgroundInk
	case c.Focused():
		return c.FocusedBackgroundInk
	default:
		return c.BackgroundInk
	}
}

func (c *CheckBox) currentMarkInk() draw.Ink {
	if c.Pressed || c.Focused() {
		return c.PressedTextInk
	}
	return c.TextInk
}

// Click makes the checkbox behave as if a user clicked on it.
func (c *CheckBox) Click() {
	c.updateState()
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

func (c *CheckBox) updateState() {
	if c.State == state.Checked {
		c.State = state.Unchecked
	} else {
		c.State = state.Checked
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (c *CheckBox) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	c.Pressed = true
	c.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (c *CheckBox) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := c.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if c.Pressed != pressed {
		c.Pressed = pressed
		c.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (c *CheckBox) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	c.Pressed = false
	c.MarkForRedraw()
	rect := c.ContentRect(false)
	if rect.ContainsPoint(where) {
		if c.ClickCallback != nil {
			c.updateState()
			c.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (c *CheckBox) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		c.Click()
		return true
	}
	return false
}
