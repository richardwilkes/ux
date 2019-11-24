package button

import (
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

// Button represents a clickable button.
type Button struct {
	selectable.Panel
	ClickCallback         func()
	Image                 *draw.Image
	Text                  string
	Font                  *draw.Font      // The font to use
	BackgroundInk         draw.Ink        // The background ink when enabled but not selected, pressed or focused
	SelectedBackgroundInk draw.Ink        // The background ink when enabled and selected, but not pressed or focused
	FocusedBackgroundInk  draw.Ink        // The background ink when enabled and focused
	PressedBackgroundInk  draw.Ink        // The background ink when enabled and pressed
	EdgeInk               draw.Ink        // The ink to use on the edges
	TextInk               draw.Ink        // The text ink to use
	PressedTextInk        draw.Ink        // The text ink when enabled and pressed
	ImageGap              float64         // The gap to put between the image and text
	CornerRadius          float64         // The amount of rounding to use on the corners
	HMargin               float64         // The margin on the left and right side of the content
	VMargin               float64         // The margin on the top and bottom of the content
	ImageOnlyHMargin      float64         // The margin on the left and right side of the content when only an image is present
	ImageOnlyVMargin      float64         // The margin on the top and bottom of the content when only an image is present
	ClickAnimationTime    time.Duration   // The amount of time to spend animating the click action
	ImageSide             side.Side       // The side of the text the image should be on
	HAlign                align.Alignment // The horizontal alignment to use
	VAlign                align.Alignment // The vertical alignment to use
	Pressed               bool
	Sticky                bool
}

// NewWithText creates a new button with the specified text.
func NewWithText(text string) *Button {
	return New(text, nil)
}

// NewWithImage creates a new button with the specified image.
func NewWithImage(image *draw.Image) *Button {
	return New("", image)
}

// New creates a new button. Both image and text are optional.
func New(text string, image *draw.Image) *Button {
	b := &Button{
		Image:                 image,
		Text:                  text,
		Font:                  draw.SystemFont,
		BackgroundInk:         draw.ControlBackgroundInk,
		SelectedBackgroundInk: draw.ControlSelectedBackgroundInk,
		FocusedBackgroundInk:  draw.ControlFocusedBackgroundInk,
		PressedBackgroundInk:  draw.ControlPressedBackgroundInk,
		EdgeInk:               draw.ControlEdgeAdjColor,
		TextInk:               draw.ControlTextColor,
		PressedTextInk:        draw.AlternateSelectedControlTextColor,
		ImageGap:              3,
		CornerRadius:          4,
		HMargin:               8,
		VMargin:               1,
		ImageOnlyHMargin:      3,
		ImageOnlyVMargin:      3,
		ClickAnimationTime:    time.Millisecond * 100,
		ImageSide:             side.Left,
		HAlign:                align.Middle,
		VAlign:                align.Middle,
	}
	b.InitTypeAndID(b)
	b.SetFocusable(true)
	b.SetSizer(b.DefaultSizes)
	b.DrawCallback = b.DefaultDraw
	b.GainedFocusCallback = b.MarkForRedraw
	b.LostFocusCallback = b.MarkForRedraw
	b.MouseDownCallback = b.DefaultMouseDown
	b.MouseDragCallback = b.DefaultMouseDrag
	b.MouseUpCallback = b.DefaultMouseUp
	b.KeyDownCallback = b.DefaultKeyDown
	return b
}

// DefaultSizes provides the default sizing.
func (b *Button) DefaultSizes(hint geom.Size) (min, pref, max geom.Size) {
	text := b.Text
	if b.Image == nil && text == "" {
		text = "M"
	}
	pref = label.Size(text, b.Font, b.Image, b.ImageSide, b.ImageGap)
	if border := b.Border(); border != nil {
		pref.AddInsets(border.Insets())
	}
	pref.Width += b.horizontalMargin()*2 + 2
	pref.Height += b.verticalMargin()*2 + 2
	pref.GrowToInteger()
	pref.ConstrainForHint(hint)
	return pref, pref, layout.MaxSize(pref)
}

func (b *Button) horizontalMargin() float64 {
	if b.Text == "" && b.Image != nil {
		return b.ImageOnlyHMargin
	}
	return b.HMargin
}

func (b *Button) verticalMargin() float64 {
	if b.Text == "" && b.Image != nil {
		return b.ImageOnlyVMargin
	}
	return b.VMargin
}

// DefaultDraw provides the default drawing.
func (b *Button) DefaultDraw(gc draw.Context, dirty geom.Rect, inLiveResize bool) {
	if !b.Enabled() {
		gc.SetOpacity(0.33)
	}
	rect := b.ContentRect(false)
	widget.DrawRoundedRectBase(gc, rect, b.CornerRadius, b.currentBackgroundInk(), b.EdgeInk)
	rect.InsetUniform(1.5)
	rect.X += b.horizontalMargin()
	rect.Y += b.verticalMargin()
	rect.Width -= b.horizontalMargin() * 2
	rect.Height -= b.verticalMargin() * 2
	label.Draw(gc, rect, b.HAlign, b.VAlign, b.Text, b.Font, b.currentTextInk(), b.Image, b.ImageSide, b.ImageGap, b.Enabled())
}

func (b *Button) currentBackgroundInk() draw.Ink {
	switch {
	case b.Pressed:
		return b.PressedBackgroundInk
	case b.Focused():
		return b.FocusedBackgroundInk
	case b.Sticky && b.Selected():
		return b.SelectedBackgroundInk
	default:
		return b.BackgroundInk
	}
}

func (b *Button) currentTextInk() draw.Ink {
	if b.Pressed || b.Focused() {
		return b.PressedTextInk
	}
	return b.TextInk
}

// Click makes the button behave as if a user clicked on it.
func (b *Button) Click() {
	b.SetSelected(true)
	pressed := b.Pressed
	b.Pressed = true
	b.MarkForRedraw()
	b.FlushDrawing()
	b.Pressed = pressed
	time.Sleep(b.ClickAnimationTime)
	b.MarkForRedraw()
	if b.ClickCallback != nil {
		b.ClickCallback()
	}
}

// DefaultMouseDown provides the default mouse down handling.
func (b *Button) DefaultMouseDown(where geom.Point, button, clickCount int, mod keys.Modifiers) bool {
	b.Pressed = true
	b.MarkForRedraw()
	return true
}

// DefaultMouseDrag provides the default mouse drag handling.
func (b *Button) DefaultMouseDrag(where geom.Point, button int, mod keys.Modifiers) {
	rect := b.ContentRect(false)
	pressed := rect.ContainsPoint(where)
	if b.Pressed != pressed {
		b.Pressed = pressed
		b.MarkForRedraw()
	}
}

// DefaultMouseUp provides the default mouse up handling.
func (b *Button) DefaultMouseUp(where geom.Point, button int, mod keys.Modifiers) {
	b.Pressed = false
	b.MarkForRedraw()
	rect := b.ContentRect(false)
	if rect.ContainsPoint(where) {
		b.SetSelected(true)
		if b.ClickCallback != nil {
			b.ClickCallback()
		}
	}
}

// DefaultKeyDown provides the default key down handling.
func (b *Button) DefaultKeyDown(keyCode int, ch rune, mod keys.Modifiers, repeat bool) bool {
	if keys.IsControlAction(keyCode) {
		b.Click()
		return true
	}
	return false
}
